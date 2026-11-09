package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (bm *BackupManager) TriggerBackup(ctx context.Context, backupConfigID string) (*BackupRecord, error) {
	cfg, err := bm.store.GetBackupConfig(backupConfigID)
	if err != nil || cfg == nil {
		return nil, fmt.Errorf("backup config %s not found: %w", backupConfigID, err)
	}
	if !cfg.BackupEnabled {
		return nil, fmt.Errorf("backup execution disabled for config %s", backupConfigID)
	}

	producer, release, err := bm.producer(ctx, cfg)
	if err != nil {
		return nil, err
	}
	defer release()
	if producer != bm {
		if cfg.Incremental {
			return nil, fmt.Errorf("incremental backups require a local Docker target")
		}
		return producer.TriggerBackup(ctx, backupConfigID)
	}
	rec := &BackupRecord{
		ID:                uuid.New().String(),
		BackupConfigID:    cfg.ID,
		DatabaseID:        cfg.DatabaseID,
		S3DestinationID:   cfg.S3DestinationID,
		SFTPDestinationID: cfg.SFTPDestinationID,
		Status:            BackupRecordStatusRunning,
		Logs:              fmt.Sprintf("Initiating automated backup '%s' at %s...\n", cfg.Name, time.Now().UTC().Format(time.RFC3339)),
	}
	if err := bm.store.CreateBackupRecord(rec); err != nil {
		return nil, fmt.Errorf("failed to create backup record: %w", err)
	}

	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 3600
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	if err := backupProgress(ctx, "PRODUCING", "Creating backup record "+rec.ID); err != nil {
		return bm.failBackupRecord(rec.ID, err.Error())
	}
	var dumpBytes []byte
	var execLogs string
	var fileExt string

	if cfg.VolumeName != "" && !bm.isNativeService(ctx, cfg.ServiceID) {
		if err := bm.validateVolumeProducer(ctx, cfg); err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		containerName := ""
		if cfg.ServiceID != "" {
			if resolved, err := bm.serviceContainer(ctx, cfg.ServiceID); err == nil {
				containerName = resolved
			}
		}
		unquiesce, err := bm.quiesceTarget(ctx, cfg, containerName)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		defer unquiesce()
		if cfg.Incremental {
			dumpBytes, execLogs, err = bm.executeVolumeBackupIncremental(ctx, cfg, cfg.VolumeName)
			if parent := bm.latestCompletedRecord(ctx, cfg.ID); parent != "" {
				rec.ParentRecordID = parent
			}
		} else {
			dumpBytes, execLogs, err = bm.executeVolumeBackup(ctx, cfg.VolumeName)
		}
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		fileExt = ".tar.gz"
	} else if cfg.FileSourcePath != "" && !bm.isNativeService(ctx, cfg.ServiceID) {
		dumpBytes, execLogs, err = bm.executeFileBackup(ctx, cfg.FileSourcePath)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		fileExt = ".tar.gz"
	} else if cfg.ServiceID != "" && bm.isNativeService(ctx, cfg.ServiceID) {
		if bm.nativeBackup == nil {
			return bm.failBackupRecord(rec.ID, "native backup transport unavailable")
		}
		var out bytes.Buffer
		if err := bm.nativeBackup(ctx, cfg.ServiceID, &out); err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		dumpBytes = out.Bytes()
		execLogs = "Native data directory archived.\n"
		fileExt = ".tar.gz"
	} else if strings.TrimSpace(cfg.CustomBackupCommand) != "" {
		containerName, err := bm.serviceContainer(ctx, cfg.ServiceID)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		unquiesce, err := bm.quiesceTarget(ctx, cfg, containerName)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		defer unquiesce()
		dumpBytes, execLogs, err = bm.executeCustomBackup(ctx, cfg)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		fileExt = ".bin"
	} else if cfg.DatabaseID == "global" || cfg.DatabaseID == "" {
		snapshotter, ok := bm.store.(interface {
			ControlPlaneSnapshot(context.Context) ([]byte, error)
		})
		if !ok {
			return bm.failBackupRecord(rec.ID, "consistent control-plane snapshots unavailable")
		}
		content, err := snapshotter.ControlPlaneSnapshot(ctx)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		dumpBytes = content
		fileExt = ".db"
		execLogs = "Global database backed up successfully.\n"
	} else {
		containerName, dumpCmd, ext, err := bm.buildDumpCommand(cfg)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		unquiesce, err := bm.quiesceTarget(ctx, cfg, containerName)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		defer unquiesce()
		fileExt = ext
		dumpBytes, execLogs, err = bm.executeDump(ctx, containerName, dumpCmd, cfg.Name)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
	}

	if err := backupProgress(ctx, "VERIFYING", "Checking produced archive"); err != nil {
		return bm.failBackupRecord(rec.ID, err.Error())
	}
	if len(dumpBytes) == 0 {
		return bm.failBackupRecord(rec.ID, "backup producer returned an empty archive")
	}
	digest := sha256.Sum256(dumpBytes)
	rec.SHA256 = hex.EncodeToString(digest[:])
	rec.VerifiedAt = time.Now().UTC().Format(time.RFC3339)
	fileName := fmt.Sprintf("backup_%s_%s%s", cfg.ID, time.Now().UTC().Format("20060102_150405")+"_"+rec.ID, fileExt)
	filePath := filepath.Join(bm.backupDir, fileName)

	if err := os.WriteFile(filePath, dumpBytes, 0o600); err != nil {
		return bm.failBackupRecord(rec.ID, fmt.Sprintf("failed to write backup archive to disk: %v", err))
	}

	sizeBytes := int64(len(dumpBytes))
	if fileInfo, err := os.Stat(filePath); err == nil && fileInfo != nil {
		sizeBytes = fileInfo.Size()
	}

	s3URL := ""
	if cfg.S3Enabled {
		if err := backupProgress(ctx, "UPLOADING", "Uploading verified archive"); err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		var s3Err error
		s3URL, execLogs, s3Err = bm.handleS3Upload(ctx, cfg, fileName, dumpBytes, execLogs)
		if s3Err != nil {
			if err := bm.store.UpdateBackupRecord(UpdateBackupRecordOpts{
				ID:            rec.ID,
				Status:        BackupRecordStatusFailed,
				FilePath:      filePath,
				FileSizeBytes: sizeBytes,
				Logs:          execLogs + fmt.Sprintf("Failed: %s\n", s3Err.Error()),
				CompletedAt:   time.Now().UTC().Format(time.RFC3339),
			}); err != nil {
				slog.Warn("failed to update backup record on s3 upload failure", "error", err)
			}
			return nil, s3Err
		}
	}

	sftpURL := ""
	if cfg.SFTPEnabled {
		if err := backupProgress(ctx, "UPLOADING", "Uploading verified archive to SFTP"); err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		var sftpErr error
		sftpURL, execLogs, sftpErr = bm.handleSFTPUpload(ctx, cfg, fileName, dumpBytes, execLogs)
		if sftpErr != nil {
			if err := bm.store.UpdateBackupRecord(UpdateBackupRecordOpts{
				ID:            rec.ID,
				Status:        BackupRecordStatusFailed,
				FilePath:      filePath,
				FileSizeBytes: sizeBytes,
				Logs:          execLogs + fmt.Sprintf("Failed: %s\n", sftpErr.Error()),
				CompletedAt:   time.Now().UTC().Format(time.RFC3339),
			}); err != nil {
				slog.Warn("failed to update backup record on sftp upload failure", "error", err)
			}
			return nil, sftpErr
		}
	}

	if cfg.DisableLocal && (s3URL != "" || sftpURL != "") {
		_ = os.Remove(filePath)
		filePath = ""
	}

	completed, err := bm.finalizeBackupRecord(FinalizeBackupOpts{
		Record:         rec,
		FilePath:       filePath,
		S3URL:          s3URL,
		SFTPURL:        sftpURL,
		ParentRecordID: rec.ParentRecordID,
		ExecLogs:       execLogs,
		SizeBytes:      sizeBytes,
	})
	if err != nil {
		return nil, err
	}
	bm.enforceRetentionPolicy(cfg)
	return completed, nil
}

