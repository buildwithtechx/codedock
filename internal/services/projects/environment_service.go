package projects

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"golang.org/x/net/publicsuffix"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type domainRepository interface {
	ListByService(ctx context.Context, serviceID string) ([]models.DomainConfig, error)
	ListAll(ctx context.Context) ([]models.DomainConfig, error)
	GetByID(ctx context.Context, id string) (*models.DomainConfig, error)
	Create(ctx context.Context, d *models.DomainConfig) error
	UpdateDNSProvisionStatus(ctx context.Context, id string, status string, provider string, provisionedIP string) error
	Delete(ctx context.Context, id string) error
}

type domainDNSService interface {
	ProvisionARecord(ctx context.Context, domain string) (string, string, error)
	DeprovisionARecord(ctx context.Context, domain, provider, value string) error
}

func canonicalDomainName(domain string) (string, error) {
	canonical := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	if canonical == "" || strings.ContainsAny(canonical, " /\\") {
		return "", fmt.Errorf("invalid domain name")
	}
	if len(canonical) > 253 {
		return "", fmt.Errorf("invalid domain name: exceeds 253 characters")
	}
	for _, label := range strings.Split(canonical, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return "", fmt.Errorf("invalid domain name: invalid label")
		}
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return "", fmt.Errorf("invalid domain name: invalid label")
			}
		}
	}
	if _, err := publicsuffix.EffectiveTLDPlusOne(canonical); err != nil {
		return "", fmt.Errorf("invalid domain name: %w", err)
	}
	return canonical, nil
}

type domainLockEntry struct {
	locked chan struct{}
	refs   atomic.Int32
}

type EnvironmentService struct {
	envRepo       repositories.EnvironmentRepository
	domainRepo    domainRepository
	varRepo       repositories.EnvRepository
	dnsService    domainDNSService
	domainLocks   sync.Map
	domainLocksMu sync.Mutex
}

func NewEnvironmentService(er repositories.EnvironmentRepository, dr domainRepository, vr repositories.EnvRepository, dnsService domainDNSService) *EnvironmentService {
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
	canonicalDomain, err := canonicalDomainName(d.DomainName)
	if err != nil {
		return nil, err
	}
	d.DomainName = canonicalDomain
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	now := time.Now()
	if d.CreatedAt.IsZero() {
		d.CreatedAt = now
	}
	d.DNSProvisionStatus = models.DNSProvisionStatusPending
	d.DNSProvider = ""
	d.DNSProvisionedIP = ""
	if err := s.domainRepo.Create(ctx, d); err != nil {
		return nil, err
	}

	if s.dnsService != nil {
		domainID := d.ID
		domainName := d.DomainName
		go func(domainID, domainName string) {
			provisionCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			unlock, err := s.lockDomain(provisionCtx, domainName)
			if err != nil {
				statusCtx, statusCancel := context.WithTimeout(context.Background(), 5*time.Second)
				statusErr := s.updateDNSProvisionStatus(statusCtx, domainID, models.DNSProvisionStatusFailed, "", "")
				statusCancel()
				if statusErr != nil {
					slog.Error("failed to mark DNS provisioning lock timeout", "domain_id", domainID, "error", statusErr)
				}
				return
			}
			defer unlock()

			var current *models.DomainConfig
			var readErr error
			for attempt := 0; attempt < 3; attempt++ {
				current, readErr = s.domainRepo.GetByID(provisionCtx, domainID)
				if readErr == nil {
					break
				}
				if !waitForDomainRetry(provisionCtx, attempt) {
					break
				}
			}
			if readErr != nil {
				statusCtx, statusCancel := context.WithTimeout(context.Background(), 5*time.Second)
				statusErr := s.updateDNSProvisionStatus(statusCtx, domainID, models.DNSProvisionStatusFailed, "", "")
				statusCancel()
				if statusErr != nil {
					slog.Error("failed to mark DNS provisioning read failure", "domain_id", domainID, "error", statusErr)
				}
				return
			}
			if current == nil || current.DomainName != domainName {
				return
			}

			provider, provisionedIP, err := s.dnsService.ProvisionARecord(provisionCtx, domainName)
			status := models.DNSProvisionStatusSuccess
			if err != nil {
				status = models.DNSProvisionStatusFailed
				provider = ""
				provisionedIP = ""
			}
			statusCtx, statusCancel := context.WithTimeout(context.Background(), 5*time.Second)
			statusErr := s.updateDNSProvisionStatus(statusCtx, domainID, status, provider, provisionedIP)
			statusCancel()
			if statusErr != nil {
				retryCtx, retryCancel := context.WithTimeout(context.Background(), 30*time.Second)
				retryErr := s.updateDNSProvisionStatus(retryCtx, domainID, status, provider, provisionedIP)
				retryCancel()
				if retryErr != nil {
					slog.Error("failed to persist DNS provisioning status", "domain_id", domainID, "error", retryErr)
					cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
					cleanupErr := s.dnsService.DeprovisionARecord(cleanupCtx, domainName, provider, provisionedIP)
					cleanupCancel()
					if cleanupErr != nil {
						slog.Error("failed to compensate DNS provisioning", "domain_id", domainID, "error", cleanupErr)
					}
				}
			}
		}(domainID, domainName)
	}

	return d, nil
}

