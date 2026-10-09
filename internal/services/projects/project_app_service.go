package projects

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type ProjectAppService struct {
	repo repositories.ProjectAppRepository
}

func NewProjectAppService(repo repositories.ProjectAppRepository) *ProjectAppService {
	return &ProjectAppService{repo: repo}
}

func (s *ProjectAppService) Create(ctx context.Context, app *models.ProjectApp) (*models.ProjectApp, error) {
	if app == nil {
		return nil, errors.New("project app is nil")
	}
	if app.OrganizationID == "" {
		return nil, errors.New("organization id is required")
	}
	if app.Name == "" {
		return nil, errors.New("name is required")
	}
	if app.Slug == "" {
		app.Slug = strings.ToLower(strings.ReplaceAll(app.Name, " ", "-"))
	}
	if app.ID == "" {
		app.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	app.CreatedAt = now
	app.UpdatedAt = now

	if err := s.repo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("failed to create project app: %w", err)
	}
	return app, nil
}

func (s *ProjectAppService) GetByID(ctx context.Context, id string) (*models.ProjectApp, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectAppService) GetBySlug(ctx context.Context, organizationID, slug string) (*models.ProjectApp, error) {
	if organizationID == "" || slug == "" {
		return nil, errors.New("organization id and slug are required")
	}
	return s.repo.GetBySlug(ctx, organizationID, slug)
}

func (s *ProjectAppService) ListByOrganization(ctx context.Context, organizationID string) ([]*models.ProjectApp, error) {
	if organizationID == "" {
		return nil, errors.New("organization id is required")
	}
	return s.repo.ListByOrganization(ctx, organizationID)
}

func (s *ProjectAppService) Update(ctx context.Context, app *models.ProjectApp) error {
	if app == nil || app.ID == "" {
		return errors.New("project app id is required")
	}
	app.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, app)
}

func (s *ProjectAppService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.repo.Delete(ctx, id)
}
