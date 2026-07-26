package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"codedock.run/codedock/pkg/types"
)

func (c *Client) Me() (*types.User, error) {
	resp, err := c.sendRequest("GET", "/auth/me", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get profile: %s", string(body))
	}

	var result struct {
		Data *types.User `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
