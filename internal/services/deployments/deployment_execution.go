package deployments

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"codedock.run/codedock/internal/engine/deploy"
	"codedock.run/codedock/internal/models"
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
		s.finishDeployment(d.ID, models.DeploymentStatusFailed, "Another deployment is already running for this service.\n", "")
		return
	}
	s.operations[d.ID] = cancel
	s.operations["service:"+d.ServiceID] = cancel
	s.operationMu.Unlock()
	go func() {
		defer cancel()
		defer func() {
			s.operationMu.Lock()
			delete(s.operations, d.ID)
			delete(s.operations, "service:"+d.ServiceID)
			s.operationMu.Unlock()
		}()
		logWriter := &deploymentLogWriter{service: s, id: d.ID, ctx: ctx, phase: "PREPARING"}
		id, err := s.executeDeployment(deploy.WithDeploymentProgress(ctx, logWriter.phaseChanged), d, logWriter)
		status := models.DeploymentStatusReady
		if err != nil {
			status = models.DeploymentStatusFailed
			if ctx.Err() != nil {
				status = models.DeploymentStatus("CANCELLED")
			}
			logWriter.logs.WriteString(fmt.Sprintf("Deployment stopped: %v\n", err))
		} else {
			logWriter.logs.WriteString("Deployment ready.\n")
		}
		s.finishDeployment(d.ID, status, logWriter.logs.String(), id)
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

func (s *DeploymentService) finishDeployment(id string, status models.DeploymentStatus, logs, containerID string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.UpdateStatus(ctx, DeployStatusOpts{ID: id, Status: status, BuildLogs: logs, ContainerID: containerID}); err != nil {
		slog.Error("save deployment result", "deployment", id, "error", err)
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
	if s.projectRepo != nil {
		project, err := s.projectRepo.Get(ctx, app.ProjectID)
		if err != nil {
			return "", err
		}
		if project.ServerID != "" {
			return "", fmt.Errorf("this deployment pipeline supports local Docker targets; SSH deployment must use its worker")
		}
	}
	if s.volumeRepo != nil {
		volumes, err := s.volumeRepo.ListByService(ctx, app.ID)
		if err != nil {
			return "", fmt.Errorf("load volumes: %w", err)
		}
		app.Volumes = volumes
	}
	if s.BeforeDeployment != nil {
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
				err := s.deployer.DependencyReady(ctx, source)
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
		if err := s.gitService.CloneOrPullAppRepository(ctx, app, sourceDir, writer); err != nil {
			return "", err
		}
		if err := writer.phaseChanged("BUILDING"); err != nil {
			return "", err
		}
	} else if err := writer.phaseChanged("PULLING"); err != nil {
		return "", err
	}
	return s.deployer.DeployAppService(ctx, app, sourceDir, writer)
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
