package backup

import (
	"codedock.run/codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"
)

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
		volume, err := bm.ValidateVolumeRestore(ctx, id)
		if err != nil {
			return nil, err
		}
		target.VolumeName = volume.VolumeName
		fingerprint = fmt.Sprintf("%s|%s|%s|%s", volume.VolumeName, volume.VolumeCreatedAt, volume.ContainerID, volume.OwnerOperation)
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
		if project.ServerID != "" {
			return nil, fmt.Errorf("database restore requires a local Docker target")
		}

		if bm.dockerClient == nil {
			return nil, fmt.Errorf("Docker runtime unavailable")
		}
		copyConfig := *cfg
		copyConfig.DatabaseID = databaseID
		if _, _, err := bm.buildRestoreCommand(&copyConfig); err != nil {
			return nil, err
		}
		inspected, err := bm.dockerClient.ContainerInspect(ctx, databaseContainerIdentity(destination))
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
		fingerprint = fmt.Sprintf("%s|%s|%s|%s|%s", databaseID, inspected.ID, destination.Engine, destination.Username, destination.DatabaseName)
	}
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
		return bm.restoreVolume(ctx, current.RecordID, current.VolumeName, current.SHA256)
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
	record, err := bm.store.GetBackupRecord(current.RecordID)
	if err != nil {
		return err
	}
	cfg, err := bm.store.GetBackupConfig(record.BackupConfigID)
	if err != nil {
		return err
	}
	copyConfig := *cfg
	copyConfig.DatabaseID = current.DatabaseID
	_, command, err := bm.buildRestoreCommand(&copyConfig)
	if err != nil {
		return err
	}
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
	if err := bm.executeRestore(ctx, current.ContainerID, command, data); err != nil {
		return err
	}
	if err := progress("VERIFYING", "Restore command succeeded; checking target runtime"); err != nil {
		return err
	}
	inspected, err := bm.dockerClient.ContainerInspect(ctx, current.ContainerID)
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
