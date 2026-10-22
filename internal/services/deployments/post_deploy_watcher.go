package deployments

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"

	"codedock/internal/models"
)

func (s *DeploymentService) watchPostDeployStability(deploymentID, projectID, containerID string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("recovered panic in post deploy stability watcher", "deployment", deploymentID, "error", r)
			}
		}()

		time.Sleep(3 * time.Second) // allow container to finish initial initialization
		ctx, cancel := context.WithTimeout(context.Background(), 65*time.Second)
		defer cancel()

		docker, release, err := s.DockerForProject(ctx, projectID)
		if err != nil {
			return
		}
		defer release()

		ticker := time.NewTicker(6 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				finding, unstable := sampleStability(ctx, docker, containerID)
				if unstable {
					logTail := fetchContainerLogTail(ctx, docker, containerID, 15)
					detail := finding
					if logTail != "" {
						detail = fmt.Sprintf("%s\nLast logs:\n%s", finding, logTail)
					}
					if err := s.repo.SetRestartLoop(context.Background(), deploymentID, detail); err != nil {
						slog.Warn("failed to record restart loop", "deployment", deploymentID, "error", err)
					}
					return
				}
			}
		}
	}()
}

func sampleStability(ctx context.Context, cli *client.Client, containerID string) (string, bool) {
	inspect, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", false
	}
	if inspect.State == nil {
		return "", false
	}

	if inspect.State.OOMKilled {
		return "Container was terminated by kernel Out-Of-Memory (OOM) killer", true
	}
	if inspect.RestartCount > 2 || inspect.State.Restarting {
		return fmt.Sprintf("Container in restart loop (restarted %d times, exit code %d)", inspect.RestartCount, inspect.State.ExitCode), true
	}
	if !inspect.State.Running && inspect.State.ExitCode != 0 {
		return fmt.Sprintf("Container stopped unexpectedly with exit code %d", inspect.State.ExitCode), true
	}
	return "", false
}

func fetchContainerLogTail(ctx context.Context, cli *client.Client, containerID string, lines int) string {
	reader, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       fmt.Sprintf("%d", lines),
		Timestamps: false,
	})
	if err != nil {
		return ""
	}
	defer reader.Close()

	var buf bytes.Buffer
	inspected, err := cli.ContainerInspect(ctx, containerID)
	if err == nil && inspected.Config != nil && inspected.Config.Tty {
		_, _ = io.Copy(&buf, reader)
	} else {
		_, _ = stdcopy.StdCopy(&buf, &buf, reader)
	}

	trimmed := strings.TrimSpace(buf.String())
	if len(trimmed) > 2000 {
		return trimmed[len(trimmed)-2000:]
	}
	return trimmed
}

// CheckDeploymentRestartLoop returns the recorded restart loop or empty string.
func CheckDeploymentRestartLoop(dep *models.Deployment) string {
	if dep == nil {
		return ""
	}
	return dep.RestartLoop
}