func (bm *BackupManager) latestCompletedRecord(ctx context.Context, configID string) string {
	records, err := bm.store.ListBackupRecords(configID)
	if err != nil {
		return ""
	}
	for _, record := range records {
		if record.Status == BackupRecordStatusCompleted {
			return record.ID
		}
	}
	return ""
}

func (bm *BackupManager) handleSFTPUpload(ctx context.Context, cfg *BackupConfig, fileName string, dumpBytes []byte, execLogs string) (string, string, error) {
	if cfg.SFTPDestinationID == "" {
		return "", execLogs, fmt.Errorf("SFTP destination ID missing for config %s", cfg.ID)
	}
	getter, ok := bm.store.(interface {
		GetSFTPDestination(string) (*SFTPDestination, error)
	})
	if !ok {
		return "", execLogs, fmt.Errorf("SFTP destination lookup unavailable")
	}
	dest, err := getter.GetSFTPDestination(cfg.SFTPDestinationID)
	if err != nil || dest == nil {
		return "", execLogs, fmt.Errorf("SFTP destination %s not found", cfg.SFTPDestinationID)
	}
	sftpURL, err := sftpPut(ctx, dest, fileName, dumpBytes)
	if err != nil {
		return "", execLogs, fmt.Errorf("SFTP upload failed: %w", err)
	}
	execLogs += fmt.Sprintf("\nUploaded backup to SFTP destination: %s", sftpURL)
	return sftpURL, execLogs, nil
}
