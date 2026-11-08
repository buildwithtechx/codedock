package backups

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type ConfigRepository interface {
	CreateConfig(ctx context.Context, cfg *BackupConfig) error
	UpdateConfig(ctx context.Context, cfg *BackupConfig) error
	GetConfigByID(ctx context.Context, id string) (*BackupConfig, error)
	GetConfigByDatabaseID(ctx context.Context, dbID string) (*BackupConfig, error)
	ListConfigs(ctx context.Context) ([]*BackupConfig, error)
	ListAllActiveConfigs(ctx context.Context) ([]*BackupConfig, error)
	DeleteConfig(ctx context.Context, id string) error
	CreateRecord(ctx context.Context, rec *BackupRecord) error
	GetRecordByID(ctx context.Context, id string) (*BackupRecord, error)
	ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*BackupRecord, error)
	ListRecordsByDatabase(ctx context.Context, databaseID string) ([]*BackupRecord, error)
	ListAllRecords(ctx context.Context, limit int) ([]*BackupRecord, error)
	ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*BackupRecord, error)
	UpdateRecord(ctx context.Context, rec *BackupRecord) error
	DeleteRecord(ctx context.Context, id string) error
}

type S3Repository interface {
	CreateS3Destination(ctx context.Context, dest *S3Destination) error
	UpdateS3Destination(ctx context.Context, dest *S3Destination) error
	ListS3Destinations(ctx context.Context) ([]*S3Destination, error)
	GetS3Destination(ctx context.Context, id string) (*S3Destination, error)
	DeleteS3Destination(ctx context.Context, id string) error
	SetDefaultDestination(ctx context.Context, id string) error
	RecordVerificationResult(ctx context.Context, id string, verified bool, errMsg string) error
}

type BatchRepository interface {
	Create(ctx context.Context, batch *BackupPolicyBatch) error
	Get(ctx context.Context, id string) (*BackupPolicyBatch, error)
	ListByProject(ctx context.Context, projectID string) ([]BackupPolicyBatch, error)
	ListMembers(ctx context.Context, batchID string) ([]*BackupConfig, error)
}

type SFTPRepository interface {
	Create(ctx context.Context, dest *SFTPDestination) error
	Get(ctx context.Context, id string) (*SFTPDestination, error)
	ListByProject(ctx context.Context, projectID string) ([]SFTPDestination, error)
	Delete(ctx context.Context, id string) error
}

type Repository interface {
	ConfigRepository
	S3Repository
	BatchRepository
	ProtectRecord(ctx context.Context, id string, until int64) error
	ClearRestoreProtection(ctx context.Context, id string, before int64) error
	ClaimRecordExpiry(ctx context.Context, id string) (bool, error)
	SetSnapshotDumper(dumper func(ctx context.Context) ([]byte, error))
	ControlPlaneSnapshot(ctx context.Context) ([]byte, error)
}

type repo struct {
	db             *sqlx.DB
	vault          Vault
	mu             sync.Mutex
	snapshotDumper func(ctx context.Context) ([]byte, error)
}

func NewRepository(db *sql.DB, v Vault) Repository {
	return &repo{db: sqlx.NewDb(db, "pgx"), vault: v}
}

