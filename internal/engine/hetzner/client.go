package hetzner

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	token   string
	http    *http.Client
	baseURL string
}

func NewClient(token string) *Client {
	return NewClientWithBaseURL(token, "https://api.hetzner.cloud/v1")
}

func NewClientWithBaseURL(token, baseURL string) *Client {
	return &Client{token: token, http: &http.Client{Timeout: 60 * time.Second}, baseURL: baseURL}
}

type APIError struct {
	Status  int
	Code    string
	Message string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("hetzner api failed: %d %s", e.Status, e.Message)
}

func (c *Client) call(ctx context.Context, method, path string, payload any) (map[string]any, error) {
	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(data)
	}
	request, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", "Bearer "+c.token)
	request.Header.Set("Content-Type", "application/json")
	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 8*1024*1024))
	if err != nil {
		return nil, err
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, fmt.Errorf("decode hetzner response: %w", err)
	}
	if response.StatusCode >= 400 {
		code := ""
		if details, ok := decoded["error"].(map[string]any); ok {
			code, _ = details["code"].(string)
		}
		return nil, &APIError{Status: response.StatusCode, Code: code, Message: string(data)}
	}
	return decoded, nil
}

func (c *Client) ListLocations(ctx context.Context) ([]map[string]any, error) {
	decoded, err := c.call(ctx, http.MethodGet, "/locations", nil)
	if err != nil {
		return nil, err
	}
	items, _ := decoded["locations"].([]any)
	return toMaps(items), nil
}

func (c *Client) ListServerTypes(ctx context.Context) ([]map[string]any, error) {
	decoded, err := c.call(ctx, http.MethodGet, "/server_types", nil)
	if err != nil {
		return nil, err
	}
	items, _ := decoded["server_types"].([]any)
	return toMaps(items), nil
}

func (c *Client) ListImages(ctx context.Context) ([]map[string]any, error) {
	decoded, err := c.call(ctx, http.MethodGet, "/images?type=system&architecture=x86", nil)
	if err != nil {
		return nil, err
	}
	items, _ := decoded["images"].([]any)
	return toMaps(items), nil
}

func (c *Client) FindServerByLabel(ctx context.Context, selector string) (map[string]any, error) {
	decoded, err := c.call(ctx, http.MethodGet, "/servers?label_selector="+selector, nil)
	if err != nil {
		return nil, err
	}
	items, _ := decoded["servers"].([]any)
	maps := toMaps(items)
	if len(maps) == 0 {
		return nil, nil
	}
	return maps[0], nil
}

func (c *Client) CreateServer(ctx context.Context, name, serverType, image, location, sshKey, userData string, labels map[string]string) (map[string]any, error) {
	decoded, err := c.call(ctx, http.MethodPost, "/servers", map[string]any{
		"name":        name,
		"server_type": serverType,
		"image":       image,
		"location":    location,
		"ssh_keys":    []string{sshKey},
		"user_data":   userData,
		"labels":      labels,
	})
	if err != nil {
		return nil, err
	}
	server, _ := decoded["server"].(map[string]any)
	return server, nil
}

func (c *Client) GetServer(ctx context.Context, id string) (map[string]any, error) {
	decoded, err := c.call(ctx, http.MethodGet, "/servers/"+id, nil)
	if err != nil {
		return nil, err
	}
	server, _ := decoded["server"].(map[string]any)
	return server, nil
}

func (c *Client) DeleteServer(ctx context.Context, id string) error {
	_, err := c.call(ctx, http.MethodDelete, "/servers/"+id, nil)
	return err
}

func (c *Client) PowerOn(ctx context.Context, id string) error {
	_, err := c.call(ctx, http.MethodPost, "/servers/"+id+"/actions/poweron", map[string]any{})
	return err
}

func (c *Client) ChangeType(ctx context.Context, id, serverType string, upgradeDisk bool) error {
	_, err := c.call(ctx, http.MethodPost, "/servers/"+id+"/actions/change_type", map[string]any{"server_type": serverType, "upgrade_disk": upgradeDisk})
	return err
}

func toMaps(items []any) []map[string]any {
	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		if object, ok := item.(map[string]any); ok {
			result = append(result, object)
		}
	}
	return result
}
