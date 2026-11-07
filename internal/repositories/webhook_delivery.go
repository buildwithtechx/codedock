package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock/internal/models"
)

type WebhookDeliveryRepository interface {
	ClaimGithub(ctx context.Context, input models.WebhookDeliveryInput) (models.ClaimResult, error)
	Record(ctx context.Context, input models.WebhookDeliveryInput) (string, error)
	MarkProcessed(ctx context.Context, id string, res models.WebhookDeliveryResult) error
	ListByProject(ctx context.Context, projectID string, limit int, cursor string) ([]models.WebhookDelivery, string, error)
	Prune(ctx context.Context, olderThan time.Duration) (int64, error)
}

type WebhookDeliveryRepo struct {
	db *sql.DB
}

func NewWebhookDeliveryRepo(db *sql.DB) *WebhookDeliveryRepo {
	return &WebhookDeliveryRepo{db: db}
}

func marshalSummary(summary any) []byte {
	if summary == nil {
		return nil
	}
	bytes, err := json.Marshal(summary)
	if err != nil {
		return nil
	}
	return bytes
}

func parseHandledProjectIDs(rawSummary []byte) []string {
	if len(rawSummary) == 0 {
		return nil
	}
	var parsed struct {
		HandledProjectIDs []string `json:"handledProjectIds"`
	}
	if err := json.Unmarshal(rawSummary, &parsed); err == nil {
		return parsed.HandledProjectIDs
	}
	return nil
}

func (r *WebhookDeliveryRepo) ClaimGithub(ctx context.Context, input models.WebhookDeliveryInput) (models.ClaimResult, error) {
	if input.DeliveryID == nil || *input.DeliveryID == "" {
		return models.ClaimResult{Claimed: true, ID: ""}, nil
	}

	id := "wdl_" + uuid.NewString()
	summaryBytes := marshalSummary(input.Summary)

	insertQuery := `
		INSERT INTO webhook_deliveries (
			id, organization_id, project_id, source, hook_id, delivery_id, event,
			auth_result, outcome, action_ref, status_code, error, summary, client_ip, user_agent, duration_ms, received_at
		) VALUES (
			$1, $2, $3, 'github', $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW()
		)
		ON CONFLICT (source, delivery_id) WHERE source = 'github' AND delivery_id IS NOT NULL
		DO NOTHING
		RETURNING id;
	`

	var claimedID string
	err := r.db.QueryRowContext(ctx, insertQuery,
		id, input.OrganizationID, input.ProjectID, input.HookID, input.DeliveryID, input.Event,
		input.AuthResult, "received", input.ActionRef, input.StatusCode, input.Error, summaryBytes,
		input.ClientIP, input.UserAgent, input.DurationMs,
	).Scan(&claimedID)

	if err == nil && claimedID != "" {
		return models.ClaimResult{Claimed: true, ID: claimedID}, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return models.ClaimResult{Claimed: false}, fmt.Errorf("failed to claim github delivery: %w", err)
	}

	// Conflict occurred. Check if this is an operator/webhook retry of a finished failure.
	retryQuery := `
		UPDATE webhook_deliveries
		SET outcome = 'received', processed_at = NULL, error = NULL, status_code = NULL, received_at = NOW()
		WHERE source = 'github'
		  AND delivery_id = $1
		  AND event = $2
		  AND outcome = 'failed'
		  AND processed_at IS NOT NULL
		RETURNING id, summary;
	`

	var retryID string
	var rawSummary []byte
	retryErr := r.db.QueryRowContext(ctx, retryQuery, *input.DeliveryID, input.Event).Scan(&retryID, &rawSummary)
	if retryErr == nil && retryID != "" {
		return models.ClaimResult{
			Claimed:           true,
			ID:                retryID,
			HandledProjectIDs: parseHandledProjectIDs(rawSummary),
		}, nil
	}
	if retryErr != nil && retryErr != sql.ErrNoRows {
		return models.ClaimResult{Claimed: false}, fmt.Errorf("failed to check failed delivery retry: %w", retryErr)
	}

	// Delivery is in-flight or already succeeded.
	return models.ClaimResult{Claimed: false, ID: ""}, nil
}

