package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock/internal/models"
	"codedock/internal/utils"
)

type ProjectAppRepository interface {
	ListByOrganization(ctx context.Context, organizationID string) ([]*models.ProjectApp, error)
	GetByID(ctx context.Context, id string) (*models.ProjectApp, error)
	GetBySlug(ctx context.Context, organizationID, slug string) (*models.ProjectApp, error)
	Create(ctx context.Context, app *models.ProjectApp) error
	Update(ctx context.Context, app *models.ProjectApp) error
	Delete(ctx context.Context, id string) error
}

type ProjectAppRepo struct {
	db *sqlx.DB
}

func NewProjectAppRepo(db *sql.DB) *ProjectAppRepo {
	return &ProjectAppRepo{db: sqlx.NewDb(db, "pgx")}
}

func (r *ProjectAppRepo) ListByOrganization(ctx context.Context, organizationID string) ([]*models.ProjectApp, error) {
	var list []*models.ProjectApp
	query := `
		SELECT id, organization_id, name, slug, git_provider, git_owner, git_repo, git_url,
		       installation_id, favicon, favicon_checked_at, deleted_at, created_at, updated_at
		FROM project_apps
		WHERE organization_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
	`
	err := r.db.SelectContext(ctx, &list, query, organizationID)
	if err != nil {
		return nil, fmt.Errorf("failed to list project apps: %w", err)
	}
	if list == nil {
		list = make([]*models.ProjectApp, 0)
	}
	return list, nil
}

func (r *ProjectAppRepo) GetByID(ctx context.Context, id string) (*models.ProjectApp, error) {
	var app models.ProjectApp
	query := `
		SELECT id, organization_id, name, slug, git_provider, git_owner, git_repo, git_url,
		       installation_id, favicon, favicon_checked_at, deleted_at, created_at, updated_at
		FROM project_apps
		WHERE id = $1 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &app, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("ProjectApp", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project app: %w", err)
	}
	return &app, nil
}

func (r *ProjectAppRepo) GetBySlug(ctx context.Context, organizationID, slug string) (*models.ProjectApp, error) {
	var app models.ProjectApp
	query := `
		SELECT id, organization_id, name, slug, git_provider, git_owner, git_repo, git_url,
		       installation_id, favicon, favicon_checked_at, deleted_at, created_at, updated_at
		FROM project_apps
		WHERE organization_id = $1 AND slug = $2 AND deleted_at IS NULL
	`
	err := r.db.GetContext(ctx, &app, query, organizationID, slug)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("ProjectApp", slug)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project app by slug: %w", err)
	}
	return &app, nil
}

func (r *ProjectAppRepo) Create(ctx context.Context, app *models.ProjectApp) error {
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	app.CreatedAt = now
	app.UpdatedAt = now
	if app.GitProvider == "" {
		app.GitProvider = "github"
	}

	query := `
		INSERT INTO project_apps (
			id, organization_id, name, slug, git_provider, git_owner, git_repo, git_url,
			installation_id, favicon, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		app.ID, app.OrganizationID, app.Name, app.Slug, app.GitProvider, app.GitOwner, app.GitRepo, app.GitURL,
		app.InstallationID, app.Favicon, app.CreatedAt, app.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create project app: %w", err)
	}
	return nil
}

func (r *ProjectAppRepo) Update(ctx context.Context, app *models.ProjectApp) error {
	app.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE project_apps SET
			name = $1, slug = $2, favicon = $3, git_provider = $4, git_owner = $5, git_repo = $6, git_url = $7,
			installation_id = $8, updated_at = $9
		WHERE id = $10
	`
	_, err := r.db.ExecContext(ctx, query,
		app.Name, app.Slug, app.Favicon, app.GitProvider, app.GitOwner, app.GitRepo, app.GitURL,
		app.InstallationID, app.UpdatedAt, app.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update project app: %w", err)
	}
	return nil
}

func (r *ProjectAppRepo) Delete(ctx context.Context, id string) error {
	now := time.Now().UTC()
	query := `UPDATE project_apps SET deleted_at = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, now, now, id)
	if err != nil {
		return fmt.Errorf("failed to delete project app: %w", err)
	}
	return nil
}
