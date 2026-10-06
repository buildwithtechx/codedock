package backups

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
)

func (s *BackupService) StartScheduledRun(ctx context.Context, id string) error {
	cfg, err := s.backupRepo.GetConfigByID(ctx, id)
	if err != nil {
		return err
	}
	if cfg.OwnerID == "" {
		return fmt.Errorf("save the backup policy to assign a scheduled-run owner")
	}
	if !cfg.BackupEnabled || cfg.Status != models.BackupConfigStatusActive || cfg.Schedule == "manual" {
		return fmt.Errorf("scheduled policy is inactive")
	}
	if s.RunAuthorization == nil {
		return fmt.Errorf("scheduled backup authorization unavailable")
	}
	if err := s.RunAuthorization(ctx, cfg.OwnerID, cfg.ProjectID); err != nil {
		return err
	}
	_, err = s.StartRun(ctx, cfg.OwnerID, cfg.ProjectID, id)
	return err
}
