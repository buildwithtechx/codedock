package system

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codedock/internal/config"
	"codedock/internal/engine/observability"
	"codedock/internal/models"
	"codedock/internal/utils"
)

func (s *serverService) ListServersByUser(ctx context.Context, userID string) ([]*models.Server, error) {
	servers, err := s.serverRepo.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, srv := range servers {
		if srv.IsLocal {
			now := time.Now().UTC()
			srv.LastSeenAt = &now
			srv.Status = models.ServerStatusOnline
			srv.IsControlPlane = true
			if len(srv.Metrics) == 0 {
				srv.Metrics = observability.GetHostMetricsPayload()
			}
		}
		applyServerMetadata(srv)
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
	hasLocal := false
	for _, srv := range servers {
		if srv.IsLocal || srv.ID == "local" || srv.ID == controlPlaneServerID {
			hasLocal = true
			break
		}
	}
	if !hasLocal {
		cp := s.controlPlaneServer(userID)
		applyServerMetadata(cp)
		return append([]*models.Server{cp}, servers...), nil
	}
	return servers, nil
}

func (s *serverService) GetServer(ctx context.Context, id, userID string) (*models.Server, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil && !utils.IsNotFound(err) {
		return nil, fmt.Errorf("load server user: %w", err)
	}
	isAdmin := user != nil && (user.Role == models.UserRoleOwner || user.Role == models.UserRoleAdmin)

	if id == controlPlaneServerID || id == "local" || id == "control-plane" {
		if !isAdmin {
			return nil, utils.NewForbiddenError("server access denied")
		}
		if srv, err := s.serverRepo.GetByID(ctx, id); err == nil && srv != nil {
			if srv.IsLocal {
				now := time.Now().UTC()
				srv.LastSeenAt = &now
				srv.Status = models.ServerStatusOnline
				srv.IsControlPlane = true
				if len(srv.Metrics) == 0 {
					srv.Metrics = observability.GetHostMetricsPayload()
				}
			}
			applyServerMetadata(srv)
			return srv, nil
		}
		cp := s.controlPlaneServer(userID)
		if id == "local" || id == "control-plane" {
			cp.ID = id
		}
		applyServerMetadata(cp)
		return cp, nil
	}
	server, err := s.serverRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, utils.NewNotFoundError("Server", id)
	}
	if !isAdmin && server.UserID != userID && !server.IsLocal {
		return nil, utils.NewForbiddenError("server access denied")
	}
	if server.IsLocal {
		now := time.Now().UTC()
		server.LastSeenAt = &now
		server.Status = models.ServerStatusOnline
		if len(server.Metrics) == 0 {
			server.Metrics = observability.GetHostMetricsPayload()
		}
	}
	applyServerMetadata(server)
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
	user, _ := s.userRepo.GetUserByID(ctx, userID)
	isAdmin := user != nil && (user.Role == models.UserRoleOwner || user.Role == models.UserRoleAdmin)
	if !isAdmin && server.UserID != userID {
		return fmt.Errorf("unauthorized to delete server")
	}
	if server.Provider != "" {
		return fmt.Errorf("managed servers must be deleted through the reviewed managed lifecycle")
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
	now := time.Now().UTC()
	return &models.Server{
		ID:             controlPlaneServerID,
		UserID:         userID,
		Name:           "Codedock Control Plane",
		IPAddress:      ipAddress,
		Status:         models.ServerStatusOnline,
		IsControlPlane: true,
		IsLocal:        true,
		LastSeenAt:     &now,
		Metrics:        observability.GetHostMetricsPayload(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func applyServerMetadata(srv *models.Server) {
	if srv == nil {
		return
	}
	isCloud := config.Get().Cloud.Enabled
	isManaged := srv.Provider != "" || srv.ServerType == "managed" || isCloud
	if isManaged {
		srv.Capabilities = &models.ServerCapabilities{
			Monitor:           true,
			Terminal:          true,
			Exec:              true,
			HostConfiguration: false,
			SSH:               true,
			NetworkSettings:   true,
		}
		if srv.Managed == nil {
			tier := "pro"
			if srv.ServerType != "" {
				tier = srv.ServerType
			}
			state := "ready"
			if srv.Status != models.ServerStatusOnline && !srv.IsLocal {
				state = string(srv.Status)
			}
			srv.Managed = &models.CloudWorkspaceSummary{
				ID:                 srv.ID,
				ServerID:           srv.ID,
				Name:               srv.Name,
				PlanTierID:         tier,
				SubscriptionStatus: "active",
				ProjectCount:       0,
				State:              state,
				Resources: &models.ManagedServerResources{
					CPUCores: 4,
					MemoryMB: 8192,
					DiskMB:   80000,
				},
				CreatedAt: srv.CreatedAt.Format(time.RFC3339),
			}
		}
	} else {
		srv.Capabilities = &models.ServerCapabilities{
			Monitor:           true,
			Terminal:          true,
			Exec:              true,
			HostConfiguration: true,
			SSH:               true,
			NetworkSettings:   false,
		}
	}
}
