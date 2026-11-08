package chatapps

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"codedock/internal/models"
	"codedock/internal/services/deployments"
	"codedock/internal/services/projects"

	"github.com/google/uuid"
)

type Dispatcher struct {
	service           *Service
	projectService    *projects.ProjectService
	appService        *projects.AppService
	deploymentService *deployments.DeploymentService
}

func NewDispatcher(
	service *Service,
	projectService *projects.ProjectService,
	appService *projects.AppService,
	deploymentService *deployments.DeploymentService,
) *Dispatcher {
	return &Dispatcher{
		service:           service,
		projectService:    projectService,
		appService:        appService,
		deploymentService: deploymentService,
	}
}

func parseCommandLine(raw string) (string, string) {
	trimmed := strings.TrimSpace(raw)
	// Strip prefix e.g. "/codedock" or "/openship" or "codedock"
	if strings.HasPrefix(trimmed, "/") {
		trimmed = trimmed[1:]
	}
	parts := strings.Fields(trimmed)
	if len(parts) == 0 {
		return "help", ""
	}
	first := strings.ToLower(parts[0])
	if first == "codedock" || first == "openship" {
		if len(parts) == 1 {
			return "help", ""
		}
		cmd := strings.ToLower(parts[1])
		args := strings.TrimSpace(strings.TrimPrefix(trimmed, parts[0]+" "+parts[1]))
		return cmd, args
	}
	cmd := first
	args := strings.TrimSpace(strings.TrimPrefix(trimmed, parts[0]))
	return cmd, args
}

