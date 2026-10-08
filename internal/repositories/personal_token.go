package repositories

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *UserRepo) GetPATByHash(ctx context.Context, hash string) (*models.PersonalAccessToken, error) {
	var pat models.PersonalAccessToken
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
