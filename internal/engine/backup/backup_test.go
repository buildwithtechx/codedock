package backup

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

type mockStore struct {
	records map[string]*models.BackupRecord
	configs map[string]*models.BackupConfig
	s3dests map[string]*models.S3Destination
}

func newMockStore() *mockStore {
	return &mockStore{
		records: make(map[string]*models.BackupRecord),
		configs: make(map[string]*models.BackupConfig),
		s3dests: make(map[string]*models.S3Destination),
	}
}

func (m *mockStore) ListAllActiveBackupConfigs() ([]*models.BackupConfig, error) {
	var list []*models.BackupConfig
	for _, c := range m.configs {
		list = append(list, c)
	}
	return list, nil
}

func (m *mockStore) GetBackupConfig(id string) (*models.BackupConfig, error) {
	if c, ok := m.configs[id]; ok {
		return c, nil
	}
	return nil, nil
}

func (m *mockStore) CreateBackupRecord(rec *models.BackupRecord) error {
	m.records[rec.ID] = rec
	return nil
}

func (m *mockStore) GetDatabase(id string) (*models.Database, error) {
	return nil, nil
}

func (m *mockStore) UpdateBackupRecord(opts models.UpdateBackupRecordOpts) error {
	if rec, ok := m.records[opts.ID]; ok {
		rec.Status = opts.Status
		rec.FilePath = opts.FilePath
		rec.S3URL = opts.S3URL
		rec.Logs = opts.Logs
		rec.CompletedAt = opts.CompletedAt
	}
	return nil
}

func (m *mockStore) GetS3Destination(id string) (*models.S3Destination, error) {
	if d, ok := m.s3dests[id]; ok {
		return d, nil
	}
	return nil, nil
}

func (m *mockStore) GetBackupRecord(id string) (*models.BackupRecord, error) {
	if r, ok := m.records[id]; ok {
		return r, nil
	}
	return nil, nil
}

func (m *mockStore) ListBackupRecords(backupConfigID string) ([]*models.BackupRecord, error) {
	var list []*models.BackupRecord
	for _, r := range m.records {
		if r.BackupConfigID == backupConfigID {
			list = append(list, r)
		}
	}
	return list, nil
}

func TestMaxStorageGBOverflow(t *testing.T) {
	store := newMockStore()
	dir := t.TempDir()
	bm := NewBackupManager(nil, store, dir)

	cfg := &models.BackupConfig{
		ID:           "cfg-1",
		MaxStorageGB: 9000000000,
	}

	rec := &models.BackupRecord{
		ID:             "rec-1",
		BackupConfigID: cfg.ID,
		Status:         models.BackupRecordStatusCompleted,
		FileSizeBytes:  1024,
		StartedAt:      time.Now().UTC().Format(time.RFC3339),
	}
	store.records[rec.ID] = rec

	bm.enforceRetentionPolicy(cfg)

	if rec.Status == models.BackupRecordStatusExpired {
		t.Fatalf("expected record to stay active when storage limit is not exceeded")
	}
}

func TestDisableLocalGuard(t *testing.T) {
	store := newMockStore()
	dir := t.TempDir()
	bm := NewBackupManager(nil, store, dir)

	dataDir := utils.GetDataDir()
	_ = os.MkdirAll(dataDir, 0o755)
	dbFile := filepath.Join(dataDir, "codedock.db")
	_ = os.WriteFile(dbFile, []byte("test db content"), 0o600)
	defer os.Remove(dbFile)

	cfg := &models.BackupConfig{
		ID:            "cfg-disable-local",
		BackupEnabled: true,
		DisableLocal:  true,
		S3Enabled:     false,
	}
	store.configs[cfg.ID] = cfg

	rec, err := bm.TriggerBackup(context.Background(), cfg.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.FilePath == "" {
		t.Fatalf("expected local file to be retained when S3 upload was not performed")
	}
	if _, err := os.Stat(rec.FilePath); os.IsNotExist(err) {
		t.Fatalf("expected local file to exist on disk when S3 upload was not performed")
	}
}

func TestEnforceRetentionPolicyS3Cleanup(t *testing.T) {
	store := newMockStore()
	dir := t.TempDir()
	bm := NewBackupManager(nil, store, dir)

	cfg := &models.BackupConfig{
		ID:              "cfg-retention",
		RetentionDays:   1,
		S3DestinationID: "dest-1",
	}

	dest := &models.S3Destination{
		ID:              "dest-1",
		Bucket:          "mybucket",
		Endpoint:        "s3.amazonaws.com",
		AccessKeyID:     "key",
		SecretAccessKey: "secret",
	}
	store.s3dests[dest.ID] = dest

	oldTime := time.Now().Add(-48 * time.Hour).Format(time.RFC3339)
	dummyFile := filepath.Join(dir, "old_backup.db")
	_ = os.WriteFile(dummyFile, []byte("data"), 0o600)

	rec := &models.BackupRecord{
		ID:             "rec-old",
		BackupConfigID: cfg.ID,
		Status:         models.BackupRecordStatusCompleted,
		FilePath:       dummyFile,
		S3URL:          "s3://mybucket/backup_old.db",
		StartedAt:      oldTime,
	}
	store.records[rec.ID] = rec

	bm.enforceRetentionPolicy(cfg)

	if rec.Status != models.BackupRecordStatusExpired {
		t.Fatalf("expected record status to be expired, got %s", rec.Status)
	}
	if _, err := os.Stat(dummyFile); !os.IsNotExist(err) {
		t.Fatalf("expected local file to be removed")
	}
}
