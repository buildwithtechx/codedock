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

type ServerService interface {
	CreateServer(ctx context.Context, userID string, req models.CreateServerRequest) (*models.Server, error)
	TestSSH(ctx context.Context, req models.TestSSHRequest) error
	ListServersByUser(ctx context.Context, userID string) ([]*models.Server, error)
	GetServer(ctx context.Context, id string) (*models.Server, error)
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
	port := req.SSHPort
	if port <= 0 {
		port = 22
	}
	user := req.SSHUser
	if user == "" {
		user = "root"
	}

	return s.sshManager.TestConnection(ctx, ssh.Config{
		Host:     host,
		Port:     port,
		User:     user,
		Key:      req.SSHKey,
		Password: req.SSHPassword,
	})
}

func (s *serverService) CreateServer(ctx context.Context, userID string, req models.CreateServerRequest) (*models.Server, error) {
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if u.PlanType != "pro" {
		servers, err := s.serverRepo.ListByUser(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check server limit: %w", err)
		}
		if len(servers) >= 1 {
			return nil, fmt.Errorf("hobby plan is limited to 1 server. please upgrade to pro")
		}
	}

	sshHost := req.SSHHost
	if sshHost == "" {
		sshHost = req.IPAddress
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
	server := &models.Server{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        req.Name,
		IPAddress:   req.IPAddress,
		SSHHost:     sshHost,
		SSHPort:     sshPort,
		SSHUser:     sshUser,
		SSHKey:      req.SSHKey,
		SSHPassword: req.SSHPassword,
		Status:      models.ServerStatusOffline,
		WorkerToken: generateWorkerToken(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if s.sshManager != nil {
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

func (s *serverService) ListServersByUser(ctx context.Context, userID string) ([]*models.Server, error) {
	return s.serverRepo.ListByUser(ctx, userID)
}

func (s *serverService) GetServer(ctx context.Context, id string) (*models.Server, error) {
	return s.serverRepo.GetByID(ctx, id)
}

func (s *serverService) DeleteServer(ctx context.Context, id, userID string) error {
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
