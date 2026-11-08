package backups

import (
	"context"
	"errors"
)

func (s *Service) GetRecord(ctx context.Context, id string) (*BackupRecord, error) {
	return s.backupRepo.GetRecordByID(ctx, id)
}

func (s *Service) DeleteRecord(ctx context.Context, id string) error {
	rec, err := s.backupRepo.GetRecordByID(ctx, id)
	if err != nil {
		return err
	}
	if s.manager != nil {
		if err := s.manager.DeleteBackupArtifacts(ctx, rec.ID); err != nil {
			return err
		}
	}
	return s.backupRepo.DeleteRecord(ctx, id)
}

func (s *Service) ListRecords(ctx context.Context, backupConfigID string) ([]*BackupRecord, error) {
	return s.backupRepo.ListRecordsByConfig(ctx, backupConfigID)
}

func (s *Service) ListAllRecords(ctx context.Context, limit int) ([]*BackupRecord, error) {
	return s.backupRepo.ListAllRecords(ctx, limit)
}

func (s *Service) ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*BackupRecord, error) {
	return s.backupRepo.ListRecordsByConfigs(ctx, configIDs, limit)
}

func (s *Service) ListRecordsByDatabase(ctx context.Context, databaseID string) ([]*BackupRecord, error) {
	return s.backupRepo.ListRecordsByDatabase(ctx, databaseID)
}

func (s *Service) GetRecordDownloadURL(ctx context.Context, recordID string) (string, error) {
	rec, err := s.backupRepo.GetRecordByID(ctx, recordID)
	if err != nil {
		return "", err
	}
	if s.manager == nil {
		return "", errors.New("backup manager not initialized")
	}
	return s.manager.GetDownloadURL(ctx, rec.ID)
}

func (s *Service) RecoverRecords(ctx context.Context) error {
	if recoverer, ok := s.backupRepo.(interface {
		RecoverRecords(context.Context) error
	}); ok {
		return recoverer.RecoverRecords(ctx)
	}
	return nil
}
