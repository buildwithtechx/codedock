package mcp

import (
	"context"
	"net/http"

	mcpgo "github.com/mark3labs/mcp-go/mcp"
)

// MCPBridge defines the interaction between Codedock REST endpoints and Model Context Protocol.
type MCPBridge interface {
	DispatchTool(ctx context.Context, name string, args map[string]any) (*mcpgo.CallToolResult, error)
	SSEHandler() http.Handler
	MessagePostHandler() http.Handler
}
