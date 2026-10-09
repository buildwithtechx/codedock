package clusters

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (s *Service) waitRelease(ctx context.Context, cluster *models.Cluster, node models.ClusterNode) error {
	timeout, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		raw, err := s.runner.Kubectl(timeout, cluster.Nodes[0], []string{"get", "node", "codedock-" + node.ServerID, "-o", "json"}, "")
		if err == nil {
			var result struct {
				Metadata struct {
					Labels map[string]string `json:"labels"`
				} `json:"metadata"`
				Status struct {
					NodeInfo struct {
						Version string `json:"kubeletVersion"`
					} `json:"nodeInfo"`
					Conditions []struct {
						Type   string `json:"type"`
						Status string `json:"status"`
					} `json:"conditions"`
				} `json:"status"`
			}
			if err := json.Unmarshal([]byte(raw), &result); err != nil {
				return err
			}
			if result.Metadata.Labels["codedock.run/cluster"] != cluster.ID {
				return fmt.Errorf("node cluster ownership changed")
			}
			for _, condition := range result.Status.Conditions {
				if condition.Type == "Ready" && condition.Status == "True" && result.Status.NodeInfo.Version == cluster.Version {
					return nil
				}
			}
		}
		select {
		case <-timeout.Done():
			return fmt.Errorf("node %s did not reach reviewed release %s: %w", node.ServerID, cluster.Version, timeout.Err())
		case <-ticker.C:
		}
	}
}
