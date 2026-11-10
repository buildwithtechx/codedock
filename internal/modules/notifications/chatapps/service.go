package chatapps

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"

	"codedock/internal/platform/vault"
)

type Service struct {
	repo  Repository
	vault *vault.Vault
}

func NewService(repo Repository, v *vault.Vault) *Service {
	return &Service{
		repo:  repo,
		vault: v,
	}
}

func (s *Service) Repo() Repository {
	return s.repo
}

func (s *Service) CreateConnection(ctx context.Context, userID string, req CreateChatAppRequest) (*ChatAppConnection, error) {
	if req.Provider == "" {
		return nil, errors.New("provider is required")
	}
	if req.Transport == "" {
		req.Transport = ChatAppTransportWebhook
	}
	if req.Purpose == "" {
		req.Purpose = ChatAppPurposeBoth
	}
	if req.Name == "" {
		req.Name = fmt.Sprintf("%s Integration", req.Provider)
	}

	encryptedCreds := ""
	if len(req.Credentials) > 0 {
		raw, err := json.Marshal(req.Credentials)
		if err != nil {
			return nil, fmt.Errorf("serialize credentials: %w", err)
		}
		if s.vault != nil {
			enc, err := s.vault.Encrypt(string(raw))
			if err != nil {
				return nil, fmt.Errorf("encrypt credentials: %w", err)
			}
			encryptedCreds = enc
		} else {
			encryptedCreds = string(raw)
		}
	}

	conn := &ChatAppConnection{
		ID:                   "chat_" + uuid.NewString()[:12],
		UserID:               userID,
		Provider:             req.Provider,
		Transport:            req.Transport,
		Purpose:              req.Purpose,
		Name:                 req.Name,
		BotUsername:          req.BotUsername,
		CredentialsEncrypted: encryptedCreds,
		Config:               req.Config,
		Enabled:              true,
		Status:               ChatAppStatusConnected,
	}

	if err := s.repo.CreateConnection(ctx, conn); err != nil {
		return nil, fmt.Errorf("create chat connection: %w", err)
	}
	return conn, nil
}

func (s *Service) GetConnection(ctx context.Context, id string) (*ChatAppConnection, error) {
	conn, err := s.repo.GetConnection(ctx, id)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, ErrConnectionNotFound
	}
	return conn, nil
}

func (s *Service) ListConnections(ctx context.Context, userID string) ([]*ChatAppConnection, error) {
	return s.repo.ListConnectionsByUser(ctx, userID)
}

func (s *Service) UpdateConnection(ctx context.Context, id string, req UpdateChatAppRequest) (*ChatAppConnection, error) {
	conn, err := s.GetConnection(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		conn.Name = *req.Name
	}
	if req.Transport != nil {
		conn.Transport = *req.Transport
	}
	if req.Purpose != nil {
		conn.Purpose = *req.Purpose
	}
	if req.BotUsername != nil {
		conn.BotUsername = *req.BotUsername
	}
	if req.Enabled != nil {
		conn.Enabled = *req.Enabled
	}
	if req.Config != nil {
		conn.Config = req.Config
	}
	if len(req.Credentials) > 0 {
		raw, err := json.Marshal(req.Credentials)
		if err != nil {
			return nil, fmt.Errorf("serialize credentials: %w", err)
		}
		if s.vault != nil {
			enc, err := s.vault.Encrypt(string(raw))
			if err != nil {
				return nil, fmt.Errorf("encrypt credentials: %w", err)
			}
			conn.CredentialsEncrypted = enc
		} else {
			conn.CredentialsEncrypted = string(raw)
		}
	}

	if err := s.repo.UpdateConnection(ctx, conn); err != nil {
		return nil, fmt.Errorf("update chat connection: %w", err)
	}
	return conn, nil
}

func (s *Service) DeleteConnection(ctx context.Context, id string) error {
	return s.repo.DeleteConnection(ctx, id)
}

func (s *Service) GetDecryptedCredentials(ctx context.Context, conn *ChatAppConnection) (map[string]string, error) {
	creds := make(map[string]string)
	if conn.CredentialsEncrypted == "" {
		return creds, nil
	}
	rawStr := conn.CredentialsEncrypted
	if s.vault != nil {
		dec, err := s.vault.Decrypt(conn.CredentialsEncrypted)
		if err == nil {
			rawStr = dec
		}
	}
	_ = json.Unmarshal([]byte(rawStr), &creds)
	return creds, nil
}

