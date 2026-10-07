package deployments

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"log/slog"

	"codedock.run/codedock/internal/models"
)

type dockerRuntimeStore interface {
	Get(context.Context, string) (*models.ServiceRuntime, error)
}

func noopRelease() {}

func (s *DeploymentService) DockerForProject(ctx context.Context, projectID string) (*client.Client, func(), error) {
	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, noopRelease, err
	}
	if project.ServerID == "" {
		if s.LocalDocker == nil {
			return nil, noopRelease, fmt.Errorf("local Docker unavailable")
		}
		return s.LocalDocker, noopRelease, nil
	}
	if s.Servers == nil || s.sshManager == nil {
		return nil, noopRelease, fmt.Errorf("SSH Docker unavailable")
	}
	server, err := s.Servers.GetByID(ctx, project.ServerID)
	if err != nil {
		return nil, noopRelease, err
	}
	if server.IsLocal || server.OrganizationID != project.OrganizationID {
		return nil, noopRelease, fmt.Errorf("remote target organization changed")
	}
	target, release, err := s.sshManager.GetDockerClient(ctx, server)
	if err != nil {
		return nil, noopRelease, err
	}
	return target, func() {
		if err := target.Close(); err != nil {
			slog.Warn("close Docker observation transport", "error", err)
		}
		release()
	}, nil
}
func (s *DeploymentService) DockerForService(ctx context.Context, id string) (*client.Client, func(), error) {
	if s.Runtime != nil {
		handles, err := s.Runtime.Handles(ctx, id)
		if err != nil {
			return nil, noopRelease, err
		}
		if handles {
			return nil, noopRelease, fmt.Errorf("use the selected native or cluster runtime controls")
		}
	}
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return nil, noopRelease, err
	}
	if app.ContainerID == "" {
		return nil, noopRelease, fmt.Errorf("service is not deployed")
	}
	target, release, err := s.DockerForProject(ctx, app.ProjectID)
	if err != nil {
		return nil, noopRelease, err
	}
	inspected, err := target.ContainerInspect(ctx, app.ContainerID)
	if err != nil {
		release()
		return nil, noopRelease, err
	}
	if inspected.Config == nil || inspected.Config.Labels["codedock.service_id"] != app.ID {
		release()
		return nil, noopRelease, fmt.Errorf("container ownership changed")
	}
	return target, release, nil
}
func (s *DeploymentService) StreamServiceLogs(ctx context.Context, id string, output io.Writer) error {
	target, release, err := s.DockerForService(ctx, id)
	if err != nil {
		return err
	}
	defer release()
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	inspected, err := target.ContainerInspect(ctx, app.ContainerID)
	if err != nil {
		return err
	}
	if inspected.Config == nil || inspected.Config.Labels["codedock.service_id"] != app.ID {
		return fmt.Errorf("log container ownership changed")
	}
	reader, err := target.ContainerLogs(ctx, app.ContainerID, container.LogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Tail: "100", Timestamps: true})
	if err != nil {
		return err
	}
	defer reader.Close()
	if inspected.Config.Tty {
		_, err = io.Copy(output, reader)
	} else {
		_, err = stdcopy.StdCopy(output, output, reader)
	}
	return err
}
func (s *DeploymentService) ScalingEligible(ctx context.Context, id string) (bool, error) {
	app, err := s.appRepo.GetByID(ctx, id)
	if err != nil {
		return false, err
	}
	project, err := s.projectRepo.Get(ctx, app.ProjectID)
	if err != nil {
		return false, err
	}
	if project.ServerID != "" {
		return false, nil
	}
	if s.Runtime != nil {
		handles, err := s.Runtime.Handles(ctx, id)
		if err != nil {
			return false, err
		}
		if handles {
			return false, nil
		}
	}
	if store, ok := s.Runtime.(dockerRuntimeStore); ok {
		runtime, err := store.Get(ctx, id)
		if err != nil {
			return false, err
		}
		if runtime.Target.Kind != "docker" || runtime.Journal != "" {
			return false, nil
		}
	} else if s.RuntimeKinds != nil {
		runtime, err := s.RuntimeKinds.Get(ctx, id)
		if err != nil {
			return false, err
		}
		if runtime.Target.Kind != "docker" || runtime.Journal != "" {
			return false, nil
		}
	}
	return true, nil
}
