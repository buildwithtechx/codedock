package deploy

import (
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func (d *DatabaseDeployer) Stop(ctx context.Context, id string) error {
	release, err := d.acquireVolume(id)
	if err != nil {
		return err
	}
	defer release()
	db, err := d.store.GetDatabase(id)
	if err != nil {
		return fmt.Errorf("load database before stopping: %w", err)
	}
	if db == nil {
		return utils.NewNotFoundError("Database", id)
	}
	if utils.IsDryRun() {
		return d.store.UpdateDatabaseStatus(id, models.DatabaseStatusStopped, db.ContainerID)
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
	return d.store.UpdateDatabaseStatus(id, models.DatabaseStatusStopped, db.ContainerID)
}
