package http

import (
	"context"
	"fmt"

	"codedock.run/codedock/internal/version"
	"github.com/mark3labs/mcp-go/mcp"
)

func getArgString(req mcp.CallToolRequest, name string) string {
	if args, ok := req.Params.Arguments.(map[string]any); ok {
		if val, ok := args[name].(string); ok {
			return val
		}
	}
	return ""
}

func (b *Bridge) handleListProjects(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, _, err := b.projectService.ListProjects(ctx, 100, 0)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list projects: %v", err)), nil
	}

	res := "Codedock Projects:\n"
	for _, p := range projects {
		res += fmt.Sprintf("- ID: %s | Name: %s | Description: %s\n", p.ID, p.Name, p.Description)
	}
	if len(projects) == 0 {
		res = "No projects found."
	}
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleGetProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := getArgString(request, "project_id")
	if projectID == "" {
		return mcp.NewToolResultError("project_id argument is required"), nil
	}

	p, err := b.projectService.GetProject(ctx, projectID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("project not found: %v", err)), nil
	}

	res := fmt.Sprintf("Project Details:\nID: %s\nName: %s\nDescription: %s\nCreatedAt: %v\n", p.ID, p.Name, p.Description, p.CreatedAt)
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleListApps(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projectID := getArgString(request, "project_id")
	if projectID == "" {
		projects, _, err := b.projectService.ListProjects(ctx, 100, 0)
		if err != nil || len(projects) == 0 {
			return mcp.NewToolResultText("No applications found."), nil
		}
		projectID = projects[0].ID
	}

	apps, err := b.appService.ListByProject(ctx, projectID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list apps: %v", err)), nil
	}

	res := fmt.Sprintf("Codedock Applications for Project %s:\n", projectID)
	for _, a := range apps {
		res += fmt.Sprintf("- ID: %s | Name: %s | Status: %s | ImageRef: %s | Repo: %s\n", a.ID, a.Name, a.Status, a.ImageRef, a.RepositoryURL)
	}
	if len(apps) == 0 {
		res = "No applications found."
	}
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleGetApp(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	serviceID := getArgString(request, "service_id")
	if serviceID == "" {
		return mcp.NewToolResultError("service_id argument is required"), nil
	}

	app, err := b.appService.GetAppService(ctx, serviceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("app service not found: %v", err)), nil
	}

	res := fmt.Sprintf("App Service Details:\nID: %s\nName: %s\nProjectID: %s\nStatus: %s\nImageRef: %s\nRepo: %s\nBranch: %s\nInternalPort: %d\nDomain: %s\n",
		app.ID, app.Name, app.ProjectID, app.Status, app.ImageRef, app.RepositoryURL, app.Branch, app.InternalPort, app.Domain)
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleRedeployApp(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	serviceID := getArgString(request, "service_id")
	if serviceID == "" {
		return mcp.NewToolResultError("service_id argument is required"), nil
	}

	app, err := b.appService.GetAppService(ctx, serviceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("app service not found: %v", err)), nil
	}

	app.Status = "building"
	if err := b.appService.UpdateAppService(ctx, app); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to update status: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Redeployment initiated for service %s (%s)", app.Name, serviceID)), nil
}

func (b *Bridge) handleRestartApp(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	serviceID := getArgString(request, "service_id")
	if serviceID == "" {
		return mcp.NewToolResultError("service_id argument is required"), nil
	}

	app, err := b.appService.GetAppService(ctx, serviceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("app service not found: %v", err)), nil
	}

	app.Status = "restarting"
	_ = b.appService.UpdateAppService(ctx, app)

	return mcp.NewToolResultText(fmt.Sprintf("Service %s (%s) restart signal sent", app.Name, serviceID)), nil
}

func (b *Bridge) handleStopApp(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	serviceID := getArgString(request, "service_id")
	if serviceID == "" {
		return mcp.NewToolResultError("service_id argument is required"), nil
	}

	app, err := b.appService.GetAppService(ctx, serviceID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("app service not found: %v", err)), nil
	}

	app.Status = "stopped"
	_ = b.appService.UpdateAppService(ctx, app)

	return mcp.NewToolResultText(fmt.Sprintf("Service %s (%s) stopped", app.Name, serviceID)), nil
}

func (b *Bridge) handleListDatabases(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dbs, err := b.dbService.ListDatabases(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list databases: %v", err)), nil
	}

	res := "Codedock Databases:\n"
	for _, d := range dbs {
		res += fmt.Sprintf("- ID: %s | Name: %s | Engine: %s | Status: %s | Port: %d\n", d.ID, d.Name, d.Engine, d.Status, d.Port)
	}
	if len(dbs) == 0 {
		res = "No databases found."
	}
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleGetDatabase(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	dbID := getArgString(request, "database_id")
	if dbID == "" {
		return mcp.NewToolResultError("database_id argument is required"), nil
	}

	db, err := b.dbService.GetDatabase(ctx, dbID)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("database not found: %v", err)), nil
	}

	res := fmt.Sprintf("Database Details:\nID: %s\nName: %s\nEngine: %s\nStatus: %s\nPort: %d\nInternalHost: %s\nUsername: %s\nDatabaseName: %s\n",
		db.ID, db.Name, db.Engine, db.Status, db.Port, db.InternalHost, db.Username, db.DatabaseName)
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleListDeployments(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	serviceID := getArgString(request, "service_id")
	if serviceID == "" {
		return mcp.NewToolResultError("service_id argument is required"), nil
	}

	deps, _, err := b.deploymentService.ListByService(ctx, serviceID, 100, 0)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list deployments: %v", err)), nil
	}

	res := fmt.Sprintf("Deployments for Service %s:\n", serviceID)
	for _, d := range deps {
		res += fmt.Sprintf("- ID: %s | Status: %s | CommitHash: %s | Trigger: %s | CreatedAt: %v\n", d.ID, d.Status, d.CommitHash, d.Trigger, d.CreatedAt)
	}
	if len(deps) == 0 {
		res = "No deployments found."
	}
	return mcp.NewToolResultText(res), nil
}

func (b *Bridge) handleGetSystemStatus(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	res := fmt.Sprintf("Codedock Platform Status: OK\nControl Plane Engine: Active\nVersion: %s", version.Version)
	return mcp.NewToolResultText(res), nil
}
