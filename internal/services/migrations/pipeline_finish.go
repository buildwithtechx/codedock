package migrations

import (
	"context"
	"fmt"
	"time"

	"codedock/internal/engine/dockerprobe"
	"codedock/internal/models"
)

func (e *pipelineExecutor) ask(ctx context.Context, prompt *models.MigrationPrompt) (promptAnswer, error) {
	e.service.appendLog(ctx, e.run.ID, "PROMPT: "+prompt.Detail)
	if err := e.service.migrations.SetRunPrompt(ctx, e.run.ID, encodePrompt(prompt)); err != nil {
		return promptAnswer{}, err
	}
	expiry := time.Until(time.Unix(prompt.ExpiresAt, 0))
	if expiry <= 0 {
		expiry = time.Minute
	}
	timer := time.NewTimer(expiry)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		e.cancelled = true
		return promptAnswer{}, ctx.Err()
	case <-timer.C:
		return promptAnswer{}, fmt.Errorf("prompt expired; resume the migration to continue")
	case answer := <-e.handle.answer:
		if answer.promptID != prompt.ID {
			return promptAnswer{}, fmt.Errorf("stale prompt answer")
		}
		if err := e.service.migrations.SetRunPrompt(ctx, e.run.ID, ""); err != nil {
			return promptAnswer{}, err
		}
		if err := e.saveProgress(ctx); err != nil {
			return promptAnswer{}, err
		}
		return answer, nil
	}
}