func (r *repo) CreateConfig(ctx context.Context, cfg *BackupConfig) error {
	if cfg.ID == "" {
		cfg.ID = uuid.New().String()
	}
	if cfg.CreatedAt == "" {
		cfg.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	cfg.UpdatedAt = cfg.CreatedAt
	if cfg.Status == "" {
		cfg.Status = ConfigStatusActive
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33)`,
		cfg.OwnerID, cfg.ProjectID, cfg.PreDeployment, cfg.ID, nullableID(cfg.DatabaseID), nullableID(cfg.ServiceID), cfg.VolumeName, nullableID(cfg.S3DestinationID), nullableID(cfg.SFTPDestinationID), nullableID(cfg.ParentBatchID), cfg.Name, cfg.Description, cfg.DbUser, cfg.DbPassword, cfg.BackupEnabled, cfg.S3Enabled, cfg.SFTPEnabled, cfg.Incremental, cfg.DisableLocal, cfg.QuiesceCommand, cfg.UnquiesceCommand, cfg.CustomBackupCommand, cfg.CustomRestoreCommand, cfg.FileSourcePath, cfg.Schedule, cfg.Timezone, cfg.Timeout, cfg.RetentionDays, cfg.MaxBackups, cfg.MaxStorageGB, string(cfg.Status), cfg.CreatedAt, cfg.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup config: %w", err)
	}
	cfg.DbPassword = "********"
	return nil
}

func (r *repo) UpdateConfig(ctx context.Context, cfg *BackupConfig) error {
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if cfg.Timeout == 0 {
		cfg.Timeout = 3600
	}
	if cfg.RetentionDays < 0 {
		cfg.RetentionDays = 0
	}
	pass := cfg.DbPassword
	if pass != "" && pass != "********" && r.vault != nil {
		enc, err := r.vault.Encrypt(pass)
		if err != nil {
			return fmt.Errorf("failed to encrypt db password: %w", err)
		}
		pass = enc
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `UPDATE backup_configs SET name=$1, description=$2, db_user=$3, db_password=CASE WHEN $4='********' THEN db_password ELSE $4 END, backup_enabled=$5, s3_enabled=$6, sftp_enabled=$7, incremental=$8, disable_local=$9, quiesce_command=$10, unquiesce_command=$11, custom_backup_command=$12, custom_restore_command=$13, file_source_path=$14, schedule=$15, timezone=$16, timeout=$17, retention_days=$18, max_backups=$19, max_storage_gb=$20, status=$21, updated_at=$22, s3_destination_id=$23, sftp_destination_id=$24, pre_deployment=$25, owner_id=$26 WHERE id=$27`,
		cfg.Name, cfg.Description, cfg.DbUser, pass, cfg.BackupEnabled, cfg.S3Enabled, cfg.SFTPEnabled, cfg.Incremental, cfg.DisableLocal, cfg.QuiesceCommand, cfg.UnquiesceCommand, cfg.CustomBackupCommand, cfg.CustomRestoreCommand, cfg.FileSourcePath, cfg.Schedule, cfg.Timezone, cfg.Timeout, cfg.RetentionDays, cfg.MaxBackups, cfg.MaxStorageGB, string(cfg.Status), cfg.UpdatedAt, nullableID(cfg.S3DestinationID), nullableID(cfg.SFTPDestinationID), cfg.PreDeployment, cfg.OwnerID, cfg.ID)
	if err != nil {
		return fmt.Errorf("failed to update backup config: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrBackupConfigNotFound
	}
	cfg.DbPassword = "********"
	return nil
}

func (r *repo) GetConfigByID(ctx context.Context, id string) (*BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg BackupConfig
	err := r.db.GetContext(ctx, &cfg, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, COALESCE(owner_id, '') as owner_id, COALESCE(project_id, '') as project_id
		FROM backup_configs WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBackupConfigNotFound
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

func (r *repo) GetConfigByDatabaseID(ctx context.Context, dbID string) (*BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg BackupConfig
	err := r.db.GetContext(ctx, &cfg, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, COALESCE(owner_id, '') as owner_id, COALESCE(project_id, '') as project_id
		FROM backup_configs WHERE database_id = $1 ORDER BY created_at DESC LIMIT 1`, dbID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBackupConfigNotFound
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

func (r *repo) ListConfigs(ctx context.Context) ([]*BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*BackupConfig
	err := r.db.SelectContext(ctx, &list, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, COALESCE(owner_id, '') as owner_id, COALESCE(project_id, '') as project_id
		FROM backup_configs ORDER BY created_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup configs: %w", err)
	}
	if list == nil {
		list = make([]*BackupConfig, 0)
	}
	return list, nil
}

func (r *repo) ListAllActiveConfigs(ctx context.Context) ([]*BackupConfig, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*BackupConfig
	err := r.db.SelectContext(ctx, &list, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, COALESCE(db_password, '') as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, COALESCE(owner_id, '') as owner_id, COALESCE(project_id, '') as project_id
		FROM backup_configs WHERE status = 'active'`)
	if err != nil {
		return nil, fmt.Errorf("failed to list active backup configs: %w", err)
	}
	if list == nil {
		list = make([]*BackupConfig, 0)
	}
	return list, nil
}

func (r *repo) DeleteConfig(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var protected int
	if err := r.db.GetContext(ctx, &protected, `SELECT COUNT(*) FROM backup_records WHERE backup_config_id=$1 AND (protected_until>$2 OR status IN ('expiring','running'))`, id, time.Now().Unix()); err != nil {
		return err
	}
	if protected > 0 {
		return fmt.Errorf("policy contains protected or active records; finish runs and remove protection before deleting")
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM backup_configs WHERE id = $1 AND NOT EXISTS(SELECT 1 FROM backup_records WHERE backup_config_id=$2 AND (protected_until>$3 OR status IN ('expiring','running')))`, id, id, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("failed to delete backup config: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return ErrBackupConfigNotFound
	}
	return nil
}

func nullableID(id string) any {
	if id == "" {
		return nil
	}
	return id
}
