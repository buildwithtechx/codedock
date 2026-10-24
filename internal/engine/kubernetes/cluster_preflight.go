package kubernetes

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type NetworkPreflight struct {
	NodesReady   int      `json:"nodesReady"`
	CoreDNSReady bool     `json:"coreDNSReady"`
	CrossNode    []string `json:"crossNode"`
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

func (r *WorkloadRuntime) CrossNodeProbe(ctx context.Context, node models.ClusterNode, peers []string) ([]string, error) {
	manifest := map[string]any{"apiVersion": "v1", "kind": "Pod", "metadata": map[string]any{"name": "codedock-net-probe", "namespace": "default", "labels": map[string]string{"codedock.run/probe": "net"}}, "spec": map[string]any{"restartPolicy": "Never", "containers": []map[string]any{{"name": "probe", "image": "busybox:1.36", "command": []string{"sh", "-ec", "trap 'exit 0' TERM; sleep 300"}}}}}
	data, err := json.Marshal(manifest)
	if err != nil {
		return nil, err
	}
	cleanup := func() {
		_, _ = r.commands.Kubectl(ctx, node, []string{"delete", "pod", "codedock-net-probe", "--ignore-not-found", "--wait=true", "--timeout=60s"}, "")
	}
	cleanup()
	if _, err := r.commands.Kubectl(ctx, node, []string{"create", "-f", "-"}, string(data)); err != nil {
		return nil, err
	}
	defer cleanup()
	if _, err := r.commands.Kubectl(ctx, node, []string{"wait", "--for=condition=Ready", "pod/codedock-net-probe", "--timeout=120s"}, ""); err != nil {
		return nil, fmt.Errorf("disposable probe did not become ready: %w", err)
	}
	results := []string{}
	for _, peer := range peers {
		if strings.TrimSpace(peer) == "" {
			continue
		}
		if _, err := r.commands.Kubectl(ctx, node, []string{"exec", "codedock-net-probe", "--", "ping", "-c", "2", "-W", "2", peer}, ""); err != nil {
			return nil, fmt.Errorf("probe cannot reach peer %s", peer)
		}
		results = append(results, peer)
	}
	if _, err := r.commands.Kubectl(ctx, node, []string{"exec", "codedock-net-probe", "--", "nslookup", "kubernetes.default"}, ""); err != nil {
		return nil, fmt.Errorf("probe cannot resolve cluster DNS")
	}
	results = append(results, "kubernetes.default")
	return results, nil
}
