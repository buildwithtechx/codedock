package engine

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/engine/build"
	"codedock.run/codedock/internal/models"
)

func (d *Deployer) getEnvironmentVariables(app *models.AppService, logWriter io.Writer) (map[string]string, error) {
	envVarsMap, err := d.store.GetEnvVars(app.ProjectID)
	if err != nil && logWriter != nil {
		fmt.Fprintf(logWriter, "⚠️ [Deployer] Warning: could not load shared project environment variables: %v\n", err)
	}
	if envVarsMap == nil {
		envVarsMap = make(map[string]string)
	}

	serviceVars, _ := d.store.ListServiceVariables(app.ID)
	for _, sv := range serviceVars {
		envVarsMap[sv.Key] = sv.Value
	}

	if d.EnvProvider != nil {
		if linkedEnvs, err := d.EnvProvider(app.ProjectID); err == nil {
			for k, v := range linkedEnvs {
				if _, exists := envVarsMap[k]; !exists {
					envVarsMap[k] = v
				}
			}
			if logWriter != nil && len(linkedEnvs) > 0 {
				fmt.Fprintf(logWriter, "🔗 [Deployer] Automatically linked %d service connection strings (DATABASE_URL, REDIS_URL, etc.)\n", len(linkedEnvs))
			}
		}
	}

	if d.EnvInterpolator != nil {
		if registry, err := d.EnvInterpolator(app.ProjectID); err == nil && len(registry) > 0 {
			envVarsMap = build.InterpolateEnvVars(envVarsMap, registry)
			if logWriter != nil {
				fmt.Fprintf(logWriter, "🔀 [Deployer] Interpolated dynamic variable references (${service.VAR_KEY} syntax).\n")
			}
		}
	}

	return envVarsMap, nil
}

func defaultAppPort() int {
	return config.Get().Defaults.AppPort
}

func defaultMemoryMB() int {
	return int(config.Get().Defaults.MemoryMB)
}

func defaultCPURequest() float64 {
	return config.Get().Defaults.CPU
}

func (d *Deployer) verifyHealthCheck(ctx context.Context, app *models.AppService, containerName string, logWriter io.Writer) error {
	if app.RuntimeMode == models.RuntimeModeWorker {
		if logWriter != nil {
			fmt.Fprintf(logWriter, "✅ [Deployer] Worker mode detected. Skipping HTTP health check.\n")
		}
		return nil
	}
	healthy := d.waitForHealthyContainer(ctx, containerName, app.HealthCheckPath, app.InternalPort)
	if !healthy {
		_ = d.containerManager.StopAndRemove(ctx, containerName)
		if logWriter != nil {
			fmt.Fprintf(logWriter, "❌ [Deployer] Health check failed. Rolling back to previous version.\n")
		}
		return fmt.Errorf("health check failed, deployment aborted")
	}
	return nil
}

func (d *Deployer) waitForHealthyContainer(ctx context.Context, containerName string, healthCheckPath string, internalPort int) bool {
	maxRetries := 30
	if timeout := config.Get().Limits.DeploymentTimeout; timeout > 0 {
		maxRetries = timeout / 2
	}
	for i := 0; i < maxRetries; i++ {
		time.Sleep(2 * time.Second)
		inspect, err := d.containerManager.Inspect(ctx, containerName)
		if err == nil {
			if !inspect.State.Running {
				if inspect.State.Status == "exited" {
					break
				}
				continue
			}
			if healthCheckPath != "" && inspect.State.Health != nil {
				if inspect.State.Health.Status == "healthy" {
					return true
				}
				if inspect.State.Health.Status == "unhealthy" {
					return false
				}
			} else if healthCheckPath != "" {
				ip := ""
				for _, net := range inspect.NetworkSettings.Networks {
					ip = net.IPAddress
					break
				}
				if ip != "" {
					resp, err := http.Get(fmt.Sprintf("http://%s:%d%s", ip, internalPort, healthCheckPath))
					if err == nil {
						resp.Body.Close()
						if resp.StatusCode >= 200 && resp.StatusCode < 400 {
							return true
						}
					}
				}
			} else {
				return true
			}
		}
	}
	return false
}
