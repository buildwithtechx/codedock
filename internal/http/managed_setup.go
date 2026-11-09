package http

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/servers"
	"codedock/internal/modules/users"
	"codedock/internal/platform/ssh"
	"codedock/internal/platform/vault"
)

func configureManaged(server *Server, db *sql.DB, vlt *vault.Vault, serverRepo servers.Repository, users users.UserRepository, orgs organizations.OrganizationRepository, sshManager *ssh.SSHManager, operationService *monitoring.OperationService) error {
	server.managedHandler = servers.NewManagedHandler(nil)
	return nil
}

func (s *Server) registerManagedRoutes(authGroup *echo.Group) {
	// Managed routes are registered in registerServerRoutes via servers.RegisterRoutes
}
