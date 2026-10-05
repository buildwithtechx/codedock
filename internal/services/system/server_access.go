package system

import (
	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"fmt"
	"strings"
	"time"
)

func (s *serverService) ListServersByUser(ctx context.Context, userID string) ([]*models.Server, error) {
	servers, err := s.serverRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(userID, "api-token-") {
		return servers, nil
	}
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("load server user: %w", err)
	}
	if user == nil {
		return nil, utils.NewNotFoundError("User", userID)
	}
	if user.Role != models.UserRoleOwner && user.Role != models.UserRoleAdmin {
		return servers, nil
	}
	return append([]*models.Server{s.controlPlaneServer(userID)}, servers...), nil
}

func (s *serverService) GetServer(ctx context.Context, id, userID string) (*models.Server, error) {
	if id == controlPlaneServerID {
		user, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("load server user: %w", err)
		}
		if user == nil || (user.Role != models.UserRoleOwner && user.Role != models.UserRoleAdmin) {
			return nil, utils.NewForbiddenError("server access denied")
		}
		return s.controlPlaneServer(userID), nil
	}
	server, err := s.serverRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if server != nil && server.UserID != userID && !(server.IsLocal && server.UserID == "system") {
		return nil, utils.NewForbiddenError("server access denied")
	}
	return server, nil
}

func (s *serverService) DeleteServer(ctx context.Context, id, userID string) error {
	if id == controlPlaneServerID {
		return fmt.Errorf("the Codedock control plane cannot be deleted")
	}
	server, err := s.serverRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found")
	}
	if server.UserID != userID {
		return fmt.Errorf("unauthorized to delete server")
	}

	if s.sshManager != nil {
		s.sshManager.RemoveClient(id)
	}
	return s.serverRepo.Delete(ctx, id)
}

func (s *serverService) controlPlaneServer(userID string) *models.Server {
	ipAddress := config.Get().Server.HostIP
	if ipAddress == "" {
		ipAddress = "127.0.0.1"
	}
	return &models.Server{
		ID:             controlPlaneServerID,
		UserID:         userID,
		Name:           "Codedock Control Plane",
		IPAddress:      ipAddress,
		Status:         models.ServerStatusOnline,
		IsControlPlane: true,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}
}
