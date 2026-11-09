package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/engine/kubernetes"
	appsmod "codedock/internal/modules/apps"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/canvas"
	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/modules/environments"
	"codedock/internal/modules/monitoring"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/servers"
	"codedock/internal/platform/observability"
	backup "codedock/internal/platform/storage"
	"codedock/internal/platform/vault"
)

type ClusterDeployer interface {
	observability.CanvasDockerTargets
	SetRuntime(r any)
}

func configureClusters(
	server *Server,
	db *sql.DB,
	vlt *vault.Vault,
	projects projectsmod.ProjectRepository,
	serversRepo servers.Repository,
	operations *monitoring.OperationService,
	gate modclusters.Gate,
	apps appsmod.AppServiceRepository,
	deployments ClusterDeployer,
	canvas *canvas.CanvasService,
	backups *backup.BackupManager,
) error {
	repository := modclusters.NewClusterRepo(db, vlt)
	if err := repository.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover cluster operations: %w", err)
	}
	runner := kubernetes.NewClusterRunner(serversRepo)
	service := modclusters.NewClusterService(repository, projects, serversRepo, runner, operations, gate, &http.Client{Timeout: time.Minute})
	if err := service.RecoverUpgrades(context.Background()); err != nil {
		return fmt.Errorf("recover cluster upgrades: %w", err)
	}
	server.clusterHandler = modclusters.NewHandler(&clusterServiceAdapter{
		service:    service,
		operations: operations,
	})
	dataService := modclusters.NewClusterDataService(
		modclusters.NewClusterDataRepo(db, vlt),
		repository,
		environments.NewEnvironmentRepo(db),
		backupmodule.NewRepository(db, vlt),
		runner,
		operations,
		gate,
		&http.Client{Timeout: time.Minute},
	)
	dataService.SetSnapshots(modclusters.NewRedisSnapshotRepo(db))
	if err := dataService.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover cluster databases: %w", err)
	}
	server.clusterDataHandler = modclusters.NewDataHandler(&clusterDataServiceAdapter{
		service:    dataService,
		operations: operations,
		apps:       apps,
		vars:       appsmod.NewServiceVarRepo(db),
	})
	return nil
}

func (s *Server) registerClusterRoutes(group *echo.Group) {
	modclusters.RegisterRoutes(group, modclusters.RoutesConfig{
		Handler:     s.clusterHandler,
		DataHandler: s.clusterDataHandler,
		Middleware: modclusters.RouteMiddleware{
			Admin:          s.authGuard.RequireRole("admin"),
			DatabaseManage: s.authGuard.RequireScope("database:manage"),
		},
	})
}
