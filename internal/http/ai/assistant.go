package ai

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"codedock/internal/config"
	"codedock/internal/http/mcp"
	aimodule "codedock/internal/modules/ai"
	"codedock/internal/modules/apps"
	"codedock/internal/modules/databases"
	"codedock/internal/modules/projects"
	"codedock/internal/version"
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
	Arguments   map[string]any `json:"arguments,omitempty"`
	Destructive bool           `json:"destructive"`
}

type QuestionOption struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

type QuestionItem struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	AllowOther  bool             `json:"allowOther,omitempty"`
	MultiSelect bool             `json:"multiSelect,omitempty"`
	Options     []QuestionOption `json:"options,omitempty"`
}

type QuestionGroup struct {
	ID        string         `json:"id"`
	Title     string         `json:"title"`
	Questions []QuestionItem `json:"questions"`
}

type ChatRequest struct {
	Messages    []ChatMessage `json:"messages"`
	Provider    string        `json:"provider,omitempty"`
	Model       string        `json:"model,omitempty"`
	Persona     string        `json:"persona,omitempty"`
	CurrentPath string        `json:"currentPath,omitempty"`
	Mentions    []ResourceTag `json:"mentions,omitempty"`
	Mode        string        `json:"mode,omitempty"`
}

type ChatResponse struct {
	Reply           string          `json:"reply"`
	Provider        string          `json:"provider"`
	Model           string          `json:"model"`
	Configured      bool            `json:"configured"`
	Steps           []ToolStep      `json:"steps,omitempty"`
	PendingActions  []PendingAction `json:"pendingActions,omitempty"`
	ExecutedActions []string        `json:"executedActions,omitempty"`
	Questions       []QuestionGroup `json:"questions,omitempty"`
}

type Dependencies struct {
	Bridge         *mcp.Bridge
	ProjectService *projects.Service
	AppService     *apps.Service
	DBService      *databases.DatabaseService
	Settings       *aimodule.AISettings
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
				appsList, _ := deps.AppService.ListByProject(ctx, p.ID)
				for _, a := range appsList {
					items = append(items, ResourceTag{ID: a.ID, Type: "service", Name: a.Name, Subtitle: p.Name + " · " + a.ImageRef, Status: string(a.Status), Href: "/services/" + a.ID})
				}
			}
		}
	}
	if deps.DBService != nil {
		dbs, _ := deps.DBService.ListDatabases(ctx, "")
		for _, d := range dbs {
			items = append(items, ResourceTag{ID: d.ID, Type: "database", Name: d.Name, Subtitle: string(d.Engine), Status: string(d.Status), Href: "/projects"})
		}
	}
	return items
}

func buildInstanceContext(ctx context.Context, currentPath string, mentions []ResourceTag, deps Dependencies, persona string) string {
	var sb strings.Builder
	mode := "self-hosted"
	if config.Get().Cloud.Enabled {
		mode = "cloud"
	}
	switch strings.ToLower(strings.TrimSpace(persona)) {
	case "logs":
		sb.WriteString("You are Codedock Logs Investigator, a specialized diagnostic engineer for this Codedock instance.\n")
		sb.WriteString("Your primary mission: trace container crashes, exit codes (e.g. 137 OOMKilled, 1 crash), stderr/stdout logs, restart loops, and failed healthchecks. Pinpoint failure lines, explain causes, and provide exact remediations.\n")
		sb.WriteString("Strict Enforcement: You operate in strict read-only mode. You cannot execute mutating operations.\n\n")
	case "security":
		sb.WriteString("You are Codedock Security Reviewer, a specialized container and platform security auditor for this Codedock instance.\n")
		sb.WriteString("Your primary mission: inspect signed runtime advisories, container digests, CVE disclosures, and precision semver patch reviews.\n")
		sb.WriteString("Strict Enforcement: You operate in strict read-only mode. You cannot execute mutating operations.\n\n")
	case "devops":
		sb.WriteString("You are Codedock DevOps Co-pilot, an active infrastructure coordinator for this Codedock instance.\n")
		sb.WriteString("Your primary mission: monitor instance health, configure domains and routing, execute authorized deployments, and scale container replicas.\n\n")
	default:
		sb.WriteString("You are Codedock AI, the built-in DevOps and infrastructure co-pilot for this Codedock instance.\n")
	}
	fmt.Fprintf(&sb, "Instance Version: %s | Mode: %s | User Current UI Route: %s | Active Persona: %s\n\n", version.Version, mode, currentPath, PersonaDisplayName(persona))

	if deps.ProjectService != nil {
		projectList, _, _ := deps.ProjectService.ListProjects(ctx, "", 30, 0)
		fmt.Fprintf(&sb, "### Live Instance Inventory (%d Projects)\n", len(projectList))
		for _, p := range projectList {
			fmt.Fprintf(&sb, "- Project %q (id=%s): %s\n", p.Name, p.ID, p.Description)
			if deps.AppService != nil {
				appsList, _ := deps.AppService.ListByProject(ctx, p.ID)
				for _, a := range appsList {
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
