package runtimes

import (
	"codedock.run/codedock/internal/engine/kubernetes"
	"codedock.run/codedock/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type DesiredStore interface {
	DesiredIDs(context.Context) ([]string, error)
	Desired(context.Context, string) (*models.DesiredRuntime, error)
	CommitDesired(context.Context, string, int, *models.DesiredRuntime) error
}

func (s *Service) RunReconciler(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			store, ok := s.store.(DesiredStore)
			if !ok {
				return
			}
			ids, err := store.DesiredIDs(ctx)
			if err != nil {
				slog.Error("list desired workloads", "error", err)
				continue
			}
			for _, id := range ids {
				operation, cancel := context.WithTimeout(ctx, 6*time.Minute)
				err := s.Reconcile(operation, id)
				cancel()
				if err != nil {
					slog.Warn("reconcile workload", "service", id, "error", err)
				}
			}
		}
	}
}
func (s *Service) Reconcile(ctx context.Context, id string) error {
	store, ok := s.store.(DesiredStore)
	if !ok {
		return fmt.Errorf("desired workload storage unavailable")
	}
	operation, release, err := s.builder.BeginServiceOperation(ctx, id)
	if err != nil {
		return err
	}
	defer release()
	ctx = operation
	runtime, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if runtime.Target.Kind != "kubernetes" || runtime.Journal != "" {
		return nil
	}
	desired, err := store.Desired(ctx, id)
	if err != nil {
		return err
	}
	if desired.Revision != runtime.Revision {
		return nil
	}
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if desired.Workload.App.ProjectID != app.ProjectID || desired.Workload.App.EnvironmentID != app.EnvironmentID {
		return fmt.Errorf("saved workload ownership changed")
	}
	cluster, err := s.cluster(ctx, app, runtime.Target)
	if err != nil {
		return err
	}
	unlock, err := s.lockCluster(cluster)
	if err != nil {
		return err
	}
	defer unlock()
	if err := s.engine.VerifyCluster(ctx, cluster, true); err != nil {
		return err
	}
	desired.Workload.App.Replicas = app.Replicas
	if desired.Workload.App.Replicas < 1 {
		desired.Workload.App.Replicas = 1
	}
	if err := kubernetes.ValidateWorkload(&desired.Workload); err != nil {
		return err
	}
	node := cluster.Nodes[0]
	if err := s.engine.Namespace(ctx, node, app); err != nil {
		return err
	}
	journal, err := s.engine.Snapshot(ctx, node, &desired.Workload)
	if err != nil {
		return err
	}
	manifest, err := desiredReplicas(desired.Manifest, app)
	if err != nil {
		return err
	}
	if err := s.engine.Validate(ctx, node, manifest); err != nil {
		return err
	}
	if err := s.store.Begin(ctx, id, runtime.Revision, journal); err != nil {
		return err
	}
	if err := s.engine.Apply(ctx, node, manifest); err != nil {
		return s.recoverFailure(app, node, journal, err)
	}
	if app.Status != models.AppServiceStatusStopped {
		if err := s.engine.Ready(ctx, node, app); err != nil {
			return s.recoverFailure(app, node, journal, err)
		}
	}
	observed, err := s.engine.Observe(ctx, node, app)
	if err != nil {
		return s.recoverFailure(app, node, journal, err)
	}
	return s.store.Observe(ctx, id, observed.Status, observed.Error, true)
}
func desiredReplicas(manifest string, app *models.AppService) (string, error) {
	var document map[string]any
	if err := json.Unmarshal([]byte(manifest), &document); err != nil {
		return "", err
	}
	items, ok := document["items"].([]any)
	if !ok {
		return "", fmt.Errorf("invalid desired resource list")
	}
	found := false
	for _, raw := range items {
		object, ok := raw.(map[string]any)
		if !ok {
			return "", fmt.Errorf("invalid desired resource")
		}
		if object["kind"] != "Deployment" {
			continue
		}
		metadata, _ := object["metadata"].(map[string]any)
		labels, _ := metadata["labels"].(map[string]any)
		if labels["codedock.run/service"] != app.ID || labels["codedock.run/project"] != app.ProjectID {
			return "", fmt.Errorf("desired workload ownership mismatch")
		}
		spec, ok := object["spec"].(map[string]any)
		if !ok {
			return "", fmt.Errorf("invalid desired workload")
		}
		replicas := app.Replicas
		if replicas < 1 {
			replicas = 1
		}
		if replicas > 64 {
			return "", fmt.Errorf("invalid desired replica count")
		}
		if app.Status == models.AppServiceStatusStopped {
			replicas = 0
		}
		spec["replicas"] = replicas
		found = true
	}
	if !found {
		return "", fmt.Errorf("desired deployment missing")
	}
	data, err := json.Marshal(document)
	return string(data), err
}
