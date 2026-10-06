package deployments

import (
	"codedock.run/codedock/internal/models"
	"context"
	"fmt"
	"io"
	"strings"
)

type TargetRuntime interface {
	Handles(context.Context, string) (bool, error)
	Deploy(context.Context, *models.AppService, string, io.Writer) (string, error)
	Observe(context.Context, string) (*models.WorkloadObservation, error)
	Lifecycle(context.Context, string, string, int) error
}

func (s *DeploymentService) deployTarget(ctx context.Context, app *models.AppService, source string, logs io.Writer) (string, error) {
	if s.Runtime != nil {
		handles, err := s.Runtime.Handles(ctx, app.ID)
		if err != nil {
			return "", err
		}
		if handles {
			return s.Runtime.Deploy(ctx, app, source, logs)
		}
	}
	runtime, release, err := s.dockerTarget(ctx, app)
	if err != nil {
		return "", err
	}
	defer release()
	return runtime.DeployAppService(ctx, app, source, logs)
}
func (s *DeploymentService) StopAppService(ctx context.Context, app *models.AppService) error {
	if s.Runtime != nil {
		handles, err := s.Runtime.Handles(ctx, app.ID)
		if err != nil {
			return err
		}
		if handles {
			return s.Runtime.Lifecycle(ctx, app.ID, "stop", 0)
		}
	}
	runtime, release, err := s.dockerTarget(ctx, app)
	if err != nil {
		return err
	}
	defer release()
	return runtime.StopAppService(ctx, app)
}
func (s *DeploymentService) RestartAppService(ctx context.Context, app *models.AppService) error {
	if s.Runtime != nil {
		handles, err := s.Runtime.Handles(ctx, app.ID)
		if err != nil {
			return err
		}
		if handles {
			return s.Runtime.Lifecycle(ctx, app.ID, "restart", 0)
		}
	}
	runtime, release, err := s.dockerTarget(ctx, app)
	if err != nil {
		return err
	}
	defer release()
	return runtime.RestartAppService(ctx, app)
}

func (s *DeploymentService) dependencyReady(ctx context.Context, source string) error {
	if id, ok := strings.CutPrefix(source, "app-"); ok && s.Runtime != nil {
		handles, err := s.Runtime.Handles(ctx, id)
		if err != nil {
			return err
		}
		if handles {
			observed, err := s.Runtime.Observe(ctx, id)
			if err != nil {
				return err
			}
			if observed.Status != "READY" {
				return fmt.Errorf("cluster prerequisite %s is not ready", source)
			}
			return nil
		}
	}
	return s.deployer.DependencyReady(ctx, source)
}

func (s *DeploymentService) RemoveAppService(ctx context.Context, app *models.AppService) error {
	if s.Runtime != nil {
		handles, err := s.Runtime.Handles(ctx, app.ID)
		if err != nil {
			return err
		}
		if handles {
			remover, ok := s.Runtime.(interface {
				Remove(context.Context, string) error
			})
			if !ok {
				return fmt.Errorf("runtime removal unavailable")
			}
			return remover.Remove(ctx, app.ID)
		}
	}
	runtime, release, err := s.dockerTarget(ctx, app)
	if err != nil {
		return err
	}
	defer release()
	return runtime.StopAppService(ctx, app)
}
