package webhooks

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateWebhook(ctx context.Context, w *Webhook) error
	ListWebhooksByService(ctx context.Context, serviceID string) ([]*Webhook, error)
	DeleteWebhook(ctx context.Context, id, serviceID string) error
	CreateLogDrain(ctx context.Context, d *LogDrain) error
	ListLogDrainsByService(ctx context.Context, serviceID string) ([]*LogDrain, error)
	DeleteLogDrain(ctx context.Context, id, serviceID string) error
}

type Repo struct {
	mu sync.RWMutex
	db *sqlx.DB
}

func NewRepository(db *sql.DB) *Repo {
	return &Repo{db: sqlx.NewDb(db, "pgx")}
}

func (r *Repo) CreateWebhook(_ context.Context, w *Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	_, err := r.db.Exec(`INSERT INTO service_webhooks (id, service_id, url, secret, events, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		w.ID, w.ServiceID, w.URL, w.Secret, w.Events, w.CreatedAt, w.UpdatedAt)
	return err
}

func (r *Repo) ListWebhooksByService(_ context.Context, serviceID string) ([]*Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*Webhook
	err := r.db.Select(&list, `SELECT id, service_id, url, secret, events, created_at, updated_at FROM service_webhooks WHERE service_id = $1 ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repo) DeleteWebhook(_ context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM service_webhooks WHERE id = $1 AND service_id = $2`, id, serviceID)
	return err
}

func (r *Repo) CreateLogDrain(_ context.Context, d *LogDrain) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.Exec(`INSERT INTO service_log_drains (id, project_id, service_id, drain_type, endpoint_url, api_key, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		d.ID, d.ProjectID, d.ServiceID, d.DrainType, d.EndpointURL, d.APIKey, d.CreatedAt, d.UpdatedAt)
	return err
}

func (r *Repo) ListLogDrainsByService(_ context.Context, serviceID string) ([]*LogDrain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*LogDrain
	err := r.db.Select(&list, `SELECT id, project_id, service_id, drain_type, endpoint_url, api_key, created_at, updated_at FROM service_log_drains WHERE service_id = $1 ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (r *Repo) DeleteLogDrain(_ context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.Exec(`DELETE FROM service_log_drains WHERE id = $1 AND service_id = $2`, id, serviceID)
	return err
}
