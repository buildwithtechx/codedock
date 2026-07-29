package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
)

type RouteRuleRepository interface {
	Create(ctx context.Context, rule *models.RouteRule) error
	GetByID(ctx context.Context, id string) (*models.RouteRule, error)
	ListByService(ctx context.Context, serviceID string) ([]*models.RouteRule, error)
	Update(ctx context.Context, id string, name *string, enabled *bool, specJSON *string) error
	Delete(ctx context.Context, id string) error
}

type sqliteRouteRuleRepository struct {
	db *sql.DB
}

func NewRouteRuleRepository(db *sql.DB) RouteRuleRepository {
	return &sqliteRouteRuleRepository{db: db}
}

func (r *sqliteRouteRuleRepository) Create(ctx context.Context, rule *models.RouteRule) error {
	if rule.ID == "" {
		rule.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	rule.CreatedAt = now
	rule.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO route_rules (id, service_id, name, enabled, rule_type, spec_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		rule.ID, rule.ServiceID, rule.Name, boolToInt(rule.Enabled),
		rule.RuleType, rule.SpecJSON, rule.CreatedAt, rule.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create route rule: %w", err)
	}
	return nil
}

func (r *sqliteRouteRuleRepository) GetByID(ctx context.Context, id string) (*models.RouteRule, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, service_id, name, enabled, rule_type, spec_json, created_at, updated_at
		FROM route_rules WHERE id = ?`, id)
	return r.scan(row)
}

func (r *sqliteRouteRuleRepository) ListByService(ctx context.Context, serviceID string) ([]*models.RouteRule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, service_id, name, enabled, rule_type, spec_json, created_at, updated_at
		FROM route_rules WHERE service_id = ? ORDER BY created_at ASC`, serviceID)
	if err != nil {
		return nil, fmt.Errorf("list route rules: %w", err)
	}
	defer rows.Close()
	var rules []*models.RouteRule
	for rows.Next() {
		rule, err := r.scan(rows)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *sqliteRouteRuleRepository) Update(ctx context.Context, id string, name *string, enabled *bool, specJSON *string) error {
	now := time.Now().UTC()
	if name != nil {
		if _, err := r.db.ExecContext(ctx, `UPDATE route_rules SET name = ?, updated_at = ? WHERE id = ?`, *name, now, id); err != nil {
			return err
		}
	}
	if enabled != nil {
		if _, err := r.db.ExecContext(ctx, `UPDATE route_rules SET enabled = ?, updated_at = ? WHERE id = ?`, boolToInt(*enabled), now, id); err != nil {
			return err
		}
	}
	if specJSON != nil {
		if _, err := r.db.ExecContext(ctx, `UPDATE route_rules SET spec_json = ?, updated_at = ? WHERE id = ?`, *specJSON, now, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *sqliteRouteRuleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM route_rules WHERE id = ?`, id)
	return err
}

func (r *sqliteRouteRuleRepository) scan(s scannable) (*models.RouteRule, error) {
	var rule models.RouteRule
	var enabledInt int
	err := s.Scan(
		&rule.ID, &rule.ServiceID, &rule.Name, &enabledInt,
		&rule.RuleType, &rule.SpecJSON, &rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
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
