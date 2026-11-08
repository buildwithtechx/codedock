package backups

import (
	"context"
	"fmt"
	"log/slog"
)

func (s *Service) CreateBatch(ctx context.Context, batch *BackupPolicyBatch) error {
	if s.batches == nil {
		return ErrBatchStorageUnavailable
	}
	if batch.Name == "" || batch.ProjectID == "" {
		return ErrBatchNameAndProjectRequired
	}
	if batch.Schedule == "" {
		batch.Schedule = DefaultSchedule
	}
	if _, err := ParseSchedule(batch.Schedule, batch.Timezone); err != nil {
		return err
	}
	if batch.Timeout == 0 {
		batch.Timeout = DefaultTimeout
	}
	return s.batches.Create(ctx, batch)
}

func (s *Service) ListBatches(ctx context.Context, projectID string) ([]BackupPolicyBatch, error) {
	if s.batches == nil {
		return nil, ErrBatchStorageUnavailable
	}
	return s.batches.ListByProject(ctx, projectID)
}

func (s *Service) GetBatch(ctx context.Context, id string) (*BackupPolicyBatch, error) {
	if s.batches == nil {
		return nil, ErrBatchStorageUnavailable
	}
	return s.batches.Get(ctx, id)
}

func (s *Service) TriggerBatch(ctx context.Context, batchID string) (*BackupPolicyBatchSummary, error) {
	if s.batches == nil {
		return nil, ErrBatchStorageUnavailable
	}
	batch, err := s.batches.Get(ctx, batchID)
	if err != nil {
		return nil, err
	}
	members, err := s.batches.ListMembers(ctx, batchID)
	if err != nil {
		return nil, err
	}
	summary := &BackupPolicyBatchSummary{
		BatchID:     batch.ID,
		TotalCount:  len(members),
		Status:      RecordStatusCompleted,
		TriggeredAt: batch.CreatedAt,
	}
	for _, member := range members {
		if _, err := s.TriggerBackup(ctx, member.ID); err != nil {
			summary.FailedCount++
			summary.Status = "partial"
			slog.Warn("batch backup member failed", "memberId", member.ID, "error", err)
		} else {
			summary.SuccessCount++
		}
	}
	return summary, nil
}

func (s *Service) StartScheduledRun(ctx context.Context, cfg *BackupConfig) error {
	if s.RunAuthorization != nil {
		if err := s.RunAuthorization(ctx, cfg.OwnerID, cfg.ProjectID); err != nil {
			return fmt.Errorf("scheduled backup authorization failed: %w", err)
		}
	}
	_, err := s.TriggerBackup(ctx, cfg.ID)
	return err
}

func (s *Service) StartRun(ctx context.Context, user, project, id string) (any, error) {
	if s.operations == nil {
		return nil, fmt.Errorf("operations service unavailable")
	}
	return s.operations.Review(ctx, user, project, "backup_run", "backup:"+id, id, "", "Run on-demand backup for policy "+id)
}

func (s *Service) BeforeProjectDeployment(ctx context.Context, projectID, serviceID string) error {
	configs, err := s.backupRepo.ListConfigs(ctx)
	if err != nil {
		return err
	}
	for _, cfg := range configs {
		if (cfg.ProjectID == projectID || (serviceID != "" && cfg.ServiceID == serviceID)) && cfg.PreDeployment && cfg.BackupEnabled {
			if _, err := s.TriggerBackup(ctx, cfg.ID); err != nil {
				return fmt.Errorf("pre-deployment backup failed: %w", err)
			}
		}
	}
	return nil
}

func (s *Service) BeforeDeployment(ctx context.Context, serviceID string) error {
	configs, err := s.backupRepo.ListConfigs(ctx)
	if err != nil {
		return err
	}
	for _, cfg := range configs {
		if cfg.ServiceID == serviceID && cfg.PreDeployment && cfg.BackupEnabled {
			if _, err := s.TriggerBackup(ctx, cfg.ID); err != nil {
				return fmt.Errorf("pre-deployment backup failed: %w", err)
			}
		}
	}
	return nil
}
