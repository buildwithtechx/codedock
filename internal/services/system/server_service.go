package system

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

const controlPlaneServerID = "codedock-control-plane"

type ServerService interface {
	CreateServer(ctx context.Context, userID string, req models.CreateServerRequest) (*models.Server, error)
	UpdateServer(ctx context.Context, id, userID string, req models.UpdateServerRequest) (*models.Server, error)
	TestSSH(ctx context.Context, req models.TestSSHRequest) error
	ListServersByUser(ctx context.Context, userID string) ([]*models.Server, error)
	GetServer(ctx context.Context, id, userID string) (*models.Server, error)
	DeleteServer(ctx context.Context, id, userID string) error
}

type serverService struct {
	serverRepo repositories.ServerRepository
	userRepo   *repositories.UserRepo
	sshManager *ssh.SSHManager
}

func NewServerService(serverRepo repositories.ServerRepository, userRepo *repositories.UserRepo, sshManager *ssh.SSHManager) ServerService {
	return &serverService{
		serverRepo: serverRepo,
		userRepo:   userRepo,
		sshManager: sshManager,
	}
}

func generateWorkerToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *serverService) TestSSH(ctx context.Context, req models.TestSSHRequest) error {
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

func (s *serverService) CreateServer(ctx context.Context, userID string, req models.CreateServerRequest) (*models.Server, error) {
	if err := validateServerConnection(req.IsLocal, req.SSHPort, req.SSHTransport, req.SSHJumpHost, serverPrivateKey(req.SSHKey, req.SSHPrivateKey), req.SSHPassword); err != nil {
		return nil, err
	}
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if u.PlanType != "pro" {
		servers, err := s.serverRepo.ListByUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check server limit: %w", err)
		}
		ownedServers := 0
		for _, server := range servers {
			if server.UserID == userID {
				ownedServers++
			}
		}
		if ownedServers >= 1 {
			return nil, fmt.Errorf("hobby plan is limited to 1 server. please upgrade to pro")
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
	status := models.ServerStatusOffline
	if req.IsLocal {
		status = models.ServerStatusOnline
	}

	server := &models.Server{
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
			server.Status = models.ServerStatusOnline
		}
	}

	if err := s.serverRepo.Create(ctx, server); err != nil {
		return nil, fmt.Errorf("failed to create server: %w", err)
	}

	return server, nil
}

func (s *serverService) UpdateServer(ctx context.Context, id, userID string, req models.UpdateServerRequest) (*models.Server, error) {
	if id == controlPlaneServerID {
		return nil, fmt.Errorf("the Codedock control plane cannot be modified")
	}

	server, err := s.serverRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if server == nil {
		return nil, fmt.Errorf("server not found")
	}
	if server.UserID != userID {
		return nil, fmt.Errorf("unauthorized to update server")
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
		if *req.SSHPort < 1 || *req.SSHPort > 65535 {
			return nil, fmt.Errorf("SSH port must be between 1 and 65535")
		}
		server.SSHPort = *req.SSHPort
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
		server.SSHKey = *req.SSHPrivateKey
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

	if err := validateServerConnection(server.IsLocal, server.SSHPort, server.SSHTransport, server.SSHJumpHost, serverPrivateKey(server.SSHKey, server.SSHPrivateKey), server.SSHPassword); err != nil {
		return nil, err
	}
	server.UpdatedAt = time.Now().UTC()

	if err := s.serverRepo.Update(ctx, server); err != nil {
		return nil, fmt.Errorf("failed to update server: %w", err)
	}

	if s.sshManager != nil {
		s.sshManager.RemoveClient(id)
	}
	return server, nil
}
