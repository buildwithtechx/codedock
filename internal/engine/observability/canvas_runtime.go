package observability

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"github.com/containerd/errdefs"
	"github.com/docker/docker/client"
	"time"
)

type CanvasRuntime struct{ docker *client.Client }

func NewCanvasRuntime(docker *client.Client) *CanvasRuntime { return &CanvasRuntime{docker: docker} }

func (r *CanvasRuntime) Observe(ctx context.Context, canvas *models.EnvironmentCanvas) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	status := func(id string) string {
		if id == "" {
			return "not deployed"
		}
		if r.docker == nil {
			return "runtime unavailable"
		}
		inspected, err := r.docker.ContainerInspect(ctx, id)
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
				node.Data["status"] = status(app.ContainerID)
			}
		}
		for _, database := range canvas.Databases {
			if node.ID == "db-"+database.ID {
				node.Data["status"] = status(database.ContainerID)
			}
		}
	}
}
