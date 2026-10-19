package deployments

import (
	"context"
)

func (s *DeploymentService) TogglePin(ctx context.Context, id string) (bool, error) {
	return s.repo.TogglePinned(ctx, id)
}
