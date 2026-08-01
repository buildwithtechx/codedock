package projects

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/system"
)

type EnvironmentService struct {
	envRepo    repositories.EnvironmentRepository
	domainRepo repositories.DomainRepository
	varRepo    repositories.EnvRepository
	dnsService *system.DNSProviderService
}

func NewEnvironmentService(er repositories.EnvironmentRepository, dr repositories.DomainRepository, vr repositories.EnvRepository, dnsService *system.DNSProviderService) *EnvironmentService {
	return &EnvironmentService{
		envRepo:    er,
		domainRepo: dr,
		varRepo:    vr,
		dnsService: dnsService,
	}
}

func (s *EnvironmentService) CreateEnvironment(ctx context.Context, env *models.EnvironmentConfig) (*models.EnvironmentConfig, error) {
	if env == nil || env.ProjectID == "" || env.Name == "" {
		return nil, errors.New("valid environment with projectId and name required")
	}
	if env.ID == "" {
		env.ID = uuid.New().String()
	}
	now := time.Now()
	if env.CreatedAt.IsZero() {
		env.CreatedAt = now
	}
	if err := s.envRepo.Create(ctx, env); err != nil {
		return nil, err
	}
	return env, nil
}

func (s *EnvironmentService) GetEnvironment(ctx context.Context, id string) (*models.EnvironmentConfig, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.envRepo.Get(ctx, id)
}

func (s *EnvironmentService) ListByProject(ctx context.Context, projectID string) ([]models.EnvironmentConfig, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	return s.envRepo.ListByProject(ctx, projectID)
}

func (s *EnvironmentService) DeleteEnvironment(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.envRepo.Delete(ctx, id)
}

func (s *EnvironmentService) CreateDomain(ctx context.Context, d *models.DomainConfig) (*models.DomainConfig, error) {
	if d == nil || d.ServiceID == "" || d.DomainName == "" {
		return nil, errors.New("valid domain with serviceId and domainName required")
	}
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	if d.DNSProvisionStatus == "" {
		d.DNSProvisionStatus = models.DNSProvisionStatusPending
	}
	if err := s.domainRepo.Create(ctx, d); err != nil {
		return nil, err
	}

	if s.dnsService != nil {
		go func() {
			err := s.dnsService.ProvisionARecord(context.Background(), d.DomainName)
			status := models.DNSProvisionStatusSuccess
			if err != nil {
				status = models.DNSProvisionStatusFailed
			}
			_ = s.domainRepo.UpdateDNSProvisionStatus(context.Background(), d.ID, status)
		}()
	}

	return d, nil
}

func (s *EnvironmentService) ListDomainsByService(ctx context.Context, serviceID string) ([]models.DomainConfig, error) {
	if serviceID == "" {
		return nil, errors.New("service id required")
	}
	return s.domainRepo.ListByService(ctx, serviceID)
}

func (s *EnvironmentService) ListAllDomains(ctx context.Context) ([]models.DomainConfig, error) {
	return s.domainRepo.ListAll(ctx)
}

func (s *EnvironmentService) GetDomain(ctx context.Context, id string) (*models.DomainConfig, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.domainRepo.GetByID(ctx, id)
}

func (s *EnvironmentService) DeleteDomain(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	domain, _ := s.domainRepo.GetByID(ctx, id)
	err := s.domainRepo.Delete(ctx, id)
	if err == nil && domain != nil && s.dnsService != nil {
		go func() {
			_ = s.dnsService.DeprovisionARecord(context.Background(), domain.DomainName)
		}()
	}
	return err
}

func (s *EnvironmentService) GetVars(ctx context.Context, projectID string) (map[string]string, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	return s.varRepo.GetVars(ctx, projectID)
}

func (s *EnvironmentService) SetVar(ctx context.Context, projectID, key, value string) error {
	if projectID == "" || key == "" {
		return errors.New("project id and key required")
	}
	return s.varRepo.SetVar(ctx, projectID, key, value)
}
