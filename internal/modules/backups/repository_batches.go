package backups

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *repo) Create(ctx context.Context, batch *BackupPolicyBatch) error {
	if batch.ID == "" {
		batch.ID = uuid.NewString()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	batch.CreatedAt = now
	batch.UpdatedAt = now
	if batch.Status == "" {
		batch.Status = "active"
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO backup_policy_batches(id,project_id,name,description,schedule,timezone,timeout,status,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, batch.ID, batch.ProjectID, batch.Name, batch.Description, batch.Schedule, batch.Timezone, batch.Timeout, batch.Status, batch.CreatedAt, batch.UpdatedAt)
	return err
}

func (r *repo) Get(ctx context.Context, id string) (*BackupPolicyBatch, error) {
	var batch BackupPolicyBatch
	if err := r.db.GetContext(ctx, &batch, `SELECT id,project_id,name,COALESCE(description,'') AS description,schedule,COALESCE(timezone,'UTC') AS timezone,timeout,status,created_at,updated_at FROM backup_policy_batches WHERE id=$1`, id); err != nil {
		return nil, err
	}
	return &batch, nil
}

func (r *repo) ListByProject(ctx context.Context, projectID string) ([]BackupPolicyBatch, error) {
	result := []BackupPolicyBatch{}
	if err := r.db.SelectContext(ctx, &result, `SELECT id,project_id,name,COALESCE(description,'') AS description,schedule,COALESCE(timezone,'UTC') AS timezone,timeout,status,created_at,updated_at FROM backup_policy_batches WHERE project_id=$1 ORDER BY name`, projectID); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repo) ListMembers(ctx context.Context, batchID string) ([]*BackupConfig, error) {
	list := []*BackupConfig{}
	if err := r.db.SelectContext(ctx, &list, `SELECT id, COALESCE(database_id, '') as database_id, COALESCE(service_id, '') as service_id, COALESCE(volume_name, '') as volume_name, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_batch_id, '') as parent_batch_id, name, COALESCE(description, '') as description, COALESCE(db_user, '') as db_user, '' as db_password, backup_enabled, s3_enabled, sftp_enabled, incremental, disable_local, COALESCE(quiesce_command, '') as quiesce_command, COALESCE(unquiesce_command, '') as unquiesce_command, COALESCE(custom_backup_command, '') as custom_backup_command, COALESCE(custom_restore_command, '') as custom_restore_command, COALESCE(file_source_path, '') as file_source_path, schedule, COALESCE(timezone, 'UTC') as timezone, timeout, retention_days, max_backups, max_storage_gb, status, created_at, updated_at, pre_deployment, COALESCE(owner_id, '') as owner_id, COALESCE(project_id, '') as project_id FROM backup_configs WHERE parent_batch_id=$1 ORDER BY name`, batchID); err != nil {
		return nil, err
	}
	if list == nil {
		list = []*BackupConfig{}
	}
	return list, nil
}

type sftpRepo struct {
	db *sqlx.DB
}

func NewSFTPRepository(db *sql.DB) SFTPRepository {
	return &sftpRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *sftpRepo) Create(ctx context.Context, dest *SFTPDestination) error {
	if dest.ID == "" {
		dest.ID = uuid.NewString()
	}
	if dest.Port <= 0 {
		dest.Port = 22
	}
	dest.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx, `INSERT INTO sftp_destinations(id,organization_id,project_id,name,description,host,port,username,password,private_key,path_prefix,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, dest.ID, dest.OrganizationID, dest.ProjectID, dest.Name, dest.Description, dest.Host, dest.Port, dest.Username, dest.Password, dest.PrivateKey, dest.PathPrefix, dest.CreatedAt)
	return err
}

func (r *sftpRepo) Get(ctx context.Context, id string) (*SFTPDestination, error) {
	var dest SFTPDestination
	if err := r.db.GetContext(ctx, &dest, `SELECT id,COALESCE(organization_id,'') AS organization_id,COALESCE(project_id,'') AS project_id,name,COALESCE(description,'') AS description,host,port,username,COALESCE(password,'') AS password,COALESCE(private_key,'') AS private_key,COALESCE(path_prefix,'') AS path_prefix,last_verified_at,COALESCE(last_verify_error,'') AS last_verify_error,created_at FROM sftp_destinations WHERE id=$1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrSFTPDestinationNotFound
		}
		return nil, err
	}
	return &dest, nil
}

func (r *sftpRepo) ListByProject(ctx context.Context, projectID string) ([]SFTPDestination, error) {
	result := []SFTPDestination{}
	if err := r.db.SelectContext(ctx, &result, `SELECT id,COALESCE(organization_id,'') AS organization_id,COALESCE(project_id,'') AS project_id,name,COALESCE(description,'') AS description,host,port,username,'' AS password,'' AS private_key,COALESCE(path_prefix,'') AS path_prefix,last_verified_at,COALESCE(last_verify_error,'') AS last_verify_error,created_at FROM sftp_destinations WHERE project_id=$1 OR organization_id IN (SELECT organization_id FROM projects WHERE id=$2) ORDER BY name`, projectID, projectID); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *sftpRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM sftp_destinations WHERE id=$1 AND NOT EXISTS(SELECT 1 FROM backup_configs WHERE sftp_destination_id=$2)`, id, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("destination is referenced or missing")
	}
	return nil
}
