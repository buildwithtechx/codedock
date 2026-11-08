package chatapps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"codedock/internal/models"
	"codedock/internal/repositories"
	"codedock/internal/utils"

	"github.com/google/uuid"
)

type Service struct {
	repo  repositories.ChatAppRepository
	vault *utils.Vault
}

func NewService(repo repositories.ChatAppRepository, vault *utils.Vault) *Service {
	return &Service{
		repo:  repo,
		vault: vault,
	}
}

func (s *Service) Repo() repositories.ChatAppRepository {
	return s.repo
}

func (s *Service) CreateConnection(ctx context.Context, userID string, req models.CreateChatAppRequest) (*models.ChatAppConnection, error) {
	if req.Provider == "" {
		return nil, errors.New("provider is required")
	}
	if req.Transport == "" {
		req.Transport = models.ChatAppTransportWebhook
	}
	if req.Purpose == "" {
		req.Purpose = models.ChatAppPurposeBoth
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

	conn := &models.ChatAppConnection{
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
		Status:               models.ChatAppStatusConnected,
	}

	if err := s.repo.CreateConnection(ctx, conn); err != nil {
		return nil, fmt.Errorf("create chat connection: %w", err)
	}
	return conn, nil
}

func (s *Service) GetConnection(ctx context.Context, id string) (*models.ChatAppConnection, error) {
	conn, err := s.repo.GetConnection(ctx, id)
	if err != nil {
		return nil, err
	}
	if conn == nil {
		return nil, errors.New("connection not found")
	}
	return conn, nil
}

func (s *Service) ListConnections(ctx context.Context, userID string) ([]*models.ChatAppConnection, error) {
	return s.repo.ListConnectionsByUser(ctx, userID)
}

func (s *Service) UpdateConnection(ctx context.Context, id string, req models.UpdateChatAppRequest) (*models.ChatAppConnection, error) {
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
	if req.Credentials != nil {
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
		return nil, fmt.Errorf("update connection: %w", err)
	}
	return conn, nil
}

func (s *Service) DeleteConnection(ctx context.Context, id string) error {
	return s.repo.DeleteConnection(ctx, id)
}

func (s *Service) GetDecryptedCredentials(ctx context.Context, conn *models.ChatAppConnection) (map[string]string, error) {
	if conn.CredentialsEncrypted == "" {
		return make(map[string]string), nil
	}
	raw := conn.CredentialsEncrypted
	if s.vault != nil {
		dec, err := s.vault.Decrypt(raw)
		if err == nil {
			raw = dec
		}
	}
	var creds map[string]string
	if err := json.Unmarshal([]byte(raw), &creds); err != nil {
		return nil, fmt.Errorf("unmarshal credentials: %w", err)
	}
	return creds, nil
}

func (s *Service) ListIdentities(ctx context.Context, connectionID string) ([]*models.ChatAppIdentity, error) {
	return s.repo.ListIdentitiesByConnection(ctx, connectionID)
}

func (s *Service) DeleteIdentity(ctx context.Context, id string) error {
	return s.repo.DeleteIdentity(ctx, id)
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
	case models.ChatAppProviderSlack:
		hasKey = creds["bot_token"] != "" || creds["app_token"] != "" || creds["signing_secret"] != ""
	case models.ChatAppProviderDiscord:
		hasKey = creds["bot_token"] != "" || creds["public_key"] != ""
	case models.ChatAppProviderTelegram:
		hasKey = creds["bot_token"] != ""
	case models.ChatAppProviderTeams:
		hasKey = creds["app_id"] != "" || creds["app_password"] != ""
	}

	now := time.Now().UTC()
	conn.LastEventAt = &now
	if hasKey {
		conn.Status = models.ChatAppStatusConnected
		conn.LastError = nil
	} else {
		conn.Status = models.ChatAppStatusPairingRequired
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
