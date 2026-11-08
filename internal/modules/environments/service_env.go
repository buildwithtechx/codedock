package environments

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type EnvironmentService struct {
	envRepo EnvironmentRepository
	varRepo ProjectVarRepository
}

func NewEnvironmentService(er EnvironmentRepository, vr ProjectVarRepository) *EnvironmentService {
	return &EnvironmentService{
		envRepo: er,
		varRepo: vr,
	}
}

func (s *EnvironmentService) CreateEnvironment(ctx context.Context, env *EnvironmentConfig) (*EnvironmentConfig, error) {
	if env == nil || env.ProjectID == "" || env.Name == "" {
		return nil, ErrInvalidEnvironment
	}
	if env.ID == "" {
		env.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	if env.CreatedAt.IsZero() {
		env.CreatedAt = now
	}
	if env.UpdatedAt.IsZero() {
		env.UpdatedAt = now
	}
	if err := s.envRepo.Create(ctx, env); err != nil {
		return nil, err
	}
	return env, nil
}

func (s *EnvironmentService) GetEnvironment(ctx context.Context, id string) (*EnvironmentConfig, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.envRepo.Get(ctx, id)
}

func (s *EnvironmentService) ListByProject(ctx context.Context, projectID string) ([]EnvironmentConfig, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	return s.envRepo.ListByProject(ctx, projectID)
}

func (s *EnvironmentService) DeleteEnvironment(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.envRepo.Delete(ctx, id)
}

func (s *EnvironmentService) GetVars(ctx context.Context, projectID string) (map[string]string, error) {
	if projectID == "" {
		return nil, errors.New("project id required")
	}
	if s.varRepo == nil {
		return make(map[string]string), nil
	}
	return s.varRepo.GetVars(ctx, projectID)
}

func (s *EnvironmentService) SetVar(ctx context.Context, projectID, key, value string) error {
	if projectID == "" || key == "" {
		return errors.New("project id and key required")
	}
	if s.varRepo == nil {
		return nil
	}
	return s.varRepo.SetVar(ctx, projectID, key, value)
}
