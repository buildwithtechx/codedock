package chatapps

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *repo) CreateChallenge(ctx context.Context, ch *ChatAppChallenge) error {
	now := time.Now().UTC()
	if ch.CreatedAt.IsZero() {
		ch.CreatedAt = now
	}
	query := `
		INSERT INTO chat_app_challenges (
			id, connection_id, provider, external_user_id, external_channel_id,
			token, code, expires_at, verified_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		ch.ID, ch.ConnectionID, string(ch.Provider), ch.ExternalUserID, ch.ExternalChannelID,
		ch.Token, ch.Code, ch.ExpiresAt, ch.VerifiedAt, ch.CreatedAt,
	)
	return err
}

func (r *repo) GetChallengeByToken(ctx context.Context, token string) (*ChatAppChallenge, error) {
	query := `
		SELECT id, connection_id, provider, external_user_id, external_channel_id,
		       token, code, expires_at, verified_at, created_at
		FROM chat_app_challenges
		WHERE token = $1
	`
	var ch ChatAppChallenge
	var provider string
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&ch.ID, &ch.ConnectionID, &provider, &ch.ExternalUserID, &ch.ExternalChannelID,
		&ch.Token, &ch.Code, &ch.ExpiresAt, &ch.VerifiedAt, &ch.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get challenge: %w", err)
	}
	ch.Provider = ChatAppProvider(provider)
	return &ch, nil
}

func (r *repo) VerifyChallenge(ctx context.Context, token string, code string) (*ChatAppChallenge, error) {
	query := `
		SELECT id, connection_id, provider, external_user_id, external_channel_id,
		       token, code, expires_at, verified_at, created_at
		FROM chat_app_challenges
		WHERE (token = $1 OR code = $2) AND expires_at > CURRENT_TIMESTAMP
	`
	var ch ChatAppChallenge
	var provider string
	err := r.db.QueryRowContext(ctx, query, token, code).Scan(
		&ch.ID, &ch.ConnectionID, &provider, &ch.ExternalUserID, &ch.ExternalChannelID,
		&ch.Token, &ch.Code, &ch.ExpiresAt, &ch.VerifiedAt, &ch.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrChallengeExpired
		}
		return nil, err
	}
	ch.Provider = ChatAppProvider(provider)

	now := time.Now().UTC()
	_, _ = r.db.ExecContext(ctx, `UPDATE chat_app_challenges SET verified_at = $1 WHERE id = $2`, now, ch.ID)
	ch.VerifiedAt = &now
	return &ch, nil
}

func (r *repo) CreateCommand(ctx context.Context, cmd *ChatAppCommand) error {
	now := time.Now().UTC()
	if cmd.CreatedAt.IsZero() {
		cmd.CreatedAt = now
	}
	query := `
		INSERT INTO chat_app_commands (
			id, connection_id, identity_id, provider, command, args, status,
			response, confirmation_token, confirmed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.ExecContext(ctx, query,
		cmd.ID, cmd.ConnectionID, cmd.IdentityID, string(cmd.Provider), cmd.Command,
		cmd.Args, cmd.Status, cmd.Response, cmd.ConfirmationToken, cmd.ConfirmedAt, cmd.CreatedAt,
	)
	return err
}

func (r *repo) GetCommandByConfirmationToken(ctx context.Context, token string) (*ChatAppCommand, error) {
	query := `
		SELECT id, connection_id, identity_id, provider, command, args, status,
		       response, confirmation_token, confirmed_at, created_at
		FROM chat_app_commands
		WHERE confirmation_token = $1 AND status = 'awaiting_confirmation'
	`
	var cmd ChatAppCommand
	var provider string
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&cmd.ID, &cmd.ConnectionID, &cmd.IdentityID, &provider, &cmd.Command, &cmd.Args,
		&cmd.Status, &cmd.Response, &cmd.ConfirmationToken, &cmd.ConfirmedAt, &cmd.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get command by token: %w", err)
	}
	cmd.Provider = ChatAppProvider(provider)
	return &cmd, nil
}

func (r *repo) MarkCommandConfirmed(ctx context.Context, id string, response string) error {
	now := time.Now().UTC()
	query := `
		UPDATE chat_app_commands
		SET status = 'executed', confirmed_at = $2, response = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, now, response)
	return err
}

func (r *repo) ListProjects(ctx context.Context, limit, offset int) ([]ProjectSummary, error) {
	if r.db == nil {
		return nil, nil
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, COALESCE(description, '') FROM projects ORDER BY created_at DESC LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []ProjectSummary
	for rows.Next() {
		var p ProjectSummary
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *repo) GetProject(ctx context.Context, id string) (*ProjectSummary, error) {
	if r.db == nil {
		return nil, nil
	}
	var p ProjectSummary
	err := r.db.QueryRowContext(ctx, `SELECT id, name, COALESCE(description, '') FROM projects WHERE id = $1`, id).Scan(&p.ID, &p.Name, &p.Description)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repo) ListAppServices(ctx context.Context, projectID string) ([]AppServiceSummary, error) {
	if r.db == nil {
		return nil, nil
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, COALESCE(status, ''), COALESCE(domain, ''), COALESCE(image_ref, ''), COALESCE(branch, ''), COALESCE(updated_at::text, '') FROM app_services WHERE project_id = $1`, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []AppServiceSummary
	for rows.Next() {
		var a AppServiceSummary
		if err := rows.Scan(&a.ID, &a.Name, &a.Status, &a.Domain, &a.ImageRef, &a.Branch, &a.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *repo) GetAppService(ctx context.Context, id string) (*AppServiceSummary, error) {
	if r.db == nil {
		return nil, nil
	}
	var a AppServiceSummary
	err := r.db.QueryRowContext(ctx, `SELECT id, name, COALESCE(status, ''), COALESCE(domain, ''), COALESCE(image_ref, ''), COALESCE(branch, ''), COALESCE(updated_at::text, '') FROM app_services WHERE id = $1`, id).Scan(&a.ID, &a.Name, &a.Status, &a.Domain, &a.ImageRef, &a.Branch, &a.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repo) RedeployApp(ctx context.Context, appID, username, provider string) (string, error) {
	if r.db == nil {
		return "", fmt.Errorf("database unavailable")
	}
	depID := "dep_" + uuid.NewString()
	_, err := r.db.ExecContext(ctx, `INSERT INTO deployments (id, service_id, status, commit_message, trigger) VALUES ($1, $2, 'BUILDING', $3, $4)`,
		depID, appID, "Chat Command Redeploy by "+username, "Chat: "+provider)
	return depID, err
}

func (r *repo) RestartApp(ctx context.Context, appID string) error {
	if r.db == nil {
		return fmt.Errorf("database unavailable")
	}
	_, err := r.db.ExecContext(ctx, `UPDATE app_services SET status = 'running' WHERE id = $1`, appID)
	return err
}