func waitForDomainRetry(ctx context.Context, attempt int) bool {
	timer := time.NewTimer(time.Duration(attempt+1) * 200 * time.Millisecond)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
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

	unlock, err := s.lockDomain(ctx, domain.DomainName)
	if err != nil {
		return err
	}
	defer unlock()
	domain, err = s.domainRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if domain == nil {
		return nil
	}

	if s.dnsService != nil && domain.DNSProvisionStatus == models.DNSProvisionStatusSuccess {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 30*time.Second)
		err := s.dnsService.DeprovisionARecord(cleanupCtx, domain.DomainName, domain.DNSProvider, domain.DNSProvisionedIP)
		cleanupCancel()
		if err != nil {
			return fmt.Errorf("deprovision domain DNS record: %w", err)
		}
	}

	deleteCtx, deleteCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer deleteCancel()
	if err := s.domainRepo.Delete(deleteCtx, id); err != nil {
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

func (s *EnvironmentService) lockDomain(ctx context.Context, domainName string) (func(), error) {
	if domainName == "" {
		return func() {}, nil
	}

	s.domainLocksMu.Lock()
	lock, loaded := s.domainLocks.Load(domainName)
	if !loaded {
		lock = &domainLockEntry{locked: make(chan struct{}, 1)}
		s.domainLocks.Store(domainName, lock)
	}
	entry := lock.(*domainLockEntry)
	entry.refs.Add(1)
	s.domainLocksMu.Unlock()

	select {
	case entry.locked <- struct{}{}:
	case <-ctx.Done():
		s.releaseDomainLock(domainName, entry)
		return nil, ctx.Err()
	}

	return func() {
		<-entry.locked
		s.releaseDomainLock(domainName, entry)
	}, nil
}

func (s *EnvironmentService) releaseDomainLock(domainName string, entry *domainLockEntry) {
	s.domainLocksMu.Lock()
	defer s.domainLocksMu.Unlock()
	if entry.refs.Add(-1) == 0 {
		s.domainLocks.CompareAndDelete(domainName, entry)
	}
}

func (s *EnvironmentService) updateDNSProvisionStatus(ctx context.Context, domainID, status, provider, provisionedIP string) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = s.domainRepo.UpdateDNSProvisionStatus(ctx, domainID, status, provider, provisionedIP)
		if err == nil {
			return nil
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if attempt == 2 {
			break
		}
		timer := time.NewTimer(time.Duration(attempt+1) * 200 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return err
}
