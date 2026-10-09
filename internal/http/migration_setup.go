package http

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/labstack/echo/v4"

	"codedock/internal/engine/ssh"
	"codedock/internal/handlers/system"
	"codedock/internal/models"
	"codedock/internal/repositories"
	deploymentservices "codedock/internal/services/deployments"
	"codedock/internal/services/migrations"
	projectservices "codedock/internal/services/projects"
	"codedock/internal/utils"
)

func configureMigration(
	server *Server,
	db *sql.DB,
	vault *utils.Vault,
	projects *repositories.ProjectRepo,
	apps *projectservices.AppService,
	appRepo *repositories.AppServiceRepo,
	envs *repositories.EnvironmentRepo,
	servers repositories.ServerRepository,
	users *repositories.UserRepo,
	orgs repositories.OrganizationRepository,
	git *deploymentservices.GitService,
	composer *projectservices.ComposeParserService,
	deployer *deploymentservices.DeploymentService,
	sshManager *ssh.SSHManager,
	runtimes *repositories.RuntimeRepo,
) error {
	repository := repositories.NewMigrationRepository(db, vault)
	service := migrations.NewService(repository, projects, apps, appRepo, repositories.NewServiceVolumeRepo(db), envs, servers, users, orgs, git, composer, deployer, sshManager, runtimes, nil, nil)
	if err := service.Recover(context.Background()); err != nil {
		return fmt.Errorf("recover migrations: %w", err)
	}
	server.migrationLifecycleHandler = system.NewMigrationLifecycleHandler(service)
	server.migrationService = service
	return nil
}

func (s *Server) registerMigrationRoutes(authGroup *echo.Group) {
	member := s.authGuard.RequireOrgRole(models.MemberPermissionMember)
	admin := s.authGuard.RequireOrgRole(models.MemberPermissionAdmin)
	serverWrite := s.authGuard.RequireScope("server:write")
	serverRead := s.authGuard.RequireScope("server:read")

	authGroup.GET("/organizations/:id/migration/sources", s.migrationLifecycleHandler.ListSources, member)
	authGroup.POST("/organizations/:id/migration/sources/test", s.migrationLifecycleHandler.TestSource, member, serverWrite)
	authGroup.POST("/organizations/:id/migration/sources", s.migrationLifecycleHandler.CreateSource, admin)
	authGroup.DELETE("/organizations/:id/migration/sources/:sourceId", s.migrationLifecycleHandler.DeleteSource, admin)
	authGroup.POST("/organizations/:id/migration/scan", s.migrationLifecycleHandler.Scan, member, serverWrite)
	authGroup.GET("/organizations/:id/migration/scan/stream", s.migrationLifecycleHandler.ScanStream, member, serverWrite)
	authGroup.POST("/organizations/:id/migration/reveal-env", s.migrationLifecycleHandler.RevealEnv, admin)
	authGroup.POST("/organizations/:id/migration/adopt", s.migrationLifecycleHandler.Adopt, member, serverWrite)
	authGroup.POST("/organizations/:id/migration/reimport", s.migrationLifecycleHandler.Reimport, member, serverWrite)
	authGroup.POST("/organizations/:id/migration/repo-compose", s.migrationLifecycleHandler.RepoCompose, member)
	authGroup.POST("/organizations/:id/migration/preview", s.migrationLifecycleHandler.Preview, member, serverWrite)
	authGroup.POST("/organizations/:id/migration/migrate", s.migrationLifecycleHandler.StartMigration, member, serverWrite)
	authGroup.POST("/organizations/:id/migration/project", s.migrationLifecycleHandler.StartProjectMove, member, serverWrite)
	authGroup.GET("/organizations/:id/migration/active", s.migrationLifecycleHandler.Active, member, serverRead)
	authGroup.GET("/organizations/:id/migration/runs", s.migrationLifecycleHandler.ListRuns, member, serverRead)
	authGroup.GET("/migrations/:runId", s.migrationLifecycleHandler.GetRun, serverRead)
	authGroup.GET("/migrations/:runId/stream", s.migrationLifecycleHandler.StreamRun, serverRead)
	authGroup.POST("/migrations/:runId/cutover", s.migrationLifecycleHandler.Cutover, serverWrite)
	authGroup.POST("/migrations/:runId/cancel", s.migrationLifecycleHandler.Cancel, serverWrite)
	authGroup.POST("/migrations/:runId/respond", s.migrationLifecycleHandler.Respond, serverWrite)
	authGroup.POST("/migrations/:runId/resume", s.migrationLifecycleHandler.Resume, serverWrite)
	authGroup.POST("/migrations/:runId/cleanup-target", s.migrationLifecycleHandler.CleanupTarget, serverWrite)
	authGroup.DELETE("/migrations/:runId", s.migrationLifecycleHandler.DeleteRun, serverWrite)
}
