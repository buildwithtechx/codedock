package git

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ProviderService struct {
	repo       GitProviderRepository
	httpClient *http.Client
}

func NewProviderService(r GitProviderRepository) *ProviderService {
	return &ProviderService{
		repo:       r,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *ProviderService) SaveProvider(ctx context.Context, gp *GitProviderConfig) error {
	if gp == nil || gp.UserID == "" || gp.Provider == "" {
		return errors.New("valid git provider config with userId and provider required")
	}
	if gp.ID == "" {
		gp.ID = uuid.New().String()
	}
	gp.UpdatedAt = time.Now().UTC()
	if gp.CreatedAt.IsZero() {
		gp.CreatedAt = gp.UpdatedAt
	}
	return s.repo.SaveProvider(ctx, gp)
}

func (s *ProviderService) ConnectProvider(ctx context.Context, userID string, req *GitConnectRequest) (*GitProviderConfig, error) {
	switch req.Provider {
	case ProviderGitHub, ProviderGitLab, ProviderBitbucket, ProviderGitea:
	default:
		return nil, ErrUnsupportedProvider
	}
	if req.AccessToken == "" {
		return nil, ErrAccessTokenRequired
	}
	gp := &GitProviderConfig{
		UserID:      userID,
		Provider:    req.Provider,
		AccessToken: req.AccessToken,
		AccountName: req.AccountName,
	}
	if err := s.SaveProvider(ctx, gp); err != nil {
		return nil, fmt.Errorf("failed to save git provider: %w", err)
	}
	gp.AccessToken = ""
	return gp, nil
}

func (s *ProviderService) GetProvider(ctx context.Context, userID, provider string) (*GitProviderConfig, error) {
	if userID == "" || provider == "" {
		return nil, errors.New("userId and provider required")
	}
	return s.repo.GetProvider(ctx, userID, provider)
}

func (s *ProviderService) GetConnectedProviders(ctx context.Context, userID string) ([]ConnectedProviderStatus, error) {
	providers, err := s.repo.ListProvidersByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	providerMap := make(map[string]*GitProviderConfig)
	for _, gp := range providers {
		providerMap[gp.Provider] = gp
	}
	var results []ConnectedProviderStatus
	for _, provider := range []string{ProviderGitHub, ProviderGitLab, ProviderBitbucket, ProviderGitea} {
		if gp, ok := providerMap[provider]; ok && gp != nil {
			results = append(results, ConnectedProviderStatus{
				Provider:    provider,
				Connected:   true,
				AccountName: gp.AccountName,
				UpdatedAt:   &gp.UpdatedAt,
			})
		} else {
			results = append(results, ConnectedProviderStatus{
				Provider:  provider,
				Connected: false,
			})
		}
	}
	return results, nil
}

func (s *ProviderService) GetAnyProviderByType(ctx context.Context, provider string) (*GitProviderConfig, error) {
	if provider == "" {
		return nil, errors.New("provider required")
	}
	return s.repo.GetAnyProviderByType(ctx, provider)
}

func (s *ProviderService) ListProvidersByUser(ctx context.Context, userID string) ([]*GitProviderConfig, error) {
	if userID == "" {
		return nil, ErrUserRequired
	}
	return s.repo.ListProvidersByUser(ctx, userID)
}

func (s *ProviderService) DisconnectProvider(ctx context.Context, userID, provider string) error {
	if userID == "" || provider == "" {
		return errors.New("userId and provider required")
	}
	return s.repo.DeleteProvider(ctx, userID, provider)
}
