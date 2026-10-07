package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func (r *WorkloadRuntime) ApplyOwnedData(ctx context.Context, node models.ClusterNode, manifest string, progress func(string, string) error) error {
	var document struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(manifest), &document); err != nil {
		return err
	}
	if len(document.Items) == 0 {
		return fmt.Errorf("empty database resource manifest")
	}
	for _, object := range document.Items {
		metadata, ok := object["metadata"].(map[string]any)
		if !ok {
			return fmt.Errorf("resource metadata missing")
		}
		name, _ := metadata["name"].(string)
		namespace, _ := metadata["namespace"].(string)
		kind, _ := object["kind"].(string)
		api, _ := object["apiVersion"].(string)
		if name == "" || kind == "" {
			return fmt.Errorf("resource identity missing")
		}
		resource := kind
		if group, _, ok := strings.Cut(api, "/"); ok {
			resource += "." + group
		}
		args := []string{}
		if namespace != "" {
			args = append(args, "-n", namespace)
		}
		raw, err := r.commands.Kubectl(ctx, node, append(args, "get", resource, name, "--ignore-not-found", "-o", "json"), "")
		if err != nil {
			return err
		}
		operation := "create"
		if raw != "" {
			var existing map[string]any
			if err := json.Unmarshal([]byte(raw), &existing); err != nil {
				return err
			}
			if err := matchingDataOwner(existing, object); err != nil {
				return err
			}
			previous := existing["metadata"].(map[string]any)
			metadata["uid"] = previous["uid"]
			metadata["resourceVersion"] = previous["resourceVersion"]
			operation = "apply"
		}
		data, err := json.Marshal(object)
		if err != nil {
			return err
		}
		call := []string{operation}
		if operation == "apply" {
			call = append(call, "--server-side", "--field-manager=codedock-databases")
		}
		call = append(call, "-f", "-")
		if _, err := r.commands.Kubectl(ctx, node, call, string(data)); err != nil {
			return fmt.Errorf("%s %s: %w", operation, name, err)
		}
		if err := progress("APPLYING", kind+" "+name); err != nil {
			return err
		}
	}
	return nil
}
func matchingDataOwner(existing, desired map[string]any) error {
	wanted, _ := desired["metadata"].(map[string]any)
	actual, _ := existing["metadata"].(map[string]any)
	labels, _ := actual["labels"].(map[string]any)
	expected, _ := wanted["labels"].(map[string]any)
	checked := false
	for _, key := range []string{"codedock.run/project", "codedock.run/database", "codedock.run/environment", "codedock.run/cluster"} {
		if value, ok := expected[key]; ok {
			checked = true
			if labels[key] != value {
				return fmt.Errorf("resource %v belongs to another owner", wanted["name"])
			}
		}
	}
	if !checked {
		return fmt.Errorf("resource ownership labels missing")
	}
	return nil
}
func (r *WorkloadRuntime) DataReady(ctx context.Context, node models.ClusterNode, record *models.ClusterData) error {
	namespace, name := DataIdentity(record.Spec)
	if record.Spec.Engine == "postgres" {
		resource := "cluster.postgresql.cnpg.io"
		args := []string{"-n", namespace, "wait", "--for=condition=Ready", resource + "/" + name, "--timeout=900s"}
		if _, err := r.commands.Kubectl(ctx, node, args, ""); err != nil {
			return err
		}
		return r.DataObserved(ctx, node, record)
	}
	shards := record.Spec.Shards
	if shards < 2 {
		shards = 1
	}
	for shard := 0; shard < shards; shard++ {
		target := name
		if record.Spec.Shards > 1 {
			target = fmt.Sprintf("%s-s%d", name, shard)
		}
		if _, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "rollout", "status", "statefulset/" + target, "--timeout=600s"}, ""); err != nil {
			return err
		}
	}
	return r.DataObserved(ctx, node, record)
}
func (r *WorkloadRuntime) DataObserved(ctx context.Context, node models.ClusterNode, record *models.ClusterData) error {
	namespace, name := DataIdentity(record.Spec)
	if record.Spec.Engine == "postgres" {
		resource := "cluster.postgresql.cnpg.io"
		raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", resource, name, "-o", "json"}, "")
		if err != nil {
			return err
		}
		var object map[string]any
		if err := json.Unmarshal([]byte(raw), &object); err != nil {
			return err
		}
		desired := map[string]any{"metadata": map[string]any{"name": name, "labels": map[string]any{"codedock.run/project": record.ProjectID, "codedock.run/database": record.ID}}}
		if err := matchingDataOwner(object, desired); err != nil {
			return err
		}
		status, _ := object["status"].(map[string]any)
		if ready := int(number(status["readyInstances"])); ready != record.Spec.Instances {
			return fmt.Errorf("database has %d/%d ready instances", ready, record.Spec.Instances)
		}
		return nil
	}
	shards := record.Spec.Shards
	if shards < 2 {
		shards = 1
	}
	for shard := 0; shard < shards; shard++ {
		target := name
		if record.Spec.Shards > 1 {
			target = fmt.Sprintf("%s-s%d", name, shard)
		}
		raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "statefulset.apps", target, "-o", "json"}, "")
		if err != nil {
			return err
		}
		var object map[string]any
		if err := json.Unmarshal([]byte(raw), &object); err != nil {
			return err
		}
		desired := map[string]any{"metadata": map[string]any{"name": target, "labels": map[string]any{"codedock.run/project": record.ProjectID, "codedock.run/database": record.ID}}}
		if err := matchingDataOwner(object, desired); err != nil {
			return err
		}
		status, _ := object["status"].(map[string]any)
		if ready := int(number(status["readyReplicas"])); ready != record.Spec.Instances {
			return fmt.Errorf("shard %d has %d/%d ready instances", shard, ready, record.Spec.Instances)
		}
	}
	return nil
}
