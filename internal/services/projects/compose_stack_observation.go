package projects

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (s *ComposeStackService) Get(ctx context.Context, project, id string) (*models.ComposeStack, error) {
	return s.store.Get(ctx, project, id)
}

func (s *ComposeStackService) ListObserved(ctx context.Context, project string) ([]models.ComposeStack, error) {
	stacks, err := s.store.List(ctx, project)
	if err != nil {
		return nil, err
	}
	for index := range stacks {
		stack := &stacks[index]
		results, err := s.runtime.Results(ctx, stack.ID)
		if err != nil {
			stack.Error = "Runtime observation unavailable: " + err.Error()
			if stack.Status == "READY" {
				stack.Status = "OBSERVATION_FAILED"
			}
			continue
		}
		data, err := json.Marshal(results)
		if err != nil {
			return nil, fmt.Errorf("encode stack observation: %w", err)
		}
		stack.Results = string(data)
		if stack.Status == "READY" {
			for _, result := range results {
				if result.State != "running" && !(result.State == "exited" && result.ExitCode == 0) || result.Health == "unhealthy" {
					stack.Status = "DEGRADED"
				}
			}
			if len(results) == 0 {
				stack.Status = "DEGRADED"
				stack.Error = "No stack containers are observed."
			}
		}
	}
	return stacks, nil
}
