package commands

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/docker/docker/client"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/engine/systemdb"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/utils"
)

type DBDeployerStore struct {
	DB    *sql.DB
	Vault *utils.Vault
}

func NewDBDeployerStore(db *sql.DB, vlt *utils.Vault) *DBDeployerStore {
	return &DBDeployerStore{DB: db, Vault: vlt}
}

func (a *DBDeployerStore) GetServerSettings() (*models.ServerSettings, error) {
	return repositories.NewSettingsRepo(a.DB, a.Vault).GetServerSettings(context.Background())
}

func (a *DBDeployerStore) ListAppServicesByProject(projectID string) ([]*models.AppService, error) {
	return repositories.NewAppServiceRepo(a.DB).ListByProject(context.Background(), projectID)
}

func (a *DBDeployerStore) GetEnvVars(projectID string) (map[string]string, error) {
	return repositories.NewEnvRepo(a.DB, a.Vault).GetVars(context.Background(), projectID)
}

func (a *DBDeployerStore) ListServiceVariables(serviceID string) ([]*models.Variable, error) {
	svVarRepo := repositories.NewServiceVarRepo(a.DB)
	return svVarRepo.ListByService(context.Background(), serviceID)
}

func (a *DBDeployerStore) ListLogDrainsByService(serviceID string) ([]*models.LogDrain, error) {
	return repositories.NewAppServiceRepo(a.DB).ListLogDrainsByService(context.Background(), serviceID)
}

func (a *DBDeployerStore) GetServerlessFunctionCode(serviceID string) (*models.ServerlessFunctionCode, error) {
	svlsRepo := repositories.NewServerlessRepository(a.DB)
	return svlsRepo.GetCodeByServiceID(context.Background(), serviceID)
}

func (a *DBDeployerStore) UpdateAppService(app *models.AppService) error {
	repo := repositories.NewAppServiceRepo(a.DB)
	return repo.Update(context.Background(), app)
}

func (a *DBDeployerStore) GetRegistry(id string) (*models.Registry, error) {
	repo := repositories.NewRegistryRepository(a.DB)
	return repo.Get(context.Background(), id)
}

func InitDataDir() (string, *sql.DB, *utils.Vault) {
	dataDir := config.Get().Server.DataDir
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		slog.Error("failed to create data directory", "err", err)
		os.Exit(1)
	}
	if err := os.Chmod(dataDir, 0o700); err != nil {
		slog.Error("failed to enforce 0700 permissions on data directory", "err", err)
		os.Exit(1)
	}
	vlt, err := utils.NewVault(dataDir)
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
		databaseURL, err = supervisor.EnsureRunning(context.Background())
		if err != nil {
			slog.Error("failed to provision embedded database", "err", err)
			os.Exit(1)
		}
	}
	db, err := repositories.OpenDatabase(databaseURL)
	if err != nil {
		slog.Error("failed to open database", "err", err)
		os.Exit(1)
	}
	if err := repositories.RunMigrations(db); err != nil {
		slog.Error("failed to run database migrations", "err", err)
		os.Exit(1)
	}
	if err := config.PrepareSelfHosted(context.Background(), config.Get(), repositories.NewSelfHostedRepo(db)); err != nil {
		slog.Error("failed to prepare self-hosted configuration", "err", err)
		os.Exit(1)
	}
	return dataDir, db, vlt
}
