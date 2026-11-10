package mcp

// ToolListResponse wraps the registered MCP tools.
type ToolListResponse struct {
	Tools []ToolDef `json:"tools"`
	Count int       `json:"count"`
}
