package projects

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProjectAppService struct {
	repo ProjectAppRepository
}

func NewProjectAppService(repo ProjectAppRepository) *ProjectAppService {
	return &ProjectAppService{repo: repo}
}

func (s *ProjectAppService) Create(ctx context.Context, app *ProjectApp) (*ProjectApp, error) {
	if app == nil {
		return nil, ErrProjectAppNil
	}
	if app.OrganizationID == "" {
		return nil, ErrOrganizationIDRequired
	}
	if app.Name == "" {
		return nil, ErrProjectAppNameRequired
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

func (s *ProjectAppService) GetByID(ctx context.Context, id string) (*ProjectApp, error) {
	if id == "" {
		return nil, ErrProjectAppIDRequired
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectAppService) GetBySlug(ctx context.Context, organizationID, slug string) (*ProjectApp, error) {
	if organizationID == "" || slug == "" {
		return nil, fmt.Errorf("organization id and slug are required")
	}
	return s.repo.GetBySlug(ctx, organizationID, slug)
}

func (s *ProjectAppService) ListByOrganization(ctx context.Context, organizationID string) ([]*ProjectApp, error) {
	if organizationID == "" {
		return nil, ErrOrganizationIDRequired
	}
	return s.repo.ListByOrganization(ctx, organizationID)
}

func (s *ProjectAppService) Update(ctx context.Context, app *ProjectApp) error {
	if app == nil || app.ID == "" {
		return ErrProjectAppIDRequired
	}
	app.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, app)
}

func (s *ProjectAppService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrProjectAppIDRequired
	}
	return s.repo.Delete(ctx, id)
}
