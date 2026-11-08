package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"codedock/internal/models"
)

type ChatAppRepository interface {
	CreateConnection(ctx context.Context, conn *models.ChatAppConnection) error
	GetConnection(ctx context.Context, id string) (*models.ChatAppConnection, error)
	ListConnectionsByUser(ctx context.Context, userID string) ([]*models.ChatAppConnection, error)
	UpdateConnection(ctx context.Context, conn *models.ChatAppConnection) error
	DeleteConnection(ctx context.Context, id string) error

	CreateIdentity(ctx context.Context, identity *models.ChatAppIdentity) error
	GetIdentityByExternal(ctx context.Context, connectionID string, extUserID string) (*models.ChatAppIdentity, error)
	ListIdentitiesByConnection(ctx context.Context, connectionID string) ([]*models.ChatAppIdentity, error)
	DeleteIdentity(ctx context.Context, id string) error

	CreateChallenge(ctx context.Context, ch *models.ChatAppChallenge) error
	GetChallengeByToken(ctx context.Context, token string) (*models.ChatAppChallenge, error)
	VerifyChallenge(ctx context.Context, token string, code string) (*models.ChatAppChallenge, error)

	CreateCommand(ctx context.Context, cmd *models.ChatAppCommand) error
	GetCommandByConfirmationToken(ctx context.Context, token string) (*models.ChatAppCommand, error)
	MarkCommandConfirmed(ctx context.Context, id string, response string) error
}

type chatAppRepo struct {
	db *sql.DB
}

func NewChatAppRepo(db *sql.DB) ChatAppRepository {
	return &chatAppRepo{db: db}
}

func (r *chatAppRepo) CreateConnection(ctx context.Context, conn *models.ChatAppConnection) error {
	configJSON, err := json.Marshal(conn.Config)
	if err != nil {
		configJSON = []byte("{}")
	}
	query := `
		INSERT INTO chat_app_connections (
			id, user_id, provider, transport, purpose, name, bot_username,
			credentials_encrypted, config, enabled, status, last_error, last_event_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	now := time.Now().UTC()
	if conn.CreatedAt.IsZero() {
		conn.CreatedAt = now
	}
	conn.UpdatedAt = now

	_, err = r.db.ExecContext(ctx, query,
		conn.ID, conn.UserID, string(conn.Provider), string(conn.Transport), string(conn.Purpose),
		conn.Name, conn.BotUsername, conn.CredentialsEncrypted, configJSON, conn.Enabled,
		string(conn.Status), conn.LastError, conn.LastEventAt, conn.CreatedAt, conn.UpdatedAt,
	)
	return err
}

func (r *chatAppRepo) GetConnection(ctx context.Context, id string) (*models.ChatAppConnection, error) {
	query := `
		SELECT id, user_id, provider, transport, purpose, name, bot_username,
		       credentials_encrypted, config, enabled, status, last_error, last_event_at, created_at, updated_at
		FROM chat_app_connections WHERE id = $1
	`
	var conn models.ChatAppConnection
	var provider, transport, purpose, status string
	var configBytes []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&conn.ID, &conn.UserID, &provider, &transport, &purpose, &conn.Name, &conn.BotUsername,
		&conn.CredentialsEncrypted, &configBytes, &conn.Enabled, &status, &conn.LastError,
		&conn.LastEventAt, &conn.CreatedAt, &conn.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get chat connection: %w", err)
	}

	conn.Provider = models.ChatAppProvider(provider)
	conn.Transport = models.ChatAppTransport(transport)
	conn.Purpose = models.ChatAppPurpose(purpose)
	conn.Status = models.ChatAppStatus(status)
	if len(configBytes) > 0 {
		_ = json.Unmarshal(configBytes, &conn.Config)
	}
	if conn.Config == nil {
		conn.Config = make(map[string]any)
	}
	return &conn, nil
}

func (r *chatAppRepo) ListConnectionsByUser(ctx context.Context, userID string) ([]*models.ChatAppConnection, error) {
	query := `
		SELECT id, user_id, provider, transport, purpose, name, bot_username,
		       credentials_encrypted, config, enabled, status, last_error, last_event_at, created_at, updated_at
		FROM chat_app_connections WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list chat connections: %w", err)
	}
	defer rows.Close()

	var list []*models.ChatAppConnection
	for rows.Next() {
		var conn models.ChatAppConnection
		var provider, transport, purpose, status string
		var configBytes []byte

		if err := rows.Scan(
			&conn.ID, &conn.UserID, &provider, &transport, &purpose, &conn.Name, &conn.BotUsername,
			&conn.CredentialsEncrypted, &configBytes, &conn.Enabled, &status, &conn.LastError,
			&conn.LastEventAt, &conn.CreatedAt, &conn.UpdatedAt,
		); err != nil {
			return nil, err
		}
		conn.Provider = models.ChatAppProvider(provider)
		conn.Transport = models.ChatAppTransport(transport)
		conn.Purpose = models.ChatAppPurpose(purpose)
		conn.Status = models.ChatAppStatus(status)
		if len(configBytes) > 0 {
			_ = json.Unmarshal(configBytes, &conn.Config)
		}
		if conn.Config == nil {
			conn.Config = make(map[string]any)
		}
		list = append(list, &conn)
	}
	return list, nil
}

