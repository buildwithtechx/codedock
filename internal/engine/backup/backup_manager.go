package backup

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

type Store interface {
	ListAllActiveBackupConfigs() ([]*models.BackupConfig, error)
	GetBackupConfig(id string) (*models.BackupConfig, error)
	CreateBackupRecord(rec *models.BackupRecord) error
	GetDatabase(id string) (*models.Database, error)
	UpdateBackupRecord(opts models.UpdateBackupRecordOpts) error
	GetS3Destination(id string) (*models.S3Destination, error)
	GetBackupRecord(id string) (*models.BackupRecord, error)
	ListBackupRecords(backupConfigID string) ([]*models.BackupRecord, error)
}

type BackupManager struct {
	scheduledRunner  func(context.Context, string) error
	dockerClient     *client.Client
	store            Store
	cronEngine       *cron.Cron
	entries          map[string]cron.EntryID
	backupDir        string
	mu               sync.Mutex
	restoreMu        sync.Mutex
	restores         map[string]context.CancelFunc
	restoreVolumes   map[string]string
	volumeOperations VolumeOperations
}

func NewBackupManager(dockerClient *client.Client, s Store, backupDir string) *BackupManager {
	if backupDir == "" {
		backupDir = filepath.Join(utils.GetDataDir(), "backups")
	}
	_ = os.MkdirAll(backupDir, 0o700)
	return &BackupManager{
		dockerClient:   dockerClient,
		store:          s,
		cronEngine:     cron.New(cron.WithSeconds()),
		entries:        make(map[string]cron.EntryID),
		restores:       make(map[string]context.CancelFunc),
		restoreVolumes: make(map[string]string),
		backupDir:      backupDir,
	}
}

func (bm *BackupManager) Start() error {
	cfgs, err := bm.store.ListAllActiveBackupConfigs()
	if err != nil {
		return fmt.Errorf("failed to load active backup configs during start: %w", err)
	}
	bm.mu.Lock()
	defer bm.mu.Unlock()
	for _, cfg := range cfgs {
		if err := bm.registerBackupLocked(cfg); err != nil {
			slog.Warn("failed to schedule backup", "name", cfg.Name, "id", cfg.ID, "err", err)
		}
	}
	bm.cronEngine.Start()
	slog.Info("backup manager started")
	return nil
}

func (bm *BackupManager) Stop() {
	if bm.cronEngine != nil {
		bm.cronEngine.Stop()
	}
}

func (bm *BackupManager) RegisterBackup(cfg *models.BackupConfig) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	return bm.registerBackupLocked(cfg)
}

func (bm *BackupManager) registerBackupLocked(cfg *models.BackupConfig) error {
	active := cfg.Status == models.BackupConfigStatusActive && cfg.BackupEnabled && strings.TrimSpace(cfg.Schedule) != "manual"
	var schedule cron.Schedule
	if active {
		var err error
		schedule, err = ParseSchedule(cfg.Schedule, cfg.Timezone)
		if err != nil {
			return fmt.Errorf("schedule backup %s: %w", cfg.Name, err)
		}
	}
	if entryID, exists := bm.entries[cfg.ID]; exists {
		bm.cronEngine.Remove(entryID)
		delete(bm.entries, cfg.ID)
	}
	if !active {
		return nil
	}
	cfgID := cfg.ID
	timeout := time.Duration(cfg.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Minute
	}
	bm.entries[cfg.ID] = bm.cronEngine.Schedule(schedule, cron.FuncJob(func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		bm.mu.Lock()
		runner := bm.scheduledRunner
		bm.mu.Unlock()
		var runErr error
		if runner != nil {
			runErr = runner(ctx, cfgID)
		} else {
			_, runErr = bm.TriggerBackup(ctx, cfgID)
		}
		if err := runErr; err != nil {
			slog.Error("scheduled backup failed", "backup_id", cfgID, "error", err)
		}
	}))
	return nil
}

func (bm *BackupManager) UnregisterBackup(backupConfigID string) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	if entryID, exists := bm.entries[backupConfigID]; exists {
		bm.cronEngine.Remove(entryID)
		delete(bm.entries, backupConfigID)
	}
}

func (bm *BackupManager) failBackupRecord(recID, errStr string) (*models.BackupRecord, error) {
	return bm.failBackupWithLogs(recID, "", errStr)
}

func (bm *BackupManager) failBackupWithLogs(recID, priorLogs, errStr string) (*models.BackupRecord, error) {
	logs := priorLogs + fmt.Sprintf("Failed: %s\n", errStr)
	if err := bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
		ID:          recID,
		Status:      models.BackupRecordStatusFailed,
		Logs:        logs,
		CompletedAt: time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		slog.Warn("failed to update backup record", "error", err)
	}
	return nil, errors.New(errStr)
}

func (bm *BackupManager) DeleteBackupRecord(ctx context.Context, recordID string) error {
	record, err := bm.store.GetBackupRecord(recordID)
	if err != nil {
		return err
	}
	cfg, err := bm.store.GetBackupConfig(record.BackupConfigID)
	if err != nil {
		return err
	}
	return bm.removeArchive(ctx, cfg, record)
}

func (bm *BackupManager) TriggerBackup(ctx context.Context, backupConfigID string) (*models.BackupRecord, error) {
	cfg, err := bm.store.GetBackupConfig(backupConfigID)
	if err != nil || cfg == nil {
		return nil, fmt.Errorf("backup config %s not found: %w", backupConfigID, err)
	}
	if !cfg.BackupEnabled {
		return nil, fmt.Errorf("backup execution disabled for config %s", backupConfigID)
	}

	rec := &models.BackupRecord{
		ID:              uuid.New().String(),
		BackupConfigID:  cfg.ID,
		DatabaseID:      cfg.DatabaseID,
		S3DestinationID: cfg.S3DestinationID,
		Status:          models.BackupRecordStatusRunning,
		Logs:            fmt.Sprintf("Initiating automated backup '%s' at %s...\n", cfg.Name, time.Now().UTC().Format(time.RFC3339)),
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

	if cfg.VolumeName != "" {
		dumpBytes, execLogs, err = bm.executeVolumeBackup(ctx, cfg.VolumeName)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
		fileExt = ".tar.gz"
	} else if cfg.DatabaseID == "global" || cfg.DatabaseID == "" {
		dbPath := filepath.Join(utils.GetDataDir(), "codedock.db")
		content, err := os.ReadFile(dbPath)
		if err != nil {
			return bm.failBackupRecord(rec.ID, fmt.Sprintf("failed to read global db: %v", err))
		}
		dumpBytes = content
		fileExt = ".db"
		execLogs = "Global database backed up successfully.\n"
	} else {
		containerName, dumpCmd, ext, err := bm.buildDumpCommand(cfg)
		if err != nil {
			return bm.failBackupRecord(rec.ID, err.Error())
		}
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
			if err := bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
				ID:            rec.ID,
				Status:        models.BackupRecordStatusFailed,
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

	if cfg.DisableLocal && s3URL != "" {
		_ = os.Remove(filePath)
		filePath = ""
	}

	completed, err := bm.finalizeBackupRecord(FinalizeBackupOpts{
		Record:    rec,
		FilePath:  filePath,
		S3URL:     s3URL,
		ExecLogs:  execLogs,
		SizeBytes: sizeBytes,
	})
	if err != nil {
		return nil, err
	}
	bm.enforceRetentionPolicy(cfg)
	return completed, nil
}

func (bm *BackupManager) SetScheduledRunner(runner func(context.Context, string) error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.scheduledRunner = runner
}
