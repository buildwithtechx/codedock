package deploy

import (
	"context"
	"fmt"
	"time"

	"codedock/internal/utils"
)

type AttachedInfrastructureGate struct {
	deployer *Deployer
}

func NewAttachedInfrastructureGate(d *Deployer) *AttachedInfrastructureGate {
	return &AttachedInfrastructureGate{deployer: d}
}

func (g *AttachedInfrastructureGate) WaitForAttachedInfrastructure(ctx context.Context, sources []string, timeout time.Duration) error {
	if utils.IsDryRun() || len(sources) == 0 {
		return nil
	}

	if timeout <= 0 {
		timeout = 60 * time.Second
	}

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for _, source := range sources {
		ready := false
		for time.Now().Before(deadline) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
			}

			err := g.deployer.DependencyReady(ctx, source)
			if err == nil {
				ready = true
				break
			}
		}

		if !ready {
			return fmt.Errorf("attached infrastructure prerequisite %s failed to become healthy within %v", source, timeout)
		}
	}

	return nil
}

func FilterAttachedInfrastructureFromTeardown(candidates []string, attachedIDs []string) []string {
	attachedSet := make(map[string]struct{}, len(attachedIDs))
	for _, id := range attachedIDs {
		attachedSet[id] = struct{}{}
	}

	filtered := make([]string, 0, len(candidates))
	for _, c := range candidates {
		if _, isAttached := attachedSet[c]; !isAttached {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
