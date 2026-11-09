package apps

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"codedock/internal/modules/apps/vars"
)

type Service struct {
	repo       AppServiceRepository
	varRepo    ServiceVarRepository
	volumeRepo ServiceVolumeRepository
}

func NewService(
	repo AppServiceRepository,
	varRepo ServiceVarRepository,
	volumeRepo ServiceVolumeRepository,
) *Service {
	return &Service{
		repo:       repo,
		varRepo:    varRepo,
		volumeRepo: volumeRepo,
	}
}

func (s *Service) CreateAppService(ctx context.Context, svc *AppService) (*AppService, error) {
	if svc == nil {
		return nil, errors.New("app service is nil")
	}
	if svc.ProjectID == "" {
		return nil, errors.New("project id is required")
	}
	if svc.Name == "" {
		return nil, errors.New("name is required")
	}
	if svc.ID == "" {
		svc.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	svc.CreatedAt = now
	svc.UpdatedAt = now
	if svc.Status == "" {
		svc.Status = AppServiceStatusCreated
	}
	if err := s.repo.Create(ctx, svc); err != nil {
		return nil, err
	}
	return svc, nil
}

func (s *Service) GetAppService(ctx context.Context, id string) (*AppService, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	app, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if s.volumeRepo != nil {
		volumes, _ := s.volumeRepo.ListByService(ctx, app.ID)
		app.Volumes = volumes
	}

	return app, nil
}

func (s *Service) ListByEnvironment(ctx context.Context, environmentID string) ([]*AppService, error) {
	if environmentID == "" {
		return nil, errors.New("environment id required")
	}
	return s.repo.ListByEnvironment(ctx, environmentID)
}

func (s *Service) ListByProject(ctx context.Context, projectID string) ([]*AppService, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	apps, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if s.volumeRepo != nil {
		for _, app := range apps {
			volumes, _ := s.volumeRepo.ListByService(ctx, app.ID)
			app.Volumes = volumes
		}
	}
	return apps, nil
}

func (s *Service) ListByOrganization(ctx context.Context, organizationID string, search string) ([]*AppService, error) {
	if organizationID == "" {
		return nil, errors.New("organization id required")
	}
	return s.repo.ListByOrganization(ctx, organizationID, search)
}

func (s *Service) UpdateAppService(ctx context.Context, svc *AppService) error {
	if svc == nil {
		return errors.New("app service is nil")
	}
	if svc.ID == "" {
		return errors.New("valid service required for update")
	}
	svc.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, svc)
}

func (s *Service) DeleteAppService(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *Service) ListAllServices(ctx context.Context, search string) ([]*AppService, error) {
	return s.repo.ListAll(ctx, search)
}

func (s *Service) SearchIDsByName(ctx context.Context, term string, limit int) ([]string, error) {
	return s.repo.SearchIDsByName(ctx, term, limit)
}

func (s *Service) CreateVariable(ctx context.Context, v *Variable) (*Variable, error) {
	return vars.NewService(s.varRepo).CreateVariable(ctx, v)
}

func (s *Service) UpdateVariable(ctx context.Context, v *Variable) error {
	return vars.NewService(s.varRepo).UpdateVariable(ctx, v)
}

func (s *Service) GetVariable(ctx context.Context, id string) (*Variable, error) {
	return vars.NewService(s.varRepo).GetVariable(ctx, id)
}

func (s *Service) ListVariablesByService(ctx context.Context, serviceID string) ([]*Variable, error) {
	return vars.NewService(s.varRepo).ListVariablesByService(ctx, serviceID)
}

func (s *Service) DeleteVariable(ctx context.Context, id string) error {
	return vars.NewService(s.varRepo).DeleteVariable(ctx, id)
}
