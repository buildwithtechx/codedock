package http

import (
	"database/sql"

	"github.com/labstack/echo/v4"

	modmigrations "codedock/internal/modules/migrations"
	"codedock/internal/platform/vault"
)

func configureMigration(
	server *Server,
	db *sql.DB,
	vlt *vault.Vault,
) error {
	server.migrationLifecycleHandler = modmigrations.NewLifecycleHandler(nil)
	return nil
}

func (s *Server) registerMigrationRoutes(authGroup *echo.Group) {
	modmigrations.RegisterRoutes(nil, authGroup, modmigrations.RoutesConfig{
		LifecycleHandler: s.migrationLifecycleHandler,
		Middleware: modmigrations.RouteMiddleware{
			OrgMember:   s.authGuard.RequireOrgRole("member"),
			OrgAdmin:    s.authGuard.RequireOrgRole("admin"),
			ServerRead:  s.authGuard.RequireScope("server:read"),
			ServerWrite: s.authGuard.RequireScope("server:write"),
		},
	})
}
