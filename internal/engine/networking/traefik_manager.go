package networking

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/containerd/errdefs"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"codedock.run/codedock/internal/config"
)

const (
	TraefikContainerName = "codedock-traefik"
	CodedockNetworkName  = "codedock-network"
)

type TraefikManager struct {
	dockerClient *client.Client
	tlsEmail     string
}

func NewTraefikManager(cli *client.Client, tlsEmail string) *TraefikManager {
	return &TraefikManager{dockerClient: cli, tlsEmail: tlsEmail}
}

func (m *TraefikManager) EnsureTraefikRunning(ctx context.Context) error {
	if err := m.ensureNetwork(ctx); err != nil {
		return fmt.Errorf("failed to ensure network: %w", err)
	}

	_, err := m.dockerClient.ContainerInspect(ctx, TraefikContainerName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			if err := m.createTraefikContainer(ctx); err != nil {
				return fmt.Errorf("failed to create traefik: %w", err)
			}
		} else {
			return err
		}
	}

	if err := m.dockerClient.ContainerStart(ctx, TraefikContainerName, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start traefik: %w", err)
	}

	slog.Info("traefik reverse proxy is running")
	return nil
}

func (m *TraefikManager) ensureNetwork(ctx context.Context) error {
	_, err := m.dockerClient.NetworkInspect(ctx, CodedockNetworkName, network.InspectOptions{})
	if err != nil {
		if errdefs.IsNotFound(err) {
			_, err = m.dockerClient.NetworkCreate(ctx, CodedockNetworkName, network.CreateOptions{
				Driver: "bridge",
			})
			return err
		}
		return err
	}
	return nil
}

func traefikImage() string {
	if img := config.Get().Traefik.Image; img != "" {
		return img
	}
	return "traefik:v3.6"
}

func dockerSocketPath() string {
	if p := config.Get().Docker.SocketPath; p != "" {
		return p
	}
	return "/var/run/docker.sock"
}

func (m *TraefikManager) createTraefikContainer(ctx context.Context) error {
	imageRef := traefikImage()
	out, err := m.dockerClient.ImagePull(ctx, imageRef, image.PullOptions{})
	if err == nil {
		defer out.Close()
		io.Copy(io.Discard, out)
	}

	cmdArgs := m.buildTraefikCmdArgs()
	hostConfig := &container.HostConfig{
		PortBindings: m.buildPortBindings(),
		Mounts:       m.buildTraefikMounts(),
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		ExtraHosts: []string{"host.docker.internal:host-gateway"},
	}

	resp, err := m.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: imageRef,
		Cmd:   cmdArgs,
		ExposedPorts: nat.PortSet{
			"80/tcp":   struct{}{},
			"443/tcp":  struct{}{},
			"443/udp":  struct{}{},
			"8080/tcp": struct{}{},
		},
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers.traefik.entrypoints":               "http",
			"traefik.http.routers.traefik.service":                   "api@internal",
			"traefik.http.services.traefik.loadbalancer.server.port": "8080",
		},
		Healthcheck: &container.HealthConfig{
			Test:     []string{"CMD", "wget", "-qO-", "http://localhost:80/ping"},
			Interval: 4 * time.Second,
			Timeout:  2 * time.Second,
			Retries:  5,
		},
	}, hostConfig, &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			CodedockNetworkName: {},
		},
	}, nil, TraefikContainerName)

	if err != nil {
		return err
	}

	slog.Info("created traefik container", "containerID", resp.ID)
	return nil
}

func (m *TraefikManager) buildTraefikCmdArgs() []string {
	cmdArgs := []string{
		"--providers.docker=true",
		"--providers.docker.exposedbydefault=false",
		"--providers.docker.network=codedock-network",
		"--entrypoints.web.address=:80",
		"--entrypoints.websecure.address=:443",
		"--api.insecure=true",
		"--log.level=INFO",
	}

	if m.tlsEmail != "" {
		cmdArgs = append(cmdArgs,
			"--certificatesresolvers.codedock-resolver.acme.email="+m.tlsEmail,
			"--certificatesresolvers.codedock-resolver.acme.storage=/letsencrypt/acme.json",
			"--certificatesresolvers.codedock-resolver.acme.httpchallenge=true",
			"--certificatesresolvers.codedock-resolver.acme.httpchallenge.entrypoint=web",
		)
	}
	if dockerHost := config.Get().Traefik.DockerHost; dockerHost != "" {
		cmdArgs = append(cmdArgs, "--providers.docker.endpoint="+dockerHost)
	}
	return cmdArgs
}

func (m *TraefikManager) buildTraefikMounts() []mount.Mount {
	if config.Get().Traefik.DockerHost != "" {
		return m.buildTraefikDataMounts()
	}

	sockPath := dockerSocketPath()
	mounts := []mount.Mount{
		{
			Type:     mount.TypeBind,
			Source:   sockPath,
			Target:   "/var/run/docker.sock",
			ReadOnly: true,
		},
	}
	return append(mounts, m.buildTraefikDataMounts()...)
}

func (m *TraefikManager) buildTraefikDataMounts() []mount.Mount {
	mounts := make([]mount.Mount, 0, 1)
	if m.tlsEmail != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: "codedock-traefik-acme",
			Target: "/letsencrypt",
		})
	}
	return mounts
}

func (m *TraefikManager) buildPortBindings() nat.PortMap {
	cfg := config.Get()
	httpPort := fmt.Sprintf("%d", cfg.Traefik.HTTPPort)
	httpsPort := fmt.Sprintf("%d", cfg.Traefik.HTTPSPort)
	apiPort := fmt.Sprintf("%d", cfg.Traefik.APIPort)
	return nat.PortMap{
		"80/tcp":   []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: httpPort}},
		"443/tcp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: httpsPort}},
		"443/udp":  []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: httpsPort}},
		"8080/tcp": []nat.PortBinding{{HostIP: "127.0.0.1", HostPort: apiPort}},
	}
}
