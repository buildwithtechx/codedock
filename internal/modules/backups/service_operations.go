package backups

import (
	"context"
	"encoding/json"
	"fmt"
)

func (s *Service) ProtectRecord(ctx context.Context, id string, until int64) error {
	store, ok := s.backupRepo.(interface {
		ProtectRecord(context.Context, string, int64) error
	})
	if !ok {
		return ErrProtectionUnavailable
	}
	return store.ProtectRecord(ctx, id, until)
}

func (s *Service) PrepareRestore(ctx context.Context, user, project, id, database string) (any, error) {
	if s.manager == nil || s.operations == nil {
		return nil, ErrRestoreRuntimeUnavailable
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

func (s *Service) ApplyRestore(ctx context.Context, user, id, confirmation string) error {
	if s.operations == nil {
		return fmt.Errorf("operations service unavailable")
	}
	payload, err := s.operations.GetPayload(ctx, id)
	if err != nil {
		return err
	}
	if payload == "" {
		return fmt.Errorf("restore operation expired or not found")
	}
	if s.manager == nil {
		return ErrRestoreRuntimeUnavailable
	}
	var target RestoreTarget
	if err := json.Unmarshal([]byte(payload), &target); err != nil {
		return fmt.Errorf("malformed restore operation: %w", err)
	}
	current, err := s.manager.PrepareRestore(ctx, target.RecordID, target.DatabaseID)
	if err != nil {
		return err
	}
	return s.operations.Apply(ctx, id, user, confirmation, current.Snapshot, func(runCtx context.Context) error {
		if target.VolumeName != "" {
			return s.manager.RestoreVolume(runCtx, target.RecordID, target.ServerID, target.VolumeName)
		}
		return s.manager.ApplyDatabaseRestore(runCtx, target.RecordID, target.DatabaseID)
	})
}

func (s *Service) VolumeRestoreTarget(ctx context.Context, recordID string) (string, string, error) {
	if s.manager == nil {
		return "", "", ErrRestoreRuntimeUnavailable
	}
	return s.manager.VolumeRestoreTarget(ctx, recordID)
}

func (s *Service) RestoreVolume(ctx context.Context, recordID, targetServerID, targetVolume string) error {
	if s.manager == nil {
		return ErrRestoreRuntimeUnavailable
	}
	return s.manager.RestoreVolume(ctx, recordID, targetServerID, targetVolume)
}

func (s *Service) CancelVolumeRestore(recordID string) bool {
	if s.manager == nil {
		return false
	}
	return s.manager.CancelVolumeRestore(recordID) == nil
}
