package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"codedock.run/codedock/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock.run/codedock/internal/models"
)

type BackupRepository interface {
	CreateConfig(ctx context.Context, cfg *models.BackupConfig) error
	UpdateConfig(ctx context.Context, cfg *models.BackupConfig) error
	GetConfigByID(ctx context.Context, id string) (*models.BackupConfig, error)
	GetConfigByDatabaseID(ctx context.Context, dbID string) (*models.BackupConfig, error)
	ListConfigs(ctx context.Context) ([]*models.BackupConfig, error)
	ListAllActiveConfigs(ctx context.Context) ([]*models.BackupConfig, error)
	DeleteConfig(ctx context.Context, id string) error
	CreateRecord(ctx context.Context, rec *models.BackupRecord) error
	GetRecordByID(ctx context.Context, id string) (*models.BackupRecord, error)
	ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*models.BackupRecord, error)
	ListRecordsByDatabase(ctx context.Context, databaseID string) ([]*models.BackupRecord, error)
	ListAllRecords(ctx context.Context, limit int) ([]*models.BackupRecord, error)
	ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*models.BackupRecord, error)
	UpdateRecord(ctx context.Context, rec *models.BackupRecord) error
	DeleteRecord(ctx context.Context, id string) error
}

type BackupRepo struct {
	db    *sqlx.DB
	mu    sync.Mutex
	vault Vault
}

func NewBackupRepo(db *sql.DB, v Vault) *BackupRepo {
	return &BackupRepo{db: sqlx.NewDb(db, "sqlite"), vault: v}
}

