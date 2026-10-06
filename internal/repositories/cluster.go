package repositories

import (
	"codedock.run/codedock/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

type ClusterRepo struct {
	db    *sqlx.DB
	vault Vault
}

func NewClusterRepo(db *sql.DB, vault Vault) *ClusterRepo {
	return &ClusterRepo{db: sqlx.NewDb(db, "sqlite"), vault: vault}
}
func (r *ClusterRepo) Get(ctx context.Context, id string) (*models.Cluster, error) {
	var cluster models.Cluster
	if err := r.db.GetContext(ctx, &cluster, `SELECT * FROM clusters WHERE id=?`, id); err != nil {
		return nil, fmt.Errorf("load cluster: %w", err)
	}
	token, err := r.vault.Decrypt(cluster.JoinToken)
	if err != nil {
		return nil, fmt.Errorf("decrypt cluster token: %w", err)
	}
	cluster.JoinToken = token
	if err := json.Unmarshal([]byte(cluster.NodesJSON), &cluster.Nodes); err != nil {
		return nil, err
	}
	return &cluster, nil
}
func (r *ClusterRepo) List(ctx context.Context, project string) ([]models.Cluster, error) {
	result := []models.Cluster{}
	if err := r.db.SelectContext(ctx, &result, `SELECT id,project_id,organization_id,name,version,nodes_json,revision,status,error,updated_at FROM clusters WHERE project_id=? ORDER BY name`, project); err != nil {
		return nil, err
	}
	for i := range result {
		if err := json.Unmarshal([]byte(result[i].NodesJSON), &result[i].Nodes); err != nil {
			return nil, err
		}
	}
	return result, nil
}
func (r *ClusterRepo) Save(ctx context.Context, cluster *models.Cluster, previous int) error {
	nodes, err := json.Marshal(cluster.Nodes)
	if err != nil {
		return err
	}
	token, err := r.vault.Encrypt(cluster.JoinToken)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	if previous == 0 {
		_, err = r.db.ExecContext(ctx, `INSERT INTO clusters(id,project_id,organization_id,name,version,nodes_json,encrypted_token,updated_at) VALUES(?,?,?,?,?,?,?,?)`, cluster.ID, cluster.ProjectID, cluster.OrganizationID, cluster.Name, cluster.Version, string(nodes), token, now)
	} else {
		result, updateErr := r.db.ExecContext(ctx, `UPDATE clusters SET status='REVIEWED',error='',name=?,version=?,nodes_json=?,revision=revision+1,updated_at=? WHERE id=? AND project_id=? AND revision=? AND NOT EXISTS(SELECT 1 FROM cluster_upgrade_journals WHERE cluster_id=clusters.id) AND NOT EXISTS(SELECT 1 FROM operations WHERE target=? AND status IN ('RUNNING','CANCELLING'))`, cluster.Name, cluster.Version, string(nodes), now, cluster.ID, cluster.ProjectID, previous, "cluster:"+cluster.ID)
		if updateErr != nil {
			return updateErr
		}
		count, updateErr := result.RowsAffected()
		if updateErr != nil {
			return updateErr
		}
		if count != 1 {
			return fmt.Errorf("cluster changed or an operation is active")
		}
	}
	return err
}
func (r *ClusterRepo) Observe(ctx context.Context, id, status, message string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE clusters SET status=?,error=?,updated_at=? WHERE id=?`, status, message, time.Now().UTC().Format(time.RFC3339Nano), id)
	return err
}

func (r *ClusterRepo) Recover(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `UPDATE clusters SET status='INTERRUPTED',error='Daemon restarted during a cluster operation; inspect saved nodes before reviewing a retry.' WHERE status='APPLYING'`)
	return err
}
