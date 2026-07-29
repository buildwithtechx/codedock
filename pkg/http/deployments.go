package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	nethttp "net/http"
	"os"
	"path/filepath"

	"codedock.run/codedock/pkg/types"
)

type ArchiveDeployResult struct {
	ContainerID string `json:"containerId"`
	AppID       string `json:"appId"`
	AppName     string `json:"appName"`
}

func (c *Client) TriggerDeployment(serviceID string) (*types.Deployment, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/services/%s/deploy", serviceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to trigger deployment (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *types.Deployment `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) DeployArchive(projectID, appName, archivePath string) (*ArchiveDeployResult, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening archive: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	_ = writer.WriteField("projectId", projectID)
	if appName != "" {
		_ = writer.WriteField("name", appName)
	}

	part, err := writer.CreateFormFile("file", filepath.Base(archivePath))
	if err != nil {
		return nil, fmt.Errorf("failed creating form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed copying archive content: %w", err)
	}
	writer.Close()

	req, err := nethttp.NewRequest("POST", c.BaseURL+"/api/deploy/archive", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("archive deploy request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusCreated {
		respBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deploy archive failed (status %d): %s", resp.StatusCode, string(respBytes))
	}

	var res struct {
		Data *ArchiveDeployResult `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Data, nil
}

func (c *Client) GetDeploymentStatus(deploymentID string) (*types.Deployment, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/deployments/%s", deploymentID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch deployment (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *types.Deployment `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) ListDeployments(serviceID string) ([]types.Deployment, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/services/%s/deployments", serviceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list deployments (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []types.Deployment `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) GetDeploymentLogs(deploymentID string) (string, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/deployments/%s/logs", deploymentID), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to fetch logs (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data, nil
}

func (c *Client) GetServiceMetrics(serviceID string) ([]types.ServiceMetric, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/services/%s/metrics", serviceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch metrics (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []types.ServiceMetric `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
