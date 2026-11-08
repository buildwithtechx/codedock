package databases

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"codedock/internal/utils"
)

func (d *dockerDeployer) Stop(ctx context.Context, id string) error {
	release, err := d.acquireVolume(id)
	if err != nil {
		return err
	}
	defer release()
	db, err := d.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("load database before stopping: %w", err)
	}
	if db == nil {
		return utils.NewNotFoundError("Database", id)
	}
	if utils.IsDryRun() {
		return d.repo.UpdateStatus(ctx, id, DatabaseStatusStopped, db.ContainerID)
	}
	if d.dockerClient == nil {
		return fmt.Errorf("Docker is unavailable")
	}
	name := db.ContainerID
	if name == "" {
		name = utils.NormalizeContainerName("codedock-db-" + db.Name)
	}
	_, err = d.dockerClient.ContainerUpdate(ctx, name, container.UpdateConfig{RestartPolicy: container.RestartPolicy{Name: "unless-stopped"}})
	if err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("set stopped database restart policy: %w", err)
	}
	if err := d.dockerClient.ContainerStop(ctx, name, container.StopOptions{}); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("stop database container: %w", err)
	}
	return d.repo.UpdateStatus(ctx, id, DatabaseStatusStopped, db.ContainerID)
}

func (d *dockerDeployer) Remove(ctx context.Context, id string, deleteRecord func() error) error {
	release, err := d.acquireVolume(id)
	if err != nil {
		return err
	}
	defer release()
	if utils.IsDryRun() {
		return deleteRecord()
	}
	if d.dockerClient == nil {
		return fmt.Errorf("Docker is unavailable")
	}
	db, err := d.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("load database for removal: %w", err)
	}
	if db == nil {
		return utils.NewNotFoundError("Database", id)
	}
	name := db.ContainerID
	if name == "" {
		name = utils.NormalizeContainerName("codedock-db-" + db.Name)
	}
	if err := d.dockerClient.ContainerStop(ctx, name, container.StopOptions{}); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("stop database before removal: %w", err)
	}
	if err := d.dockerClient.ContainerRemove(ctx, name, container.RemoveOptions{}); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("remove database container: %w", err)
	}
	if err := d.dockerClient.VolumeRemove(ctx, "codedock-db-data-"+id, false); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("remove database volume: %w", err)
	}
	return deleteRecord()
}

func (d *dockerDeployer) ImportData(ctx context.Context, dbConfig *Database, sourceURL string) error {
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
		if err := d.dockerClient.ContainerExecStart(ctx, execID.ID, container.ExecStartOptions{}); err != nil {
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
		if err := d.dockerClient.ContainerExecStart(ctx, execID.ID, container.ExecStartOptions{}); err != nil {
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
