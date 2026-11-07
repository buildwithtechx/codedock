package client

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock/internal/models"
)

func (c *Client) ListProjects() ([]*models.ProjectConfig, error) {
	resp, err := c.sendRequest("GET", "/projects", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list projects (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.ProjectConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) CreateProject(req *models.CreateProjectRequest) (*models.ProjectConfig, error) {
	resp, err := c.sendRequest("POST", "/projects", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create project (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.ProjectConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) GetProject(id string) (*models.ProjectConfig, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/projects/%s", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get project (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.ProjectConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) DeleteProject(id string) error {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("/projects/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete project (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) ListRouteProtection(projectID string) ([]*models.RouteProtectionStatus, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/projects/%s/route-protection", projectID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list route protection (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.RouteProtectionStatus `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) ChangeRouteProtection(projectID string, input *models.RouteProtectionInput) (*models.RouteProtectionStatus, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/projects/%s/route-protection", projectID), input)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update route protection (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.RouteProtectionStatus `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
