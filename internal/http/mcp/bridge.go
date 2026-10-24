package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"codedock/internal/services/databases"
	"codedock/internal/services/projects"
	"codedock/internal/version"
	mcpgo "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Bridge struct {
	server         *server.MCPServer
	sseServer      *server.SSEServer
	dispatcher     http.Handler
	tools          []ToolDef
	toolByName     map[string]ToolDef
	projectService *projects.ProjectService
	dbService      *databases.DatabaseService
}

func NewBridge(ps *projects.ProjectService, db *databases.DatabaseService) *Bridge {
	mcpServer := server.NewMCPServer(
		"codedock-mcp",
		version.Version,
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)
	allTools := AllTools()
	byName := make(map[string]ToolDef, len(allTools))
	for _, t := range allTools {
		byName[t.Name] = t
	}
	b := &Bridge{
		server:         mcpServer,
		sseServer:      server.NewSSEServer(mcpServer),
		tools:          allTools,
		toolByName:     byName,
		projectService: ps,
		dbService:      db,
	}
	b.registerTools()
	return b
}

func (b *Bridge) SetDispatcher(h http.Handler) {
	b.dispatcher = h
}

func (b *Bridge) MCPServer() *server.MCPServer {
	return b.server
}

func (b *Bridge) SSEServer() *server.SSEServer {
	return b.sseServer
}

func (b *Bridge) Tools() []ToolDef {
	return b.tools
}

func (b *Bridge) FindTool(name string) (ToolDef, bool) {
	t, ok := b.toolByName[name]
	return t, ok
}

func (b *Bridge) DatabaseService() *databases.DatabaseService {
	return b.dbService
}

func (b *Bridge) registerTools() {
	for _, t := range b.tools {
		toolDef := t
		rawSchema, err := json.Marshal(toolDef.InputSchema)
		var mcpTool mcpgo.Tool
		if err == nil {
			mcpTool = mcpgo.NewToolWithRawSchema(toolDef.Name, toolDef.Description, rawSchema)
		} else {
			mcpTool = mcpgo.NewTool(toolDef.Name, mcpgo.WithDescription(toolDef.Description))
		}

		b.server.AddTool(mcpTool, func(ctx context.Context, req mcpgo.CallToolRequest) (*mcpgo.CallToolResult, error) {
			if b.dispatcher == nil {
				return mcpgo.NewToolResultError(fmt.Sprintf("tool %s requires HTTP dispatcher", toolDef.Name)), nil
			}
			args, _ := req.Params.Arguments.(map[string]any)
			res := b.DispatchTool(ctx, toolDef, args, DispatchMeta{})
			if !res.OK {
				return mcpgo.NewToolResultError(res.Text), nil
			}
			return mcpgo.NewToolResultText(res.Text), nil
		})
	}
}
