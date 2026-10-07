package http

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/handlers/system"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/operations"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

func configureManaged(server *Server, db *sql.DB, vault *utils.Vault, servers repositories.ServerRepository, users *repositories.UserRepo, orgs repositories.OrganizationRepository, sshManager *ssh.SSHManager, operationService *operations.Service) error {
	repository := repositories.NewManagedRepository(db, vault)
	service := systemservices.NewManagedService(repository, servers, orgs, users, sshManager, operationService, nil)
	if err := service.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover managed servers: %w", err)
	}
	server.managedHandler = system.NewManagedHandler(service)
	return nil
}

func (s *Server) registerManagedRoutes(authGroup *echo.Group) {
	member := s.authGuard.RequireOrgRole(models.MemberPermissionMember)
	admin := s.authGuard.RequireOrgRole(models.MemberPermissionAdmin)

	authGroup.GET("/managed/catalog", s.managedHandler.Catalog)
	authGroup.GET("/organizations/:id/managed/credentials", s.managedHandler.ListCredentials, member)
	authGroup.POST("/organizations/:id/managed/credentials", s.managedHandler.CreateCredential, admin)
	authGroup.DELETE("/organizations/:id/managed/credentials/:credentialId", s.managedHandler.DeleteCredential, admin)
	authGroup.GET("/organizations/:id/managed/quota", s.managedHandler.GetQuota, member)
	authGroup.PUT("/organizations/:id/managed/quota", s.managedHandler.SetQuota, admin)
	authGroup.POST("/organizations/:id/managed/servers/review", s.managedHandler.ReviewProvision, member, s.authGuard.RequireScope("server:write"))
	authGroup.POST("/organizations/:id/managed/servers/:serverId/refresh", s.managedHandler.RefreshServer, member, s.authGuard.RequireScope("server:read"))
	authGroup.POST("/organizations/:id/managed/servers/:serverId/resize/review", s.managedHandler.ReviewResize, admin)
	authGroup.POST("/organizations/:id/managed/servers/:serverId/delete/review", s.managedHandler.ReviewDelete, admin)
	authGroup.POST("/managed-operations/:operationId/apply", s.managedHandler.ApplyOperation, s.authGuard.RequireScope("server:write"))
}
