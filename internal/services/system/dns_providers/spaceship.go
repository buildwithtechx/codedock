package dnsproviders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type spaceshipRecord struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
	TTL     int    `json:"ttl"`
}

func fetchSpaceshipRecords(ctx context.Context, client *http.Client, key, secret, domain string) ([]spaceshipRecord, error) {
	rootDomain := getRootDomain(domain)
	query := url.Values{"take": {"500"}, "skip": {"0"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://spaceship.dev/api/v1/dns/records/"+url.PathEscape(rootDomain)+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", key)
	req.Header.Set("X-API-Secret", secret)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			return nil, fmt.Errorf("spaceship returned status %d: %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("spaceship returned status %d", resp.StatusCode)
	}

	var response struct {
		Items []spaceshipRecord `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	for i := range response.Items {
		if response.Items[i].Name == "@" {
			response.Items[i].Name = rootDomain
		} else {
			response.Items[i].Name += "." + rootDomain
		}
	}
	return response.Items, nil
}

func spaceshipRecordExists(records []spaceshipRecord, domain, recordType, value string) bool {
	for _, record := range records {
		if record.Name == domain && record.Type == recordType && (value == "" || record.Address == value) {
			return true
		}
	}
	return false
}

func (s *Service) provisionSpaceship(ctx context.Context, key, secret, domain, recordType, value string) error {
	if key == "" || secret == "" {
		return fmt.Errorf("spaceship api key and secret are required")
	}
	client := newProviderHTTPClient()
	records, err := fetchSpaceshipRecords(ctx, client, key, secret, domain)
	if err != nil {
		return err
	}
	if spaceshipRecordExists(records, domain, recordType, value) {
		return nil
	}

	rootDomain := getRootDomain(domain)
	recordName := "@"
	if domain != rootDomain {
		recordName = domain[:len(domain)-len(rootDomain)-1]
	}
	items := make([]spaceshipRecord, 0, len(records)+1)
	for _, record := range records {
		record.Name = relativeSpaceshipName(record.Name, rootDomain)
		items = append(items, record)
	}
	items = append(items, spaceshipRecord{Type: recordType, Name: recordName, Address: value, TTL: 3600})
	payload := map[string]any{"force": false, "items": items}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "https://spaceship.dev/api/v1/dns/records/"+url.PathEscape(rootDomain), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", key)
	req.Header.Set("X-API-Secret", secret)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("spaceship returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *Service) deprovisionSpaceship(ctx context.Context, key, secret, domain, recordType, value string) error {
	if key == "" || secret == "" {
		return fmt.Errorf("spaceship api key and secret are required")
	}
	client := newProviderHTTPClient()
	records, err := fetchSpaceshipRecords(ctx, client, key, secret, domain)
	if err != nil {
		return err
	}

	deletions := make([]spaceshipRecord, 0)
	for _, record := range records {
		if record.Name == domain && record.Type == recordType && (value == "" || record.Address == value) {
			record.Name = relativeSpaceshipName(record.Name, getRootDomain(domain))
			deletions = append(deletions, record)
		}
	}
	if len(deletions) == 0 {
		return nil
	}
	b, _ := json.Marshal(deletions)
	rootDomain := getRootDomain(domain)
	reqDel, err := http.NewRequestWithContext(ctx, http.MethodDelete, "https://spaceship.dev/api/v1/dns/records/"+url.PathEscape(rootDomain), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	reqDel.Header.Set("X-API-Key", key)
	reqDel.Header.Set("X-API-Secret", secret)
	reqDel.Header.Set("Content-Type", "application/json")
	respDel, err := client.Do(reqDel)
	if err != nil {
		return err
	}
	defer respDel.Body.Close()
	if respDel.StatusCode >= 400 {
		return fmt.Errorf("spaceship returned status %d", respDel.StatusCode)
	}
	return nil
}

func relativeSpaceshipName(name, rootDomain string) string {
	if name == rootDomain {
		return "@"
	}
	return name[:len(name)-len(rootDomain)-1]
}
