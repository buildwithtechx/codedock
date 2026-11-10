package databases

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	compose "codedock/internal/platform/docker"
	"codedock/internal/utils"
)

type dockerDeployer struct {
	dockerClient *client.Client
	repo         Repository
	volumes      VolumeGate
}

func NewDockerDeployer(dockerClient *client.Client, repo Repository, volumes VolumeGate) Deployer {
	return &dockerDeployer{
		dockerClient: dockerClient,
		repo:         repo,
		volumes:      volumes,
	}
}

func (d *dockerDeployer) acquireVolume(id string) (func(), error) {
	if d.volumes == nil {
		return func() {}, nil
	}
	return d.volumes.AcquireVolume("codedock-db-data-" + id)
}

func (d *dockerDeployer) SpinUp(ctx context.Context, dbConfig *Database) (string, error) {
	if d.dockerClient == nil {
		return "", fmt.Errorf("docker daemon connection is not available")
	}

	release, err := d.acquireVolume(dbConfig.ID)
	if err != nil {
		return "", err
	}
	defer release()
	if utils.IsDryRun() {
		return fmt.Sprintf("codedock-db-%s-dryrun", dbConfig.Name), nil
	}

	if d.repo != nil {
		fresh, err := d.repo.GetByID(ctx, dbConfig.ID)
		if err != nil {
			return "", fmt.Errorf("reload database before startup: %w", err)
		}
		if fresh == nil {
			return "", utils.NewNotFoundError("Database", dbConfig.ID)
		}
		*dbConfig = *fresh
	}

	containerName := utils.NormalizeContainerName(fmt.Sprintf("codedock-db-%s", dbConfig.Name))
	_ = d.dockerClient.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})

	imageName, envVars, cmd, containerMountPath, err := d.buildContainerConfig(dbConfig)
	if err != nil {
		return "", err
	}

	pullResp, err := d.dockerClient.ImagePull(ctx, imageName, image.PullOptions{})
	if err == nil {
		_, _ = io.Copy(io.Discard, pullResp)
		_ = pullResp.Close()
	}

	if err := utils.EnsureCodedockNetwork(ctx, d.dockerClient); err != nil {
		return "", fmt.Errorf("failed to ensure Docker network: %w", err)
	}

	containerCfg, hostCfg, netCfg := d.createContainerSettings(dbConfig, containerName, imageName, envVars, cmd, containerMountPath)

	if d.repo != nil {
		resolvers := strings.TrimSpace(d.repo.GetCustomDNSResolvers(ctx))
		if resolvers != "" {
			for _, part := range strings.FieldsFunc(resolvers, func(r rune) bool { return r == ',' || r == ';' || r == ' ' || r == '\n' }) {
				if p := strings.TrimSpace(part); p != "" {
					hostCfg.DNS = append(hostCfg.DNS, p)
				}
			}
		}
	}

	created, err := d.dockerClient.ContainerCreate(ctx, containerCfg, hostCfg, netCfg, nil, containerName)
	if err != nil {
		return "", fmt.Errorf("failed to create database container: %w", err)
	}
	if err := d.dockerClient.ContainerStart(ctx, created.ID, container.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start database container: %w", err)
	}
	internalDNS := fmt.Sprintf("%s:%d", containerName, dbConfig.Port)
	if d.repo != nil {
		if err := d.repo.UpdateStatus(ctx, dbConfig.ID, DatabaseStatusRunning, created.ID); err != nil {
			slog.Warn("failed to update database status", "error", err)
		}
	}
	dbConfig.ContainerID = created.ID
	dbConfig.Status = DatabaseStatusRunning
	dbConfig.InternalDNS = internalDNS
	return created.ID, nil
}

