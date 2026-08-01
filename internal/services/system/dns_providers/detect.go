package dnsproviders

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"codedock.run/codedock/internal/models"
)

func (s *Service) detectProvider(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) string {
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
		parts := strings.Split(rootDomain, ".")
		if len(parts) == 2 {
			client := newProviderHTTPClient()
			getQ := url.Values{}
			getQ.Set("ApiUser", cfg.NamecheapAPIUser)
			getQ.Set("ApiKey", cfg.NamecheapAPIKey)
			getQ.Set("UserName", cfg.NamecheapAPIUser)
			getQ.Set("Command", "namecheap.domains.dns.getHosts")
			getQ.Set("ClientIp", cfg.NamecheapClientIP)
			getQ.Set("SLD", parts[0])
			getQ.Set("TLD", parts[1])

			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.namecheap.com/xml.response?"+getQ.Encode(), nil)
			if resp, err := client.Do(req); err == nil {
				body, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				if namecheapRecordExists(body, subDomain, recordType, value) {
					return dnsProviderNamecheap
				}
			}
		}
	}

	if cfg.SpaceshipAPIKey != "" {
		client := newProviderHTTPClient()
		reqGet, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://spaceship.dev/api/v1/dns", nil)
		reqGet.Header.Set("X-Api-Key", cfg.SpaceshipAPIKey)
		if respGet, err := client.Do(reqGet); err == nil {
			var records []struct {
				Name    string `json:"name"`
				Type    string `json:"type"`
				Address string `json:"address"`
			}
			if json.NewDecoder(respGet.Body).Decode(&records) == nil {
				for _, r := range records {
					if r.Name == domain && r.Type == recordType && (value == "" || r.Address == value) {
						respGet.Body.Close()
						return dnsProviderSpaceship
					}
				}
			}
			respGet.Body.Close()
		}
	}

	return ""
}
