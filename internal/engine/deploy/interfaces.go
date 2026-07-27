package deploy

import (
	"codedock.run/codedock/internal/models"
)

type DeployerStore interface {
	ContainerManagerStore
	ListAppServicesByProject(projectID string) ([]*models.AppService, error)
	GetEnvVars(projectID string) (map[string]string, error)
	ListServiceVariables(serviceID string) ([]*models.Variable, error)
	GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error)
	UpdateAppService(app *models.AppService) error
	ListLogDrainsByService(serviceID string) ([]*models.LogDrain, error)
	GetRegistry(id string) (*models.Registry, error)
}

type DatabaseDeployerStore interface {
	GetServerSettings() (*models.ServerSettings, error)
	UpdateDatabaseStatus(id string, status models.DatabaseStatus, containerID string) error
	GetDatabase(id string) (*models.Database, error)
}

type ContainerManagerStore interface {
	GetServerSettings() (*models.ServerSettings, error)
}
