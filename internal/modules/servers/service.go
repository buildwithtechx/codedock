package servers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock/internal/config"
	"codedock/internal/platform/ssh"
)

const controlPlaneServerID = "codedock-control-plane"

type UserInfo struct {
	ID       string
	Role     string
	PlanType string
}

type UserLookup interface {
	GetUserByID(ctx context.Context, id string) (*UserInfo, error)
}

type SSHManager interface {
	TestConnection(ctx context.Context, cfg ssh.Config) error
	RemoveClient(serverID string)
}

type Service interface {
	CreateServer(ctx context.Context, userID string, req CreateServerRequest) (*Server, error)
	UpdateServer(ctx context.Context, id, userID string, req UpdateServerRequest) (*Server, error)
	TestSSH(ctx context.Context, req TestSSHRequest) error
	ListServersByUser(ctx context.Context, userID string, search string) ([]*Server, error)
	GetServer(ctx context.Context, id, userID string) (*Server, error)
	DeleteServer(ctx context.Context, id, userID string) error
}

type service struct {
	repo       Repository
	users      UserLookup
	sshManager SSHManager
}

func NewService(repo Repository, users UserLookup, sshManager SSHManager) Service {
	return &service{
		repo:       repo,
		users:      users,
		sshManager: sshManager,
	}
}

func generateWorkerToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *service) TestSSH(ctx context.Context, req TestSSHRequest) error {
	if s.sshManager == nil {
		return fmt.Errorf("ssh manager is not initialized")
	}

	host := req.SSHHost
	if host == "" {
		return fmt.Errorf("sshHost is required")
	}
	if req.SSHPort < 0 || req.SSHPort > 65535 {
		return fmt.Errorf("SSH port must be between 1 and 65535")
	}
	port := req.SSHPort
	if port <= 0 {
		port = 22
	}
	user := req.SSHUser
	if user == "" {
		user = "root"
	}
	key := req.SSHKey
	if key == "" && req.SSHPrivateKey != "" {
		key = req.SSHPrivateKey
	}

	return s.sshManager.TestConnection(ctx, ssh.Config{
		Host:     host,
		Port:     port,
		User:     user,
		Key:      key,
		Password: req.SSHPassword,
	})
}

func (s *service) CreateServer(ctx context.Context, userID string, req CreateServerRequest) (*Server, error) {
	if err := validateServerConnection(req.IsLocal, req.SSHPort, req.SSHTransport, req.SSHJumpHost, serverPrivateKey(req.SSHKey, req.SSHPrivateKey), req.SSHPassword); err != nil {
		return nil, err
	}
	if s.users != nil {
		u, err := s.users.GetUserByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user: %w", err)
		}
		if u != nil && config.Get().Cloud.Enabled && u.PlanType != "pro" {
			servers, err := s.repo.ListByUser(ctx, userID, "")
			if err != nil {
				return nil, fmt.Errorf("failed to check server limit: %w", err)
			}
			owned := 0
			for _, srv := range servers {
				if srv.UserID == userID {
					owned++
				}
			}
			if owned >= 1 {
				return nil, fmt.Errorf("hobby plan is limited to 1 server. please upgrade to pro")
			}
		}
	}

	key := req.SSHKey
	if key == "" && req.SSHPrivateKey != "" {
		key = req.SSHPrivateKey
	}

	sshAuthMethod := req.SSHAuthMethod
	if sshAuthMethod == "" {
		if req.SSHPassword != "" {
			sshAuthMethod = "password"
		} else {
			sshAuthMethod = "key"
		}
	}

	sshTransport := req.SSHTransport
	if sshTransport == "" {
		sshTransport = "direct"
	}

	ipAddress := req.IPAddress
	sshHost := req.SSHHost
	if req.IsLocal {
		if ipAddress == "" {
			ipAddress = "127.0.0.1"
		}
		sshHost = ipAddress
	} else if sshHost == "" {
		sshHost = ipAddress
	}

	sshPort := req.SSHPort
	if sshPort <= 0 {
		sshPort = 22
	}
	sshUser := req.SSHUser
	if sshUser == "" {
		sshUser = "root"
	}

	now := time.Now().UTC()
	status := ServerStatusOffline
	if req.IsLocal {
		status = ServerStatusOnline
	}

	server := &Server{
		ID:            uuid.New().String(),
		UserID:        userID,
		Name:          req.Name,
		IPAddress:     ipAddress,
		IsLocal:       req.IsLocal,
		SSHHost:       sshHost,
		SSHPort:       sshPort,
		SSHUser:       sshUser,
		SSHAuthMethod: sshAuthMethod,
		SSHKey:        key,
		SSHPrivateKey: key,
		SSHPassword:   req.SSHPassword,
		SSHTransport:  sshTransport,
		SSHJumpHost:   req.SSHJumpHost,
		Status:        status,
		WorkerToken:   generateWorkerToken(),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if !req.IsLocal && s.sshManager != nil {
		if err := s.sshManager.TestConnection(ctx, ssh.Config{
			Host:     server.SSHHost,
			Port:     server.SSHPort,
			User:     server.SSHUser,
			Key:      server.SSHKey,
			Password: server.SSHPassword,
		}); err == nil {
			server.Status = ServerStatusOnline
		}
	}

	if err := s.repo.Create(ctx, server); err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	return server, nil
}

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

func validateServerConnection(local bool, port int, transport, jumpHost, key, password string) error {
	if local {
		return nil
	}
	if port < 0 || port > 65535 {
		return fmt.Errorf("SSH port must be between 1 and 65535")
	}
	if transport == "jump" && jumpHost == "" {
		return fmt.Errorf("jump host is required when using jump transport")
	}
	if key == "" && password == "" {
		return fmt.Errorf("SSH key or password is required for remote servers")
	}
	return nil
}

func serverPrivateKey(key, privateKey string) string {
	if key != "" {
		return key
	}
	return privateKey
}
