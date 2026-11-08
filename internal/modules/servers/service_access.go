package servers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"codedock/internal/config"
	"codedock/internal/platform/observability"
)

func (s *service) ListServersByUser(ctx context.Context, userID string, search string) ([]*Server, error) {
	servers, err := s.repo.ListByUser(ctx, userID, search)
	if err != nil {
		return nil, err
	}
	for _, srv := range servers {
		if srv.IsLocal {
			now := time.Now().UTC()
			srv.LastSeenAt = &now
			srv.Status = ServerStatusOnline
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
	var isAdmin bool
	if s.users != nil {
		if user, err := s.users.GetUserByID(ctx, userID); err == nil && user != nil {
			isAdmin = user.Role == "owner" || user.Role == "admin"
		}
	}
	if !isAdmin {
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
		return append([]*Server{cp}, servers...), nil
	}
	return servers, nil
}

func (s *service) GetServer(ctx context.Context, id, userID string) (*Server, error) {
	var isAdmin bool
	if s.users != nil {
		if user, _ := s.users.GetUserByID(ctx, userID); user != nil {
			isAdmin = user.Role == "owner" || user.Role == "admin"
		}
	}

	if id == controlPlaneServerID || id == "local" || id == "control-plane" {
		if !isAdmin {
			return nil, ErrServerUnauthorized
		}
		if srv, err := s.repo.GetByID(ctx, id); err == nil && srv != nil {
			if srv.IsLocal {
				now := time.Now().UTC()
				srv.LastSeenAt = &now
				srv.Status = ServerStatusOnline
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
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, ErrServerNotFound
	}
	if !isAdmin && server.UserID != userID && !server.IsLocal {
		return nil, ErrServerUnauthorized
	}
	if server.IsLocal {
		now := time.Now().UTC()
		server.LastSeenAt = &now
		server.Status = ServerStatusOnline
		if len(server.Metrics) == 0 {
			server.Metrics = observability.GetHostMetricsPayload()
		}
	}
	applyServerMetadata(server)
	return server, nil
}

func (s *service) DeleteServer(ctx context.Context, id, userID string) error {
	if id == controlPlaneServerID {
		return fmt.Errorf("the Codedock control plane cannot be deleted")
	}
	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if server == nil {
		return ErrServerNotFound
	}
	var isAdmin bool
	if s.users != nil {
		if user, _ := s.users.GetUserByID(ctx, userID); user != nil {
			isAdmin = user.Role == "owner" || user.Role == "admin"
		}
	}
	if !isAdmin && server.UserID != userID {
		return ErrServerUnauthorized
	}
	if server.Provider != "" {
		return fmt.Errorf("managed servers must be deleted through the reviewed managed lifecycle")
	}

	if s.sshManager != nil {
		s.sshManager.RemoveClient(id)
	}
	return s.repo.Delete(ctx, id)
}

func (s *service) controlPlaneServer(userID string) *Server {
	ipAddress := config.Get().Server.HostIP
	if ipAddress == "" {
		ipAddress = "127.0.0.1"
	}
	now := time.Now().UTC()
	return &Server{
		ID:             controlPlaneServerID,
		UserID:         userID,
		Name:           "Codedock Control Plane",
		IPAddress:      ipAddress,
		Status:         ServerStatusOnline,
		IsControlPlane: true,
		IsLocal:        true,
		LastSeenAt:     &now,
		Metrics:        observability.GetHostMetricsPayload(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

func applyServerMetadata(srv *Server) {
	if srv == nil {
		return
	}
	isCloud := config.Get().Cloud.Enabled
	isManaged := srv.Provider != "" || srv.ServerType == "managed" || isCloud
	buildCacheBudgetSupported := !srv.IsLocal && !srv.IsControlPlane && srv.ID != controlPlaneServerID
	if isManaged {
		srv.Capabilities = &ServerCapabilities{
			Monitor:           true,
			Terminal:          true,
			Exec:              true,
			HostConfiguration: false,
			SSH:               true,
			NetworkSettings:   true,
			BuildCacheBudget:  buildCacheBudgetSupported,
		}
		if srv.Managed == nil {
			tier := "pro"
			if srv.ServerType != "" {
				tier = srv.ServerType
			}
			state := "ready"
			if srv.Status != ServerStatusOnline && !srv.IsLocal {
				state = string(srv.Status)
			}
			srv.Managed = &CloudWorkspaceSummary{
				ID:                 srv.ID,
				ServerID:           srv.ID,
				Name:               srv.Name,
				PlanTierID:         tier,
				SubscriptionStatus: "active",
				ProjectCount:       0,
				State:              state,
				Resources: &ManagedServerResources{
					CPUCores: 4,
					MemoryMB: 8192,
					DiskMB:   80000,
				},
				CreatedAt: srv.CreatedAt.Format(time.RFC3339),
			}
		}
	} else {
		srv.Capabilities = &ServerCapabilities{
			Monitor:           true,
			Terminal:          true,
			Exec:              true,
			HostConfiguration: true,
			SSH:               true,
			NetworkSettings:   false,
			BuildCacheBudget:  buildCacheBudgetSupported,
		}
	}
}
