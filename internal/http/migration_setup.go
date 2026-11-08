package http

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	modmigrations "codedock/internal/modules/migrations"
	"codedock/internal/platform/ssh"
	"codedock/internal/platform/vault"
	"codedock/internal/repositories"
	"codedock/internal/services/migrations"
	projectservices "codedock/internal/services/projects"
)

func configureMigration(
	server *Server,
	db *sql.DB,
	vlt *vault.Vault,
	projects *repositories.ProjectRepo,
	apps *projectservices.AppService,
	appRepo *repositories.AppServiceRepo,
	envs *repositories.EnvironmentRepo,
	servers repositories.ServerRepository,
	users *repositories.UserRepo,
	orgs repositories.OrganizationRepository,
	git migrations.GitSyncer,
	composer *projectservices.ComposeParserService,
	deployer migrations.AppDeployer,
	sshManager *ssh.SSHManager,
	runtimes *repositories.RuntimeRepo,
) error {
	repository := repositories.NewMigrationRepository(db, vlt)
	service := migrations.NewService(repository, projects, apps, appRepo, repositories.NewServiceVolumeRepo(db), envs, servers, users, orgs, git, composer, deployer, sshManager, runtimes, nil, nil)
	if err := service.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover migrations: %w", err)
	}
	server.migrationLifecycleHandler = modmigrations.NewLifecycleHandler(&migrationLifecycleAdapter{service: service})
	server.migrationService = service
	return nil
}

func (s *Server) registerMigrationRoutes(authGroup *echo.Group) {
	modmigrations.RegisterRoutes(nil, authGroup, modmigrations.RoutesConfig{
		LifecycleHandler: s.migrationLifecycleHandler,
		Middleware: modmigrations.RouteMiddleware{
			OrgMember:   s.authGuard.RequireOrgRole(models.MemberPermissionMember),
			OrgAdmin:    s.authGuard.RequireOrgRole(models.MemberPermissionAdmin),
			ServerWrite: s.authGuard.RequireScope("server:write"),
			ServerRead:  s.authGuard.RequireScope("server:read"),
		},
	})
}
