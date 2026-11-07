package deploy

import (
	"context"
	"testing"
	"time"
)

func TestFilterAttachedInfrastructureFromTeardown(t *testing.T) {
	candidates := []string{"container-app-1", "container-db-1", "container-redis-1", "container-app-old"}
	attached := []string{"container-db-1", "container-redis-1"}

	filtered := FilterAttachedInfrastructureFromTeardown(candidates, attached)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 filtered containers, got %d", len(filtered))
	}

	for _, id := range filtered {
		if id == "container-db-1" || id == "container-redis-1" {
			t.Errorf("attached infrastructure %s should have been filtered out of teardown", id)
		}
	}
}

func TestWaitForAttachedInfrastructureEmpty(t *testing.T) {
	gate := NewAttachedInfrastructureGate(&Deployer{})
	err := gate.WaitForAttachedInfrastructure(context.Background(), nil, 5*time.Second)
	if err != nil {
		t.Fatalf("expected no error for empty sources, got: %v", err)
	}
}
