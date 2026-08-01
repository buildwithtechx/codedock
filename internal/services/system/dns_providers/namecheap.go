package dnsproviders

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"codedock.run/codedock/internal/models"
)

type namecheapApiResponse struct {
	Status string `xml:"Status,attr"`
	Errors struct {
		Errors []struct {
			Message string `xml:",chardata"`
		} `xml:"Error"`
	} `xml:"Errors"`
}

type namecheapLock struct {
	mutex sync.Mutex
	refs  int
}

func validateNamecheapResponse(body []byte) error {
	var res namecheapApiResponse
	if err := xml.Unmarshal(body, &res); err != nil {
		return err
	}
	status := strings.ToUpper(strings.TrimSpace(res.Status))
	if status == "OK" || status == "SUCCESS" {
		return nil
	}

	messages := make([]string, 0, len(res.Errors.Errors))
	for _, err := range res.Errors.Errors {
		msg := strings.TrimSpace(err.Message)
		if msg != "" {
			messages = append(messages, msg)
		}
	}
	if len(messages) == 0 {
		return fmt.Errorf("namecheap api error")
	}
	return fmt.Errorf("namecheap api error: %s", strings.Join(messages, "; "))
}

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
	if err := validateNamecheapResponse(body); err != nil {
		return nil, err
	}
	return body, nil
}

var namecheapLocks = struct {
	sync.Mutex
	entries map[string]*namecheapLock
}{entries: make(map[string]*namecheapLock)}

func lockNamecheap(domain string) func() {
	namecheapLocks.Lock()
	lock := namecheapLocks.entries[domain]
	if lock == nil {
		lock = &namecheapLock{}
		namecheapLocks.entries[domain] = lock
	}
	lock.refs++
	namecheapLocks.Unlock()

	lock.mutex.Lock()
	return func() {
		lock.mutex.Unlock()
		namecheapLocks.Lock()
		lock.refs--
		if lock.refs == 0 {
			delete(namecheapLocks.entries, domain)
		}
		namecheapLocks.Unlock()
	}
}

func (s *Service) provisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) error {
	rootDomain := getRootDomain(domain)
	unlock := lockNamecheap(rootDomain)
	defer unlock()

	normalizedDomain := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	subDomain := strings.TrimSuffix(normalizedDomain, "."+rootDomain)
	if subDomain == normalizedDomain {
		subDomain = "@"
	}
	sld, tld, err := splitNamecheapDomain(normalizedDomain)
	if err != nil {
		return err
	}

	client := s.httpClient()
	body, err := fetchNamecheapHosts(ctx, client, cfg, sld, tld)
	if err != nil {
		return err
	}
	hosts, err := parseNamecheapHosts(body)
	if err != nil {
		return err
	}
	for _, host := range hosts {
		if host.Name == subDomain && host.Type == recordType {
			return fmt.Errorf("namecheap record conflict for %s", domain)
		}
	}
	hosts = append(hosts, namecheapHost{Name: subDomain, Type: recordType, Address: value})
	setQ := buildNamecheapSetHostsQuery([]string{sld, tld}, cfg, hosts)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.namecheap.com/xml.response?"+setQ.Encode(), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("namecheap returned status %d", resp.StatusCode)
	}
	if err := validateNamecheapResponse(body); err != nil {
		return err
	}
	return nil
}

func (s *Service) deprovisionNamecheap(ctx context.Context, cfg *models.ServerSettings, domain, recordType, value string) error {
	normalizedDomain := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))
	rootDomain := getRootDomain(normalizedDomain)
	unlock := lockNamecheap(rootDomain)
	defer unlock()

	subDomain := strings.TrimSuffix(normalizedDomain, "."+rootDomain)
	if subDomain == normalizedDomain {
		subDomain = "@"
	}
	sld, tld, err := splitNamecheapDomain(normalizedDomain)
	if err != nil {
		return err
	}

	client := s.httpClient()
	body, err := fetchNamecheapHosts(ctx, client, cfg, sld, tld)
	if err != nil {
		return err
	}

	hosts, err := parseNamecheapHosts(body)
	if err != nil {
		return err
	}

	remaining := make([]namecheapHost, 0, len(hosts))
	matched := false
	for _, host := range hosts {
		if host.Name == subDomain && host.Type == recordType && (value == "" || host.Address == value) {
			matched = true
			continue
		}
		remaining = append(remaining, host)
	}
	if !matched {
		return nil
	}
	if len(remaining) == 0 {
		return s.setNamecheapDefaultDNS(ctx, client, cfg, sld, tld)
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.namecheap.com/xml.response?"+buildNamecheapSetHostsQuery([]string{sld, tld}, cfg, remaining).Encode(), nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("namecheap returned status %d", resp.StatusCode)
	}
	if err := validateNamecheapResponse(body); err != nil {
		return err
	}
	return nil
}

func (s *Service) setNamecheapDefaultDNS(ctx context.Context, client *http.Client, cfg *models.ServerSettings, sld, tld string) error {
	query := url.Values{}
	query.Set("ApiUser", cfg.NamecheapAPIUser)
	query.Set("ApiKey", cfg.NamecheapAPIKey)
	query.Set("UserName", cfg.NamecheapAPIUser)
	query.Set("Command", "namecheap.domains.dns.setDefault")
	query.Set("ClientIp", cfg.NamecheapClientIP)
	query.Set("SLD", sld)
	query.Set("TLD", tld)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.namecheap.com/xml.response?"+query.Encode(), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("namecheap returned status %d", resp.StatusCode)
	}
	return validateNamecheapResponse(body)
}
