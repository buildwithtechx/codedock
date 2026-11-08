package backups

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *repo) ProtectRecord(ctx context.Context, id string, until int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if until < 0 || until > time.Now().AddDate(10, 0, 0).Unix() {
		return fmt.Errorf("protection expiry must be within ten years")
	}
	result, err := r.db.ExecContext(ctx, `UPDATE backup_records SET protected_until=$1 WHERE id=$2 AND status='completed'`, until, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("only completed backup records can be protected")
	}
	return nil
}

func (r *repo) ClearRestoreProtection(ctx context.Context, id string, before int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE backup_records SET protected_until=0 WHERE id=$1 AND protected_until<=$2`, id, before)
	return err
}

func (r *repo) ClaimRecordExpiry(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.ExecContext(ctx, `UPDATE backup_records SET status='expiring' WHERE id=$1 AND status='completed' AND protected_until<=$2`, id, time.Now().Unix())
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count == 1, err
}

func (r *repo) SetSnapshotDumper(dumper func(ctx context.Context) ([]byte, error)) {
	r.snapshotDumper = dumper
}

func (r *repo) ControlPlaneSnapshot(ctx context.Context) ([]byte, error) {
	if r.snapshotDumper == nil {
		return nil, fmt.Errorf("control-plane snapshotter is not configured")
	}
	return r.snapshotDumper(ctx)
}

func (r *repo) ResolveProjectID(ctx context.Context, databaseID, serviceID string) string {
	if r.db == nil {
		return ""
	}
	var projectID string
	if databaseID != "" && databaseID != "global" {
		_ = r.db.GetContext(ctx, &projectID, `SELECT COALESCE(project_id, '') FROM databases WHERE id = $1`, databaseID)
		return projectID
	}
	if serviceID != "" {
		_ = r.db.GetContext(ctx, &projectID, `SELECT COALESCE(project_id, '') FROM app_services WHERE id = $1`, serviceID)
		return projectID
	}
	return ""
}

func (r *repo) HasProjectPermission(ctx context.Context, projectID, userID, role, permission string) bool {
	if role == "admin" || role == "owner" {
		return true
	}
	if r.db == nil || projectID == "" {
		return false
	}
	var perm string
	err := r.db.GetContext(ctx, &perm, `SELECT permission FROM project_members WHERE project_id = $1 AND user_id = $2`, projectID, userID)
	if err != nil {
		err = r.db.GetContext(ctx, &perm, `SELECT om.permission FROM organization_members om JOIN projects p ON p.organization_id = om.organization_id WHERE p.id = $1 AND om.user_id = $2`, projectID, userID)
		if err != nil {
			return false
		}
	}
	if permission == "admin" {
		return perm == "admin" || perm == "owner"
	}
	return perm != ""
}

func (r *repo) LogAudit(ctx context.Context, userID, action, resource, ip string) {
	if r.db == nil {
		return
	}
	_, _ = r.db.ExecContext(ctx, `INSERT INTO audit_logs (id, user_id, action, resource, ip_address, created_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		uuid.NewString(), userID, action, resource, ip, time.Now().UTC().Format(time.RFC3339))
}
