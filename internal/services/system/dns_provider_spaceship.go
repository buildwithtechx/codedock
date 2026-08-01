package system

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *DNSProviderService) provisionSpaceship(ctx context.Context, key, domain, recordType, value string) error {
	client := newProviderHTTPClient()
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
	client := newProviderHTTPClient()
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
