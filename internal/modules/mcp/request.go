package mcp

// ToolCallRequest represents a request to execute an MCP tool directly.
type ToolCallRequest struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}
