package observability

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"github.com/containerd/errdefs"
	"github.com/docker/docker/client"
	"time"
)

type ClusterObservations interface {
	Handles(context.Context, string) (bool, error)
	Observe(context.Context, string) (*models.WorkloadObservation, error)
	Logs(context.Context, string) (string, error)
}
type CanvasDockerTargets interface {
	DockerForProject(context.Context, string) (*client.Client, func(), error)
}
type CanvasRuntime struct {
	Targets CanvasDockerTargets
	docker  *client.Client
	Cluster ClusterObservations
}

func NewCanvasRuntime(docker *client.Client) *CanvasRuntime { return &CanvasRuntime{docker: docker} }

func (r *CanvasRuntime) Observe(ctx context.Context, canvas *models.EnvironmentCanvas) {
	if canvas == nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	docker, release, targetErr := r.canvasDocker(ctx, canvas)
	defer release()
	status := func(id string) string {
		if id == "" {
			return "not deployed"
		}
		if targetErr != nil || docker == nil {
			return "runtime unavailable"
		}
		inspected, err := docker.ContainerInspect(ctx, id)
		if errdefs.IsNotFound(err) {
			return "container missing"
		}
		if err != nil || inspected.State == nil {
			return "runtime unavailable"
		}
		if inspected.State.Health != nil && inspected.State.Health.Status != "healthy" {
			return inspected.State.Health.Status
		}
		return inspected.State.Status
	}
	for index := range canvas.Nodes {
		node := &canvas.Nodes[index]
		if id, ok := node.Data["containerId"].(string); ok {
			node.Data["status"] = status(id)
			node.Data["label"] = fmt.Sprint(node.Data["label"]) + " / " + fmt.Sprint(node.Data["status"])
		}
		for _, app := range canvas.Apps {
			if node.ID == "app-"+app.ID {
				observedStatus := status(app.ContainerID)
				if r.Cluster != nil {
					handles, err := r.Cluster.Handles(ctx, app.ID)
					if err != nil {
						observedStatus = "runtime unavailable"
					} else if handles {
						observed, err := r.Cluster.Observe(ctx, app.ID)
						if err != nil {
							observedStatus = "runtime unavailable"
							node.Data["runtimeError"] = err.Error()
						} else {
							observedStatus = observed.Status
							node.Data["runtimeKind"] = observed.Kind
							node.Data["availableReplicas"] = observed.Available
							node.Data["desiredReplicas"] = observed.Desired
						}
					}
				}
				node.Data["status"] = observedStatus
				node.Data["label"] = app.Name + " / " + observedStatus
			}
		}
		for _, database := range canvas.Databases {
			if node.ID == "db-"+database.ID {
				node.Data["status"] = status(database.ContainerID)
			}
		}
		for _, record := range canvas.ClusterDatabases {
			if node.ID == "clusterdb-"+record.ID {
				node.Data["status"] = record.Status
				if record.Error != "" {
					node.Data["runtimeError"] = record.Error
				}
				node.Data["label"] = record.Spec.Name + " / " + record.Status
			}
		}
	}
}

func (r *CanvasRuntime) canvasDocker(ctx context.Context, canvas *models.EnvironmentCanvas) (*client.Client, func(), error) {
	if r.Targets == nil {
		return r.docker, func() {}, nil
	}
	if canvas == nil || canvas.Environment == nil || canvas.Environment.ProjectID == "" {
		return nil, func() {}, fmt.Errorf("canvas project identity unavailable")
	}
	target, release, err := r.Targets.DockerForProject(ctx, canvas.Environment.ProjectID)
	if err != nil {
		return nil, func() {}, err
	}
	return target, release, nil
}
