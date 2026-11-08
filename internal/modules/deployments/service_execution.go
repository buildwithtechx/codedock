package deployments

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"codedock/internal/utils"
)

type deploymentLogWriter struct {
	service *DeploymentService
	id      string
	ctx     context.Context
	mu      sync.Mutex
	logs    bytes.Buffer
}

func (w *deploymentLogWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	n, err = w.logs.Write(p)
	return n, err
}

func (s *DeploymentService) ExecuteDeploymentAsync(d *Deployment) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	s.operationMu.Lock()
	if s.operations == nil {
		s.operations = make(map[string]context.CancelFunc)
	}
	if _, busy := s.operations["service:"+d.ServiceID]; busy {
		s.operationMu.Unlock()
		cancel()
		s.finishDeployment(d.ID, d.ProjectID, DeploymentStatusFailed, "Another deployment is already running for this service.\n", "", "")
		return
	}
	if s.admissionGuard != nil {
		if err := s.admissionGuard.Acquire(d.ServiceID, ""); err != nil {
			s.operationMu.Unlock()
			cancel()
			s.finishDeployment(d.ID, d.ProjectID, DeploymentStatusFailed, fmt.Sprintf("Build admission limit exceeded: %v\n", err), "", "")
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

		logWriter := &deploymentLogWriter{service: s, id: d.ID, ctx: ctx}
		containerID, err := s.executeDeployment(ctx, d, logWriter)
		status := DeploymentStatusReady
		resourceIssue := ""
		if err != nil {
			status = DeploymentStatusFailed
			if ctx.Err() != nil {
				status = DeploymentStatusCancelled
			}
			fmt.Fprintf(&logWriter.logs, "Deployment stopped: %v\n", err)
			resourceIssue = "execution_error"
		} else {
			logWriter.logs.WriteString("Deployment ready.\n")
		}
		s.finishDeployment(d.ID, d.ProjectID, status, logWriter.logs.String(), containerID, resourceIssue)
	}()
}

func (s *DeploymentService) executeDeployment(ctx context.Context, d *Deployment, logWriter *deploymentLogWriter) (string, error) {
	if s.appLookup == nil {
		return "", ErrAppServiceRequired
	}
	app, err := s.appLookup.GetAppService(ctx, d.ServiceID)
	if err != nil {
		return "", fmt.Errorf("load app service: %w", err)
	}

	sourceDir := filepath.Join(utils.GetDataDir(), "builds", d.ID)
	_ = os.MkdirAll(sourceDir, 0o755)
	defer os.RemoveAll(sourceDir)

	if s.codeSync != nil && app.RepositoryURL != "" {
		fmt.Fprintf(logWriter, "📥 Syncing codebase from %s (branch: %s)...\n", app.RepositoryURL, d.Branch)
		if err := s.codeSync.SyncCodebase(ctx, app.RepositoryURL, d.Branch, d.CommitHash, sourceDir, logWriter); err != nil {
			return "", fmt.Errorf("sync codebase: %w", err)
		}
	}

	if s.deployEngine != nil {
		fmt.Fprintf(logWriter, "🚀 Launching build & deploy for %s...\n", app.Name)
		containerID, err := s.deployEngine.DeployApp(ctx, app, sourceDir, logWriter)
		if err != nil {
			return "", err
		}
		if s.appLookup != nil && containerID != "" {
			_ = s.appLookup.UpdateAppContainerID(ctx, app.ID, containerID)
		}
		return containerID, nil
	}

	return "", nil
}

func (s *DeploymentService) CancelDeployment(id string) error {
	s.operationMu.Lock()
	defer s.operationMu.Unlock()
	cancel, ok := s.operations[id]
	if !ok {
		return ErrDeploymentNotRunning
	}
	cancel()
	return nil
}

func (s *DeploymentService) finishDeployment(id, projectID string, status DeploymentStatus, logs, containerID, resourceIssue string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = s.UpdateStatus(ctx, DeployStatusOpts{
		ID:            id,
		Status:        status,
		BuildLogs:     logs,
		ContainerID:   containerID,
		ResourceIssue: resourceIssue,
	})
}