func (d *dockerDeployer) buildContainerConfig(dbConfig *Database) (string, []string, []string, string, error) {
	var imageName string
	var envVars []string
	var cmd []string
	var containerMountPath string

	tmplMgr, err := compose.NewTemplateManager()
	if err != nil {
		return "", nil, nil, "", fmt.Errorf("failed to initialize template manager: %w", err)
	}

	composeFile, err := tmplMgr.GetTemplate(strings.ToLower(string(dbConfig.Engine)))
	if err != nil {
		return "", nil, nil, "", fmt.Errorf("unsupported database engine %s: %w", dbConfig.Engine, err)
	}

	tmplService, exists := composeFile.Services[strings.ToLower(string(dbConfig.Engine))]
	if !exists {
		for _, s := range composeFile.Services {
			tmplService = s
			break
		}
	}

	imageName = tmplService.Image
	if dbConfig.Version != "" && !strings.Contains(imageName, ":") {
		imageName = imageName + ":" + dbConfig.Version
	} else if dbConfig.Version != "" {
		parts := strings.Split(imageName, ":")
		imageName = parts[0] + ":" + dbConfig.Version
	}

	for _, ev := range tmplService.Environment {
		resolved := strings.ReplaceAll(ev, "${db.password}", dbConfig.Password)
		resolved = strings.ReplaceAll(resolved, "${db.username}", dbConfig.Username)
		resolved = strings.ReplaceAll(resolved, "${db.database_name}", dbConfig.DatabaseName)
		envVars = append(envVars, resolved)
	}

	for i := 0; i < len(tmplService.Command); i++ {
		c := tmplService.Command[i]
		if c == "--requirepass" && dbConfig.Password == "" {
			i++
			continue
		}
		resolved := strings.ReplaceAll(c, "${db.password}", dbConfig.Password)
		resolved = strings.ReplaceAll(resolved, "${db.username}", dbConfig.Username)
		if resolved != "" {
			cmd = append(cmd, resolved)
		}
	}

	if dbConfig.CustomArgs != "" {
		cmd = append(cmd, strings.Fields(dbConfig.CustomArgs)...)
	}

	if dbConfig.LogicalReplication && (strings.Contains(imageName, "postgres") || strings.Contains(imageName, "timescaledb")) {
		cmd = append(cmd, "-c", "wal_level=logical", "-c", "max_replication_slots=10", "-c", "max_wal_senders=10")
	}

	if len(tmplService.Volumes) > 0 {
		parts := strings.Split(tmplService.Volumes[0], ":")
		if len(parts) >= 2 {
			containerMountPath = parts[1]
		}
	} else {
		containerMountPath = "/data"
	}

	return imageName, envVars, cmd, containerMountPath, nil
}

func (d *dockerDeployer) createContainerSettings(dbConfig *Database, containerName, imageName string, envVars, cmd []string, containerMountPath string) (*container.Config, *container.HostConfig, *network.NetworkingConfig) {
	volumeName := fmt.Sprintf("codedock-db-data-%s", dbConfig.ID)

	labels := map[string]string{"codedock.database_id": dbConfig.ID}
	if dbConfig.ExternalDNS != "" {
		labels["traefik.enable"] = "true"
		labels[fmt.Sprintf("traefik.tcp.routers.%s.rule", containerName)] = fmt.Sprintf("HostSNI(`%s`)", dbConfig.ExternalDNS)
		labels[fmt.Sprintf("traefik.tcp.routers.%s.tls", containerName)] = "true"
		labels[fmt.Sprintf("traefik.tcp.routers.%s.tls.certresolver", containerName)] = "letsencrypt"
		labels[fmt.Sprintf("traefik.tcp.routers.%s.entrypoints", containerName)] = "websecure"
		labels[fmt.Sprintf("traefik.tcp.services.%s.loadbalancer.server.port", containerName)] = fmt.Sprintf("%d", dbConfig.Port)
	}

	containerCfg := &container.Config{
		Image:  imageName,
		Env:    envVars,
		Cmd:    cmd,
		Labels: labels,
	}
	memMB := dbConfig.MemoryLimit
	if memMB <= 0 {
		memMB = utils.DefaultDBMemoryMB()
	}
	cpuReq := dbConfig.CPULimit
	if cpuReq <= 0 {
		cpuReq = utils.DefaultDBCPURequest()
	}

	hostCfg := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
		Resources: container.Resources{
			Memory:   utils.MegaBytesToBytes(memMB),
			NanoCPUs: utils.CPURequestToNanoCPUs(cpuReq),
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: volumeName,
				Target: containerMountPath,
			},
		},
	}
	netCfg := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			utils.GetRuntimeNetwork(): {
				Aliases: []string{containerName, dbConfig.Name},
			},
		},
	}

	return containerCfg, hostCfg, netCfg
}
