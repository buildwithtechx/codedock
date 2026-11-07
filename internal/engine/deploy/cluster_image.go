package deploy

import (
	"codedock/internal/engine/build"
	"codedock/internal/models"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/google/uuid"
	"io"
	"net/url"
	"strings"
)

func (d *Deployer) PrepareClusterWorkload(ctx context.Context, app *models.AppService, target models.RuntimeTarget, source string, logs io.Writer) (*models.KubernetesWorkload, error) {
	variables, err := d.getEnvironmentVariables(ctx, app, logs)
	if err != nil {
		return nil, err
	}
	workload := &models.KubernetesWorkload{App: *app, Target: target, Image: app.ImageRef, Variables: variables}
	registryID := target.RegistryID
	if registryID == "" && app.RegistryID != nil {
		registryID = *app.RegistryID
	}
	if registryID != "" {
		credentials, err := d.store.GetRegistry(registryID)
		if err != nil {
			return nil, err
		}
		if credentials == nil || credentials.ProjectID != app.ProjectID {
			return nil, fmt.Errorf("registry must belong to this project")
		}
		workload.Registry = credentials
	}
	if workload.App.Replicas <= 0 {
		workload.App.Replicas = 1
	}
	if workload.App.InternalPort <= 0 {
		workload.App.InternalPort = defaultAppPort()
	}
	if workload.Image != "" {
		return workload, nil
	}
	if workload.Registry == nil || target.ImageRepository == "" {
		return nil, fmt.Errorf("Git cluster deployments require a project registry and image repository")
	}
	host := strings.TrimSuffix(workload.Registry.RegistryURL, "/")
	if strings.Contains(host, "://") {
		parsed, err := url.Parse(host)
		if err != nil {
			return nil, err
		}
		host = parsed.Host
	}
	if !strings.HasPrefix(target.ImageRepository, host+"/") || strings.ContainsAny(target.ImageRepository, " \t\r\n@") {
		return nil, fmt.Errorf("image repository must be under the selected registry host")
	}
	if strings.Contains(strings.TrimPrefix(target.ImageRepository, host+"/"), ":") {
		return nil, fmt.Errorf("image repository cannot contain a tag")
	}
	if d.containerManager == nil || d.containerManager.dockerClient == nil {
		return nil, fmt.Errorf("Docker builder is unavailable")
	}
	source, err = build.ResolveSourceRoot(source, app.RootDirectory)
	if err != nil {
		return nil, err
	}
	if err := d.prepareServerlessCode(app, source, logs); err != nil {
		return nil, err
	}
	buildEnv := map[string]string{}
	for key, value := range variables {
		buildEnv[key] = value
	}
	buildEnv["BUILD_COMMAND"] = app.BuildCommand
	local, err := d.builder.Build(ctx, build.BuildOptions{ProjectID: app.ProjectID, ServiceID: app.ID, SourceDir: source, DockerfilePath: app.DockerfilePath, LogWriter: logs, AppConfig: app, EnvVars: buildEnv})
	if err != nil {
		return nil, fmt.Errorf("build cluster image: %w", err)
	}
	workload.Image = target.ImageRepository + ":" + uuid.NewString()
	docker := d.containerManager.dockerClient
	if err := docker.ImageTag(ctx, local, workload.Image); err != nil {
		return nil, err
	}
	auth, err := json.Marshal(registry.AuthConfig{Username: workload.Registry.Username, Password: workload.Registry.PasswordToken, ServerAddress: host})
	if err != nil {
		return nil, err
	}
	stream, err := docker.ImagePush(ctx, workload.Image, image.PushOptions{RegistryAuth: base64.URLEncoding.EncodeToString(auth)})
	if err != nil {
		return nil, fmt.Errorf("publish cluster image: %w", err)
	}
	defer stream.Close()
	decoder := json.NewDecoder(stream)
	for {
		var event struct {
			Error  string `json:"error"`
			Status string `json:"status"`
		}
		err := decoder.Decode(&event)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if event.Error != "" {
			return nil, fmt.Errorf("publish cluster image: %s", event.Error)
		}
		if logs != nil && event.Status != "" {
			if _, err := fmt.Fprintln(logs, event.Status); err != nil {
				return nil, err
			}
		}
	}
	if workload.App.StaticOutput != "" {
		workload.App.InternalPort = 80
	}
	return workload, nil
}
func (d *Deployer) DockerWorkloadActive(ctx context.Context, service string) (bool, error) {
	if d.containerManager == nil || d.containerManager.dockerClient == nil {
		return false, fmt.Errorf("Docker observation unavailable")
	}
	containers, err := d.containerManager.dockerClient.ContainerList(ctx, container.ListOptions{Filters: filters.NewArgs(filters.Arg("label", "codedock.service_id="+service))})
	if err != nil {
		return false, fmt.Errorf("observe Docker workload: %w", err)
	}
	return len(containers) > 0, nil
}
