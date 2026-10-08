package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (r *RuntimeRepo) DesiredIDs(ctx context.Context) ([]string, error) {
	ids := []string{}
	if err := r.db.SelectContext(ctx, &ids, `SELECT d.service_id FROM runtime_desired d JOIN service_runtimes r ON r.service_id=d.service_id WHERE r.encrypted_journal='' ORDER BY r.updated_at`); err != nil {
		return nil, err
	}
	return ids, nil
}
func (r *RuntimeRepo) Desired(ctx context.Context, id string) (*models.DesiredRuntime, error) {
	var encrypted string
	if err := r.db.GetContext(ctx, &encrypted, `SELECT encrypted_workload FROM runtime_desired WHERE service_id=$1`, id); err != nil {
		return nil, err
	}
	data, err := r.vault.Decrypt(encrypted)
	if err != nil {
		return nil, err
	}
	var desired models.DesiredRuntime
	if err := json.Unmarshal([]byte(data), &desired); err != nil {
		return nil, err
	}
	return &desired, nil
}
func (r *RuntimeRepo) CommitDesired(ctx context.Context, id string, revision int, desired *models.DesiredRuntime) error {
	data, err := json.Marshal(desired)
	if err != nil {
		return err
	}
	encrypted, err := r.vault.Encrypt(string(data))
	if err != nil {
		return err
	}
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE service_runtimes SET encrypted_journal='',status='READY',error='' WHERE service_id=$1 AND revision=$2 AND encrypted_journal!=''`, id, revision)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("runtime deployment claim changed")
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO runtime_desired(service_id,encrypted_workload) VALUES($1,$2) ON CONFLICT(service_id) DO UPDATE SET encrypted_workload=excluded.encrypted_workload`, id, encrypted); err != nil {
		return err
	}
	return tx.Commit()
}
