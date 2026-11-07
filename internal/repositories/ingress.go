package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"codedock/internal/models"
)

type IngressRepository interface {
	ListConnections(ctx context.Context, serverID string) ([]*models.IngressConnection, error)
	GetConnection(ctx context.Context, id string) (*models.IngressConnection, error)
	SaveConnection(ctx context.Context, conn *models.IngressConnection) error
	DeleteConnection(ctx context.Context, id string) error

	GetHostname(ctx context.Context, domainID string) (*models.IngressHostname, error)
	SaveHostname(ctx context.Context, host *models.IngressHostname) error
	DeleteHostname(ctx context.Context, domainID string) error
}

type sqlIngressRepo struct {
	db *sql.DB
	mu sync.RWMutex
	// In-memory fallback if db is nil
	memConns map[string]*models.IngressConnection
	memHosts map[string]*models.IngressHostname
}

func NewIngressRepository(db *sql.DB) IngressRepository {
	return &sqlIngressRepo{
		db:       db,
		memConns: make(map[string]*models.IngressConnection),
		memHosts: make(map[string]*models.IngressHostname),
	}
}

func (r *sqlIngressRepo) ListConnections(ctx context.Context, serverID string) ([]*models.IngressConnection, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		var res []*models.IngressConnection
		for _, c := range r.memConns {
			if serverID == "" || c.ServerID == serverID {
				res = append(res, c)
			}
		}
		if res == nil {
			res = make([]*models.IngressConnection, 0)
		}
		return res, nil
	}

	query := `SELECT id, server_id, credential_id, account_id, tunnel_id, name, connector, edge_port, managed_tunnel, default_for_domains, default_access_policy_id, status, error, created_at, updated_at FROM ingress_connections`
	var rows *sql.Rows
	var err error
	if serverID != "" {
		query += ` WHERE server_id = $1`
		rows, err = r.db.QueryContext(ctx, query, serverID)
	} else {
		rows, err = r.db.QueryContext(ctx, query)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.IngressConnection
	for rows.Next() {
		var c models.IngressConnection
		var tCreatedAt, tUpdatedAt time.Time
		var tunnelID, defaultPolicy, errText sql.NullString
		if err := rows.Scan(&c.ID, &c.ServerID, &c.CredentialID, &c.AccountID, &tunnelID, &c.Name, &c.Connector, &c.EdgePort, &c.ManagedTunnel, &c.DefaultForDomains, &defaultPolicy, &c.Status, &errText, &tCreatedAt, &tUpdatedAt); err != nil {
			return nil, err
		}
		if tunnelID.Valid {
			c.TunnelID = &tunnelID.String
		}
		if defaultPolicy.Valid {
			c.DefaultAccessPolicyID = &defaultPolicy.String
		}
		if errText.Valid {
			c.Error = &errText.String
		}
		c.CreatedAt = tCreatedAt.Format(time.RFC3339)
		c.UpdatedAt = tUpdatedAt.Format(time.RFC3339)
		list = append(list, &c)
	}
	if list == nil {
		list = make([]*models.IngressConnection, 0)
	}
	return list, nil
}

func (r *sqlIngressRepo) GetConnection(ctx context.Context, id string) (*models.IngressConnection, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		c, ok := r.memConns[id]
		if !ok {
			return nil, fmt.Errorf("connection not found")
		}
		return c, nil
	}

	query := `SELECT id, server_id, credential_id, account_id, tunnel_id, name, connector, edge_port, managed_tunnel, default_for_domains, default_access_policy_id, status, error, created_at, updated_at FROM ingress_connections WHERE id = $1`
	var c models.IngressConnection
	var tCreatedAt, tUpdatedAt time.Time
	var tunnelID, defaultPolicy, errText sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(&c.ID, &c.ServerID, &c.CredentialID, &c.AccountID, &tunnelID, &c.Name, &c.Connector, &c.EdgePort, &c.ManagedTunnel, &c.DefaultForDomains, &defaultPolicy, &c.Status, &errText, &tCreatedAt, &tUpdatedAt)
	if err != nil {
		return nil, err
	}
	if tunnelID.Valid {
		c.TunnelID = &tunnelID.String
	}
	if defaultPolicy.Valid {
		c.DefaultAccessPolicyID = &defaultPolicy.String
	}
	if errText.Valid {
		c.Error = &errText.String
	}
	c.CreatedAt = tCreatedAt.Format(time.RFC3339)
	c.UpdatedAt = tUpdatedAt.Format(time.RFC3339)
	return &c, nil
}

