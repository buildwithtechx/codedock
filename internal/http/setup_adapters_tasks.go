package http

import (
	"context"
	"errors"

	"codedock/internal/modules/apps"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/projects"
	"codedock/internal/modules/users"
)

type scheduledTaskServiceAdapter struct {
	svc *monitoring.ScheduledTaskService
}

func (s *scheduledTaskServiceAdapter) ListByService(ctx context.Context, serviceID string) ([]monitoring.ScheduledTask, error) {
	if s.svc == nil {
		return []monitoring.ScheduledTask{}, nil
	}
	return s.svc.ListScheduledTasksByService(ctx, serviceID)
}

func (s *scheduledTaskServiceAdapter) ListByProject(ctx context.Context, projectID string) ([]monitoring.ScheduledTask, error) {
	if s.svc == nil {
		return []monitoring.ScheduledTask{}, nil
	}
	return s.svc.ListScheduledTasksByProject(ctx, projectID)
}

func (s *scheduledTaskServiceAdapter) CreateScheduledTask(ctx context.Context, task *monitoring.ScheduledTask) (*monitoring.ScheduledTask, error) {
	if s.svc == nil {
		return nil, errors.New("task service unavailable")
	}
	return s.svc.CreateScheduledTask(ctx, task)
}

func (s *scheduledTaskServiceAdapter) GetScheduledTask(ctx context.Context, id string) (*monitoring.ScheduledTask, error) {
	if s.svc == nil {
		return nil, errors.New("task service unavailable")
	}
	return s.svc.GetScheduledTask(ctx, id)
}

func (s *scheduledTaskServiceAdapter) UpdateScheduledTask(ctx context.Context, task *monitoring.ScheduledTask) error {
	if s.svc == nil {
		return errors.New("task service unavailable")
	}
	return s.svc.UpdateScheduledTask(ctx, task)
}

func (s *scheduledTaskServiceAdapter) DeleteScheduledTask(ctx context.Context, id string) error {
	if s.svc == nil {
		return errors.New("task service unavailable")
	}
	return s.svc.DeleteScheduledTask(ctx, id)
}

func (s *scheduledTaskServiceAdapter) ExecuteScheduledTask(ctx context.Context, id string) (string, error) {
	if s.svc == nil {
		return "", errors.New("task service unavailable")
	}
	return s.svc.ExecuteScheduledTask(ctx, id)
}

type scheduledTaskAccessAdapter struct {
	appService     *apps.Service
	projectService *projects.Service
}

func (a *scheduledTaskAccessAdapter) HasServiceAccess(ctx context.Context, userID, role, serviceID string) bool {
	if a.appService == nil || a.projectService == nil || serviceID == "" {
		return false
	}
	svc, err := a.appService.GetAppService(ctx, serviceID)
	if err != nil || svc == nil || svc.ProjectID == "" {
		return false
	}
	return a.projectService.HasPermission(ctx, svc.ProjectID, userID, users.UserRole(role), "")
}

func (a *scheduledTaskAccessAdapter) HasProjectAccess(ctx context.Context, userID, role, projectID string) bool {
	if a.projectService == nil || projectID == "" {
		return false
	}
	return a.projectService.HasPermission(ctx, projectID, userID, users.UserRole(role), "")
}
