package backups

import (
	"context"
	"encoding/json"
	"fmt"
)

type Manager interface {
	TriggerBackup(ctx context.Context, backupConfigID string) (*BackupRecord, error)
	RestoreBackup(ctx context.Context, recordID string) error
	UnregisterBackup(backupConfigID string)
	GetDownloadURL(ctx context.Context, recordID string) (string, error)
	DeleteBackupArtifacts(ctx context.Context, recordID string) error
	PrepareRestore(ctx context.Context, id, databaseID string) (*RestoreTarget, error)
	ApplyDatabaseRestore(ctx context.Context, recordID, databaseID string) error
	RestoreVolume(ctx context.Context, recordID, serverID, volumeName string) error
	VolumeRestoreTarget(ctx context.Context, recordID string) (string, string, error)
	CancelVolumeRestore(recordID string) error
}

type engineBridge interface {
	TriggerBackup(ctx context.Context, backupConfigID string) (any, error)
	RestoreBackup(ctx context.Context, recordID string) error
	UnregisterBackup(backupConfigID string)
	GetDownloadURL(ctx context.Context, recordID string) (string, error)
	DeleteBackupArtifacts(ctx context.Context, recordID string) error
	PrepareRestore(ctx context.Context, id, databaseID string) (any, error)
	ApplyDatabaseRestore(ctx context.Context, recordID, databaseID string) error
	RestoreVolume(ctx context.Context, recordID, serverID, volumeName string) error
	VolumeRestoreTarget(ctx context.Context, recordID string) (string, string, error)
	CancelVolumeRestore(recordID string) error
}

type wrappedManager struct {
	engine any
}

func NewManager(engine any) Manager {
	if m, ok := engine.(Manager); ok {
		return m
	}
	return &wrappedManager{engine: engine}
}

func (m *wrappedManager) TriggerBackup(ctx context.Context, backupConfigID string) (*BackupRecord, error) {
	if m.engine == nil {
		return nil, fmt.Errorf("backup engine unavailable")
	}
	type triggerer interface {
		TriggerBackup(context.Context, string) (any, error)
	}
	t, ok := m.engine.(triggerer)
	if !ok {
		return nil, fmt.Errorf("engine does not support TriggerBackup")
	}
	res, err := t.TriggerBackup(ctx, backupConfigID)
	if err != nil {
		return nil, err
	}
	return convertToBackupRecord(res)
}

func (m *wrappedManager) RestoreBackup(ctx context.Context, recordID string) error {
	type restorer interface {
		RestoreBackup(context.Context, string) error
	}
	if r, ok := m.engine.(restorer); ok {
		return r.RestoreBackup(ctx, recordID)
	}
	return fmt.Errorf("engine does not support RestoreBackup")
}

func (m *wrappedManager) UnregisterBackup(backupConfigID string) {
	type unregisterer interface {
		UnregisterBackup(string)
	}
	if u, ok := m.engine.(unregisterer); ok {
		u.UnregisterBackup(backupConfigID)
	}
}

func (m *wrappedManager) GetDownloadURL(ctx context.Context, recordID string) (string, error) {
	type downloader interface {
		GetDownloadURL(context.Context, string) (string, error)
	}
	if d, ok := m.engine.(downloader); ok {
		return d.GetDownloadURL(ctx, recordID)
	}
	return "", fmt.Errorf("engine does not support GetDownloadURL")
}

func (m *wrappedManager) DeleteBackupArtifacts(ctx context.Context, recordID string) error {
	type deleter interface {
		DeleteBackupArtifacts(context.Context, string) error
	}
	if d, ok := m.engine.(deleter); ok {
		return d.DeleteBackupArtifacts(ctx, recordID)
	}
	return nil
}

func (m *wrappedManager) PrepareRestore(ctx context.Context, id, databaseID string) (*RestoreTarget, error) {
	type preparer interface {
		PrepareRestore(context.Context, string, string) (any, error)
	}
	p, ok := m.engine.(preparer)
	if !ok {
		return nil, fmt.Errorf("engine does not support PrepareRestore")
	}
	res, err := p.PrepareRestore(ctx, id, databaseID)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	var target RestoreTarget
	if err := json.Unmarshal(data, &target); err != nil {
		return nil, err
	}
	return &target, nil
}

func (m *wrappedManager) ApplyDatabaseRestore(ctx context.Context, recordID, databaseID string) error {
	type applier interface {
		ApplyDatabaseRestore(context.Context, string, string) error
	}
	if a, ok := m.engine.(applier); ok {
		return a.ApplyDatabaseRestore(ctx, recordID, databaseID)
	}
	return fmt.Errorf("engine does not support ApplyDatabaseRestore")
}

func (m *wrappedManager) RestoreVolume(ctx context.Context, recordID, serverID, volumeName string) error {
	type restorer interface {
		RestoreVolume(context.Context, string, string, string) error
	}
	if r, ok := m.engine.(restorer); ok {
		return r.RestoreVolume(ctx, recordID, serverID, volumeName)
	}
	return fmt.Errorf("engine does not support RestoreVolume")
}

func (m *wrappedManager) VolumeRestoreTarget(ctx context.Context, recordID string) (string, string, error) {
	type resolver interface {
		VolumeRestoreTarget(context.Context, string) (string, string, error)
	}
	if r, ok := m.engine.(resolver); ok {
		return r.VolumeRestoreTarget(ctx, recordID)
	}
	return "", "", fmt.Errorf("engine does not support VolumeRestoreTarget")
}

func (m *wrappedManager) CancelVolumeRestore(recordID string) error {
	type canceller interface {
		CancelVolumeRestore(string) error
	}
	if c, ok := m.engine.(canceller); ok {
		return c.CancelVolumeRestore(recordID)
	}
	return nil
}

func convertToBackupRecord(src any) (*BackupRecord, error) {
	if src == nil {
		return nil, nil
	}
	if rec, ok := src.(*BackupRecord); ok {
		return rec, nil
	}
	data, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	var rec BackupRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}