func (r *sqlIngressRepo) SaveConnection(ctx context.Context, conn *models.IngressConnection) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		r.memConns[conn.ID] = conn
		return nil
	}

	query := `INSERT INTO ingress_connections (id, server_id, credential_id, account_id, tunnel_id, name, connector, edge_port, managed_tunnel, default_for_domains, default_access_policy_id, status, error, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET
	default_for_domains = EXCLUDED.default_for_domains,
	default_access_policy_id = EXCLUDED.default_access_policy_id,
	status = EXCLUDED.status,
	error = EXCLUDED.error,
	updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, query, conn.ID, conn.ServerID, conn.CredentialID, conn.AccountID, conn.TunnelID, conn.Name, conn.Connector, conn.EdgePort, conn.ManagedTunnel, conn.DefaultForDomains, conn.DefaultAccessPolicyID, conn.Status, conn.Error)
	return err
}

func (r *sqlIngressRepo) DeleteConnection(ctx context.Context, id string) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.memConns, id)
		return nil
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM ingress_connections WHERE id = $1`, id)
	return err
}

func (r *sqlIngressRepo) GetHostname(ctx context.Context, domainID string) (*models.IngressHostname, error) {
	if r.db == nil {
		r.mu.RLock()
		defer r.mu.RUnlock()
		h, ok := r.memHosts[domainID]
		if !ok {
			return nil, nil
		}
		return h, nil
	}

	query := `SELECT id, domain_id, connection_id, hostname, access_policy_id, status, phase, error, updated_at FROM ingress_hostnames WHERE domain_id = $1`
	var h models.IngressHostname
	var dID, polID, errText sql.NullString
	var tUpdatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, domainID).Scan(&h.ID, &dID, &h.ConnectionID, &h.Hostname, &polID, &h.Status, &h.Phase, &errText, &tUpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if dID.Valid {
		h.DomainID = &dID.String
	}
	if polID.Valid {
		h.AccessPolicyID = &polID.String
	}
	if errText.Valid {
		h.Error = &errText.String
	}
	h.UpdatedAt = tUpdatedAt.Format(time.RFC3339)
	return &h, nil
}

func (r *sqlIngressRepo) SaveHostname(ctx context.Context, host *models.IngressHostname) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		if host.DomainID != nil {
			r.memHosts[*host.DomainID] = host
		}
		return nil
	}

	query := `INSERT INTO ingress_hostnames (id, domain_id, connection_id, hostname, access_policy_id, status, phase, error, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
ON CONFLICT (id) DO UPDATE SET
	connection_id = EXCLUDED.connection_id,
	access_policy_id = EXCLUDED.access_policy_id,
	status = EXCLUDED.status,
	phase = EXCLUDED.phase,
	error = EXCLUDED.error,
	updated_at = NOW()`
	_, err := r.db.ExecContext(ctx, query, host.ID, host.DomainID, host.ConnectionID, host.Hostname, host.AccessPolicyID, host.Status, host.Phase, host.Error)
	return err
}

func (r *sqlIngressRepo) DeleteHostname(ctx context.Context, domainID string) error {
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()
		delete(r.memHosts, domainID)
		return nil
	}
	_, err := r.db.ExecContext(ctx, `DELETE FROM ingress_hostnames WHERE domain_id = $1`, domainID)
	return err
}
