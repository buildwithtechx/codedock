package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/engine/bare"
	"codedock/internal/engine/kubernetes"
	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/platform/observability"
	backup "codedock/internal/platform/storage"
	"codedock/internal/platform/vault"
	"codedock/internal/repositories"
	"codedock/internal/services/clusterdata"
	"codedock/internal/services/clusters"
	"codedock/internal/services/operations"
	projectservices "codedock/internal/services/projects"
	"codedock/internal/services/runtimes"
)

type ClusterDeployer interface {
	observability.CanvasDockerTargets
	SetRuntime(r any)
}

func configureClusters(server *Server, db *sql.DB, vlt *vault.Vault, projects *repositories.ProjectRepo, servers repositories.ServerRepository, operations *operations.Service, gate clusters.Gate, apps repositories.AppServiceRepository, deployments ClusterDeployer, canvas *projectservices.CanvasService, backups *backup.BackupManager) error {
	repository := repositories.NewClusterRepo(db, vlt)
	if err := repository.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover cluster operations: %w", err)
	}
	runner := kubernetes.NewClusterRunner(servers)
	service := clusters.NewService(repository, projects, servers, runner, operations, gate, &http.Client{Timeout: time.Minute})
	if err := service.RecoverUpgrades(context.Background()); err != nil {
		return fmt.Errorf("recover cluster upgrades: %w", err)
	}
	server.clusterHandler = modclusters.NewHandler(&clusterServiceAdapter{
		service:    service,
		operations: operations,
	})
	dataService := clusterdata.NewService(repositories.NewClusterDataRepo(db, vlt), repository, repositories.NewEnvironmentRepo(db), repositories.NewS3DestinationRepo(db, vlt), runner, operations, gate, &http.Client{Timeout: time.Minute})
	dataService.SetSnapshots(repositories.NewRedisSnapshotRepo(db))
	if err := dataService.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover cluster databases: %w", err)
	}
	server.clusterDataHandler = modclusters.NewDataHandler(&clusterDataServiceAdapter{
		service:    dataService,
		operations: operations,
		apps:       apps,
		vars:       repositories.NewServiceVarRepo(db),
	})
	runtime := runtimes.NewService(repositories.NewRuntimeRepo(db, vlt), apps, repository, servers, kubernetes.NewWorkloadRuntime(runner), server.deployer, operations, gate)
	runtime.SetStreamer(runner)
	runtime.SetBare(bare.NewRuntime(runner), projects)
	if err := runtime.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover Kubernetes application rollouts: %w", err)
	}
	reconcileCtx, cancelReconciliation := context.WithCancel(context.Background())
	server.router.Server.RegisterOnShutdown(cancelReconciliation)
	go runtime.RunReconciler(reconcileCtx)
	if deployments != nil {
		deployments.SetRuntime(runtime)
	}
	backups.SetNativeData(runtime.BackupNativeData, runtime.RestoreNativeData)
	canvasRuntime := observability.NewCanvasRuntime(server.dockerClient)
	canvasRuntime.Cluster = runtime
	canvasRuntime.Targets = deployments
	canvas.SetRuntime(canvasRuntime)
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
