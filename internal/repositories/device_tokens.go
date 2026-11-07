package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
)

type DeviceTokenRepository interface {
	RegisterToken(ctx context.Context, userID, token, platform string) (*models.DeviceToken, error)
	DeleteToken(ctx context.Context, userID, token string) error
	ListTokensByUser(ctx context.Context, userID string) ([]*models.DeviceToken, error)
	ListAllTokens(ctx context.Context) ([]*models.DeviceToken, error)
}

type DeviceTokenRepo struct {
	db *sqlx.DB
}

func NewDeviceTokenRepo(db *sql.DB) *DeviceTokenRepo {
	return &DeviceTokenRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *DeviceTokenRepo) RegisterToken(ctx context.Context, userID, token, platform string) (*models.DeviceToken, error) {
	if token == "" {
		return nil, fmt.Errorf("device token is required")
	}
	if platform == "" {
		platform = "expo"
	}
	id := uuid.NewString()
	now := time.Now().UTC()

	query := `INSERT INTO notification_device_tokens (id, user_id, token, platform, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $5)
	          ON CONFLICT(token) DO UPDATE SET
	              user_id = EXCLUDED.user_id,
	              platform = EXCLUDED.platform,
	              updated_at = EXCLUDED.updated_at
	          RETURNING id, user_id, token, platform, created_at, updated_at`

	var dt models.DeviceToken
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, id, userID, token, platform, now).
		Scan(&dt.ID, &dt.UserID, &dt.Token, &dt.Platform, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("register device token: %w", err)
	}
	dt.CreatedAt = createdAt.Format(time.RFC3339)
	dt.UpdatedAt = updatedAt.Format(time.RFC3339)
	return &dt, nil
}

func (r *DeviceTokenRepo) DeleteToken(ctx context.Context, userID, token string) error {
	query := `DELETE FROM notification_device_tokens WHERE token = $1`
	args := []any{token}
	if userID != "" {
		query += ` AND user_id = $2`
		args = append(args, userID)
	}
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete device token: %w", err)
	}
	return nil
}

func (r *DeviceTokenRepo) ListTokensByUser(ctx context.Context, userID string) ([]*models.DeviceToken, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, token, platform, created_at, updated_at FROM notification_device_tokens WHERE user_id = $1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, fmt.Errorf("list device tokens by user: %w", err)
	}
	defer rows.Close()

	tokens := make([]*models.DeviceToken, 0)
	for rows.Next() {
		var dt models.DeviceToken
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&dt.ID, &dt.UserID, &dt.Token, &dt.Platform, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		dt.CreatedAt = createdAt.Format(time.RFC3339)
		dt.UpdatedAt = updatedAt.Format(time.RFC3339)
		tokens = append(tokens, &dt)
	}
	return tokens, rows.Err()
}

func (r *DeviceTokenRepo) ListAllTokens(ctx context.Context) ([]*models.DeviceToken, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, token, platform, created_at, updated_at FROM notification_device_tokens ORDER BY updated_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("list all device tokens: %w", err)
	}
	defer rows.Close()

	tokens := make([]*models.DeviceToken, 0)
	for rows.Next() {
		var dt models.DeviceToken
		var createdAt, updatedAt time.Time
		if err := rows.Scan(&dt.ID, &dt.UserID, &dt.Token, &dt.Platform, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		dt.CreatedAt = createdAt.Format(time.RFC3339)
		dt.UpdatedAt = updatedAt.Format(time.RFC3339)
		tokens = append(tokens, &dt)
	}
	return tokens, rows.Err()
}
