package deploy

import (
	"codedock.run/codedock/internal/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func (d *DatabaseDeployer) Remove(ctx context.Context, id string, deleteRecord func() error) error {
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
	db, err := d.store.GetDatabase(id)
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
