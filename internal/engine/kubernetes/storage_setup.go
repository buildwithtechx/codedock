package kubernetes

import (
	"codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
)

func (r *WorkloadRuntime) ApplyStorage(ctx context.Context, node models.ClusterNode, manifest string, progress func(string, string) error) error {
	var document map[string]any
	if err := json.Unmarshal([]byte(manifest), &document); err != nil {
		items := map[string]any{}
		if err := json.Unmarshal([]byte(manifest), &items); err != nil {
			return fmt.Errorf("invalid storage manifest")
		}
	}
	if _, err := r.commands.Kubectl(ctx, node, []string{"apply", "--server-side", "--field-manager=codedock-storage", "-f", "-"}, manifest); err != nil {
		return err
	}
	if err := progress("APPLYING", "Storage provisioner manifests applied"); err != nil {
		return err
	}
	for _, deployment := range []string{"longhorn-driver-deployer", "longhorn-manager"} {
		if _, err := r.commands.Kubectl(ctx, node, []string{"-n", "longhorn-system", "rollout", "status", "deployment/" + deployment, "--timeout=600s"}, ""); err != nil {
			return err
		}
	}
	return progress("VERIFIED", "Shared storage provisioner is ready")
}

func (r *WorkloadRuntime) VolumeSnapshot(ctx context.Context, node models.ClusterNode, namespace, claim, snapshot string) error {
	object := map[string]any{"apiVersion": "snapshot.storage.k8s.io/v1", "kind": "VolumeSnapshot", "metadata": map[string]any{"name": snapshot, "namespace": namespace, "labels": map[string]string{"codedock.run/snapshot": snapshot}}, "spec": map[string]any{"volumeSnapshotClassName": "longhorn", "source": map[string]string{"persistentVolumeClaimName": claim}}}
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	if _, err := r.commands.Kubectl(ctx, node, []string{"-n", namespace, "apply", "--server-side", "--field-manager=codedock-storage", "-f", "-"}, string(data)); err != nil {
		return err
	}
	_, err = r.commands.Kubectl(ctx, node, []string{"-n", namespace, "wait", "--for=jsonpath={.status.readyToUse}=true", "volumesnapshot/" + snapshot, "--timeout=600s"}, "")
	return err
}

func (r *WorkloadRuntime) VolumeRestore(ctx context.Context, node models.ClusterNode, namespace, snapshot, claim string) error {
	object := map[string]any{"apiVersion": "v1", "kind": "PersistentVolumeClaim", "metadata": map[string]any{"name": claim, "namespace": namespace, "labels": map[string]string{"codedock.run/restored": snapshot}}, "spec": map[string]any{"accessModes": []string{"ReadWriteOnce"}, "storageClassName": "longhorn", "dataSource": map[string]string{"name": snapshot, "kind": "VolumeSnapshot", "apiGroup": "snapshot.storage.k8s.io"}, "resources": map[string]any{"requests": map[string]string{"storage": "1Gi"}}}}
	data, err := json.Marshal(object)
	if err != nil {
		return err
	}
	_, err = r.commands.Kubectl(ctx, node, []string{"-n", namespace, "apply", "--server-side", "--field-manager=codedock-storage", "-f", "-"}, string(data))
	return err
}
