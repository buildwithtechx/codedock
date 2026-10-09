package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
	"codedock/internal/utils"
)

type DeploymentRepository interface {
	Create(ctx context.Context, d *models.Deployment) error
	GetByID(ctx context.Context, id string) (*models.Deployment, error)
	ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*models.Deployment, int, error)
	ListByOrganization(ctx context.Context, filter models.DeploymentListFilter) ([]models.DeploymentListItem, int, error)
	Update(ctx context.Context, d *models.Deployment) error
	UpdateStatus(ctx context.Context, id string, status models.DeploymentStatus, buildLogs, containerID string) error
}

type DeploymentRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewDeploymentRepo(db *sql.DB) *DeploymentRepo {
	return &DeploymentRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *DeploymentRepo) Create(ctx context.Context, d *models.Deployment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.Status == "" {
		d.Status = "BUILDING"
	}
	if err := r.db.GetContext(ctx, &d.OrganizationID, `SELECT organization_id FROM projects WHERE id = $1`, d.ProjectID); err != nil {
		return fmt.Errorf("resolve deployment organization: %w", err)
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO deployments (
		id, service_id, organization_id, project_id, status, commit_hash,
		commit_message, branch, image_ref, "trigger", build_logs, container_id, created_at, updated_at, finished_at
	) VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`,
		d.ID, d.ServiceID, d.OrganizationID, d.ProjectID, d.Status, d.CommitHash,
		d.CommitMessage, d.Branch, d.ImageRef, d.Trigger, d.BuildLogs, d.ContainerID, d.CreatedAt, d.UpdatedAt, d.FinishedAt)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}
	return nil
}

func (r *DeploymentRepo) GetByID(ctx context.Context, id string) (*models.Deployment, error) {
	var d models.Deployment
	err := r.db.GetContext(ctx, &d, `SELECT d.id, COALESCE(d.service_id, '') AS service_id, d.organization_id,
		COALESCE(s.environment_id, '') AS environment_id, d.project_id, d.status, d.commit_hash,
		d.commit_message, d.branch, d."trigger", COALESCE(d.image_ref, '') AS image_ref, d.build_logs, d.container_id, d.created_at, d.updated_at, d.finished_at
		FROM deployments d LEFT JOIN app_services s ON s.id = d.service_id WHERE d.id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Deployment", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan deployment: %w", err)
	}
	return &d, nil
}

func (r *DeploymentRepo) ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*models.Deployment, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM deployments WHERE service_id = $1`, serviceID); err != nil {
		return nil, 0, err
	}

	var deps []*models.Deployment
	err := r.db.SelectContext(ctx, &deps, `SELECT d.id, COALESCE(d.service_id, '') AS service_id, d.organization_id,
		COALESCE(s.environment_id, '') AS environment_id, d.project_id, d.status, d.commit_hash,
		d.commit_message, d.branch, d."trigger", COALESCE(d.image_ref, '') AS image_ref, d.build_logs, d.container_id, d.created_at, d.updated_at, d.finished_at
		FROM deployments d LEFT JOIN app_services s ON s.id = d.service_id
		WHERE d.service_id = $1 ORDER BY d.created_at DESC LIMIT $2 OFFSET $3`, serviceID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query service deployments: %w", err)
	}
	if deps == nil {
		deps = make([]*models.Deployment, 0)
	}
	return deps, total, nil
}

func (r *DeploymentRepo) ListByOrganization(ctx context.Context, filter models.DeploymentListFilter) ([]models.DeploymentListItem, int, error) {
	search := strings.TrimSpace(filter.Search)
	searchPattern := "%" + search + "%"
	where := `
		FROM deployments d
		INNER JOIN projects p ON p.id = d.project_id
		LEFT JOIN app_services s ON s.id = d.service_id
		WHERE p.organization_id = $1
			AND ($2 = '' OR d.project_id = $3)
			AND ($4 = '' OR d.service_id = $5)
			AND ($6 = '' OR lower(d.status) = lower($7))
			AND (
				$8 = '' OR s.name ILIKE $9 OR p.name ILIKE $10
				OR d.branch ILIKE $11 OR d.commit_hash ILIKE $12
			)`
	args := []any{
		filter.OrganizationID,
		filter.ProjectID,
		filter.ProjectID,
		filter.ServiceID,
		filter.ServiceID,
		filter.Status,
		filter.Status,
		search,
		searchPattern,
		searchPattern,
		searchPattern,
		searchPattern,
	}

	var total int
	if err := r.db.GetContext(ctx, &total, "SELECT COUNT(*) "+where, args...); err != nil {
		return nil, 0, fmt.Errorf("count organization deployments: %w", err)
	}

	deployments := make([]models.DeploymentListItem, 0)
	query := `SELECT d.id, COALESCE(d.service_id, '') AS service_id, COALESCE(s.name, '') AS service_name,
		COALESCE(s.environment_id, '') AS environment_id,
		d.project_id, p.name AS project_name, d.status, d.commit_hash, d.commit_message,
		d.branch, d."trigger", d.container_id, d.created_at, d.updated_at, d.finished_at ` + where + `
		ORDER BY d.created_at DESC LIMIT $13 OFFSET $14`
	args = append(args, filter.Limit, filter.Offset)
	if err := r.db.SelectContext(ctx, &deployments, query, args...); err != nil {
		return nil, 0, fmt.Errorf("list organization deployments: %w", err)
	}
	return deployments, total, nil
}

func (r *DeploymentRepo) Update(_ context.Context, d *models.Deployment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d.UpdatedAt = time.Now().UTC()
	_, err := r.db.Exec(`UPDATE deployments SET status = $1, commit_hash = $2, commit_message = $3,
		branch = $4, image_ref = $5, "trigger" = $6, build_logs = $7, container_id = $8, updated_at = $9, finished_at = $10 WHERE id = $11`,
		d.Status, d.CommitHash, d.CommitMessage, d.Branch, d.ImageRef, d.Trigger, d.BuildLogs, d.ContainerID, d.UpdatedAt, d.FinishedAt, d.ID)
	return err
}

func (r *DeploymentRepo) UpdateStatus(_ context.Context, id string, status models.DeploymentStatus, buildLogs, containerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	if status == models.DeploymentStatusReady || status == models.DeploymentStatus("CANCELLED") || status == models.DeploymentStatusActive || status == models.DeploymentStatusFailed || status == models.DeploymentStatusRemoved || status == models.DeploymentStatusSlept {
		_, err := r.db.Exec(`UPDATE deployments SET status = $1, build_logs = $2, container_id = $3, updated_at = $4, finished_at = $5 WHERE id = $6`,
			status, buildLogs, containerID, now, now, id)
		return err
	}
	_, err := r.db.Exec(`UPDATE deployments SET status = $1, build_logs = $2, container_id = $3, updated_at = $4 WHERE id = $5`,
		status, buildLogs, containerID, now, id)
	return err
}
