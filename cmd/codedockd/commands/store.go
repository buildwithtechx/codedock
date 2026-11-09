package commands

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"

	cdkdb "codedock/db"
	"codedock/internal/config"
	"codedock/internal/modules/apps"
	"codedock/internal/modules/deployments"
	"codedock/internal/modules/projects"
	"codedock/internal/modules/settings"
	systemdb "codedock/internal/platform/database"
	"codedock/internal/platform/vault"
)

type DBDeployerStore struct {
	DB    *sql.DB
	Vault *vault.Vault
}

func NewDBDeployerStore(db *sql.DB, vlt *vault.Vault) *DBDeployerStore {
	return &DBDeployerStore{DB: db, Vault: vlt}
}

func (a *DBDeployerStore) GetServerSettings() (*settings.ServerSettings, error) {
	return settings.NewSettingsRepo(a.DB, a.Vault).GetServerSettings(context.Background())
}

func (a *DBDeployerStore) ListAppServicesByProject(projectID string) ([]*apps.AppService, error) {
	return apps.NewAppServiceRepo(a.DB).ListByProject(context.Background(), projectID)
}

func (a *DBDeployerStore) GetEnvVars(projectID string) (map[string]string, error) {
	return projects.NewEnvRepo(a.DB, a.Vault).GetVars(context.Background(), projectID)
}

func (a *DBDeployerStore) ListServiceVariables(serviceID string) ([]*apps.Variable, error) {
	return apps.NewServiceVarRepo(a.DB).ListByService(context.Background(), serviceID)
}

func (a *DBDeployerStore) ListLogDrainsByService(serviceID string) ([]*apps.LogDrain, error) {
	return apps.NewAppServiceRepo(a.DB).ListLogDrainsByService(context.Background(), serviceID)
}

func (a *DBDeployerStore) GetServerlessFunctionCode(serviceID string) (*apps.ServerlessFunctionCode, error) {
	return apps.NewServerlessRepository(a.DB).GetCodeByServiceID(context.Background(), serviceID)
}

func (a *DBDeployerStore) UpdateAppService(app *apps.AppService) error {
	return apps.NewAppServiceRepo(a.DB).Update(context.Background(), app)
}

func (a *DBDeployerStore) GetRegistry(id string) (*deployments.Registry, error) {
	return deployments.NewRegistryRepository(a.DB).Get(context.Background(), id)
}

func InitDataDir() (string, *sql.DB, *vault.Vault) {
	dataDir := config.Get().Server.DataDir
	if abs, err := filepath.Abs(dataDir); err == nil {
		dataDir = abs
	}
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		slog.Error("failed to create data directory", "err", err)
		os.Exit(1)
	}
	if err := os.Chmod(dataDir, 0o700); err != nil {
		slog.Error("failed to enforce 0700 permissions on data directory", "err", err)
		os.Exit(1)
	}
	vlt, err := vault.NewVault(dataDir)
	if err != nil {
		slog.Error("failed to initialize secrets vault", "err", err)
		os.Exit(1)
	}
	databaseURL := config.Get().Database.URL
	if databaseURL == "" {
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			slog.Error("docker is required for the embedded database", "err", err, "hint", "start docker or set CODEDOCK_DATABASE_URL")
			os.Exit(1)
		}
		supervisor := systemdb.NewSupervisor(dockerClient, systemdb.ContainerName, dataDir, config.Get().Postgres.Image, config.Get().Postgres.Port)
		supervisor.SetConnectHost(config.Get().Postgres.Host)
		supervisor.SetNetwork(config.Get().Postgres.Network)
		databaseURL, err = supervisor.EnsureRunning(context.Background())
		if err != nil {
			slog.Error("failed to provision embedded database", "err", err)
			os.Exit(1)
		}
	}
	sqlDB, err := systemdb.OpenDatabase(databaseURL)
	if err != nil {
		slog.Error("failed to open database", "err", err)
		os.Exit(1)
	}
	if err := cdkdb.RunMigrations(sqlDB); err != nil {
		slog.Error("failed to run database migrations", "err", err)
		os.Exit(1)
	}
	if err := config.PrepareSelfHosted(context.Background(), config.Get(), settings.NewSelfHostedRepo(sqlDB)); err != nil {
		slog.Error("failed to prepare self-hosted configuration", "err", err)
		os.Exit(1)
	}
	return dataDir, sqlDB, vlt
}
