package deployments

import (
	"codedock.run/codedock/internal/engine/deploy"
	"codedock.run/codedock/internal/engine/networking"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"context"
	"fmt"
	"github.com/containerd/errdefs"
	"log/slog"
	"strings"
)

func (s *DeploymentService) dockerTarget(ctx context.Context, app *models.AppService) (*deploy.Deployer, func(), error) {
	if s.projectRepo == nil {
		return s.deployer, func() {}, nil
	}
	project, err := s.projectRepo.Get(ctx, app.ProjectID)
	if err != nil {
		return nil, nil, err
	}
	if project.ServerID == "" {
		return s.deployer, func() {}, nil
	}
	if s.Servers == nil || s.sshManager == nil || s.HostGate == nil {
		return nil, nil, fmt.Errorf("SSH Docker runtime unavailable")
	}
	server, err := s.Servers.GetByID(ctx, project.ServerID)
	if err != nil {
		return nil, nil, err
	}
	if server.IsLocal || server.OrganizationID != project.OrganizationID {
		return nil, nil, fmt.Errorf("remote target organization changed")
	}
	unlock, err := s.HostGate.AcquireVolume("server:" + server.ID)
	if err != nil {
		return nil, nil, err
	}
	client, release, err := s.sshManager.GetDockerClient(ctx, server)
	if err != nil {
		unlock()
		return nil, nil, err
	}
	cleanup := func() {
		if err := client.Close(); err != nil {
			slog.Warn("close remote Docker transport", "error", err)
		}
		release()
		unlock()
	}
	proxy, inspectErr := client.ContainerInspect(ctx, networking.TraefikContainerName)
	if inspectErr != nil && !errdefs.IsNotFound(inspectErr) {
		cleanup()
		return nil, nil, inspectErr
	}
	if inspectErr == nil && (proxy.Config == nil || !strings.HasPrefix(proxy.Config.Image, "traefik:")) {
		cleanup()
		return nil, nil, fmt.Errorf("remote proxy name belongs to another runtime")
	}
	if errdefs.IsNotFound(inspectErr) {
		if err := networking.NewTraefikManager(client, "").EnsureTraefikRunning(ctx); err != nil {
			cleanup()
			return nil, nil, err
		}
	}
	runtime, err := s.deployer.ForDockerHost(ctx, client, server.ID)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	return runtime, cleanup, nil
}
func (s *DeploymentService) SetRemoteTargets(servers repositories.ServerRepository, gate deploy.VolumeOperations) {
	s.Servers = servers
	s.HostGate = gate
}
