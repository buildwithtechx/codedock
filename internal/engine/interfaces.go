package engine

import (
	"time"

	"codedock/internal/engine/deploy"
	appsmod "codedock/internal/modules/apps"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/tasks"
	backup "codedock/internal/platform/storage"
)

type DeployerStore = deploy.DeployerStore

type DatabaseDeployerStore = deploy.DatabaseDeployerStore

type ContainerManagerStore = deploy.ContainerManagerStore

type CronManagerStore interface {
	ListScheduledTasks() ([]tasks.ScheduledTask, error)
	GetScheduledTask(id string) (*tasks.ScheduledTask, error)
	GetProject(id string) (*projectsmod.ProjectConfig, error)
	GetAppService(id string) (*appsmod.AppService, error)
	UpdateScheduledTaskStatusAndOutput(id string, status tasks.ScheduledTaskStatus, lastRunAt *time.Time, output string) error
}

type BackupManagerStore = backup.Store