func (r *WebhookDeliveryRepo) Record(ctx context.Context, input models.WebhookDeliveryInput) (string, error) {
	id := "wdl_" + uuid.NewString()
	summaryBytes := marshalSummary(input.Summary)

	query := `
		INSERT INTO webhook_deliveries (
			id, organization_id, project_id, source, hook_id, delivery_id, event,
			auth_result, outcome, action_ref, status_code, error, summary, client_ip, user_agent, duration_ms, received_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW()
		)
		RETURNING id;
	`
	var returnedID string
	err := r.db.QueryRowContext(ctx, query,
		id, input.OrganizationID, input.ProjectID, input.Source, input.HookID, input.DeliveryID, input.Event,
		input.AuthResult, input.Outcome, input.ActionRef, input.StatusCode, input.Error, summaryBytes,
		input.ClientIP, input.UserAgent, input.DurationMs,
	).Scan(&returnedID)
	if err != nil {
		return "", fmt.Errorf("failed to record webhook delivery: %w", err)
	}
	return returnedID, nil
}

func (r *WebhookDeliveryRepo) MarkProcessed(ctx context.Context, id string, res models.WebhookDeliveryResult) error {
	summaryBytes := marshalSummary(res.Summary)
	query := `
		UPDATE webhook_deliveries
		SET outcome = $1, action_ref = $2, status_code = $3, error = $4, summary = $5, duration_ms = $6, processed_at = NOW()
		WHERE id = $7;
	`
	_, err := r.db.ExecContext(ctx, query,
		res.Outcome, res.ActionRef, res.StatusCode, res.Error, summaryBytes, res.DurationMs, id,
	)
	if err != nil {
		return fmt.Errorf("failed to mark webhook delivery processed: %w", err)
	}
	return nil
}

func (r *WebhookDeliveryRepo) ListByProject(ctx context.Context, projectID string, limit int, cursor string) ([]models.WebhookDelivery, string, error) {
	if limit <= 0 || limit > 100 {
		limit = 25
	}

	query := `
		SELECT id, organization_id, project_id, source, hook_id, delivery_id, event,
		       auth_result, outcome, action_ref, status_code, error, summary, client_ip,
		       user_agent, duration_ms, received_at, processed_at
		FROM webhook_deliveries
		WHERE project_id = $1
		ORDER BY received_at DESC, id DESC
		LIMIT $2;
	`

	rows, err := r.db.QueryContext(ctx, query, projectID, limit+1)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list webhook deliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []models.WebhookDelivery
	for rows.Next() {
		var d models.WebhookDelivery
		var summaryBytes []byte
		if err := rows.Scan(
			&d.ID, &d.OrganizationID, &d.ProjectID, &d.Source, &d.HookID, &d.DeliveryID, &d.Event,
			&d.AuthResult, &d.Outcome, &d.ActionRef, &d.StatusCode, &d.Error, &summaryBytes,
			&d.ClientIP, &d.UserAgent, &d.DurationMs, &d.ReceivedAt, &d.ProcessedAt,
		); err != nil {
			return nil, "", fmt.Errorf("failed to scan webhook delivery: %w", err)
		}
		if len(summaryBytes) > 0 {
			d.Summary = json.RawMessage(summaryBytes)
		}
		deliveries = append(deliveries, d)
	}

	var nextCursor string
	if len(deliveries) > limit {
		last := deliveries[limit-1]
		nextCursor = last.ID
		deliveries = deliveries[:limit]
	}

	return deliveries, nextCursor, nil
}

func (r *WebhookDeliveryRepo) Prune(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-olderThan)
	query := `DELETE FROM webhook_deliveries WHERE received_at < $1;`
	res, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to prune webhook deliveries: %w", err)
	}
	return res.RowsAffected()
}
