package backups

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type OperationsService interface {
	Review(ctx context.Context, user, project, kind, target, payload, snapshot, effects string) (any, error)
	Apply(ctx context.Context, id, user, confirmation, snapshot string, runner func(ctx context.Context) error) error
	GetPayload(ctx context.Context, id string) (string, error)
}

type Service struct {
	RunAuthorization func(context.Context, string, string) error
	operations       OperationsService
	backupRepo       ConfigRepository
	s3Repo           S3Repository
	batches          BatchRepository
	sftp             SFTPRepository
	manager          Manager
}

func NewService(br ConfigRepository, sr S3Repository, m Manager) *Service {
	svc := &Service{
		backupRepo: br,
		s3Repo:     sr,
		manager:    m,
	}
	if authz, ok := br.(interface {
		AuthorizeScheduledRun(context.Context, string, string) error
	}); ok {
		svc.RunAuthorization = authz.AuthorizeScheduledRun
	}
	return svc
}

func (s *Service) SetOperations(service any) {
	if os, ok := service.(OperationsService); ok {
		s.operations = os
	}
}

func (s *Service) SetBatches(store BatchRepository) { s.batches = store }
func (s *Service) SetSFTP(store SFTPRepository)     { s.sftp = store }

func (s *Service) CreateConfig(ctx context.Context, cfg *BackupConfig) error {
	if cfg == nil {
		return ErrValidBackupConfigRequired
	}
	if cfg.ID == "" {
		cfg.ID = "bkc_" + uuid.New().String()
	}
	if cfg.Schedule == "" {
		cfg.Schedule = DefaultSchedule
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = DefaultTimeout
	}
	setScheduleStatus(cfg)
	return s.backupRepo.CreateConfig(ctx, cfg)
}

func (s *Service) GetConfig(ctx context.Context, id string) (*BackupConfig, error) {
	if id == "" {
		return nil, ErrBackupIDRequired
	}
	cfg, err := s.backupRepo.GetConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, ErrBackupConfigNotFound
	}
	return cfg, nil
}

func (s *Service) UpdateConfig(ctx context.Context, cfg *BackupConfig) error {
	if cfg == nil || cfg.ID == "" {
		return ErrValidBackupConfigRequired
	}
	setScheduleStatus(cfg)
	return s.backupRepo.UpdateConfig(ctx, cfg)
}

func (s *Service) DeleteConfig(ctx context.Context, id string) error {
	if id == "" {
		return ErrBackupIDRequired
	}
	return s.backupRepo.DeleteConfig(ctx, id)
}

func (s *Service) ListConfigs(ctx context.Context) ([]*BackupConfig, error) {
	return s.backupRepo.ListConfigs(ctx)
}

func (s *Service) ListAllConfigs(ctx context.Context) ([]*BackupConfig, error) {
	return s.backupRepo.ListConfigs(ctx)
}

func (s *Service) TriggerBackup(ctx context.Context, id string) (*BackupRecord, error) {
	if s.manager == nil {
		return nil, errors.New("backup manager not initialized")
	}
	return s.manager.TriggerBackup(ctx, id)
}

func (s *Service) RestoreBackup(ctx context.Context, id string) error {
	if s.manager == nil {
		return errors.New("backup manager not initialized")
	}
	return s.manager.RestoreBackup(ctx, id)
}

func setScheduleStatus(cfg *BackupConfig) {
	if cfg == nil {
		return
	}
	if cfg.NextRunAt == "" && cfg.BackupEnabled {
		if next, err := NextRun(cfg.Schedule, cfg.Timezone); err == nil && !next.IsZero() {
			cfg.NextRunAt = next.UTC().Format(time.RFC3339)
		}
	}
	if cfg.Status == "" {
		if cfg.BackupEnabled {
			cfg.Status = ConfigStatusActive
		} else {
			cfg.Status = ConfigStatusInactive
		}
	}
}
