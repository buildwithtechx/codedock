package deploy

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
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
	rolloutDirectory      string
	builder               build.Builder
	containerManager      *ContainerManager
	store                 DeployerStore
	ScopedEnvProvider     func(context.Context, string, string) (map[string]string, error)
	ScopedEnvInterpolator func(context.Context, string, string) (map[string]map[string]string, error)
	EnvProvider           func(projectID string) (map[string]string, error)
	EnvInterpolator       func(projectID string) (map[string]map[string]string, error)
	RouteRuleFetcher      func(ctx context.Context, serviceID, serviceName string) (map[string]string, error)
}

func NewDeployer(dockerClient *client.Client, s DeployerStore) *Deployer {
	return &Deployer{
		builder:          build.NewBuilder(dockerClient),
		containerManager: NewContainerManager(dockerClient, s),
		store:            s,
	}
}

func (d *Deployer) DeployAppService(ctx context.Context, app *models.AppService, sourceDir string, logWriter io.Writer) (string, error) {
	release, err := d.serviceOperation(ctx, app.ID)
	if err != nil {
		return "", err
	}
	defer release()
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

	if d.containerManager.dockerClient == nil {
		return "", fmt.Errorf("Docker runtime is unavailable")
	}
	if app.ImageRef == "" {
		sourceDir, err = build.ResolveSourceRoot(sourceDir, app.RootDirectory)
		if err != nil {
			return "", err
		}
	}
	if err := d.prepareServerlessCode(app, sourceDir, logWriter); err != nil {
		return "", err
	}

	envVarsMap, err := d.getEnvironmentVariables(ctx, app, logWriter)
	if err != nil {
		return "", fmt.Errorf("failed resolving service env vars: %w", err)
	}

	appDomain := app.Domain
	if appDomain == "" {
		appDomain = utils.GenerateAppDomain(app.ID+"-"+app.Name, "", "")
	}

	if strings.Contains(appDomain, "://") {
		parsed, err := url.Parse(appDomain)
		if err != nil || parsed.Hostname() == "" {
			return "", fmt.Errorf("invalid application domain")
		}
		appDomain = parsed.Hostname()
	}

	var envSlice []string
	for k, v := range envVarsMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}

	internalPort := app.InternalPort
	if app.StaticOutput != "" && app.ImageRef == "" {
		internalPort = 80
	}
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
			if err != nil {
				return "", fmt.Errorf("load image registry: %w", err)
			}
			if reg == nil {
				return "", fmt.Errorf("image registry not found")
			}
			{
				authConfig := registry.AuthConfig{
					Username:      reg.Username,
					Password:      reg.PasswordToken,
					ServerAddress: reg.RegistryURL,
				}
				encodedJSON, err := json.Marshal(authConfig)
				if err != nil {
					return "", fmt.Errorf("encode registry authentication: %w", err)
				}
				pullOpts.RegistryAuth = base64.URLEncoding.EncodeToString(encodedJSON)
			}
		}
		if d.containerManager != nil && d.containerManager.dockerClient != nil {
			out, err := d.containerManager.dockerClient.ImagePull(ctx, imageTag, pullOpts)
			if err != nil {
				return "", fmt.Errorf("pull image: %w", err)
			}
			decoder := json.NewDecoder(out)
			for {
				var event struct {
					Error string `json:"error"`
				}
				err := decoder.Decode(&event)
				if err == io.EOF {
					break
				}
				if err != nil {
					out.Close()
					return "", fmt.Errorf("read image pull: %w", err)
				}
				if event.Error != "" {
					out.Close()
					return "", fmt.Errorf("pull image: %s", event.Error)
				}
			}
			if err := out.Close(); err != nil {
				return "", fmt.Errorf("close image pull: %w", err)
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

	logDrains, _ := d.store.ListLogDrainsByService(app.ID)
	app.InternalPort = internalPort

	extraLabels := map[string]string{}
	if d.RouteRuleFetcher != nil {
		if ml, err := d.RouteRuleFetcher(ctx, app.ID, app.Name); err != nil {
			if logWriter != nil {
				fmt.Fprintf(logWriter, "Failed to fetch route rules: %v\n", err)
			}
			return "", fmt.Errorf("failed to fetch route rules: %w", err)
		} else {
			extraLabels = ml
		}
	}

	return d.replaceReplicas(ctx, app, ContainerRunOptions{
		ImageTag: imageTag, ServiceID: app.ID, Domain: appDomain, InternalPort: internalPort,
		RuntimeMode: app.RuntimeMode, Envs: envSlice, Cmd: applicationCommand(app),
		MemoryLimitMB: memoryLimit, CPURequest: cpuRequest, HealthCheckPath: app.HealthCheckPath,
		Volumes: app.Volumes, MaintenanceMode: app.MaintenanceMode, LogDrains: logDrains, ExtraLabels: extraLabels,
	}, replicas, logWriter)
}

func (d *Deployer) StopAppService(ctx context.Context, app *models.AppService) error {
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

func applicationCommand(app *models.AppService) []string {
	if app.StartCommand == "" || app.StaticOutput != "" && app.ImageRef == "" {
		return nil
	}
	return []string{"/bin/sh", "-c", app.StartCommand}
}
