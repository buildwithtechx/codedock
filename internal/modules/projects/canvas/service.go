package canvas

import (
	"context"
	"fmt"
)

type CanvasService struct {
	repo    CanvasRepository
	runtime interface {
		Observe(context.Context, *EnvironmentCanvas)
	}
}

func NewCanvasService(r CanvasRepository) *CanvasService {
	return &CanvasService{repo: r}
}

func (s *CanvasService) ListSummaries(ctx context.Context, organizationID string) ([]CanvasSummary, error) {
	if organizationID == "" {
		return nil, fmt.Errorf("organization is required")
	}
	return s.repo.ListCanvasSummaries(ctx, organizationID)
}

func (s *CanvasService) GetSummary(ctx context.Context, id string) (*CanvasSummary, error) {
	if id == "" {
		return nil, fmt.Errorf("canvas id is required")
	}
	return s.repo.GetCanvasSummary(ctx, id)
}

func (s *CanvasService) GetEnvironmentCanvas(ctx context.Context, id string) (*EnvironmentCanvas, error) {
	if id == "" {
		return nil, fmt.Errorf("environment id is required")
	}
	canvas, err := s.repo.GetEnvironmentCanvas(ctx, id)
	if err != nil {
		return nil, err
	}
	if s.runtime != nil && canvas != nil {
		s.runtime.Observe(ctx, canvas)
	}
	return canvas, nil
}

func (s *CanvasService) SetRuntime(runtime interface {
	Observe(context.Context, *EnvironmentCanvas)
}) {
	s.runtime = runtime
}

func (s *CanvasService) ApplyTopology(ctx context.Context, environment string, request TopologyApplyRequest) error {
	repo, ok := s.repo.(interface {
		ApplyTopology(context.Context, string, TopologyApplyRequest) error
	})
	if !ok {
		return fmt.Errorf("topology modification unavailable")
	}
	return repo.ApplyTopology(ctx, environment, request)
}

func (s *CanvasService) ReadObservation(ctx context.Context, canvas *EnvironmentCanvas, node string) (*RuntimeObservation, error) {
	runtime, ok := s.runtime.(interface {
		Read(context.Context, *EnvironmentCanvas, string) (*RuntimeObservation, error)
	})
	if !ok {
		return nil, fmt.Errorf("runtime observation unavailable")
	}
	return runtime.Read(ctx, canvas, node)
}