func (r *BackupRepo) EnsureTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS backup_configs (
			id TEXT PRIMARY KEY,
			database_id TEXT,
			service_id TEXT,
			volume_name TEXT,
			s3_destination_id TEXT,
			sftp_destination_id TEXT,
			parent_batch_id TEXT,
			name TEXT NOT NULL,
			description TEXT,
			db_user TEXT,
			db_password TEXT,
			backup_enabled INTEGER DEFAULT 1,
			s3_enabled INTEGER DEFAULT 0,
			sftp_enabled INTEGER DEFAULT 0,
			incremental INTEGER DEFAULT 0,
			disable_local INTEGER DEFAULT 0,
			quiesce_command TEXT DEFAULT '',
			unquiesce_command TEXT DEFAULT '',
			custom_backup_command TEXT DEFAULT '',
			custom_restore_command TEXT DEFAULT '',
			file_source_path TEXT DEFAULT '',
			schedule TEXT NOT NULL,
			timezone TEXT DEFAULT 'UTC',
			timeout INTEGER DEFAULT 3600,
			retention_days INTEGER DEFAULT 7,
			max_backups INTEGER DEFAULT 0,
			max_storage_gb INTEGER DEFAULT 0,
			status TEXT DEFAULT 'active',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS backup_records (
			id TEXT PRIMARY KEY,
			backup_config_id TEXT NOT NULL,
			database_id TEXT,
			service_id TEXT,
			volume_name TEXT,
			s3_destination_id TEXT,
			sftp_destination_id TEXT,
			sftp_url TEXT DEFAULT '',
			parent_record_id TEXT DEFAULT '',
			status TEXT DEFAULT 'running',
			file_path TEXT,
			file_size_bytes INTEGER DEFAULT 0,
			s3_url TEXT,
			logs TEXT,
			started_at TEXT NOT NULL,
			completed_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS s3_destinations (
			id TEXT PRIMARY KEY,
			organization_id TEXT DEFAULT '',
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			provider TEXT DEFAULT 's3',
			endpoint TEXT NOT NULL,
			bucket TEXT NOT NULL,
			region TEXT,
			path_prefix TEXT DEFAULT '',
			is_default INTEGER DEFAULT 0,
			last_verified_at TEXT,
			last_verify_error TEXT DEFAULT '',
			access_key_id TEXT,
			secret_access_key TEXT,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sftp_destinations (
			id TEXT PRIMARY KEY,
			organization_id TEXT DEFAULT '',
			project_id TEXT DEFAULT '',
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			host TEXT NOT NULL,
			port INTEGER DEFAULT 22,
			username TEXT NOT NULL,
			password TEXT DEFAULT '',
			private_key TEXT DEFAULT '',
			path_prefix TEXT DEFAULT '',
			last_verified_at TEXT,
			last_verify_error TEXT DEFAULT '',
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS backup_policy_batches (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT DEFAULT '',
			schedule TEXT NOT NULL,
			timezone TEXT DEFAULT 'UTC',
			timeout INTEGER DEFAULT 3600,
			status TEXT DEFAULT 'active',
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
	}
	for _, q := range queries {
		if _, err := r.db.Exec(q); err != nil {
			return fmt.Errorf("failed to create backup table: %w", err)
		}
	}
	for _, alter := range []string{
		`ALTER TABLE backup_configs ADD COLUMN sftp_destination_id TEXT`,
		`ALTER TABLE backup_configs ADD COLUMN parent_batch_id TEXT`,
		`ALTER TABLE backup_configs ADD COLUMN sftp_enabled INTEGER DEFAULT 0`,
		`ALTER TABLE backup_configs ADD COLUMN incremental INTEGER DEFAULT 0`,
		`ALTER TABLE backup_configs ADD COLUMN quiesce_command TEXT DEFAULT ''`,
		`ALTER TABLE backup_configs ADD COLUMN unquiesce_command TEXT DEFAULT ''`,
		`ALTER TABLE backup_configs ADD COLUMN custom_backup_command TEXT DEFAULT ''`,
		`ALTER TABLE backup_configs ADD COLUMN custom_restore_command TEXT DEFAULT ''`,
		`ALTER TABLE backup_configs ADD COLUMN file_source_path TEXT DEFAULT ''`,
		`ALTER TABLE backup_records ADD COLUMN sftp_destination_id TEXT`,
		`ALTER TABLE backup_records ADD COLUMN sftp_url TEXT DEFAULT ''`,
		`ALTER TABLE backup_records ADD COLUMN parent_record_id TEXT DEFAULT ''`,
	} {
		_, _ = r.db.Exec(alter)
	}
	_, _ = r.db.Exec(`ALTER TABLE s3_destinations ADD COLUMN path_prefix TEXT DEFAULT ''`)
	_, _ = r.db.Exec(`ALTER TABLE s3_destinations ADD COLUMN is_default INTEGER DEFAULT 0`)
	_, _ = r.db.Exec(`ALTER TABLE s3_destinations ADD COLUMN last_verified_at TEXT`)
	_, _ = r.db.Exec(`ALTER TABLE s3_destinations ADD COLUMN last_verify_error TEXT DEFAULT ''`)
	return r.ensureRecoveryColumns()
}

func (r *BackupRepo) CreateConfig(ctx context.Context, cfg *models.BackupConfig) error {
	if cfg.ID == "" {
		cfg.ID = uuid.New().String()
	}
	if cfg.CreatedAt == "" {
		cfg.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	cfg.UpdatedAt = cfg.CreatedAt
	if cfg.Status == "" {
		cfg.Status = "active"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 3600
	}
	if cfg.RetentionDays < 0 {
		cfg.RetentionDays = 0
	}
	if cfg.DbPassword != "" && r.vault != nil {
		enc, err := r.vault.Encrypt(cfg.DbPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt db password: %w", err)
		}
		cfg.DbPassword = enc
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO backup_configs (owner_id,project_id,pre_deployment,id, database_id, service_id, volume_name, s3_destination_id, sftp_destination_id, parent_batch_id, name, description, db_user, db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, quiesce_command, unquiesce_command, custom_backup_command, custom_restore_command, file_source_path, schedule, timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		cfg.OwnerID, cfg.ProjectID, cfg.PreDeployment, cfg.ID, nullableID(cfg.DatabaseID), nullableID(cfg.ServiceID), cfg.VolumeName, nullableID(cfg.S3DestinationID), nullableID(cfg.SFTPDestinationID), nullableID(cfg.ParentBatchID), cfg.Name, cfg.Description, cfg.DbUser, cfg.DbPassword, cfg.BackupEnabled, cfg.S3Enabled, cfg.SFTPEnabled, cfg.Incremental, cfg.DisableLocal, cfg.QuiesceCommand, cfg.UnquiesceCommand, cfg.CustomBackupCommand, cfg.CustomRestoreCommand, cfg.FileSourcePath, cfg.Schedule, cfg.Timezone, cfg.Timeout, cfg.RetentionDays, cfg.MaxBackups, cfg.MaxStorageGB, cfg.Status, cfg.CreatedAt, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup config: %w", err)
	}
	cfg.DbPassword = "********"
	return nil
}

func (r *BackupRepo) GetConfigByID(ctx context.Context, id string) (*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.BackupConfig
	err := r.db.GetContext(ctx, &cfg, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, owner_id, project_id
		FROM backup_configs WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Config", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get backup config %s: %w", id, err)
	}
	if cfg.DbPassword != "" && r.vault != nil {
		dec, err := r.vault.Decrypt(cfg.DbPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt db password: %w", err)
		}
		cfg.DbPassword = dec
	}
	return &cfg, nil
}

func (r *BackupRepo) UpdateConfig(ctx context.Context, cfg *models.BackupConfig) error {
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)

	if cfg.DbPassword != "" && cfg.DbPassword != "********" && r.vault != nil {
		enc, err := r.vault.Encrypt(cfg.DbPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt db password: %w", err)
		}
		cfg.DbPassword = enc
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if cfg.DbPassword == "********" || cfg.DbPassword == "" {
		res, err := r.db.ExecContext(ctx, `UPDATE backup_configs SET owner_id=?,project_id=?,pre_deployment=?, database_id=?, service_id=?, volume_name=?, s3_destination_id=?, sftp_destination_id=?, parent_batch_id=?, name=?, description=?, db_user=?, backup_enabled=?, s3_enabled=?, sftp_enabled=?, incremental=?, disable_local=?, quiesce_command=?, unquiesce_command=?, custom_backup_command=?, custom_restore_command=?, file_source_path=?, schedule=?, timezone=?, timeout=?, retention_days=?, max_backups=?, max_storage_gb=?, status=?, updated_at=? WHERE id=?`,
			cfg.OwnerID, cfg.ProjectID, cfg.PreDeployment, nullableID(cfg.DatabaseID), nullableID(cfg.ServiceID), cfg.VolumeName, nullableID(cfg.S3DestinationID), nullableID(cfg.SFTPDestinationID), nullableID(cfg.ParentBatchID), cfg.Name, cfg.Description, cfg.DbUser, cfg.BackupEnabled, cfg.S3Enabled, cfg.SFTPEnabled, cfg.Incremental, cfg.DisableLocal, cfg.QuiesceCommand, cfg.UnquiesceCommand, cfg.CustomBackupCommand, cfg.CustomRestoreCommand, cfg.FileSourcePath, cfg.Schedule, cfg.Timezone, cfg.Timeout, cfg.RetentionDays, cfg.MaxBackups, cfg.MaxStorageGB, cfg.Status, cfg.UpdatedAt, cfg.ID)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}
		if affected == 0 {
			return utils.NewNotFoundError("BackupConfig", cfg.ID)
		}
		return nil
	}

	res, err := r.db.ExecContext(ctx, `UPDATE backup_configs SET owner_id=?,project_id=?,pre_deployment=?, database_id=?, service_id=?, volume_name=?, s3_destination_id=?, sftp_destination_id=?, parent_batch_id=?, name=?, description=?, db_user=?, db_password=?, backup_enabled=?, s3_enabled=?, sftp_enabled=?, incremental=?, disable_local=?, quiesce_command=?, unquiesce_command=?, custom_backup_command=?, custom_restore_command=?, file_source_path=?, schedule=?, timezone=?, timeout=?, retention_days=?, max_backups=?, max_storage_gb=?, status=?, updated_at=? WHERE id=?`,
		cfg.OwnerID, cfg.ProjectID, cfg.PreDeployment, nullableID(cfg.DatabaseID), nullableID(cfg.ServiceID), cfg.VolumeName, nullableID(cfg.S3DestinationID), nullableID(cfg.SFTPDestinationID), nullableID(cfg.ParentBatchID), cfg.Name, cfg.Description, cfg.DbUser, cfg.DbPassword, cfg.BackupEnabled, cfg.S3Enabled, cfg.SFTPEnabled, cfg.Incremental, cfg.DisableLocal, cfg.QuiesceCommand, cfg.UnquiesceCommand, cfg.CustomBackupCommand, cfg.CustomRestoreCommand, cfg.FileSourcePath, cfg.Schedule, cfg.Timezone, cfg.Timeout, cfg.RetentionDays, cfg.MaxBackups, cfg.MaxStorageGB, cfg.Status, cfg.UpdatedAt, cfg.ID)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if affected == 0 {
		return utils.NewNotFoundError("BackupConfig", cfg.ID)
	}
	return nil
}

func (r *BackupRepo) ListConfigs(ctx context.Context) ([]*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.BackupConfig
	err := r.db.SelectContext(ctx, &list, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, owner_id, project_id
		FROM backup_configs ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup configs: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupConfig, 0)
	}
	if r.vault != nil {
		for _, cfg := range list {
			if cfg.DbPassword != "" {
				dec, err := r.vault.Decrypt(cfg.DbPassword)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt db password: %w", err)
				}
				cfg.DbPassword = dec
			}
		}
	}
	return list, nil
}

func (r *BackupRepo) ListAllActiveConfigs(ctx context.Context) ([]*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.BackupConfig
	err := r.db.SelectContext(ctx, &list, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, owner_id, project_id
		FROM backup_configs WHERE status = 'active'`)
	if err != nil {
		return nil, fmt.Errorf("failed to list active backup configs: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupConfig, 0)
	}
	if r.vault != nil {
		for _, cfg := range list {
			if cfg.DbPassword != "" {
				dec, err := r.vault.Decrypt(cfg.DbPassword)
				if err != nil {
					return nil, fmt.Errorf("failed to decrypt db password: %w", err)
				}
				cfg.DbPassword = dec
			}
		}
	}
	return list, nil
}

func (r *BackupRepo) GetConfigByDatabaseID(ctx context.Context, dbID string) (*models.BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.BackupConfig
	err := r.db.GetContext(ctx, &cfg, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, owner_id, project_id
		FROM backup_configs WHERE database_id = ? ORDER BY created_at DESC LIMIT 1`, dbID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Config", dbID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get backup config for database %s: %w", dbID, err)
	}
	if cfg.DbPassword != "" && r.vault != nil {
		dec, err := r.vault.Decrypt(cfg.DbPassword)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt db password: %w", err)
		}
		cfg.DbPassword = dec
	}
	return &cfg, nil
}

func (r *BackupRepo) DeleteConfig(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var protected int
	if err := r.db.GetContext(ctx, &protected, `SELECT COUNT(*) FROM backup_records WHERE backup_config_id=? AND (protected_until>? OR status IN ('expiring','running'))`, id, time.Now().Unix()); err != nil {
		return err
	}
	if protected > 0 {
		return fmt.Errorf("policy contains protected or active records; finish runs and remove protection before deleting")
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM backup_configs WHERE id = ? AND NOT EXISTS(SELECT 1 FROM backup_records WHERE backup_config_id=? AND (protected_until>? OR status IN ('expiring','running')))`, id, id, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to delete backup config: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return utils.NewNotFoundError("BackupConfig", id)
	}
	return nil
}

func nullableID(id string) any {
	if id == "" {
		return nil
	}
	return id
}
