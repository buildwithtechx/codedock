package projects

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ProjectSettingsService struct {
	repo ProjectSettingsRepository
}

func NewProjectSettingsService(r ProjectSettingsRepository) *ProjectSettingsService {
	return &ProjectSettingsService{
		repo: r,
	}
}

func (s *ProjectSettingsService) CreateToken(ctx context.Context, t *ProjectToken) (*ProjectToken, string, error) {
	if t == nil || t.ProjectID == "" || t.Name == "" {
		return nil, "", ErrTokenInvalid
	}
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	t.CreatedAt = time.Now().UTC()

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, "", fmt.Errorf("generate token bytes: %w", err)
	}
	rawSecret := hex.EncodeToString(randomBytes)
	fullToken := fmt.Sprintf(ProjectTokenPrefixTemplate, rawSecret)
	if len(fullToken) >= 16 {
		t.TokenPrefix = fullToken[:16]
	} else {
		t.TokenPrefix = fullToken
	}

	err := s.repo.CreateToken(ctx, t, fullToken)
	if err != nil {
		return nil, "", err
	}
	return t, fullToken, nil
}

func (s *ProjectSettingsService) GetTokenByHash(ctx context.Context, tokenHash string) (*ProjectToken, error) {
	return s.repo.GetTokenByHash(ctx, tokenHash)
}

func (s *ProjectSettingsService) UpdateTokenLastUsed(ctx context.Context, id string) error {
	return s.repo.UpdateTokenLastUsed(ctx, id)
}

func (s *ProjectSettingsService) ListTokens(ctx context.Context, projectID string) ([]*ProjectToken, error) {
	if projectID == "" {
		return nil, ErrProjectIDRequired
	}
	return s.repo.ListTokensByProject(ctx, projectID)
}

func (s *ProjectSettingsService) DeleteToken(ctx context.Context, id, projectID string) error {
	if id == "" || projectID == "" {
		return ErrProjectIDRequired
	}
	return s.repo.DeleteToken(ctx, id, projectID)
}
