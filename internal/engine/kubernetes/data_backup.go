package kubernetes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

func DataBackupManifest(record *models.ClusterData) (string, string, error) {
	namespace, name := DataIdentity(record.Spec)
	backup := "backup-" + uuid.NewString()
	object := map[string]any{"apiVersion": "postgresql.cnpg.io/v1", "kind": "Backup", "metadata": map[string]any{"name": backup, "namespace": namespace, "labels": map[string]string{"codedock.run/project": record.ProjectID, "codedock.run/database": record.ID}}, "spec": map[string]any{"cluster": map[string]string{"name": name}, "method": "plugin", "pluginConfiguration": map[string]string{"name": "barman-cloud.cloudnative-pg.io"}}}
	data, err := json.Marshal(map[string]any{"apiVersion": "v1", "kind": "List", "items": []map[string]any{object}})
	return string(data), backup, err
}
func (r *WorkloadRuntime) WaitDataBackup(ctx context.Context, node models.ClusterNode, record *models.ClusterData, manifest string) error {
	var list struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(manifest), &list); err != nil {
		return err
	}
	if len(list.Items) != 1 {
		return fmt.Errorf("invalid backup manifest")
	}
	namespace, _ := DataIdentity(record.Spec)
	_, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "wait", "--for=jsonpath={.status.phase}=completed", "backup.postgresql.cnpg.io/" + list.Items[0].Metadata.Name, "--timeout=1800s"}, "")
	return err
}
func (r *WorkloadRuntime) DataBackups(ctx context.Context, node models.ClusterNode, record *models.ClusterData) ([]map[string]any, error) {
	namespace, _ := DataIdentity(record.Spec)
	raw, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "get", "backups.postgresql.cnpg.io", "-o", "json"}, "")
	if err != nil {
		return nil, err
	}
	var list struct {
		Items []map[string]any `json:"items"`
	}
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return nil, err
	}
	_, name := DataIdentity(record.Spec)
	result := []map[string]any{}
	for _, item := range list.Items {
		spec, _ := item["spec"].(map[string]any)
		cluster, _ := spec["cluster"].(map[string]any)
		if cluster["name"] != name {
			continue
		}
		metadata, _ := item["metadata"].(map[string]any)
		status, _ := item["status"].(map[string]any)
		result = append(result, map[string]any{"name": metadata["name"], "uid": metadata["uid"], "phase": status["phase"], "startedAt": status["startedAt"], "stoppedAt": status["stoppedAt"], "error": status["error"]})
	}
	return result, nil
}
