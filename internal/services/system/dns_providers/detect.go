package dnsproviders

import (
	"context"
	"strings"

	"codedock.run/codedock/internal/models"
)

func (s *Service) detectProvider(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) string {
	domain = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	if cfg.CloudflareAPIToken != "" {
		rootDomain := getRootDomain(domain)
		client := newProviderHTTPClient()
		zoneID, err := s.getCloudflareZoneID(ctx, client, cfg.CloudflareAPIToken, rootDomain)
		if err == nil {
			if exists, _, err := checkCloudflareRecordExists(ctx, client, cfg.CloudflareAPIToken, zoneID, domain, recordType, value); err == nil && exists {
				return dnsProviderCloudflare
			}
		}
	}

	if cfg.NamecheapAPIKey != "" && cfg.NamecheapAPIUser != "" {
		rootDomain := getRootDomain(domain)
		subDomain := strings.TrimSuffix(domain, "."+rootDomain)
		if subDomain == domain {
			subDomain = "@"
		}
		sld, tld, err := splitNamecheapDomain(domain)
		if err == nil {
			client := newProviderHTTPClient()
			body, err := fetchNamecheapHosts(ctx, client, cfg, sld, tld)
			if err == nil && namecheapRecordExists(body, subDomain, recordType, value) {
				return dnsProviderNamecheap
			}
		}
	}

	if cfg.SpaceshipAPIKey != "" {
		client := newProviderHTTPClient()
		if records, err := fetchSpaceshipRecords(ctx, client, cfg.SpaceshipAPIKey); err == nil && spaceshipRecordExists(records, domain, recordType, value) {
			return dnsProviderSpaceship
		}
	}

	return ""
}
