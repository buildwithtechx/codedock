package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
)

type MailRepository interface {
	Create(ctx context.Context, server *models.MailServer) error
	GetByID(ctx context.Context, id string) (*models.MailServer, error)
	ListByOrganization(ctx context.Context, orgID string) ([]*models.MailServer, error)
	Delete(ctx context.Context, id string) error
}

type MailRepo struct {
	db    *sqlx.DB
	mu    sync.Mutex
	vault Vault
}

func NewMailRepo(db *sql.DB, v Vault) *MailRepo {
	return &MailRepo{db: sqlx.NewDb(db, "pgx"), vault: v}
}

func (r *MailRepo) Create(ctx context.Context, server *models.MailServer) error {
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

func (r *MailRepo) GetByID(ctx context.Context, id string) (*models.MailServer, error) {
	var server models.MailServer
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

func (r *MailRepo) ListByOrganization(ctx context.Context, orgID string) ([]*models.MailServer, error) {
	servers := []*models.MailServer{}
	if err := r.db.SelectContext(ctx, &servers, `SELECT id, organization_id, server_id, domain, maild_url, webmail_url, api_token, created_at, updated_at FROM mail_servers WHERE organization_id = $1 ORDER BY domain`, orgID); err != nil {
		return nil, fmt.Errorf("failed to list mail servers: %w", err)
	}
	for _, server := range servers {
		server.APIToken = ""
	}
	return servers, nil
}

func (r *MailRepo) Delete(ctx context.Context, id string) error {
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
