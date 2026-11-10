package clusters

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/clusters/kubernetes"
	"codedock/internal/modules/environments"
	"codedock/internal/modules/operations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/servers"
	"codedock/internal/platform/vault"
)

type Services struct {
	Service     *ClusterService
	DataService *ClusterDataService
}

func InitClusterServices(
	db *sql.DB,
	vlt *vault.Vault,
	projects projectsmod.ProjectRepository,
	serversRepo servers.Repository,
	operations *operations.OperationService,
	gate Gate,
) (*Services, error) {
	repository := NewClusterRepo(db, vlt)
	if err := repository.Recover(context.Background()); err != nil {
		return nil, fmt.Errorf("recover cluster operations: %w", err)
	}

	runner := kubernetes.NewClusterRunner(serversRepo)
	service := NewClusterService(repository, projects, serversRepo, runner, operations, gate, &http.Client{Timeout: time.Minute})
	if err := service.RecoverUpgrades(context.Background()); err != nil {
		return nil, fmt.Errorf("recover cluster upgrades: %w", err)
	}

	dataService := NewClusterDataService(
		NewClusterDataRepo(db, vlt),
		repository,
		environments.NewEnvironmentRepo(db),
		backupmodule.NewRepository(db, vlt),
		runner,
		operations,
		gate,
		&http.Client{Timeout: time.Minute},
	)
	dataService.SetSnapshots(NewRedisSnapshotRepo(db))
	if err := dataService.Recover(context.Background()); err != nil {
		return nil, fmt.Errorf("recover cluster databases: %w", err)
	}

	return &Services{
		Service:     service,
		DataService: dataService,
	}, nil
}
