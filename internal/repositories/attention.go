package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"codedock.run/codedock/internal/models"
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

type sqliteAttentionRepository struct {
	db *sql.DB
}

func NewAttentionRepository(db *sql.DB) AttentionRepository {
	return &sqliteAttentionRepository{db: db}
}

const attentionColumns = `id, organization_id, kind, subject, project_id, service_id, severity, status, title, detail, remediation, action, action_params, occurrences, first_seen, last_seen, updated_at`

func (r *sqliteAttentionRepository) Upsert(ctx context.Context, issue *models.AttentionIssue) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO attention_issues (id, organization_id, kind, subject, project_id, service_id, severity, status, title, detail, remediation, action, action_params, occurrences, first_seen, last_seen, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'open', ?, ?, ?, ?, ?, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
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

func (r *sqliteAttentionRepository) Get(ctx context.Context, id string) (*models.AttentionIssue, error) {
	row := r.db.QueryRowContext(ctx, `SELECT `+attentionColumns+` FROM attention_issues WHERE id = ?`, id)
	return r.scanRow(row)
}

func (r *sqliteAttentionRepository) List(ctx context.Context, orgID, status string, limit int) ([]*models.AttentionIssue, error) {
	clause := `WHERE organization_id = ?`
	args := []any{orgID}
	if status != "" {
		clause += ` AND status = ?`
		args = append(args, status)
	}
	args = append(args, limit)
	rows, err := r.db.QueryContext(ctx, `SELECT `+attentionColumns+` FROM attention_issues `+clause+` ORDER BY last_seen DESC LIMIT ?`, args...)
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

func (r *sqliteAttentionRepository) CountOpen(ctx context.Context, orgID string) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM attention_issues WHERE organization_id = ? AND status = 'open'`, orgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count open issues: %w", err)
	}
	return count, nil
}

func (r *sqliteAttentionRepository) SetStatus(ctx context.Context, id, orgID string, status models.AttentionStatus) error {
	result, err := r.db.ExecContext(ctx, `UPDATE attention_issues SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND organization_id = ?`, string(status), id, orgID)
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

func (r *sqliteAttentionRepository) ResolveMissing(ctx context.Context, orgID string, live []string) error {
	if len(live) == 0 {
		_, err := r.db.ExecContext(ctx, `UPDATE attention_issues SET status = 'resolved', updated_at = CURRENT_TIMESTAMP WHERE organization_id = ? AND status = 'open'`, orgID)
		return err
	}
	placeholders := ""
	args := []any{orgID}
	for index, key := range live {
		if index > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, key)
	}
	_, err := r.db.ExecContext(ctx, `UPDATE attention_issues SET status = 'resolved', updated_at = CURRENT_TIMESTAMP WHERE organization_id = ? AND status = 'open' AND (kind || ':' || subject) NOT IN (`+placeholders+`)`, args...)
	return err
}

func (r *sqliteAttentionRepository) AppendDetail(ctx context.Context, id, line string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE attention_issues SET detail = substr(detail || ?, -8192), updated_at = CURRENT_TIMESTAMP WHERE id = ?`, "\n"+line, id)
	if err != nil {
		return fmt.Errorf("append issue detail: %w", err)
	}
	return nil
}

func (r *sqliteAttentionRepository) scanRow(row *sql.Row) (*models.AttentionIssue, error) {
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
