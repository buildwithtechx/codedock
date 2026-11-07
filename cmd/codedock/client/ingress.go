package client

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock/internal/models"
)

func (c *Client) ListIngressConnections() ([]*models.IngressConnection, error) {
	resp, err := c.sendRequest("GET", "/dns/ingress", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list ingress connections (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*models.IngressConnection `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) ConnectIngress(serverId, name, connector string) (*models.IngressConnection, error) {
	req := map[string]string{
		"serverId":  serverId,
		"name":      name,
		"connector": connector,
	}
	resp, err := c.sendRequest("POST", "/dns/ingress", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to connect ingress (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.IngressConnection `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) RetryIngress(id string) (*models.IngressConnection, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/dns/ingress/%s/retry", id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to retry ingress reconciliation (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *models.IngressConnection `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) DisconnectIngress(id string) error {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("/dns/ingress/%s", id), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to disconnect ingress (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