func (d *Dispatcher) DispatchCommand(ctx context.Context, req models.ChatCommandRequest) (*models.ChatCommandResponse, error) {
	cmd, args := parseCommandLine(req.RawCommand)

	// 1. Unauthenticated / Pairing Commands
	if cmd == "link" {
		ch, err := d.service.CreatePairingChallenge(ctx, req.ConnectionID, req.Provider, req.ExternalUserID, req.ExternalChannelID)
		if err != nil {
			return nil, fmt.Errorf("create pairing challenge: %w", err)
		}
		msg := fmt.Sprintf(
			"🔗 **Pair your Codedock Account**\n\n"+
				"Open this link in your browser to authorize access:\n"+
				"`/settings?tab=chat-apps&link=%s`\n\n"+
				"After approving on the dashboard, reply in this channel with:\n"+
				"`/codedock verify %s`\n*(Expires in 15 minutes)*",
			ch.Token, ch.Code,
		)
		return &models.ChatCommandResponse{Text: msg}, nil
	}

	if cmd == "verify" {
		code := strings.TrimSpace(args)
		if code == "" {
			return &models.ChatCommandResponse{Text: "Please provide your 6-digit verification code: `/codedock verify <code>`"}, nil
		}
		conn, err := d.service.GetConnection(ctx, req.ConnectionID)
		if err != nil {
			return nil, fmt.Errorf("load connection: %w", err)
		}
		identity, err := d.service.VerifyPairing(ctx, conn.UserID, "", code)
		if err != nil {
			return &models.ChatCommandResponse{Text: fmt.Sprintf("❌ Verification failed: %v", err)}, nil
		}
		msg := fmt.Sprintf(
			"✅ **Account Successfully Linked!**\n\n"+
				"Your %s identity (`%s`) is now paired with Codedock (User ID: `%s`).\n"+
				"Type `/codedock help` to view available commands.",
			identity.Provider, identity.ExternalUserID, identity.UserID,
		)
		return &models.ChatCommandResponse{Text: msg}, nil
	}

	// 2. Identity Verification (Strict Access Control)
	identity, err := d.service.Repo().GetIdentityByExternal(ctx, req.ConnectionID, req.ExternalUserID)
	if err != nil || identity == nil || !identity.Active {
		return &models.ChatCommandResponse{
			Text: "🔒 **Authentication Required**\n\n" +
				"You must pair your Codedock account before inspecting projects or triggering actions.\n" +
				"Type `/codedock link` to begin pairing.",
		}, nil
	}

	// 3. Authenticated Commands
	switch cmd {
	case "help":
		helpText := "### 🚢 Codedock Chat Assistant Commands\n\n" +
			"| Command | Description |\n" +
			"|---|---|\n" +
			"| `/codedock projects` | List all accessible projects and current health |\n" +
			"| `/codedock status <id>` | Show service details, container health, and domain |\n" +
			"| `/codedock logs <id>` | View container status and get secure logs deep link |\n" +
			"| `/codedock redeploy <id>` | Trigger redeployment (requires confirmation) |\n" +
			"| `/codedock restart <id>` | Restart container service (requires confirmation) |\n" +
			"| `/codedock confirm <token>` | Confirm a pending mutating action |\n" +
			"| `/codedock ask <query>` | Ask Codedock AI infrastructure co-pilot |\n" +
			"| `/codedock link` | Pair a new account or refresh credentials |\n"
		return &models.ChatCommandResponse{Text: helpText}, nil

	case "projects":
		if d.projectService == nil {
			return &models.ChatCommandResponse{Text: "Project service unavailable"}, nil
		}
		projectsList, _, err := d.projectService.ListProjects(ctx, "", 25, 0)
		if err != nil {
			return nil, err
		}
		if len(projectsList) == 0 {
			return &models.ChatCommandResponse{Text: "No projects found in this Codedock instance."}, nil
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("### 📂 Codedock Projects (%d)\n\n", len(projectsList)))
		for _, p := range projectsList {
			fmt.Fprintf(&sb, "- **%s** (`%s`)\n", p.Name, p.ID)
			if d.appService != nil {
				apps, _ := d.appService.ListByProject(ctx, p.ID)
				for _, a := range apps {
					domain := a.Domain
					if domain == "" {
						domain = "no public domain"
					}
					fmt.Fprintf(&sb, "  • Service `%s`: **%s** (%s)\n", a.Name, a.Status, domain)
				}
			}
		}
		return &models.ChatCommandResponse{Text: sb.String()}, nil

	case "status":
		target := strings.TrimSpace(args)
		if target == "" {
			return &models.ChatCommandResponse{Text: "Usage: `/codedock status <project_or_service_id>`"}, nil
		}
		if d.appService != nil {
			if app, err := d.appService.GetAppService(ctx, target); err == nil && app != nil {
				msg := fmt.Sprintf(
					"### 📦 Service Status: %s\n\n"+
						"- **ID**: `%s`\n"+
						"- **Status**: `%s`\n"+
						"- **Image**: `%s`\n"+
						"- **Branch**: `%s`\n"+
						"- **Domain**: %s\n"+
						"- **Updated**: %s\n",
					app.Name, app.ID, app.Status, app.ImageRef, app.Branch,
					app.Domain, app.UpdatedAt.Format(time.RFC3339),
				)
				return &models.ChatCommandResponse{Text: msg}, nil
			}
		}
		if d.projectService != nil {
			if proj, err := d.projectService.GetProject(ctx, target); err == nil && proj != nil {
				msg := fmt.Sprintf("### 📂 Project: %s (`%s`)\n%s\n", proj.Name, proj.ID, proj.Description)
				return &models.ChatCommandResponse{Text: msg}, nil
			}
		}
		return &models.ChatCommandResponse{Text: fmt.Sprintf("Resource `%s` not found.", target)}, nil

	case "logs":
		target := strings.TrimSpace(args)
		if target == "" {
			return &models.ChatCommandResponse{Text: "Usage: `/codedock logs <service_id>`"}, nil
		}
		// Security safeguard: Do not dump raw text or secrets in chat channels
		deepLink := fmt.Sprintf("/services/%s", target)
		status := "Unknown"
		name := target
		if d.appService != nil {
			if app, err := d.appService.GetAppService(ctx, target); err == nil && app != nil {
				status = string(app.Status)
				name = app.Name
			}
		}
		msg := fmt.Sprintf(
			"📋 **Service Logs for %s**\n\n"+
				"- **Current Status**: `%s`\n\n"+
				"To prevent accidental leakage of sensitive tokens or internal environment variables, "+
				"raw container logs are gated in the dashboard:\n"+
				"[Open Secure Container Logs](%s)",
			name, status, deepLink,
		)
		return &models.ChatCommandResponse{Text: msg, DeepLink: &deepLink}, nil

	case "redeploy", "restart":
		target := strings.TrimSpace(args)
		if target == "" {
			return &models.ChatCommandResponse{Text: fmt.Sprintf("Usage: `/codedock %s <service_id>`", cmd)}, nil
		}
		var appName string = target
		if d.appService != nil {
			if app, err := d.appService.GetAppService(ctx, target); err == nil && app != nil {
				appName = app.Name
			}
		}

		// Mutating operation requires confirmation token
		tokenBytes := make([]byte, 8)
		_, _ = rand.Read(tokenBytes)
		confirmToken := "cf_" + hex.EncodeToString(tokenBytes)

		dbCmd := &models.ChatAppCommand{
			ID:                "cmd_" + uuid.NewString()[:12],
			ConnectionID:      req.ConnectionID,
			IdentityID:        &identity.ID,
			Provider:          req.Provider,
			Command:           cmd,
			Args:              target,
			Status:            "awaiting_confirmation",
			ConfirmationToken: &confirmToken,
		}
		_ = d.service.Repo().CreateCommand(ctx, dbCmd)

		msg := fmt.Sprintf(
			"⚠️ **Confirmation Required**\n\n"+
				"You are about to **%s** service `%s` (`%s`).\n\n"+
				"To execute this action, reply in this channel within 5 minutes:\n"+
				"`/codedock confirm %s`",
			cmd, appName, target, confirmToken,
		)
		return &models.ChatCommandResponse{
			Text:                 msg,
			RequiresConfirmation: true,
			ConfirmationToken:    &confirmToken,
		}, nil

	case "confirm":
		token := strings.TrimSpace(args)
		if token == "" {
			return &models.ChatCommandResponse{Text: "Usage: `/codedock confirm <token>`"}, nil
		}
		pendingCmd, err := d.service.Repo().GetCommandByConfirmationToken(ctx, token)
		if err != nil || pendingCmd == nil {
			return &models.ChatCommandResponse{Text: "❌ Invalid or expired confirmation token."}, nil
		}

		// Execute the confirmed mutating action
		var resultMsg string
		if pendingCmd.Command == "redeploy" {
			if d.deploymentService != nil && d.appService != nil {
				if app, err := d.appService.GetAppService(ctx, pendingCmd.Args); err == nil && app != nil {
					newDep := &models.Deployment{
						ServiceID:     app.ID,
						EnvironmentID: app.EnvironmentID,
						ProjectID:     app.ProjectID,
						Status:        "BUILDING",
						CommitMessage: "Chat Command Redeploy by " + req.ExternalUsername,
						Branch:        app.Branch,
						Trigger:       "Chat: " + string(req.Provider),
					}
					created, err := d.deploymentService.CreateDeployment(ctx, newDep)
					if err == nil && created != nil {
						d.deploymentService.ExecuteDeploymentAsync(created)
						resultMsg = fmt.Sprintf("🚀 **Redeployment Triggered** for `%s` (Deployment `%s`).", app.Name, created.ID)
					}
				}
			}
		} else if pendingCmd.Command == "restart" {
			if d.deploymentService != nil && d.appService != nil {
				if app, err := d.appService.GetAppService(ctx, pendingCmd.Args); err == nil && app != nil {
					if err := d.deploymentService.RestartAppService(ctx, app); err == nil {
						app.Status = models.AppServiceStatusRunning
						_ = d.appService.UpdateAppService(ctx, app)
						resultMsg = fmt.Sprintf("🔄 **Service Restarted** successfully for `%s`.", app.Name)
					}
				}
			}
		}

		if resultMsg == "" {
			resultMsg = fmt.Sprintf("Action `%s` executed for resource `%s`.", pendingCmd.Command, pendingCmd.Args)
		}

		_ = d.service.Repo().MarkCommandConfirmed(ctx, pendingCmd.ID, resultMsg)
		return &models.ChatCommandResponse{Text: resultMsg}, nil

	case "ask":
		query := strings.TrimSpace(args)
		if query == "" {
			return &models.ChatCommandResponse{Text: "Usage: `/codedock ask <question>`"}, nil
		}
		return &models.ChatCommandResponse{
			Text: fmt.Sprintf("🤖 **Codedock AI** (User: `%s`):\nI received your query: %q\nTo view full interactive analysis and tool execution cards, open the Codedock dashboard.", identity.ExternalUsername, query),
		}, nil

	default:
		return &models.ChatCommandResponse{
			Text: fmt.Sprintf("Unknown command: `%s`. Type `/codedock help` for available commands.", cmd),
		}, nil
	}
}
