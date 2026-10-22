package deployments

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"codedock/internal/engine/deploy"
	"codedock/internal/models"
)

func (s *DeploymentService) ExecuteDeploymentAsync(d *models.Deployment) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	s.operationMu.Lock()
	if s.operations == nil {
		s.operations = make(map[string]context.CancelFunc)
	}
	if _, busy := s.operations["service:"+d.ServiceID]; busy {
		s.operationMu.Unlock()
		cancel()
		s.finishDeployment(d.ID, models.DeploymentStatusFailed, "Another deployment is already running for this service.\n", "", "")
		return
	}
	if s.admissionGuard != nil {
		if err := s.admissionGuard.Acquire(d.ServiceID, ""); err != nil {
			s.operationMu.Unlock()
			cancel()
			s.finishDeployment(d.ID, models.DeploymentStatusFailed, fmt.Sprintf("Build admission limit exceeded: %v\n", err), "", "")
			return
		}
	}
	s.operations[d.ID] = cancel
	s.operations["service:"+d.ServiceID] = cancel
	s.operationMu.Unlock()
	go func() {
		defer cancel()
		defer func() {
			if s.admissionGuard != nil {
				s.admissionGuard.Release(d.ServiceID)
			}
			s.operationMu.Lock()
			delete(s.operations, d.ID)
			delete(s.operations, "service:"+d.ServiceID)
			s.operationMu.Unlock()
		}()
		logWriter := &deploymentLogWriter{service: s, id: d.ID, ctx: ctx, phase: "PREPARING"}
		execCtx := deploy.WithDeploymentProgress(ctx, logWriter.phaseChanged)
		execCtx = deploy.WithImageBuiltCallback(execCtx, func(tag string) {
			d.ImageRef = tag
			if err := s.repo.Update(context.Background(), d); err != nil {
				slog.Warn("failed to update deployment image ref", "deployment", d.ID, "error", err)
			}
		})
		id, err := s.executeDeployment(execCtx, d, logWriter)
		status := models.DeploymentStatusReady
		resourceIssue := ""
		if err != nil {
			status = models.DeploymentStatusFailed
			if ctx.Err() != nil {
				status = models.DeploymentStatus("CANCELLED")
			}
			logWriter.logs.WriteString(fmt.Sprintf("Deployment stopped: %v\n", err))
			var rErr *deploy.ResourceFailureError
			if errors.As(err, &rErr) {
				resourceIssue = rErr.Issue
			} else {
				resourceIssue = deploy.ClassifyResourceIssue(err.Error() + " " + logWriter.logs.String())
			}
		} else {
			logWriter.logs.WriteString("Deployment ready.\n")
		}
		s.finishDeployment(d.ID, status, logWriter.logs.String(), id, resourceIssue)
	}()
}

func (s *DeploymentService) CancelDeployment(id string) error {
	s.operationMu.Lock()
	defer s.operationMu.Unlock()
	cancel, ok := s.operations[id]
	if !ok {
		return fmt.Errorf("deployment is no longer running")
	}
	cancel()
	return nil
}

func (s *DeploymentService) finishDeployment(id string, status models.DeploymentStatus, logs, containerID, resourceIssue string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.UpdateStatus(ctx, DeployStatusOpts{
		ID:            id,
		Status:        status,
		BuildLogs:     logs,
		ContainerID:   containerID,
		ResourceIssue: resourceIssue,
	}); err != nil {
		slog.Error("save deployment result", "deployment", id, "error", err)
	}
	if status == models.DeploymentStatusReady {
		now := time.Now().UTC()
		_ = s.repo.SetArtifactRetainedAt(ctx, id, &now)
	}
}

