package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock.run/codedock/pkg/types"
)

func (c *Client) ListDatabases(projectID string) ([]*types.Database, error) {
	url := "/databases"
	if projectID != "" {
		url = fmt.Sprintf("/databases?projectId=%s", projectID)
	}
	resp, err := c.sendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list databases (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*types.Database `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) GetDatabase(id string) (*types.Database, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/databases/%s", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get database (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *types.Database `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) CreateDatabase(req *types.CreateDatabaseRequest) (*types.Database, error) {
	resp, err := c.sendRequest("POST", "/databases", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create database (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *types.Database `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) DeleteDatabase(id string) error {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("/databases/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete database (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) ImportDatabase(id string, req *types.ImportDatabaseRequest) error {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/databases/%s/import", id), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to import database (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