func (r *chatAppRepo) UpdateConnection(ctx context.Context, conn *models.ChatAppConnection) error {
	configJSON, err := json.Marshal(conn.Config)
	if err != nil {
		configJSON = []byte("{}")
	}
	query := `
		UPDATE chat_app_connections
		SET transport = $2, purpose = $3, name = $4, bot_username = $5,
		    credentials_encrypted = $6, config = $7, enabled = $8, status = $9,
		    last_error = $10, last_event_at = $11, updated_at = $12
		WHERE id = $1
	`
	conn.UpdatedAt = time.Now().UTC()
	_, err = r.db.ExecContext(ctx, query,
		conn.ID, string(conn.Transport), string(conn.Purpose), conn.Name, conn.BotUsername,
		conn.CredentialsEncrypted, configJSON, conn.Enabled, string(conn.Status),
		conn.LastError, conn.LastEventAt, conn.UpdatedAt,
	)
	return err
}

func (r *chatAppRepo) DeleteConnection(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chat_app_connections WHERE id = $1`, id)
	return err
}

func (r *chatAppRepo) CreateIdentity(ctx context.Context, identity *models.ChatAppIdentity) error {
	permsJSON, err := json.Marshal(identity.Permissions)
	if err != nil {
		permsJSON = []byte("[]")
	}
	now := time.Now().UTC()
	if identity.CreatedAt.IsZero() {
		identity.CreatedAt = now
	}
	if identity.PairedAt.IsZero() {
		identity.PairedAt = now
	}
	identity.UpdatedAt = now

	query := `
		INSERT INTO chat_app_identities (
			id, connection_id, user_id, provider, external_user_id, external_username,
			role, permissions, active, paired_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (connection_id, external_user_id) DO UPDATE
		SET user_id = EXCLUDED.user_id, external_username = EXCLUDED.external_username,
		    role = EXCLUDED.role, permissions = EXCLUDED.permissions, active = EXCLUDED.active,
		    paired_at = EXCLUDED.paired_at, updated_at = EXCLUDED.updated_at
	`
	_, err = r.db.ExecContext(ctx, query,
		identity.ID, identity.ConnectionID, identity.UserID, string(identity.Provider),
		identity.ExternalUserID, identity.ExternalUsername, identity.Role, permsJSON,
		identity.Active, identity.PairedAt, identity.CreatedAt, identity.UpdatedAt,
	)
	return err
}

func (r *chatAppRepo) GetIdentityByExternal(ctx context.Context, connectionID string, extUserID string) (*models.ChatAppIdentity, error) {
	query := `
		SELECT id, connection_id, user_id, provider, external_user_id, external_username,
		       role, permissions, active, paired_at, created_at, updated_at
		FROM chat_app_identities
		WHERE connection_id = $1 AND external_user_id = $2 AND active = TRUE
	`
	var idt models.ChatAppIdentity
	var provider string
	var permsBytes []byte

	err := r.db.QueryRowContext(ctx, query, connectionID, extUserID).Scan(
		&idt.ID, &idt.ConnectionID, &idt.UserID, &provider, &idt.ExternalUserID,
		&idt.ExternalUsername, &idt.Role, &permsBytes, &idt.Active, &idt.PairedAt,
		&idt.CreatedAt, &idt.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get chat identity: %w", err)
	}
	idt.Provider = models.ChatAppProvider(provider)
	if len(permsBytes) > 0 {
		_ = json.Unmarshal(permsBytes, &idt.Permissions)
	}
	return &idt, nil
}

func (r *chatAppRepo) ListIdentitiesByConnection(ctx context.Context, connectionID string) ([]*models.ChatAppIdentity, error) {
	query := `
		SELECT id, connection_id, user_id, provider, external_user_id, external_username,
		       role, permissions, active, paired_at, created_at, updated_at
		FROM chat_app_identities
		WHERE connection_id = $1 ORDER BY paired_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, connectionID)
	if err != nil {
		return nil, fmt.Errorf("list chat identities: %w", err)
	}
	defer rows.Close()

	var list []*models.ChatAppIdentity
	for rows.Next() {
		var idt models.ChatAppIdentity
		var provider string
		var permsBytes []byte

		if err := rows.Scan(
			&idt.ID, &idt.ConnectionID, &idt.UserID, &provider, &idt.ExternalUserID,
			&idt.ExternalUsername, &idt.Role, &permsBytes, &idt.Active, &idt.PairedAt,
			&idt.CreatedAt, &idt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		idt.Provider = models.ChatAppProvider(provider)
		if len(permsBytes) > 0 {
			_ = json.Unmarshal(permsBytes, &idt.Permissions)
		}
		list = append(list, &idt)
	}
	return list, nil
}

func (r *chatAppRepo) DeleteIdentity(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chat_app_identities WHERE id = $1`, id)
	return err
}

func (r *chatAppRepo) CreateChallenge(ctx context.Context, ch *models.ChatAppChallenge) error {
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

func (r *chatAppRepo) GetChallengeByToken(ctx context.Context, token string) (*models.ChatAppChallenge, error) {
	query := `
		SELECT id, connection_id, provider, external_user_id, external_channel_id,
		       token, code, expires_at, verified_at, created_at
		FROM chat_app_challenges
		WHERE token = $1
	`
	var ch models.ChatAppChallenge
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
	ch.Provider = models.ChatAppProvider(provider)
	return &ch, nil
}

func (r *chatAppRepo) VerifyChallenge(ctx context.Context, token string, code string) (*models.ChatAppChallenge, error) {
	query := `
		SELECT id, connection_id, provider, external_user_id, external_channel_id,
		       token, code, expires_at, verified_at, created_at
		FROM chat_app_challenges
		WHERE (token = $1 OR code = $2) AND expires_at > CURRENT_TIMESTAMP
	`
	var ch models.ChatAppChallenge
	var provider string
	err := r.db.QueryRowContext(ctx, query, token, code).Scan(
		&ch.ID, &ch.ConnectionID, &provider, &ch.ExternalUserID, &ch.ExternalChannelID,
		&ch.Token, &ch.Code, &ch.ExpiresAt, &ch.VerifiedAt, &ch.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("challenge expired or invalid")
		}
		return nil, err
	}
	ch.Provider = models.ChatAppProvider(provider)

	now := time.Now().UTC()
	_, _ = r.db.ExecContext(ctx, `UPDATE chat_app_challenges SET verified_at = $1 WHERE id = $2`, now, ch.ID)
	ch.VerifiedAt = &now
	return &ch, nil
}

func (r *chatAppRepo) CreateCommand(ctx context.Context, cmd *models.ChatAppCommand) error {
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

func (r *chatAppRepo) GetCommandByConfirmationToken(ctx context.Context, token string) (*models.ChatAppCommand, error) {
	query := `
		SELECT id, connection_id, identity_id, provider, command, args, status,
		       response, confirmation_token, confirmed_at, created_at
		FROM chat_app_commands
		WHERE confirmation_token = $1 AND status = 'awaiting_confirmation'
	`
	var cmd models.ChatAppCommand
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
	cmd.Provider = models.ChatAppProvider(provider)
	return &cmd, nil
}

func (r *chatAppRepo) MarkCommandConfirmed(ctx context.Context, id string, response string) error {
	now := time.Now().UTC()
	query := `
		UPDATE chat_app_commands
		SET status = 'executed', confirmed_at = $2, response = $3
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, now, response)
	return err
}
