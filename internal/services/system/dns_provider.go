package system

import (
	"context"
	"fmt"

	"codedock.run/codedock/internal/repositories"
)

type DNSProviderService struct {
	settingsRepo repositories.SettingsRepository
}

func NewDNSProviderService(repo repositories.SettingsRepository) *DNSProviderService {
	return &DNSProviderService{settingsRepo: repo}
}

func (s *DNSProviderService) ProvisionARecord(ctx context.Context, domain string) (string, error) {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil {
		return "", err
	}
	if cfg == nil {
		return "", fmt.Errorf("server settings not configured")
	}
	targetIP := cfg.PublicIPv4
	if targetIP == "" {
		return "", fmt.Errorf("PublicIPv4 is not set in server settings")
	}
	return s.ProvisionRecord(ctx, domain, "A", targetIP)
}

func (s *DNSProviderService) DeprovisionARecord(ctx context.Context, domain, provider string) error {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("server settings not configured")
	}
	targetIP := cfg.PublicIPv4
	return s.DeprovisionRecord(ctx, domain, "A", targetIP, provider)
}

func (s *DNSProviderService) ProvisionRecord(ctx context.Context, domain, recordType, value string) (string, error) {
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

func (s *DNSProviderService) DeprovisionRecord(ctx context.Context, domain, recordType, value, provider string) error {
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
			return nil
		}
		return s.deprovisionCloudflare(ctx, cfg.CloudflareAPIToken, domain, recordType, value)
	case dnsProviderNamecheap:
		if cfg.NamecheapAPIKey == "" || cfg.NamecheapAPIUser == "" {
			return nil
		}
		return s.deprovisionNamecheap(ctx, cfg, domain, recordType, value)
	case dnsProviderSpaceship:
		if cfg.SpaceshipAPIKey == "" {
			return nil
		}
		return s.deprovisionSpaceship(ctx, cfg.SpaceshipAPIKey, domain, recordType, value)
	}

	return nil
}
