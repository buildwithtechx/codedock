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
	ListFiltered(ctx context.Context, prefixes []string, limit, offset int) ([]models.AuditLog, error)
	Facets(ctx context.Context) (map[string]int, int, error)
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
	return r.ListFiltered(ctx, nil, limit, offset)
}

func (r *AuditLogRepo) ListFiltered(ctx context.Context, prefixes []string, limit, offset int) ([]models.AuditLog, error) {
	query := `
		SELECT id, user_id, action, resource, details, ip_address, created_at
		FROM audit_logs
	`
	args := []any{}
	if len(prefixes) > 0 {
		conditions := make([]string, 0, len(prefixes))
		for i, prefix := range prefixes {
			conditions = append(conditions, fmt.Sprintf("action LIKE $%d", i+1))
			args = append(args, prefix+"%")
		}
		query += " WHERE " + strings.Join(conditions, " OR ")
	}
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

func (r *AuditLogRepo) Facets(ctx context.Context) (map[string]int, int, error) {
	caseParts := []string{}
	for _, category := range models.AuditCategories() {
		prefixes := models.CategoryPrefixes(category.ID)
		if len(prefixes) == 0 {
			continue
		}
		likes := make([]string, 0, len(prefixes))
		for _, prefix := range prefixes {
			likes = append(likes, fmt.Sprintf("action LIKE '%s%%'", strings.ReplaceAll(prefix, "'", "''")))
		}
		caseParts = append(caseParts, fmt.Sprintf("WHEN %s THEN '%s'", strings.Join(likes, " OR "), category.ID))
	}
	query := fmt.Sprintf(`SELECT CASE %s ELSE 'system' END AS category, COUNT(*) FROM audit_logs GROUP BY 1`, strings.Join(caseParts, " "))
	rows, err := r.db.QueryContext(ctx, query)
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
