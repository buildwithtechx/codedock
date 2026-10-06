package deploy

import (
	"context"
	"fmt"
	"strings"

	"codedock.run/codedock/internal/utils"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

func (d *Deployer) DependencyReady(ctx context.Context, source string) error {
	if utils.IsDryRun() {
		return nil
	}
	if d.containerManager.dockerClient == nil {
		return fmt.Errorf("runtime unavailable")
	}
	if after, ok := strings.CutPrefix(source, "app-"); ok {
		name := utils.NormalizeContainerName(after)
		inspected, err := d.containerManager.Inspect(ctx, name)
		if err != nil {
			inspected, err = d.containerManager.Inspect(ctx, name+"-1")
			if err != nil {
				return fmt.Errorf("prerequisite %s is unavailable: %w", source, err)
			}
		}
		if inspected.State == nil || !inspected.State.Running {
			return fmt.Errorf("prerequisite %s is not running", source)
		}
		if inspected.State.Health != nil && inspected.State.Health.Status != "healthy" {
			return fmt.Errorf("prerequisite %s is not healthy", source)
		}
		return nil
	}
	if after, ok := strings.CutPrefix(source, "db-"); ok {
		id := after
		current, err := d.containerManager.dockerClient.ContainerList(ctx, container.ListOptions{All: true, Filters: filters.NewArgs(filters.Arg("volume", "codedock-db-data-"+id))})
		if err != nil {
			return err
		}
		for _, database := range current {
			owned := database.Labels["codedock.database_id"] == id
			for _, name := range database.Names {
				if strings.HasPrefix(name, "/codedock-db-") {
					owned = true
				}
			}
			if owned && database.State == "running" {
				return nil
			}
		}
		return fmt.Errorf("database prerequisite %s is not running", source)
	}
	return fmt.Errorf("unsupported prerequisite")
}
