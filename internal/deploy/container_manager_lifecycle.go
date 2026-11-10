package deploy

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

func (c *ContainerManager) StopAndRemove(ctx context.Context, containerIDOrName string) error {
	stopTimeout := 10
	_ = c.dockerClient.ContainerStop(ctx, containerIDOrName, container.StopOptions{Timeout: &stopTimeout})
	err := c.dockerClient.ContainerRemove(ctx, containerIDOrName, container.RemoveOptions{Force: true})
	if err != nil && !errdefs.IsNotFound(err) {
		return err
	}
	return nil
}

func (c *ContainerManager) Inspect(ctx context.Context, containerIDOrName string) (types.ContainerJSON, error) {
	return c.dockerClient.ContainerInspect(ctx, containerIDOrName)
}

func (c *ContainerManager) StreamLogs(ctx context.Context, containerIDOrName string, out io.Writer) error {
	logsReader, err := c.dockerClient.ContainerLogs(ctx, containerIDOrName, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "100",
	})
	if err != nil {
		return fmt.Errorf("failed to open container logs stream: %w", err)
	}
	defer logsReader.Close()
	_, err = io.Copy(out, logsReader)
	return err
}

func (c *ContainerManager) CleanupOrphanedContainers(ctx context.Context, prefix string, excludeContainerNames []string) error {
	containers, err := c.dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return fmt.Errorf("failed to list containers: %w", err)
	}

	excludeMap := make(map[string]bool)
	for _, n := range excludeContainerNames {
		excludeMap[n] = true
		excludeMap["/"+n] = true
	}

	for _, ctn := range containers {
		for _, name := range ctn.Names {
			if strings.HasPrefix(name, "/"+prefix+"-") {
				if !excludeMap[ctn.ID] && !excludeMap[name] {
					_ = c.StopAndRemove(ctx, ctn.ID)
				}
				break
			}
		}
	}
	return nil
}

func validateHostPath(path string, serviceID string) error {
	forbidden := []string{"/var/run/docker.sock", "/proc", "/sys", "/etc", "/root", "/boot"}
	for _, f := range forbidden {
		if strings.Contains(path, f) {
			return fmt.Errorf("forbidden path")
		}
	}
	expectedPrefix := fmt.Sprintf("/data/codedock/%s/", serviceID)
	if !strings.HasPrefix(path, expectedPrefix) {
		return fmt.Errorf("must start with %s", expectedPrefix)
	}
	return nil
}
