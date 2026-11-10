package git

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type webhookDeliveryRepo struct {
	db *sql.DB
}

func NewWebhookDeliveryRepository(db *sql.DB) WebhookDeliveryRepository {
	return &webhookDeliveryRepo{db: db}
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

func (r *webhookDeliveryRepo) ClaimGithub(ctx context.Context, input WebhookDeliveryInput) (ClaimResult, error) {
	if input.DeliveryID == nil || *input.DeliveryID == "" {
		return ClaimResult{Claimed: true, ID: ""}, nil
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
		return ClaimResult{Claimed: true, ID: claimedID}, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return ClaimResult{Claimed: false}, fmt.Errorf("failed to claim github delivery: %w", err)
	}

	// Conflict occurred. Check if this is an operator retry of a finished failure.
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
	var retriedID string
	var rawSummary []byte
	err = r.db.QueryRowContext(ctx, retryQuery, input.DeliveryID, input.Event).Scan(&retriedID, &rawSummary)
	if err == nil && retriedID != "" {
		handled := parseHandledProjectIDs(rawSummary)
		return ClaimResult{Claimed: true, ID: retriedID, HandledProjectIDs: handled}, nil
	}

	return ClaimResult{Claimed: false}, nil
}

func (r *webhookDeliveryRepo) Record(ctx context.Context, input WebhookDeliveryInput) (string, error) {
	id := "wdl_" + uuid.NewString()
	summaryBytes := marshalSummary(input.Summary)

	query := `
		INSERT INTO webhook_deliveries (
			id, organization_id, project_id, source, hook_id, delivery_id, event,
			auth_result, outcome, action_ref, status_code, error, summary, client_ip, user_agent, duration_ms, received_at, processed_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, NOW(), NOW()
		)
		RETURNING id;
	`
	var recordedID string
	err := r.db.QueryRowContext(ctx, query,
		id, input.OrganizationID, input.ProjectID, input.Source, input.HookID, input.DeliveryID, input.Event,
		input.AuthResult, input.Outcome, input.ActionRef, input.StatusCode, input.Error, summaryBytes,
		input.ClientIP, input.UserAgent, input.DurationMs,
	).Scan(&recordedID)

	return recordedID, err
}

func (r *webhookDeliveryRepo) MarkProcessed(ctx context.Context, id string, res WebhookDeliveryResult) error {
	if id == "" {
		return nil
	}
	summaryBytes := marshalSummary(res.Summary)

	query := `
		UPDATE webhook_deliveries
		SET outcome = $1, action_ref = $2, status_code = $3, error = $4, summary = $5, duration_ms = $6, processed_at = NOW()
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		res.Outcome, res.ActionRef, res.StatusCode, res.Error, summaryBytes, res.DurationMs, id,
	)
	return err
}

func (r *webhookDeliveryRepo) ListByProject(ctx context.Context, projectID string, limit int, cursor string) ([]WebhookDelivery, string, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	query := `
		SELECT id, organization_id, project_id, source, hook_id, delivery_id, event,
		       auth_result, outcome, action_ref, status_code, error, summary, client_ip, user_agent, duration_ms, received_at, processed_at
		FROM webhook_deliveries
		WHERE project_id = $1
	`
	args := []any{projectID}

	if cursor != "" {
		query += " AND received_at < $2 ORDER BY received_at DESC LIMIT $3"
		args = append(args, cursor, limit+1)
	} else {
		query += " ORDER BY received_at DESC LIMIT $2"
		args = append(args, limit+1)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to query webhook deliveries: %w", err)
	}
	defer rows.Close()

	var deliveries []WebhookDelivery
	for rows.Next() {
		var d WebhookDelivery
		err := rows.Scan(
			&d.ID, &d.OrganizationID, &d.ProjectID, &d.Source, &d.HookID, &d.DeliveryID, &d.Event,
			&d.AuthResult, &d.Outcome, &d.ActionRef, &d.StatusCode, &d.Error, &d.Summary,
			&d.ClientIP, &d.UserAgent, &d.DurationMs, &d.ReceivedAt, &d.ProcessedAt,
		)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan webhook delivery: %w", err)
		}
		deliveries = append(deliveries, d)
	}

	nextCursor := ""
	if len(deliveries) > limit {
		nextCursor = deliveries[limit-1].ReceivedAt.Format(time.RFC3339Nano)
		deliveries = deliveries[:limit]
	}

	return deliveries, nextCursor, nil
}

func (r *webhookDeliveryRepo) Prune(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoff := time.Now().Add(-olderThan)
	res, err := r.db.ExecContext(ctx, "DELETE FROM webhook_deliveries WHERE received_at < $1", cutoff)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
