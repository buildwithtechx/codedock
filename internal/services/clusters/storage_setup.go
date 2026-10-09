package clusters

import (
	"codedock/internal/engine/kubernetes"
	"codedock/internal/models"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

var allowedStorageVersions = map[string]string{
	"v1.6.2": "https://raw.githubusercontent.com/longhorn/longhorn/v1.6.2/deploy/longhorn.yaml",
	"v1.7.0": "https://raw.githubusercontent.com/longhorn/longhorn/v1.7.0/deploy/longhorn.yaml",
}

func (s *Service) ReviewStorageSetup(ctx context.Context, user, project, id string, request models.StorageSetupRequest) (*models.OperationReview, error) {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if cluster.ProjectID != project || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return nil, fmt.Errorf("select a ready cluster in this project")
	}
	url, ok := allowedStorageVersions[request.Version]
	if !ok || !regexp.MustCompile(`^v1\.[67]\.[0-9]+$`).MatchString(request.Version) {
		return nil, fmt.Errorf("select a reviewed Longhorn release")
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := s.http.Do(httpRequest)
	if err != nil {
		return nil, fmt.Errorf("download storage installer: %w", err)
	}
	defer response.Body.Close()
	data, err := io.ReadAll(io.LimitReader(response.Body, 8*1024*1024))
	if err != nil || response.StatusCode != http.StatusOK || len(data) == 0 {
		return nil, fmt.Errorf("storage installer unavailable")
	}
	sum := sha256.Sum256(data)
	plan := models.StorageSetupPlan{ClusterID: id, ProjectID: project, Version: request.Version, Manifest: string(data), SHA256: hex.EncodeToString(sum[:]), Action: "setup"}
	payload, err := json.Marshal(plan)
	if err != nil {
		return nil, err
	}
	effects := fmt.Sprintf("Install Longhorn %s as the shared storage provisioner with driver.longhorn.io and default RWX support. Existing PVCs are preserved; new file volumes can snapshot and restore. Installer SHA256: %s.", request.Version, plan.SHA256)
	return s.operations.Review(ctx, user, project, "cluster-storage", "cluster:"+id, string(payload), clusterSnapshot(cluster), effects)
}

func (s *Service) ApplyStorageSetup(ctx context.Context, user, id, confirmation string) error {
	op, err := s.operations.Get(ctx, id)
	if err != nil {
		return err
	}
	if op.Kind != "cluster-storage" || op.UserID != user {
		return fmt.Errorf("storage operation not found")
	}
	var plan models.StorageSetupPlan
	if err := json.Unmarshal([]byte(op.Payload), &plan); err != nil {
		return err
	}
	cluster, err := s.store.Get(ctx, plan.ClusterID)
	if err != nil {
		return err
	}
	if cluster.ProjectID != plan.ProjectID {
		return fmt.Errorf("cluster project changed")
	}
	sum := sha256.Sum256([]byte(plan.Manifest))
	if hex.EncodeToString(sum[:]) != plan.SHA256 {
		return fmt.Errorf("storage installer checksum mismatch")
	}
	return s.operations.Apply(ctx, id, user, confirmation, clusterSnapshot(cluster), func(ctx context.Context, op *models.Operation, progress func(string, string) error) error {
		release, err := s.gate.AcquireVolume("cluster:" + cluster.ID)
		if err != nil {
			return err
		}
		defer release()
		runtime := kubernetes.NewWorkloadRuntime(s.runner)
		return runtime.ApplyStorage(ctx, cluster.Nodes[0], plan.Manifest, progress)
	})
}

func (s *Service) SnapshotVolume(ctx context.Context, project, id, namespace, claim, snapshot string) error {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if cluster.ProjectID != project || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return fmt.Errorf("select a ready cluster in this project")
	}
	if !regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,62}$`).MatchString(snapshot) || !regexp.MustCompile(`^[a-z0-9][a-z0-9.-]{0,62}$`).MatchString(claim) {
		return fmt.Errorf("select a valid claim and snapshot name")
	}
	runtime := kubernetes.NewWorkloadRuntime(s.runner)
	return runtime.VolumeSnapshot(ctx, cluster.Nodes[0], namespace, claim, snapshot)
}

func (s *Service) RestoreVolume(ctx context.Context, project, id, namespace, snapshot, claim string) error {
	cluster, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if cluster.ProjectID != project || cluster.Status != "READY" || len(cluster.Nodes) == 0 {
		return fmt.Errorf("select a ready cluster in this project")
	}
	runtime := kubernetes.NewWorkloadRuntime(s.runner)
	return runtime.VolumeRestore(ctx, cluster.Nodes[0], namespace, snapshot, claim)
}
