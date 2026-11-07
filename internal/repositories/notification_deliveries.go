package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
)

type NotificationDeliveryRepository interface {
	RecordDelivery(ctx context.Context, delivery *models.NotificationDelivery) error
	ListDeliveries(ctx context.Context, limit int) ([]models.NotificationDelivery, error)
}

type NotificationDeliveryRepo struct {
	db *sqlx.DB
}

func NewNotificationDeliveryRepo(db *sql.DB) *NotificationDeliveryRepo {
	return &NotificationDeliveryRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *NotificationDeliveryRepo) RecordDelivery(ctx context.Context, delivery *models.NotificationDelivery) error {
	if delivery.ID == "" {
		delivery.ID = uuid.NewString()
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO notification_deliveries (id, event_type, channel, title, success, error) VALUES ($1, $2, $3, $4, $5, $6)`,
		delivery.ID, delivery.EventType, delivery.Channel, delivery.Title, delivery.Success, delivery.Error)
	if err != nil {
		return fmt.Errorf("failed to record notification delivery: %w", err)
	}
	return nil
}

func (r *NotificationDeliveryRepo) ListDeliveries(ctx context.Context, limit int) ([]models.NotificationDelivery, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	deliveries := make([]models.NotificationDelivery, 0)
	err := r.db.SelectContext(ctx, &deliveries, `SELECT id, event_type, channel, title, success, COALESCE(error, '') AS error, created_at FROM notification_deliveries ORDER BY created_at DESC, id DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list notification deliveries: %w", err)
	}
	return deliveries, nil
}
