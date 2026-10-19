package repositories

import (
	"context"
	"fmt"
	"time"

	"codedock/internal/models"
)

func (r *DeploymentRepo) ListForRetention(ctx context.Context, projectID string) ([]*models.Deployment, error) {
	var deps []*models.Deployment
	query := `SELECT d.id, COALESCE(d.service_id, '') AS service_id, d.organization_id,
		COALESCE(s.environment_id, '') AS environment_id, d.project_id, d.status, d.commit_hash,
		d.commit_message, d.branch, d."trigger", COALESCE(d.image_ref, '') AS image_ref, d.build_logs, d.container_id,
		d.artifact_retained_at, COALESCE(d.pinned, FALSE) AS pinned, d.created_at, d.updated_at, d.finished_at
		FROM deployments d LEFT JOIN app_services s ON s.id = d.service_id
		WHERE d.project_id = $1 AND (
			UPPER(d.status) IN ('READY', 'ACTIVE', 'SUCCESS', 'PARTIAL_FAILURE', 'BUILDING', 'CLONING', 'PULLING', 'PENDING')
			OR d.artifact_retained_at IS NOT NULL
			OR d.pinned = TRUE
		)
		ORDER BY d.created_at DESC, d.id DESC`

	if err := r.db.SelectContext(ctx, &deps, query, projectID); err != nil {
		return nil, fmt.Errorf("list deployments for retention: %w", err)
	}
	if deps == nil {
		deps = make([]*models.Deployment, 0)
	}
	return deps, nil
}

func (r *DeploymentRepo) SetArtifactRetainedAt(ctx context.Context, id string, at *time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE deployments SET artifact_retained_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, at, id)
	return err
}

func (r *DeploymentRepo) SetPinned(ctx context.Context, id string, pinned bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE deployments SET pinned = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, pinned, id)
	return err
}

func (r *DeploymentRepo) TogglePinned(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var newPinned bool
	err := r.db.GetContext(ctx, &newPinned, `UPDATE deployments SET pinned = NOT COALESCE(pinned, FALSE), updated_at = CURRENT_TIMESTAMP WHERE id = $1 RETURNING pinned`, id)
	return newPinned, err
}
