package app

import (
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/adapters"
	"codedock/internal/config"
	"codedock/internal/http/middleware"
	"codedock/internal/modules/audit/operations"
	backupmodule "codedock/internal/modules/backups"
	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/modules/databases"
	mailmodule "codedock/internal/modules/mail"
	modmigrations "codedock/internal/modules/migrations"
	modmonitoring "codedock/internal/modules/monitoring"
	"codedock/internal/modules/organizations"
	serversmod "codedock/internal/modules/servers"
	"codedock/internal/modules/servers/managed"
	"codedock/internal/modules/servers/networks"
	"codedock/internal/modules/servers/takeover"
)

func (a *App) registerServerRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	serversUserLookup := &adapters.ServerUserLookupAdapter{Users: s.userRepo}
	serversService := serversmod.NewService(s.serverRepo, serversUserLookup, s.sshManager)
	serverAudit := &adapters.ServerAuditLoggerAdapter{Audit: &adapters.AuditDomainAdapter{Audit: s.auditService}}
	serversHandler := serversmod.NewHandler(serversService, serverAudit)

	takeoverRepo := takeover.NewTakeoverRepository(a.db, a.vault)
	takeoverCreator := &adapters.ProjectCreatorAdapter{Projects: s.projectRepo, Apps: s.appRepo}
	takeoverService := takeover.NewTakeoverService(takeoverRepo, takeoverCreator)
	serversTakeoverHandler := takeover.NewTakeoverHandler(takeoverService, takeoverRepo)

	serversNetworksHandler := networks.NewServerNetworksHandler()
	serversWSHandler := serversmod.NewServerMetricsWSHandler(&adapters.ServerTokenValidatorAdapter{Tokens: s.authMod.TokenService}, serversService, serversUserLookup)
	managedHandler := managed.NewManagedHandler(nil)

	serversmod.RegisterRoutes(apiGroup, authGroup, serversmod.RoutesConfig{
		Handler:         serversHandler,
		ManagedHandler:  managedHandler,
		TakeoverHandler: serversTakeoverHandler,
		NetworksHandler: serversNetworksHandler,
		WSHandler:       serversWSHandler,
		Middleware: serversmod.RouteMiddleware{
			ServerRead:  a.authGuard.RequireScope("server:read"),
			ServerWrite: a.authGuard.RequireScope("server:write"),
			Admin:       a.authGuard.RequireRole("admin"),
			OrgMember:   a.authGuard.RequireOrgRole(organizations.MemberPermissionMember),
			OrgAdmin:    a.authGuard.RequireOrgRole(organizations.MemberPermissionAdmin),
		},
	})
}

func (a *App) registerClusterRoutes(group *echo.Group, s *appServices) {
	clusterHandler := modclusters.NewHandler(&adapters.ClusterServiceAdapter{
		Service:    s.clusterService,
		Operations: s.operationService,
	})
	clusterDataHandler := modclusters.NewDataHandler(&adapters.ClusterDataServiceAdapter{
		Service: s.clusterDataService,
	})

	modclusters.RegisterRoutes(group, modclusters.RoutesConfig{
		Handler:     clusterHandler,
		DataHandler: clusterDataHandler,
		Middleware: modclusters.RouteMiddleware{
			Admin:          a.authGuard.RequireRole("admin"),
			DatabaseManage: a.authGuard.RequireScope("database:manage"),
		},
	})
}

func (a *App) registerDatabaseRoutes(authGroup *echo.Group, s *appServices) {
	databases.RegisterRoutes(authGroup, databases.NewHandler(s.databaseService), a.authGuard.RequireScope)
}

func (a *App) registerMailRoutes(authGroup *echo.Group, s *appServices) {
	mailHandler := mailmodule.NewHandler(mailmodule.NewService(mailmodule.NewRepository(a.db, a.vault)), nil, nil)
	mailmodule.RegisterRoutes(authGroup, mailHandler)
}

func (a *App) registerMigrationRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	migrationService := modmigrations.NewBundleService(a.db, s.dbRepo, a.dataDir, config.Get().Database.URL)
	lifecycleHandler := modmigrations.NewLifecycleHandler(nil)
	bundleHandler := modmigrations.NewBundleHandler(
		&adapters.MigrationBundleAdapter{Service: migrationService},
		&adapters.MigrationUserCounterAdapter{Users: s.userRepo},
		middleware.NewRateLimiter(5, time.Minute),
	)

	modmigrations.RegisterRoutes(apiGroup, authGroup, modmigrations.RoutesConfig{
		LifecycleHandler: lifecycleHandler,
		BundleHandler:    bundleHandler,
		Middleware: modmigrations.RouteMiddleware{
			OrgMember:   a.authGuard.RequireOrgRole("member"),
			OrgAdmin:    a.authGuard.RequireOrgRole("admin"),
			ServerRead:  a.authGuard.RequireScope("server:read"),
			ServerWrite: a.authGuard.RequireScope("server:write"),
			Admin:       a.authGuard.RequireRole("admin"),
		},
	})
}

func (a *App) registerBackupRoutes(authGroup *echo.Group, s *appServices) {
	operationHandler := operations.NewOperationHandler(&adapters.OperationsServiceAdapter{Svc: s.operationService})
	operations.RegisterRoutes(authGroup, operationHandler, a.authGuard.RequireScope("backup:read"), a.authGuard.RequireScope("backup:write"))

	backupHandler := backupmodule.NewHandler(s.backupService)
	backupHandler.SetOperations(s.operationService)
	backupmodule.RegisterRoutes(authGroup, backupHandler, backupmodule.RouteMiddleware{
		Admin:          a.authGuard.RequireRole("admin"),
		BackupWrite:    a.authGuard.RequireScope("backup:write"),
		DatabaseManage: a.authGuard.RequireScope("database:manage"),
	})
}

func (a *App) registerAnalyticsRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	wsHandler := modmonitoring.NewServiceLogsWSHandler(
		s.authMod.TokenService,
		&adapters.WSAccessCheckerAdapter{AppService: s.appService, ProjectService: s.projectService},
		&adapters.WSUserValidatorAdapter{UserRepo: s.userRepo},
	)
	wsHandler.Streams = s.deploymentService

	analyticsHandler := modmonitoring.NewAnalyticsHandler(&adapters.AnalyticsServiceAdapter{Svc: s.analyticsService})
	attentionHandler := modmonitoring.NewAttentionHandler(&adapters.AttentionServiceAdapter{Svc: s.attentionService})
	metricsHandler := modmonitoring.NewMetricsHandler(&adapters.MetricsServiceAdapter{Svc: s.metricsService})
	logHandler := modmonitoring.NewLogHandler(&adapters.LogServiceAdapter{Svc: s.logService})
	exampleHandler := modmonitoring.NewExampleHandler(&adapters.ExampleServiceAdapter{Svc: s.exampleService})

	serviceAuth := a.RequireServiceRole("")
	modmonitoring.RegisterRoutes(authGroup, apiGroup, modmonitoring.RoutesConfig{
		Analytics:     analyticsHandler,
		Attention:     attentionHandler,
		Metrics:       metricsHandler,
		Logs:          logHandler,
		ServiceLogsWS: wsHandler,
		Example:       exampleHandler,
		Middleware: modmonitoring.RouteMiddleware{
			OrgMember:   a.authGuard.RequireOrgRole("member"),
			OrgAdmin:    a.authGuard.RequireOrgRole("admin"),
			ReadScope:   a.authGuard.RequireScope("server:read"),
			WriteScope:  a.authGuard.RequireScope("server:write"),
			ServiceAuth: serviceAuth,
		},
	})
}
