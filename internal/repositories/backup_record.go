package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"github.com/google/uuid"
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
	_, err := r.db.ExecContext(ctx, `INSERT INTO backup_records (id, backup_config_id, database_id, status, file_path, file_size_bytes, s3_url, logs, started_at, completed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.ID, rec.BackupConfigID, rec.DatabaseID, rec.Status, rec.FilePath, rec.FileSizeBytes, rec.S3URL, rec.Logs, rec.StartedAt, rec.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to create backup record: %w", err)
	}
	return nil
}

func (r *BackupRepo) ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*models.BackupRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var list []*models.BackupRecord
	err := r.db.SelectContext(ctx, &list, `SELECT id, backup_config_id, COALESCE(database_id, '') as database_id, status, COALESCE(file_path, '') as file_path, file_size_bytes, COALESCE(s3_url, '') as s3_url, COALESCE(logs, '') as logs, started_at, COALESCE(completed_at, '') as completed_at
		FROM backup_records WHERE backup_config_id = ? ORDER BY started_at DESC`, backupConfigID)
	if err != nil {
		return nil, fmt.Errorf("failed to list backup records: %w", err)
	}
	if list == nil {
		list = make([]*models.BackupRecord, 0)
	}
	return list, nil
}

func (r *BackupRepo) GetRecordByID(ctx context.Context, id string) (*models.BackupRecord, error) {
	var rec models.BackupRecord
	err := r.db.GetContext(ctx, &rec, `
		SELECT id, backup_config_id, COALESCE(database_id, '') as database_id, status, COALESCE(file_path, '') as file_path, file_size_bytes, COALESCE(s3_url, '') as s3_url, COALESCE(logs, '') as logs, started_at, COALESCE(completed_at, '') as completed_at
		FROM backup_records WHERE id = ?`, id)
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
		SET status = ?, file_path = ?, s3_url = ?, logs = ?, file_size_bytes = ?, completed_at = ?
		WHERE id = ?`,
		rec.Status, rec.FilePath, rec.S3URL, rec.Logs, rec.FileSizeBytes, rec.CompletedAt, rec.ID)
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

func (r *BackupRepo) DeleteRecord(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, "DELETE FROM backup_records WHERE id=?", id)
	return err
}
