package clusters

import (
	"codedock.run/codedock/internal/engine/kubernetes"
	"context"
	"fmt"
)

func (s *Service) StorageInventory(ctx context.Context, project, id string) (*kubernetes.StorageInventory, error) {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if cluster.ProjectID != project || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return nil, fmt.Errorf("select a ready cluster in this project")
	}
	runtime := kubernetes.NewWorkloadRuntime(s.runner)
	return runtime.StorageInventory(ctx, cluster.Nodes[0])
}

func (s *Service) StorageHealth(ctx context.Context, project, id string) error {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if cluster.ProjectID != project || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return fmt.Errorf("select a ready cluster in this project")
	}
	runtime := kubernetes.NewWorkloadRuntime(s.runner)
	return runtime.StorageHealth(ctx, cluster.Nodes[0])
}

func (s *Service) NetworkPreflight(ctx context.Context, project, id string) (*kubernetes.NetworkPreflight, error) {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if cluster.ProjectID != project || len(cluster.Nodes) == 0 {
		return nil, fmt.Errorf("select a cluster in this project")
	}
	runtime := kubernetes.NewWorkloadRuntime(s.runner)
	preflight, err := runtime.NetworkPreflight(ctx, cluster.Nodes[0])
	if err != nil {
		return nil, err
	}
	peers := make([]string, 0, len(cluster.Nodes))
	for _, node := range cluster.Nodes {
		peers = append(peers, node.PrivateIP)
	}
	reached, err := runtime.CrossNodeProbe(ctx, cluster.Nodes[0], peers)
	if err != nil {
		return nil, err
	}
	preflight.CrossNode = reached
	preflight.Checks = append(preflight.Checks, fmt.Sprintf("disposable probe reached %d peers and cluster DNS", len(reached)))
	return preflight, nil
}
