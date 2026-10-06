package http

import (
	"codedock.run/codedock/internal/engine/bare"
	"codedock.run/codedock/internal/engine/kubernetes"
	"codedock.run/codedock/internal/engine/observability"
	"codedock.run/codedock/internal/handlers/system"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/clusterdata"
	"codedock.run/codedock/internal/services/clusters"
	"codedock.run/codedock/internal/services/deployments"
	"codedock.run/codedock/internal/services/operations"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/services/runtimes"
	"codedock.run/codedock/internal/utils"
	"context"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func configureClusters(server *Server, db *sql.DB, vault *utils.Vault, projects *repositories.ProjectRepo, servers repositories.ServerRepository, operations *operations.Service, gate clusters.Gate, apps repositories.AppServiceRepository, deployments *deployments.DeploymentService, canvas *projectservices.CanvasService) error {
	repository := repositories.NewClusterRepo(db, vault)
	if err := repository.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover cluster operations: %w", err)
	}
	runner := kubernetes.NewClusterRunner(servers)
	service := clusters.NewService(repository, projects, servers, runner, operations, gate, &http.Client{Timeout: time.Minute})
	if err := service.RecoverUpgrades(context.Background()); err != nil {
		return fmt.Errorf("recover cluster upgrades: %w", err)
	}
	server.clusterHandler = system.NewClusterHandler(service, operations)
	dataService := clusterdata.NewService(repositories.NewClusterDataRepo(db, vault), repository, repositories.NewEnvironmentRepo(db), repositories.NewS3DestinationRepo(db, vault), runner, operations, gate, &http.Client{Timeout: time.Minute})
	if err := dataService.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover cluster databases: %w", err)
	}
	server.clusterDataHandler = system.NewClusterDataHandler(dataService, operations)
	runtime := runtimes.NewService(repositories.NewRuntimeRepo(db, vault), apps, repository, servers, kubernetes.NewWorkloadRuntime(runner), server.deployer, operations, gate)
	runtime.SetBare(bare.NewRuntime(runner), projects)
	if err := runtime.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover Kubernetes application rollouts: %w", err)
	}
	reconcileCtx, cancelReconciliation := context.WithCancel(context.Background())
	server.router.Server.RegisterOnShutdown(cancelReconciliation)
	go runtime.RunReconciler(reconcileCtx)
	deployments.Runtime = runtime
	server.appServiceHandler.Runtime = runtime
	canvasRuntime := observability.NewCanvasRuntime(server.dockerClient)
	canvasRuntime.Cluster = runtime
	canvas.SetRuntime(canvasRuntime)
	return nil
}

func (s *Server) registerClusterRoutes(group *echo.Group) {
	admin := s.authGuard.RequireRole("admin")
	group.GET("/projects/:id/clusters/:clusterId/databases", s.clusterDataHandler.List, admin)
	group.POST("/projects/:id/clusters/:clusterId/databases/review", s.clusterDataHandler.Review, admin)
	group.GET("/projects/:id/clusters/:clusterId/databases/:databaseId/credentials", s.clusterDataHandler.Credentials, admin, s.authGuard.RequireScope("database:manage"))
	group.GET("/projects/:id/clusters/:clusterId/databases/:databaseId/backups", s.clusterDataHandler.Backups, admin)
	group.POST("/cluster-data-operations/:operationId/apply", s.clusterDataHandler.Apply, admin)
	group.GET("/cluster-data-operations/:operationId", s.clusterDataHandler.Operation, admin)
	group.POST("/cluster-data-operations/:operationId/cancel", s.clusterDataHandler.Cancel, admin)
	group.GET("/projects/:id/clusters", s.clusterHandler.List, admin)
	group.POST("/projects/:id/clusters/review", s.clusterHandler.Review, admin)
	group.POST("/cluster-operations/:operationId/apply", s.clusterHandler.Apply, admin)
	group.GET("/cluster-operations/:operationId", s.clusterHandler.GetOperation, admin)
	group.POST("/cluster-operations/:operationId/cancel", s.clusterHandler.Cancel, admin)
}
