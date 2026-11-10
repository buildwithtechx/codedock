package autoscaling

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type autoscalingRepo struct {
	db *sqlx.DB
}

func NewAutoscalingRepository(db *sql.DB) AutoscalingRepository {
	return &autoscalingRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *autoscalingRepo) GetByServiceID(ctx context.Context, serviceID string) (*AutoscalingPolicy, error) {
	var p AutoscalingPolicy
	query := `SELECT service_id, enabled, min_replicas, max_replicas, scale_up_cpu, scale_down_cpu, cooldown_seconds, last_cpu, last_decision, COALESCE(to_char(last_evaluated_at, 'YYYY-MM-DD"T"HH24:MI:SSZ'), '') as last_evaluated_at, COALESCE(to_char(last_scaled_at, 'YYYY-MM-DD"T"HH24:MI:SSZ'), '') as last_scaled_at, COALESCE(to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SSZ'), '') as updated_at FROM autoscaling_policies WHERE service_id = $1`
	err := r.db.GetContext(ctx, &p, query, serviceID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	p.Supported = true
	return &p, nil
}

func (r *autoscalingRepo) Upsert(ctx context.Context, p *AutoscalingPolicy) error {
	query := `INSERT INTO autoscaling_policies (service_id, enabled, min_replicas, max_replicas, scale_up_cpu, scale_down_cpu, cooldown_seconds, last_cpu, last_decision, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP)
		ON CONFLICT(service_id) DO UPDATE SET
			enabled = excluded.enabled,
			min_replicas = excluded.min_replicas,
			max_replicas = excluded.max_replicas,
			scale_up_cpu = excluded.scale_up_cpu,
			scale_down_cpu = excluded.scale_down_cpu,
			cooldown_seconds = excluded.cooldown_seconds,
			last_cpu = excluded.last_cpu,
			last_decision = excluded.last_decision,
			updated_at = CURRENT_TIMESTAMP`
	_, err := r.db.ExecContext(ctx, query, p.ServiceID, p.Enabled, p.MinReplicas, p.MaxReplicas, p.ScaleUpCPU, p.ScaleDownCPU, p.CooldownSeconds, p.LastCPU, p.LastDecision)
	return err
}

type serverlessRepo struct {
	db *sqlx.DB
}

func NewServerlessRepository(db *sql.DB) ServerlessRepository {
	return &serverlessRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *serverlessRepo) SaveCode(ctx context.Context, serviceID, runtime, codeContent string) (*ServerlessFunctionCode, error) {
	existing, err := r.GetCodeByServiceID(ctx, serviceID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	now := time.Now()
	if existing != nil {
		query := `UPDATE serverless_functions_code SET runtime = $1, code_content = $2, updated_at = $3 WHERE service_id = $4`
		_, err := r.db.ExecContext(ctx, query, runtime, codeContent, now, serviceID)
		if err != nil {
			return nil, err
		}
		existing.Runtime = runtime
		existing.CodeContent = codeContent
		existing.UpdatedAt = now
		return existing, nil
	}

	id := uuid.New().String()
	query := `INSERT INTO serverless_functions_code (id, service_id, runtime, code_content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.db.ExecContext(ctx, query, id, serviceID, runtime, codeContent, now, now)
	if err != nil {
		return nil, err
	}

	return &ServerlessFunctionCode{
		ID:          id,
		ServiceID:   serviceID,
		Runtime:     runtime,
		CodeContent: codeContent,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *serverlessRepo) GetCodeByServiceID(ctx context.Context, serviceID string) (*ServerlessFunctionCode, error) {
	query := `SELECT id, service_id, runtime, code_content, created_at, updated_at FROM serverless_functions_code WHERE service_id = $1`
	var code ServerlessFunctionCode
	err := r.db.GetContext(ctx, &code, query, serviceID)
	if err != nil {
		return nil, err
	}
	return &code, nil
}
