package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	nethttp "net/http"
)

type MCPStatusResponse struct {
	Enabled         bool     `json:"enabled"`
	ProtocolVersion string   `json:"protocolVersion"`
	ServerName      string   `json:"serverName"`
	ServerVersion   string   `json:"serverVersion"`
	Transports      []string `json:"transports"`
	ToolCount       int      `json:"toolCount"`
	PromptCount     int      `json:"promptCount"`
}

type MCPToolDescriptor struct {
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	ReadOnly    bool   `json:"readOnly"`
	Destructive bool   `json:"destructive"`
}

type AIToolStep struct {
	ID     string `json:"id"`
	Tool   string `json:"tool"`
	Label  string `json:"label"`
	Status string `json:"status"`
}

type AIPendingAction struct {
	ActionID    string `json:"actionId"`
	Tool        string `json:"tool"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Destructive bool   `json:"destructive"`
}

type AIChatResponse struct {
	Reply          string            `json:"reply"`
	Provider       string            `json:"provider"`
	Model          string            `json:"model"`
	Mode           string            `json:"mode"`
	Steps          []AIToolStep      `json:"steps"`
	PendingActions []AIPendingAction `json:"pendingActions"`
}

func (c *Client) GetMCPStatus() (*MCPStatusResponse, error) {
	resp, err := c.sendRequest("GET", "/mcp/status", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mcp status failed (status %d): %s", resp.StatusCode, string(body))
	}

	var wrapped struct {
		Data *MCPStatusResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapped); err != nil {
		return nil, err
	}
	return wrapped.Data, nil
}

func (c *Client) ListMCPTools() ([]MCPToolDescriptor, error) {
	resp, err := c.sendRequest("GET", "/mcp/tools", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("mcp tools failed (status %d): %s", resp.StatusCode, string(body))
	}

	var wrapped struct {
		Data struct {
			Tools []MCPToolDescriptor `json:"tools"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapped); err != nil {
		return nil, err
	}
	return wrapped.Data.Tools, nil
}

func (c *Client) SendMCPJSONRPC(rawRequest []byte) ([]byte, error) {
	req, err := nethttp.NewRequest("POST", fmt.Sprintf("%s/api/v1/mcp", c.BaseURL), bytes.NewReader(rawRequest))
	if err != nil {
		return nil, err
	}
	if c.Token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == nethttp.StatusNoContent {
		return nil, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("mcp request failed (status %d): %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) AskAI(prompt, mode string) (*AIChatResponse, error) {
	if mode == "" {
		mode = "agent"
	}
	payload := map[string]any{
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"mode":        mode,
		"currentPath": "/cli",
	}
	resp, err := c.sendRequest("POST", "/ai/assistant/chat", payload)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ai chat failed (status %d): %s", resp.StatusCode, string(body))
	}

	var wrapped struct {
		Data *AIChatResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapped); err != nil {
		return nil, err
	}
	return wrapped.Data, nil
}

func (c *Client) ConfirmAIAction(actionID string) (map[string]any, error) {
	resp, err := c.sendRequest("POST", fmt.Sprintf("/ai/assistant/actions/%s/confirm", actionID), map[string]any{})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != nethttp.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("confirm action failed (status %d): %s", resp.StatusCode, string(body))
	}

	var wrapped struct {
		Data map[string]any `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapped); err != nil {
		return nil, err
	}
	return wrapped.Data, nil
}
