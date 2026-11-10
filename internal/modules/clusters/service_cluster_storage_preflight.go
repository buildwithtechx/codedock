package clusters

import (
	"codedock/internal/modules/clusters/kubernetes"
	"context"
	"fmt"
)

type kubeRunnerAdapter struct {
	runner Runner
}

func (a kubeRunnerAdapter) Kubectl(ctx context.Context, node kubernetes.ClusterNode, args []string, input string) (string, error) {
	return a.runner.Kubectl(ctx, ClusterNode(node), args, input)
}

func toKubeNode(n ClusterNode) kubernetes.ClusterNode {
	return kubernetes.ClusterNode{
		ServerID:    n.ServerID,
		PrivateIP:   n.PrivateIP,
		Interface:   n.Interface,
		Fingerprint: n.Fingerprint,
	}
}

func (s *ClusterService) StorageInventory(ctx context.Context, project, id string) (*kubernetes.StorageInventory, error) {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if cluster.ProjectID != project || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return nil, fmt.Errorf("select a ready cluster in this project")
	}
	runtime := kubernetes.NewWorkloadRuntime(kubeRunnerAdapter{runner: s.runner})
	return runtime.StorageInventory(ctx, toKubeNode(cluster.Nodes[0]))
}

func (s *ClusterService) StorageHealth(ctx context.Context, project, id string) error {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if cluster.ProjectID != project || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return fmt.Errorf("select a ready cluster in this project")
	}
	runtime := kubernetes.NewWorkloadRuntime(kubeRunnerAdapter{runner: s.runner})
	return runtime.StorageHealth(ctx, toKubeNode(cluster.Nodes[0]))
}

func (s *ClusterService) NetworkPreflight(ctx context.Context, project, id string) (*kubernetes.NetworkPreflight, error) {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if cluster.ProjectID != project || len(cluster.Nodes) == 0 {
		return nil, fmt.Errorf("select a cluster in this project")
	}
	runtime := kubernetes.NewWorkloadRuntime(kubeRunnerAdapter{runner: s.runner})
	preflight, err := runtime.NetworkPreflight(ctx, toKubeNode(cluster.Nodes[0]))
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
