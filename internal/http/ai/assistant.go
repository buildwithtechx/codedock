package ai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"codedock/internal/config"
	"codedock/internal/http/mcp"
	"codedock/internal/models"
	"codedock/internal/services/databases"
	"codedock/internal/services/projects"
	"codedock/internal/version"

	"github.com/google/uuid"
)

type ResourceTag struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Subtitle string `json:"subtitle,omitempty"`
	Status   string `json:"status,omitempty"`
	Href     string `json:"href,omitempty"`
}

type ToolStep struct {
	ID     string `json:"id"`
	Tool   string `json:"tool"`
	Label  string `json:"label"`
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

type PendingAction struct {
	ActionID    string         `json:"actionId"`
	Tool        string         `json:"tool"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Method      string         `json:"method"`
	Path        string         `json:"path"`
	Arguments   map[string]any `json:"arguments"`
	Destructive bool           `json:"destructive"`
}

type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	Mode        string        `json:"mode"`
	Provider    string        `json:"provider,omitempty"`
	Model       string        `json:"model,omitempty"`
	CurrentPath string        `json:"currentPath,omitempty"`
	Mentions    []ResourceTag `json:"mentions,omitempty"`
}

type ChatResponse struct {
	Reply           string          `json:"reply"`
	Provider        string          `json:"provider"`
	Model           string          `json:"model"`
	Configured      bool            `json:"configured"`
	Steps           []ToolStep      `json:"steps,omitempty"`
	PendingActions  []PendingAction `json:"pendingActions,omitempty"`
	ExecutedActions []string        `json:"executedActions,omitempty"`
}

type Dependencies struct {
	Bridge         *mcp.Bridge
	ProjectService *projects.ProjectService
	AppService     *projects.AppService
	DBService      *databases.DatabaseService
	Settings       *models.AISettings
}

var (
	pendingActionsMu sync.Mutex
	pendingActions   = map[string]PendingAction{}
)

func ListResources(ctx context.Context, deps Dependencies) []ResourceTag {
	var items []ResourceTag
	if deps.ProjectService != nil {
		projectList, _, _ := deps.ProjectService.ListProjects(ctx, "", 50, 0)
		for _, p := range projectList {
			items = append(items, ResourceTag{ID: p.ID, Type: "project", Name: p.Name, Subtitle: p.Description, Href: "/projects/" + p.ID})
			if deps.AppService != nil {
				apps, _ := deps.AppService.ListByProject(ctx, p.ID)
				for _, a := range apps {
					items = append(items, ResourceTag{ID: a.ID, Type: "service", Name: a.Name, Subtitle: p.Name + " · " + a.ImageRef, Status: string(a.Status), Href: "/services/" + a.ID})
				}
			}
		}
	}
	if deps.DBService != nil {
		dbs, _ := deps.DBService.ListDatabases(ctx, "")
		for _, d := range dbs {
			items = append(items, ResourceTag{ID: d.ID, Type: "database", Name: d.Name, Subtitle: d.Engine, Status: d.Status, Href: "/projects"})
		}
	}
	return items
}

func buildInstanceContext(ctx context.Context, currentPath string, mentions []ResourceTag, deps Dependencies) string {
	var sb strings.Builder
	mode := "self-hosted"
	if config.Get().Cloud.Enabled {
		mode = "cloud"
	}
	sb.WriteString("You are Codedock AI, the built-in DevOps and infrastructure co-pilot for this Codedock instance.\n")
	fmt.Fprintf(&sb, "Instance Version: %s | Mode: %s | User Current UI Route: %s\n\n", version.Version, mode, currentPath)

	if deps.ProjectService != nil {
		projectList, _, _ := deps.ProjectService.ListProjects(ctx, "", 30, 0)
		fmt.Fprintf(&sb, "### Live Instance Inventory (%d Projects)\n", len(projectList))
		for _, p := range projectList {
			fmt.Fprintf(&sb, "- Project %q (id=%s): %s\n", p.Name, p.ID, p.Description)
			if deps.AppService != nil {
				apps, _ := deps.AppService.ListByProject(ctx, p.ID)
				for _, a := range apps {
					fmt.Fprintf(&sb, "  • Service %q (id=%s, status=%s, image=%s, repo=%s, branch=%s, domain=%s)\n",
						a.Name, a.ID, a.Status, a.ImageRef, a.RepositoryURL, a.Branch, a.Domain)
				}
			}
		}
	}

	if deps.DBService != nil {
		if dbs, _ := deps.DBService.ListDatabases(ctx, ""); len(dbs) > 0 {
			fmt.Fprintf(&sb, "\n### Managed Databases (%d)\n", len(dbs))
			for _, d := range dbs {
				fmt.Fprintf(&sb, "- Database %q (id=%s, engine=%s, status=%s, port=%d)\n", d.Name, d.ID, d.Engine, d.Status, d.Port)
			}
		}
	}

	if len(mentions) > 0 {
		sb.WriteString("\n### Tagged @Resources in User Prompt\n")
		for _, m := range mentions {
			fmt.Fprintf(&sb, "- @%s:%s (id=%s, status=%s, href=%s)\n", m.Type, m.Name, m.ID, m.Status, m.Href)
		}
	}

	if deps.Bridge != nil {
		if featuresJSON, ok := deps.Bridge.ReadResource(ctx, "codedock://instance/features"); ok {
			sb.WriteString("\n### Codedock Platform Capabilities & UI Routes\n")
			sb.WriteString(featuresJSON)
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\n### Guidelines\n- Reference real projects, services, IDs, and statuses from this instance.\n- Provide clickable Markdown links to dashboard routes when helpful (e.g. `[Projects](/projects)`).\n- Use MCP tools to inspect logs, check stats, or trigger actions (`post_apps_by_id_redeploy`, `post_apps_by_id_restart`, `post_apps_by_id_stop`).\n")
	return sb.String()
}

func curatedTools(bridge *mcp.Bridge) []map[string]any {
	if bridge == nil {
		return nil
	}
	names := []string{
		"get_projects", "get_projects_by_id", "get_apps", "get_apps_by_id",
		"post_apps_by_id_redeploy", "post_apps_by_id_restart", "post_apps_by_id_stop",
		"get_databases", "get_databases_by_id", "get_deployments", "get_system_stats",
		"get_attention", "get_issues", "get_one_click", "get_backups", "get_jobs",
	}
	var out []map[string]any
	for _, n := range names {
		if t, ok := bridge.FindTool(n); ok {
			out = append(out, map[string]any{
				"type":     "function",
				"function": map[string]any{"name": t.Name, "description": t.Description, "parameters": t.InputSchema},
			})
		}
	}
	return out
}

func HandleChat(ctx context.Context, req ChatRequest, meta mcp.DispatchMeta, deps Dependencies) (ChatResponse, error) {
	sysPrompt := buildInstanceContext(ctx, req.CurrentPath, req.Mentions, deps)
	ep, err := resolveProviderEndpoint(deps.Settings, req.Provider, req.Model)
	if err != nil {
		reply, steps, pending := buildLocalFallback(ctx, req, deps)
		return ChatResponse{
			Reply: reply, Provider: "codedock-local", Model: "instance-copilot",
			Configured: false, Steps: steps, PendingActions: pending,
		}, nil
	}

	tools := curatedTools(deps.Bridge)
	msgs := append([]ChatMessage{}, req.Messages...)
	var steps []ToolStep
	var pending []PendingAction
	var executed []string

	for hop := 0; hop < 3; hop++ {
		assistantMsg, callErr := callCompletion(ctx, ep, sysPrompt, msgs, tools)
		if callErr != nil {
			return ChatResponse{}, callErr
		}
		if len(assistantMsg.ToolCalls) == 0 {
			return ChatResponse{
				Reply: assistantMsg.Content, Provider: ep.Provider, Model: ep.Model,
				Configured: true, Steps: steps, PendingActions: pending, ExecutedActions: executed,
			}, nil
		}

		msgs = append(msgs, assistantMsg)
		for _, tc := range assistantMsg.ToolCalls {
			toolDef, ok := deps.Bridge.FindTool(tc.Function.Name)
			if !ok {
				continue
			}
			var args map[string]any
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)

			if !toolDef.Annotations.ReadOnlyHint && (req.Mode != "agent" || toolDef.Annotations.DestructiveHint) {
				actID := "act_" + uuid.NewString()[:8]
				act := PendingAction{
					ActionID: actID, Tool: toolDef.Name, Title: "Run " + toolDef.Name,
					Description: toolDef.Description, Method: toolDef.Method, Path: toolDef.Path,
					Arguments: args, Destructive: toolDef.Annotations.DestructiveHint,
				}
				pendingActionsMu.Lock()
				pendingActions[actID] = act
				pendingActionsMu.Unlock()
				pending = append(pending, act)
				steps = append(steps, ToolStep{ID: tc.ID, Tool: toolDef.Name, Label: "Awaiting approval: " + toolDef.Name, Status: "pending"})
				msgs = append(msgs, ChatMessage{
					Role: "tool", ToolCallID: tc.ID,
					Content: fmt.Sprintf(`{"status":"awaiting_user_confirmation","actionId":%q}`, actID),
				})
				continue
			}

			res := deps.Bridge.DispatchTool(ctx, toolDef, args, meta)
			status := "done"
			if !res.OK {
				status = "error"
			} else if !toolDef.Annotations.ReadOnlyHint {
				executed = append(executed, toolDef.Name)
			}
			steps = append(steps, ToolStep{ID: tc.ID, Tool: toolDef.Name, Label: toolDef.Description, Status: status})
			msgs = append(msgs, ChatMessage{Role: "tool", ToolCallID: tc.ID, Content: res.Text})
		}
	}

	finalMsg, err := callCompletion(ctx, ep, sysPrompt, msgs, nil)
	if err != nil {
		return ChatResponse{}, err
	}
	return ChatResponse{
		Reply: finalMsg.Content, Provider: ep.Provider, Model: ep.Model,
		Configured: true, Steps: steps, PendingActions: pending, ExecutedActions: executed,
	}, nil
}

func ConfirmAction(ctx context.Context, actionID string, meta mcp.DispatchMeta, deps Dependencies) (map[string]any, error) {
	pendingActionsMu.Lock()
	act, ok := pendingActions[actionID]
	if ok {
		delete(pendingActions, actionID)
	}
	pendingActionsMu.Unlock()
	if !ok {
		return nil, errors.New("action expired or already executed")
	}
	toolDef, exists := deps.Bridge.FindTool(act.Tool)
	if !exists {
		return nil, errors.New("tool not found")
	}
	res := deps.Bridge.DispatchTool(ctx, toolDef, act.Arguments, meta)
	return map[string]any{
		"actionId": act.ActionID, "tool": act.Tool, "ok": res.OK, "status": res.Status, "result": res.Data,
		"executedAt": time.Now().UTC().Format(time.RFC3339),
	}, nil
}
