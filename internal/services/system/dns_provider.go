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

func (s *DNSProviderService) ProvisionARecord(ctx context.Context, domain string) error {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil || cfg == nil {
		return err
	}
	targetIP := cfg.PublicIPv4
	if targetIP == "" {
		return fmt.Errorf("PublicIPv4 is not set in server settings")
	}
	return s.ProvisionRecord(ctx, domain, "A", targetIP)
}

func (s *DNSProviderService) DeprovisionARecord(ctx context.Context, domain string) error {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil || cfg == nil {
		return err
	}
	targetIP := cfg.PublicIPv4
	return s.DeprovisionRecord(ctx, domain, "A", targetIP)
}

func (s *DNSProviderService) ProvisionRecord(ctx context.Context, domain, recordType, value string) error {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil || cfg == nil {
		return err
	}

	if cfg.CloudflareAPIToken != "" {
		if err := s.provisionCloudflare(ctx, cfg.CloudflareAPIToken, domain, recordType, value); err == nil {
			return nil
		}
	}

	if cfg.NamecheapAPIKey != "" && cfg.NamecheapAPIUser != "" {
		if err := s.provisionNamecheap(ctx, cfg, domain, recordType, value); err == nil {
			return nil
		}
	}

	if cfg.SpaceshipAPIKey != "" {
		if err := s.provisionSpaceship(ctx, cfg.SpaceshipAPIKey, domain, recordType, value); err == nil {
			return nil
		}
	}

	return nil
}

func (s *DNSProviderService) DeprovisionRecord(ctx context.Context, domain, recordType, value string) error {
	cfg, err := s.settingsRepo.GetServerSettings(ctx)
	if err != nil || cfg == nil {
		return err
	}

	if cfg.CloudflareAPIToken != "" {
		_ = s.deprovisionCloudflare(ctx, cfg.CloudflareAPIToken, domain, recordType, value)
	}

	if cfg.NamecheapAPIKey != "" && cfg.NamecheapAPIUser != "" {
		_ = s.deprovisionNamecheap(ctx, cfg, domain, recordType, value)
	}

	if cfg.SpaceshipAPIKey != "" {
		_ = s.deprovisionSpaceship(ctx, cfg.SpaceshipAPIKey, domain, recordType, value)
	}

	return nil
}

func getRootDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return domain
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}

func (s *DNSProviderService) provisionCloudflare(ctx context.Context, token, domain, recordType, value string) error {
	rootDomain := getRootDomain(domain)
	client := &http.Client{}

	zoneID, err := s.getCloudflareZoneID(ctx, client, token, rootDomain)
	if err != nil {
		return err
	}

	exists, _, err := s.checkCloudflareRecordExists(ctx, client, token, zoneID, domain, recordType, value)
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
	client := &http.Client{}

	zoneID, err := s.getCloudflareZoneID(ctx, client, token, rootDomain)
	if err != nil {
		return err
	}

	_, recordIDs, err := s.checkCloudflareRecordExists(ctx, client, token, zoneID, domain, recordType, value)
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

func (s *DNSProviderService) checkCloudflareRecordExists(ctx context.Context, client *http.Client, token, zoneID, domain, recordType, value string) (bool, []string, error) {
	u := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records?name=%s&type=%s", zoneID, domain, recordType)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return false, nil, err
	}
	defer resp.Body.Close()

	var res struct {
		Result []struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Type    string `json:"type"`
			Content string `json:"content"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, nil, err
	}

	var matchingIDs []string
	exists := false
	for _, rec := range res.Result {
		matchingIDs = append(matchingIDs, rec.ID)
		if value == "" || rec.Content == value {
			exists = true
		}
	}
	return exists, matchingIDs, nil
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

	client := &http.Client{}
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
		bodyStr := string(body)
		if strings.Contains(bodyStr, fmt.Sprintf(`Name="%s"`, subDomain)) && strings.Contains(bodyStr, fmt.Sprintf(`Type="%s"`, recordType)) {
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

	client := &http.Client{}
	u := "https://api.namecheap.com/xml.response"
	q := url.Values{}
	q.Set("ApiUser", cfg.NamecheapAPIUser)
	q.Set("ApiKey", cfg.NamecheapAPIKey)
	q.Set("UserName", cfg.NamecheapAPIUser)
	q.Set("Command", "namecheap.domains.dns.setHosts")
	q.Set("ClientIp", cfg.NamecheapClientIP)
	q.Set("SLD", parts[0])
	q.Set("TLD", parts[1])

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+q.Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *DNSProviderService) provisionSpaceship(ctx context.Context, key, domain, recordType, value string) error {
	client := &http.Client{}
	reqGet, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://spaceship.dev/api/v1/dns", nil)
	reqGet.Header.Set("X-Api-Key", key)
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
					return nil
				}
			}
		}
		respGet.Body.Close()
	}

	payload := map[string]any{
		"type":    recordType,
		"name":    domain,
		"address": value,
	}
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://spaceship.dev/api/v1/dns", bytes.NewBuffer(b))
	req.Header.Set("X-Api-Key", key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (s *DNSProviderService) deprovisionSpaceship(ctx context.Context, key, domain, recordType, value string) error {
	client := &http.Client{}
	reqGet, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://spaceship.dev/api/v1/dns", nil)
	reqGet.Header.Set("X-Api-Key", key)
	respGet, err := client.Do(reqGet)
	if err != nil {
		return err
	}

	var records []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Address string `json:"address"`
	}
	_ = json.NewDecoder(respGet.Body).Decode(&records)
	respGet.Body.Close()

	for _, r := range records {
		if r.Name == domain && r.Type == recordType && (value == "" || r.Address == value) {
			targetURL := "https://spaceship.dev/api/v1/dns"
			if r.ID != "" {
				targetURL = fmt.Sprintf("%s/%s", targetURL, r.ID)
			}
			reqDel, _ := http.NewRequestWithContext(ctx, http.MethodDelete, targetURL, nil)
			reqDel.Header.Set("X-Api-Key", key)
			if respDel, err := client.Do(reqDel); err == nil {
				respDel.Body.Close()
			}
		}
	}
	return nil
}
