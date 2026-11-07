package ai

import (
	"context"
	"fmt"
	"strings"

	"codedock/internal/version"
	"github.com/google/uuid"
)

func buildLocalFallback(ctx context.Context, req ChatRequest, deps Dependencies) (string, []ToolStep, []PendingAction) {
	lastPrompt := strings.ToLower(req.Messages[len(req.Messages)-1].Content)
	var projectNames []string
	var projectsCount, dbsCount int
	var targetServiceID, targetServiceName string

	for _, m := range req.Mentions {
		if m.Type == "service" {
			targetServiceID = m.ID
			targetServiceName = m.Name
			break
		}
	}

	if deps.ProjectService != nil {
		if projectList, _, err := deps.ProjectService.ListProjects(ctx, "", 50, 0); err == nil {
			projectsCount = len(projectList)
			for _, p := range projectList {
				projectNames = append(projectNames, fmt.Sprintf("[%s](/projects/%s)", p.Name, p.ID))
				if targetServiceID == "" && deps.AppService != nil {
					apps, _ := deps.AppService.ListByProject(ctx, p.ID)
					for _, a := range apps {
						if strings.Contains(lastPrompt, strings.ToLower(a.Name)) {
							targetServiceID = a.ID
							targetServiceName = a.Name
							break
						}
					}
				}
			}
		}
	}
	if deps.DBService != nil {
		if dbList, err := deps.DBService.ListDatabases(ctx, ""); err == nil {
			dbsCount = len(dbList)
		}
	}

	steps := []ToolStep{
		{
			ID:     "step_inv",
			Tool:   "get_projects",
			Label:  fmt.Sprintf("Inspected %d projects and %d databases in instance", projectsCount, dbsCount),
			Status: "done",
		},
	}
	var pending []PendingAction

	if targetServiceID != "" && deps.Bridge != nil && (strings.Contains(lastPrompt, "redeploy") || strings.Contains(lastPrompt, "restart") || strings.Contains(lastPrompt, "stop")) {
		toolName := "post_apps_by_id_redeploy"
		title := "Redeploy service " + targetServiceName
		destructive := false
		if strings.Contains(lastPrompt, "restart") {
			toolName = "post_apps_by_id_restart"
			title = "Restart service " + targetServiceName
			destructive = true
		} else if strings.Contains(lastPrompt, "stop") {
			toolName = "post_apps_by_id_stop"
			title = "Stop service " + targetServiceName
			destructive = true
		}
		if toolDef, ok := deps.Bridge.FindTool(toolName); ok {
			actID := "act_" + uuid.NewString()[:8]
			act := PendingAction{
				ActionID:    actID,
				Tool:        toolDef.Name,
				Title:       title,
				Description: toolDef.Description,
				Method:      toolDef.Method,
				Path:        toolDef.Path,
				Arguments:   map[string]any{"id": targetServiceID},
				Destructive: destructive,
			}
			pendingActionsMu.Lock()
			pendingActions[actID] = act
			pendingActionsMu.Unlock()
			pending = append(pending, act)
			return fmt.Sprintf("I found service **%s** (`%s`). Confirm the action card below to execute `%s` directly on this instance.", targetServiceName, targetServiceID, toolName), steps, pending
		}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Here is the live state of your **Codedock `%s`** instance:\n\n", version.Version)
	fmt.Fprintf(&sb, "- **Projects (%d)**: ", projectsCount)
	if len(projectNames) == 0 {
		sb.WriteString("None created yet — head to [Projects](/projects) or [Explore Apps](/apps/new) to deploy your first workload.\n")
	} else {
		sb.WriteString(strings.Join(projectNames, ", "))
		sb.WriteString("\n")
	}
	fmt.Fprintf(&sb, "- **Managed Databases (%d)**: %d provisioned\n", dbsCount, dbsCount)
	toolsCount := 0
	if deps.Bridge != nil {
		toolsCount = len(deps.Bridge.Tools())
	}
	fmt.Fprintf(&sb, "- **MCP Tools Active**: `%d` route-backed tools exposed at `/api/mcp`\n\n", toolsCount)
	sb.WriteString("> **No External AI Provider Configured:** Running in **Local Instance Copilot** mode. To unlock full multi-turn LLM reasoning, add an API key and select a default provider in **[Settings → AI Models](/settings?tab=ai)**.")
	return sb.String(), steps, pending
}
