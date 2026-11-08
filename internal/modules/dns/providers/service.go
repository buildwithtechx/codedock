package dnsproviders

import (
	"context"
	"fmt"
	"net/http"
)

type Service struct {
	settingsReader SettingsReader
	newHTTPClient  func() *http.Client
}

func New(reader SettingsReader) *Service {
	return &Service{settingsReader: reader, newHTTPClient: newProviderHTTPClient}
}

func NewWithHTTPClient(reader SettingsReader, newHTTPClient func() *http.Client) *Service {
	return &Service{settingsReader: reader, newHTTPClient: newHTTPClient}
}

func (s *Service) httpClient() *http.Client {
	if s.newHTTPClient == nil {
		return newProviderHTTPClient()
	}
	return s.newHTTPClient()
}

func (s *Service) ProvisionARecord(ctx context.Context, domain string) (string, string, error) {
	if s.settingsReader == nil {
		return "", "", fmt.Errorf("settings reader not configured")
	}
	cfg, err := s.settingsReader.GetDNSConfig(ctx)
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
	if s.settingsReader == nil {
		return "", fmt.Errorf("settings reader not configured")
	}
	cfg, err := s.settingsReader.GetDNSConfig(ctx)
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

	if cfg.SpaceshipAPIKey != "" && cfg.SpaceshipAPISecret != "" {
		if err := s.provisionSpaceship(ctx, cfg.SpaceshipAPIKey, cfg.SpaceshipAPISecret, domain, recordType, value); err == nil {
			return dnsProviderSpaceship, nil
		}
	}

	return "", fmt.Errorf("no configured DNS provider could provision record for %s", domain)
}

func (s *Service) DeprovisionRecord(ctx context.Context, domain, recordType, value, provider string) error {
	if s.settingsReader == nil {
		return fmt.Errorf("settings reader not configured")
	}
	cfg, err := s.settingsReader.GetDNSConfig(ctx)
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
			return fmt.Errorf("cloudflare API token not set")
		}
		return s.deprovisionCloudflare(ctx, cfg.CloudflareAPIToken, domain, recordType, value)
	case dnsProviderNamecheap:
		if cfg.NamecheapAPIKey == "" || cfg.NamecheapAPIUser == "" {
			return fmt.Errorf("namecheap API credentials not set")
		}
		return s.deprovisionNamecheap(ctx, cfg, domain, recordType, value)
	case dnsProviderSpaceship:
		if cfg.SpaceshipAPIKey == "" {
			return fmt.Errorf("spaceship API key not set")
		}
		if cfg.SpaceshipAPISecret == "" {
			return fmt.Errorf("spaceship API secret not set")
		}
		return s.deprovisionSpaceship(ctx, cfg.SpaceshipAPIKey, cfg.SpaceshipAPISecret, domain, recordType, value)
	default:
		return fmt.Errorf("unsupported or undetectable DNS provider %q", provider)
	}
}
