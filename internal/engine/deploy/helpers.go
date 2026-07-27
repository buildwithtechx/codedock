package deploy

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/engine/build"
	"codedock.run/codedock/internal/models"
)

func ApplyCustomDNS(hostCfg *container.HostConfig, customDNS string) {
	if strings.TrimSpace(customDNS) == "" {
		return
	}
	parts := strings.Split(customDNS, ",")
	var dnsList []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			dnsList = append(dnsList, p)
		}
	}
	if len(dnsList) > 0 {
		hostCfg.DNS = dnsList
	}
}

func (d *Deployer) getEnvironmentVariables(app *models.AppService, logWriter io.Writer) (map[string]string, error) {
	envVarsMap, err := d.store.GetEnvVars(app.ProjectID)
	if err != nil && logWriter != nil {
		fmt.Fprintf(logWriter, "⚠️ [Deployer] Warning: could not load shared project environment variables: %v\n", err)
	}
	if envVarsMap == nil {
		envVarsMap = make(map[string]string)
	}

	if d.EnvProvider != nil {
		providerVars, err := d.EnvProvider(app.ProjectID)
		if err == nil {
			for k, v := range providerVars {
				if _, exists := envVarsMap[k]; !exists {
					envVarsMap[k] = v
				}
			}
		}
	}

	appVars, err := d.store.ListServiceVariables(app.ID)
	if err == nil {
		for _, v := range appVars {
			envVarsMap[v.Key] = v.Value
		}
	}

	var registry map[string]map[string]string
	if d.EnvInterpolator != nil {
		registry, _ = d.EnvInterpolator(app.ProjectID)
	}
	envVarsMap = build.InterpolateEnvVars(envVarsMap, registry)

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
		if err != nil || !inspect.State.Running {
			continue
		}

		checkPath := healthCheckPath
		if checkPath == "" {
			checkPath = "/"
		}
		containerIP := ""
		if inspect.NetworkSettings != nil {
			for _, net := range inspect.NetworkSettings.Networks {
				if net.IPAddress != "" {
					containerIP = net.IPAddress
					break
				}
			}
		}
		if containerIP == "" {
			containerIP = containerName
		}
		targetURL := fmt.Sprintf("http://%s:%d%s", containerIP, internalPort, checkPath)
		client := http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get(targetURL)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode >= 200 && resp.StatusCode < 400 {
				return true
			}
		}
	}
	return false
}
