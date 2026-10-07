package backup

import (
	"bytes"
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"
)

func (bm *BackupManager) forConfig(ctx context.Context, cfg *models.BackupConfig) (*BackupManager, func(), error) {
	return bm.producer(ctx, cfg)
}

func (bm *BackupManager) protectRestoreSource(ctx context.Context, id string) (int64, error) {
	record, err := bm.store.GetBackupRecord(id)
	if err != nil {
		return 0, err
	}
	deadline := time.Now().Add(24 * time.Hour).Unix()
	if record.ProtectedUntil >= deadline {
		return 0, nil
	}
	protector, ok := bm.store.(interface {
		ProtectRecord(context.Context, string, int64) error
	})
	if !ok {
		return 0, fmt.Errorf("record protection unavailable")
	}
	if err := protector.ProtectRecord(ctx, id, deadline); err != nil {
		return 0, err
	}
	return deadline, nil
}

func (bm *BackupManager) clearRestoreProtection(ctx context.Context, id string, deadline int64) {
	if deadline <= 0 {
		return
	}
	record, err := bm.store.GetBackupRecord(id)
	if err != nil || record.ProtectedUntil != deadline {
		return
	}
	clearer, ok := bm.store.(interface {
		ClearRestoreProtection(context.Context, string, int64) error
	})
	if !ok {
		return
	}
	_ = clearer.ClearRestoreProtection(ctx, id, deadline)
}

func (bm *BackupManager) PrepareRestore(ctx context.Context, id, databaseID string) (*models.RestoreTarget, error) {
	digest, err := bm.VerifyArchive(ctx, id)
	if err != nil {
		return nil, err
	}
	record, err := bm.store.GetBackupRecord(id)
	if err != nil {
		return nil, err
	}
	cfg, err := bm.store.GetBackupConfig(record.BackupConfigID)
	if err != nil {
		return nil, err
	}
	target := &models.RestoreTarget{RecordID: id, SHA256: digest, Mode: "in-place"}
	var fingerprint string
	if cfg.VolumeName != "" {
		if databaseID != "" {
			return nil, fmt.Errorf("volume archives cannot be restored as database dumps")
		}
		resolved, release, err := bm.forConfig(ctx, cfg)
		if err != nil {
			return nil, err
		}
		defer release()
		volume, err := resolved.ValidateVolumeRestore(ctx, id)
		if err != nil {
			return nil, err
		}
		target.VolumeName = volume.VolumeName
		fingerprint = fmt.Sprintf("%s|%s|%s|%s|%s", cfg.ProjectID, volume.VolumeName, volume.VolumeCreatedAt, volume.ContainerID, volume.OwnerOperation)
	} else if strings.TrimSpace(cfg.FileSourcePath) != "" && cfg.DatabaseID == "" && strings.TrimSpace(cfg.CustomBackupCommand) == "" {
		cleaned := filepath.Clean("/" + strings.Trim(cfg.FileSourcePath, "/"))
		if cleaned == "/" || strings.Contains(cleaned, "..") {
			return nil, fmt.Errorf("invalid file source path")
		}
		target.Engine = "file"
		fingerprint = fmt.Sprintf("%s|%s|file", cfg.ProjectID, cleaned)
	} else if strings.TrimSpace(cfg.CustomRestoreCommand) != "" || (strings.TrimSpace(cfg.CustomBackupCommand) != "" && cfg.DatabaseID == "" && strings.TrimSpace(cfg.FileSourcePath) == "") {
		containerName, err := bm.serviceContainer(ctx, cfg.ServiceID)
		if err != nil {
			return nil, err
		}
		copyConfig := *cfg
		resolved, release, err := bm.forConfig(ctx, &copyConfig)
		if err != nil {
			return nil, err
		}
		defer release()
		if resolved.dockerClient == nil {
			return nil, fmt.Errorf("Docker runtime unavailable")
		}
		inspected, err := resolved.dockerClient.ContainerInspect(ctx, containerName)
		if err != nil {
			return nil, err
		}
		if inspected.State == nil || !inspected.State.Running {
			return nil, fmt.Errorf("custom restore target must be running")
		}
		target.Engine = "custom"
		target.ContainerID = inspected.ID
		target.DatabaseID = ""
		fingerprint = fmt.Sprintf("%s|%s|custom", cfg.ServiceID, inspected.ID)
	} else {
		source, err := bm.store.GetDatabase(cfg.DatabaseID)
		if err != nil {
			return nil, err
		}
		if databaseID == "" {
			databaseID = cfg.DatabaseID
		}
		destination, err := bm.store.GetDatabase(databaseID)
		if err != nil {
			return nil, err
		}
		if source == nil || destination == nil || source.Engine != destination.Engine {
			return nil, fmt.Errorf("restore target must use the same database engine")
		}
		projects, ok := bm.store.(interface {
			GetProject(string) (*models.ProjectConfig, error)
		})
		if !ok {
			return nil, fmt.Errorf("restore target validation unavailable")
		}
		project, err := projects.GetProject(destination.ProjectID)
		if err != nil || project == nil {
			return nil, fmt.Errorf("restore target project unavailable")
		}
		copyConfig := *cfg
		copyConfig.DatabaseID = databaseID
		resolved, release, err := bm.forConfig(ctx, &copyConfig)
		if err != nil {
			return nil, err
		}
		defer release()
		if resolved.dockerClient == nil {
			return nil, fmt.Errorf("Docker runtime unavailable")
		}
		if _, _, err := resolved.buildRestoreCommand(&copyConfig); err != nil {
			return nil, err
		}
		inspected, err := resolved.dockerClient.ContainerInspect(ctx, databaseContainerIdentity(destination))
		if err != nil {
			return nil, err
		}
		if inspected.State == nil || !inspected.State.Running {
			return nil, fmt.Errorf("database restore target must be running")
		}
		target.DatabaseID = databaseID
		target.Engine = string(destination.Engine)
		target.ContainerID = inspected.ID
		if databaseID != source.ID {
			target.Mode = "new-target"
		}
		fingerprint = fmt.Sprintf("%s|%s|%s|%s|%s|%s", project.ID, project.ServerID, databaseID, inspected.ID, destination.Engine, destination.DatabaseName)
	}
	deadline, err := bm.protectRestoreSource(ctx, id)
	if err != nil {
		return nil, err
	}
	target.ProtectedUntil = deadline
	hash := sha256.Sum256([]byte(id + "|" + digest + "|" + cfg.UpdatedAt + "|" + fingerprint))
	target.Snapshot = hex.EncodeToString(hash[:])
	return target, nil
}

