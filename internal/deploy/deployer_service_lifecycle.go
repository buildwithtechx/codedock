package deploy

import (
	"context"
	"fmt"
	"io"

	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/git/stacks"
	"codedock/internal/utils"
)

func (d *Deployer) StopAppService(ctx context.Context, app *appsmod.AppService) error {
	release, err := d.serviceOperation(ctx, app.ID)
	if err != nil {
		return err
	}
	defer release()
	if utils.IsDryRun() {
		return nil
	}

	if d.containerManager.dockerClient == nil {
		return fmt.Errorf("Docker runtime is unavailable")
	}

	containerName := utils.NormalizeContainerName(app.ID)
	var stopErr error
	if app.Replicas <= 1 {
		stopErr = d.containerManager.StopAndRemove(ctx, containerName)
	} else {
		for i := 1; i <= app.Replicas; i++ {
			replicaName := fmt.Sprintf("%s-%d", containerName, i)
			if err := d.containerManager.StopAndRemove(ctx, replicaName); err != nil && stopErr == nil {
				stopErr = err
			}
		}
		_ = d.containerManager.StopAndRemove(ctx, containerName)
	}
	if stopErr != nil {
		return fmt.Errorf("failed stopping app service: %w", stopErr)
	}

	app.Status = appsmod.AppServiceStatusStopped
	return d.store.UpdateAppService(app)
}

func (d *Deployer) StreamServiceLogs(ctx context.Context, app *appsmod.AppService, out io.Writer) error {
	containerName := utils.NormalizeContainerName(app.ID)
	return d.containerManager.StreamLogs(ctx, containerName, out)
}

func (d *Deployer) InspectServiceContainer(ctx context.Context, app *appsmod.AppService) (map[string]any, error) {
	containerName := utils.NormalizeContainerName(app.ID)
	inspect, err := d.containerManager.Inspect(ctx, containerName)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"id":         inspect.ID,
		"name":       inspect.Name,
		"status":     inspect.State.Status,
		"running":    inspect.State.Running,
		"exit_code":  inspect.State.ExitCode,
		"started_at": inspect.State.StartedAt,
		"ip_address": inspect.NetworkSettings.IPAddress,
	}, nil
}

func resolvedStartCommand(app *appsmod.AppService, sourceDir string) string {
	if app.StartCommand == "" || !app.Optimize {
		return app.StartCommand
	}
	if def, ok := stacks.Registry[stacks.DetectStack(sourceDir)]; ok && def.OptimizeCommand != "" {
		return stacks.ResolveWebCommand(app.StartCommand, true, def.OptimizeCommand)
	}
	return app.StartCommand
}

func containerCommand(app *appsmod.AppService, command string) []string {
	if command == "" || app.StaticOutput != "" && app.ImageRef == "" {
		return nil
	}
	return []string{"/bin/sh", "-c", command}
}
