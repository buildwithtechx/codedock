package backups

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockBackupConfigRepo struct {
	configs map[string]*BackupConfig
	records map[string]*BackupRecord
}

func newMockBackupConfigRepo() *mockBackupConfigRepo {
	return &mockBackupConfigRepo{
		configs: make(map[string]*BackupConfig),
		records: make(map[string]*BackupRecord),
	}
}

func (m *mockBackupConfigRepo) CreateConfig(ctx context.Context, cfg *BackupConfig) error {
	m.configs[cfg.ID] = cfg
	return nil
}

func (m *mockBackupConfigRepo) UpdateConfig(ctx context.Context, cfg *BackupConfig) error {
	m.configs[cfg.ID] = cfg
	return nil
}

func (m *mockBackupConfigRepo) GetConfigByID(ctx context.Context, id string) (*BackupConfig, error) {
	cfg, ok := m.configs[id]
	if !ok {
		return nil, ErrBackupConfigNotFound
	}
	cp := *cfg
	return &cp, nil
}

func (m *mockBackupConfigRepo) GetConfigByDatabaseID(ctx context.Context, dbID string) (*BackupConfig, error) {
	for _, cfg := range m.configs {
		if cfg.DatabaseID == dbID {
			cp := *cfg
			return &cp, nil
		}
	}
	return nil, ErrBackupConfigNotFound
}

func (m *mockBackupConfigRepo) ListConfigs(ctx context.Context) ([]*BackupConfig, error) {
	var list []*BackupConfig
	for _, cfg := range m.configs {
		cp := *cfg
		list = append(list, &cp)
	}
	return list, nil
}

func (m *mockBackupConfigRepo) ListAllActiveConfigs(ctx context.Context) ([]*BackupConfig, error) {
	var list []*BackupConfig
	for _, cfg := range m.configs {
		if cfg.Status == ConfigStatusActive {
			cp := *cfg
			list = append(list, &cp)
		}
	}
	return list, nil
}

func (m *mockBackupConfigRepo) DeleteConfig(ctx context.Context, id string) error {
	delete(m.configs, id)
	return nil
}

func (m *mockBackupConfigRepo) CreateRecord(ctx context.Context, rec *BackupRecord) error {
	m.records[rec.ID] = rec
	return nil
}

func (m *mockBackupConfigRepo) GetRecordByID(ctx context.Context, id string) (*BackupRecord, error) {
	rec, ok := m.records[id]
	if !ok {
		return nil, ErrBackupRecordNotFound
	}
	cp := *rec
	return &cp, nil
}

func (m *mockBackupConfigRepo) ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*BackupRecord, error) {
	var list []*BackupRecord
	for _, rec := range m.records {
		if rec.BackupConfigID == backupConfigID {
			cp := *rec
			list = append(list, &cp)
		}
	}
	return list, nil
}

func (m *mockBackupConfigRepo) ListRecordsByDatabase(ctx context.Context, databaseID string) ([]*BackupRecord, error) {
	var list []*BackupRecord
	for _, rec := range m.records {
		if rec.DatabaseID == databaseID {
			cp := *rec
			list = append(list, &cp)
		}
	}
	return list, nil
}

func (m *mockBackupConfigRepo) ListAllRecords(ctx context.Context, limit int) ([]*BackupRecord, error) {
	var list []*BackupRecord
	for _, rec := range m.records {
		cp := *rec
		list = append(list, &cp)
		if len(list) >= limit {
			break
		}
	}
	return list, nil
}

func (m *mockBackupConfigRepo) ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*BackupRecord, error) {
	match := make(map[string]bool)
	for _, id := range configIDs {
		match[id] = true
	}
	var list []*BackupRecord
	for _, rec := range m.records {
		if match[rec.BackupConfigID] {
			cp := *rec
			list = append(list, &cp)
			if len(list) >= limit {
				break
			}
		}
	}
	return list, nil
}

func (m *mockBackupConfigRepo) UpdateRecord(ctx context.Context, rec *BackupRecord) error {
	m.records[rec.ID] = rec
	return nil
}

func (m *mockBackupConfigRepo) DeleteRecord(ctx context.Context, id string) error {
	delete(m.records, id)
	return nil
}

func TestRedactBackupModels(t *testing.T) {
	assert.Nil(t, RedactBackupConfig(nil))
	assert.Nil(t, RedactS3Destination(nil))
	assert.Nil(t, RedactSFTPDestination(nil))

	cfg := &BackupConfig{
		ID:         "cfg-1",
		DbPassword: "super-secret-password",
	}
	RedactBackupConfig(cfg)
	assert.Equal(t, "********", cfg.DbPassword)

	s3 := &S3Destination{
		ID:              "s3-1",
		SecretAccessKey: "super-secret-s3-key",
	}
	RedactS3Destination(s3)
	assert.Equal(t, "********", s3.SecretAccessKey)

	sftp := &SFTPDestination{
		ID:         "sftp-1",
		Password:   "secret-pass",
		PrivateKey: "secret-key",
	}
	RedactSFTPDestination(sftp)
	assert.Equal(t, "********", sftp.Password)
	assert.Equal(t, "********", sftp.PrivateKey)
}

func TestBackupService_ConfigLifecycle(t *testing.T) {
	repo := newMockBackupConfigRepo()
	svc := NewService(repo, nil, nil)

	cfg := &BackupConfig{
		ID:            "b-1",
		Name:          "Nightly DB Backup",
		Schedule:      "0 2 * * *",
		BackupEnabled: true,
	}

	err := svc.CreateConfig(context.Background(), cfg)
	require.NoError(t, err)

	fetched, err := svc.GetConfig(context.Background(), "b-1")
	require.NoError(t, err)
	assert.Equal(t, "Nightly DB Backup", fetched.Name)
	assert.Equal(t, ConfigStatusActive, fetched.Status)

	fetched.Name = "Updated DB Backup"
	err = svc.UpdateConfig(context.Background(), fetched)
	require.NoError(t, err)

	updated, err := svc.GetConfig(context.Background(), "b-1")
	require.NoError(t, err)
	assert.Equal(t, "Updated DB Backup", updated.Name)

	err = svc.DeleteConfig(context.Background(), "b-1")
	require.NoError(t, err)

	_, err = svc.GetConfig(context.Background(), "b-1")
	assert.Error(t, err)
}
