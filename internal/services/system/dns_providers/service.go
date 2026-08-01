package dnsproviders

import (
	"context"
	"fmt"

	"codedock.run/codedock/internal/repositories"
)

type Service struct {
	settingsRepo repositories.SettingsRepository
}

func New(repo repositories.SettingsRepository) *Service {
	return &Service{settingsRepo: repo}
}

func (s *Service) ProvisionARecord(ctx context.Context, domain string) (string, string, error) {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil {
		return "", "", err
	}
	if cfg == nil {
		return "", "", fmt.Errorf("server settings not configured")
	}
	targetIP := cfg.PublicIPv4
	if targetIP == "" {
		return "", "", fmt.Errorf("PublicIPv4 is not set in server settings")
	}
	provider, err := s.ProvisionRecord(ctx, domain, "A", targetIP)
	return provider, targetIP, err
}

func (s *Service) DeprovisionARecord(ctx context.Context, domain, provider, value string) error {
	if value == "" {
		return fmt.Errorf("provisioned A record value is missing")
	}
	return s.DeprovisionRecord(ctx, domain, "A", value, provider)
}

func (s *Service) ProvisionRecord(ctx context.Context, domain, recordType, value string) (string, error) {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil {
		return "", err
	}
	if cfg == nil {
		return "", fmt.Errorf("server settings not configured")
	}

	if cfg.CloudflareAPIToken != "" {
		if err := s.provisionCloudflare(ctx, cfg.CloudflareAPIToken, domain, recordType, value); err == nil {
			return dnsProviderCloudflare, nil
		}
	}

	if cfg.NamecheapAPIKey != "" && cfg.NamecheapAPIUser != "" {
		if err := s.provisionNamecheap(ctx, cfg, domain, recordType, value); err == nil {
			return dnsProviderNamecheap, nil
		}
	}

	if cfg.SpaceshipAPIKey != "" {
		if err := s.provisionSpaceship(ctx, cfg.SpaceshipAPIKey, domain, recordType, value); err == nil {
			return dnsProviderSpaceship, nil
		}
	}

	return "", fmt.Errorf("failed to provision dns record for %s", domain)
}

func (s *Service) DeprovisionRecord(ctx context.Context, domain, recordType, value, provider string) error {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("server settings not configured")
	}

	if provider == "" {
		provider = s.detectProvider(ctx, cfg, domain, recordType, value)
	}

	switch provider {
	case dnsProviderCloudflare:
		if cfg.CloudflareAPIToken == "" {
			return fmt.Errorf("cloudflare credentials missing during deprovision")
		}
		return s.deprovisionCloudflare(ctx, cfg.CloudflareAPIToken, domain, recordType, value)
	case dnsProviderNamecheap:
		if cfg.NamecheapAPIKey == "" || cfg.NamecheapAPIUser == "" {
			return fmt.Errorf("namecheap credentials missing during deprovision")
		}
		return s.deprovisionNamecheap(ctx, cfg, domain, recordType, value)
	case dnsProviderSpaceship:
		if cfg.SpaceshipAPIKey == "" {
			return fmt.Errorf("spaceship credentials missing during deprovision")
		}
		return s.deprovisionSpaceship(ctx, cfg.SpaceshipAPIKey, domain, recordType, value)
	}

	return fmt.Errorf("dns provider not found for %s", domain)
}
