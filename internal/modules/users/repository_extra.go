package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"codedock/internal/utils"
)

func (r *UserRepo) GetUserTOTPSecret(ctx context.Context, userID string) (string, []string, error) {
	var secret string
	var recovery string
	err := r.db.QueryRowContext(ctx, `SELECT COALESCE(totp_secret, ''), COALESCE(recovery_codes, '') FROM users WHERE id = $1`, userID).Scan(&secret, &recovery)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get totp secret: %w", err)
	}
	var codes []string
	if recovery != "" {
		for _, part := range strings.Split(recovery, ",") {
			if part = strings.TrimSpace(part); part != "" {
				codes = append(codes, part)
			}
		}
	}
	return secret, codes, nil
}

func (r *UserRepo) FindFoundingAdmin(ctx context.Context) (*User, error) {
	var u User
	query := `SELECT id, name, email, password_hash, role, plan_type, email_verified, totp_enabled, created_at, updated_at
	          FROM users
	          WHERE role = 'admin'
	          ORDER BY created_at ASC
	          LIMIT 1`
	err := r.db.GetContext(ctx, &u, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepo) GetAppearance(ctx context.Context, userID string) (*AppearancePreferences, error) {
	var raw sql.NullString
	err := r.db.QueryRowContext(ctx, `SELECT appearance FROM users WHERE id = $1`, userID).Scan(&raw)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if !raw.Valid || strings.TrimSpace(raw.String) == "" {
		return &AppearancePreferences{
			ActiveTheme: "default",
			Mode:        "dark",
		}, nil
	}
	var pref AppearancePreferences
	if err := json.Unmarshal([]byte(raw.String), &pref); err != nil {
		return nil, fmt.Errorf("decode appearance preferences: %w", err)
	}
	return &pref, nil
}

func (r *UserRepo) SetAppearance(ctx context.Context, userID string, pref *AppearancePreferences) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var val any
	if pref != nil {
		data, err := json.Marshal(pref)
		if err != nil {
			return fmt.Errorf("encode appearance preferences: %w", err)
		}
		val = string(data)
	}

	res, err := r.db.ExecContext(ctx, `UPDATE users SET appearance = $1, updated_at = NOW() WHERE id = $2`, val, userID)
	if err != nil {
		return fmt.Errorf("update appearance: %w", err)
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return utils.NewNotFoundError("User", userID)
	}
	return nil
}

func (r *UserRepo) GetPATByHash(ctx context.Context, hash string) (*PersonalAccessToken, error) {
	var pat PersonalAccessToken
	err := r.db.QueryRowContext(ctx, `SELECT id, user_id, name, token_hash, prefix, access_level, project_scope, allowed_projects, expires_at, created_at FROM personal_access_tokens WHERE token_hash = $1`, hash).Scan(&pat.ID, &pat.UserID, &pat.Name, &pat.TokenHash, &pat.Prefix, &pat.AccessLevel, &pat.ProjectScope, &pat.AllowedProjects, &pat.ExpiresAt, &pat.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("Personal access token", "")
	}
	if err != nil {
		return nil, fmt.Errorf("load personal token: %w", err)
	}
	return &pat, nil
}

type PersonalTokenResources struct{ db *sql.DB }

func NewPersonalTokenResources(db *sql.DB) *PersonalTokenResources {
	return &PersonalTokenResources{db: db}
}

func (r *PersonalTokenResources) ProjectForResource(ctx context.Context, kind, id string) (string, error) {
	queries := map[string]string{
		"operations":        `SELECT project_id FROM operations WHERE id=$1`,
		"backup-operations": `SELECT project_id FROM operations WHERE id=$1`,
		"stacks":            `SELECT project_id FROM compose_stacks WHERE id=$1`,
		"projects":          `SELECT id FROM projects WHERE id = $1`,
		"apps":              `SELECT project_id FROM app_services WHERE id = $1`,
		"services":          `SELECT project_id FROM app_services WHERE id = $1`,
		"environments":      `SELECT project_id FROM environments WHERE id = $1`,
		"databases":         `SELECT project_id FROM databases WHERE id = $1`,
		"deployments":       `SELECT project_id FROM deployments WHERE id = $1`,
		"domains":           `SELECT COALESCE(d.project_id,a.project_id) FROM domains d LEFT JOIN app_services a ON a.id=d.service_id WHERE d.id=$1`,
		"volumes":           `SELECT a.project_id FROM service_volumes v JOIN app_services a ON a.id=v.service_id WHERE v.id=$1`,
		"route-rules":       `SELECT a.project_id FROM route_rules r JOIN app_services a ON a.id=r.service_id WHERE r.id=$1`,
		"scheduled-tasks":   `SELECT a.project_id FROM scheduled_tasks t JOIN app_services a ON a.id=t.service_id WHERE t.id=$1`,
		"log-drains":        `SELECT a.project_id FROM log_drains l JOIN app_services a ON a.id=l.service_id WHERE l.id=$1`,
		"webhooks":          `SELECT a.project_id FROM service_webhooks w JOIN app_services a ON a.id=w.service_id WHERE w.id=$1`,
		"registries":        `SELECT project_id FROM registries WHERE id=$1 AND project_id IS NOT NULL`,
		"variables":         `SELECT a.project_id FROM service_vars v JOIN app_services a ON a.id = v.service_id WHERE v.id = $1`,
		"backups":           `SELECT COALESCE(a.project_id,d.project_id) FROM backup_configs b LEFT JOIN app_services a ON a.id=b.service_id LEFT JOIN databases d ON d.id=b.database_id WHERE b.id = $1`,
		"backup-records":    `SELECT COALESCE(a.project_id,d.project_id) FROM backup_records r JOIN backup_configs b ON b.id=r.backup_config_id LEFT JOIN app_services a ON a.id=b.service_id LEFT JOIN databases d ON d.id=b.database_id WHERE r.id = $1`,
	}
	query, ok := queries[kind]
	if !ok {
		return "", errors.New("resource is not project scoped")
	}
	var project string
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&project); err != nil {
		return "", fmt.Errorf("resolve token project: %w", err)
	}
	return project, nil
}
