package dns

import (
	"context"

	dnsproviders "codedock/internal/modules/dns/providers"
)

type ProviderService struct {
	providers *dnsproviders.Service
}

func NewProviderService(reader dnsproviders.SettingsReader) *ProviderService {
	return &ProviderService{providers: dnsproviders.New(reader)}
}

func (s *ProviderService) ProvisionARecord(ctx context.Context, domain string) (string, string, error) {
	return s.providers.ProvisionARecord(ctx, domain)
}

func (s *ProviderService) DeprovisionARecord(ctx context.Context, domain, provider, value string) error {
	return s.providers.DeprovisionARecord(ctx, domain, provider, value)
}

func (s *ProviderService) ProvisionRecord(ctx context.Context, domain, recordType, value string) (string, error) {
	return s.providers.ProvisionRecord(ctx, domain, recordType, value)
}

func (s *ProviderService) DeprovisionRecord(ctx context.Context, domain, recordType, value, provider string) error {
	return s.providers.DeprovisionRecord(ctx, domain, recordType, value, provider)
}
