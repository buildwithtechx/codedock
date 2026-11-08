package deployments

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

type RegistryService struct {
	repo RegistryRepository
}

func NewRegistryService(repo RegistryRepository) *RegistryService {
	return &RegistryService{repo: repo}
}

func (s *RegistryService) CreateRegistry(ctx context.Context, projectID string, req CreateRegistryRequest) (*Registry, error) {
	if projectID == "" {
		return nil, ErrProjectRequired
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.New("registry name is required")
	}
	url := strings.TrimSpace(req.RegistryURL)
	if url == "" {
		return nil, errors.New("registry URL is required")
	}
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, errors.New("username is required")
	}

	reg := &Registry{
		ID:            uuid.NewString(),
		ProjectID:     projectID,
		Name:          name,
		RegistryURL:   url,
		Username:      username,
		PasswordToken: req.PasswordToken,
	}

	if err := s.repo.Create(ctx, reg); err != nil {
		return nil, err
	}
	reg.PasswordToken = ""
	return reg, nil
}

func (s *RegistryService) ListRegistries(ctx context.Context, projectID string) ([]*Registry, error) {
	if projectID == "" {
		return nil, ErrProjectRequired
	}
	list, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	for _, r := range list {
		r.PasswordToken = ""
	}
	return list, nil
}

func (s *RegistryService) DeleteRegistry(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("registry id required")
	}
	return s.repo.Delete(ctx, id)
}

func (s *RegistryService) GetRegistry(ctx context.Context, id string) (*Registry, error) {
	if id == "" {
		return nil, errors.New("registry id required")
	}
	reg, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if reg != nil {
		reg.PasswordToken = ""
	}
	return reg, nil
}
