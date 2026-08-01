package system

import (
	"context"

	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/system/dns_providers"
)

type DNSProviderService struct {
	providers *dnsproviders.Service
}

func NewDNSProviderService(repo repositories.SettingsRepository) *DNSProviderService {
	return &DNSProviderService{providers: dnsproviders.New(repo)}
}

func (s *DNSProviderService) ProvisionARecord(ctx context.Context, domain string) (string, string, error) {
	return s.providers.ProvisionARecord(ctx, domain)
}

func (s *DNSProviderService) DeprovisionARecord(ctx context.Context, domain, provider, value string) error {
	return s.providers.DeprovisionARecord(ctx, domain, provider, value)
}

func (s *DNSProviderService) ProvisionRecord(ctx context.Context, domain, recordType, value string) (string, error) {
	return s.providers.ProvisionRecord(ctx, domain, recordType, value)
}

func (s *DNSProviderService) DeprovisionRecord(ctx context.Context, domain, recordType, value, provider string) error {
	return s.providers.DeprovisionRecord(ctx, domain, recordType, value, provider)
}
