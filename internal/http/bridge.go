package http

import (
	"codedock.run/codedock/internal/services/databases"
	"codedock.run/codedock/internal/services/deployments"
	"codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/version"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type Bridge struct {
	server            *server.MCPServer
	projectService    *projects.ProjectService
	appService        *projects.AppService
	dbService         *databases.DatabaseService
	deploymentService *deployments.DeploymentService
}

func NewBridge(ps *projects.ProjectService, as *projects.AppService, db *databases.DatabaseService, ds *deployments.DeploymentService) *Bridge {
	mcpServer := server.NewMCPServer("codedock-mcp", version.Version, server.WithResourceCapabilities(true, true), server.WithPromptCapabilities(true))
	b := &Bridge{
		server:            mcpServer,
		projectService:    ps,
		appService:        as,
		dbService:         db,
		deploymentService: ds,
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
		mcp.NewTool("get_project",
			mcp.WithDescription("Get detailed information about a project by ID."),
			mcp.WithString("project_id", mcp.Required(), mcp.Description("The ID of the project")),
		),
		b.handleGetProject,
	)

	b.server.AddTool(
		mcp.NewTool("list_apps",
			mcp.WithDescription("List all application services for a project or all projects."),
			mcp.WithString("project_id", mcp.Description("Optional project ID filter")),
		),
		b.handleListApps,
	)

	b.server.AddTool(
		mcp.NewTool("get_app",
			mcp.WithDescription("Get detailed service configuration for an application."),
			mcp.WithString("service_id", mcp.Required(), mcp.Description("The ID of the application service")),
		),
		b.handleGetApp,
	)

	b.server.AddTool(
		mcp.NewTool("redeploy_app",
			mcp.WithDescription("Trigger a redeploy for an application service."),
			mcp.WithString("service_id", mcp.Required(), mcp.Description("The ID of the service to redeploy")),
		),
		b.handleRedeployApp,
	)

	b.server.AddTool(
		mcp.NewTool("restart_app",
			mcp.WithDescription("Restart an application service container."),
			mcp.WithString("service_id", mcp.Required(), mcp.Description("The ID of the service to restart")),
		),
		b.handleRestartApp,
	)

	b.server.AddTool(
		mcp.NewTool("stop_app",
			mcp.WithDescription("Stop an application service container."),
			mcp.WithString("service_id", mcp.Required(), mcp.Description("The ID of the service to stop")),
		),
		b.handleStopApp,
	)

	b.server.AddTool(
		mcp.NewTool("list_databases",
			mcp.WithDescription("List all managed databases registered in this Codedock instance."),
		),
		b.handleListDatabases,
	)

	b.server.AddTool(
		mcp.NewTool("get_database",
			mcp.WithDescription("Get connection and health details for a database by ID."),
			mcp.WithString("database_id", mcp.Required(), mcp.Description("The ID of the database")),
		),
		b.handleGetDatabase,
	)

	b.server.AddTool(
		mcp.NewTool("list_deployments",
			mcp.WithDescription("List deployment history for an application service."),
			mcp.WithString("service_id", mcp.Required(), mcp.Description("The ID of the application service")),
		),
		b.handleListDeployments,
	)

	b.server.AddTool(
		mcp.NewTool("get_system_status",
			mcp.WithDescription("Check operational status and metrics of the Codedock platform."),
		),
		b.handleGetSystemStatus,
	)
}
