package engine

import (
	"time"

	"codedock/internal/engine/deploy"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/monitoring"
	projectsmod "codedock/internal/modules/projects"
	backup "codedock/internal/platform/storage"
)

type DeployerStore = deploy.DeployerStore

type DatabaseDeployerStore = deploy.DatabaseDeployerStore

type ContainerManagerStore = deploy.ContainerManagerStore

type CronManagerStore interface {
	ListScheduledTasks() ([]monitoring.ScheduledTask, error)
	GetScheduledTask(id string) (*monitoring.ScheduledTask, error)
	GetProject(id string) (*projectsmod.ProjectConfig, error)
	GetAppService(id string) (*appsmod.AppService, error)
	UpdateScheduledTaskStatusAndOutput(id string, status monitoring.ScheduledTaskStatus, lastRunAt *time.Time, output string) error
}

type BackupManagerStore = backup.Store
