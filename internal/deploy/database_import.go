package deploy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"

	"codedock/internal/modules/databases"
	"codedock/internal/utils"
)

func (d *DatabaseDeployer) ImportData(ctx context.Context, dbConfig *databases.Database, sourceURL string) error {
	if d.dockerClient == nil {
		return fmt.Errorf("docker daemon connection is not available")
	}
	release, err := d.acquireVolume(dbConfig.ID)
	if err != nil {
		return err
	}
	defer release()
	if utils.IsDryRun() {
		return nil
	}

	containerName := utils.NormalizeContainerName(fmt.Sprintf("codedock-db-%s", dbConfig.Name))

	switch strings.ToLower(string(dbConfig.Engine)) {
	case "postgres":
		cmd := fmt.Sprintf("pg_dump -Fc \"%s\" | pg_restore -U %s -d %s -1", sourceURL, dbConfig.Username, dbConfig.DatabaseName)
		execConfig := container.ExecOptions{
			Cmd:          []string{"sh", "-c", cmd},
			AttachStderr: true,
			AttachStdout: true,
		}
		execID, err := d.dockerClient.ContainerExecCreate(ctx, containerName, execConfig)
		if err != nil {
			return fmt.Errorf("failed to create exec for pg_dump: %w", err)
		}
		err = d.dockerClient.ContainerExecStart(ctx, execID.ID, container.ExecStartOptions{})
		if err != nil {
			return fmt.Errorf("failed to start exec for pg_dump: %w", err)
		}
		for {
			inspect, err := d.dockerClient.ContainerExecInspect(ctx, execID.ID)
			if err != nil {
				return fmt.Errorf("failed to inspect exec: %w", err)
			}
			if !inspect.Running {
				if inspect.ExitCode != 0 {
					return fmt.Errorf("pg_dump/restore failed with exit code %d", inspect.ExitCode)
				}
				break
			}
			time.Sleep(1 * time.Second)
		}
		return nil

	case "redis":
		cmd := fmt.Sprintf("redis-cli -u \"%s\" --rdb /data/dump.rdb", sourceURL)
		execConfig := container.ExecOptions{
			Cmd:          []string{"sh", "-c", cmd},
			AttachStderr: true,
			AttachStdout: true,
		}
		execID, err := d.dockerClient.ContainerExecCreate(ctx, containerName, execConfig)
		if err != nil {
			return fmt.Errorf("failed to create exec for redis-cli: %w", err)
		}
		err = d.dockerClient.ContainerExecStart(ctx, execID.ID, container.ExecStartOptions{})
		if err != nil {
			return fmt.Errorf("failed to start exec for redis-cli: %w", err)
		}
		for {
			inspect, err := d.dockerClient.ContainerExecInspect(ctx, execID.ID)
			if err != nil {
				return fmt.Errorf("failed to inspect exec: %w", err)
			}
			if !inspect.Running {
				if inspect.ExitCode != 0 {
					return fmt.Errorf("redis-cli failed with exit code %d", inspect.ExitCode)
				}
				break
			}
			time.Sleep(1 * time.Second)
		}
		return d.dockerClient.ContainerRestart(ctx, containerName, container.StopOptions{})

	default:
		return fmt.Errorf("data import not supported for engine %s", dbConfig.Engine)
	}
}
