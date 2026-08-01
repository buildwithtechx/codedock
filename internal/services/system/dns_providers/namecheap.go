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

func (s *Service) provisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) error {
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

func (s *Service) deprovisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) error {
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
