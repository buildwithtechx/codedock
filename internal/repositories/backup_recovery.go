package repositories

import (
	"context"
	"fmt"
	"time"
)

func (r *BackupRepo) ensureRecoveryColumns() error {
	for _, column := range []struct{ table, name, definition string }{
		{"backup_records", "protected_until", "INTEGER NOT NULL DEFAULT 0"},
		{"backup_records", "sha256", "TEXT NOT NULL DEFAULT ''"},
		{"backup_records", "verified_at", "TEXT NOT NULL DEFAULT ''"},
		{"backup_configs", "pre_deployment", "INTEGER NOT NULL DEFAULT 0"},
		{"backup_configs", "owner_id", "TEXT NOT NULL DEFAULT ''"},
		{"backup_configs", "project_id", "TEXT NOT NULL DEFAULT ''"},
	} {
		rows, err := r.db.Query("PRAGMA table_info(" + column.table + ")")
		if err != nil {
			return err
		}
		exists := false
		for rows.Next() {
			var id, notNull, key int
			var name, kind string
			var defaultValue any
			if err := rows.Scan(&id, &name, &kind, &notNull, &defaultValue, &key); err != nil {
				rows.Close()
				return err
			}
			exists = exists || name == column.name
		}
		rowErr := rows.Err()
		closeErr := rows.Close()
		if rowErr != nil {
			return rowErr
		}
		if closeErr != nil {
			return closeErr
		}
		if !exists {
			if _, err := r.db.Exec("ALTER TABLE " + column.table + " ADD COLUMN " + column.name + " " + column.definition); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *BackupRepo) ProtectRecord(ctx context.Context, id string, until int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if until < 0 || until > time.Now().AddDate(10, 0, 0).Unix() {
		return fmt.Errorf("protection expiry must be within ten years")
	}
	result, err := r.db.ExecContext(ctx, `UPDATE backup_records SET protected_until=? WHERE id=? AND status='completed'`, until, id)
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
	_, err := r.db.ExecContext(ctx, `UPDATE backup_records SET protected_until=0 WHERE id=? AND protected_until<=?`, id, before)
	return err
}

func (r *BackupRepo) ClaimRecordExpiry(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.ExecContext(ctx, `UPDATE backup_records SET status='expiring' WHERE id=? AND status='completed' AND protected_until<=?`, id, time.Now().Unix())
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count == 1, err
}
func (r *BackupRepo) RecoverRecords(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE backup_records SET status='failed',logs=COALESCE(logs,'')||'\nDaemon restarted before backup completion. Retry the policy.',completed_at=? WHERE status IN ('running','expiring')`, time.Now().UTC().Format(time.RFC3339))
	return err
}

func (r *BackupRepo) ClaimRecordDeletion(ctx context.Context, id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	result, err := r.db.ExecContext(ctx, `UPDATE backup_records SET status='expiring' WHERE id=? AND status IN ('completed','failed','expired') AND protected_until<=?`, id, time.Now().Unix())
	if err != nil {
		return false, err
	}
	count, err := result.RowsAffected()
	return count == 1, err
}
