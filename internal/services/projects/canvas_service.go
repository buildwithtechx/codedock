package projects

import (
	"context"
	"errors"
	"fmt"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type CanvasService struct {
	repo    repositories.CanvasRepository
	runtime interface {
		Observe(context.Context, *models.EnvironmentCanvas)
	}
}

func NewCanvasService(r repositories.CanvasRepository) *CanvasService {
	return &CanvasService{repo: r}
}

func (s *CanvasService) ListSummaries(ctx context.Context, organizationID string) ([]models.CanvasSummary, error) {
	if organizationID == "" {
		return nil, errors.New("organization id required")
	}
	return s.repo.ListCanvasSummaries(ctx, organizationID)
}

func (s *CanvasService) GetSummary(ctx context.Context, id string) (*models.CanvasSummary, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetCanvasSummary(ctx, id)
}

func (s *CanvasService) GetEnvironmentCanvas(ctx context.Context, id string) (*models.EnvironmentCanvas, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	canvas, err := s.repo.GetEnvironmentCanvas(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.runtime != nil {
		s.runtime.Observe(ctx, canvas)
	}
	return canvas, nil
}

func (s *CanvasService) SetRuntime(runtime interface {
	Observe(context.Context, *models.EnvironmentCanvas)
}) {
	s.runtime = runtime
}

func (s *CanvasService) ApplyTopology(ctx context.Context, environment string, request models.TopologyApplyRequest) error {
	repo, ok := s.repo.(interface {
		ApplyTopology(context.Context, string, models.TopologyApplyRequest) error
	})
	if !ok {
		return fmt.Errorf("topology persistence unavailable")
	}
	return repo.ApplyTopology(ctx, environment, request)
}

func (s *CanvasService) ReadObservation(ctx context.Context, canvas *models.EnvironmentCanvas, node string) (*models.RuntimeObservation, error) {
	runtime, ok := s.runtime.(interface {
		Read(context.Context, *models.EnvironmentCanvas, string) (*models.RuntimeObservation, error)
	})
	if !ok {
		return nil, fmt.Errorf("runtime observation unavailable")
	}
	return runtime.Read(ctx, canvas, node)
}
