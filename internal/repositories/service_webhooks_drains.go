package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func (r *AppServiceRepo) CreateWebhook(ctx context.Context, w *models.Webhook) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if w.ID == "" {
		w.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	eventTypesStr := strings.Join(w.EventTypes, ",")
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO service_webhooks (id, service_id, url, event_types, include_pr_environments, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		w.ID, w.ServiceID, w.URL, eventTypesStr, w.IncludePREnvironments, w.CreatedAt, w.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create webhook: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) ListWebhooksByService(ctx context.Context, serviceID string) ([]*models.Webhook, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, service_id, url, event_types, include_pr_environments, created_at, updated_at
		 FROM service_webhooks WHERE service_id = $1 ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list webhooks: %w", err)
	}
	defer rows.Close()
	var out []*models.Webhook
	for rows.Next() {
		var w models.Webhook
		var eventsStr string
		if err := rows.Scan(&w.ID, &w.ServiceID, &w.URL, &eventsStr, &w.IncludePREnvironments, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan webhook: %w", err)
		}
		if eventsStr != "" {
			w.EventTypes = strings.Split(eventsStr, ",")
		} else {
			w.EventTypes = []string{}
		}
		out = append(out, &w)
	}
	return out, rows.Err()
}

func (r *AppServiceRepo) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := "DELETE FROM service_webhooks WHERE id = $1 AND service_id = $2"
	res, err := r.db.ExecContext(ctx, query, id, serviceID)
	if err != nil {
		return fmt.Errorf("delete webhook: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return utils.NewNotFoundError("Webhook", id)
	}
	return nil
}

func (r *AppServiceRepo) CreateLogDrain(ctx context.Context, d *models.LogDrain) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO log_drains (id, service_id, project_id, drain_type, endpoint_url, auth_token, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		d.ID, d.ServiceID, d.ProjectID, d.DrainType, d.EndpointURL, d.AuthToken, d.CreatedAt, d.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create log drain: %w", err)
	}
	return nil
}

func (r *AppServiceRepo) ListLogDrainsByService(ctx context.Context, serviceID string) ([]*models.LogDrain, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, service_id, project_id, drain_type, endpoint_url, COALESCE(auth_token, '') AS auth_token, created_at, updated_at
		 FROM log_drains WHERE service_id = $1 ORDER BY created_at DESC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list log drains: %w", err)
	}
	defer rows.Close()
	var out []*models.LogDrain
	for rows.Next() {
		var d models.LogDrain
		if err := rows.Scan(&d.ID, &d.ServiceID, &d.ProjectID, &d.DrainType, &d.EndpointURL, &d.AuthToken, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan log drain: %w", err)
		}
		out = append(out, &d)
	}
	return out, rows.Err()
}

func (r *AppServiceRepo) DeleteLogDrain(ctx context.Context, id, serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := "DELETE FROM log_drains WHERE id = $1 AND service_id = $2"
	res, err := r.db.ExecContext(ctx, query, id, serviceID)
	if err != nil {
		return fmt.Errorf("delete log drain: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return utils.NewNotFoundError("LogDrain", id)
	}
	return nil
}
