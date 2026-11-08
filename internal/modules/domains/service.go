package domains

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

	"codedock/internal/utils"
)

type DomainDNSService interface {
	ProvisionARecord(ctx context.Context, domain string) (string, string, error)
	DeprovisionARecord(ctx context.Context, domain, provider, value string) error
}

func CanonicalDomainName(domain string) (string, error) {
	canonical := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	if canonical == "" || strings.ContainsAny(canonical, " /\\") {
		return "", utils.NewValidationError("invalid domain name")
	}
	if len(canonical) > 253 {
		return "", utils.NewValidationError("invalid domain name: exceeds 253 characters")
	}
	for _, label := range strings.Split(canonical, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return "", utils.NewValidationError("invalid domain name: invalid label")
		}
		for _, char := range label {
			if (char < 'a' || char > 'z') && (char < '0' || char > '9') && char != '-' {
				return "", utils.NewValidationError("invalid domain name: invalid label")
			}
		}
	}
	if _, err := publicsuffix.EffectiveTLDPlusOne(canonical); err != nil {
		return "", utils.NewValidationError(fmt.Sprintf("invalid domain name: %v", err))
	}
	return canonical, nil
}

type domainLockEntry struct {
	locked chan struct{}
	refs   atomic.Int32
}

type domainProvisionOperation struct {
	cancel context.CancelFunc
}

type DomainService struct {
	domainRepo    DomainRepository
	dnsService    DomainDNSService
	domainLocks   sync.Map
	domainLocksMu sync.Mutex
	provisionOps  sync.Map
}

func NewDomainService(dr DomainRepository, dns DomainDNSService) *DomainService {
	return &DomainService{
		domainRepo: dr,
		dnsService: dns,
	}
}

func (s *DomainService) ListDomainsByService(ctx context.Context, serviceID string) ([]DomainConfig, error) {
	if serviceID == "" {
		return nil, errors.New("service id required")
	}
	return s.domainRepo.ListByService(ctx, serviceID)
}

func (s *DomainService) ListAllDomains(ctx context.Context) ([]DomainConfig, error) {
	return s.domainRepo.ListAll(ctx)
}

func (s *DomainService) GetDomain(ctx context.Context, id string) (*DomainConfig, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.domainRepo.GetByID(ctx, id)
}

func (s *DomainService) CreateDomain(ctx context.Context, d *DomainConfig) (*DomainConfig, error) {
	return s.createDomain(ctx, d, true)
}

func (s *DomainService) CreateGeneratedDomain(ctx context.Context, d *DomainConfig) (*DomainConfig, error) {
	return s.createDomain(ctx, d, false)
}

func (s *DomainService) createDomain(ctx context.Context, d *DomainConfig, provisionDNS bool) (*DomainConfig, error) {
	if d == nil || d.ServiceID == "" || d.DomainName == "" {
		return nil, ErrInvalidDomain
	}
	canonicalDomain, err := CanonicalDomainName(d.DomainName)
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
	d.DNSProvisionStatus = DNSProvisionStatusPending
	d.DNSProvider = ""
	d.DNSProvisionedIP = ""
	if err := s.domainRepo.Create(ctx, d); err != nil {
		return nil, err
	}

	if provisionDNS && s.dnsService != nil {
		domainID := d.ID
		domainName := d.DomainName
		provisionCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		op := &domainProvisionOperation{cancel: cancel}
		s.provisionOps.Store(domainName, op)
		go s.runDNSProvisioning(domainID, domainName, provisionCtx, cancel, op)
	}

	return d, nil
}

func (s *DomainService) runDNSProvisioning(domainID, domainName string, ctx context.Context, cancel context.CancelFunc, op *domainProvisionOperation) {
	defer cancel()
	defer s.provisionOps.CompareAndDelete(domainName, op)

	unlock, err := s.lockDomain(ctx, domainName)
	if err != nil {
		s.safeUpdateStatus(domainID, DNSProvisionStatusFailed, "", "")
		return
	}
	defer unlock()

	current, readErr := s.domainRepo.GetByID(ctx, domainID)
	if readErr != nil || current == nil || current.DomainName != domainName {
		s.safeUpdateStatus(domainID, DNSProvisionStatusFailed, "", "")
		return
	}

	provider, provisionedIP, err := s.dnsService.ProvisionARecord(ctx, domainName)
	status := DNSProvisionStatusSuccess
	if err != nil {
		status = DNSProvisionStatusFailed
		provider = ""
		provisionedIP = ""
	}
	if err := s.safeUpdateStatus(domainID, status, provider, provisionedIP); err != nil {
		slog.Error("failed to persist DNS provisioning status", "domain_id", domainID, "error", err)
	}
}

func (s *DomainService) safeUpdateStatus(id, status, provider, ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.domainRepo.UpdateDNSProvisionStatus(ctx, id, status, provider, ip)
}

func (s *DomainService) DeleteDomain(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	deleteCtx, deleteCancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer deleteCancel()
	domain, err := s.domainRepo.GetByID(deleteCtx, id)
	if err != nil || domain == nil {
		return err
	}

	if op, ok := s.provisionOps.Load(domain.DomainName); ok {
		op.(*domainProvisionOperation).cancel()
	}

	unlock, err := s.lockDomain(deleteCtx, domain.DomainName)
	if err != nil {
		return err
	}
	defer unlock()

	if s.dnsService != nil && domain.DNSProvisionStatus == DNSProvisionStatusSuccess {
		_ = s.dnsService.DeprovisionARecord(deleteCtx, domain.DomainName, domain.DNSProvider, domain.DNSProvisionedIP)
	}

	dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dbCancel()
	return s.domainRepo.Delete(dbCtx, id)
}

func (s *DomainService) lockDomain(ctx context.Context, domainName string) (func(), error) {
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

func (s *DomainService) releaseDomainLock(domainName string, entry *domainLockEntry) {
	s.domainLocksMu.Lock()
	defer s.domainLocksMu.Unlock()
	if entry.refs.Add(-1) == 0 {
		s.domainLocks.CompareAndDelete(domainName, entry)
	}
}
