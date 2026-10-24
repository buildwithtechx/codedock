package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock/internal/models"
)

type RouteProtectionRepository interface {
	ListByProject(ctx context.Context, projectID string) ([]models.RouteProtection, error)
	Get(ctx context.Context, projectID, hostname string) (*models.RouteProtection, error)
	GetByHostname(ctx context.Context, hostname string) (*models.RouteProtection, error)
	Save(ctx context.Context, rp *models.RouteProtection) error
	MarkApplied(ctx context.Context, projectID, hostname, revision string) error
	Delete(ctx context.Context, projectID, hostname string) error
}

type postgresRouteProtectionRepository struct {
	db *sql.DB
}

func NewRouteProtectionRepository(db *sql.DB) RouteProtectionRepository {
	return &postgresRouteProtectionRepository{db: db}
}

type rpScanner interface {
	Scan(dest ...any) error
}

func (r *postgresRouteProtectionRepository) scan(s rpScanner) (*models.RouteProtection, error) {
	var (
		rp              models.RouteProtection
		secret          sql.NullString
		appliedRevision sql.NullString
	)
	err := s.Scan(
		&rp.ID,
		&rp.ProjectID,
		&rp.Hostname,
		&rp.Enabled,
		&secret,
		&rp.Revision,
		&appliedRevision,
		&rp.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan route protection: %w", err)
	}
	if secret.Valid {
		rp.Secret = &secret.String
	}
	if appliedRevision.Valid {
		rp.AppliedRevision = &appliedRevision.String
	}
	return &rp, nil
}

func (r *postgresRouteProtectionRepository) ListByProject(ctx context.Context, projectID string) ([]models.RouteProtection, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, project_id, hostname, enabled, secret, revision, applied_revision, updated_at
		FROM route_protection WHERE project_id = $1 ORDER BY hostname ASC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list route protection: %w", err)
	}
	defer rows.Close()

	var list []models.RouteProtection
	for rows.Next() {
		rp, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		if rp != nil {
			list = append(list, *rp)
		}
	}
	return list, rows.Err()
}

func (r *postgresRouteProtectionRepository) Get(ctx context.Context, projectID, hostname string) (*models.RouteProtection, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, project_id, hostname, enabled, secret, revision, applied_revision, updated_at
		FROM route_protection WHERE project_id = $1 AND hostname = $2`, projectID, hostname)
	return r.scan(row)
}

func (r *postgresRouteProtectionRepository) GetByHostname(ctx context.Context, hostname string) (*models.RouteProtection, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, project_id, hostname, enabled, secret, revision, applied_revision, updated_at
		FROM route_protection WHERE hostname = $1`, hostname)
	return r.scan(row)
}

func (r *postgresRouteProtectionRepository) Save(ctx context.Context, rp *models.RouteProtection) error {
	if rp.ID == "" {
		rp.ID = uuid.NewString()
	}
	rp.UpdatedAt = time.Now().UTC()

	var secretVal *string
	if rp.Secret != nil && *rp.Secret != "" {
		secretVal = rp.Secret
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO route_protection (id, project_id, hostname, enabled, secret, revision, applied_revision, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (project_id, hostname) DO UPDATE SET
			enabled = EXCLUDED.enabled,
			secret = COALESCE(EXCLUDED.secret, route_protection.secret),
			revision = EXCLUDED.revision,
			applied_revision = EXCLUDED.applied_revision,
			updated_at = EXCLUDED.updated_at`,
		rp.ID, rp.ProjectID, rp.Hostname, rp.Enabled, secretVal, rp.Revision, rp.AppliedRevision, rp.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save route protection: %w", err)
	}
	return nil
}

func (r *postgresRouteProtectionRepository) MarkApplied(ctx context.Context, projectID, hostname, revision string) error {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `
		UPDATE route_protection
		SET applied_revision = $1, updated_at = $2
		WHERE project_id = $3 AND hostname = $4 AND revision = $1`,
		revision, now, projectID, hostname,
	)
	if err != nil {
		return fmt.Errorf("mark applied route protection: %w", err)
	}
	return nil
}

func (r *postgresRouteProtectionRepository) Delete(ctx context.Context, projectID, hostname string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM route_protection WHERE project_id = $1 AND hostname = $2`,
		projectID, hostname,
	)
	if err != nil {
		return fmt.Errorf("delete route protection: %w", err)
	}
	return nil
}
