package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func validateExistingClaim(existing, desired map[string]any) error {
	current, _ := existing["spec"].(map[string]any)
	requested, _ := desired["spec"].(map[string]any)
	for _, key := range []string{"storageClassName", "accessModes", "resources"} {
		if key == "storageClassName" && requested[key] == nil {
			continue
		}
		a, err := json.Marshal(current[key])
		if err != nil {
			return err
		}
		b, err := json.Marshal(requested[key])
		if err != nil {
			return err
		}
		if string(a) != string(b) {
			return fmt.Errorf("persistent claim %s cannot change during an application rollout; migrate storage separately", key)
		}
	}
	return nil
}
func (r *WorkloadRuntime) Namespace(ctx context.Context, node models.ClusterNode, app *models.AppService) error {
	namespace, _, err := WorkloadIdentity(app)
	if err != nil {
		return err
	}
	response, err := r.commands.Kubectl(ctx, node, []string{"get", "namespace", namespace, "--ignore-not-found", "-o", "json"}, "")
	if err != nil {
		return err
	}
	if response != "" {
		var resource map[string]any
		if err := json.Unmarshal([]byte(response), &resource); err != nil {
			return err
		}
		metadata, _ := resource["metadata"].(map[string]any)
		labels, _ := metadata["labels"].(map[string]any)
		if labels["codedock.run/project"] != app.ProjectID || labels["codedock.run/environment"] != app.EnvironmentID {
			return fmt.Errorf("namespace belongs to another project or environment")
		}
		return nil
	}
	data, err := json.Marshal(map[string]any{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]any{"name": namespace, "labels": map[string]string{"codedock.run/project": app.ProjectID, "codedock.run/environment": app.EnvironmentID}}})
	if err != nil {
		return err
	}
	_, err = r.commands.Kubectl(ctx, node, []string{"create", "-f", "-"}, string(data))
	return err
}
func (r *WorkloadRuntime) deleteOwned(ctx context.Context, node models.ClusterNode, app *models.AppService, kind, name string) error {
	namespace, _, err := WorkloadIdentity(app)
	if err != nil {
		return err
	}
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
	if err := OwnedResource(object, app.ProjectID, app.ID); err != nil {
		return err
	}
	metadata, _ := object["metadata"].(map[string]any)
	uid, _ := metadata["uid"].(string)
	if uid == "" {
		return fmt.Errorf("resource UID missing")
	}
	routes := map[string]string{"deployment": "/apis/apps/v1/namespaces/", "service": "/api/v1/namespaces/", "ingress": "/apis/networking.k8s.io/v1/namespaces/", "secret": "/api/v1/namespaces/"}
	prefix, ok := routes[kind]
	if !ok {
		return fmt.Errorf("unsupported resource cleanup")
	}
	data, err := json.Marshal(map[string]any{"apiVersion": "v1", "kind": "DeleteOptions", "preconditions": map[string]string{"uid": uid}, "propagationPolicy": "Foreground"})
	if err != nil {
		return err
	}
	plurals := map[string]string{"deployment": "deployments", "service": "services", "ingress": "ingresses", "secret": "secrets"}
	_, err = r.commands.Kubectl(ctx, node, []string{"delete", "--raw", prefix + namespace + "/" + plurals[kind] + "/" + name, "-f", "-"}, string(data))
	return err
}

func (r *WorkloadRuntime) Remove(ctx context.Context, node models.ClusterNode, app *models.AppService) error {
	namespace, name, err := WorkloadIdentity(app)
	if err != nil {
		return err
	}
	if err := r.deleteOwned(ctx, node, app, "deployment", name); err != nil {
		return err
	}
	if _, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "wait", "--for=delete", "deployment/" + name, "--timeout=60s"}, ""); err != nil {
		return err
	}
	for _, kind := range []string{"service", "ingress"} {
		if err := r.deleteOwned(ctx, node, app, kind, name); err != nil {
			return err
		}
	}
	for _, suffix := range []string{"-env", "-registry"} {
		if err := r.deleteOwned(ctx, node, app, "secret", name+suffix); err != nil {
			return err
		}
	}
	return nil
}
