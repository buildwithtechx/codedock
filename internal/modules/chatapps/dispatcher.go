package chatapps

import (
	"context"
	"fmt"
	"strings"
)

type ProjectSummary struct {
	ID          string
	Name        string
	Description string
}

type AppServiceSummary struct {
	ID        string
	Name      string
	Status    string
	Domain    string
	ImageRef  string
	Branch    string
	UpdatedAt string
}

type CommandRunner interface {
	ListProjects(ctx context.Context, limit, offset int) ([]ProjectSummary, error)
	GetProject(ctx context.Context, id string) (*ProjectSummary, error)
	ListAppServices(ctx context.Context, projectID string) ([]AppServiceSummary, error)
	GetAppService(ctx context.Context, id string) (*AppServiceSummary, error)
	RedeployApp(ctx context.Context, appID, username, provider string) (string, error)
	RestartApp(ctx context.Context, appID string) error
}

type Dispatcher struct {
	service *Service
	runner  CommandRunner
}

func NewDispatcher(service *Service, runner CommandRunner) *Dispatcher {
	if runner == nil && service != nil && service.repo != nil {
		if r, ok := service.repo.(CommandRunner); ok {
			runner = r
		}
	}
	return &Dispatcher{
		service: service,
		runner:  runner,
	}
}

func parseCommandLine(raw string) (string, string) {
	trimmed := strings.TrimSpace(raw)
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

func (d *Dispatcher) DispatchCommand(ctx context.Context, req ChatCommandRequest) (*ChatCommandResponse, error) {
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
		return &ChatCommandResponse{Text: msg}, nil
	}

	if cmd == "verify" {
		code := strings.TrimSpace(args)
		if code == "" {
			return &ChatCommandResponse{Text: "Please provide your 6-digit verification code: `/codedock verify <code>`"}, nil
		}
		conn, err := d.service.GetConnection(ctx, req.ConnectionID)
		if err != nil {
			return nil, fmt.Errorf("load connection: %w", err)
		}
		identity, err := d.service.VerifyPairing(ctx, conn.UserID, "", code)
		if err != nil {
			return &ChatCommandResponse{Text: fmt.Sprintf("❌ Verification failed: %v", err)}, nil
		}
		msg := fmt.Sprintf(
			"✅ **Account Successfully Linked!**\n\n"+
				"Your %s identity (`%s`) is now paired with Codedock (User ID: `%s`).\n"+
				"Type `/codedock help` to view available commands.",
			identity.Provider, identity.ExternalUserID, identity.UserID,
		)
		return &ChatCommandResponse{Text: msg}, nil
	}

	// 2. Identity Verification (Strict Access Control)
	identity, err := d.service.Repo().GetIdentityByExternal(ctx, req.ConnectionID, req.ExternalUserID)
	if err != nil || identity == nil || !identity.Active {
		return &ChatCommandResponse{
			Text: "🔒 **Authentication Required**\n\n" +
				"You must pair your Codedock account before inspecting projects or triggering actions.\n" +
				"Type `/codedock link` to begin pairing.",
		}, nil
	}

	return d.executeAuthenticatedCommand(ctx, req, identity, cmd, args)
}
