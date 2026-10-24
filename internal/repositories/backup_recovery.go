package repositories

import (
	"context"
	"fmt"
	"time"
)

func (r *BackupRepo) ProtectRecord(ctx context.Context, id string, until int64) error {
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

func (r *BackupRepo) ClearRestoreProtection(ctx context.Context, id string, before int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE backup_records SET protected_until=0 WHERE id=$1 AND protected_until<=$2`, id, before)
	return err
}

func (r *BackupRepo) ClaimRecordExpiry(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.ExecContext(ctx, `UPDATE backup_records SET status='expiring' WHERE id=$1 AND status='completed' AND protected_until<=$2`, id, time.Now().Unix())
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count == 1, err
}
func (r *BackupRepo) RecoverRecords(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE backup_records SET status='failed',logs=COALESCE(logs,'')||'\nDaemon restarted before backup completion. Retry the policy.',completed_at=$1 WHERE status IN ('running','expiring')`, time.Now().UTC().Format(time.RFC3339))
	return err
}

func (r *BackupRepo) ClaimRecordDeletion(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.ExecContext(ctx, `UPDATE backup_records SET status='expiring' WHERE id=$1 AND status IN ('completed','failed','expired') AND protected_until<=$2`, id, time.Now().Unix())
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count == 1, err
}
