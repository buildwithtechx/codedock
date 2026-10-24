package repositories

import (
	"codedock/internal/models"
	"context"
	"fmt"
)

func (r *ServiceVarRepo) createReviewedVariable(ctx context.Context, variable *models.Variable) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	revision, err := topologyRevision(ctx, tx, variable.EnvironmentID)
	if err != nil {
		return err
	}
	if variable.ExpectedTopologyRevision != revision {
		return fmt.Errorf("topology changed since binding review; reload and review again")
	}
	_, err = tx.ExecContext(ctx, `INSERT INTO service_vars(id,service_id,environment_id,key,value,is_secret,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8)
 ON CONFLICT(service_id,key) DO UPDATE SET value=excluded.value,is_secret=excluded.is_secret,updated_at=excluded.updated_at`, variable.ID, variable.ServiceID, variable.EnvironmentID, variable.Key, variable.Value, variable.IsSecret, variable.CreatedAt, variable.UpdatedAt)
	if err != nil {
		return err
	}
	return tx.Commit()
}
