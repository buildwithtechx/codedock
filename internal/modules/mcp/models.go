package mcp

// ToolAnnotations marks capabilities and side-effect expectations for an MCP tool.
type ToolAnnotations struct {
	ReadOnlyHint    bool `json:"readOnlyHint"`
	DestructiveHint bool `json:"destructiveHint"`
}

// ToolDef defines a registered Model Context Protocol tool mapped to a Codedock endpoint.
type ToolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Method      string          `json:"method"`
	Path        string          `json:"path"`
	Category    string          `json:"category"`
	Auth        string          `json:"auth"`
	PathParams  []string        `json:"pathParams"`
	HasBody     bool            `json:"hasBody"`
	Annotations ToolAnnotations `json:"annotations"`
	InputSchema map[string]any  `json:"inputSchema"`
}
