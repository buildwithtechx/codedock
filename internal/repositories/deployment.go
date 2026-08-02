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

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
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
	return &DeploymentRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *DeploymentRepo) Create(_ context.Context, d *models.Deployment) error {
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
	_, err := r.db.Exec(`INSERT INTO deployments (
		id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.ID, d.ServiceID, d.EnvironmentID, d.ProjectID, d.Status, d.CommitHash,
		d.CommitMessage, d.Branch, d.Trigger, d.BuildLogs, d.ContainerID, d.CreatedAt, d.UpdatedAt, d.FinishedAt)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}
	return nil
}

func (r *DeploymentRepo) GetByID(ctx context.Context, id string) (*models.Deployment, error) {
	var d models.Deployment
	err := r.db.GetContext(ctx, &d, `SELECT id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
		FROM deployments WHERE id = ?`, id)
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
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM deployments WHERE service_id = ?`, serviceID); err != nil {
		return nil, 0, err
	}

	var deps []*models.Deployment
	err := r.db.SelectContext(ctx, &deps, `SELECT id, service_id, environment_id, project_id, status, commit_hash,
		commit_message, branch, trigger, build_logs, container_id, created_at, updated_at, finished_at
		FROM deployments WHERE service_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`, serviceID, limit, offset)
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
		WHERE p.organization_id = ?
			AND (? = '' OR d.project_id = ?)
			AND (? = '' OR d.service_id = ?)
			AND (? = '' OR lower(d.status) = lower(?))
			AND (
				? = '' OR s.name LIKE ? COLLATE NOCASE OR p.name LIKE ? COLLATE NOCASE
				OR d.branch LIKE ? COLLATE NOCASE OR d.commit_hash LIKE ? COLLATE NOCASE
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
	query := `SELECT d.id, d.service_id, COALESCE(s.name, '' ) AS service_name, d.environment_id,
		d.project_id, p.name AS project_name, d.status, d.commit_hash, d.commit_message,
		d.branch, d.trigger, d.build_logs, d.container_id, d.created_at, d.updated_at, d.finished_at ` + where + `
		ORDER BY d.created_at DESC LIMIT ? OFFSET ?`
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
	_, err := r.db.Exec(`UPDATE deployments SET status = ?, commit_hash = ?, commit_message = ?,
		branch = ?, trigger = ?, build_logs = ?, container_id = ?, updated_at = ?, finished_at = ? WHERE id = ?`,
		d.Status, d.CommitHash, d.CommitMessage, d.Branch, d.Trigger, d.BuildLogs, d.ContainerID, d.UpdatedAt, d.FinishedAt, d.ID)
	return err
}

func (r *DeploymentRepo) UpdateStatus(_ context.Context, id string, status models.DeploymentStatus, buildLogs, containerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	if status == models.DeploymentStatusActive || status == models.DeploymentStatusFailed || status == models.DeploymentStatusRemoved || status == models.DeploymentStatusSlept {
		_, err := r.db.Exec(`UPDATE deployments SET status = ?, build_logs = ?, container_id = ?, updated_at = ?, finished_at = ? WHERE id = ?`,
			status, buildLogs, containerID, now, now, id)
		return err
	}
	_, err := r.db.Exec(`UPDATE deployments SET status = ?, build_logs = ?, container_id = ?, updated_at = ? WHERE id = ?`,
		status, buildLogs, containerID, now, id)
	return err
}
