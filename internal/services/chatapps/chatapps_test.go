package chatapps_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"codedock/internal/models"
	chatappservices "codedock/internal/services/chatapps"
)

type memoryChatRepo struct {
	mu          sync.Mutex
	connections map[string]*models.ChatAppConnection
	identities  map[string]*models.ChatAppIdentity
	challenges  map[string]*models.ChatAppChallenge
	commands    map[string]*models.ChatAppCommand
}

func newMemoryRepo() *memoryChatRepo {
	return &memoryChatRepo{
		connections: make(map[string]*models.ChatAppConnection),
		identities:  make(map[string]*models.ChatAppIdentity),
		challenges:  make(map[string]*models.ChatAppChallenge),
		commands:    make(map[string]*models.ChatAppCommand),
	}
}

func (m *memoryChatRepo) CreateConnection(ctx context.Context, conn *models.ChatAppConnection) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[conn.ID] = conn
	return nil
}

func (m *memoryChatRepo) GetConnection(ctx context.Context, id string) (*models.ChatAppConnection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.connections[id]
	if !ok {
		return nil, nil
	}
	return c, nil
}

func (m *memoryChatRepo) ListConnectionsByUser(ctx context.Context, userID string) ([]*models.ChatAppConnection, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*models.ChatAppConnection
	for _, c := range m.connections {
		if c.UserID == userID {
			list = append(list, c)
		}
	}
	return list, nil
}

func (m *memoryChatRepo) UpdateConnection(ctx context.Context, conn *models.ChatAppConnection) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[conn.ID] = conn
	return nil
}

func (m *memoryChatRepo) DeleteConnection(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.connections, id)
	return nil
}

func (m *memoryChatRepo) CreateIdentity(ctx context.Context, identity *models.ChatAppIdentity) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.identities[identity.ConnectionID+":"+identity.ExternalUserID] = identity
	return nil
}

func (m *memoryChatRepo) GetIdentityByExternal(ctx context.Context, connectionID string, extUserID string) (*models.ChatAppIdentity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	idt, ok := m.identities[connectionID+":"+extUserID]
	if !ok {
		return nil, nil
	}
	return idt, nil
}

func (m *memoryChatRepo) ListIdentitiesByConnection(ctx context.Context, connectionID string) ([]*models.ChatAppIdentity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var list []*models.ChatAppIdentity
	for _, idt := range m.identities {
		if idt.ConnectionID == connectionID {
			list = append(list, idt)
		}
	}
	return list, nil
}

func (m *memoryChatRepo) DeleteIdentity(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, idt := range m.identities {
		if idt.ID == id {
			delete(m.identities, k)
			break
		}
	}
	return nil
}

func (m *memoryChatRepo) CreateChallenge(ctx context.Context, ch *models.ChatAppChallenge) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.challenges[ch.Token] = ch
	return nil
}

func (m *memoryChatRepo) GetChallengeByToken(ctx context.Context, token string) (*models.ChatAppChallenge, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ch, ok := m.challenges[token]
	if !ok {
		return nil, nil
	}
	return ch, nil
}

func (m *memoryChatRepo) VerifyChallenge(ctx context.Context, token string, code string) (*models.ChatAppChallenge, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ch := range m.challenges {
		if (token != "" && ch.Token == token) || (code != "" && ch.Code == code) {
			if time.Now().UTC().After(ch.ExpiresAt) {
				return nil, errors.New("challenge expired")
			}
			now := time.Now().UTC()
			ch.VerifiedAt = &now
			return ch, nil
		}
	}
	return nil, errors.New("invalid code or token")
}

func (m *memoryChatRepo) CreateCommand(ctx context.Context, cmd *models.ChatAppCommand) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.commands[cmd.ID] = cmd
	return nil
}

func (m *memoryChatRepo) GetCommandByConfirmationToken(ctx context.Context, token string) (*models.ChatAppCommand, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, cmd := range m.commands {
		if cmd.ConfirmationToken != nil && *cmd.ConfirmationToken == token && cmd.Status == "awaiting_confirmation" {
			return cmd, nil
		}
	}
	return nil, nil
}

func (m *memoryChatRepo) MarkCommandConfirmed(ctx context.Context, id string, response string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cmd, ok := m.commands[id]
	if ok {
		cmd.Status = "executed"
		now := time.Now().UTC()
		cmd.ConfirmedAt = &now
		cmd.Response = response
	}
	return nil
}

