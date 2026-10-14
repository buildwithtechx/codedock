package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"codedock/internal/models"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *models.AuditLog) error
	List(ctx context.Context, limit, offset int) ([]models.AuditLog, error)
	ListFiltered(ctx context.Context, filter models.AuditLogFilter, limit, offset int) ([]models.AuditLog, error)
	Facets(ctx context.Context, filter models.AuditLogFilter) (map[string]int, int, error)
	GetSettings(ctx context.Context, orgID string) (*models.AuditSettings, error)
	UpdateSettings(ctx context.Context, settings *models.AuditSettings) error
}

type AuditLogRepo struct {
	db *sql.DB
}

func NewAuditLogRepo(db *sql.DB) *AuditLogRepo {
	return &AuditLogRepo{db: db}
}

func (r *AuditLogRepo) Create(ctx context.Context, log *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, user_id, action, resource, details, ip_address)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.UserID, log.Action, log.Resource, log.Details, log.IPAddress)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}

func (r *AuditLogRepo) List(ctx context.Context, limit, offset int) ([]models.AuditLog, error) {
	return r.ListFiltered(ctx, models.AuditLogFilter{}, limit, offset)
}

func escapeLikePrefix(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

func auditLogFilterClause(filter models.AuditLogFilter) (string, []any) {
	conditions := []string{}
	args := []any{}
	if len(filter.Prefixes) > 0 {
		likes := make([]string, 0, len(filter.Prefixes))
		for _, prefix := range filter.Prefixes {
			args = append(args, escapeLikePrefix(prefix)+"%")
			likes = append(likes, fmt.Sprintf("action ILIKE $%d", len(args)))
		}
		conditions = append(conditions, "("+strings.Join(likes, " OR ")+")")
	}
	if q := strings.TrimSpace(filter.Q); q != "" {
		args = append(args, ContainsLikePattern(q))
		n := len(args)
		conditions = append(conditions, fmt.Sprintf(`(action ILIKE $%d OR resource ILIKE $%d OR EXISTS (SELECT 1 FROM users WHERE users.id = audit_logs.user_id AND (users.email ILIKE $%d OR users.name ILIKE $%d)))`, n, n, n, n))
	}
	if len(filter.QResourceIDs) > 0 {
		marks := make([]string, 0, len(filter.QResourceIDs))
		for _, id := range filter.QResourceIDs {
			args = append(args, id)
			marks = append(marks, fmt.Sprintf("$%d", len(args)))
		}
		conditions = append(conditions, "resource IN ("+strings.Join(marks, ", ")+")")
	}
	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *AuditLogRepo) ListFiltered(ctx context.Context, filter models.AuditLogFilter, limit, offset int) ([]models.AuditLog, error) {
	where, args := auditLogFilterClause(filter)
	query := `
		SELECT id, user_id, action, resource, details, ip_address, created_at
		FROM audit_logs
	` + where
	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var details, ipAddress sql.NullString
		if err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Resource, &details, &ipAddress, &log.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		if details.Valid {
			log.Details = details.String
		}
		if ipAddress.Valid {
			log.IPAddress = ipAddress.String
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return logs, nil
}

func (r *AuditLogRepo) Facets(ctx context.Context, filter models.AuditLogFilter) (map[string]int, int, error) {
	caseParts := []string{}
	for _, category := range models.AuditCategories() {
		prefixes := models.CategoryPrefixes(category.ID)
		if len(prefixes) == 0 {
			continue
		}
		likes := make([]string, 0, len(prefixes))
		for _, prefix := range prefixes {
			likes = append(likes, fmt.Sprintf("action ILIKE '%s%%'", strings.ReplaceAll(escapeLikePrefix(prefix), "'", "''")))
		}
		caseParts = append(caseParts, fmt.Sprintf("WHEN %s THEN '%s'", strings.Join(likes, " OR "), category.ID))
	}
	where, args := auditLogFilterClause(filter)
	query := fmt.Sprintf(`SELECT CASE %s ELSE 'system' END AS category, COUNT(*) FROM audit_logs%s GROUP BY 1`, strings.Join(caseParts, " "), where)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to load audit facets: %w", err)
	}
	defer rows.Close()

	counts := map[string]int{}
	total := 0
	for rows.Next() {
		var category string
		var count int
		if err := rows.Scan(&category, &count); err != nil {
			return nil, 0, fmt.Errorf("failed to scan audit facet: %w", err)
		}
		counts[category] = count
		total += count
	}
	return counts, total, nil
}

func (r *AuditLogRepo) GetSettings(ctx context.Context, orgID string) (*models.AuditSettings, error) {
	query := `
		SELECT organization_id, enabled, retention_days, updated_at
		FROM audit_settings
		WHERE organization_id = $1
	`
	var s models.AuditSettings
	var updatedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(&s.OrganizationID, &s.Enabled, &s.RetentionDays, &updatedAt)
	if err == sql.ErrNoRows {
		return &models.AuditSettings{
			OrganizationID: orgID,
			Enabled:        true,
			RetentionDays:  90,
		}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get audit settings: %w", err)
	}
	if updatedAt.Valid {
		s.UpdatedAt = updatedAt.Time.String()
	}
	return &s, nil
}

func (r *AuditLogRepo) UpdateSettings(ctx context.Context, settings *models.AuditSettings) error {
	query := `
		INSERT INTO audit_settings (organization_id, enabled, retention_days, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (organization_id)
		DO UPDATE SET enabled = $2, retention_days = $3, updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.ExecContext(ctx, query, settings.OrganizationID, settings.Enabled, settings.RetentionDays)
	if err != nil {
		return fmt.Errorf("failed to update audit settings: %w", err)
	}
	return nil
}
