package dnsproviders

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type spaceshipRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Address string `json:"address"`
}

func fetchSpaceshipRecords(ctx context.Context, client *http.Client, key string) ([]spaceshipRecord, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://spaceship.dev/api/v1/dns", nil)
	req.Header.Set("X-Api-Key", key)
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

	records := make([]spaceshipRecord, 0)
	if err := json.NewDecoder(resp.Body).Decode(&records); err != nil {
		return nil, err
	}
	return records, nil
}

func spaceshipRecordExists(records []spaceshipRecord, domain, recordType, value string) bool {
	for _, record := range records {
		if record.Name == domain && record.Type == recordType && (value == "" || record.Address == value) {
			return true
		}
	}
	return false
}

func (s *Service) provisionSpaceship(ctx context.Context, key, domain, recordType, value string) error {
	if key == "" {
		return fmt.Errorf("spaceship api key not provided")
	}
	client := newProviderHTTPClient()
	records, err := fetchSpaceshipRecords(ctx, client, key)
	if err != nil {
		return err
	}
	if spaceshipRecordExists(records, domain, recordType, value) {
		return nil
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
	if resp.StatusCode >= 400 {
		return fmt.Errorf("spaceship returned status %d", resp.StatusCode)
	}
	return nil
}

func (s *Service) deprovisionSpaceship(ctx context.Context, key, domain, recordType, value string) error {
	if key == "" {
		return fmt.Errorf("spaceship api key not provided")
	}
	client := newProviderHTTPClient()
	records, err := fetchSpaceshipRecords(ctx, client, key)
	if err != nil {
		return err
	}

	for _, record := range records {
		if record.Name == domain && record.Type == recordType && (value == "" || record.Address == value) {
			targetURL := "https://spaceship.dev/api/v1/dns"
			if record.ID != "" {
				targetURL = fmt.Sprintf("%s/%s", targetURL, record.ID)
			}
			reqDel, _ := http.NewRequestWithContext(ctx, http.MethodDelete, targetURL, nil)
			reqDel.Header.Set("X-Api-Key", key)
			respDel, err := client.Do(reqDel)
			if err != nil {
				return err
			}
			respDel.Body.Close()
			if respDel.StatusCode >= 400 {
				return fmt.Errorf("spaceship returned status %d", respDel.StatusCode)
			}
		}
	}
	return nil
}
