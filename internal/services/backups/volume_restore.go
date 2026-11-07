package backups

import (
	"codedock/internal/models"
	"context"
	"errors"
)

func (s *BackupService) ValidateVolumeRestore(ctx context.Context, id string) (*models.VolumeRestoreTarget, error) {
	if s.manager == nil {
		return nil, errors.New("backup manager unavailable")
	}
	return s.manager.ValidateVolumeRestore(ctx, id)
}
func (s *BackupService) RestoreVolume(ctx context.Context, id, volume string) error {
	if s.manager == nil {
		return errors.New("backup manager unavailable")
	}
	return s.manager.RestoreVolume(ctx, id, volume)
}
func (s *BackupService) CancelVolumeRestore(id string) bool {
	return s.manager != nil && s.manager.CancelVolumeRestore(id)
}
