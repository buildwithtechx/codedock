package migrations

import (
	"context"
	"fmt"

	"codedock/internal/models"
)

func (s *Service) CutoverRun(ctx context.Context, userID, runID string, req models.CutoverMigrationRequest) error {
	run, err := s.orgRun(ctx, userID, runID)
	if err != nil {
		return err
	}
	if err := s.requireOrg(ctx, userID, run.OrganizationID, true); err != nil {
		return err
	}
	if run.Status != models.MigrationStatusAwaitingCutover {
		return fmt.Errorf("migration is not awaiting cutover")
	}
	if !verifyConfirmation(run.TokenHash, req.Confirmation) {
		return fmt.Errorf("invalid cutover confirmation")
	}
	selection := decodeSelection(run.Selection)
	want := "false"
	if req.Kill {
		want = "true"
	}
	if stored, ok := selection.Decisions["cutover-kill"]; ok && stored != want {
		return fmt.Errorf("failed cutover can only resume the same choice")
	}
	selection.Decisions["cutover-kill"] = want
	run.Selection = encodeSelection(selection)
	if err := s.migrations.UpdateRun(ctx, run); err != nil {
		return err
	}
	if err := s.cutoverKill(ctx, run, req.Kill); err != nil {
		s.setStatus(ctx, run.ID, models.MigrationStatusFailed, "CUTOVER", err.Error())
		s.appendLog(ctx, run.ID, "CUTOVER FAILED: "+err.Error())
		return err
	}
	s.setStatus(ctx, run.ID, models.MigrationStatusCompleted, "CUTOVER", "")
	if req.Kill {
		s.appendLog(ctx, run.ID, "COMPLETED: cutover destroyed originals; source volumes remain")
	} else {
		s.appendLog(ctx, run.ID, "COMPLETED: cutover retained originals; source volumes remain")
	}
	return nil
}

func (s *Service) cutoverKill(ctx context.Context, run *models.MigrationRun, kill bool) error {
	selection := decodeSelection(run.Selection)
	progress := decodeProgress(run.Progress)
	if run.SourceKind == "project" {
		return s.cutoverProject(ctx, run, selection, progress, kill)
	}
	source, err := s.migrations.GetSource(ctx, run.SourceID)
	if err != nil {
		return err
	}
	if source == nil {
		return fmt.Errorf("migration source not found")
	}
	runner := s.sources(source)
	for _, entry := range progress {
		if entry.Skipped || entry.ContainerID == "" {
			continue
		}
		if !entry.Verified {
			continue
		}
		if kill {
			if _, err := runner.Run(ctx, "docker stop "+entry.ContainerID); err != nil {
				return fmt.Errorf("stop original %s: %w", entry.Service, err)
			}
			if _, err := runner.Run(ctx, "docker rm "+entry.ContainerID); err != nil {
				return fmt.Errorf("remove original %s: %w", entry.Service, err)
			}
			continue
		}
		_, _ = runner.Run(ctx, "docker start "+entry.ContainerID)
	}
	return nil
}

func (s *Service) cutoverProject(ctx context.Context, run *models.MigrationRun, selection models.MigrationSelection, progress []models.MigrationServiceProgress, kill bool) error {
	_, sourceRunner, err := s.targetServer(ctx, run.OrganizationID, selection.SourceServerID)
	if err != nil {
		return err
	}
	for _, entry := range progress {
		if entry.Skipped || entry.ContainerID == "" {
			continue
		}
		if !entry.Verified {
			continue
		}
		if _, err := sourceRunner.Run(ctx, "docker stop "+entry.ContainerID); err != nil {
			return fmt.Errorf("stop original %s: %w", entry.Service, err)
		}
		if kill {
			if _, err := sourceRunner.Run(ctx, "docker rm "+entry.ContainerID); err != nil {
				return fmt.Errorf("remove original %s: %w", entry.Service, err)
			}
		}
	}
	return nil
}
