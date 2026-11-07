package kubernetes

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (r *WorkloadRuntime) Lifecycle(ctx context.Context, node models.ClusterNode, app *models.AppService, action string, replicas int) error {
	namespace, name, err := WorkloadIdentity(app)
	if err != nil {
		return err
	}
	raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "deployment", name, "-o", "json"}, "")
	if err != nil {
		return err
	}
	var object map[string]any
	if err := json.Unmarshal([]byte(raw), &object); err != nil {
		return err
	}
	if err := OwnedResource(object, app.ProjectID, app.ID); err != nil {
		return err
	}
	args := []string{"-n", namespace}
	switch action {
	case "stop":
		args = append(args, "scale", "deployment/"+name, "--replicas=0")
	case "scale":
		if replicas < 1 || replicas > 64 {
			return fmt.Errorf("replicas must be between one and sixty-four")
		}
		args = append(args, "scale", "deployment/"+name, fmt.Sprintf("--replicas=%d", replicas))
	case "restart":
		args = append(args, "rollout", "restart", "deployment/"+name)
	default:
		return fmt.Errorf("unsupported workload action")
	}
	if _, err := r.commands.Kubectl(ctx, node, args, ""); err != nil {
		return err
	}
	if action != "stop" {
		return r.Ready(ctx, node, app)
	}
	return nil
}
func (r *WorkloadRuntime) Observe(ctx context.Context, node models.ClusterNode, app *models.AppService) (*models.WorkloadObservation, error) {
	namespace, name, err := WorkloadIdentity(app)
	if err != nil {
		return nil, err
	}
	result := &models.WorkloadObservation{Kind: "kubernetes", Status: "NOT_DEPLOYED", Pods: []models.RuntimePod{}}
	raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "deployment", name, "--ignore-not-found", "-o", "json"}, "")
	if err != nil {
		return nil, err
	}
	if raw == "" {
		return result, nil
	}
	var object map[string]any
	if err := json.Unmarshal([]byte(raw), &object); err != nil {
		return nil, err
	}
	if err := OwnedResource(object, app.ProjectID, app.ID); err != nil {
		return nil, err
	}
	spec, _ := object["spec"].(map[string]any)
	status, _ := object["status"].(map[string]any)
	result.Desired = int(number(spec["replicas"]))
	result.Available = int(number(status["availableReplicas"]))
	result.Status = "PENDING"
	if result.Desired == 0 {
		result.Status = "STOPPED"
	} else if result.Available == result.Desired {
		result.Status = "READY"
	}
	metadata, _ := object["metadata"].(map[string]any)
	if number(status["observedGeneration"]) < number(metadata["generation"]) {
		result.Status = "PENDING"
	}
	conditions, _ := status["conditions"].([]any)
	for _, raw := range conditions {
		condition, _ := raw.(map[string]any)
		if condition["type"] == "Progressing" && condition["status"] == "False" {
			result.Status = "FAILED"
			result.Error = value(condition["message"])
		}
	}
	raw, err = r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "pods", "-l", "codedock.run/service=" + app.ID, "-o", "json"}, "")
	if err != nil {
		return nil, err
	}
	var list struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return nil, err
	}
	for _, pod := range list.Items {
		if err := OwnedResource(pod, app.ProjectID, app.ID); err != nil {
			return nil, err
		}
		metadata, _ := pod["metadata"].(map[string]any)
		spec, _ := pod["spec"].(map[string]any)
		status, _ := pod["status"].(map[string]any)
		item := models.RuntimePod{Name: value(metadata["name"]), Node: value(spec["nodeName"]), Phase: value(status["phase"])}
		containers, _ := status["containerStatuses"].([]any)
		item.Ready = len(containers) > 0
		for _, raw := range containers {
			container, _ := raw.(map[string]any)
			ready, _ := container["ready"].(bool)
			item.Ready = item.Ready && ready
			item.Restarts += int(number(container["restartCount"]))
		}
		result.Pods = append(result.Pods, item)
	}
	raw, err = r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "--raw", "/apis/metrics.k8s.io/v1beta1/namespaces/" + namespace + "/pods"}, "")
	if err != nil {
		result.MetricsError = "Cluster metrics API unavailable"
		return result, nil
	}
	if err := applyPodMetrics(result, raw); err != nil {
		result.MetricsError = err.Error()
		return result, nil
	}
	result.MetricsAvailable = true
	return result, nil
}
func number(value any) float64 { result, _ := value.(float64); return result }
func value(raw any) string     { result, _ := raw.(string); return result }
func (r *WorkloadRuntime) Logs(ctx context.Context, node models.ClusterNode, app *models.AppService) (string, error) {
	observation, err := r.Observe(ctx, node, app)
	if err != nil {
		return "", err
	}
	namespace, _, err := WorkloadIdentity(app)
	if err != nil {
		return "", err
	}
	output := strings.Builder{}
	for _, pod := range observation.Pods {
		raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "logs", pod.Name, "-c", "app", "--tail=200", "--timestamps"}, "")
		if err != nil {
			return "", err
		}
		fmt.Fprintf(&output, "[%s]\n%s\n", pod.Name, raw)
	}
	return output.String(), nil
}
func (r *WorkloadRuntime) Exec(ctx context.Context, node models.ClusterNode, app *models.AppService, request models.RuntimeExecRequest) (string, error) {
	if len(request.Command) == 0 || len(request.Command) > 64 {
		return "", fmt.Errorf("provide one to sixty-four command arguments")
	}
	for _, arg := range request.Command {
		if len(arg) > 8192 || strings.ContainsRune(arg, 0) {
			return "", fmt.Errorf("invalid command argument")
		}
	}
	observation, err := r.Observe(ctx, node, app)
	if err != nil {
		return "", err
	}
	pod := ""
	for _, candidate := range observation.Pods {
		if candidate.Ready && (request.Pod == "" || request.Pod == candidate.Name) {
			pod = candidate.Name
			break
		}
	}
	if pod == "" {
		return "", fmt.Errorf("select a ready pod owned by this service")
	}
	namespace, _, err := WorkloadIdentity(app)
	if err != nil {
		return "", err
	}
	args := append([]string{"-n", namespace, "exec", pod, "-c", "app", "--"}, request.Command...)
	return r.commands.Kubectl(ctx, node, args, "")
}
