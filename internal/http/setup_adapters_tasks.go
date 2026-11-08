package http

import (
	"context"
	"errors"

	"codedock/internal/models"
	"codedock/internal/modules/monitoring"
	projectservices "codedock/internal/services/projects"
	systemservices "codedock/internal/services/system"
)

type scheduledTaskServiceAdapter struct {
	svc *systemservices.ScheduledTaskService
}

func (s *scheduledTaskServiceAdapter) ListByService(ctx context.Context, serviceID string) ([]monitoring.ScheduledTask, error) {
	if s.svc == nil {
		return []monitoring.ScheduledTask{}, nil
	}
	tasks, err := s.svc.ListScheduledTasksByService(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	res := make([]monitoring.ScheduledTask, len(tasks))
	for i, t := range tasks {
		res[i] = monitoring.ScheduledTask{
			ID:         t.ID,
			ServiceID:  t.ServiceID,
			Name:       t.Name,
			Schedule:   t.Schedule,
			Command:    t.Command,
			Status:     t.Status,
			LastRunAt:  t.LastRunAt,
			LastOutput: t.LastOutput,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}
	}
	return res, nil
}

func (s *scheduledTaskServiceAdapter) ListByProject(ctx context.Context, projectID string) ([]monitoring.ScheduledTask, error) {
	if s.svc == nil {
		return []monitoring.ScheduledTask{}, nil
	}
	tasks, err := s.svc.ListScheduledTasksByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	res := make([]monitoring.ScheduledTask, len(tasks))
	for i, t := range tasks {
		res[i] = monitoring.ScheduledTask{
			ID:         t.ID,
			ServiceID:  t.ServiceID,
			Name:       t.Name,
			Schedule:   t.Schedule,
			Command:    t.Command,
			Status:     t.Status,
			LastRunAt:  t.LastRunAt,
			LastOutput: t.LastOutput,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}
	}
	return res, nil
}

func (s *scheduledTaskServiceAdapter) CreateScheduledTask(ctx context.Context, task *monitoring.ScheduledTask) (*monitoring.ScheduledTask, error) {
	if s.svc == nil {
		return nil, errors.New("task service unavailable")
	}
	m := &models.ScheduledTask{
		ServiceID: task.ServiceID,
		Name:      task.Name,
		Schedule:  task.Schedule,
		Command:   task.Command,
	}
	created, err := s.svc.CreateScheduledTask(ctx, m)
	if err != nil || created == nil {
		return nil, err
	}
	task.ID = created.ID
	task.Status = created.Status
	task.CreatedAt = created.CreatedAt
	task.UpdatedAt = created.UpdatedAt
	return task, nil
}

func (s *scheduledTaskServiceAdapter) GetScheduledTask(ctx context.Context, id string) (*monitoring.ScheduledTask, error) {
	if s.svc == nil {
		return nil, errors.New("task service unavailable")
	}
	t, err := s.svc.GetScheduledTask(ctx, id)
	if err != nil || t == nil {
		return nil, err
	}
	return &monitoring.ScheduledTask{
		ID:         t.ID,
		ServiceID:  t.ServiceID,
		Name:       t.Name,
		Schedule:   t.Schedule,
		Command:    t.Command,
		Status:     t.Status,
		LastRunAt:  t.LastRunAt,
		LastOutput: t.LastOutput,
		CreatedAt:  t.CreatedAt,
		UpdatedAt:  t.UpdatedAt,
	}, nil
}

func (s *scheduledTaskServiceAdapter) UpdateScheduledTask(ctx context.Context, task *monitoring.ScheduledTask) error {
	if s.svc == nil {
		return errors.New("task service unavailable")
	}
	return s.svc.UpdateScheduledTask(ctx, &models.ScheduledTask{
		ID:        task.ID,
		ServiceID: task.ServiceID,
		Name:      task.Name,
		Schedule:  task.Schedule,
		Command:   task.Command,
		Status:    task.Status,
	})
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
	appService     *projectservices.AppService
	projectService *projectservices.ProjectService
}

func (a *scheduledTaskAccessAdapter) HasServiceAccess(ctx context.Context, userID, role, serviceID string) bool {
	if a.appService == nil || a.projectService == nil || serviceID == "" {
		return false
	}
	svc, err := a.appService.GetAppService(ctx, serviceID)
	if err != nil || svc == nil || svc.ProjectID == "" {
		return false
	}
	return a.projectService.HasPermission(ctx, svc.ProjectID, userID, models.UserRole(role), "")
}

func (a *scheduledTaskAccessAdapter) HasProjectAccess(ctx context.Context, userID, role, projectID string) bool {
	if a.projectService == nil || projectID == "" {
		return false
	}
	return a.projectService.HasPermission(ctx, projectID, userID, models.UserRole(role), "")
}
