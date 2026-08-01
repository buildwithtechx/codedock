package dnsproviders

import (
	"context"
	"strings"
	"time"

	"codedock.run/codedock/internal/models"
)

func (s *Service) detectProvider(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) string {
	domain = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	if cfg.CloudflareAPIToken != "" {
		probeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		rootDomain := getRootDomain(domain)
		client := s.httpClient()
		zoneID, err := s.getCloudflareZoneID(probeCtx, client, cfg.CloudflareAPIToken, rootDomain)
		if err == nil {
			if exists, _, err := checkCloudflareRecordExists(probeCtx, client, cfg.CloudflareAPIToken, zoneID, domain, recordType, value); err == nil && exists {
				cancel()
				return dnsProviderCloudflare
			}
		}
		cancel()
	}

	if cfg.NamecheapAPIKey != "" && cfg.NamecheapAPIUser != "" {
		rootDomain := getRootDomain(domain)
		subDomain := strings.TrimSuffix(domain, "."+rootDomain)
		if subDomain == domain {
			subDomain = "@"
		}
		sld, tld, err := splitNamecheapDomain(domain)
		if err == nil {
			probeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			client := s.httpClient()
			body, err := fetchNamecheapHosts(probeCtx, client, cfg, sld, tld)
			cancel()
			if err == nil && namecheapRecordExists(body, subDomain, recordType, value) {
				return dnsProviderNamecheap
			}
		}
	}

	if cfg.SpaceshipAPIKey != "" && cfg.SpaceshipAPISecret != "" {
		probeCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		client := s.httpClient()
		if records, err := fetchSpaceshipRecords(probeCtx, client, cfg.SpaceshipAPIKey, cfg.SpaceshipAPISecret, domain); err == nil && spaceshipRecordExists(records, domain, recordType, value) {
			cancel()
			return dnsProviderSpaceship
		}
		cancel()
	}

	return ""
}
