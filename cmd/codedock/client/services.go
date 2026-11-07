package client

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock/internal/models"
)

func (c *Client) ListServices(environmentID string) ([]*models.AppService, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/environments/%s/services", environmentID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list services (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.AppService `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) CreateService(service *models.AppService) (*models.AppService, error) {
	resp, err := c.sendRequest("POST", "/services", service)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create service (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.AppService `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) GetService(id string) (*models.AppService, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/services/%s", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get service (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.AppService `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) DeleteService(id string) error {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("/services/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete service (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
