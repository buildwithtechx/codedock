package deploy

import (
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	settingsmodule "codedock/internal/modules/settings"
)

type DeployerStore interface {
	ContainerManagerStore
	ListAppServicesByProject(projectID string) ([]*appsmod.AppService, error)
	GetEnvVars(projectID string) (map[string]string, error)
	ListServiceVariables(serviceID string) ([]*appsmod.Variable, error)
	GetServerlessFunctionCode(serviceID string) (*appsmod.ServerlessFunctionCode, error)
	UpdateAppService(app *appsmod.AppService) error
	ListLogDrainsByService(serviceID string) ([]*appsmod.LogDrain, error)
	GetRegistry(id string) (*deploymentsmod.Registry, error)
}

type DatabaseDeployerStore interface {
	GetServerSettings() (*settingsmodule.ServerSettings, error)
	UpdateDatabaseStatus(id string, status databases.DatabaseStatus, containerID string) error
	GetDatabase(id string) (*databases.Database, error)
}

type ContainerManagerStore interface {
	GetServerSettings() (*settingsmodule.ServerSettings, error)
}
