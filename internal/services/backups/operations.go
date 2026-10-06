package backups

import (
	"codedock.run/codedock/internal/engine/backup"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services/operations"
	"context"
	"encoding/json"
	"fmt"
)

func (s *BackupService) SetOperations(service *operations.Service) {
	s.operations = service
	if s.manager != nil {
		s.manager.SetScheduledRunner(s.StartScheduledRun)
	}
}
func (s *BackupService) ProtectRecord(ctx context.Context, id string, until int64) error {
	store, ok := s.backupRepo.(interface {
		ProtectRecord(context.Context, string, int64) error
	})
	if !ok {
		return fmt.Errorf("backup protection storage unavailable")
	}
	return store.ProtectRecord(ctx, id, until)
}
func (s *BackupService) PrepareRestore(ctx context.Context, user, project, id, database string) (*models.OperationReview, error) {
	if s.manager == nil || s.operations == nil {
		return nil, fmt.Errorf("restore runtime unavailable")
	}
	target, err := s.manager.PrepareRestore(ctx, id, database)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(target)
	if err != nil {
		return nil, err
	}
	resource := target.DatabaseID
	if target.VolumeName != "" {
		resource = target.VolumeName
	}
	effects := fmt.Sprintf("Restore archive %s (SHA256 %s) into %s using %s mode. Target data may be overwritten. Interruption may leave partial data. This confirmation expires in ten minutes.", id, target.SHA256, resource, target.Mode)
	return s.operations.Review(ctx, user, project, "restore", "restore:"+resource, string(payload), target.Snapshot, effects)
}
func (s *BackupService) ApplyRestore(ctx context.Context, user, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Kind != "restore" || op.UserID != user {
		return fmt.Errorf("restore operation not found")
	}
	var target models.RestoreTarget
	if err := json.Unmarshal([]byte(op.Payload), &target); err != nil {
		return err
	}
	target.Snapshot = op.Snapshot
	current, err := s.manager.PrepareRestore(ctx, target.RecordID, target.DatabaseID)
	if err != nil {
		return err
	}
	return s.operations.Apply(ctx, id, user, confirmation, current.Snapshot, func(ctx context.Context, _ *models.Operation, progress func(string, string) error) error {
		return s.manager.ApplyRestore(ctx, &target, progress)
	})
}
func (s *BackupService) StartRun(ctx context.Context, user, project, id string) (*models.Operation, error) {
	if s.manager == nil || s.operations == nil {
		return nil, fmt.Errorf("backup runtime unavailable")
	}
	cfg, err := s.backupRepo.GetConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if !cfg.BackupEnabled {
		return nil, fmt.Errorf("backup policy is disabled")
	}
	reviewed, err := s.operations.Review(ctx, user, project, "backup", "backup-config:"+id, id, cfg.UpdatedAt, "Create and verify a backup using the saved producer and destination; preserve protected records during retention.")
	if err != nil {
		return nil, err
	}
	if err := s.operations.Apply(ctx, reviewed.Operation.ID, user, reviewed.Confirmation, cfg.UpdatedAt, func(ctx context.Context, op *models.Operation, progress func(string, string) error) error {
		latest, err := s.backupRepo.GetConfigByID(ctx, id)
		if err != nil {
			return err
		}
		if latest.UpdatedAt != op.Snapshot {
			return fmt.Errorf("backup policy changed before execution")
		}
		record, err := s.manager.TriggerBackup(backup.WithProgress(ctx, progress), id)
		if err != nil {
			return err
		}
		if _, err := s.manager.VerifyArchive(ctx, record.ID); err != nil {
			return err
		}
		return progress("VERIFIED", "Verified backup record "+record.ID)
	}); err != nil {
		return nil, err
	}
	return s.operations.Get(ctx, reviewed.Operation.ID)
}
func (s *BackupService) BeforeDeployment(ctx context.Context, serviceID string) error {
	configs, err := s.backupRepo.ListConfigs(ctx)
	if err != nil {
		return err
	}
	for _, cfg := range configs {
		if !cfg.PreDeployment || cfg.ServiceID != serviceID {
			continue
		}
		if !cfg.BackupEnabled {
			return fmt.Errorf("required pre-deployment backup policy %s is disabled", cfg.Name)
		}
		if s.manager == nil {
			return fmt.Errorf("required backup runtime unavailable")
		}
		record, err := s.manager.TriggerBackup(ctx, cfg.ID)
		if err != nil {
			return fmt.Errorf("required pre-deployment backup %s: %w", cfg.Name, err)
		}
		if _, err := s.manager.VerifyArchive(ctx, record.ID); err != nil {
			return fmt.Errorf("verify pre-deployment backup: %w", err)
		}
	}
	return nil
}
