package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"codedock/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
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
	db             *sqlx.DB
	mu             sync.Mutex
	vault          Vault
	snapshotDumper func(ctx context.Context) ([]byte, error)
}

func NewBackupRepo(db *sql.DB, v Vault) *BackupRepo {
	return &BackupRepo{db: sqlx.NewDb(db, "pgx"), vault: v}
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33)`,
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
		FROM backup_configs WHERE id = $1`, id)
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
		res, err := r.db.ExecContext(ctx, `UPDATE backup_configs SET owner_id=$1,project_id=$2,pre_deployment=$3, database_id=$4, service_id=$5, volume_name=$6, s3_destination_id=$7, sftp_destination_id=$8, parent_batch_id=$9, name=$10, description=$11, db_user=$12, backup_enabled=$13, s3_enabled=$14, sftp_enabled=$15, incremental=$16, disable_local=$17, quiesce_command=$18, unquiesce_command=$19, custom_backup_command=$20, custom_restore_command=$21, file_source_path=$22, schedule=$23, timezone=$24, timeout=$25, retention_days=$26, max_backups=$27, max_storage_gb=$28, status=$29, updated_at=$30 WHERE id=$31`,
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

	res, err := r.db.ExecContext(ctx, `UPDATE backup_configs SET owner_id=$1,project_id=$2,pre_deployment=$3, database_id=$4, service_id=$5, volume_name=$6, s3_destination_id=$7, sftp_destination_id=$8, parent_batch_id=$9, name=$10, description=$11, db_user=$12, db_password=$13, backup_enabled=$14, s3_enabled=$15, sftp_enabled=$16, incremental=$17, disable_local=$18, quiesce_command=$19, unquiesce_command=$20, custom_backup_command=$21, custom_restore_command=$22, file_source_path=$23, schedule=$24, timezone=$25, timeout=$26, retention_days=$27, max_backups=$28, max_storage_gb=$29, status=$30, updated_at=$31 WHERE id=$32`,
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
		FROM backup_configs WHERE database_id = $1 ORDER BY created_at DESC LIMIT 1`, dbID)
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
