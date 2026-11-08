package deployments

import (
	"context"
	"fmt"
	"time"
)

func (r *deploymentRepo) ListForRetention(ctx context.Context, projectID string) ([]*Deployment, error) {
	var deps []*Deployment
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
		deps = make([]*Deployment, 0)
	}
	return deps, nil
}

func (r *deploymentRepo) SetArtifactRetainedAt(ctx context.Context, id string, at *time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE deployments SET artifact_retained_at = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, at, id)
	return err
}

func (r *deploymentRepo) SetRestartLoop(ctx context.Context, id string, finding string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE deployments SET restart_loop = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`, finding, id)
	return err
}

func (r *deploymentRepo) RecoverInterrupted(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE deployments SET status='FAILED', finished_at=CURRENT_TIMESTAMP,
		build_logs=COALESCE(build_logs,'')||chr(10)||'Daemon restarted during deployment. Retained rollout containers have been reconciled; inspect runtime before retrying.'
		WHERE status IN ('pending','PREPARING','CLONING','PULLING','BUILDING','STARTING','READINESS','ROUTING')`)
	return err
}
