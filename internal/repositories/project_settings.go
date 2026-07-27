package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"codedock.run/codedock/internal/utils"

	"github.com/jmoiron/sqlx"

	"codedock.run/codedock/internal/models"
)

type ProjectSettingsRepository interface {
	CreateToken(ctx context.Context, t *models.ProjectToken, fullToken string) error
	ListTokensByProject(ctx context.Context, projectID string) ([]*models.ProjectToken, error)
	DeleteToken(ctx context.Context, id, projectID string) error
	GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error)
	UpdateTokenLastUsed(ctx context.Context, id string) error
}

type ProjectSettingsRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewProjectSettingsRepo(db *sql.DB) *ProjectSettingsRepo {
	return &ProjectSettingsRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *ProjectSettingsRepo) deleteByIDAndProject(ctx context.Context, table, id, projectID, entityName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	query := fmt.Sprintf("DELETE FROM %s WHERE id = ? AND project_id = ?", table)
	res, err := r.db.ExecContext(ctx, query, id, projectID)
	if err != nil {
		return fmt.Errorf("delete %s: %w", entityName, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("%s not found", entityName)
	}
	return nil
}

func (r *ProjectSettingsRepo) CreateToken(ctx context.Context, t *models.ProjectToken, fullToken string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	scopesStr := strings.Join(t.Scopes, ",")
	ipStr := strings.Join(t.IPAllowlist, ",")
	var expiresAtVal any
	if t.ExpiresAt != nil {
		expiresAtVal = t.ExpiresAt.Format(time.RFC3339)
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO project_tokens (id, project_id, environment_id, name, token_prefix, token_hash, scopes, ip_allowlist, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID, t.ProjectID, t.EnvironmentID, t.Name, t.TokenPrefix, fullToken, scopesStr, ipStr, expiresAtVal, t.CreatedAt.Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	return nil
}

func (r *ProjectSettingsRepo) ListTokensByProject(ctx context.Context, projectID string) ([]*models.ProjectToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, environment_id, name, token_prefix, scopes, ip_allowlist, expires_at, created_at
		 FROM project_tokens WHERE project_id = ? ORDER BY created_at DESC`, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tokens: %w", err)
	}
	defer rows.Close()
	var out []*models.ProjectToken
	for rows.Next() {
		var t models.ProjectToken
		var scopesStr, ipStr string
		var expiresAtStr sql.NullString
		var createdAtStr string
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.EnvironmentID, &t.Name, &t.TokenPrefix, &scopesStr, &ipStr, &expiresAtStr, &createdAtStr); err != nil {
			return nil, fmt.Errorf("scan token: %w", err)
		}
		if scopesStr != "" {
			t.Scopes = strings.Split(scopesStr, ",")
		} else {
			t.Scopes = []string{}
		}
		if ipStr != "" {
			t.IPAllowlist = strings.Split(ipStr, ",")
		} else {
			t.IPAllowlist = []string{}
		}
		if expiresAtStr.Valid && expiresAtStr.String != "" {
			parsed, _ := time.Parse(time.RFC3339, expiresAtStr.String)
			t.ExpiresAt = &parsed
		}
		parsedCreated, _ := time.Parse(time.RFC3339, createdAtStr)
		t.CreatedAt = parsedCreated
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *ProjectSettingsRepo) DeleteToken(ctx context.Context, id, projectID string) error {
	return r.deleteByIDAndProject(ctx, "project_tokens", id, projectID, "token")
}

func (r *ProjectSettingsRepo) GetTokenByHash(ctx context.Context, tokenHash string) (*models.ProjectToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var t models.ProjectToken
	var scopesStr, ipStr string
	var expiresAtStr sql.NullString
	var createdAtStr string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, environment_id, name, token_prefix, scopes, ip_allowlist, expires_at, created_at
		 FROM project_tokens WHERE token_hash = ?`, tokenHash).
		Scan(&t.ID, &t.ProjectID, &t.EnvironmentID, &t.Name, &t.TokenPrefix, &scopesStr, &ipStr, &expiresAtStr, &createdAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, utils.NewNotFoundError("Token", tokenHash)
		}
		return nil, fmt.Errorf("get token by hash: %w", err)
	}
	if scopesStr != "" {
		t.Scopes = strings.Split(scopesStr, ",")
	}
	if ipStr != "" {
		t.IPAllowlist = strings.Split(ipStr, ",")
	}
	if expiresAtStr.Valid && expiresAtStr.String != "" {
		parsed, _ := time.Parse(time.RFC3339, expiresAtStr.String)
		t.ExpiresAt = &parsed
	}
	parsedCreated, _ := time.Parse(time.RFC3339, createdAtStr)
	t.CreatedAt = parsedCreated
	return &t, nil
}

func (r *ProjectSettingsRepo) UpdateTokenLastUsed(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, err := r.db.ExecContext(ctx, `UPDATE project_tokens SET last_used_at = ? WHERE id = ?`, time.Now().Format(time.RFC3339), id)
	return err
}