func (s *DeploymentService) executeDeployment(ctx context.Context, d *models.Deployment, writer *deploymentLogWriter) (string, error) {
	if s.deployer == nil || s.appRepo == nil {
		return "", fmt.Errorf("deployment runtime unavailable")
	}
	operation, release, err := s.deployer.BeginServiceOperation(ctx, d.ServiceID)
	if err != nil {
		return "", err
	}
	defer release()
	ctx = operation
	if err := writer.phaseChanged("PREPARING"); err != nil {
		return "", err
	}
	app, err := s.appRepo.GetByID(ctx, d.ServiceID)
	if err != nil {
		return "", fmt.Errorf("load service: %w", err)
	}

	if s.volumeRepo != nil {
		volumes, err := s.volumeRepo.ListByService(ctx, app.ID)
		if err != nil {
			return "", fmt.Errorf("load volumes: %w", err)
		}
		app.Volumes = volumes
	}
	if err := deploy.AssertRequiredServiceEnvironment(ctx, s.deployer, app); err != nil {
		return "", err
	}
	if s.BeforeProjectDeployment != nil {
		if err := s.BeforeProjectDeployment(ctx, app.ProjectID, app.ID); err != nil {
			return "", fmt.Errorf("pre-deployment backup: %w", err)
		}
	} else if s.BeforeDeployment != nil {
		if err := s.BeforeDeployment(ctx, app.ID); err != nil {
			return "", fmt.Errorf("pre-deployment backup: %w", err)
		}
	}
	if dependencies, ok := s.appRepo.(interface {
		DeploymentDependencies(context.Context, *models.AppService) ([]string, error)
	}); ok {
		sources, err := dependencies.DeploymentDependencies(ctx, app)
		if err != nil {
			return "", fmt.Errorf("load prerequisites: %w", err)
		}
		for _, source := range sources {
			for {
				err := s.dependencyReady(ctx, source)
				if err == nil {
					break
				}
				if _, writeErr := writer.Write([]byte("Waiting for " + source + " to become ready.\n")); writeErr != nil {
					return "", writeErr
				}
				select {
				case <-ctx.Done():
					return "", fmt.Errorf("prerequisite readiness: %w", err)
				case <-time.After(2 * time.Second):
				}
			}
		}
	}
	isRetryWithArtifact := strings.EqualFold(d.Trigger, "retry") && d.ImageRef != ""
	if isRetryWithArtifact {
		if _, writeErr := writer.Write([]byte(fmt.Sprintf("⚡ Reusing verified build artifact: %s\nSkipping source checkout and Docker build.\n", d.ImageRef))); writeErr != nil {
			return "", writeErr
		}
		appCopy := *app
		appCopy.ImageRef = d.ImageRef
		return s.deployTarget(ctx, &appCopy, "", writer)
	}

	sourceDir := ""
	if app.ImageRef == "" {
		if s.gitService == nil {
			return "", fmt.Errorf("repository integration unavailable")
		}
		sourceDir = filepath.Join("data", "builds", app.ID, d.ID)
		defer func() {
			if err := os.RemoveAll(sourceDir); err != nil {
				slog.Warn("remove deployment checkout", "error", err)
			}
		}()
		if err := writer.phaseChanged("CLONING"); err != nil {
			return "", err
		}
		if err := s.gitService.CloneOrPullAppRepositoryAtCommit(ctx, app, d.CommitHash, sourceDir, writer); err != nil {
			return "", err
		}
		if err := writer.phaseChanged("BUILDING"); err != nil {
			return "", err
		}
	} else if err := writer.phaseChanged("PULLING"); err != nil {
		return "", err
	}
	return s.deployTarget(ctx, app, sourceDir, writer)
}

type deploymentLogWriter struct {
	mu        sync.Mutex
	service   *DeploymentService
	id        string
	ctx       context.Context
	phase     string
	logs      strings.Builder
	lastFlush time.Time
}

func (w *deploymentLogWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.logs.Len() < 1024*1024 {
		w.logs.Write(data)
	}
	if time.Since(w.lastFlush) > time.Second {
		w.lastFlush = time.Now()
		if err := w.service.UpdateStatus(w.ctx, DeployStatusOpts{ID: w.id, Status: models.DeploymentStatus(w.phase), BuildLogs: w.logs.String()}); err != nil {
			return 0, err
		}
	}
	return len(data), nil
}

func (w *deploymentLogWriter) phaseChanged(phase string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.phase = phase
	w.logs.WriteString("Phase: ")
	w.logs.WriteString(phase)
	w.logs.WriteByte('\n')
	return w.service.UpdateStatus(w.ctx, DeployStatusOpts{ID: w.id, Status: models.DeploymentStatus(phase), BuildLogs: w.logs.String()})
}
