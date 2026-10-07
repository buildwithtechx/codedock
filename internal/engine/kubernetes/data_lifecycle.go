package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func dataOwnedLabels(record *models.ClusterData) map[string]string {
	return map[string]string{
		"codedock.run/project":     record.ProjectID,
		"codedock.run/database":    record.Spec.ID,
		"codedock.run/environment": record.Spec.EnvironmentID,
	}
}

func (r *WorkloadRuntime) verifyDataOwned(object map[string]any, record *models.ClusterData) error {
	metadata, _ := object["metadata"].(map[string]any)
	labels, _ := metadata["labels"].(map[string]any)
	for key, value := range dataOwnedLabels(record) {
		if labels[key] != value {
			return fmt.Errorf("database resource belongs to another project or database")
		}
	}
	return nil
}

func (r *WorkloadRuntime) deleteDataKind(ctx context.Context, node models.ClusterNode, namespace, kind, name string, record *models.ClusterData) error {
	response, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", kind, name, "--ignore-not-found", "-o", "json"}, "")
	if err != nil {
		return err
	}
	if response == "" {
		return nil
	}
	var object map[string]any
	if err := json.Unmarshal([]byte(response), &object); err != nil {
		return err
	}
	if err := r.verifyDataOwned(object, record); err != nil {
		return err
	}
	_, err = r.commands.Kubectl(ctx, node, []string{"-n", namespace, "delete", kind, name, "--wait=true", "--timeout=120s"}, "")
	return err
}

func (r *WorkloadRuntime) deleteDataPVCs(ctx context.Context, node models.ClusterNode, namespace string, record *models.ClusterData) error {
	response, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "pvc", "-l", "codedock.run/database=" + record.Spec.ID, "-o", "json"}, "")
	if err != nil {
		return err
	}
	var list map[string]any
	if err := json.Unmarshal([]byte(response), &list); err != nil {
		return err
	}
	items, _ := list["items"].([]any)
	for _, item := range items {
		object, _ := item.(map[string]any)
		if err := r.verifyDataOwned(object, record); err != nil {
			return err
		}
		metadata, _ := object["metadata"].(map[string]any)
		name, _ := metadata["name"].(string)
		if name == "" {
			continue
		}
		if _, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "delete", "pvc", name, "--wait=true", "--timeout=120s"}, ""); err != nil {
			return err
		}
	}
	return nil
}

func (r *WorkloadRuntime) StopOwnedData(ctx context.Context, node models.ClusterNode, record *models.ClusterData) error {
	namespace, name := DataIdentity(record.Spec)
	if record.Spec.Engine == "postgres" {
		if err := r.deleteDataKind(ctx, node, namespace, "cluster.postgresql.cnpg.io", name, record); err != nil {
			return err
		}
		return nil
	}
	if err := r.deleteDataKind(ctx, node, namespace, "statefulset", name, record); err != nil {
		return err
	}
	if err := r.deleteDataKind(ctx, node, namespace, "service", name, record); err != nil {
		return err
	}
	return nil
}

func (r *WorkloadRuntime) RemoveOwnedData(ctx context.Context, node models.ClusterNode, record *models.ClusterData) error {
	namespace, name := DataIdentity(record.Spec)
	if record.Spec.Engine == "postgres" {
		if err := r.deleteDataKind(ctx, node, namespace, "cluster.postgresql.cnpg.io", name, record); err != nil {
			return err
		}
		for _, suffix := range []string{"-rw", "-ro", "-r"} {
			_ = r.deleteDataKind(ctx, node, namespace, "service", name+suffix, record)
		}
	} else {
		if err := r.deleteDataKind(ctx, node, namespace, "statefulset", name, record); err != nil {
			return err
		}
		_ = r.deleteDataKind(ctx, node, namespace, "service", name, record)
	}
	if err := r.deleteDataKind(ctx, node, namespace, "secret", name+"-auth", record); err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
	}
	return r.deleteDataPVCs(ctx, node, namespace, record)
}
