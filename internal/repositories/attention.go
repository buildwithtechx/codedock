package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"codedock/internal/models"
)

type AttentionRepository interface {
	Upsert(ctx context.Context, issue *models.AttentionIssue) error
	Get(ctx context.Context, id string) (*models.AttentionIssue, error)
	List(ctx context.Context, orgID, status string, limit int) ([]*models.AttentionIssue, error)
	CountOpen(ctx context.Context, orgID string) (int, error)
	SetStatus(ctx context.Context, id, orgID string, status models.AttentionStatus) error
	ResolveMissing(ctx context.Context, orgID string, live []string) error
	AppendDetail(ctx context.Context, id, line string) error
}

type postgresAttentionRepository struct {
	db *sql.DB
}

func NewAttentionRepository(db *sql.DB) AttentionRepository {
	return &postgresAttentionRepository{db: db}
}

const attentionColumns = `id, organization_id, kind, subject, project_id, service_id, severity, status, title, detail, remediation, action, action_params, occurrences, first_seen, last_seen, updated_at`

func (r *postgresAttentionRepository) Upsert(ctx context.Context, issue *models.AttentionIssue) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO attention_issues (id, organization_id, kind, subject, project_id, service_id, severity, status, title, detail, remediation, action, action_params, occurrences, first_seen, last_seen, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, 'open', $8, $9, $10, $11, $12, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT(organization_id, kind, subject) DO UPDATE SET
			project_id = excluded.project_id, service_id = excluded.service_id, severity = excluded.severity,
			status = CASE WHEN attention_issues.status = 'resolved' THEN 'open' ELSE attention_issues.status END,
			title = excluded.title, detail = excluded.detail, remediation = excluded.remediation,
			action = excluded.action, action_params = excluded.action_params,
			occurrences = attention_issues.occurrences + 1, last_seen = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
	`, issue.ID, issue.OrganizationID, issue.Kind, issue.Subject, issue.ProjectID, issue.ServiceID, string(issue.Severity), issue.Title, issue.Detail, issue.Remediation, issue.Action, issue.ActionParams)
	if err != nil {
		return fmt.Errorf("upsert attention issue: %w", err)
	}
	return nil
}

func (r *postgresAttentionRepository) Get(ctx context.Context, id string) (*models.AttentionIssue, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+attentionColumns+` FROM attention_issues WHERE id = $1`, id)
	return r.scanRow(row)
}

func (r *postgresAttentionRepository) List(ctx context.Context, orgID, status string, limit int) ([]*models.AttentionIssue, error) {
	clause := `WHERE organization_id = $1`
	args := []any{orgID}
	if status != "" {
		clause += fmt.Sprintf(` AND status = $%d`, len(args)+1)
		args = append(args, status)
	}
	args = append(args, limit)
	rows, err := r.db.QueryContext(ctx, `SELECT `+attentionColumns+` FROM attention_issues `+clause+fmt.Sprintf(` ORDER BY last_seen DESC LIMIT $%d`, len(args)), args...)
	if err != nil {
		return nil, fmt.Errorf("list attention issues: %w", err)
	}
	defer rows.Close()
	var issues []*models.AttentionIssue
	for rows.Next() {
		var issue models.AttentionIssue
		var severity, statusValue string
		if err := rows.Scan(&issue.ID, &issue.OrganizationID, &issue.Kind, &issue.Subject, &issue.ProjectID, &issue.ServiceID, &severity, &statusValue, &issue.Title, &issue.Detail, &issue.Remediation, &issue.Action, &issue.ActionParams, &issue.Occurrences, &issue.FirstSeen, &issue.LastSeen, &issue.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan attention issue: %w", err)
		}
		issue.Severity = models.AttentionSeverity(severity)
		issue.Status = models.AttentionStatus(statusValue)
		issues = append(issues, &issue)
	}
	return issues, rows.Err()
}

func (r *postgresAttentionRepository) CountOpen(ctx context.Context, orgID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM attention_issues WHERE organization_id = $1 AND status = 'open'`, orgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count open issues: %w", err)
	}
	return count, nil
}

func (r *postgresAttentionRepository) SetStatus(ctx context.Context, id, orgID string, status models.AttentionStatus) error {
	result, err := r.db.ExecContext(ctx, `UPDATE attention_issues SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND organization_id = $3`, string(status), id, orgID)
	if err != nil {
		return fmt.Errorf("set issue status: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected != 1 {
		return fmt.Errorf("attention issue not found")
	}
	return nil
}

func (r *postgresAttentionRepository) ResolveMissing(ctx context.Context, orgID string, live []string) error {
	if len(live) == 0 {
		_, err := r.db.ExecContext(ctx, `UPDATE attention_issues SET status = 'resolved', updated_at = CURRENT_TIMESTAMP WHERE organization_id = $1 AND status = 'open'`, orgID)
		return err
	}
	placeholders := ""
	args := []any{orgID}
	for index, key := range live {
		if index > 0 {
			placeholders += ","
		}
		placeholders += fmt.Sprintf("$%d", index+2)
		args = append(args, key)
	}
	_, err := r.db.ExecContext(ctx, `UPDATE attention_issues SET status = 'resolved', updated_at = CURRENT_TIMESTAMP WHERE organization_id = $1 AND status = 'open' AND (kind || ':' || subject) NOT IN (`+placeholders+`)`, args...)
	return err
}

func (r *postgresAttentionRepository) AppendDetail(ctx context.Context, id, line string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE attention_issues SET detail = right(detail || $1, 8192), updated_at = CURRENT_TIMESTAMP WHERE id = $2`, "\n"+line, id)
	if err != nil {
		return fmt.Errorf("append issue detail: %w", err)
	}
	return nil
}

func (r *postgresAttentionRepository) scanRow(row *sql.Row) (*models.AttentionIssue, error) {
	var issue models.AttentionIssue
	var severity, status string
	if err := row.Scan(&issue.ID, &issue.OrganizationID, &issue.Kind, &issue.Subject, &issue.ProjectID, &issue.ServiceID, &severity, &status, &issue.Title, &issue.Detail, &issue.Remediation, &issue.Action, &issue.ActionParams, &issue.Occurrences, &issue.FirstSeen, &issue.LastSeen, &issue.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan attention issue: %w", err)
	}
	issue.Severity = models.AttentionSeverity(severity)
	issue.Status = models.AttentionStatus(status)
	return &issue, nil
}
