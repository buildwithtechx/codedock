package http

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/servers"
	"codedock/internal/platform/ssh"
	"codedock/internal/platform/vault"
	"codedock/internal/repositories"
	"codedock/internal/services/operations"
	systemservices "codedock/internal/services/system"
)

func configureManaged(server *Server, db *sql.DB, vlt *vault.Vault, serverRepo repositories.ServerRepository, users *repositories.UserRepo, orgs repositories.OrganizationRepository, sshManager *ssh.SSHManager, operationService *operations.Service) error {
	repository := repositories.NewManagedRepository(db, vlt)
	service := systemservices.NewManagedService(repository, serverRepo, orgs, users, sshManager, operationService, nil)
	if err := service.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover managed servers: %w", err)
	}
	server.managedHandler = servers.NewManagedHandler(&managedServiceAdapter{svc: service})
	return nil
}

func (s *Server) registerManagedRoutes(authGroup *echo.Group) {
	// Managed routes are registered in registerServerRoutes via servers.RegisterRoutes
}
