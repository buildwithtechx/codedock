package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Commands interface {
	Kubectl(context.Context, models.ClusterNode, []string, string) (string, error)
}
type WorkloadRuntime struct{ commands Commands }

func NewWorkloadRuntime(commands Commands) *WorkloadRuntime { return &WorkloadRuntime{commands} }
func (r *WorkloadRuntime) Apply(ctx context.Context, node models.ClusterNode, manifest string) error {
	_, err := r.commands.Kubectl(ctx, node, []string{"apply", "--server-side", "--field-manager=codedock", "-f", "-"}, manifest)
	if err != nil {
		return fmt.Errorf("apply Kubernetes resources: %w", err)
	}
	return nil
}
func (r *WorkloadRuntime) Validate(ctx context.Context, node models.ClusterNode, manifest string) error {
	_, err := r.commands.Kubectl(ctx, node, []string{"apply", "--server-side", "--dry-run=server", "--field-manager=codedock", "-f", "-"}, manifest)
	if err != nil {
		return fmt.Errorf("validate Kubernetes target capabilities and resource configuration: %w", err)
	}
	return nil
}
func (r *WorkloadRuntime) Ready(ctx context.Context, node models.ClusterNode, app *models.AppService) error {
	namespace, name, err := WorkloadIdentity(app)
	if err != nil {
		return err
	}
	_, err = r.commands.Kubectl(ctx, node, []string{"-n", namespace, "rollout", "status", "deployment/" + name, "--timeout=300s"}, "")
	if err != nil {
		return fmt.Errorf("Kubernetes workload readiness: %w", err)
	}
	return nil
}
func (r *WorkloadRuntime) Snapshot(ctx context.Context, node models.ClusterNode, workload *models.KubernetesWorkload) (string, error) {
	manifest, err := WorkloadManifest(workload)
	if err != nil {
		return "", err
	}
	var desired map[string]any
	if err := json.Unmarshal([]byte(manifest), &desired); err != nil {
		return "", err
	}
	saved := []map[string]any{}
	namespace, _, err := WorkloadIdentity(&workload.App)
	if err != nil {
		return "", err
	}
	desiredItems, ok := desired["items"].([]any)
	if !ok {
		return "", fmt.Errorf("invalid workload manifest")
	}
	for _, raw := range desiredItems {
		object, ok := raw.(map[string]any)
		if !ok {
			return "", fmt.Errorf("invalid workload resource")
		}
		if object["kind"] == "Namespace" {
			continue
		}
		metadata := object["metadata"].(map[string]any)
		name := metadata["name"].(string)
		response, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", object["kind"].(string), name, "--ignore-not-found", "-o", "json"}, "")
		if err != nil {
			return "", fmt.Errorf("snapshot previous owned resource: %w", err)
		}
		if response == "" {
			continue
		}
		var existing map[string]any
		if err := json.Unmarshal([]byte(response), &existing); err != nil {
			return "", err
		}
		if err := OwnedResource(existing, workload.App.ProjectID, workload.App.ID); err != nil {
			return "", err
		}
		if object["kind"] == "PersistentVolumeClaim" {
			if err := validateExistingClaim(existing, object); err != nil {
				return "", err
			}
		}
	}
	response, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "deployments,services,secrets,persistentvolumeclaims,ingresses", "-l", "codedock.run/service=" + workload.App.ID, "-o", "json"}, "")
	if err != nil {
		return "", err
	}
	var resources struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(response), &resources); err != nil {
		return "", err
	}
	if resources.Items == nil {
		return "", fmt.Errorf("invalid resource snapshot")
	}
	for _, existing := range resources.Items {
		if err := OwnedResource(existing, workload.App.ProjectID, workload.App.ID); err != nil {
			return "", err
		}
		sanitizeResource(existing)
		saved = append(saved, existing)
	}
	data, err := json.Marshal(map[string]any{"apiVersion": "v1", "kind": "List", "items": saved})
	if err != nil {
		return "", err
	}
	return string(data), nil
}
func OwnedResource(resource map[string]any, project, service string) error {
	metadata, ok := resource["metadata"].(map[string]any)
	if !ok {
		return fmt.Errorf("resource ownership metadata missing")
	}
	labels, ok := metadata["labels"].(map[string]any)
	if !ok || labels["codedock.run/project"] != project || labels["codedock.run/service"] != service {
		return fmt.Errorf("resource is not owned by this Codedock service")
	}
	return nil
}
func sanitizeResource(resource map[string]any) {
	delete(resource, "status")
	metadata := resource["metadata"].(map[string]any)
	for _, key := range []string{"uid", "resourceVersion", "generation", "creationTimestamp", "managedFields", "selfLink", "deletionTimestamp", "deletionGracePeriodSeconds"} {
		delete(metadata, key)
	}
	if resource["kind"] == "Service" {
		spec := resource["spec"].(map[string]any)
		for _, key := range []string{"clusterIP", "clusterIPs", "ipFamilies", "ipFamilyPolicy", "healthCheckNodePort"} {
			delete(spec, key)
		}
	}
}
func (r *WorkloadRuntime) Rollback(ctx context.Context, node models.ClusterNode, app *models.AppService, journal string) error {
	_, name, err := WorkloadIdentity(app)
	if err != nil {
		return err
	}
	var previous map[string]any
	if err := json.Unmarshal([]byte(journal), &previous); err != nil {
		return err
	}
	items, ok := previous["items"].([]any)
	if !ok {
		return fmt.Errorf("invalid workload recovery journal")
	}
	restore := []map[string]any{}
	namespace, _, err := WorkloadIdentity(app)
	if err != nil {
		return err
	}
	for _, raw := range items {
		object, ok := raw.(map[string]any)
		if !ok {
			return fmt.Errorf("invalid recovery resource")
		}
		if err := OwnedResource(object, app.ProjectID, app.ID); err != nil {
			return err
		}
		metadata, _ := object["metadata"].(map[string]any)
		kind, _ := object["kind"].(string)
		resourceName, _ := metadata["name"].(string)
		response, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", kind, resourceName, "--ignore-not-found", "-o", "json"}, "")
		if err != nil {
			return err
		}
		if response != "" {
			var current map[string]any
			if err := json.Unmarshal([]byte(response), &current); err != nil {
				return err
			}
			if err := OwnedResource(current, app.ProjectID, app.ID); err != nil {
				return err
			}
			currentMetadata, _ := current["metadata"].(map[string]any)
			if currentMetadata["uid"] == nil || currentMetadata["resourceVersion"] == nil {
				return fmt.Errorf("recovery resource identity missing")
			}
			metadata["uid"] = currentMetadata["uid"]
			metadata["resourceVersion"] = currentMetadata["resourceVersion"]
		}
		if kind != "PersistentVolumeClaim" {
			if response == "" {
				data, err := json.Marshal(object)
				if err != nil {
					return err
				}
				if _, err := r.commands.Kubectl(ctx, node, []string{"create", "--field-manager=codedock", "-f", "-"}, string(data)); err != nil {
					return err
				}
			} else {
				restore = append(restore, object)
			}
		}
	}
	if len(restore) > 0 {
		data, err := json.Marshal(map[string]any{"apiVersion": "v1", "kind": "List", "items": restore})
		if err != nil {
			return err
		}
		if err := r.Apply(ctx, node, string(data)); err != nil {
			return err
		}
	}
	hadDeployment := false
	hadIngress := false
	for _, raw := range items {
		object := raw.(map[string]any)
		if object["kind"] == "Deployment" {
			hadDeployment = true
		}
		if object["kind"] == "Ingress" {
			hadIngress = true
		}
	}
	var cleanupErr error
	if !hadDeployment {
		cleanupErr = errors.Join(r.deleteOwned(ctx, node, app, "deployment", name), r.deleteOwned(ctx, node, app, "service", name))
	}
	if !hadIngress {
		cleanupErr = errors.Join(cleanupErr, r.deleteOwned(ctx, node, app, "ingress", name))
	}
	if hadDeployment {
		ready, finish := context.WithTimeout(ctx, 5*time.Minute)
		defer finish()
		cleanupErr = errors.Join(cleanupErr, r.Ready(ready, node, app))
	}
	return cleanupErr
}
