package runtimes

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"strings"
)

func (s *Service) target(ctx context.Context, id string) (*models.AppService, *models.Cluster, error) {
	app, err := s.apps.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	runtime, err := s.store.Get(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if runtime.Target.Kind != "kubernetes" {
		return nil, nil, fmt.Errorf("application uses Docker")
	}
	cluster, err := s.cluster(ctx, app, runtime.Target)
	return app, cluster, err
}
func (s *Service) Observe(ctx context.Context, id string) (*models.WorkloadObservation, error) {
	nativeApp, nativeRuntime, nativeErr := s.nativeTarget(ctx, id)
	if nativeErr != nil {
		return nil, nativeErr
	}
	if nativeApp != nil {
		if !strings.HasPrefix(nativeApp.ContainerID, "bare:") {
			return &models.WorkloadObservation{Kind: "bare", Status: "NOT_DEPLOYED", Pods: []models.RuntimePod{}}, nil
		}
		return s.native.Observe(ctx, nativeApp, nativeRuntime.Target)
	}
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.engine.Observe(ctx, cluster.Nodes[0], app)
}
func (s *Service) Logs(ctx context.Context, id string) (string, error) {
	nativeApp, nativeRuntime, nativeErr := s.nativeTarget(ctx, id)
	if nativeErr != nil {
		return "", nativeErr
	}
	if nativeApp != nil {
		return s.native.Logs(ctx, nativeApp, nativeRuntime.Target)
	}
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return "", err
	}
	return s.engine.Logs(ctx, cluster.Nodes[0], app)
}
func (s *Service) Exec(ctx context.Context, id string, request models.RuntimeExecRequest) (string, error) {
	nativeApp, nativeRuntime, nativeErr := s.nativeTarget(ctx, id)
	if nativeErr != nil {
		return "", nativeErr
	}
	if nativeApp != nil {
		return s.native.Exec(ctx, nativeApp, nativeRuntime.Target, request)
	}
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return "", err
	}
	return s.engine.Exec(ctx, cluster.Nodes[0], app, request)
}
func (s *Service) Lifecycle(ctx context.Context, id, action string, replicas int) error {
	operation, release, err := s.builder.BeginServiceOperation(ctx, id)
	if err != nil {
		return err
	}
	defer release()
	ctx = operation
	nativeApp, nativeRuntime, nativeErr := s.nativeTarget(ctx, id)
	if nativeErr != nil {
		return nativeErr
	}
	if nativeApp != nil {
		return s.nativeLifecycle(ctx, nativeApp, nativeRuntime, action)
	}
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return err
	}
	unlock, err := s.lockCluster(cluster)
	if err != nil {
		return err
	}
	defer unlock()
	if err := s.engine.VerifyCluster(ctx, cluster, action != "stop"); err != nil {
		return err
	}
	runtime, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if runtime.Journal != "" {
		return fmt.Errorf("recover the interrupted deployment before changing the workload")
	}
	if action == "scale" && replicas > 1 {
		for _, volume := range runtime.Target.Volumes {
			if !volume.Shared {
				return fmt.Errorf("multiple replicas require shared storage")
			}
		}
	}
	if err := s.engine.Lifecycle(ctx, cluster.Nodes[0], app, action, replicas); err != nil {
		return err
	}
	if action == "stop" {
		app.Status = models.AppServiceStatusStopped
	} else {
		app.Status = models.AppServiceStatusRunning
	}
	if action == "scale" {
		app.Replicas = replicas
	}
	return s.apps.Update(ctx, app)
}

func (s *Service) Remove(ctx context.Context, id string) error {
	operation, release, err := s.builder.BeginServiceOperation(ctx, id)
	if err != nil {
		return err
	}
	defer release()
	ctx = operation
	nativeApp, nativeRuntime, nativeErr := s.nativeTarget(ctx, id)
	if nativeErr != nil {
		return nativeErr
	}
	if nativeApp != nil {
		return s.nativeLifecycle(ctx, nativeApp, nativeRuntime, "remove")
	}
	app, cluster, err := s.target(ctx, id)
	if err != nil {
		return err
	}
	unlock, err := s.lockCluster(cluster)
	if err != nil {
		return err
	}
	defer unlock()
	runtime, err := s.store.Get(ctx, id)
	if err != nil {
		return err
	}
	if runtime.Journal != "" {
		return fmt.Errorf("complete deployment recovery before deleting this application")
	}
	if err := s.engine.VerifyCluster(ctx, cluster, false); err != nil {
		return err
	}
	return s.engine.Remove(ctx, cluster.Nodes[0], app)
}
