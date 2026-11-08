package observability

import (
	"bytes"
	"codedock/internal/models"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"time"
)

func (r *CanvasRuntime) Read(ctx context.Context, canvas *models.EnvironmentCanvas, nodeID string) (*models.RuntimeObservation, error) {
	if canvas == nil {
		return nil, fmt.Errorf("canvas unavailable")
	}
	if r.Cluster != nil {
		for _, app := range canvas.Apps {
			if nodeID != "app-"+app.ID {
				continue
			}
			handles, err := r.Cluster.Handles(ctx, app.ID)
			if err != nil {
				return nil, err
			}
			if handles {
				observed, err := r.Cluster.Observe(ctx, app.ID)
				if err != nil {
					return nil, err
				}
				logs, err := r.Cluster.Logs(ctx, app.ID)
				if err != nil {
					return nil, err
				}
				return &models.RuntimeObservation{Logs: logs, Metrics: map[string]any{"runtimeKind": observed.Kind, "status": observed.Status, "desiredReplicas": observed.Desired, "availableReplicas": observed.Available, "metricsAvailable": observed.MetricsAvailable, "metricsError": observed.MetricsError, "pods": observed.Pods}}, nil
			}
		}
	}
	docker, release, err := r.canvasDocker(ctx, canvas)
	if err != nil {
		return nil, err
	}
	defer release()
	if docker == nil {
		return nil, fmt.Errorf("runtime unavailable")
	}
	containerID := ""
	for _, node := range canvas.Nodes {
		if node.ID == nodeID {
			if id, ok := node.Data["containerId"].(string); ok {
				containerID = id
			}
			if id, ok := node.Data["serviceId"].(string); ok {
				for _, app := range canvas.Apps {
					if app.ID == id {
						containerID = app.ContainerID
					}
				}
			}
			if id, ok := node.Data["databaseId"].(string); ok {
				for _, database := range canvas.Databases {
					if database.ID == id {
						containerID = database.ContainerID
					}
				}
			}
		}
	}
	if containerID == "" {
		return nil, fmt.Errorf("resource has no observed container")
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	health, err := NewStatsMonitor(docker).GetHealth(ctx, containerID)
	if err != nil {
		return nil, err
	}
	stream, err := docker.ContainerLogs(ctx, containerID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Tail: "100", Timestamps: true})
	if err != nil {
		return nil, fmt.Errorf("read runtime logs: %w", err)
	}
	defer stream.Close()
	inspected, err := docker.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	var logs bytes.Buffer
	reader := io.LimitReader(stream, 1024*1024)
	if inspected.Config != nil && inspected.Config.Tty {
		_, err = io.Copy(&logs, reader)
	} else {
		_, err = stdcopy.StdCopy(&logs, &logs, reader)
	}
	if err != nil {
		return nil, fmt.Errorf("decode runtime logs: %w", err)
	}
	return &models.RuntimeObservation{Logs: logs.String(), Metrics: map[string]any{"status": health.Status, "cpuPercent": health.CPUUsagePercentage, "memoryBytes": health.MemoryUsageBytes, "memoryLimitBytes": health.MemoryLimitBytes, "uptimeSeconds": health.UptimeSeconds}}, nil
}
