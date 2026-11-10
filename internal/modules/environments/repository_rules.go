package environments

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type postgresRouteRuleRepository struct {
	db *sql.DB
}

func NewRouteRuleRepository(db *sql.DB) RouteRuleRepository {
	return &postgresRouteRuleRepository{db: db}
}

type scannable interface {
	Scan(dest ...any) error
}

func (r *postgresRouteRuleRepository) Create(ctx context.Context, rule *RouteRule) error {
	if rule.ID == "" {
		rule.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO route_rules (id, service_id, name, enabled, rule_type, spec_json, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		rule.ID, rule.ServiceID, rule.Name, boolToInt(rule.Enabled),
		rule.RuleType, rule.SpecJSON, rule.CreatedAt, rule.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create route rule: %w", err)
	}
	return nil
}

func (r *postgresRouteRuleRepository) GetByID(ctx context.Context, id string) (*RouteRule, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, service_id, name, enabled, rule_type, spec_json, created_at, updated_at
		FROM route_rules WHERE id = $1`, id)
	return r.scan(row)
}

func (r *postgresRouteRuleRepository) ListByService(ctx context.Context, serviceID string) ([]*RouteRule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, service_id, name, enabled, rule_type, spec_json, created_at, updated_at
		FROM route_rules WHERE service_id = $1 ORDER BY created_at ASC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list route rules: %w", err)
	}
	defer rows.Close()
	var rules []*RouteRule
	for rows.Next() {
		rule, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("route rules iteration: %w", err)
	}
	return rules, nil
}

func (r *postgresRouteRuleRepository) Update(ctx context.Context, id string, name *string, enabled *bool, specJSON *string) error {
	now := time.Now().UTC()
	updates := []string{"updated_at = $1"}
	args := []any{now}

	if name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, *name)
	}
	if enabled != nil {
		updates = append(updates, fmt.Sprintf("enabled = $%d", len(args)+1))
		args = append(args, boolToInt(*enabled))
	}
	if specJSON != nil {
		updates = append(updates, fmt.Sprintf("spec_json = $%d", len(args)+1))
		args = append(args, *specJSON)
	}

	if len(updates) == 1 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE route_rules SET %s WHERE id = $%d", strings.Join(updates, ", "), len(args))
	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update route rule: %w", err)
	}
	return nil
}

func (r *postgresRouteRuleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM route_rules WHERE id = $1`, id)
	return err
}

func (r *postgresRouteRuleRepository) scan(s scannable) (*RouteRule, error) {
	var rule RouteRule
	var enabledInt int
	err := s.Scan(
		&rule.ID, &rule.ServiceID, &rule.Name, &enabledInt,
		&rule.RuleType, &rule.SpecJSON, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRouteRuleNotFound
		}
		return nil, fmt.Errorf("scan route rule: %w", err)
	}
	rule.Enabled = enabledInt == 1
	return &rule, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
