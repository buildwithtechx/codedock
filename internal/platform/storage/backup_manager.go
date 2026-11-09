package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/robfig/cron/v3"

	"codedock/internal/utils"
)

type Store interface {
	ListAllActiveBackupConfigs() ([]*BackupConfig, error)
	GetBackupConfig(id string) (*BackupConfig, error)
	CreateBackupRecord(rec *BackupRecord) error
	GetDatabase(id string) (*Database, error)
	UpdateBackupRecord(opts UpdateBackupRecordOpts) error
	GetS3Destination(id string) (*S3Destination, error)
	GetBackupRecord(id string) (*BackupRecord, error)
	ListBackupRecords(backupConfigID string) ([]*BackupRecord, error)
}

type BackupManager struct {
	dockerTarget     DockerTarget
	serviceRuntime   func(context.Context, string) (string, error)
	nativeBackup     func(context.Context, string, io.Writer) error
	nativeRestore    func(ctx context.Context, serviceID string, archive io.Reader) error
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
	OnBackupFailed   func(recordID, errStr string)
}

func (bm *BackupManager) SetServiceRuntime(lookup func(context.Context, string) (string, error)) {
	bm.serviceRuntime = lookup
}

func (bm *BackupManager) SetNativeData(backup func(context.Context, string, io.Writer) error, restore func(context.Context, string, io.Reader) error) {
	bm.nativeBackup = backup
	bm.nativeRestore = restore
}

func (bm *BackupManager) isNativeService(ctx context.Context, serviceID string) bool {
	if serviceID == "" || bm.serviceRuntime == nil {
		return false
	}
	kind, err := bm.serviceRuntime(ctx, serviceID)
	return err == nil && kind == "bare"
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
	if err := bm.Reconcile(); err != nil {
		return err
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

func (bm *BackupManager) RegisterBackup(cfg *BackupConfig) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	return bm.registerBackupLocked(cfg)
}

func (bm *BackupManager) registerBackupLocked(cfg *BackupConfig) error {
	active := cfg.Status == BackupConfigStatusActive && cfg.BackupEnabled && strings.TrimSpace(cfg.Schedule) != "manual"
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

func (bm *BackupManager) failBackupRecord(recID, errStr string) (*BackupRecord, error) {
	return bm.failBackupWithLogs(recID, "", errStr)
}

func (bm *BackupManager) failBackupWithLogs(recID, priorLogs, errStr string) (*BackupRecord, error) {
	if bm.OnBackupFailed != nil {
		bm.OnBackupFailed(recID, errStr)
	}
	logs := priorLogs + fmt.Sprintf("Failed: %s\n", errStr)
	if err := bm.store.UpdateBackupRecord(UpdateBackupRecordOpts{
		ID:          recID,
		Status:      BackupRecordStatusFailed,
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

func (bm *BackupManager) SetScheduledRunner(runner func(context.Context, string) error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	bm.scheduledRunner = runner
}
