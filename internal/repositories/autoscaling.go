package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"codedock.run/codedock/internal/models"
)

type AutoscalingRepo struct{ db *sql.DB }

func NewAutoscalingRepo(db *sql.DB) *AutoscalingRepo { return &AutoscalingRepo{db: db} }

const scalingColumns = `service_id,enabled,min_replicas,max_replicas,scale_up_cpu,scale_down_cpu,cooldown_seconds,last_cpu,last_decision,last_evaluated_at,last_scaled_at,updated_at`

type scalingScanner interface{ Scan(...any) error }

func scanScaling(row scalingScanner) (*models.AutoscalingPolicy, error) {
	p := &models.AutoscalingPolicy{}
	err := row.Scan(&p.ServiceID, &p.Enabled, &p.MinReplicas, &p.MaxReplicas, &p.ScaleUpCPU, &p.ScaleDownCPU, &p.CooldownSeconds, &p.LastCPU, &p.LastDecision, &p.LastEvaluatedAt, &p.LastScaledAt, &p.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("scan autoscaling policy: %w", err)
	}
	return p, nil
}

func (r *AutoscalingRepo) Get(ctx context.Context, id string) (*models.AutoscalingPolicy, error) {
	p, err := scanScaling(r.db.QueryRowContext(ctx, `SELECT `+scalingColumns+` FROM autoscaling_policies WHERE service_id=?`, id))
	if errors.Is(err, sql.ErrNoRows) {
		return models.DefaultAutoscalingPolicy(id), nil
	}
	return p, err
}

func (r *AutoscalingRepo) Save(ctx context.Context, p *models.AutoscalingPolicy) error {
	p.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.db.ExecContext(ctx, `INSERT INTO autoscaling_policies(service_id,enabled,min_replicas,max_replicas,scale_up_cpu,scale_down_cpu,cooldown_seconds,updated_at) VALUES(?,?,?,?,?,?,?,?) ON CONFLICT(service_id) DO UPDATE SET enabled=excluded.enabled,min_replicas=excluded.min_replicas,max_replicas=excluded.max_replicas,scale_up_cpu=excluded.scale_up_cpu,scale_down_cpu=excluded.scale_down_cpu,cooldown_seconds=excluded.cooldown_seconds,updated_at=excluded.updated_at`, p.ServiceID, p.Enabled, p.MinReplicas, p.MaxReplicas, p.ScaleUpCPU, p.ScaleDownCPU, p.CooldownSeconds, p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("save autoscaling policy: %w", err)
	}
	return nil
}

func (r *AutoscalingRepo) ListEnabled(ctx context.Context) ([]*models.AutoscalingPolicy, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+scalingColumns+` FROM autoscaling_policies WHERE enabled=1`)
	if err != nil {
		return nil, fmt.Errorf("list enabled autoscaling policies: %w", err)
	}
	defer rows.Close()
	var list []*models.AutoscalingPolicy
	for rows.Next() {
		p, err := scanScaling(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *AutoscalingRepo) Observe(ctx context.Context, p *models.AutoscalingPolicy, cpu float64, decision string, scaled bool) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	_, err := r.db.ExecContext(ctx, `UPDATE autoscaling_policies SET last_cpu=?,last_decision=?,last_evaluated_at=?,last_scaled_at=CASE WHEN ? THEN ? ELSE last_scaled_at END WHERE service_id=?`, cpu, decision, now, scaled, now, p.ServiceID)
	if err != nil {
		return fmt.Errorf("record autoscaling decision: %w", err)
	}
	return nil
}

func (r *AutoscalingRepo) SetReplicas(ctx context.Context, p *models.AutoscalingPolicy, from, to int) (bool, error) {
	result, err := r.db.ExecContext(ctx, `UPDATE app_services SET replicas=?,updated_at=? WHERE id=? AND CASE WHEN replicas<1 THEN 1 ELSE replicas END=? AND status='running' AND EXISTS(SELECT 1 FROM projects WHERE projects.id=app_services.project_id AND COALESCE(server_id,'')='') AND NOT EXISTS(SELECT 1 FROM service_runtimes WHERE service_id=app_services.id AND (runtime_kind!='docker' OR encrypted_journal!='')) AND EXISTS(SELECT 1 FROM autoscaling_policies WHERE service_id=app_services.id AND enabled=1 AND updated_at=?) AND NOT EXISTS(SELECT 1 FROM deployments WHERE service_id=app_services.id AND status IN ('pending','PREPARING','CLONING','PULLING','BUILDING','STARTING','READINESS','ROUTING'))`, to, time.Now().UTC(), p.ServiceID, from, p.UpdatedAt)
	if err != nil {
		return false, fmt.Errorf("reserve autoscaling replicas: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("read autoscaling reservation: %w", err)
	}
	return count == 1, nil
}

func (r *AutoscalingRepo) RollbackReplicas(ctx context.Context, id string, from, to int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE app_services SET replicas=? WHERE id=? AND replicas=?`, to, id, from)
	if err != nil {
		return fmt.Errorf("rollback autoscaling replicas: %w", err)
	}
	return nil
}
