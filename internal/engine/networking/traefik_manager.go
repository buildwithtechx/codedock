package networking

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
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
	DynamicDir           = "/etc/codedock/traefik-dynamic"
)

type TraefikManager struct {
	dockerClient *client.Client
	tlsEmail     string
	logDir       string
}

func (m *TraefikManager) SetLogDir(dir string) {
	m.logDir = dir
}

func (m *TraefikManager) AccessLogPath() string {
	if m.logDir == "" {
		return ""
	}
	return m.logDir + "/access.json"
}

func NewTraefikManager(cli *client.Client, tlsEmail string) *TraefikManager {
	return &TraefikManager{dockerClient: cli, tlsEmail: tlsEmail}
}

func (m *TraefikManager) EnsureTraefikRunning(ctx context.Context) error {
	if err := m.ensureNetwork(ctx); err != nil {
		return fmt.Errorf("failed to ensure network: %w", err)
	}

	existing, err := m.dockerClient.ContainerInspect(ctx, TraefikContainerName)
	if err != nil {
		if errdefs.IsNotFound(err) {
			if err := m.createTraefikContainer(ctx); err != nil {
				return fmt.Errorf("failed to create traefik: %w", err)
			}
		} else {
			return err
		}
	}

	if err == nil && existing.Config != nil && !equalStrings(existing.Config.Cmd, m.buildTraefikCmdArgs()) {
		if existing.State != nil && existing.State.Running {
			if err := m.dockerClient.ContainerStop(ctx, TraefikContainerName, container.StopOptions{}); err != nil {
				return fmt.Errorf("stop proxy to apply configuration: %w", err)
			}
		}
		if err := m.dockerClient.ContainerRemove(ctx, TraefikContainerName, container.RemoveOptions{}); err != nil {
			return fmt.Errorf("replace proxy configuration: %w", err)
		}
		if err := m.createTraefikContainer(ctx); err != nil {
			return fmt.Errorf("apply proxy configuration: %w", err)
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
		"--providers.file.directory=/dynamic",
		"--providers.file.watch=true",
		"--entrypoints.web.address=:80",
		"--entrypoints.websecure.address=:443",
		"--api.insecure=true",
		"--ping=true",
		"--ping.entrypoint=web",
		"--log.level=INFO",
	}

	if m.tlsEmail != "" {
		cmdArgs = append(cmdArgs,
			"--certificatesresolvers.letsencrypt.acme.email="+m.tlsEmail,
			"--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json",
			"--certificatesresolvers.letsencrypt.acme.httpchallenge=true",
			"--certificatesresolvers.letsencrypt.acme.httpchallenge.entrypoint=web",
		)
	}
	if dockerHost := config.Get().Traefik.DockerHost; dockerHost != "" {
		cmdArgs = append(cmdArgs, "--providers.docker.endpoint="+dockerHost)
	}
	if m.logDir != "" {
		cmdArgs = append(cmdArgs, "--accesslog=true", "--accesslog.format=json", "--accesslog.filepath=/var/log/traefik/access.json")
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
		{
			Type:   mount.TypeBind,
			Source: DynamicDir,
			Target: "/dynamic",
		},
	}
	return append(mounts, m.buildTraefikDataMounts()...)
}

func (m *TraefikManager) buildTraefikDataMounts() []mount.Mount {
	mounts := make([]mount.Mount, 0, 2)
	if m.tlsEmail != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeVolume,
			Source: "codedock-traefik-acme",
			Target: "/letsencrypt",
		})
	}
	if m.logDir != "" {
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: m.logDir,
			Target: "/var/log/traefik",
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

func traefikCertificateEmail(args []string) string {
	for _, arg := range args {
		if value, found := strings.CutPrefix(arg, "--certificatesresolvers.letsencrypt.acme.email="); found {
			return value
		}
	}
	return ""
}

func equalStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