func (s *Service) TestConnection(ctx context.Context, id string) (map[string]any, error) {
	conn, err := s.GetConnection(ctx, id)
	if err != nil {
		return nil, err
	}

	creds, err := s.GetDecryptedCredentials(ctx, conn)
	if err != nil {
		return nil, err
	}

	// Validate presence of key tokens depending on provider
	hasKey := false
	switch conn.Provider {
	case ChatAppProviderSlack:
		hasKey = creds["bot_token"] != "" || creds["app_token"] != "" || creds["signing_secret"] != ""
	case ChatAppProviderDiscord:
		hasKey = creds["bot_token"] != "" || creds["public_key"] != ""
	case ChatAppProviderTelegram:
		hasKey = creds["bot_token"] != ""
	case ChatAppProviderTeams:
		hasKey = creds["app_id"] != "" || creds["app_password"] != ""
	}

	now := time.Now().UTC()
	conn.LastEventAt = &now
	if hasKey {
		conn.Status = ChatAppStatusConnected
		conn.LastError = nil
	} else {
		conn.Status = ChatAppStatusPairingRequired
		errStr := "incomplete credentials"
		conn.LastError = &errStr
	}
	_ = s.repo.UpdateConnection(ctx, conn)

	return map[string]any{
		"id":         conn.ID,
		"provider":   conn.Provider,
		"status":     conn.Status,
		"ok":         hasKey,
		"checked_at": now.Format(time.RFC3339),
	}, nil
}

func (s *Service) ListIdentities(ctx context.Context, connectionID string) ([]*ChatAppIdentity, error) {
	return s.repo.ListIdentitiesByConnection(ctx, connectionID)
}

func (s *Service) DeleteIdentity(ctx context.Context, id string) error {
	return s.repo.DeleteIdentity(ctx, id)
}

func (s *Service) CreatePairingChallenge(
	ctx context.Context,
	connectionID string,
	provider ChatAppProvider,
	extUserID string,
	extChannelID string,
) (*ChatAppChallenge, error) {
	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, fmt.Errorf("generate challenge token: %w", err)
	}
	token := hex.EncodeToString(tokenBytes)

	// 6-digit verification code
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return nil, fmt.Errorf("generate verification code: %w", err)
	}
	code := fmt.Sprintf("%06d", n.Int64()+100000)

	challenge := &ChatAppChallenge{
		ID:                "chal_" + uuid.NewString()[:12],
		ConnectionID:      connectionID,
		Provider:          provider,
		ExternalUserID:    extUserID,
		ExternalChannelID: extChannelID,
		Token:             token,
		Code:              code,
		ExpiresAt:         time.Now().UTC().Add(15 * time.Minute),
	}

	if err := s.repo.CreateChallenge(ctx, challenge); err != nil {
		return nil, fmt.Errorf("save pairing challenge: %w", err)
	}
	return challenge, nil
}

func (s *Service) GetChallenge(ctx context.Context, token string) (*ChatAppChallenge, error) {
	ch, err := s.repo.GetChallengeByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if ch == nil {
		return nil, ErrChallengeNotFound
	}
	if time.Now().UTC().After(ch.ExpiresAt) {
		return nil, ErrChallengeExpired
	}
	return ch, nil
}

func (s *Service) VerifyPairing(ctx context.Context, userID string, token string, code string) (*ChatAppIdentity, error) {
	ch, err := s.repo.VerifyChallenge(ctx, token, code)
	if err != nil {
		return nil, fmt.Errorf("verify pairing: %w", err)
	}

	identity := &ChatAppIdentity{
		ID:               "ident_" + uuid.NewString()[:12],
		ConnectionID:     ch.ConnectionID,
		UserID:           userID,
		Provider:         ch.Provider,
		ExternalUserID:   ch.ExternalUserID,
		ExternalUsername: ch.ExternalUserID,
		Role:             "operator",
		Permissions:      []string{"projects:read", "deployments:write", "logs:read"},
		Active:           true,
		PairedAt:         time.Now().UTC(),
	}

	if err := s.repo.CreateIdentity(ctx, identity); err != nil {
		return nil, fmt.Errorf("create paired identity: %w", err)
	}
	return identity, nil
}
