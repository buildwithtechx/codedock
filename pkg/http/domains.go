package http

import (
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"

	"codedock.run/codedock/pkg/types"
)

func (c *Client) ListDomains(serviceID string) ([]*types.DomainConfig, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("/services/%s/domains", serviceID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to list domains (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data []*types.DomainConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) AddDomain(serviceID string, req *types.DomainConfig) (*types.DomainConfig, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/services/%s/domains", serviceID), req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusCreated && resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to add domain (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Data *types.DomainConfig `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

func (c *Client) RemoveDomain(domainID string) error {
	resp, err := c.sendRequest("DELETE", fmt.Sprintf("/domains/%s", domainID), nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK && resp.StatusCode != nethttp.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to remove domain (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
