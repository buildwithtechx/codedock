package backup

import (
	"context"
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
	dockerClient *client.Client
	store        Store
	cronEngine   *cron.Cron
	entries      map[string]cron.EntryID
	backupDir    string
	mu           sync.Mutex
}

func NewBackupManager(dockerClient *client.Client, s Store, backupDir string) *BackupManager {
	if backupDir == "" {
		backupDir = filepath.Join(utils.GetDataDir(), "backups")
	}
	_ = os.MkdirAll(backupDir, 0o700)
	return &BackupManager{
		dockerClient: dockerClient,
		store:        s,
		cronEngine:   cron.New(cron.WithSeconds()),
		entries:      make(map[string]cron.EntryID),
		backupDir:    backupDir,
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
	if entryID, exists := bm.entries[cfg.ID]; exists {
		bm.cronEngine.Remove(entryID)
		delete(bm.entries, cfg.ID)
	}
	if cfg.Status != "active" || !cfg.BackupEnabled {
		return nil
	}
	schedule := strings.TrimSpace(cfg.Schedule)
	if len(strings.Fields(schedule)) == 5 && !strings.HasPrefix(schedule, "@") {
		schedule = "0 " + schedule
	}
	cfgID := cfg.ID
	entryID, err := bm.cronEngine.AddFunc(schedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()
		_, _ = bm.TriggerBackup(ctx, cfgID)
	})
	if err != nil {
		return fmt.Errorf("invalid cron schedule '%s' for backup %s: %w", cfg.Schedule, cfg.Name, err)
	}
	bm.entries[cfg.ID] = entryID
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

func (bm *BackupManager) DeleteBackupRecord(ctx context.Context, recordID string) {
	rec, err := bm.store.GetBackupRecord(recordID)
	if err == nil && rec != nil && rec.FilePath != "" {
		_ = os.Remove(rec.FilePath)
	}
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

	fileName := fmt.Sprintf("backup_%s_%s%s", cfg.ID, time.Now().UTC().Format("20060102_150405"), fileExt)
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
		var s3Err error
		s3URL, execLogs, s3Err = bm.handleS3Upload(ctx, cfg, fileName, dumpBytes, execLogs)
		if s3Err != nil {
			_ = bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
				ID:            rec.ID,
				Status:        models.BackupRecordStatusFailed,
				FilePath:      filePath,
				FileSizeBytes: sizeBytes,
				Logs:          execLogs + fmt.Sprintf("Failed: %s\n", s3Err.Error()),
				CompletedAt:   time.Now().UTC().Format(time.RFC3339),
			})
			return nil, s3Err
		}
	}

	if cfg.DisableLocal && s3URL != "" {
		_ = os.Remove(filePath)
		filePath = ""
	}

	bm.enforceRetentionPolicy(cfg)

	return bm.finalizeBackupRecord(FinalizeBackupOpts{
		Record:    rec,
		FilePath:  filePath,
		S3URL:     s3URL,
		ExecLogs:  execLogs,
		SizeBytes: sizeBytes,
	})
}

type FinalizeBackupOpts struct {
	Record    *models.BackupRecord
	FilePath  string
	S3URL     string
	ExecLogs  string
	SizeBytes int64
}

func (bm *BackupManager) finalizeBackupRecord(opts FinalizeBackupOpts) (*models.BackupRecord, error) {
	nowStr := time.Now().UTC().Format(time.RFC3339)
	finalLogs := opts.Record.Logs + opts.ExecLogs + "\nBackup run completed successfully."
	if err := bm.store.UpdateBackupRecord(models.UpdateBackupRecordOpts{
		ID:            opts.Record.ID,
		Status:        models.BackupRecordStatusCompleted,
		FilePath:      opts.FilePath,
		S3URL:         opts.S3URL,
		Logs:          finalLogs,
		FileSizeBytes: opts.SizeBytes,
		CompletedAt:   nowStr,
	}); err != nil {
		slog.Warn("failed to update backup record", "error", err)
	}

	opts.Record.Status = models.BackupRecordStatusCompleted
	opts.Record.FilePath = opts.FilePath
	opts.Record.FileSizeBytes = opts.SizeBytes
	opts.Record.S3URL = opts.S3URL
	opts.Record.Logs = finalLogs
	opts.Record.CompletedAt = nowStr
	return opts.Record, nil
}
