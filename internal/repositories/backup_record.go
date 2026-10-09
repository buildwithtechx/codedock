package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"codedock/internal/models"
	"codedock/internal/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func (r *BackupRepo) CreateRecord(ctx context.Context, rec *models.BackupRecord) error {
	if rec.ID == "" {
		rec.ID = uuid.New().String()
	}
	if rec.StartedAt == "" {
		rec.StartedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if rec.Status == "" {
		rec.Status = models.BackupRecordStatusRunning
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO backup_records (id, backup_config_id, database_id, s3_destination_id, sftp_destination_id, parent_record_id, status, file_path, file_size_bytes, s3_url, sftp_url, logs, started_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		rec.ID, rec.BackupConfigID, rec.DatabaseID, rec.S3DestinationID, rec.SFTPDestinationID, rec.ParentRecordID, rec.Status, rec.FilePath, rec.FileSizeBytes, rec.S3URL, rec.SFTPURL, rec.Logs, rec.StartedAt, nullIfEmpty(rec.CompletedAt))
	if err != nil {
		return fmt.Errorf("failed to create backup record: %w", err)
	}
	return nil
}

func (r *BackupRepo) ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*models.BackupRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.BackupRecord
	err := r.db.SelectContext(ctx, &list, `SELECT id, backup_config_id, COALESCE(database_id, '') as database_id, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_record_id, '') as parent_record_id, status, COALESCE(file_path, '') as file_path, file_size_bytes, COALESCE(s3_url, '') as s3_url, COALESCE(sftp_url, '') as sftp_url, COALESCE(logs, '') as logs, started_at, COALESCE(completed_at::text, '') as completed_at, protected_until, sha256, verified_at
		FROM backup_records WHERE backup_config_id = $1 ORDER BY started_at DESC`, backupConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup records: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupRecord, 0)
	}
	return list, nil
}

func (r *BackupRepo) ListRecordsByDatabase(ctx context.Context, databaseID string) ([]*models.BackupRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.BackupRecord
	err := r.db.SelectContext(ctx, &list, `SELECT id, backup_config_id, COALESCE(database_id, '') as database_id, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_record_id, '') as parent_record_id, status, COALESCE(file_path, '') as file_path, file_size_bytes, COALESCE(s3_url, '') as s3_url, COALESCE(sftp_url, '') as sftp_url, COALESCE(logs, '') as logs, started_at, COALESCE(completed_at::text, '') as completed_at, protected_until, sha256, verified_at
		FROM backup_records WHERE database_id = $1 ORDER BY started_at DESC`, databaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup records by database: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupRecord, 0)
	}
	return list, nil
}

func (r *BackupRepo) ListAllRecords(ctx context.Context, limit int) ([]*models.BackupRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var list []*models.BackupRecord
	err := r.db.SelectContext(ctx, &list, `SELECT id, backup_config_id, COALESCE(database_id, '') as database_id, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_record_id, '') as parent_record_id, status, COALESCE(file_path, '') as file_path, file_size_bytes, COALESCE(s3_url, '') as s3_url, COALESCE(sftp_url, '') as sftp_url, COALESCE(logs, '') as logs, started_at, COALESCE(completed_at::text, '') as completed_at, protected_until, sha256, verified_at
		FROM backup_records WHERE EXISTS (SELECT 1 FROM backup_configs WHERE backup_configs.id = backup_records.backup_config_id) ORDER BY started_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list all backup records: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupRecord, 0)
	}
	return list, nil
}

func (r *BackupRepo) GetRecordByID(ctx context.Context, id string) (*models.BackupRecord, error) {
	var rec models.BackupRecord
	err := r.db.GetContext(ctx, &rec, `
		SELECT id, backup_config_id, COALESCE(database_id, '') as database_id, COALESCE(s3_destination_id, '') as s3_destination_id, COALESCE(sftp_destination_id, '') as sftp_destination_id, COALESCE(parent_record_id, '') as parent_record_id, status, COALESCE(file_path, '') as file_path, file_size_bytes, COALESCE(s3_url, '') as s3_url, COALESCE(sftp_url, '') as sftp_url, COALESCE(logs, '') as logs, started_at, COALESCE(completed_at::text, '') as completed_at, protected_until, sha256, verified_at
		FROM backup_records WHERE id = $1`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("Record", id)
		}
		return nil, err
	}
	return &rec, nil
}

func (r *BackupRepo) UpdateRecord(ctx context.Context, rec *models.BackupRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	res, err := r.db.ExecContext(ctx, `
		UPDATE backup_records
		SET sha256 = $1, verified_at = $2, status = $3, file_path = $4, s3_url = $5, sftp_url = $6, s3_destination_id = $7, sftp_destination_id = $8, parent_record_id = $9, logs = $10, file_size_bytes = $11, completed_at = $12
		WHERE id = $13 AND (status!='expiring' OR $14 IN ('expired','failed'))`,
		rec.SHA256, rec.VerifiedAt, rec.Status, rec.FilePath, rec.S3URL, rec.SFTPURL, rec.S3DestinationID, rec.SFTPDestinationID, rec.ParentRecordID, rec.Logs, rec.FileSizeBytes, nullIfEmpty(rec.CompletedAt), rec.ID, rec.Status)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return utils.NewNotFoundError("BackupRecord", rec.ID)
	}
	return nil
}

func nullIfEmpty(value string) any {
	if value == "" {
		return nil
	}
	return value
}

func (r *BackupRepo) DeleteRecord(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.ExecContext(ctx, "DELETE FROM backup_records WHERE id=$1 AND protected_until<=$2 AND status NOT IN ('running','completed')", id, time.Now().Unix())
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("record is protected, active or has not been claimed for deletion")
	}
	return nil
}

func (r *BackupRepo) ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*models.BackupRecord, error) {
	list := make([]*models.BackupRecord, 0)
	if len(configIDs) == 0 {
		return list, nil
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	query, args, err := sqlx.In(`SELECT id, backup_config_id, COALESCE(database_id, '') AS database_id, COALESCE(s3_destination_id, '') AS s3_destination_id, COALESCE(sftp_destination_id, '') AS sftp_destination_id, COALESCE(parent_record_id, '') AS parent_record_id, status, COALESCE(file_path, '') AS file_path, file_size_bytes, COALESCE(s3_url, '') AS s3_url, COALESCE(sftp_url, '') AS sftp_url, COALESCE(logs, '') AS logs, started_at, COALESCE(completed_at::text, '') AS completed_at, protected_until, sha256, verified_at FROM backup_records WHERE backup_config_id IN (?) ORDER BY started_at DESC LIMIT ?`, configIDs, limit)
	if err != nil {
		return nil, fmt.Errorf("prepare scoped backup query: %w", err)
	}
	query = r.db.Rebind(query)
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.db.SelectContext(ctx, &list, query, args...); err != nil {
		return nil, fmt.Errorf("list accessible backup records: %w", err)
	}
	return list, nil
}
