package deploy

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/google/uuid"

	"codedock.run/codedock/internal/engine/build"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

type Deployer struct {
	builder          build.Builder
	containerManager *ContainerManager
	store            DeployerStore
	EnvProvider      func(projectID string) (map[string]string, error)
	EnvInterpolator  func(projectID string) (map[string]map[string]string, error)
}

func NewDeployer(dockerClient *client.Client, s DeployerStore) *Deployer {
	return &Deployer{
		builder:          build.NewBuilder(dockerClient),
		containerManager: NewContainerManager(dockerClient, s),
		store:            s,
	}
}

func (d *Deployer) DeployAppService(ctx context.Context, app *models.AppService, sourceDir string, logWriter io.Writer) (string, error) {
	if logWriter != nil {
		fmt.Fprintf(logWriter, "🚀 [Deployer] Starting deployment for service: %s (ID: %s)\n", app.Name, app.ID)
	}

	if utils.IsDryRun() {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "🚀 [Deployer] Dry-run mode is enabled. Skipping actual build and run steps.\n")
		}
		newContainerName := fmt.Sprintf("%s-dryrun", utils.NormalizeContainerName(app.ID))
		return newContainerName, nil
	}

	if err := d.prepareServerlessCode(app, sourceDir, logWriter); err != nil {
		return "", err
	}

	envVarsMap, err := d.getEnvironmentVariables(app, logWriter)
	if err != nil {
		return "", fmt.Errorf("failed resolving service env vars: %w", err)
	}

	appDomain := app.Domain
	if appDomain == "" {
		appDomain = utils.GenerateAppDomain(app.Name, "", "")
	}

	var envSlice []string
	for k, v := range envVarsMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

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

	imageTag := fmt.Sprintf("codedock-app-%s:%s", app.ID, uuid.NewString()[:8])
	if app.ImageRef != "" {
		imageTag = app.ImageRef
		if app.RegistryID != nil && *app.RegistryID != "" {
			reg, err := d.store.GetRegistry(*app.RegistryID)
			if err == nil && reg != nil && d.containerManager != nil && d.containerManager.dockerClient != nil {
				authConfig := registry.AuthConfig{
					Username:      reg.Username,
					Password:      reg.PasswordToken,
					ServerAddress: reg.RegistryURL,
				}
				encodedJSON, err := json.Marshal(authConfig)
				if err == nil {
					authStr := base64.URLEncoding.EncodeToString(encodedJSON)
					out, err := d.containerManager.dockerClient.ImagePull(ctx, imageTag, image.PullOptions{RegistryAuth: authStr})
					if err == nil {
						io.Copy(io.Discard, out)
						out.Close()
					}
				}
			}
		}
	} else {
		var buildEnv map[string]string
		if app.BuildCommand != "" {
			buildEnv = map[string]string{"BUILD_COMMAND": app.BuildCommand}
		}

		buildOpts := build.BuildOptions{
			ProjectID:      app.ProjectID,
			ServiceID:      app.ID,
			SourceDir:      sourceDir,
			DockerfilePath: app.DockerfilePath,
			LogWriter:      logWriter,
			AppConfig:      app,
			EnvVars:        buildEnv,
		}

		builtTag, err := d.builder.Build(ctx, buildOpts)
		if err != nil {
			return "", fmt.Errorf("build failed: %w", err)
		}
		imageTag = builtTag
	}

	replicas := app.Replicas
	if replicas <= 0 {
		replicas = 1
	}

	primaryContainerName := utils.NormalizeContainerName(app.ID)
	logDrains, _ := d.store.ListLogDrainsByService(app.ID)

	for i := 0; i < replicas; i++ {
		containerName := primaryContainerName
		if replicas > 1 {
			containerName = fmt.Sprintf("%s-%d", primaryContainerName, i+1)
		}

		runOpts := ContainerRunOptions{
			Name:            containerName,
			ImageTag:        imageTag,
			ServiceID:       app.ID,
			Domain:          appDomain,
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

		containerID, err := d.containerManager.CreateAndStart(ctx, runOpts)
		if err != nil {
			return "", fmt.Errorf("failed starting container %s: %w", containerName, err)
		}

		if i == 0 {
			if err := d.verifyHealthCheck(ctx, app, containerName, logWriter); err != nil {
				return "", err
			}
			app.ContainerID = containerID
			app.Status = models.AppServiceStatusRunning
			app.Domain = appDomain
			_ = d.store.UpdateAppService(app)
		}
	}

	var excludeNames []string
	for i := 0; i < replicas; i++ {
		if replicas == 1 {
			excludeNames = append(excludeNames, primaryContainerName)
		} else {
			excludeNames = append(excludeNames, fmt.Sprintf("%s-%d", primaryContainerName, i+1))
		}
	}

	_ = d.containerManager.CleanupOrphanedContainers(ctx, primaryContainerName, excludeNames)

	return primaryContainerName, nil
}

func (d *Deployer) StopAppService(ctx context.Context, app *models.AppService) error {
	if utils.IsDryRun() {
		return nil
	}

	containerName := utils.NormalizeContainerName(app.ID)
	if err := d.containerManager.StopAndRemove(ctx, containerName); err != nil {
		return err
	}

	app.Status = models.AppServiceStatusStopped
	return d.store.UpdateAppService(app)
}

func (d *Deployer) StreamServiceLogs(ctx context.Context, app *models.AppService, out io.Writer) error {
	containerName := utils.NormalizeContainerName(app.ID)
	return d.containerManager.StreamLogs(ctx, containerName, out)
}

func (d *Deployer) InspectServiceContainer(ctx context.Context, app *models.AppService) (map[string]any, error) {
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

func (d *Deployer) ExecuteRollingUpdate(ctx context.Context, app *models.AppService, newImageTag string, logWriter io.Writer) error {
	replicas := app.Replicas
	if replicas <= 1 {
		_, err := d.DeployAppService(ctx, app, "", logWriter)
		return err
	}

	primaryName := utils.NormalizeContainerName(app.ID)
	logDrains, _ := d.store.ListLogDrainsByService(app.ID)

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

		runOpts := ContainerRunOptions{
			Name:            containerName,
			ImageTag:        newImageTag,
			ServiceID:       app.ID,
			Domain:          app.Domain,
			InternalPort:    app.InternalPort,
			RuntimeMode:     app.RuntimeMode,
			Envs:            envSlice,
			Cmd:             strings.Fields(app.StartCommand),
			MemoryLimitMB:   app.MemoryLimit,
			CPURequest:      app.CPULimit,
			HealthCheckPath: app.HealthCheckPath,
			Volumes:         app.Volumes,
			MaintenanceMode: app.MaintenanceMode,
			LogDrains:       logDrains,
		}

		newID, err := d.containerManager.CreateAndStart(ctx, runOpts)
		if err != nil {
			return fmt.Errorf("rolling update failed on replica %d: %w", i+1, err)
		}

		if err := d.verifyHealthCheck(ctx, app, containerName, logWriter); err != nil {
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
		imageTag = fmt.Sprintf("codedock-app-%s:latest", app.ID)
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

	return nil
}

func (d *Deployer) RestartAppService(ctx context.Context, app *models.AppService) error {
	containerName := utils.NormalizeContainerName(app.ID)
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

	app.Status = models.AppServiceStatusRunning
	return d.store.UpdateAppService(app)
}
