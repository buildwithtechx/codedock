package dockerprobe

import (
	"context"
	"io"
	"testing"

	"codedock/internal/models"
)

type mockSwarmRunner struct {
	output string
	err    error
}

func (m *mockSwarmRunner) Run(_ context.Context, _ string) (string, error) {
	return m.output, m.err
}

func (m *mockSwarmRunner) RunPipe(_ context.Context, _ string, _ io.Reader, _ io.Writer) error {
	return nil
}

func TestInspectSwarmActiveManager(t *testing.T) {
	runner := &mockSwarmRunner{
		output: `{"NodeID":"node-123","LocalNodeState":"active","ControlAvailable":true,"Nodes":3,"Managers":1,"Cluster":{"ID":"cluster-abc"}}`,
	}

	inspection, err := InspectSwarm(context.Background(), runner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inspection.State != models.SwarmStateActive {
		t.Fatalf("expected state active, got %s", inspection.State)
	}
	if inspection.Role == nil || *inspection.Role != "manager" {
		t.Fatalf("expected role manager, got %v", inspection.Role)
	}
	if inspection.NodeCount == nil || *inspection.NodeCount != 3 {
		t.Fatalf("expected 3 nodes, got %v", inspection.NodeCount)
	}
}

func TestInspectSwarmInactive(t *testing.T) {
	runner := &mockSwarmRunner{
		output: `{"NodeID":"","LocalNodeState":"inactive","ControlAvailable":false,"Nodes":0,"Managers":0}`,
	}

	inspection, err := InspectSwarm(context.Background(), runner)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if inspection.State != models.SwarmStateInactive {
		t.Fatalf("expected state inactive, got %s", inspection.State)
	}
	if inspection.Role != nil {
		t.Fatalf("expected nil role for inactive swarm, got %v", inspection.Role)
	}
}