func (bm *BackupManager) ApplyRestore(ctx context.Context, reviewed *models.RestoreTarget, progress func(string, string) error) error {
	duration, err := bm.restoreTimeout(reviewed.RecordID)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()
	current, err := bm.PrepareRestore(ctx, reviewed.RecordID, reviewed.DatabaseID)
	if err != nil {
		return err
	}
	if current.Snapshot != reviewed.Snapshot {
		return fmt.Errorf("restore target or archive changed since review")
	}
	if err := progress("RESTORING", "Restore may overwrite target data. Interruption can leave partial data."); err != nil {
		return err
	}
	if current.VolumeName != "" {
		record, err := bm.store.GetBackupRecord(current.RecordID)
		if err != nil {
			return err
		}
		cfg, err := bm.store.GetBackupConfig(record.BackupConfigID)
		if err != nil {
			return err
		}
		resolved, release, err := bm.forConfig(ctx, cfg)
		if err != nil {
			return err
		}
		defer release()
		defer bm.clearRestoreProtection(context.Background(), current.RecordID, reviewed.ProtectedUntil)
		return resolved.restoreVolume(ctx, current.RecordID, current.VolumeName, current.SHA256)
	}
	record, err := bm.store.GetBackupRecord(current.RecordID)
	if err != nil {
		return err
	}
	nativeCfg, err := bm.store.GetBackupConfig(record.BackupConfigID)
	if err != nil {
		return err
	}
	if nativeCfg.ServiceID != "" && bm.isNativeService(ctx, nativeCfg.ServiceID) {
		if bm.nativeRestore == nil {
			return fmt.Errorf("native restore transport unavailable")
		}
		defer bm.clearRestoreProtection(context.Background(), current.RecordID, reviewed.ProtectedUntil)
		stream, _, err := bm.OpenBackupArchive(ctx, current.RecordID)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(io.LimitReader(stream, 1024*1024*1024+1))
		closeErr := stream.Close()
		if readErr != nil {
			return readErr
		}
		if closeErr != nil {
			return closeErr
		}
		digest := sha256.Sum256(data)
		if hex.EncodeToString(digest[:]) != current.SHA256 {
			return fmt.Errorf("restore archive checksum changed")
		}
		if err := progress("RESTORING", "Restore may overwrite native data. Interruption can leave partial data."); err != nil {
			return err
		}
		if err := bm.nativeRestore(ctx, nativeCfg.ServiceID, bytes.NewReader(data)); err != nil {
			return err
		}
		return progress("VERIFIED", "Native data directory restored")
	}
	if current.Engine == "file" {
		defer bm.clearRestoreProtection(context.Background(), current.RecordID, reviewed.ProtectedUntil)
		stream, _, err := bm.OpenBackupArchive(ctx, current.RecordID)
		if err != nil {
			return err
		}
		data, readErr := io.ReadAll(io.LimitReader(stream, 1024*1024*1024+1))
		closeErr := stream.Close()
		if readErr != nil {
			return readErr
		}
		if closeErr != nil {
			return closeErr
		}
		digest := sha256.Sum256(data)
		if hex.EncodeToString(digest[:]) != current.SHA256 {
			return fmt.Errorf("restore archive checksum changed")
		}
		if err := progress("RESTORING", "Restore may overwrite target files. Interruption can leave partial data."); err != nil {
			return err
		}
		record, err := bm.store.GetBackupRecord(current.RecordID)
		if err != nil {
			return err
		}
		cfg, err := bm.store.GetBackupConfig(record.BackupConfigID)
		if err != nil {
			return err
		}
		resolved, release, err := bm.forConfig(ctx, cfg)
		if err != nil {
			return err
		}
		defer release()
		if err := resolved.restoreFile(ctx, cfg.FileSourcePath, data); err != nil {
			return err
		}
		return progress("VERIFIED", "File source restored")
	}
	if bm.volumeOperations != nil {
		release, err := bm.volumeOperations.AcquireVolume("codedock-db-data-" + current.DatabaseID)
		if err != nil {
			return err
		}
		defer release()
	}
	locked, err := bm.PrepareRestore(ctx, current.RecordID, current.DatabaseID)
	if err != nil {
		return err
	}
	if locked.Snapshot != reviewed.Snapshot {
		return fmt.Errorf("restore target changed before execution")
	}
	record, err = bm.store.GetBackupRecord(current.RecordID)
	if err != nil {
		return err
	}
	cfg, err := bm.store.GetBackupConfig(record.BackupConfigID)
	if err != nil {
		return err
	}
	copyConfig := *cfg
	copyConfig.DatabaseID = current.DatabaseID
	resolved, release, err := bm.forConfig(ctx, &copyConfig)
	if err != nil {
		return err
	}
	defer release()
	defer bm.clearRestoreProtection(context.Background(), current.RecordID, reviewed.ProtectedUntil)
	if resolved.dockerClient == nil {
		return fmt.Errorf("Docker runtime unavailable")
	}
	_, command, err := resolved.buildRestoreCommand(&copyConfig)
	if err != nil {
		return err
	}
	stream, _, err := resolved.OpenBackupArchive(ctx, current.RecordID)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	data, readErr := io.ReadAll(io.LimitReader(stream, 1024*1024*1024+1))
	closeErr := stream.Close()
	if readErr != nil {
		return readErr
	}
	if closeErr != nil {
		return closeErr
	}
	digest := sha256.Sum256(data)
	if hex.EncodeToString(digest[:]) != current.SHA256 {
		return fmt.Errorf("restore archive checksum changed")
	}
	if err := resolved.executeRestore(ctx, current.ContainerID, command, data); err != nil {
		return err
	}
	if err := progress("VERIFYING", "Restore command succeeded; checking target runtime"); err != nil {
		return err
	}
	inspected, err := resolved.dockerClient.ContainerInspect(ctx, current.ContainerID)
	if err != nil {
		return err
	}
	if inspected.State == nil || !inspected.State.Running {
		return fmt.Errorf("restored database runtime is not running")
	}
	if inspected.State.Health != nil && inspected.State.Health.Status != "healthy" {
		return fmt.Errorf("restored database is not healthy")
	}
	return nil
}

func (bm *BackupManager) restoreTimeout(id string) (time.Duration, error) {
	record, err := bm.store.GetBackupRecord(id)
	if err != nil {
		return 0, err
	}
	cfg, err := bm.store.GetBackupConfig(record.BackupConfigID)
	if err != nil {
		return 0, err
	}
	seconds := cfg.Timeout
	if seconds <= 0 {
		seconds = 3600
	}
	if seconds > 86400 {
		seconds = 86400
	}
	return time.Duration(seconds) * time.Second, nil
}
