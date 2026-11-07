package kubernetes

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type Streamer interface {
	StreamKubectl(context.Context, models.ClusterNode, []string, io.Writer) error
	StreamKubectlPTY(context.Context, models.ClusterNode, []string, io.Reader, io.Writer) error
}

func (r *WorkloadRuntime) resolvePod(ctx context.Context, node models.ClusterNode, app *models.AppService, pod string) (string, string, error) {
	observation, err := r.Observe(ctx, node, app)
	if err != nil {
		return "", "", err
	}
	for _, candidate := range observation.Pods {
		if candidate.Ready && (pod == "" || pod == candidate.Name) {
			namespace, _, err := WorkloadIdentity(app)
			if err != nil {
				return "", "", err
			}
			return namespace, candidate.Name, nil
		}
	}
	return "", "", fmt.Errorf("select a ready pod owned by this service")
}

func (r *WorkloadRuntime) StreamLogs(ctx context.Context, streamer Streamer, node models.ClusterNode, app *models.AppService, pod string, output io.Writer) error {
	namespace, name, err := r.resolvePod(ctx, node, app, pod)
	if err != nil {
		return err
	}
	return streamer.StreamKubectl(ctx, node, []string{"-n", namespace, "logs", name, "-c", "app", "--follow", "--tail=200", "--timestamps"}, output)
}

func (r *WorkloadRuntime) StreamExec(ctx context.Context, streamer Streamer, node models.ClusterNode, app *models.AppService, pod string, command []string, input io.Reader, output io.Writer) error {
	if len(command) == 0 || len(command) > 64 {
		return fmt.Errorf("provide one to sixty-four command arguments")
	}
	for _, arg := range command {
		if len(arg) > 8192 || strings.ContainsRune(arg, 0) {
			return fmt.Errorf("invalid command argument")
		}
	}
	namespace, name, err := r.resolvePod(ctx, node, app, pod)
	if err != nil {
		return err
	}
	args := append([]string{"-n", namespace, "exec", "-i", "-t", name, "-c", "app", "--"}, command...)
	return streamer.StreamKubectlPTY(ctx, node, args, input, output)
}

type InstanceNode struct {
	ID     string `json:"id"`
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Node   string `json:"node,omitempty"`
}

type InstanceEdge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
}

type InstanceGraph struct {
	Nodes []InstanceNode `json:"nodes"`
	Edges []InstanceEdge `json:"edges"`
}

func (r *WorkloadRuntime) InstanceGraph(ctx context.Context, node models.ClusterNode, app *models.AppService) (*InstanceGraph, error) {
	observation, err := r.Observe(ctx, node, app)
	if err != nil {
		return nil, err
	}
	graph := &InstanceGraph{}
	namespace, name, err := WorkloadIdentity(app)
	if err != nil {
		return nil, err
	}
	graph.Nodes = append(graph.Nodes, InstanceNode{ID: "deployment/" + name, Kind: "deployment", Name: name, Status: observation.Status})
	for _, pod := range observation.Pods {
		status := "running"
		if !pod.Ready {
			status = pod.Phase
		}
		graph.Nodes = append(graph.Nodes, InstanceNode{ID: "pod/" + pod.Name, Kind: "pod", Name: pod.Name, Status: status, Node: pod.Node})
		graph.Edges = append(graph.Edges, InstanceEdge{Source: "deployment/" + name, Target: "pod/" + pod.Name, Label: "schedules"})
		if pod.Node != "" {
			graph.Nodes = append(graph.Nodes, InstanceNode{ID: "node/" + pod.Node, Kind: "node", Name: pod.Node, Status: "ready"})
			graph.Edges = append(graph.Edges, InstanceEdge{Source: "pod/" + pod.Name, Target: "node/" + pod.Node, Label: "runs on"})
		}
	}
	raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "pvc", "-l", "codedock.run/service=" + app.ID, "-o", "json"}, "")
	if err == nil {
		var list struct {
			Items []map[string]any `json:"items"`
		}
		if err := json.Unmarshal([]byte(raw), &list); err == nil {
			for _, item := range list.Items {
				metadata, _ := item["metadata"].(map[string]any)
				claim, _ := metadata["name"].(string)
				status, _ := item["status"].(map[string]any)
				graph.Nodes = append(graph.Nodes, InstanceNode{ID: "pvc/" + claim, Kind: "volume", Name: claim, Status: value(status["phase"])})
				graph.Edges = append(graph.Edges, InstanceEdge{Source: "deployment/" + name, Target: "pvc/" + claim, Label: "mounts"})
			}
		}
	}
	return graph, nil
}

type ResourceControl struct {
	Kind    string `json:"kind"`
	Name    string `json:"name"`
	Ready   string `json:"ready"`
	Message string `json:"message"`
}

func (r *WorkloadRuntime) ResourceSummary(ctx context.Context, node models.ClusterNode, app *models.AppService) ([]ResourceControl, error) {
	namespace, name, err := WorkloadIdentity(app)
	if err != nil {
		return nil, err
	}
	summary := []ResourceControl{}
	for _, kind := range []string{"deployment", "service", "ingress", "secret", "persistentvolumeclaim"} {
		raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", kind, "-l", "codedock.run/service=" + app.ID, "-o", "json"}, "")
		if err != nil {
			continue
		}
		var list struct {
			Items []map[string]any `json:"items"`
		}
		if err := json.Unmarshal([]byte(raw), &list); err != nil {
			continue
		}
		for _, item := range list.Items {
			if err := OwnedResource(item, app.ProjectID, app.ID); err != nil {
				continue
			}
			metadata, _ := item["metadata"].(map[string]any)
			status, _ := item["status"].(map[string]any)
			control := ResourceControl{Kind: kind, Name: value(metadata["name"]), Ready: "unknown"}
			switch kind {
			case "deployment":
				control.Ready = fmt.Sprintf("%v/%v available", status["availableReplicas"], status["replicas"])
			case "service":
				spec, _ := item["spec"].(map[string]any)
				control.Ready = value(spec["clusterIP"])
			case "ingress":
				control.Ready = "routed"
			case "secret":
				control.Ready = "stored"
			case "persistentvolumeclaim":
				control.Ready = value(status["phase"])
			}
			_ = name
			summary = append(summary, control)
		}
	}
	return summary, nil
}
