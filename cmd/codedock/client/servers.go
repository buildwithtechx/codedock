package client

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock/internal/models"
)

func (c *Client) ListServers() ([]*models.Server, error) {
	resp, err := c.sendRequest("GET", "/servers", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list servers (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.Server `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) CreateServer(req *models.CreateServerRequest) (*models.Server, error) {
	resp, err := c.sendRequest("POST", "/servers", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create server (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.Server `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) DeleteServer(id string) error {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("/servers/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete server (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
