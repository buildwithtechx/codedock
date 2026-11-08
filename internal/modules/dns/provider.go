package dns

import (
	"context"

	dnsproviders "codedock/internal/modules/dns/providers"
	"codedock/internal/repositories"
)

type settingsAdapter struct {
	repo repositories.SettingsRepository
}

func (a *settingsAdapter) GetDNSConfig(ctx context.Context) (*dnsproviders.DNSConfig, error) {
	if a.repo == nil {
		return nil, nil
	}
	cfg, err := a.repo.GetServerSettings(ctx)
	if err != nil || cfg == nil {
		return nil, err
	}
	return &dnsproviders.DNSConfig{
		PublicIPv4:         cfg.PublicIPv4,
		CloudflareAPIToken: cfg.CloudflareAPIToken,
		NamecheapAPIKey:    cfg.NamecheapAPIKey,
		NamecheapAPIUser:   cfg.NamecheapAPIUser,
		NamecheapClientIP:  cfg.NamecheapClientIP,
		SpaceshipAPIKey:    cfg.SpaceshipAPIKey,
		SpaceshipAPISecret: cfg.SpaceshipAPISecret,
	}, nil
}

type ProviderService struct {
	providers *dnsproviders.Service
}

func NewProviderService(repo repositories.SettingsRepository) *ProviderService {
	return &ProviderService{providers: dnsproviders.New(&settingsAdapter{repo: repo})}
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
