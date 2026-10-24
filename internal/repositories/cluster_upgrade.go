package repositories

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (r *ClusterRepo) BeginUpgrade(ctx context.Context, plan *models.ClusterPlan) error {
	data, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	encrypted, err := r.vault.Encrypt(string(data))
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `INSERT INTO cluster_upgrade_journals(cluster_id,encrypted_plan) VALUES($1,$2)`, plan.Cluster.ID, encrypted)
	return err
}
func (r *ClusterRepo) PendingUpgrades(ctx context.Context) ([]models.ClusterPlan, error) {
	encrypted := []string{}
	if err := r.db.SelectContext(ctx, &encrypted, `SELECT encrypted_plan FROM cluster_upgrade_journals`); err != nil {
		return nil, err
	}
	plans := []models.ClusterPlan{}
	for _, value := range encrypted {
		data, err := r.vault.Decrypt(value)
		if err != nil {
			return nil, err
		}
		var plan models.ClusterPlan
		if err := json.Unmarshal([]byte(data), &plan); err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}
func (r *ClusterRepo) FinishUpgrade(ctx context.Context, id, version, status, message string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `UPDATE clusters SET version=$1,status=$2,error=$3,revision=revision+1 WHERE id=$4 AND EXISTS(SELECT 1 FROM cluster_upgrade_journals WHERE cluster_id=$5)`, version, status, message, id, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("cluster upgrade journal missing")
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM cluster_upgrade_journals WHERE cluster_id=$1`, id); err != nil {
		return err
	}
	return tx.Commit()
}
