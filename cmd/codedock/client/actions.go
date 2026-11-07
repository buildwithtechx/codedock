package client

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock/internal/models"
)

func (c *Client) ListWorkflows() ([]*models.ActionWorkflow, error) {
	resp, err := c.sendRequest("GET", "/actions/workflows", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list workflows (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.ActionWorkflow `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) ListRuns() ([]*models.ActionRun, error) {
	resp, err := c.sendRequest("GET", "/actions/runs", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list action runs (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.ActionRun `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) RetryRun(id string) (*models.ActionRun, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/actions/runs/%s/retry", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to retry action run (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.ActionRun `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) CancelRun(id string) (*models.ActionRun, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/actions/runs/%s/cancel", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to cancel action run (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.ActionRun `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) ListRunners() ([]*models.ActionRunner, error) {
	resp, err := c.sendRequest("GET", "/actions/runners", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list action runners (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.ActionRunner `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
