package deployments

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
)

type deploymentRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewDeploymentRepository(db *sql.DB) DeploymentRepository {
	return &deploymentRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *deploymentRepo) Create(ctx context.Context, d *Deployment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if d.ID == "" {
		d.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	if d.Status == "" {
		d.Status = DeploymentStatusBuilding
	}
	if d.OrganizationID == "" && d.ProjectID != "" {
		_ = r.db.GetContext(ctx, &d.OrganizationID, `SELECT organization_id FROM projects WHERE id = $1`, d.ProjectID)
	}
	_, err := r.db.ExecContext(ctx, `INSERT INTO deployments (
		id, service_id, organization_id, project_id, status, commit_hash,
		commit_message, branch, image_ref, "trigger", build_logs, container_id, resource_issue, created_at, updated_at, finished_at
	) VALUES ($1, NULLIF($2, ''), $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		d.ID, d.ServiceID, d.OrganizationID, d.ProjectID, d.Status, d.CommitHash,
		d.CommitMessage, d.Branch, d.ImageRef, d.Trigger, d.BuildLogs, d.ContainerID, d.ResourceIssue, d.CreatedAt, d.UpdatedAt, d.FinishedAt)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}
	return nil
}

func (r *deploymentRepo) GetByID(ctx context.Context, id string) (*Deployment, error) {
	var d Deployment
	err := r.db.GetContext(ctx, &d, `SELECT d.id, COALESCE(d.service_id, '') AS service_id, d.organization_id,
		COALESCE(s.environment_id, '') AS environment_id, d.project_id, d.status, d.commit_hash,
		d.commit_message, d.branch, d."trigger", COALESCE(d.image_ref, '') AS image_ref, d.build_logs, d.container_id,
		COALESCE(d.resource_issue, '') AS resource_issue, COALESCE(d.restart_loop, '') AS restart_loop,
		d.artifact_retained_at, COALESCE(d.pinned, FALSE) AS pinned, d.created_at, d.updated_at, d.finished_at
		FROM deployments d LEFT JOIN app_services s ON s.id = d.service_id WHERE d.id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrDeploymentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan deployment: %w", err)
	}
	return &d, nil
}

func (r *deploymentRepo) ListByService(ctx context.Context, serviceID string, limit, offset int) ([]*Deployment, int, error) {
	var total int
	if err := r.db.GetContext(ctx, &total, `SELECT COUNT(*) FROM deployments WHERE service_id = $1`, serviceID); err != nil {
		return nil, 0, err
	}

	var deps []*Deployment
	err := r.db.SelectContext(ctx, &deps, `SELECT d.id, COALESCE(d.service_id, '') AS service_id, d.organization_id,
		COALESCE(s.environment_id, '') AS environment_id, d.project_id, d.status, d.commit_hash,
		d.commit_message, d.branch, d."trigger", COALESCE(d.image_ref, '') AS image_ref, d.build_logs, d.container_id,
		COALESCE(d.resource_issue, '') AS resource_issue, COALESCE(d.restart_loop, '') AS restart_loop,
		d.artifact_retained_at, COALESCE(d.pinned, FALSE) AS pinned, d.created_at, d.updated_at, d.finished_at
		FROM deployments d LEFT JOIN app_services s ON s.id = d.service_id
		WHERE d.service_id = $1 ORDER BY d.created_at DESC LIMIT $2 OFFSET $3`, serviceID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if deps == nil {
		deps = make([]*Deployment, 0)
	}
	return deps, total, nil
}

func (r *deploymentRepo) ListByOrganization(ctx context.Context, filter DeploymentListFilter) ([]DeploymentListItem, int, error) {
	conditions := []string{"d.organization_id = $1"}
	args := []any{filter.OrganizationID}
	argIdx := 2

	if filter.ProjectID != "" {
		conditions = append(conditions, fmt.Sprintf("d.project_id = $%d", argIdx))
		args = append(args, filter.ProjectID)
		argIdx++
	}
	if filter.ServiceID != "" {
		conditions = append(conditions, fmt.Sprintf("d.service_id = $%d", argIdx))
		args = append(args, filter.ServiceID)
		argIdx++
	}
	if filter.Status != "" {
		conditions = append(conditions, fmt.Sprintf("UPPER(d.status) = UPPER($%d)", argIdx))
		args = append(args, filter.Status)
		argIdx++
	}
	if filter.Days > 0 {
		conditions = append(conditions, fmt.Sprintf("d.created_at >= NOW() - INTERVAL '%d days'", filter.Days))
	}
	if filter.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(d.commit_message ILIKE $%d OR d.branch ILIKE $%d OR s.name ILIKE $%d)", argIdx, argIdx, argIdx))
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	whereClause := strings.Join(conditions, " AND ")
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM deployments d LEFT JOIN app_services s ON s.id = d.service_id WHERE %s`, whereClause)
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	query := fmt.Sprintf(`SELECT d.id, COALESCE(d.service_id, '') AS service_id, COALESCE(s.name, '') AS service_name,
		COALESCE(s.environment_id, '') AS environment_id, d.project_id, COALESCE(p.name, '') AS project_name,
		d.status, d.branch, d.commit_hash, d.commit_message, d."trigger", d.container_id,
		d.artifact_retained_at, COALESCE(d.pinned, FALSE) AS pinned, COALESCE(d.resource_issue, '') AS resource_issue,
		COALESCE(d.restart_loop, '') AS restart_loop, d.created_at, d.updated_at, d.finished_at
		FROM deployments d
		LEFT JOIN app_services s ON s.id = d.service_id
		LEFT JOIN projects p ON p.id = d.project_id
		WHERE %s ORDER BY d.created_at DESC LIMIT $%d OFFSET $%d`, whereClause, argIdx, argIdx+1)

	args = append(args, filter.Limit, filter.Offset)
	var items []DeploymentListItem
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return nil, 0, err
	}
	if items == nil {
		items = make([]DeploymentListItem, 0)
	}
	return items, total, nil
}

func (r *deploymentRepo) Update(ctx context.Context, d *Deployment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	d.UpdatedAt = time.Now().UTC()
	query := `UPDATE deployments SET status = $1, commit_hash = $2, commit_message = $3, branch = $4,
		image_ref = $5, "trigger" = $6, build_logs = $7, container_id = $8, resource_issue = $9,
		updated_at = $10, finished_at = $11 WHERE id = $12`
	_, err := r.db.ExecContext(ctx, query, d.Status, d.CommitHash, d.CommitMessage, d.Branch,
		d.ImageRef, d.Trigger, d.BuildLogs, d.ContainerID, d.ResourceIssue, d.UpdatedAt, d.FinishedAt, d.ID)
	return err
}

func (r *deploymentRepo) UpdateStatus(ctx context.Context, id string, status DeploymentStatus, buildLogs, containerID, resourceIssue string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now().UTC()
	var finishedAt *time.Time
	if status == DeploymentStatusReady || status == DeploymentStatusActive || status == DeploymentStatusFailed || status == DeploymentStatusCancelled {
		finishedAt = &now
	}
	query := `UPDATE deployments SET status = $1, build_logs = $2, container_id = $3, resource_issue = $4,
		updated_at = $5, finished_at = COALESCE($6, finished_at) WHERE id = $7`
	_, err := r.db.ExecContext(ctx, query, status, buildLogs, containerID, resourceIssue, now, finishedAt, id)
	return err
}

func (r *deploymentRepo) SetPinned(ctx context.Context, id string, pinned bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE deployments SET pinned = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, pinned, id)
	return err
}

func (r *deploymentRepo) TogglePinned(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var newPinned bool
	err := r.db.GetContext(ctx, &newPinned, `UPDATE deployments SET pinned = NOT COALESCE(pinned, FALSE), updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING pinned`, id)
	return newPinned, err
}
