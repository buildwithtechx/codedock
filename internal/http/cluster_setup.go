package http

import (
	"codedock.run/codedock/internal/engine/kubernetes"
	"codedock.run/codedock/internal/handlers/system"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/clusters"
	"codedock.run/codedock/internal/services/operations"
	"codedock.run/codedock/internal/utils"
	"context"
	"database/sql"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func configureClusters(server *Server, db *sql.DB, vault *utils.Vault, projects *repositories.ProjectRepo, servers repositories.ServerRepository, operations *operations.Service, gate clusters.Gate) error {
	repository := repositories.NewClusterRepo(db, vault)
	if err := repository.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover cluster operations: %w", err)
	}
	runner := kubernetes.NewClusterRunner(servers)
	service := clusters.NewService(repository, projects, servers, runner, operations, gate, &http.Client{Timeout: time.Minute})
	server.clusterHandler = system.NewClusterHandler(service, operations)
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
