package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"codedock/internal/models"
)

type ServerRepository interface {
	Create(ctx context.Context, server *models.Server) error
	Update(ctx context.Context, server *models.Server) error
	GetByID(ctx context.Context, id string) (*models.Server, error)
	GetByToken(ctx context.Context, token string) (*models.Server, error)
	ListByUser(ctx context.Context, userID string) ([]*models.Server, error)
	UpdateStatus(ctx context.Context, id string, status models.ServerStatus) error
	UpdateMetrics(ctx context.Context, id string, metricsJSON []byte) error
	Delete(ctx context.Context, id string) error
}

type postgresServerRepository struct {
	db    *sql.DB
	vault Vault
}

func NewServerRepository(db *sql.DB, vault Vault) ServerRepository {
	return &postgresServerRepository{db: db, vault: vault}
}

func (r *postgresServerRepository) Create(ctx context.Context, server *models.Server) error {
	sshKey, err := r.encryptSecret(server.SSHKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt SSH key: %w", err)
	}
	sshPrivateKey, err := r.encryptSecret(server.SSHPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt SSH private key: %w", err)
	}
	sshPassword, err := r.encryptSecret(server.SSHPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt SSH password: %w", err)
	}
	var lastSeenAt any
	if server.LastSeenAt != nil {
		lastSeenAt = *server.LastSeenAt
	}
	query := `
		INSERT INTO servers (
			id, user_id, name, ip_address, is_local,
			ssh_host, ssh_port, ssh_user, ssh_auth_method,
			ssh_key, ssh_private_key, ssh_password, ssh_transport, ssh_jump_host,
			status, provider, external_id, region, server_type, worker_token, last_seen_at, metrics, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
	`
	_, err = r.db.ExecContext(ctx, query,
		server.ID, server.UserID, server.Name, server.IPAddress, server.IsLocal,
		server.SSHHost, server.SSHPort, server.SSHUser, server.SSHAuthMethod,
		sshKey, sshPrivateKey, sshPassword, server.SSHTransport, server.SSHJumpHost,
		server.Status, server.Provider, server.ExternalID, server.Region, server.ServerType, server.WorkerToken, lastSeenAt, string(server.Metrics), server.CreatedAt, server.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	return nil
}

func (r *postgresServerRepository) Update(ctx context.Context, server *models.Server) error {
	sshKey, err := r.encryptSecret(server.SSHKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt SSH key: %w", err)
	}
	sshPrivateKey, err := r.encryptSecret(server.SSHPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt SSH private key: %w", err)
	}
	sshPassword, err := r.encryptSecret(server.SSHPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt SSH password: %w", err)
	}
	query := `
		UPDATE servers SET
			name = $1,
			ip_address = $2,
			is_local = $3,
			ssh_host = $4,
			ssh_port = $5,
			ssh_user = $6,
			ssh_auth_method = $7,
			ssh_key = $8,
			ssh_private_key = $9,
			ssh_password = $10,
			ssh_transport = $11,
			ssh_jump_host = $12,
			status = $13,
			provider = $14,
			external_id = $15,
			region = $16,
			server_type = $17,
			updated_at = $18
		WHERE id = $19
	`
	_, err = r.db.ExecContext(ctx, query,
		server.Name, server.IPAddress, server.IsLocal,
		server.SSHHost, server.SSHPort, server.SSHUser, server.SSHAuthMethod,
		sshKey, sshPrivateKey, sshPassword, server.SSHTransport, server.SSHJumpHost,
		server.Status, server.Provider, server.ExternalID, server.Region, server.ServerType, server.UpdatedAt, server.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update server: %w", err)
	}
	return nil
}

const serverSelectColumns = `
	id, user_id, name, ip_address, COALESCE(is_local, FALSE),
	COALESCE(ssh_host, ''), COALESCE(ssh_port, 22), COALESCE(ssh_user, 'root'),
	COALESCE(ssh_auth_method, 'key'), COALESCE(ssh_key, ''), COALESCE(ssh_private_key, ''),
	COALESCE(ssh_password, ''), COALESCE(ssh_transport, 'direct'), COALESCE(ssh_jump_host, ''),
	status, COALESCE(provider, ''), COALESCE(external_id, ''), COALESCE(region, ''), COALESCE(server_type, ''),
	worker_token, last_seen_at, metrics, created_at, updated_at
`

func (r *postgresServerRepository) GetByID(ctx context.Context, id string) (*models.Server, error) {
	query := fmt.Sprintf(`SELECT %s FROM servers WHERE id = $1`, serverSelectColumns)
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanRow(row)
}

func (r *postgresServerRepository) GetByToken(ctx context.Context, token string) (*models.Server, error) {
	query := fmt.Sprintf(`SELECT %s FROM servers WHERE worker_token = $1`, serverSelectColumns)
	row := r.db.QueryRowContext(ctx, query, token)
	return r.scanRow(row)
}

func (r *postgresServerRepository) ListByUser(ctx context.Context, userID string) ([]*models.Server, error) {
	query := fmt.Sprintf(`SELECT %s FROM servers WHERE user_id = $1 OR (user_id = 'system' AND is_local = TRUE) ORDER BY created_at DESC`, serverSelectColumns)
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list servers: %w", err)
	}
	defer rows.Close()

	var servers []*models.Server
	for rows.Next() {
		server, err := r.scanRows(rows)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, nil
}

func (r *postgresServerRepository) UpdateStatus(ctx context.Context, id string, status models.ServerStatus) error {
	query := `UPDATE servers SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update server status: %w", err)
	}
	return nil
}

func (r *postgresServerRepository) UpdateMetrics(ctx context.Context, id string, metricsJSON []byte) error {
	query := `UPDATE servers SET metrics = $1, last_seen_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, string(metricsJSON), id)
	if err != nil {
		return fmt.Errorf("failed to update server metrics: %w", err)
	}
	return nil
}

