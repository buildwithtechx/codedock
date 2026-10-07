package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (r *WorkloadRuntime) dataPods(ctx context.Context, node models.ClusterNode, namespace, name string) ([]map[string]any, error) {
	response, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "pods", "-l", "codedock.run/database=" + name, "-o", "json"}, "")
	if err != nil {
		return nil, err
	}
	var list map[string]any
	if err := json.Unmarshal([]byte(response), &list); err != nil {
		return nil, err
	}
	items, _ := list["items"].([]any)
	pods := make([]map[string]any, 0, len(items))
	for _, item := range items {
		object, _ := item.(map[string]any)
		pods = append(pods, object)
	}
	return pods, nil
}

func (r *WorkloadRuntime) DataInspect(ctx context.Context, node models.ClusterNode, record *models.ClusterData) ([]string, []string, error) {
	namespace, name := DataIdentity(record.Spec)
	roles := []string{}
	volumes := []string{}
	if record.Spec.Engine == "postgres" {
		raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "cluster.postgresql.cnpg.io", name, "-o", "json"}, "")
		if err != nil {
			return nil, nil, err
		}
		var object map[string]any
		if err := json.Unmarshal([]byte(raw), &object); err != nil {
			return nil, nil, err
		}
		status, _ := object["status"].(map[string]any)
		if primary, _ := status["currentPrimary"].(string); primary != "" {
			roles = append(roles, "primary:"+primary)
		}
		if instances, _ := status["instances"].([]any); len(instances) > 0 {
			for _, instance := range instances {
				if name, _ := instance.(map[string]any)["podName"].(string); name != "" {
					roles = append(roles, "standby:"+name)
				}
			}
		}
	} else {
		pods, err := r.dataPods(ctx, node, namespace, name)
		if err != nil {
			return nil, nil, err
		}
		for _, pod := range pods {
			metadata, _ := pod["metadata"].(map[string]any)
			podName, _ := metadata["name"].(string)
			role := "replica"
			if strings.HasSuffix(podName, "-0") || strings.Contains(podName, name+"-0") {
				role = "primary"
			}
			roles = append(roles, role+":"+podName)
		}
	}
	raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "pvc", "-l", "codedock.run/database=" + record.Spec.ID, "-o", "json"}, "")
	if err != nil {
		return roles, volumes, err
	}
	var list map[string]any
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return roles, volumes, err
	}
	items, _ := list["items"].([]any)
	for _, item := range items {
		object, _ := item.(map[string]any)
		metadata, _ := object["metadata"].(map[string]any)
		if name, _ := metadata["name"].(string); name != "" {
			volumes = append(volumes, name)
		}
	}
	return roles, volumes, nil
}

func (r *WorkloadRuntime) DataConnectionCheck(ctx context.Context, node models.ClusterNode, record *models.ClusterData, password string) error {
	namespace, name := DataIdentity(record.Spec)
	if password == "" {
		return fmt.Errorf("database credentials unavailable")
	}
	if record.Spec.Engine == "postgres" {
		pod := name + "-1"
		if record.Spec.Instances == 1 {
			pods, err := r.dataPods(ctx, node, namespace, name)
			if err != nil || len(pods) == 0 {
				return fmt.Errorf("no database pods available")
			}
			metadata, _ := pods[0]["metadata"].(map[string]any)
			pod, _ = metadata["name"].(string)
		}
		_, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "exec", pod, "-c", "postgres", "--", "env", "PGPASSWORD=" + password, "pg_isready", "-U", "app", "-d", "app", "-h", "127.0.0.1"}, "")
		return err
	}
	pods, err := r.dataPods(ctx, node, namespace, name)
	if err != nil || len(pods) == 0 {
		return fmt.Errorf("no database pods available")
	}
	metadata, _ := pods[0]["metadata"].(map[string]any)
	pod, _ := metadata["name"].(string)
	_, err = r.commands.Kubectl(ctx, node, []string{"-n", namespace, "exec", pod, "-c", "redis", "--", "redis-cli", "-a", password, "ping"}, "")
	return err
}
