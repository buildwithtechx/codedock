package dnsproviders

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"codedock.run/codedock/internal/models"
)

func fetchNamecheapHosts(ctx context.Context, client *http.Client, cfg *models.ServerSettings, sld, tld string) ([]byte, error) {
	u := "https://api.namecheap.com/xml.response"
	getQ := url.Values{}
	getQ.Set("ApiUser", cfg.NamecheapAPIUser)
	getQ.Set("ApiKey", cfg.NamecheapAPIKey)
	getQ.Set("UserName", cfg.NamecheapAPIUser)
	getQ.Set("Command", "namecheap.domains.dns.getHosts")
	getQ.Set("ClientIp", cfg.NamecheapClientIP)
	getQ.Set("SLD", sld)
	getQ.Set("TLD", tld)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+getQ.Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("namecheap returned status %d", resp.StatusCode)
	}
	return body, nil
}

func (s *Service) provisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) error {
	rootDomain := getRootDomain(domain)
	subDomain := strings.TrimSuffix(domain, "."+rootDomain)
	if subDomain == domain {
		subDomain = "@"
	}
	sld, tld, err := splitNamecheapDomain(domain)
	if err != nil {
		return err
	}

	client := newProviderHTTPClient()
	if body, err := fetchNamecheapHosts(ctx, client, cfg, sld, tld); err == nil && namecheapRecordExists(body, subDomain, recordType, value) {
		return nil
	}

	addQ := url.Values{}
	addQ.Set("ApiUser", cfg.NamecheapAPIUser)
	addQ.Set("ApiKey", cfg.NamecheapAPIKey)
	addQ.Set("UserName", cfg.NamecheapAPIUser)
	addQ.Set("Command", "namecheap.domains.dns.addHost")
	addQ.Set("ClientIp", cfg.NamecheapClientIP)
	addQ.Set("SLD", sld)
	addQ.Set("TLD", tld)
	addQ.Set("HostName1", subDomain)
	addQ.Set("RecordType1", recordType)
	addQ.Set("Address1", value)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.namecheap.com/xml.response?"+addQ.Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("namecheap returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *Service) deprovisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) error {
	rootDomain := getRootDomain(domain)
	subDomain := strings.TrimSuffix(domain, "."+rootDomain)
	if subDomain == domain {
		subDomain = "@"
	}
	sld, tld, err := splitNamecheapDomain(domain)
	if err != nil {
		return err
	}

	client := newProviderHTTPClient()
	body, err := fetchNamecheapHosts(ctx, client, cfg, sld, tld)
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

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.namecheap.com/xml.response?"+buildNamecheapSetHostsQuery([]string{sld, tld}, cfg, remaining).Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("namecheap returned status %d", resp.StatusCode)
	}
	return nil
}
