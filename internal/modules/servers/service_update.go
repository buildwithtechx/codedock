package servers

import (
	"context"
	"fmt"
	"time"
)

func (s *service) UpdateServer(ctx context.Context, id, userID string, req UpdateServerRequest) (*Server, error) {
	if id == controlPlaneServerID {
		return nil, fmt.Errorf("the Codedock control plane cannot be modified")
	}

	server, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, ErrServerNotFound
	}
	var isAdmin bool
	if s.users != nil {
		if u, _ := s.users.GetUserByID(ctx, userID); u != nil {
			isAdmin = u.Role == "owner" || u.Role == "admin"
		}
	}
	if !isAdmin && server.UserID != userID {
		return nil, ErrServerUnauthorized
	}

	if req.Name != nil && *req.Name != "" {
		server.Name = *req.Name
	}
	if req.IPAddress != nil {
		server.IPAddress = *req.IPAddress
	}
	if req.IsLocal != nil {
		server.IsLocal = *req.IsLocal
	}
	if req.SSHHost != nil {
		server.SSHHost = *req.SSHHost
	}
	if req.SSHPort != nil {
		if *req.SSHPort < 0 || *req.SSHPort > 65535 {
			return nil, fmt.Errorf("SSH port must be between 1 and 65535")
		}
		if *req.SSHPort != 0 {
			server.SSHPort = *req.SSHPort
		}
	}
	if req.SSHUser != nil && *req.SSHUser != "" {
		server.SSHUser = *req.SSHUser
	}
	if req.SSHAuthMethod != nil {
		server.SSHAuthMethod = *req.SSHAuthMethod
	}
	if req.SSHKey != nil {
		server.SSHKey = *req.SSHKey
	}
	if req.SSHPrivateKey != nil {
		server.SSHPrivateKey = *req.SSHPrivateKey
		if req.SSHKey == nil {
			server.SSHKey = *req.SSHPrivateKey
		}
	}
	if req.SSHPassword != nil {
		server.SSHPassword = *req.SSHPassword
	}
	if req.SSHTransport != nil {
		server.SSHTransport = *req.SSHTransport
	}
	if req.SSHJumpHost != nil {
		server.SSHJumpHost = *req.SSHJumpHost
	}
	if req.BuildCacheBudgetBytes != nil {
		if server.IsLocal || server.IsControlPlane || server.ID == controlPlaneServerID {
			return nil, fmt.Errorf("build-cache budgets apply to remote deployment servers")
		}
		if *req.BuildCacheBudgetBytes <= 0 {
			server.BuildCacheBudgetBytes = nil
		} else {
			const minBudget = int64(1073741824)
			const maxBudget = int64(70368744177664)
			if *req.BuildCacheBudgetBytes < minBudget || *req.BuildCacheBudgetBytes > maxBudget {
				return nil, fmt.Errorf("build-cache budget must be between 1 GiB and 64 TiB")
			}
			server.BuildCacheBudgetBytes = req.BuildCacheBudgetBytes
		}
	}

	if err := validateServerConnection(server.IsLocal, server.SSHPort, server.SSHTransport, server.SSHJumpHost, serverPrivateKey(server.SSHKey, server.SSHPrivateKey), server.SSHPassword); err != nil {
		return nil, err
	}
	server.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, server); err != nil {
		return nil, fmt.Errorf("failed to update server: %w", err)
	}

	connChanged := req.SSHHost != nil || req.SSHPort != nil || req.SSHUser != nil || req.SSHKey != nil || req.SSHPassword != nil
	if s.sshManager != nil && connChanged {
		s.sshManager.RemoveClient(id)
	}
	applyServerMetadata(server)
	return server, nil
}
