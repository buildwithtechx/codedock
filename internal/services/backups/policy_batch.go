package backups

import (
	"codedock.run/codedock/internal/engine/backup"
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type BatchStore interface {
	Create(ctx context.Context, batch *models.BackupPolicyBatch) error
	Get(ctx context.Context, id string) (*models.BackupPolicyBatch, error)
	ListByProject(ctx context.Context, projectID string) ([]models.BackupPolicyBatch, error)
	ListMembers(ctx context.Context, batchID string) ([]*models.BackupConfig, error)
}

type SFTPStore interface {
	Create(ctx context.Context, dest *models.SFTPDestination) error
	Get(ctx context.Context, id string) (*models.SFTPDestination, error)
	ListByProject(ctx context.Context, projectID string) ([]models.SFTPDestination, error)
	Delete(ctx context.Context, id string) error
}

func (s *BackupService) SetBatches(store BatchStore) { s.batches = store }
func (s *BackupService) SetSFTP(store SFTPStore)     { s.sftp = store }

func (s *BackupService) CreateBatch(ctx context.Context, batch *models.BackupPolicyBatch) error {
	if s.batches == nil {
		return fmt.Errorf("policy batch storage unavailable")
	}
	if batch.Name == "" || batch.ProjectID == "" {
		return fmt.Errorf("batch name and project are required")
	}
	if batch.Schedule == "" {
		batch.Schedule = "0 2 * * *"
	}
	if _, err := backup.ParseSchedule(batch.Schedule, batch.Timezone); err != nil {
		return err
	}
	if batch.Timeout == 0 {
		batch.Timeout = 3600
	}
	return s.batches.Create(ctx, batch)
}

func (s *BackupService) ListBatches(ctx context.Context, projectID string) ([]models.BackupPolicyBatch, error) {
	if s.batches == nil {
		return nil, fmt.Errorf("policy batch storage unavailable")
	}
	return s.batches.ListByProject(ctx, projectID)
}

func (s *BackupService) GetBatch(ctx context.Context, batchID string) (*models.BackupPolicyBatch, error) {
	if s.batches == nil {
		return nil, fmt.Errorf("policy batch storage unavailable")
	}
	return s.batches.Get(ctx, batchID)
}

func (s *BackupService) TriggerBatch(ctx context.Context, batchID string) ([]*models.BackupRecord, error) {
	if s.batches == nil || s.manager == nil {
		return nil, fmt.Errorf("policy batch runtime unavailable")
	}
	members, err := s.batches.ListMembers(ctx, batchID)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, fmt.Errorf("policy batch has no member policies")
	}
	records := make([]*models.BackupRecord, 0, len(members))
	for _, member := range members {
		if !member.BackupEnabled {
			return nil, fmt.Errorf("member policy %s is disabled", member.Name)
		}
		record, err := s.manager.TriggerBackup(ctx, member.ID)
		if err != nil {
			return nil, fmt.Errorf("batch member %s: %w", member.Name, err)
		}
		if _, err := s.manager.VerifyArchive(ctx, record.ID); err != nil {
			return nil, fmt.Errorf("verify batch member %s: %w", member.Name, err)
		}
		records = append(records, record)
	}
	return records, nil
}

func (s *BackupService) expandInheritedPolicies(ctx context.Context, configs []*models.BackupConfig) ([]*models.BackupConfig, error) {
	if s.batches == nil {
		return configs, nil
	}
	expanded := make([]*models.BackupConfig, 0, len(configs))
	for _, cfg := range configs {
		if cfg.ParentBatchID == "" || hasProducer(cfg) {
			expanded = append(expanded, cfg)
			continue
		}
		members, err := s.batches.ListMembers(ctx, cfg.ParentBatchID)
		if err != nil {
			return nil, err
		}
		expanded = append(expanded, members...)
	}
	return expanded, nil
}

func hasProducer(cfg *models.BackupConfig) bool {
	return cfg.DatabaseID != "" || cfg.VolumeName != "" || cfg.FileSourcePath != "" || cfg.CustomBackupCommand != ""
}

func (s *BackupService) CreateSFTPDestination(ctx context.Context, dest *models.SFTPDestination) error {
	if s.sftp == nil {
		return fmt.Errorf("sftp storage unavailable")
	}
	if dest.ID == "" {
		dest.ID = uuid.New().String()
	}
	if dest.Host == "" || dest.Username == "" {
		return fmt.Errorf("sftp host and username are required")
	}
	if dest.Password == "" && dest.PrivateKey == "" {
		return fmt.Errorf("sftp destination needs a password or private key")
	}
	return s.sftp.Create(ctx, dest)
}

func (s *BackupService) VerifySFTPDestination(ctx context.Context, id string) error {
	if s.sftp == nil {
		return fmt.Errorf("sftp storage unavailable")
	}
	dest, err := s.sftp.Get(ctx, id)
	if err != nil {
		return err
	}
	return backup.VerifySFTP(ctx, dest)
}

func (s *BackupService) ListSFTPDestinations(ctx context.Context, projectID string) ([]models.SFTPDestination, error) {
	if s.sftp == nil {
		return nil, fmt.Errorf("sftp storage unavailable")
	}
	return s.sftp.ListByProject(ctx, projectID)
}

func (s *BackupService) DeleteSFTPDestination(ctx context.Context, id string) error {
	if s.sftp == nil {
		return fmt.Errorf("sftp storage unavailable")
	}
	return s.sftp.Delete(ctx, id)
}

func VerifySFTPDestination(ctx context.Context, dest *models.SFTPDestination) error {
	return backup.VerifySFTP(ctx, dest)
}
