package engine

import (
	"time"

	"codedock.run/codedock/internal/engine/backup"
	"codedock.run/codedock/internal/engine/deploy"
	"codedock.run/codedock/internal/models"
)

type DeployerStore = deploy.DeployerStore

type DatabaseDeployerStore = deploy.DatabaseDeployerStore

type ContainerManagerStore = deploy.ContainerManagerStore

type CronManagerStore interface {
	ListScheduledTasks() ([]models.ScheduledTask, error)
	GetScheduledTask(id string) (*models.ScheduledTask, error)
	GetProject(id string) (*models.ProjectConfig, error)
	GetAppService(id string) (*models.AppService, error)
	UpdateScheduledTaskStatusAndOutput(id string, status models.ScheduledTaskStatus, lastRunAt *time.Time, output string) error
}

type BackupManagerStore = backup.Store
