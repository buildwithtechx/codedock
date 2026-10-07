package deployments

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"strings"
)

func (s *DeploymentService) dockerDependencyReady(ctx context.Context, source string) error {
	var target *client.Client
	release := noopRelease
	var id, key, owner string
	var err error
	if service, ok := strings.CutPrefix(source, "app-"); ok {
		target, release, err = s.DockerForService(ctx, service)
		if err != nil {
			return err
		}
		app, readErr := s.appRepo.GetByID(ctx, service)
		if readErr != nil {
			release()
			return readErr
		}
		if app.ContainerID == "" {
			release()
			return fmt.Errorf("prerequisite is not running")
		}
		id, key, owner = app.ContainerID, "codedock.service_id", service
	} else if database, ok := strings.CutPrefix(source, "db-"); ok {
		if s.Databases == nil {
			return fmt.Errorf("database prerequisite unavailable")
		}
		record, readErr := s.Databases.GetByID(ctx, database)
		if readErr != nil {
			return readErr
		}
		if record.ContainerID == "" {
			return fmt.Errorf("prerequisite is not running")
		}
		target, release, err = s.DockerForProject(ctx, record.ProjectID)
		if err != nil {
			return err
		}
		id, key, owner = record.ContainerID, "codedock.database_id", database
	} else {
		return fmt.Errorf("unsupported prerequisite")
	}
	defer release()
	inspected, err := target.ContainerInspect(ctx, id)
	if err != nil {
		return err
	}
	if inspected.Config == nil || inspected.Config.Labels[key] != owner {
		return fmt.Errorf("prerequisite ownership changed")
	}
	if inspected.State == nil || !inspected.State.Running {
		return fmt.Errorf("prerequisite is not running")
	}
	if inspected.State.Health != nil && inspected.State.Health.Status != "healthy" {
		return fmt.Errorf("prerequisite is not healthy")
	}
	return nil
}
