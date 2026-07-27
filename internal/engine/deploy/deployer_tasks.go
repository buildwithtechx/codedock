package deploy

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/errdefs"
	"github.com/google/uuid"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func (d *Deployer) ExecuteRollingUpdate(ctx context.Context, app *models.AppService, newImageTag string, logWriter io.Writer) error {
	replicas := app.Replicas
	if replicas <= 1 {
		appCopy := *app
		if newImageTag != "" {
			appCopy.ImageRef = newImageTag
		}
		_, err := d.DeployAppService(ctx, &appCopy, "", logWriter)
		return err
	}

	primaryName := utils.NormalizeContainerName(app.ID)
	logDrains, _ := d.store.ListLogDrainsByService(app.ID)

	internalPort := app.InternalPort
	if internalPort <= 0 {
		internalPort = defaultAppPort()
	}
	memoryLimit := app.MemoryLimit
	if memoryLimit <= 0 {
		memoryLimit = defaultMemoryMB()
	}
	cpuRequest := app.CPULimit
	if cpuRequest <= 0 {
		cpuRequest = defaultCPURequest()
	}

	for i := 0; i < replicas; i++ {
		containerName := fmt.Sprintf("%s-%d", primaryName, i+1)
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🔄 [Deployer] Rolling update: replacing replica %d/%d (%s)...\n", i+1, replicas, containerName)
		}

		envVarsMap, _ := d.getEnvironmentVariables(app, logWriter)
		var envSlice []string
		for k, v := range envVarsMap {
			envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
		}

		effectiveImageTag := newImageTag
		if effectiveImageTag == "" {
			effectiveImageTag = app.ImageRef
		}
		if effectiveImageTag == "" {
			effectiveImageTag = fmt.Sprintf("codedock-app-%s:latest", app.ID)
		}

		runOpts := ContainerRunOptions{
			Name:            containerName,
			ImageTag:        effectiveImageTag,
			ServiceID:       app.ID,
			Domain:          app.Domain,
			InternalPort:    internalPort,
			RuntimeMode:     app.RuntimeMode,
			Envs:            envSlice,
			Cmd:             strings.Fields(app.StartCommand),
			MemoryLimitMB:   memoryLimit,
			CPURequest:      cpuRequest,
			HealthCheckPath: app.HealthCheckPath,
			Volumes:         app.Volumes,
			MaintenanceMode: app.MaintenanceMode,
			LogDrains:       logDrains,
		}

		newID, err := d.containerManager.CreateAndStart(ctx, runOpts)
		if err != nil {
			return fmt.Errorf("rolling update failed on replica %d: %w", i+1, err)
		}

		appCopy := *app
		appCopy.InternalPort = internalPort
		if err := d.verifyHealthCheck(ctx, &appCopy, containerName, logWriter); err != nil {
			return fmt.Errorf("health check failed on replica %d during rolling update: %w", i+1, err)
		}

		if i == 0 {
			app.ContainerID = newID
			_ = d.store.UpdateAppService(app)
		}
	}

	return nil
}

func (d *Deployer) ExecuteOneOffTask(ctx context.Context, app *models.AppService, command string, logWriter io.Writer) error {
	containerName := fmt.Sprintf("%s-task-%s", utils.NormalizeContainerName(app.ID), uuid.NewString()[:6])
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🏃 [Deployer] Running one-off task in container %s: %s\n", containerName, command)
	}

	envVarsMap, _ := d.getEnvironmentVariables(app, logWriter)
	var envSlice []string
	for k, v := range envVarsMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	imageTag := app.ImageRef
	if imageTag == "" {
		primaryName := utils.NormalizeContainerName(app.ID)
		if inspect, err := d.containerManager.Inspect(ctx, primaryName); err == nil && inspect.Config != nil {
			imageTag = inspect.Config.Image
		} else if inspect, err := d.containerManager.Inspect(ctx, fmt.Sprintf("%s-1", primaryName)); err == nil && inspect.Config != nil {
			imageTag = inspect.Config.Image
		} else if app.ContainerID != "" {
			if inspect, err := d.containerManager.Inspect(ctx, app.ContainerID); err == nil && inspect.Config != nil {
				imageTag = inspect.Config.Image
			}
		}
		if imageTag == "" {
			imageTag = fmt.Sprintf("codedock-app-%s:latest", app.ID)
		}
	}

	runOpts := ContainerRunOptions{
		Name:            containerName,
		ImageTag:        imageTag,
		ServiceID:       app.ID,
		RuntimeMode:     models.RuntimeModeWorker,
		Envs:            envSlice,
		Cmd:             strings.Fields(command),
		MemoryLimitMB:   app.MemoryLimit,
		CPURequest:      app.CPULimit,
		MaintenanceMode: false,
	}

	cid, err := d.containerManager.CreateAndStart(ctx, runOpts)
	if err != nil {
		return fmt.Errorf("failed starting task container: %w", err)
	}

	defer d.containerManager.StopAndRemove(ctx, cid)

	if logWriter != nil {
		_ = d.containerManager.StreamLogs(ctx, cid, logWriter)
	}

	statusCh, errCh := d.containerManager.dockerClient.ContainerWait(ctx, cid, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("error waiting for task container: %w", err)
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("task container exited with non-zero code %d", status.StatusCode)
		}
	}

	return nil
}

func (d *Deployer) RestartAppService(ctx context.Context, app *models.AppService) error {
	containerName := utils.NormalizeContainerName(app.ID)
	replicas := app.Replicas
	if replicas <= 1 {
		inspect, err := d.containerManager.Inspect(ctx, containerName)
		if err != nil {
			if errdefs.IsNotFound(err) {
				return utils.NewNotFoundError("Container", containerName)
			}
			return fmt.Errorf("failed to inspect container for restart: %w", err)
		}
		if err := d.containerManager.dockerClient.ContainerRestart(ctx, inspect.ID, container.StopOptions{}); err != nil {
			return fmt.Errorf("failed to restart container %s: %w", containerName, err)
		}
	} else {
		var restartErrs []string
		for i := 1; i <= replicas; i++ {
			repName := fmt.Sprintf("%s-%d", containerName, i)
			inspect, err := d.containerManager.Inspect(ctx, repName)
			if err != nil {
				restartErrs = append(restartErrs, fmt.Sprintf("%s: %v", repName, err))
				continue
			}
			if err := d.containerManager.dockerClient.ContainerRestart(ctx, inspect.ID, container.StopOptions{}); err != nil {
				restartErrs = append(restartErrs, fmt.Sprintf("%s: %v", repName, err))
			}
		}
		if len(restartErrs) > 0 {
			return fmt.Errorf("failed restarting replicas: %s", strings.Join(restartErrs, "; "))
		}
	}

	app.Status = models.AppServiceStatusRunning
	return d.store.UpdateAppService(app)
}