func TestChatApps_PairingFlow(t *testing.T) {
	ctx := context.Background()
	repo := newMemoryRepo()
	svc := chatappservices.NewService(repo, nil)

	conn, err := svc.CreateConnection(ctx, "user_1", models.CreateChatAppRequest{
		Provider: models.ChatAppProviderSlack,
		Name:     "Test Slack Bot",
	})
	if err != nil {
		t.Fatalf("create connection: %v", err)
	}

	dispatcher := chatappservices.NewDispatcher(svc, nil, nil, nil)

	// 1. Unpaired user attempts command
	resp, err := dispatcher.DispatchCommand(ctx, models.ChatCommandRequest{
		ConnectionID:   conn.ID,
		Provider:       conn.Provider,
		ExternalUserID: "slack_user_999",
		RawCommand:     "/codedock projects",
	})
	if err != nil {
		t.Fatalf("dispatch error: %v", err)
	}
	if !strings.Contains(resp.Text, "Authentication Required") {
		t.Fatalf("expected authentication required, got: %s", resp.Text)
	}

	// 2. User initiates pairing
	resp, err = dispatcher.DispatchCommand(ctx, models.ChatCommandRequest{
		ConnectionID:      conn.ID,
		Provider:          conn.Provider,
		ExternalUserID:    "slack_user_999",
		ExternalChannelID: "C123",
		RawCommand:        "/codedock link",
	})
	if err != nil {
		t.Fatalf("link error: %v", err)
	}
	if !strings.Contains(resp.Text, "Pair your Codedock Account") || !strings.Contains(resp.Text, "/codedock verify") {
		t.Fatalf("expected pairing link prompt, got: %s", resp.Text)
	}

	// Extract code from memory challenges
	var code string
	for _, ch := range repo.challenges {
		if ch.ExternalUserID == "slack_user_999" {
			code = ch.Code
			break
		}
	}
	if code == "" {
		t.Fatalf("expected challenge in repo")
	}

	// 3. Complete verification with wrong code
	resp, _ = dispatcher.DispatchCommand(ctx, models.ChatCommandRequest{
		ConnectionID:   conn.ID,
		Provider:       conn.Provider,
		ExternalUserID: "slack_user_999",
		RawCommand:     "/codedock verify 000000",
	})
	if !strings.Contains(resp.Text, "Verification failed") {
		t.Fatalf("expected verification failure for bad code, got: %s", resp.Text)
	}

	// 4. Complete verification with correct code
	resp, _ = dispatcher.DispatchCommand(ctx, models.ChatCommandRequest{
		ConnectionID:   conn.ID,
		Provider:       conn.Provider,
		ExternalUserID: "slack_user_999",
		RawCommand:     "/codedock verify " + code,
	})
	if !strings.Contains(resp.Text, "Account Successfully Linked") {
		t.Fatalf("expected success, got: %s", resp.Text)
	}

	// 5. Help command now works
	resp, _ = dispatcher.DispatchCommand(ctx, models.ChatCommandRequest{
		ConnectionID:   conn.ID,
		Provider:       conn.Provider,
		ExternalUserID: "slack_user_999",
		RawCommand:     "/codedock help",
	})
	if !strings.Contains(resp.Text, "Codedock Chat Assistant Commands") {
		t.Fatalf("expected help table, got: %s", resp.Text)
	}

	// 6. Logs command gives deep link, not secrets
	resp, _ = dispatcher.DispatchCommand(ctx, models.ChatCommandRequest{
		ConnectionID:   conn.ID,
		Provider:       conn.Provider,
		ExternalUserID: "slack_user_999",
		RawCommand:     "/codedock logs svc_alpha",
	})
	if !strings.Contains(resp.Text, "Open Secure Container Logs") || resp.DeepLink == nil {
		t.Fatalf("expected secure deep link for logs, got: %s", resp.Text)
	}

	// 7. Redeploy requires confirmation token
	resp, _ = dispatcher.DispatchCommand(ctx, models.ChatCommandRequest{
		ConnectionID:   conn.ID,
		Provider:       conn.Provider,
		ExternalUserID: "slack_user_999",
		RawCommand:     "/codedock redeploy svc_alpha",
	})
	if !resp.RequiresConfirmation || resp.ConfirmationToken == nil {
		t.Fatalf("expected confirmation required, got: %+v", resp)
	}

	// 8. Confirm execution
	resp, _ = dispatcher.DispatchCommand(ctx, models.ChatCommandRequest{
		ConnectionID:   conn.ID,
		Provider:       conn.Provider,
		ExternalUserID: "slack_user_999",
		RawCommand:     "/codedock confirm " + *resp.ConfirmationToken,
	})
	if !strings.Contains(resp.Text, "executed") && !strings.Contains(resp.Text, "Redeployment") {
		t.Fatalf("expected executed message, got: %s", resp.Text)
	}
}
