package deploy

import (
	"codedock.run/codedock/internal/utils"
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func (d *DatabaseDeployer) Remove(ctx context.Context, id string) error {
	if utils.IsDryRun() {
		return nil
	}
	if d.dockerClient == nil {
		return fmt.Errorf("Docker is unavailable")
	}
	db, err := d.store.GetDatabase(id)
	if err != nil {
		return fmt.Errorf("load database for removal: %w", err)
	}
	release, err := d.acquireVolume(id)
	if err != nil {
		return err
	}
	defer release()
	name := db.ContainerID
	if name == "" {
		name = utils.NormalizeContainerName("codedock-db-" + db.Name)
	}
	if err := d.dockerClient.ContainerRemove(ctx, name, container.RemoveOptions{Force: true}); err != nil && !client.IsErrNotFound(err) {
		return fmt.Errorf("remove database container: %w", err)
	}
	return nil
}
