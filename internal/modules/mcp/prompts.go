package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"codedock/internal/config"
	"codedock/internal/version"
)

type PromptArg struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

type PromptDef struct {
	Name        string      `json:"name"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Arguments   []PromptArg `json:"arguments,omitempty"`
}

type ResourceDef struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

func ListPrompts() []PromptDef {
	return []PromptDef{
		{
			Name:        "codedock-overview",
			Title:       "Codedock: How to drive the platform via MCP",
			Description: "Orientation on Codedock tool groups, workspace scoping, and safe operational workflows.",
		},
		{
			Name:        "deploy-from-git",
			Title:       "Deploy an application from a Git repository",
			Description: "End-to-end flow: verify Git connection, inspect repository, create project/service, deploy, and tail logs.",
			Arguments: []PromptArg{
				{Name: "repo", Description: "Repository owner/name or URL"},
				{Name: "branch", Description: "Branch to deploy (defaults to main)"},
			},
		},
		{
			Name:        "install-catalog-app",
			Title:       "Install a one-click catalog application",
			Description: "Browse the one-click template catalog, review configuration, and deploy.",
			Arguments:   []PromptArg{{Name: "app", Description: "Template ID (e.g. postgres, redis, wordpress, n8n)"}},
		},
		{
			Name:        "diagnose-and-fix-issues",
			Title:       "Diagnose unhealthy services and failed deployments",
			Description: "Inspect active platform issues, read deployment/runtime logs, and remediate or rollback.",
		},
		{
			Name:        "backup-and-restore",
			Title:       "Configure and run S3/SFTP backups and restores",
			Description: "Verify backup storage destinations, trigger backup runs, and review/apply restores.",
		},
	}
}

func GetPrompt(name string, args map[string]string) (map[string]any, bool) {
	prompts := map[string]string{
		"codedock-overview": "You are connected to a Codedock platform instance via MCP.\n" +
			"1. Start by calling `get_projects` or `get_attention` to understand the current workspace state.\n" +
			"2. Every tool maps onto a real permission-checked Codedock REST endpoint.\n" +
			"3. Read state first (`get_*`), make changes deliberately, and verify deployment/runtime health afterward.",
		"deploy-from-git": fmt.Sprintf("Deploy a Git-hosted application (repo=%q, branch=%q):\n"+
			"1. Check Git providers with `get_git_status` and list repos with `get_git_repos`.\n"+
			"2. Create or select a project (`get_projects` / `post_projects`) and environment (`get_projects_by_id_environments`).\n"+
			"3. Create the service (`post_environments_by_id_apps`) and trigger deployment (`post_services_by_serviceId_deploy`).\n"+
			"4. Poll `get_deployments_by_id` and `get_deployments_by_id_logs` until READY.", args["repo"], args["branch"]),
		"install-catalog-app": fmt.Sprintf("Install one-click catalog app %q:\n"+
			"1. Call `get_one_click` to list templates and `get_one_click_by_id` to inspect variables.\n"+
			"2. Call `post_one_click_review` then `post_one_click_deploy` with target projectId and environment variables.", args["app"]),
		"diagnose-and-fix-issues": "Diagnose platform issues:\n" +
			"1. Call `get_attention` and `get_issues` to see all active outages, crash loops, or failed builds.\n" +
			"2. For any failing deployment or service, call `get_deployments_by_id_logs` or `get_apps_by_id_runtime_logs`.\n" +
			"3. Remediate by updating environment variables (`put_services_by_serviceId_variables_by_id`), restarting (`post_apps_by_id_restart`), or rolling back (`post_deployments_by_id_rollback`).",
		"backup-and-restore": "Backup and restore workflow:\n" +
			"1. List destinations with `get_s3_destinations` / `get_sftp_destinations` and policies with `get_backups`.\n" +
			"2. Trigger an on-demand backup with `post_backups_by_id_runs` and poll `get_operations_by_operationId`.\n" +
			"3. To restore, prepare with `post_backup_records_by_recordId_restore_review` and apply with `post_backup_operations_by_operationId_apply`.",
	}
	text, ok := prompts[name]
	if !ok {
		return nil, false
	}
	return map[string]any{
		"description": "Codedock guided workflow: " + name,
		"messages": []map[string]any{
			{
				"role":    "user",
				"content": map[string]any{"type": "text", "text": text},
			},
		},
	}, true
}

func ListResources() []ResourceDef {
	return []ResourceDef{
		{
			URI:         "codedock://instance/overview",
			Name:        "Instance Overview & Topology",
			Description: "Live snapshot of Codedock version, runtime mode, projects, services, and databases.",
			MimeType:    "application/json",
		},
		{
			URI:         "codedock://instance/features",
			Name:        "Codedock Platform Capabilities & Routes",
			Description: "Reference guide of Codedock features, build engines, networking, and UI paths.",
			MimeType:    "application/json",
		},
	}
}

func (b *Bridge) ReadResource(ctx context.Context, uri string) (string, bool) {
	switch uri {
	case "codedock://instance/overview":
		var projectsList any
		var dbsList any
		var pCount, dCount int
		if b.projectService != nil {
			if ps, _, err := b.projectService.ListProjects(ctx, "", 100, 0); err == nil {
				projectsList = ps
				pCount = len(ps)
			}
		}
		if b.dbService != nil {
			if ds, err := b.dbService.ListDatabases(ctx, ""); err == nil {
				dbsList = ds
				dCount = len(ds)
			}
		}
		mode := "self-hosted"
		if config.Get().Cloud.Enabled {
			mode = "cloud"
		}
		payload, _ := json.MarshalIndent(map[string]any{
			"version":        version.Version,
			"mode":           mode,
			"projectsCount":  pCount,
			"databasesCount": dCount,
			"projects":       projectsList,
			"databases":      dbsList,
		}, "", "  ")
		return string(payload), true
	case "codedock://instance/features":
		payload, _ := json.MarshalIndent(map[string]any{
			"platform":      "Codedock",
			"orchestrators": []string{"Native Multi-Host Docker + Traefik", "K3s Compute Clusters"},
			"builders":      []string{"Railpack", "Nixpacks", "Dockerfile", "Docker Compose", "One-Click App Catalog", "Serverless"},
			"features": map[string]string{
				"dashboard":   "/ (Projects, Activity, Updates, Quick Tips)",
				"projects":    "/projects (Environments, Services, Compose Stacks, Canvas, Domains, Git Webhooks)",
				"apps":        "/apps (Installed Apps & One-Click Template Catalog)",
				"deployments": "/deployments (Live Build & Runtime Logs, Instant Rollback, Image Pinning)",
				"actions":     "/actions (GitHub-Compatible Workflows, Runs, Self-Hosted & Cloud Runners)",
				"monitoring":  "/monitoring (Unified Issues Feed, Container Health Watch, Traffic & Geo Analytics)",
				"servers":     "/servers (Multi-Host SSH Fleet, WireGuard Private Networks, K3s Clusters, Takeover)",
				"domains":     "/domains (Custom Domains, Wildcard DNS-01 TLS, Cloudflare Tunnels & Access)",
				"emails":      "/emails (Self-Hosted Mail Server, DKIM/SPF/DMARC, Mailboxes, Aliases, Webmail)",
				"backups":     "/backups (S3 & SFTP Destinations, Scheduled Policies, Volume & DB Restores)",
				"jobs":        "/jobs (Recurring Cron Jobs & Command Execution)",
				"team":        "/team (Organization Members, Roles & Integrated Audit Logs)",
				"tokens":      "/tokens (SSH Keys, Cloud Provider Tokens, Registry Credentials, Scoped PATs)",
				"settings":    "/settings (Appearance Studio, Server Routing, DNS Providers, AI Models, MCP, Updates, Maintenance)",
			},
		}, "", "  ")
		return string(payload), true
	default:
		return "", false
	}
}
