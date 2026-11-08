package mail

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type Repository interface {
	Create(ctx context.Context, server *MailServer) error
	GetByID(ctx context.Context, id string) (*MailServer, error)
	ListByOrganization(ctx context.Context, orgID string) ([]*MailServer, error)
	Delete(ctx context.Context, id string) error
}

type repo struct {
	db    *sqlx.DB
	mu    sync.Mutex
	vault Vault
}

func NewRepository(db *sql.DB, v Vault) Repository {
	return &repo{db: sqlx.NewDb(db, "pgx"), vault: v}
}

func (r *repo) Create(ctx context.Context, server *MailServer) error {
	if server.ID == "" {
		server.ID = uuid.New().String()
	}
	now := time.Now().UTC().Format(time.RFC3339)
	server.CreatedAt = now
	server.UpdatedAt = now
	if server.APIToken != "" && r.vault != nil {
		enc, err := r.vault.Encrypt(server.APIToken)
		if err != nil {
			return fmt.Errorf("failed to encrypt mail api token: %w", err)
		}
		server.APIToken = enc
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `INSERT INTO mail_servers (id, organization_id, server_id, domain, maild_url, webmail_url, api_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		server.ID, server.OrganizationID, server.ServerID, server.Domain, server.MaildURL, server.WebmailURL, server.APIToken, server.CreatedAt, server.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create mail server: %w", err)
	}
	server.APIToken = ""
	return nil
}

func (r *repo) GetByID(ctx context.Context, id string) (*MailServer, error) {
	var server MailServer
	if err := r.db.GetContext(ctx, &server, `SELECT id, organization_id, server_id, domain, maild_url, webmail_url, api_token, created_at, updated_at FROM mail_servers WHERE id = $1`, id); err != nil {
		return nil, fmt.Errorf("failed to get mail server: %w", err)
	}
	if server.APIToken != "" && r.vault != nil {
		dec, err := r.vault.Decrypt(server.APIToken)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt mail api token: %w", err)
		}
		server.APIToken = dec
	}
	return &server, nil
}

func (r *repo) ListByOrganization(ctx context.Context, orgID string) ([]*MailServer, error) {
	servers := []*MailServer{}
	if err := r.db.SelectContext(ctx, &servers, `SELECT id, organization_id, server_id, domain, maild_url, webmail_url, api_token, created_at, updated_at FROM mail_servers WHERE organization_id = $1 ORDER BY domain`, orgID); err != nil {
		return nil, fmt.Errorf("failed to list mail servers: %w", err)
	}
	for _, server := range servers {
		server.APIToken = ""
	}
	return servers, nil
}

func (r *repo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.ExecContext(ctx, `DELETE FROM mail_servers WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete mail server: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to read delete result: %w", err)
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
