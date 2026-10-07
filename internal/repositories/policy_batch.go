package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PolicyBatchRepo struct {
	db *sqlx.DB
}

func NewPolicyBatchRepo(db *sql.DB) *PolicyBatchRepo {
	return &PolicyBatchRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *PolicyBatchRepo) Create(ctx context.Context, batch *models.BackupPolicyBatch) error {
	if batch.ID == "" {
		batch.ID = uuid.NewString()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	batch.CreatedAt = now
	batch.UpdatedAt = now
	if batch.Status == "" {
		batch.Status = "active"
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO backup_policy_batches(id,project_id,name,description,schedule,timezone,timeout,status,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?,?)`, batch.ID, batch.ProjectID, batch.Name, batch.Description, batch.Schedule, batch.Timezone, batch.Timeout, batch.Status, batch.CreatedAt, batch.UpdatedAt)
	return err
}

func (r *PolicyBatchRepo) Get(ctx context.Context, id string) (*models.BackupPolicyBatch, error) {
	var batch models.BackupPolicyBatch
	if err := r.db.GetContext(ctx, &batch, `SELECT id,project_id,name,COALESCE(description,'') AS description,schedule,COALESCE(timezone,'UTC') AS timezone,timeout,status,created_at,updated_at FROM backup_policy_batches WHERE id=?`, id); err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *PolicyBatchRepo) ListByProject(ctx context.Context, projectID string) ([]models.BackupPolicyBatch, error) {
	result := []models.BackupPolicyBatch{}
	if err := r.db.SelectContext(ctx, &result, `SELECT id,project_id,name,COALESCE(description,'') AS description,schedule,COALESCE(timezone,'UTC') AS timezone,timeout,status,created_at,updated_at FROM backup_policy_batches WHERE project_id=? ORDER BY name`, projectID); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *PolicyBatchRepo) ListMembers(ctx context.Context, batchID string) ([]*models.BackupConfig, error) {
	list := []*models.BackupConfig{}
	if err := r.db.SelectContext(ctx, &list, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, '' as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, owner_id, project_id FROM backup_configs WHERE parent_batch_id=? ORDER BY name`, batchID); err != nil {
		return nil, err
	}
	if list == nil {
		list = []*models.BackupConfig{}
	}
	return list, nil
}
