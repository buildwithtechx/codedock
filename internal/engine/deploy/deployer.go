package deploy

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

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
		var pullOpts image.PullOptions
		if app.RegistryID != nil && *app.RegistryID != "" {
			reg, err := d.store.GetRegistry(*app.RegistryID)
			if err == nil && reg != nil {
				authConfig := registry.AuthConfig{
					Username:      reg.Username,
					Password:      reg.PasswordToken,
					ServerAddress: reg.RegistryURL,
				}
				encodedJSON, err := json.Marshal(authConfig)
				if err == nil {
					pullOpts.RegistryAuth = base64.URLEncoding.EncodeToString(encodedJSON)
				}
			}
		}
		if d.containerManager != nil && d.containerManager.dockerClient != nil {
			out, err := d.containerManager.dockerClient.ImagePull(ctx, imageTag, pullOpts)
			if err == nil {
				io.Copy(io.Discard, out)
				out.Close()
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
	app.InternalPort = internalPort

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

	if replicas > 1 {
		_ = d.containerManager.StopAndRemove(ctx, primaryContainerName)
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
