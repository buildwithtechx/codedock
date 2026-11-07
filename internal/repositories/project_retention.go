package repositories

import (
	"context"
	"time"
)

func (r *ProjectRepo) UpdateRollbackSettings(ctx context.Context, projectID string, strategy *string, window *int) error {
	now := time.Now().UTC()
	if strategy != nil && window != nil {
		_, err := r.db.ExecContext(ctx, `UPDATE projects SET default_rollback_strategy = $1, rollback_window = $2, updated_at = $3 WHERE id = $4`, *strategy, *window, now, projectID)
		return err
	} else if strategy != nil && window == nil {
		_, err := r.db.ExecContext(ctx, `UPDATE projects SET default_rollback_strategy = $1, rollback_window = NULL, updated_at = $2 WHERE id = $3`, *strategy, now, projectID)
		return err
	} else if strategy == nil && window != nil {
		_, err := r.db.ExecContext(ctx, `UPDATE projects SET rollback_window = $1, updated_at = $2 WHERE id = $3`, *window, now, projectID)
		return err
	}
	_, err := r.db.ExecContext(ctx, `UPDATE projects SET rollback_window = NULL, updated_at = $1 WHERE id = $2`, now, projectID)
	return err
}

func (r *ProjectRepo) UpdateCapacityMeasurement(ctx context.Context, projectID string, snapshotSizeBytes int64) error {
	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, `UPDATE projects SET snapshot_size_bytes = $1, capacity_measured_at = $2, updated_at = $3 WHERE id = $4`, snapshotSizeBytes, now, now, projectID)
	return err
}
