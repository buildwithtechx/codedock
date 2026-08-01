package projects

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/system"
)

type domainRepository interface {
	ListByService(ctx context.Context, serviceID string) ([]models.DomainConfig, error)
	ListAll(ctx context.Context) ([]models.DomainConfig, error)
	GetByID(ctx context.Context, id string) (*models.DomainConfig, error)
	Create(ctx context.Context, d *models.DomainConfig) error
	UpdateDNSProvisionStatus(ctx context.Context, id string, status string, provider string) error
	Delete(ctx context.Context, id string) error
}

type domainLockEntry struct {
	mutex sync.Mutex
	refs  atomic.Int32
}

type EnvironmentService struct {
	envRepo     repositories.EnvironmentRepository
	domainRepo  domainRepository
	varRepo     repositories.EnvRepository
	dnsService  *system.DNSProviderService
	domainLocks sync.Map
}

func NewEnvironmentService(er repositories.EnvironmentRepository, dr domainRepository, vr repositories.EnvRepository, dnsService *system.DNSProviderService) *EnvironmentService {
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
		domainID := d.ID
		domainName := d.DomainName
		go func(domainID, domainName string) {
			provisionCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			unlock := s.lockDomain(domainName)
			defer unlock()

			current, err := s.domainRepo.GetByID(provisionCtx, domainID)
			if err != nil || current == nil || current.DomainName != domainName {
				return
			}

			provider, err := s.dnsService.ProvisionARecord(provisionCtx, domainName)
			status := models.DNSProvisionStatusSuccess
			if err != nil {
				status = models.DNSProvisionStatusFailed
				provider = ""
			}
			_ = s.updateDNSProvisionStatus(provisionCtx, domainID, status, provider)
		}(domainID, domainName)
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
	domain, err := s.domainRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if domain == nil {
		return nil
	}

	unlock := s.lockDomain(domain.DomainName)
	defer unlock()

	if s.dnsService != nil {
		if err := s.dnsService.DeprovisionARecord(ctx, domain.DomainName, domain.DNSProvider); err != nil {
			return err
		}
	}

	if err := s.domainRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
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

func (s *EnvironmentService) lockDomain(domainName string) func() {
	if domainName == "" {
		return func() {}
	}

	lock, _ := s.domainLocks.LoadOrStore(domainName, &domainLockEntry{})
	entry := lock.(*domainLockEntry)
	entry.refs.Add(1)
	entry.mutex.Lock()
	return func() {
		entry.mutex.Unlock()
		if entry.refs.Add(-1) == 0 {
			s.domainLocks.CompareAndDelete(domainName, entry)
		}
	}
}

func (s *EnvironmentService) updateDNSProvisionStatus(ctx context.Context, domainID, status, provider string) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = s.domainRepo.UpdateDNSProvisionStatus(ctx, domainID, status, provider)
		if err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		time.Sleep(time.Duration(attempt+1) * 200 * time.Millisecond)
	}
	return err
}
