package system

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"codedock/internal/models"
)

type cfListResponse[T any] struct {
	Success bool  `json:"success"`
	Errors  []any `json:"errors"`
	Result  []T   `json:"result"`
}

type cfSingleResponse[T any] struct {
	Success bool  `json:"success"`
	Errors  []any `json:"errors"`
	Result  T     `json:"result"`
}

func (s *CloudflareIngressService) ListAccounts(ctx context.Context, token string) ([]models.IngressDiscoveryAccount, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.baseURL+"/accounts", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("cloudflare accounts returned status %d", resp.StatusCode)
	}

	var data cfListResponse[models.IngressDiscoveryAccount]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Result, nil
}

func (s *CloudflareIngressService) ListTunnels(ctx context.Context, token, accountID string) ([]models.IngressDiscoveryTunnel, error) {
	path := fmt.Sprintf("%s/accounts/%s/cfd_tunnel?is_deleted=false", s.baseURL, url.PathEscape(accountID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("cloudflare tunnels returned status %d", resp.StatusCode)
	}

	var raw cfListResponse[map[string]any]
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	result := make([]models.IngressDiscoveryTunnel, 0, len(raw.Result))
	for _, item := range raw.Result {
		id, _ := item["id"].(string)
		name, _ := item["name"].(string)
		status, _ := item["status"].(string)
		remote := false
		if remoteCfg, ok := item["remote_config"].(bool); ok && remoteCfg {
			remote = true
		} else if cfgSrc, ok := item["config_src"].(string); ok && cfgSrc == "cloudflare" {
			remote = true
		}
		result = append(result, models.IngressDiscoveryTunnel{
			ID:     id,
			Name:   name,
			Status: status,
			Remote: remote,
		})
	}
	return result, nil
}

func (s *CloudflareIngressService) ListAccessPolicies(ctx context.Context, token, accountID string) ([]models.IngressDiscoveryPolicy, error) {
	path := fmt.Sprintf("%s/accounts/%s/access/policies", s.baseURL, url.PathEscape(accountID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("cloudflare policies returned status %d", resp.StatusCode)
	}

	var data cfListResponse[models.IngressDiscoveryPolicy]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.Result, nil
}

func (s *CloudflareIngressService) CreateTunnel(ctx context.Context, token, accountID, name string) (string, error) {
	path := fmt.Sprintf("%s/accounts/%s/cfd_tunnel", s.baseURL, url.PathEscape(accountID))
	payload := map[string]string{
		"name":       name,
		"config_src": "cloudflare",
	}
	b, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("cloudflare tunnel create failed with status %d", resp.StatusCode)
	}

	var data cfSingleResponse[map[string]any]
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	id, _ := data.Result["id"].(string)
	return id, nil
}
