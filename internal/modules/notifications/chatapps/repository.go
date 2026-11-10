package chatapps

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Repository interface {
	CreateConnection(ctx context.Context, conn *ChatAppConnection) error
	GetConnection(ctx context.Context, id string) (*ChatAppConnection, error)
	ListConnectionsByUser(ctx context.Context, userID string) ([]*ChatAppConnection, error)
	UpdateConnection(ctx context.Context, conn *ChatAppConnection) error
	DeleteConnection(ctx context.Context, id string) error

	CreateIdentity(ctx context.Context, identity *ChatAppIdentity) error
	GetIdentityByExternal(ctx context.Context, connectionID string, extUserID string) (*ChatAppIdentity, error)
	ListIdentitiesByConnection(ctx context.Context, connectionID string) ([]*ChatAppIdentity, error)
	DeleteIdentity(ctx context.Context, id string) error

	CreateChallenge(ctx context.Context, ch *ChatAppChallenge) error
	GetChallengeByToken(ctx context.Context, token string) (*ChatAppChallenge, error)
	VerifyChallenge(ctx context.Context, token string, code string) (*ChatAppChallenge, error)

	CreateCommand(ctx context.Context, cmd *ChatAppCommand) error
	GetCommandByConfirmationToken(ctx context.Context, token string) (*ChatAppCommand, error)
	MarkCommandConfirmed(ctx context.Context, id string, response string) error
}

type repo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repo{db: db}
}

func (r *repo) CreateConnection(ctx context.Context, conn *ChatAppConnection) error {
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

func (r *repo) GetConnection(ctx context.Context, id string) (*ChatAppConnection, error) {
	query := `
		SELECT id, user_id, provider, transport, purpose, name, bot_username,
		       credentials_encrypted, config, enabled, status, last_error, last_event_at, created_at, updated_at
		FROM chat_app_connections WHERE id = $1
	`
	var conn ChatAppConnection
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

	conn.Provider = ChatAppProvider(provider)
	conn.Transport = ChatAppTransport(transport)
	conn.Purpose = ChatAppPurpose(purpose)
	conn.Status = ChatAppStatus(status)
	if len(configBytes) > 0 {
		_ = json.Unmarshal(configBytes, &conn.Config)
	}
	if conn.Config == nil {
		conn.Config = make(map[string]any)
	}
	return &conn, nil
}

func (r *repo) ListConnectionsByUser(ctx context.Context, userID string) ([]*ChatAppConnection, error) {
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

	var list []*ChatAppConnection
	for rows.Next() {
		var conn ChatAppConnection
		var provider, transport, purpose, status string
		var configBytes []byte

		if err := rows.Scan(
			&conn.ID, &conn.UserID, &provider, &transport, &purpose, &conn.Name, &conn.BotUsername,
			&conn.CredentialsEncrypted, &configBytes, &conn.Enabled, &status, &conn.LastError,
			&conn.LastEventAt, &conn.CreatedAt, &conn.UpdatedAt,
		); err != nil {
			return nil, err
		}
		conn.Provider = ChatAppProvider(provider)
		conn.Transport = ChatAppTransport(transport)
		conn.Purpose = ChatAppPurpose(purpose)
		conn.Status = ChatAppStatus(status)
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

func (r *repo) UpdateConnection(ctx context.Context, conn *ChatAppConnection) error {
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

func (r *repo) DeleteConnection(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chat_app_connections WHERE id = $1`, id)
	return err
}

func (r *repo) CreateIdentity(ctx context.Context, identity *ChatAppIdentity) error {
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

func (r *repo) GetIdentityByExternal(ctx context.Context, connectionID string, extUserID string) (*ChatAppIdentity, error) {
	query := `
		SELECT id, connection_id, user_id, provider, external_user_id, external_username,
		       role, permissions, active, paired_at, created_at, updated_at
		FROM chat_app_identities
		WHERE connection_id = $1 AND external_user_id = $2 AND active = TRUE
	`
	var idt ChatAppIdentity
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
	idt.Provider = ChatAppProvider(provider)
	if len(permsBytes) > 0 {
		_ = json.Unmarshal(permsBytes, &idt.Permissions)
	}
	return &idt, nil
}

func (r *repo) ListIdentitiesByConnection(ctx context.Context, connectionID string) ([]*ChatAppIdentity, error) {
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

	var list []*ChatAppIdentity
	for rows.Next() {
		var idt ChatAppIdentity
		var provider string
		var permsBytes []byte

		if err := rows.Scan(
			&idt.ID, &idt.ConnectionID, &idt.UserID, &provider, &idt.ExternalUserID,
			&idt.ExternalUsername, &idt.Role, &permsBytes, &idt.Active, &idt.PairedAt,
			&idt.CreatedAt, &idt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		idt.Provider = ChatAppProvider(provider)
		if len(permsBytes) > 0 {
			_ = json.Unmarshal(permsBytes, &idt.Permissions)
		}
		list = append(list, &idt)
	}
	return list, nil
}

func (r *repo) DeleteIdentity(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM chat_app_identities WHERE id = $1`, id)
	return err
}
