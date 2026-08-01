package dnsproviders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Service) provisionCloudflare(ctx context.Context, token, domain, recordType, value string) error {
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
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records", bytes.NewBuffer(b))
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

func (s *Service) deprovisionCloudflare(ctx context.Context, token, domain, recordType, value string) error {
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
		req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records/"+recID, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
		}
	}
	return nil
}

func (s *Service) getCloudflareZoneID(ctx context.Context, client *http.Client, token, rootDomain string) (string, error) {
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