func (e *pipelineExecutor) deployPhase(ctx context.Context) error {
	e.phase = "DEPLOY"
	e.service.setStatus(ctx, e.run.ID, models.MigrationStatusTransferring, e.phase, "")
	if err := e.service.projects.SetServer(ctx, e.run.ProjectID, e.run.TargetServerID); err != nil {
		return fmt.Errorf("bind project to target: %w", err)
	}
	for index := range e.progress {
		entry := &e.progress[index]
		if entry.Skipped || entry.Deployed {
			continue
		}
		if err := e.checkCancel(ctx); err != nil {
			return err
		}
		if err := e.checkPort(ctx, entry); err != nil {
			return err
		}
		e.service.appendLog(ctx, e.run.ID, "DEPLOYING "+entry.Service)
		containerID, err := e.service.deployer.DeployAppService(ctx, entry.ServiceID, "", e.service.logWriter(e.run.ID))
		if err != nil {
			entry.Error = err.Error()
			_ = e.saveProgress(ctx)
			return fmt.Errorf("deploy %s: %w", entry.Service, err)
		}
		entry.Deployed = true
		entry.Error = ""
		e.service.appendLog(ctx, e.run.ID, "DEPLOYED "+entry.Service+" as "+containerID)
		if err := e.saveProgress(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (e *pipelineExecutor) checkPort(ctx context.Context, entry *models.MigrationServiceProgress) error {
	service, err := e.service.appRepo.GetByID(ctx, entry.ServiceID)
	if err != nil {
		return err
	}
	if service.InternalPort <= 0 {
		return nil
	}
	used, err := usedHostPorts(ctx, e.target)
	if err != nil {
		return err
	}
	port := fmt.Sprintf("%d", service.InternalPort)
	if !used[port] {
		return nil
	}
	if decision, ok := e.selection.Decisions["port:"+port]; ok {
		if decision == "skip-service" {
			entry.Skipped = true
			return nil
		}
		return nil
	}
	answer, err := e.ask(ctx, newPrompt(models.MigrationPromptPortConflict, port,
		fmt.Sprintf("host port %s on the target is already published; %s wants it", port, entry.Service),
		[]models.MigrationPromptOption{
			{ID: "deploy-anyway", Label: "Deploy anyway", Description: "Attempt deployment despite the conflict"},
			{ID: "skip-service", Label: "Skip service", Description: "Exclude this service from the move"},
		}))
	if err != nil {
		return err
	}
	e.selection.Decisions["port:"+port] = answer.optionID
	if answer.optionID == "skip-service" {
		entry.Skipped = true
	}
	return nil
}

func (e *pipelineExecutor) verifyPhase(ctx context.Context) error {
	e.phase = "VERIFY"
	e.service.setStatus(ctx, e.run.ID, models.MigrationStatusVerifying, e.phase, "")
	for index := range e.progress {
		entry := &e.progress[index]
		if entry.Skipped || entry.Verified {
			continue
		}
		if err := e.checkCancel(ctx); err != nil {
			return err
		}
		service, err := e.service.appRepo.GetByID(ctx, entry.ServiceID)
		if err != nil {
			return err
		}
		if service.ContainerID == "" {
			return fmt.Errorf("verify %s: deployment recorded no container", entry.Service)
		}
		state, err := dockerprobe.ContainerState(ctx, e.target, service.ContainerID)
		if err != nil {
			entry.Error = err.Error()
			_ = e.saveProgress(ctx)
			return fmt.Errorf("verify %s: %w", entry.Service, err)
		}
		if state != "running" {
			return fmt.Errorf("verify %s: target container is %s", entry.Service, state)
		}
		entry.Verified = true
		entry.Error = ""
		e.service.appendLog(ctx, e.run.ID, "VERIFIED "+entry.Service)
		if err := e.saveProgress(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (e *pipelineExecutor) finishPhase(ctx context.Context) error {
	if e.run.Mode == models.MigrationModeCopy {
		e.phase = "DONE"
		e.service.setStatus(ctx, e.run.ID, models.MigrationStatusCompleted, e.phase, "")
		e.service.appendLog(ctx, e.run.ID, "COMPLETED: copy finished; originals untouched")
		return nil
	}
	if e.selection.KillOriginals {
		if err := e.service.cutoverKill(context.Background(), e.run, true); err != nil {
			return fmt.Errorf("automatic cutover: %w", err)
		}
		e.service.setStatus(ctx, e.run.ID, models.MigrationStatusCompleted, "CUTOVER", "")
		e.service.appendLog(ctx, e.run.ID, "COMPLETED: originals retired after verification")
		return nil
	}
	e.phase = "AWAIT"
	e.service.setStatus(ctx, e.run.ID, models.MigrationStatusAwaitingCutover, e.phase, "")
	e.service.appendLog(ctx, e.run.ID, "AWAITING CUTOVER: review target health, then confirm")
	return nil
}

func (e *pipelineExecutor) rollback(ctx context.Context) {
	for _, entry := range e.progress {
		if !entry.Deployed || entry.ServiceID == "" {
			continue
		}
		service, err := e.service.appRepo.GetByID(ctx, entry.ServiceID)
		if err != nil || service == nil {
			continue
		}
		_ = e.service.deployer.RemoveAppService(ctx, service)
		for _, volume := range entry.CopiedVolumes {
			_ = dockerprobe.RemoveVolume(ctx, e.target, volume)
		}
	}
}

func (s *Service) ResumeRun(ctx context.Context, userID, runID string, req models.ResumeMigrationRequest) (*StartMigrationResult, error) {
	run, err := s.orgRun(ctx, userID, runID)
	if err != nil {
		return nil, err
	}
	if run.Status == models.MigrationStatusCompleted {
		return nil, fmt.Errorf("migration already completed")
	}
	if s.runs.get(run.ID) != nil {
		return nil, fmt.Errorf("migration worker is still running")
	}
	selection := decodeSelection(run.Selection)
	for key, value := range req.Overrides {
		selection.Overrides[key] = value
	}
	seen := map[string]bool{}
	for _, skip := range selection.Skips {
		seen[skip] = true
	}
	for _, skip := range req.Skips {
		if !seen[skip] {
			selection.Skips = append(selection.Skips, skip)
			seen[skip] = true
		}
	}
	token, hash, err := newConfirmationToken()
	if err != nil {
		return nil, err
	}
	run.Selection = encodeSelection(selection)
	run.Prompt = ""
	run.Error = ""
	run.TokenHash = hash
	run.CancelRequested = false
	if run.Status == models.MigrationStatusAwaitingCutover {
		run.Status = models.MigrationStatusAwaitingCutover
		if err := s.migrations.UpdateRun(ctx, run); err != nil {
			return nil, err
		}
		return &StartMigrationResult{RunID: run.ID, Confirmation: token}, nil
	}
	run.Status = models.MigrationStatusTransferring
	if err := s.migrations.UpdateRun(ctx, run); err != nil {
		return nil, err
	}
	s.launch(run.ID)
	return &StartMigrationResult{RunID: run.ID, Confirmation: token}, nil
}

func (s *Service) CleanupTargetData(ctx context.Context, userID, runID string) error {
	run, err := s.orgRun(ctx, userID, runID)
	if err != nil {
		return err
	}
	if err := s.requireOrg(ctx, userID, run.OrganizationID, true); err != nil {
		return err
	}
	if run.Status != models.MigrationStatusFailed && run.Status != models.MigrationStatusCancelled {
		return fmt.Errorf("only failed or cancelled runs can clean the target")
	}
	_, targetRunner, err := s.targetServer(ctx, run.OrganizationID, run.TargetServerID)
	if err != nil {
		return err
	}
	progress := decodeProgress(run.Progress)
	for index := range progress {
		entry := &progress[index]
		for _, volume := range entry.CopiedVolumes {
			if err := dockerprobe.RemoveVolume(ctx, targetRunner, volume); err != nil {
				return err
			}
		}
		entry.CopiedVolumes = nil
		entry.Transferred = false
		entry.Deployed = false
		entry.Verified = false
		if entry.ServiceID == "" {
			continue
		}
		service, err := s.appRepo.GetByID(ctx, entry.ServiceID)
		if err != nil || service == nil || service.ContainerID == "" {
			continue
		}
		if _, err := targetRunner.Run(ctx, "docker stop "+service.ContainerID); err != nil {
			continue
		}
		_, _ = targetRunner.Run(ctx, "docker rm "+service.ContainerID)
	}
	run.Progress = encodeProgress(progress)
	if err := s.migrations.UpdateRun(ctx, run); err != nil {
		return err
	}
	s.appendLog(ctx, run.ID, "CLEANED target volumes and containers from failed run")
	return nil
}
