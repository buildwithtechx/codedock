package http

import (
	"codedock.run/codedock/internal/engine/kubernetes"
	"codedock.run/codedock/internal/engine/observability"
	"codedock.run/codedock/internal/handlers/system"
	"codedock.run/codedock/internal/repositories"
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
	runtime := runtimes.NewService(repositories.NewRuntimeRepo(db, vault), apps, repository, servers, kubernetes.NewWorkloadRuntime(runner), server.deployer, operations, gate)
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
	group.GET("/projects/:id/clusters", s.clusterHandler.List, admin)
	group.POST("/projects/:id/clusters/review", s.clusterHandler.Review, admin)
	group.POST("/cluster-operations/:operationId/apply", s.clusterHandler.Apply, admin)
	group.GET("/cluster-operations/:operationId", s.clusterHandler.GetOperation, admin)
	group.POST("/cluster-operations/:operationId/cancel", s.clusterHandler.Cancel, admin)
}
