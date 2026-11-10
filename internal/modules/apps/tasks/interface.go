package tasks

import (
	"context"
	"time"
)

// ScheduledTaskRepository defines database persistence for periodic cron tasks.
type ScheduledTaskRepository interface {
	Create(ctx context.Context, j *ScheduledTask) error
	GetByID(ctx context.Context, id string) (*ScheduledTask, error)
	ListAll(ctx context.Context) ([]ScheduledTask, error)
	ListByProject(ctx context.Context, projectID string) ([]ScheduledTask, error)
	ListByService(ctx context.Context, serviceID string) ([]ScheduledTask, error)
	Update(ctx context.Context, j *ScheduledTask) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status ScheduledTaskStatus, lastRunAt *time.Time, output string) error
}

// ScheduledTaskOperations orchestrates scheduled task lifecycles and manual executions.
type ScheduledTaskOperations interface {
	CreateScheduledTask(ctx context.Context, task *ScheduledTask) (*ScheduledTask, error)
	GetScheduledTask(ctx context.Context, id string) (*ScheduledTask, error)
	ListAllScheduledTasks(ctx context.Context) ([]ScheduledTask, error)
	ListScheduledTasksByProject(ctx context.Context, projectID string) ([]ScheduledTask, error)
	ListScheduledTasksByService(ctx context.Context, serviceID string) ([]ScheduledTask, error)
	UpdateScheduledTask(ctx context.Context, task *ScheduledTask) error
	DeleteScheduledTask(ctx context.Context, id string) error
	ExecuteScheduledTask(ctx context.Context, id string) (string, error)
}

// ScheduledTaskServiceFull provides full service operations used by HTTP handlers.
type ScheduledTaskServiceFull interface {
	ListByService(ctx context.Context, serviceID string) ([]ScheduledTask, error)
	ListByProject(ctx context.Context, projectID string) ([]ScheduledTask, error)
	CreateScheduledTask(ctx context.Context, task *ScheduledTask) (*ScheduledTask, error)
	GetScheduledTask(ctx context.Context, id string) (*ScheduledTask, error)
	UpdateScheduledTask(ctx context.Context, task *ScheduledTask) error
	DeleteScheduledTask(ctx context.Context, id string) error
	ExecuteScheduledTask(ctx context.Context, id string) (string, error)
}
