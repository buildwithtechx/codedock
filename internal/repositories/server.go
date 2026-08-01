package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"codedock.run/codedock/internal/models"
)

type ServerRepository interface {
	Create(ctx context.Context, server *models.Server) error
	GetByID(ctx context.Context, id string) (*models.Server, error)
	GetByToken(ctx context.Context, token string) (*models.Server, error)
	ListByUser(ctx context.Context, userID string) ([]*models.Server, error)
	UpdateStatus(ctx context.Context, id string, status models.ServerStatus) error
	UpdateMetrics(ctx context.Context, id string, metricsJSON []byte) error
	Delete(ctx context.Context, id string) error
}

type sqliteServerRepository struct {
	db    *sql.DB
	vault Vault
}

func NewServerRepository(db *sql.DB, vault Vault) ServerRepository {
	return &sqliteServerRepository{db: db, vault: vault}
}

func (r *sqliteServerRepository) Create(ctx context.Context, server *models.Server) error {
	sshKey, err := r.encryptSecret(server.SSHKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt SSH key: %w", err)
	}
	sshPassword, err := r.encryptSecret(server.SSHPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt SSH password: %w", err)
	}
	query := `
		INSERT INTO servers (id, user_id, name, ip_address, ssh_host, ssh_port, ssh_user, ssh_key, ssh_password, status, worker_token, last_seen_at, metrics, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = r.db.ExecContext(ctx, query,
		server.ID, server.UserID, server.Name, server.IPAddress,
		server.SSHHost, server.SSHPort, server.SSHUser, sshKey, sshPassword,
		server.Status, server.WorkerToken, server.LastSeenAt, server.Metrics, server.CreatedAt, server.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}
	return nil
}

func (r *sqliteServerRepository) GetByID(ctx context.Context, id string) (*models.Server, error) {
	query := `SELECT id, user_id, name, ip_address, COALESCE(ssh_host, ''), COALESCE(ssh_port, 22), COALESCE(ssh_user, 'root'), COALESCE(ssh_key, ''), COALESCE(ssh_password, ''), status, worker_token, last_seen_at, metrics, created_at, updated_at FROM servers WHERE id = ?`
	row := r.db.QueryRowContext(ctx, query, id)
	return r.scanRow(row)
}

func (r *sqliteServerRepository) GetByToken(ctx context.Context, token string) (*models.Server, error) {
	query := `SELECT id, user_id, name, ip_address, COALESCE(ssh_host, ''), COALESCE(ssh_port, 22), COALESCE(ssh_user, 'root'), COALESCE(ssh_key, ''), COALESCE(ssh_password, ''), status, worker_token, last_seen_at, metrics, created_at, updated_at FROM servers WHERE worker_token = ?`
	row := r.db.QueryRowContext(ctx, query, token)
	return r.scanRow(row)
}

func (r *sqliteServerRepository) ListByUser(ctx context.Context, userID string) ([]*models.Server, error) {
	query := `SELECT id, user_id, name, ip_address, COALESCE(ssh_host, ''), COALESCE(ssh_port, 22), COALESCE(ssh_user, 'root'), COALESCE(ssh_key, ''), COALESCE(ssh_password, ''), status, worker_token, last_seen_at, metrics, created_at, updated_at FROM servers WHERE user_id = ? ORDER BY created_at DESC`
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

func (r *sqliteServerRepository) UpdateStatus(ctx context.Context, id string, status models.ServerStatus) error {
	query := `UPDATE servers SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update server status: %w", err)
	}
	return nil
}

func (r *sqliteServerRepository) UpdateMetrics(ctx context.Context, id string, metricsJSON []byte) error {
	query := `UPDATE servers SET metrics = ?, last_seen_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, metricsJSON, id)
	if err != nil {
		return fmt.Errorf("failed to update server metrics: %w", err)
	}
	return nil
}

func (r *sqliteServerRepository) scanRow(row *sql.Row) (*models.Server, error) {
	var s models.Server
	err := row.Scan(&s.ID, &s.UserID, &s.Name, &s.IPAddress, &s.SSHHost, &s.SSHPort, &s.SSHUser, &s.SSHKey, &s.SSHPassword, &s.Status, &s.WorkerToken, &s.LastSeenAt, &s.Metrics, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan server: %w", err)
	}
	if err := r.decryptSecrets(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sqliteServerRepository) scanRows(rows *sql.Rows) (*models.Server, error) {
	var s models.Server
	err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.IPAddress, &s.SSHHost, &s.SSHPort, &s.SSHUser, &s.SSHKey, &s.SSHPassword, &s.Status, &s.WorkerToken, &s.LastSeenAt, &s.Metrics, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan server: %w", err)
	}
	if err := r.decryptSecrets(&s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *sqliteServerRepository) encryptSecret(value string) (string, error) {
	if value == "" || r.vault == nil {
		return value, nil
	}
	if _, err := r.vault.Decrypt(value); err == nil {
		return value, nil
	}
	return r.vault.Encrypt(value)
}

func (r *sqliteServerRepository) decryptSecrets(server *models.Server) error {
	if r.vault == nil {
		return nil
	}
	if server.SSHKey != "" {
		value, err := r.vault.Decrypt(server.SSHKey)
		if err == nil {
			server.SSHKey = value
		}
	}
	if server.SSHPassword != "" {
		value, err := r.vault.Decrypt(server.SSHPassword)
		if err == nil {
			server.SSHPassword = value
		}
	}
	return nil
}

func (r *sqliteServerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM servers WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
