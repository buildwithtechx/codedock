package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"codedock.run/codedock/internal/models"
)

const (
	defaultManagedMaxServers  = 5
	defaultManagedMaxMemoryGB = 32
)

type ManagedRepository interface {
	CreateCredential(ctx context.Context, credential *models.ManagedCredential) error
	GetCredential(ctx context.Context, id string) (*models.ManagedCredential, error)
	ListCredentialsByOrg(ctx context.Context, orgID string) ([]*models.ManagedCredential, error)
	DeleteCredential(ctx context.Context, id string) error
	GetQuota(ctx context.Context, orgID string) (*models.ManagedQuota, error)
	SetQuota(ctx context.Context, quota *models.ManagedQuota) error
	SaveLink(ctx context.Context, link *models.ManagedServerLink) error
	GetLink(ctx context.Context, serverID string) (*models.ManagedServerLink, error)
	DeleteLink(ctx context.Context, serverID string) error
	ListOrgServerIDs(ctx context.Context, orgID string) ([]string, error)
	ListProvisioningServerIDs(ctx context.Context) ([]string, error)
}

type postgresManagedRepository struct {
	db    *sql.DB
	vault Vault
}

func NewManagedRepository(db *sql.DB, vault Vault) ManagedRepository {
	return &postgresManagedRepository{db: db, vault: vault}
}

func (r *postgresManagedRepository) CreateCredential(ctx context.Context, credential *models.ManagedCredential) error {
	token, err := r.encryptSecret(credential.Token)
	if err != nil {
		return fmt.Errorf("failed to encrypt provider token: %w", err)
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO managed_providers (id, organization_id, provider, label, encrypted_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, credential.ID, credential.OrganizationID, string(credential.Provider), credential.Label, token)
	if err != nil {
		return fmt.Errorf("failed to create managed credential: %w", err)
	}
	return nil
}

func (r *postgresManagedRepository) GetCredential(ctx context.Context, id string) (*models.ManagedCredential, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, organization_id, provider, label, encrypted_token, created_at, updated_at
		FROM managed_providers WHERE id = $1
	`, id)
	return r.scanCredential(row)
}

func (r *postgresManagedRepository) ListCredentialsByOrg(ctx context.Context, orgID string) ([]*models.ManagedCredential, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, organization_id, provider, label, encrypted_token, created_at, updated_at
		FROM managed_providers WHERE organization_id = $1 ORDER BY created_at DESC
	`, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list managed credentials: %w", err)
	}
	defer rows.Close()
	var credentials []*models.ManagedCredential
	for rows.Next() {
		var credential models.ManagedCredential
		var provider string
		if err := rows.Scan(&credential.ID, &credential.OrganizationID, &provider, &credential.Label, &credential.Token, &credential.CreatedAt, &credential.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan managed credential: %w", err)
		}
		credential.Provider = models.ManagedProvider(provider)
		credential.Token = ""
		credentials = append(credentials, &credential)
	}
	return credentials, rows.Err()
}

func (r *postgresManagedRepository) DeleteCredential(ctx context.Context, id string) error {
	var linked int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM managed_servers WHERE credential_id = $1`, id).Scan(&linked); err != nil {
		return fmt.Errorf("failed to check credential usage: %w", err)
	}
	if linked > 0 {
		return fmt.Errorf("credential still backs %d managed server(s)", linked)
	}
	if _, err := r.db.ExecContext(ctx, `DELETE FROM managed_providers WHERE id = $1`, id); err != nil {
		return fmt.Errorf("failed to delete managed credential: %w", err)
	}
	return nil
}

func (r *postgresManagedRepository) GetQuota(ctx context.Context, orgID string) (*models.ManagedQuota, error) {
	quota := &models.ManagedQuota{OrganizationID: orgID, MaxServers: defaultManagedMaxServers, MaxMemoryGB: defaultManagedMaxMemoryGB}
	err := r.db.QueryRowContext(ctx, `
		SELECT max_servers, max_memory_gb, updated_at FROM managed_quotas WHERE organization_id = $1
	`, orgID).Scan(&quota.MaxServers, &quota.MaxMemoryGB, &quota.UpdatedAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to load managed quota: %w", err)
	}
	return quota, nil
}

func (r *postgresManagedRepository) SetQuota(ctx context.Context, quota *models.ManagedQuota) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO managed_quotas (organization_id, max_servers, max_memory_gb, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT(organization_id) DO UPDATE SET max_servers = excluded.max_servers, max_memory_gb = excluded.max_memory_gb, updated_at = CURRENT_TIMESTAMP
	`, quota.OrganizationID, quota.MaxServers, quota.MaxMemoryGB)
	if err != nil {
		return fmt.Errorf("failed to save managed quota: %w", err)
	}
	return nil
}

func (r *postgresManagedRepository) SaveLink(ctx context.Context, link *models.ManagedServerLink) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO managed_servers (server_id, organization_id, credential_id, ssh_key_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(server_id) DO UPDATE SET credential_id = excluded.credential_id, ssh_key_name = excluded.ssh_key_name, updated_at = CURRENT_TIMESTAMP
	`, link.ServerID, link.OrganizationID, link.CredentialID, link.SSHKeyName)
	if err != nil {
		return fmt.Errorf("failed to save managed server link: %w", err)
	}
	return nil
}

func (r *postgresManagedRepository) GetLink(ctx context.Context, serverID string) (*models.ManagedServerLink, error) {
	var link models.ManagedServerLink
	err := r.db.QueryRowContext(ctx, `
		SELECT server_id, organization_id, credential_id, ssh_key_name, created_at, updated_at
		FROM managed_servers WHERE server_id = $1
	`, serverID).Scan(&link.ServerID, &link.OrganizationID, &link.CredentialID, &link.SSHKeyName, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to load managed server link: %w", err)
	}
	return &link, nil
}

func (r *postgresManagedRepository) DeleteLink(ctx context.Context, serverID string) error {
	if _, err := r.db.ExecContext(ctx, `DELETE FROM managed_servers WHERE server_id = $1`, serverID); err != nil {
		return fmt.Errorf("failed to delete managed server link: %w", err)
	}
	return nil
}

func (r *postgresManagedRepository) ListOrgServerIDs(ctx context.Context, orgID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT server_id FROM managed_servers WHERE organization_id = $1`, orgID)
	if err != nil {
		return nil, fmt.Errorf("failed to list managed servers: %w", err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan managed server id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *postgresManagedRepository) ListProvisioningServerIDs(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT m.server_id FROM managed_servers m
		JOIN servers s ON s.id = m.server_id
		WHERE s.status = $1
	`, string(models.ServerStatusProvisioning))
	if err != nil {
		return nil, fmt.Errorf("failed to list provisioning servers: %w", err)
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan provisioning server id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *postgresManagedRepository) scanCredential(row *sql.Row) (*models.ManagedCredential, error) {
	var credential models.ManagedCredential
	var provider string
	if err := row.Scan(&credential.ID, &credential.OrganizationID, &provider, &credential.Label, &credential.Token, &credential.CreatedAt, &credential.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan managed credential: %w", err)
	}
	credential.Provider = models.ManagedProvider(provider)
	if credential.Token != "" && r.vault != nil {
		if plain, err := r.vault.Decrypt(credential.Token); err == nil {
			credential.Token = plain
		}
	}
	return &credential, nil
}

func (r *postgresManagedRepository) encryptSecret(value string) (string, error) {
	if value == "" || r.vault == nil {
		return value, nil
	}
	if _, err := r.vault.Decrypt(value); err == nil {
		return value, nil
	}
	return r.vault.Encrypt(value)
}
