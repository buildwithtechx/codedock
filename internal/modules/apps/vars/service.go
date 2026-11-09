package vars

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type VariableService interface {
	CreateVariable(ctx context.Context, v *Variable) (*Variable, error)
	UpdateVariable(ctx context.Context, v *Variable) error
	GetVariable(ctx context.Context, id string) (*Variable, error)
	ListVariablesByService(ctx context.Context, serviceID string) ([]*Variable, error)
	DeleteVariable(ctx context.Context, id string) error
}

type Service struct {
	repo ServiceVarRepository
}

func NewService(repo ServiceVarRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateVariable(ctx context.Context, v *Variable) (*Variable, error) {
	if v.ServiceID == "" || v.Key == "" {
		return nil, errors.New("serviceId and key required")
	}
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	v.CreatedAt = now
	v.UpdatedAt = now
	if err := s.repo.Create(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *Service) UpdateVariable(ctx context.Context, v *Variable) error {
	if v.ID == "" {
		return errors.New("variable ID required")
	}
	v.UpdatedAt = time.Now().UTC()
	return s.repo.Update(ctx, v)
}

func (s *Service) GetVariable(ctx context.Context, id string) (*Variable, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListVariablesByService(ctx context.Context, serviceID string) ([]*Variable, error) {
	if serviceID == "" {
		return nil, errors.New("service ID required")
	}
	return s.repo.ListByService(ctx, serviceID)
}

func (s *Service) DeleteVariable(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	return s.repo.Delete(ctx, id)
}
