package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"crypto/sha256"

	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(ur repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

func (s *UserService) CountUsers(ctx context.Context) (int, error) {
	return s.userRepo.CountUsers(ctx)
}

func (s *UserService) CreateUser(ctx context.Context, u *models.User) error {
	if u.Email == "" {
		return errors.New("user email is required")
	}
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	u.UpdatedAt = now
	return s.userRepo.CreateUser(ctx, u)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}
	return s.userRepo.GetUserByEmail(ctx, email)
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("user id is required")
	}
	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]models.User, int, error) {
	return s.userRepo.ListUsers(ctx, limit, offset)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.DeleteUser(ctx, id)
}

func (s *UserService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return utils.NewNotFoundError("User", userID)
	}

	if !utils.CheckPasswordHash(oldPassword, u.PasswordHash) {
		return errors.New("invalid old password")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	u.PasswordHash = string(hashed)
	return s.userRepo.UpdateUser(ctx, u)
}

func (s *UserService) UpdateUser(ctx context.Context, u *models.User) error {
	if u == nil || u.ID == "" {
		return errors.New("valid user required for update")
	}
	u.UpdatedAt = time.Now()
	return s.userRepo.UpdateUser(ctx, u)
}

func (s *UserService) CreatePAT(ctx context.Context, userID, name string, accessLevel string, projectScope string, allowedProjects []string, expiresAt *time.Time) (*models.PersonalAccessToken, string, error) {
	if userID == "" || name == "" {
		return nil, "", errors.New("userId and name are required")
	}
	if accessLevel == "" {
		accessLevel = "read_write"
	}
	if projectScope == "" {
		projectScope = "all"
	}
	if accessLevel != "read" && accessLevel != "read_write" {
		return nil, "", errors.New("access level must be read or read_write")
	}
	if projectScope != "all" && projectScope != "specific" {
		return nil, "", errors.New("project scope must be all or specific")
	}
	if projectScope == "specific" && len(allowedProjects) == 0 {
		return nil, "", errors.New("specific project scope requires allowed projects")
	}
	for _, id := range allowedProjects {
		if strings.TrimSpace(id) == "" {
			return nil, "", errors.New("allowed project id cannot be empty")
		}
	}
	if expiresAt != nil && !expiresAt.After(time.Now()) {
		return nil, "", errors.New("expiry must be in the future")
	}
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, "", fmt.Errorf("generate personal token: %w", err)
	}
	rawToken := "vpt_" + hex.EncodeToString(bytes)

	hasher := sha256.New()
	hasher.Write([]byte(rawToken))
	hash := hex.EncodeToString(hasher.Sum(nil))

	var allowedProjectsJSON *string
	if projectScope == "specific" && len(allowedProjects) > 0 {
		if j, err := json.Marshal(allowedProjects); err == nil {
			s := string(j)
			allowedProjectsJSON = &s
		}
	}

	pat := &models.PersonalAccessToken{
		ID:              uuid.New().String(),
		UserID:          userID,
		Name:            name,
		TokenHash:       hash,
		Prefix:          rawToken[:8],
		AccessLevel:     accessLevel,
		ProjectScope:    projectScope,
		AllowedProjects: allowedProjectsJSON,
		CreatedAt:       time.Now(),
	}
	if expiresAt != nil {
		pat.ExpiresAt = expiresAt
	}
	if err := s.userRepo.CreatePAT(ctx, pat); err != nil {
		return nil, "", err
	}
	return pat, rawToken, nil
}

func (s *UserService) ListPATs(ctx context.Context, userID string) ([]*models.PersonalAccessToken, error) {
	if userID == "" {
		return nil, errors.New("userId is required")
	}
	return s.userRepo.ListPATs(ctx, userID)
}

func (s *UserService) DeletePAT(ctx context.Context, id, userID string) error {
	if id == "" || userID == "" {
		return errors.New("id and userId are required")
	}
	return s.userRepo.DeletePAT(ctx, id, userID)
}
