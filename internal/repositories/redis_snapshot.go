package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RedisSnapshotRepo struct {
	db *sqlx.DB
}

func NewRedisSnapshotRepo(db *sql.DB) *RedisSnapshotRepo {
	return &RedisSnapshotRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *RedisSnapshotRepo) Create(ctx context.Context, snapshot *models.RedisSnapshot) error {
	if snapshot.ID == "" {
		snapshot.ID = uuid.NewString()
	}
	if snapshot.CreatedAt == "" {
		snapshot.CreatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO redis_snapshots(id,database_id,project_id,cluster_id,s3_destination_id,s3_key,size_bytes,status,error,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`, snapshot.ID, snapshot.DatabaseID, snapshot.ProjectID, snapshot.ClusterID, nullableID(snapshot.S3DestinationID), snapshot.S3Key, snapshot.SizeBytes, snapshot.Status, snapshot.Error, snapshot.CreatedAt)
	return err
}

func (r *RedisSnapshotRepo) List(ctx context.Context, databaseID string) ([]models.RedisSnapshot, error) {
	result := []models.RedisSnapshot{}
	if err := r.db.SelectContext(ctx, &result, `SELECT id,database_id,project_id,cluster_id,COALESCE(s3_destination_id,'') AS s3_destination_id,s3_key,size_bytes,status,error,created_at FROM redis_snapshots WHERE database_id=$1 ORDER BY created_at DESC`, databaseID); err != nil {
		return nil, fmt.Errorf("list redis snapshots: %w", err)
	}
	return result, nil
}

func (r *RedisSnapshotRepo) Get(ctx context.Context, id string) (*models.RedisSnapshot, error) {
	var snapshot models.RedisSnapshot
	if err := r.db.GetContext(ctx, &snapshot, `SELECT id,database_id,project_id,cluster_id,COALESCE(s3_destination_id,'') AS s3_destination_id,s3_key,size_bytes,status,error,created_at FROM redis_snapshots WHERE id=$1`, id); err != nil {
		return nil, err
	}
	return &snapshot, nil
}