func (r *postgresServerRepository) scanRow(row *sql.Row) (*models.Server, error) {
	var s models.Server
	var metricsRaw sql.NullString
	err := row.Scan(
		&s.ID, &s.UserID, &s.Name, &s.IPAddress, &s.IsLocal,
		&s.SSHHost, &s.SSHPort, &s.SSHUser, &s.SSHAuthMethod,
		&s.SSHKey, &s.SSHPrivateKey, &s.SSHPassword, &s.SSHTransport, &s.SSHJumpHost,
		&s.Status, &s.Provider, &s.ExternalID, &s.Region, &s.ServerType,
		&s.WorkerToken, &s.LastSeenAt, &metricsRaw, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan server: %w", err)
	}
	if metricsRaw.Valid && metricsRaw.String != "" {
		s.Metrics = json.RawMessage(metricsRaw.String)
	} else {
		s.Metrics = json.RawMessage("{}")
	}
	if err := r.decryptSecrets(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *postgresServerRepository) scanRows(rows *sql.Rows) (*models.Server, error) {
	var s models.Server
	var metricsRaw sql.NullString
	err := rows.Scan(
		&s.ID, &s.UserID, &s.Name, &s.IPAddress, &s.IsLocal,
		&s.SSHHost, &s.SSHPort, &s.SSHUser, &s.SSHAuthMethod,
		&s.SSHKey, &s.SSHPrivateKey, &s.SSHPassword, &s.SSHTransport, &s.SSHJumpHost,
		&s.Status, &s.Provider, &s.ExternalID, &s.Region, &s.ServerType,
		&s.WorkerToken, &s.LastSeenAt, &metricsRaw, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan server: %w", err)
	}
	if metricsRaw.Valid && metricsRaw.String != "" {
		s.Metrics = json.RawMessage(metricsRaw.String)
	} else {
		s.Metrics = json.RawMessage("{}")
	}
	if err := r.decryptSecrets(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *postgresServerRepository) encryptSecret(value string) (string, error) {
	if value == "" || r.vault == nil {
		return value, nil
	}
	if _, err := r.vault.Decrypt(value); err == nil {
		return value, nil
	}
	return r.vault.Encrypt(value)
}

func (r *postgresServerRepository) decryptSecrets(server *models.Server) error {
	if r.vault == nil {
		return nil
	}
	if server.SSHKey != "" {
		if value, err := r.vault.Decrypt(server.SSHKey); err == nil {
			server.SSHKey = value
		}
	}
	if server.SSHPrivateKey != "" {
		if value, err := r.vault.Decrypt(server.SSHPrivateKey); err == nil {
			server.SSHPrivateKey = value
		}
	}
	if server.SSHPassword != "" {
		if value, err := r.vault.Decrypt(server.SSHPassword); err == nil {
			server.SSHPassword = value
		}
	}
	return nil
}

func (r *postgresServerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM servers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
