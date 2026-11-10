package chatapps

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func (d *Dispatcher) executeAuthenticatedCommand(
	ctx context.Context,
	req ChatCommandRequest,
	identity *ChatAppIdentity,
	cmd, args string,
) (*ChatCommandResponse, error) {
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
		return &ChatCommandResponse{Text: helpText}, nil

	case "projects":
		if d.runner == nil {
			return &ChatCommandResponse{Text: "Project service unavailable"}, nil
		}
		projectsList, err := d.runner.ListProjects(ctx, 25, 0)
		if err != nil {
			return nil, err
		}
		if len(projectsList) == 0 {
			return &ChatCommandResponse{Text: "No projects found in this Codedock instance."}, nil
		}

		var sb strings.Builder
		sb.WriteString("### 📂 Codedock Projects\n\n")
		for _, p := range projectsList {
			fmt.Fprintf(&sb, "**%s** (`%s`)\n", p.Name, p.ID)
			apps, err := d.runner.ListAppServices(ctx, p.ID)
			if err == nil {
				for _, a := range apps {
					domain := a.Domain
					if domain == "" {
						domain = "no public domain"
					}
					fmt.Fprintf(&sb, "  • Service `%s`: **%s** (%s)\n", a.Name, a.Status, domain)
				}
			}
		}
		return &ChatCommandResponse{Text: sb.String()}, nil

	case "status":
		target := strings.TrimSpace(args)
		if target == "" {
			return &ChatCommandResponse{Text: "Usage: `/codedock status <project_or_service_id>`"}, nil
		}
		if d.runner != nil {
			if app, err := d.runner.GetAppService(ctx, target); err == nil && app != nil {
				msg := fmt.Sprintf(
					"### 📦 Service Status: %s\n\n"+
						"- **ID**: `%s`\n"+
						"- **Status**: `%s`\n"+
						"- **Image**: `%s`\n"+
						"- **Branch**: `%s`\n"+
						"- **Domain**: %s\n"+
						"- **Updated**: %s\n",
					app.Name, app.ID, app.Status, app.ImageRef, app.Branch,
					app.Domain, app.UpdatedAt,
				)
				return &ChatCommandResponse{Text: msg}, nil
			}
			if proj, err := d.runner.GetProject(ctx, target); err == nil && proj != nil {
				msg := fmt.Sprintf("### 📂 Project: %s (`%s`)\n%s\n", proj.Name, proj.ID, proj.Description)
				return &ChatCommandResponse{Text: msg}, nil
			}
		}
		return &ChatCommandResponse{Text: fmt.Sprintf("Resource `%s` not found.", target)}, nil

	case "logs":
		target := strings.TrimSpace(args)
		if target == "" {
			return &ChatCommandResponse{Text: "Usage: `/codedock logs <service_id>`"}, nil
		}
		deepLink := fmt.Sprintf("/services/%s", target)
		status := "Unknown"
		name := target
		if d.runner != nil {
			if app, err := d.runner.GetAppService(ctx, target); err == nil && app != nil {
				status = app.Status
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
		return &ChatCommandResponse{Text: msg, DeepLink: &deepLink}, nil

	case "redeploy", "restart":
		target := strings.TrimSpace(args)
		if target == "" {
			return &ChatCommandResponse{Text: fmt.Sprintf("Usage: `/codedock %s <service_id>`", cmd)}, nil
		}
		var appName string = target
		if d.runner != nil {
			if app, err := d.runner.GetAppService(ctx, target); err == nil && app != nil {
				appName = app.Name
			}
		}

		tokenBytes := make([]byte, 8)
		_, _ = rand.Read(tokenBytes)
		confirmToken := "cf_" + hex.EncodeToString(tokenBytes)

		dbCmd := &ChatAppCommand{
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
		return &ChatCommandResponse{
			Text:                 msg,
			RequiresConfirmation: true,
			ConfirmationToken:    &confirmToken,
		}, nil

	case "confirm":
		token := strings.TrimSpace(args)
		if token == "" {
			return &ChatCommandResponse{Text: "Usage: `/codedock confirm <token>`"}, nil
		}
		pendingCmd, err := d.service.Repo().GetCommandByConfirmationToken(ctx, token)
		if err != nil || pendingCmd == nil {
			return &ChatCommandResponse{Text: "❌ Invalid or expired confirmation token."}, nil
		}

		var resultMsg string
		if d.runner != nil {
			if pendingCmd.Command == "redeploy" {
				depID, err := d.runner.RedeployApp(ctx, pendingCmd.Args, req.ExternalUsername, string(req.Provider))
				if err == nil {
					resultMsg = fmt.Sprintf("🚀 **Redeployment Triggered** for `%s` (Deployment `%s`).", pendingCmd.Args, depID)
				}
			} else if pendingCmd.Command == "restart" {
				if err := d.runner.RestartApp(ctx, pendingCmd.Args); err == nil {
					resultMsg = fmt.Sprintf("🔄 **Service Restarted** successfully for `%s`.", pendingCmd.Args)
				}
			}
		}

		if resultMsg == "" {
			resultMsg = fmt.Sprintf("Action `%s` executed for resource `%s`.", pendingCmd.Command, pendingCmd.Args)
		}

		_ = d.service.Repo().MarkCommandConfirmed(ctx, pendingCmd.ID, resultMsg)
		return &ChatCommandResponse{Text: resultMsg}, nil

	case "ask":
		query := strings.TrimSpace(args)
		if query == "" {
			return &ChatCommandResponse{Text: "Usage: `/codedock ask <question>`"}, nil
		}
		return &ChatCommandResponse{
			Text: fmt.Sprintf("🤖 **Codedock AI** (User: `%s`):\nI received your query: %q\nTo view full interactive analysis and tool execution cards, open the Codedock dashboard.", identity.ExternalUsername, query),
		}, nil

	default:
		return &ChatCommandResponse{
			Text: fmt.Sprintf("Unknown command: `%s`. Type `/codedock help` for available commands.", cmd),
		}, nil
	}
}
