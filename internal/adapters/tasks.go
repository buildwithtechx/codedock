package adapters

import (
	"context"
	"errors"

	"codedock/internal/modules/apps"
	"codedock/internal/modules/apps/tasks"
	"codedock/internal/modules/projects"
	"codedock/internal/modules/users"
)

type ScheduledTaskServiceAdapter struct {
	Svc *tasks.ScheduledTaskService
}

func (s *ScheduledTaskServiceAdapter) ListByService(ctx context.Context, serviceID string) ([]tasks.ScheduledTask, error) {
	if s.Svc == nil {
		return []tasks.ScheduledTask{}, nil
	}
	return s.Svc.ListScheduledTasksByService(ctx, serviceID)
}

func (s *ScheduledTaskServiceAdapter) ListByProject(ctx context.Context, projectID string) ([]tasks.ScheduledTask, error) {
	if s.Svc == nil {
		return []tasks.ScheduledTask{}, nil
	}
	return s.Svc.ListScheduledTasksByProject(ctx, projectID)
}

func (s *ScheduledTaskServiceAdapter) CreateScheduledTask(ctx context.Context, task *tasks.ScheduledTask) (*tasks.ScheduledTask, error) {
	if s.Svc == nil {
		return nil, errors.New("task service unavailable")
	}
	return s.Svc.CreateScheduledTask(ctx, task)
}

func (s *ScheduledTaskServiceAdapter) GetScheduledTask(ctx context.Context, id string) (*tasks.ScheduledTask, error) {
	if s.Svc == nil {
		return nil, errors.New("task service unavailable")
	}
	return s.Svc.GetScheduledTask(ctx, id)
}

func (s *ScheduledTaskServiceAdapter) UpdateScheduledTask(ctx context.Context, task *tasks.ScheduledTask) error {
	if s.Svc == nil {
		return errors.New("task service unavailable")
	}
	return s.Svc.UpdateScheduledTask(ctx, task)
}

func (s *ScheduledTaskServiceAdapter) DeleteScheduledTask(ctx context.Context, id string) error {
	if s.Svc == nil {
		return errors.New("task service unavailable")
	}
	return s.Svc.DeleteScheduledTask(ctx, id)
}

func (s *ScheduledTaskServiceAdapter) ExecuteScheduledTask(ctx context.Context, id string) (string, error) {
	if s.Svc == nil {
		return "", errors.New("task service unavailable")
	}
	return s.Svc.ExecuteScheduledTask(ctx, id)
}

type ScheduledTaskAccessAdapter struct {
	AppService     *apps.Service
	ProjectService *projects.Service
}

func (a *ScheduledTaskAccessAdapter) HasServiceAccess(ctx context.Context, userID, role, serviceID string) bool {
	if a.AppService == nil || a.ProjectService == nil || serviceID == "" {
		return false
	}
	svc, err := a.AppService.GetAppService(ctx, serviceID)
	if err != nil || svc == nil || svc.ProjectID == "" {
		return false
	}
	return a.ProjectService.HasPermission(ctx, svc.ProjectID, userID, users.UserRole(role), "")
}

func (a *ScheduledTaskAccessAdapter) HasProjectAccess(ctx context.Context, userID, role, projectID string) bool {
	if a.ProjectService == nil || projectID == "" {
		return false
	}
	return a.ProjectService.HasPermission(ctx, projectID, userID, users.UserRole(role), "")
}
