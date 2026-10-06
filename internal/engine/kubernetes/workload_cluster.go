package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (r *WorkloadRuntime) VerifyCluster(ctx context.Context, cluster *models.Cluster, requireReady bool) error {
	if len(cluster.Nodes) == 0 {
		return fmt.Errorf("cluster has no control plane")
	}
	raw, err := r.commands.Kubectl(ctx, cluster.Nodes[0], []string{"get", "nodes", "-o", "json"}, "")
	if err != nil {
		return err
	}
	var response struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(raw), &response); err != nil {
		return err
	}
	observed := map[string]bool{}
	for _, object := range response.Items {
		metadata, _ := object["metadata"].(map[string]any)
		labels, _ := metadata["labels"].(map[string]any)
		if labels["codedock.run/cluster"] != cluster.ID {
			return fmt.Errorf("control plane contains an unregistered cluster node")
		}
		name, _ := metadata["name"].(string)
		status, _ := object["status"].(map[string]any)
		ready := false
		conditions, _ := status["conditions"].([]any)
		for _, raw := range conditions {
			condition, _ := raw.(map[string]any)
			if condition["type"] == "Ready" && condition["status"] == "True" {
				ready = true
			}
		}
		observed[name] = ready
	}
	if len(observed) != len(cluster.Nodes) {
		return fmt.Errorf("observed node count differs from registered cluster")
	}
	for _, node := range cluster.Nodes {
		ready, exists := observed["codedock-"+node.ServerID]
		if !exists || (requireReady && !ready) {
			return fmt.Errorf("registered cluster node is missing or not ready")
		}
	}
	return nil
}
