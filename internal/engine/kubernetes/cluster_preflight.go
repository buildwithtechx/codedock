package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type NetworkPreflight struct {
	NodesReady   int      `json:"nodesReady"`
	CoreDNSReady bool     `json:"coreDNSReady"`
	Checks       []string `json:"checks"`
}

func (r *WorkloadRuntime) NetworkPreflight(ctx context.Context, node models.ClusterNode) (*NetworkPreflight, error) {
	preflight := &NetworkPreflight{}
	raw, err := r.commands.Kubectl(ctx, node, []string{"get", "nodes", "-o", "json"}, "")
	if err != nil {
		return nil, err
	}
	var nodes map[string]any
	if err := json.Unmarshal([]byte(raw), &nodes); err != nil {
		return nil, err
	}
	items, _ := nodes["items"].([]any)
	for _, item := range items {
		object, _ := item.(map[string]any)
		status, _ := object["status"].(map[string]any)
		conditions, _ := status["conditions"].([]any)
		for _, condition := range conditions {
			entry, _ := condition.(map[string]any)
			if entry["type"] == "Ready" && entry["status"] == "True" {
				preflight.NodesReady++
			}
		}
	}
	preflight.Checks = append(preflight.Checks, fmt.Sprintf("%d nodes report Ready", preflight.NodesReady))
	dns, err := r.commands.Kubectl(ctx, node, []string{"-n", "kube-system", "get", "deployment", "coredns", "-o", "json"}, "")
	if err != nil {
		return nil, fmt.Errorf("CoreDNS deployment unavailable: %w", err)
	}
	var deployment map[string]any
	if err := json.Unmarshal([]byte(dns), &deployment); err != nil {
		return nil, err
	}
	status, _ := deployment["status"].(map[string]any)
	if available, _ := status["availableReplicas"].(float64); available >= 1 {
		preflight.CoreDNSReady = true
		preflight.Checks = append(preflight.Checks, "CoreDNS has available replicas")
	} else {
		return nil, fmt.Errorf("CoreDNS has no available replicas")
	}
	services, err := r.commands.Kubectl(ctx, node, []string{"get", "endpoints", "kubernetes", "-o", "json"}, "")
	if err != nil || !strings.Contains(services, "443") {
		return nil, fmt.Errorf("Kubernetes API service endpoints unavailable")
	}
	preflight.Checks = append(preflight.Checks, "Kubernetes API endpoints resolve")
	return preflight, nil
}
