package adapters

import (
	"context"
	"errors"
	"fmt"

	"codedock/internal/modules/backups"
	"codedock/internal/platform/storage"
	"codedock/internal/utils"
)

func (a *EngineAdapter) ClaimRecordExpiry(ctx context.Context, id string) (bool, error) {
	return a.backupRepo.ClaimRecordExpiry(ctx, id)
}

func (a *EngineAdapter) ControlPlaneSnapshot(ctx context.Context) ([]byte, error) {
	return a.backupRepo.ControlPlaneSnapshot(ctx)
}

func (a *EngineAdapter) VolumeRestoreOwner(ctx context.Context, cfg *storage.BackupConfig) (string, error) {
	return a.volumeOwner(ctx, cfg, false)
}

func (a *EngineAdapter) VolumeBackupOwner(ctx context.Context, cfg *storage.BackupConfig) (string, error) {
	return a.volumeOwner(ctx, cfg, true)
}

func (a *EngineAdapter) volumeOwner(ctx context.Context, cfg *storage.BackupConfig, allowRemote bool) (string, error) {
	var projectID, containerID string
	if cfg.ServiceID != "" && cfg.DatabaseID != "" {
		return "", errors.New("volume backup must have exactly one owner")
	}
	if cfg.ServiceID != "" {
		app, err := a.appRepo.GetByID(ctx, cfg.ServiceID)
		if err != nil {
			return "", fmt.Errorf("load volume service: %w", err)
		}
		projectID = app.ProjectID
		containerID = utils.NormalizeContainerName(app.ID)
		if app.ContainerID != "" {
			containerID = app.ContainerID
		}
	} else if cfg.DatabaseID != "" && cfg.DatabaseID != "global" {
		db, err := a.dbRepo.GetByID(ctx, cfg.DatabaseID)
		if err != nil {
			return "", fmt.Errorf("load volume database: %w", err)
		}
		projectID = db.ProjectID
		containerID = db.ContainerID
		if containerID == "" {
			containerID = utils.NormalizeContainerName("codedock-db-" + db.Name)
		}
		if cfg.VolumeName != "codedock-db-data-"+db.ID {
			return "", errors.New("volume is not the database data volume")
		}
	} else {
		return "", errors.New("volume backup requires a registered service or database owner")
	}
	project, err := a.projectRepo.Get(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("load volume project: %w", err)
	}
	if project.ServerID != "" && !allowRemote {
		return "", errors.New("volume restore supports local Docker targets only")
	}
	return containerID, nil
}

func toStorageBackupConfig(c *backups.BackupConfig) *storage.BackupConfig {
	if c == nil {
		return nil
	}
	return &storage.BackupConfig{
		ID:                   c.ID,
		OwnerID:              c.OwnerID,
		ProjectID:            c.ProjectID,
		PreDeployment:        c.PreDeployment,
		DatabaseID:           c.DatabaseID,
		ServiceID:            c.ServiceID,
		VolumeName:           c.VolumeName,
		S3DestinationID:      c.S3DestinationID,
		SFTPDestinationID:    c.SFTPDestinationID,
		ParentBatchID:        c.ParentBatchID,
		Name:                 c.Name,
		Description:          c.Description,
		DbUser:               c.DbUser,
		DbPassword:           c.DbPassword,
		BackupEnabled:        c.BackupEnabled,
		S3Enabled:            c.S3Enabled,
		SFTPEnabled:          c.SFTPEnabled,
		Incremental:          c.Incremental,
		DisableLocal:         c.DisableLocal,
		QuiesceCommand:       c.QuiesceCommand,
		UnquiesceCommand:     c.UnquiesceCommand,
		CustomBackupCommand:  c.CustomBackupCommand,
		CustomRestoreCommand: c.CustomRestoreCommand,
		FileSourcePath:       c.FileSourcePath,
		Schedule:             c.Schedule,
		Timezone:             c.Timezone,
		Timeout:              c.Timeout,
		RetentionDays:        c.RetentionDays,
		MaxBackups:           c.MaxBackups,
		MaxStorageGB:         c.MaxStorageGB,
		Status:               storage.BackupConfigStatus(c.Status),
		CreatedAt:            c.CreatedAt,
		UpdatedAt:            c.UpdatedAt,
	}
}

func toStorageBackupRecord(r *backups.BackupRecord) *storage.BackupRecord {
	if r == nil {
		return nil
	}
	return &storage.BackupRecord{
		ID:                r.ID,
		BackupConfigID:    r.BackupConfigID,
		DatabaseID:        r.DatabaseID,
		S3DestinationID:   r.S3DestinationID,
		SFTPDestinationID: r.SFTPDestinationID,
		Status:            storage.BackupRecordStatus(r.Status),
		FilePath:          r.FilePath,
		FileSizeBytes:     r.FileSizeBytes,
		S3URL:             r.S3URL,
		SFTPURL:           r.SFTPURL,
		Logs:              r.Logs,
		StartedAt:         r.StartedAt,
		CompletedAt:       r.CompletedAt,
		ProtectedUntil:    r.ProtectedUntil,
		ParentRecordID:    r.ParentRecordID,
		SHA256:            r.SHA256,
		VerifiedAt:        r.VerifiedAt,
	}
}

func toModuleBackupRecord(r *storage.BackupRecord) *backups.BackupRecord {
	if r == nil {
		return nil
	}
	return &backups.BackupRecord{
		ID:                r.ID,
		BackupConfigID:    r.BackupConfigID,
		DatabaseID:        r.DatabaseID,
		S3DestinationID:   r.S3DestinationID,
		SFTPDestinationID: r.SFTPDestinationID,
		Status:            backups.RecordStatus(r.Status),
		FilePath:          r.FilePath,
		FileSizeBytes:     r.FileSizeBytes,
		S3URL:             r.S3URL,
		SFTPURL:           r.SFTPURL,
		Logs:              r.Logs,
		StartedAt:         r.StartedAt,
		CompletedAt:       r.CompletedAt,
		ProtectedUntil:    r.ProtectedUntil,
		ParentRecordID:    r.ParentRecordID,
		SHA256:            r.SHA256,
		VerifiedAt:        r.VerifiedAt,
	}
}
