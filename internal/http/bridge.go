package http

import (
	"context"
	"fmt"

	"codedock.run/codedock/internal/services/databases"
	"codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/version"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Bridge struct {
	server         *server.MCPServer
	projectService *projects.ProjectService
	appService     *projects.AppService
	dbService      *databases.DatabaseService
}

func NewBridge(ps *projects.ProjectService, as *projects.AppService, db *databases.DatabaseService) *Bridge {
	mcpServer := server.NewMCPServer("codedock-mcp", version.Version, server.WithResourceCapabilities(true, true), server.WithPromptCapabilities(true))
	b := &Bridge{
		server:         mcpServer,
		projectService: ps,
		appService:     as,
		dbService:      db,
	}
	b.registerTools()
	return b
}

func (b *Bridge) MCPServer() *server.MCPServer {
	return b.server
}

func (b *Bridge) registerTools() {
	b.server.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("List all deployment projects registered in this Codedock instance."),
		),
		b.handleListProjects,
	)

	b.server.AddTool(
		mcp.NewTool("list_databases",
			mcp.WithDescription("List all managed databases registered in this Codedock instance."),
		),
		b.handleListDatabases,
	)

	b.server.AddTool(
		mcp.NewTool("get_system_status",
			mcp.WithDescription("Check basic operational and health metrics of the Codedock platform."),
		),
		b.handleGetSystemStatus,
	)
}

func (b *Bridge) handleListProjects(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, _, err := b.projectService.ListProjects(ctx, 100, 0)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res := "Codedock Projects:\n"
	for _, p := range projects {
		res += fmt.Sprintf("- ID: %s | Name: %s\n", p.ID, p.Name)
	}
	if len(projects) == 0 {
		res = "No projects found."
	}
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleListDatabases(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dbs, err := b.dbService.ListDatabases(ctx)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res := "Codedock Databases:\n"
	for _, d := range dbs {
		res += fmt.Sprintf("- ID: %s | Name: %s | Engine: %s | Status: %s\n", d.ID, d.Name, d.Engine, d.Status)
	}
	if len(dbs) == 0 {
		res = "No databases found."
	}
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleGetSystemStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	res := fmt.Sprintf("Codedock Status: OK\nEngine: Active\nVersion: %s", version.Version)
	return mcp.NewToolResultText(res), nil
}
