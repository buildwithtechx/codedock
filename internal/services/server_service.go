package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type ServerService interface {
	CreateServer(ctx context.Context, userID, name, ipAddress string) (*models.Server, error)
	ListServersByUser(ctx context.Context, userID string) ([]*models.Server, error)
	GetServer(ctx context.Context, id string) (*models.Server, error)
	DeleteServer(ctx context.Context, id string) error
}

type serverService struct {
	serverRepo repositories.ServerRepository
	userRepo   *repositories.UserRepo
}

func NewServerService(serverRepo repositories.ServerRepository, userRepo *repositories.UserRepo) ServerService {
	return &serverService{
		serverRepo: serverRepo,
		userRepo:   userRepo,
	}
}

func generateWorkerToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (s *serverService) CreateServer(ctx context.Context, userID, name, ipAddress string) (*models.Server, error) {
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

	now := time.Now().UTC()
	server := &models.Server{
		ID:          uuid.New().String(),
		UserID:      userID,
		Name:        name,
		IPAddress:   ipAddress,
		Status:      models.ServerStatusProvisioning,
		WorkerToken: generateWorkerToken(),
		CreatedAt:   now,
		UpdatedAt:   now,
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

func (s *serverService) DeleteServer(ctx context.Context, id string) error {
	return s.serverRepo.Delete(ctx, id)
}
