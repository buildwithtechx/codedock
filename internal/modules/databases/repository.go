package databases

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/platform/vault"
)

type sqlRepository struct {
	db    *sqlx.DB
	vault *vault.Vault
}

func NewRepository(db *sql.DB, v *vault.Vault) Repository {
	if db == nil {
		return &sqlRepository{vault: v}
	}
	return &sqlRepository{db: sqlx.NewDb(db, "pgx"), vault: v}
}

func containsLikePattern(term string) string {
	escaped := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(term)
	return "%" + escaped + "%"
}

func (r *sqlRepository) Create(ctx context.Context, d *Database) error {
	if d.Password != "" && r.vault != nil {
		enc, err := r.vault.Encrypt(d.Password)
		if err != nil {
			return fmt.Errorf("failed to encrypt database password: %w", err)
		}
		d.EncryptedPassword = enc
	}

	const query = `
		INSERT INTO databases (
			id, project_id, environment_id, name, engine, version,
			username, encrypted_password, database_name, internal_host, port, external_port,
			internal_dns, external_dns, volume_path, container_id, status,
			custom_args, logical_replication, cpu_limit, memory_limit,
			created_at, updated_at
		) VALUES (
			:id, :project_id, :environment_id, :name, :engine, :version,
			:username, :encrypted_password, :database_name, :internal_host, :port, :external_port,
			:internal_dns, :external_dns, :volume_path, :container_id, :status,
			:custom_args, :logical_replication, :cpu_limit, :memory_limit,
			:created_at, :updated_at
		)`
	_, err := r.db.NamedExecContext(ctx, query, d)
	return err
}

const selectDatabaseCols = `
	id, project_id, COALESCE(environment_id, '') as environment_id, name, engine, version,
	username, encrypted_password, database_name, internal_host, port, external_port,
	COALESCE(internal_dns, '') as internal_dns, COALESCE(external_dns, '') as external_dns,
	COALESCE(volume_path, '') as volume_path, COALESCE(container_id, '') as container_id, status,
	COALESCE(custom_args, '') as custom_args, COALESCE(logical_replication, FALSE) as logical_replication,
	COALESCE(cpu_limit, 0) as cpu_limit, COALESCE(memory_limit, 0) as memory_limit,
	created_at, updated_at
`

func (r *sqlRepository) decryptPassword(d *Database) {
	if d.EncryptedPassword != "" && r.vault != nil {
		if dec, err := r.vault.Decrypt(d.EncryptedPassword); err == nil {
			d.Password = dec
		}
	}
}

func (r *sqlRepository) GetByID(ctx context.Context, id string) (*Database, error) {
	var d Database
	query := `SELECT ` + selectDatabaseCols + ` FROM databases WHERE id = $1`
	if err := r.db.GetContext(ctx, &d, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDatabaseNotFound
		}
		return nil, err
	}
	r.decryptPassword(&d)
	return &d, nil
}

func (r *sqlRepository) ListAll(ctx context.Context, search string) ([]*Database, error) {
	var list []*Database
	search = strings.TrimSpace(search)
	pattern := containsLikePattern(search)
	query := `
		SELECT ` + selectDatabaseCols + `
		FROM databases
		WHERE ($1 = '' OR name ILIKE $2 OR engine ILIKE $2 OR database_name ILIKE $2 OR internal_dns ILIKE $2 OR external_dns ILIKE $2)
		ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &list, query, search, pattern); err != nil {
		return nil, err
	}
	for _, d := range list {
		r.decryptPassword(d)
	}
	return list, nil
}

func (r *sqlRepository) ListByProject(ctx context.Context, projectID string) ([]*Database, error) {
	var list []*Database
	query := `SELECT ` + selectDatabaseCols + ` FROM databases WHERE project_id = $1 ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &list, query, projectID); err != nil {
		return nil, err
	}
	for _, d := range list {
		r.decryptPassword(d)
	}
	return list, nil
}

func (r *sqlRepository) Update(ctx context.Context, d *Database) error {
	if d.Password != "" && r.vault != nil {
		enc, err := r.vault.Encrypt(d.Password)
		if err != nil {
			return fmt.Errorf("failed to encrypt database password: %w", err)
		}
		d.EncryptedPassword = enc
	}

	const query = `
		UPDATE databases SET
			name = :name, version = :version, username = :username, encrypted_password = :encrypted_password,
			database_name = :database_name, internal_host = :internal_host, port = :port,
			external_port = :external_port, internal_dns = :internal_dns, external_dns = :external_dns,
			volume_path = :volume_path, container_id = :container_id, status = :status,
			custom_args = :custom_args, logical_replication = :logical_replication,
			cpu_limit = :cpu_limit, memory_limit = :memory_limit, updated_at = :updated_at
		WHERE id = :id`
	_, err := r.db.NamedExecContext(ctx, query, d)
	return err
}

func (r *sqlRepository) UpdateStatus(ctx context.Context, id string, status DatabaseStatus, containerID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE databases SET status = $1, container_id = $2, updated_at = $3 WHERE id = $4`,
		status, containerID, time.Now(), id,
	)
	return err
}

func (r *sqlRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM databases WHERE id = $1`, id)
	return err
}

func (r *sqlRepository) GetNextAvailablePort(ctx context.Context) (int, error) {
	var maxPort sql.NullInt64
	if err := r.db.GetContext(ctx, &maxPort, `SELECT MAX(port) FROM databases`); err != nil {
		return 0, err
	}
	if !maxPort.Valid || maxPort.Int64 < 5432 {
		return 5432, nil
	}
	return int(maxPort.Int64) + 1, nil
}

func (r *sqlRepository) ProjectExists(ctx context.Context, projectID string) (bool, error) {
	if r.db == nil || projectID == "" {
		return false, nil
	}
	var count int
	if err := r.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM projects WHERE id = $1`, projectID); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *sqlRepository) IsProjectMemberOrOwner(ctx context.Context, projectID, userID, userRole string) bool {
	if userRole == "owner" || userRole == "admin" {
		return true
	}
	if r.db == nil || projectID == "" || userID == "" {
		return false
	}
	var count int
	query := `
		SELECT COUNT(*) FROM projects p
		INNER JOIN organization_members om ON p.organization_id = om.organization_id
		WHERE p.id = $1 AND om.user_id = $2 AND om.status = 'accepted'`
	if err := r.db.GetContext(ctx, &count, query, projectID, userID); err != nil {
		return false
	}
	return count > 0
}

func (r *sqlRepository) GetCustomDNSResolvers(ctx context.Context) string {
	if r.db == nil {
		return ""
	}
	var dns string
	_ = r.db.GetContext(ctx, &dns, `SELECT COALESCE(custom_dns_resolvers, '') FROM server_settings WHERE id = 'default'`)
	return dns
}

func (r *sqlRepository) LogAction(ctx context.Context, entry AuditLogEntry) {
	if r.db == nil {
		return
	}
	detailsJSON, _ := json.Marshal(entry.Details)
	_, _ = r.db.ExecContext(ctx, `
		INSERT INTO audit_logs (id, user_id, action, resource, ip_address, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		uuid.NewString(), entry.UserID, entry.Action, entry.Resource, entry.IPAddress, string(detailsJSON), time.Now(),
	)
}
