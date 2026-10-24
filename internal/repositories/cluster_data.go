package repositories

import (
	"codedock/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type ClusterDataRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewClusterDataRepo(db *sql.DB, vault Vault) *ClusterDataRepo {
	return &ClusterDataRepo{sqlx.NewDb(db, "pgx"), vault}
}
func (r *ClusterDataRepo) Create(ctx context.Context, plan *models.ClusterDataPlan) error {
	data, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	encrypted, err := r.vault.Encrypt(string(data))
	if err != nil {
		return err
	}
	var organization string
	if err := r.db.GetContext(ctx, &organization, `SELECT organization_id FROM clusters WHERE id=$1 AND project_id=$2`, plan.Record.ClusterID, plan.Record.ProjectID); err != nil {
		return fmt.Errorf("validate database cluster ownership: %w", err)
	}
	secret, err := r.vault.Encrypt(plan.Password)
	if err != nil {
		return err
	}
	port, username, database := 6379, "default", "0"
	if plan.Record.Spec.Engine == "postgres" {
		port, username, database = 5432, "app", "app"
	}
	_, version, _ := strings.Cut(plan.Record.Spec.Image, ":")
	_, err = r.db.ExecContext(ctx, `INSERT INTO cluster_databases(id,organization_id,cluster_id,project_id,name,engine,version,port,username,database_name,secret_encrypted,encrypted_config,status,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,'REVIEWED',$13)`, plan.Record.ID, organization, plan.Record.ClusterID, plan.Record.ProjectID, plan.Record.Spec.Name, plan.Record.Spec.Engine, version, port, username, database, secret, encrypted, time.Now().UTC().Format(time.RFC3339Nano))
	return err
}
func (r *ClusterDataRepo) Get(ctx context.Context, id string) (*models.ClusterDataPlan, error) {
	var record models.ClusterData
	if err := r.db.GetContext(ctx, &record, `SELECT id,cluster_id,project_id,encrypted_config,status,error,updated_at FROM cluster_databases WHERE id=$1 AND encrypted_config!=''`, id); err != nil {
		return nil, err
	}
	data, err := r.vault.Decrypt(record.Config)
	if err != nil {
		return nil, err
	}
	var plan models.ClusterDataPlan
	if err := json.Unmarshal([]byte(data), &plan); err != nil {
		return nil, err
	}
	record.Spec = plan.Record.Spec
	plan.Record = record
	return &plan, nil
}
func (r *ClusterDataRepo) List(ctx context.Context, cluster string) ([]models.ClusterData, error) {
	ids := []string{}
	if err := r.db.SelectContext(ctx, &ids, `SELECT id FROM cluster_databases WHERE cluster_id=$1 ORDER BY updated_at DESC`, cluster); err != nil {
		return nil, err
	}
	result := []models.ClusterData{}
	for _, id := range ids {
		plan, err := r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, plan.Record)
	}
	return result, nil
}

func (r *ClusterDataRepo) ListByProject(ctx context.Context, project string) ([]models.ClusterData, error) {
	ids := []string{}
	if err := r.db.SelectContext(ctx, &ids, `SELECT id FROM cluster_databases WHERE project_id=$1 ORDER BY updated_at DESC`, project); err != nil {
		return nil, err
	}
	result := []models.ClusterData{}
	for _, id := range ids {
		plan, err := r.Get(ctx, id)
		if err != nil {
			return nil, err
		}
		result = append(result, plan.Record)
	}
	return result, nil
}
func (r *ClusterDataRepo) Observe(ctx context.Context, id, status, message string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE cluster_databases SET status=$1,error=$2,updated_at=$3 WHERE id=$4`, status, message, time.Now().UTC().Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("cluster database disappeared")
	}
	return nil
}

func (r *ClusterDataRepo) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM cluster_databases WHERE id=$1`, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return fmt.Errorf("cluster database disappeared")
	}
	return nil
}
func (r *ClusterDataRepo) Recover(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE cluster_databases SET status='INTERRUPTED',error='Daemon restarted; review recovery to reconcile owned resources and retained storage' WHERE status='APPLYING'`)
	return err
}
