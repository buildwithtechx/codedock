package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type OperationRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewOperationRepo(db *sql.DB, vault Vault) *OperationRepo {
	return &OperationRepo{db: sqlx.NewDb(db, "pgx"), vault: vault}
}
func (r *OperationRepo) Create(ctx context.Context, op *models.Operation) error {
	encrypted, err := r.vault.Encrypt(op.Payload)
	if err != nil {
		return fmt.Errorf("encrypt operation: %w", err)
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO operations(id,user_id,project_id,kind,target,status,effects,encrypted_payload,snapshot,token_hash,expires_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`, op.ID, op.UserID, op.ProjectID, op.Kind, op.Target, op.Status, op.Effects, encrypted, op.Snapshot, op.TokenHash, op.ExpiresAt, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
func (r *OperationRepo) Get(ctx context.Context, id string) (*models.Operation, error) {
	var op models.Operation
	if err := r.db.GetContext(ctx, &op, `SELECT * FROM operations WHERE id=$1`, id); err != nil {
		return nil, fmt.Errorf("load operation: %w", err)
	}
	payload, err := r.vault.Decrypt(op.Payload)
	if err != nil {
		return nil, fmt.Errorf("decrypt operation: %w", err)
	}
	op.Payload = payload
	return &op, nil
}
func (r *OperationRepo) List(ctx context.Context, project, user string) ([]models.Operation, error) {
	result := []models.Operation{}
	err := r.db.SelectContext(ctx, &result, `SELECT id,user_id,project_id,kind,target,status,phase,effects,error,logs,expires_at,updated_at FROM operations WHERE (project_id=$1 OR $2='') AND user_id=$3 ORDER BY updated_at DESC LIMIT 100`, project, project, user)
	return result, err
}
func (r *OperationRepo) Claim(ctx context.Context, id, snapshot string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE operations SET status='RUNNING',phase='PREPARING',updated_at=$1 WHERE id=$2 AND status='REVIEWED' AND snapshot=$3 AND expires_at>$4`, time.Now().UTC().Format(time.RFC3339Nano), id, snapshot, time.Now().Unix())
	if err != nil {
		return fmt.Errorf("claim operation; another operation may own this target: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("review expired or target changed; prepare again")
	}
	return nil
}
func (r *OperationRepo) Observe(ctx context.Context, id, status, phase, message, logs string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE operations SET status=$1,phase=$2,error=$3,logs=right(logs||$4, 65536),updated_at=$5 WHERE id=$6 AND status IN ('RUNNING','CANCELLING')`, status, phase, message, logs, time.Now().UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("operation is no longer active")
	}
	return nil
}
func (r *OperationRepo) Recover(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE operations SET status='INTERRUPTED',error='Daemon restarted during execution. Inspect the target and prepare again before retrying.',updated_at=$1 WHERE status IN ('RUNNING','CANCELLING')`, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
