package tasks

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	cron "codedock/internal/platform/scheduler"
)

type ScheduledTaskService struct {
	repo    ScheduledTaskRepository
	manager *cron.CronManager
}

func NewScheduledTaskService(r ScheduledTaskRepository, m *cron.CronManager) *ScheduledTaskService {
	return &ScheduledTaskService{
		repo:    r,
		manager: m,
	}
}

func (s *ScheduledTaskService) CreateScheduledTask(ctx context.Context, j *ScheduledTask) (*ScheduledTask, error) {
	if j == nil || j.ServiceID == "" || j.Name == "" {
		return nil, errors.New("valid scheduledTask with serviceId and name required")
	}
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	if j.Schedule == "" {
		j.Schedule = "0 0 * * *"
	}
	if j.Status == "" {
		j.Status = string(ScheduledTaskStatusActive)
	}
	now := time.Now()
	j.CreatedAt = now
	j.UpdatedAt = now
	if err := s.repo.Create(ctx, j); err != nil {
		return nil, err
	}
	if s.manager != nil && j.Status == string(ScheduledTaskStatusActive) {
		_ = s.manager.RegisterScheduledTask(&cron.ScheduledTask{
			ID: j.ID, ServiceID: j.ServiceID, Name: j.Name, Schedule: j.Schedule, Command: j.Command, Status: cron.ScheduledTaskStatus(j.Status),
		})
	}
	return j, nil
}

func (s *ScheduledTaskService) GetScheduledTask(ctx context.Context, id string) (*ScheduledTask, error) {
	if id == "" {
		return nil, errors.New("id required")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ScheduledTaskService) ListAllScheduledTasks(ctx context.Context) ([]ScheduledTask, error) {
	return s.repo.ListAll(ctx)
}

func (s *ScheduledTaskService) ListScheduledTasksByProject(ctx context.Context, projectID string) ([]ScheduledTask, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *ScheduledTaskService) ListByProject(ctx context.Context, projectID string) ([]ScheduledTask, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *ScheduledTaskService) ListScheduledTasksByService(ctx context.Context, serviceID string) ([]ScheduledTask, error) {
	return s.repo.ListByService(ctx, serviceID)
}

func (s *ScheduledTaskService) ListByService(ctx context.Context, serviceID string) ([]ScheduledTask, error) {
	return s.repo.ListByService(ctx, serviceID)
}

func (s *ScheduledTaskService) UpdateScheduledTask(ctx context.Context, j *ScheduledTask) error {
	if j == nil || j.ID == "" {
		return errors.New("scheduledTask with id required")
	}
	j.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, j); err != nil {
		return err
	}
	if s.manager != nil {
		s.manager.UnregisterScheduledTask(j.ID)
		if j.Status == string(ScheduledTaskStatusActive) {
			_ = s.manager.RegisterScheduledTask(&cron.ScheduledTask{
				ID: j.ID, ServiceID: j.ServiceID, Name: j.Name, Schedule: j.Schedule, Command: j.Command, Status: cron.ScheduledTaskStatus(j.Status),
			})
		}
	}
	return nil
}

func (s *ScheduledTaskService) DeleteScheduledTask(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id required")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	if s.manager != nil {
		s.manager.UnregisterScheduledTask(id)
	}
	return nil
}

func (s *ScheduledTaskService) ExecuteScheduledTask(ctx context.Context, id string) (string, error) {
	if s.manager == nil {
		return "", errors.New("cron manager not configured")
	}
	return s.manager.ExecuteScheduledTask(ctx, id)
}
