package system

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"codedock.run/codedock/internal/models"
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

func (s *DNSProviderService) provisionCloudflare(ctx context.Context, token, domain, recordType, value string) error {
	rootDomain := getRootDomain(domain)
	client := newProviderHTTPClient()

	zoneID, err := s.getCloudflareZoneID(ctx, client, token, rootDomain)
	if err != nil {
		return err
	}

	exists, _, err := checkCloudflareRecordExists(ctx, client, token, zoneID, domain, recordType, value)
	if err == nil && exists {
		return nil
	}

	payload := map[string]any{
		"type":    recordType,
		"name":    domain,
		"content": value,
		"proxied": false,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneID), bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("cloudflare returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *DNSProviderService) deprovisionCloudflare(ctx context.Context, token, domain, recordType, value string) error {
	rootDomain := getRootDomain(domain)
	client := newProviderHTTPClient()

	zoneID, err := s.getCloudflareZoneID(ctx, client, token, rootDomain)
	if err != nil {
		return err
	}

	_, recordIDs, err := checkCloudflareRecordExists(ctx, client, token, zoneID, domain, recordType, value)
	if err != nil {
		return err
	}

	for _, recID := range recordIDs {
		req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zoneID, recID), nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}
	return nil
}

func (s *DNSProviderService) getCloudflareZoneID(ctx context.Context, client *http.Client, token, rootDomain string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.cloudflare.com/client/v4/zones?name="+rootDomain, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var zoneRes struct {
		Result []struct {
			ID string `json:"id"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&zoneRes); err != nil || len(zoneRes.Result) == 0 {
		return "", fmt.Errorf("cloudflare zone not found for %s", rootDomain)
	}
	return zoneRes.Result[0].ID, nil
}

func (s *DNSProviderService) provisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) error {
	rootDomain := getRootDomain(domain)
	subDomain := strings.TrimSuffix(domain, "."+rootDomain)
	if subDomain == domain {
		subDomain = "@"
	}
	parts := strings.Split(rootDomain, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid root domain for namecheap")
	}

	client := newProviderHTTPClient()
	getU := "https://api.namecheap.com/xml.response"
	getQ := url.Values{}
	getQ.Set("ApiUser", cfg.NamecheapAPIUser)
	getQ.Set("ApiKey", cfg.NamecheapAPIKey)
	getQ.Set("UserName", cfg.NamecheapAPIUser)
	getQ.Set("Command", "namecheap.domains.dns.getHosts")
	getQ.Set("ClientIp", cfg.NamecheapClientIP)
	getQ.Set("SLD", parts[0])
	getQ.Set("TLD", parts[1])

	reqGet, _ := http.NewRequestWithContext(ctx, http.MethodGet, getU+"?"+getQ.Encode(), nil)
	if respGet, err := client.Do(reqGet); err == nil {
		body, _ := io.ReadAll(respGet.Body)
		respGet.Body.Close()
		if namecheapRecordExists(body, subDomain, recordType, value) {
			return nil
		}
	}

	addQ := url.Values{}
	addQ.Set("ApiUser", cfg.NamecheapAPIUser)
	addQ.Set("ApiKey", cfg.NamecheapAPIKey)
	addQ.Set("UserName", cfg.NamecheapAPIUser)
	addQ.Set("Command", "namecheap.domains.dns.addHost")
	addQ.Set("ClientIp", cfg.NamecheapClientIP)
	addQ.Set("SLD", parts[0])
	addQ.Set("TLD", parts[1])
	addQ.Set("HostName1", subDomain)
	addQ.Set("RecordType1", recordType)
	addQ.Set("Address1", value)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, getU+"?"+addQ.Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *DNSProviderService) deprovisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) error {
	rootDomain := getRootDomain(domain)
	subDomain := strings.TrimSuffix(domain, "."+rootDomain)
	if subDomain == domain {
		subDomain = "@"
	}
	parts := strings.Split(rootDomain, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid root domain for namecheap")
	}

	client := newProviderHTTPClient()
	u := "https://api.namecheap.com/xml.response"
	getQ := url.Values{}
	getQ.Set("ApiUser", cfg.NamecheapAPIUser)
	getQ.Set("ApiKey", cfg.NamecheapAPIKey)
	getQ.Set("UserName", cfg.NamecheapAPIUser)
	getQ.Set("Command", "namecheap.domains.dns.getHosts")
	getQ.Set("ClientIp", cfg.NamecheapClientIP)
	getQ.Set("SLD", parts[0])
	getQ.Set("TLD", parts[1])

	reqGet, _ := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+getQ.Encode(), nil)
	resp, err := client.Do(reqGet)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	hosts, err := parseNamecheapHosts(body)
	if err != nil {
		return err
	}

	remaining := make([]namecheapHost, 0, len(hosts))
	for _, host := range hosts {
		if host.Name == subDomain && host.Type == recordType && (value == "" || host.Address == value) {
			continue
		}
		remaining = append(remaining, host)
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+buildNamecheapSetHostsQuery(parts, cfg, remaining).Encode(), nil)
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("namecheap returned status %d", resp.StatusCode)
	}
	return nil
}
