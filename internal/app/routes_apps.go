package app

import (
	"log/slog"

	"github.com/labstack/echo/v4"

	"codedock/internal/adapters"
	cdkhttp "codedock/internal/http"
	"codedock/internal/modules/actions"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/apps/tasks"
	"codedock/internal/modules/audit"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	"codedock/internal/modules/dns/ingress"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/notifications/chatapps"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/security"
	"codedock/internal/modules/support"
)

func (a *App) registerAppRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	appAccess := &adapters.AppAccessChecker{Ps: s.projectService}
	envLookup := &adapters.EnvLookupAdapter{Repo: s.environmentRepo}
	appsHandler := appsmod.NewHandler(s.appService, appAccess, envLookup)
	appsHandler.SetDeployer(&adapters.AppLifecycleDeployerAdapter{Ds: s.deploymentService, AppService: s.appService})
	appVarHandler := appsmod.NewVariableHandler(s.appService, appAccess)
	appVolumeHandler := appsmod.NewVolumeHandler(s.appService, appAccess)
	appWebhookHandler := appsmod.NewWebhookHandler(s.appService)
	autoscalingHandler := appsmod.NewAutoscalingHandler(s.autoscalingService)

	appsmod.RegisterRoutes(authGroup, appsmod.RoutesConfig{
		Handler:            appsHandler,
		VariableHandler:    appVarHandler,
		VolumeHandler:      appVolumeHandler,
		WebhookHandler:     appWebhookHandler,
		AutoscalingHandler: autoscalingHandler,
	})
}

func (a *App) registerDeploymentRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	serviceAuth := a.RequireServiceRole("")

	deploymentHandler := deploymentsmod.NewDeploymentHandler(
		s.deploymentService,
		&adapters.AppLookupAdapter{AppService: s.appService},
		&adapters.ProjectLookupAdapter{ProjectService: s.projectService},
		&adapters.AIExplainerAdapter{Svc: s.aiAnalysisService},
		&adapters.PRPreviewLookupAdapter{Pr: s.prPreviewService},
	)
	lifecycleHandler := deploymentsmod.NewLifecycleHandler(
		s.deploymentService,
		&adapters.AppLookupAdapter{AppService: s.appService},
		&adapters.AuditLoggerAdapter{Audit: s.auditService},
	)
	registryRepo := deploymentsmod.NewRegistryRepository(a.db)
	registryHandler := deploymentsmod.NewRegistryHandler(deploymentsmod.NewRegistryService(registryRepo))
	archiveHandler := deploymentsmod.NewArchiveHandler(
		deploymentsmod.NewArchiveService(s.deploymentService),
		&adapters.ProjectLookupAdapter{ProjectService: s.projectService},
	)
	terminalHandler := deploymentsmod.NewTerminalHandler(
		a.dockerClient,
		s.authMod.TokenService,
		&adapters.AppLookupAdapter{AppService: s.appService},
		&adapters.ProjectLookupAdapter{ProjectService: s.projectService},
		&adapters.UserStatusAdapter{Us: s.userService},
	)
	terminalHandler.Targets = s.deploymentService

	deploymentsmod.RegisterRoutes(
		apiGroup,
		authGroup,
		deploymentHandler,
		lifecycleHandler,
		archiveHandler,
		registryHandler,
		terminalHandler,
		deploymentsmod.RouteMiddleware{
			Admin:              a.authGuard.RequireRole("admin"),
			ServiceAuth:        serviceAuth,
			ServiceAdmin:       a.RequireServiceRole(organizations.MemberPermissionAdmin),
			ProjectAuth:        a.RequireProjectRole(""),
			ProjectAdmin:       a.RequireProjectRole(organizations.MemberPermissionAdmin),
			RequireProjectRole: func(r string) echo.MiddlewareFunc { return a.RequireProjectRole(organizations.MemberPermission(r)) },
			RequireScope:       a.authGuard.RequireScope,
		},
	)

	gitmod.RegisterRoutes(
		apiGroup,
		authGroup,
		gitmod.NewProvidersHandler(s.gitService),
		gitmod.NewInspectionHandler(s.gitService),
		gitmod.NewAppsHandler(s.gitAppsService),
		gitmod.NewWebhookHandler(s.webhookDispatcher, s.webhookDeliveryService, s.prPreviewService, &adapters.GitServiceLookupAdapter{AppRepo: s.appRepo}),
		gitmod.RouteMiddleware{
			Admin:              a.authGuard.RequireRole("admin"),
			RequireProjectRole: func(r string) echo.MiddlewareFunc { return a.RequireProjectRole(organizations.MemberPermission(r)) },
		},
	)
}

func (a *App) registerMiscRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	apiGroup.GET("/meta", a.getMeta)
	authGroup.GET("/docs/openapi.json", cdkhttp.ServeOpenAPISpec)

	securityVerifier, err := security.NewVerifier("")
	if err != nil {
		slog.Error("failed to initialize security verifier", "error", err)
	}
	securityService := security.NewSecurityService(a.db, security.NewFeedManager(a.db, securityVerifier, ""), nil, nil, nil, nil)
	security.RegisterRoutes(authGroup, security.NewHandler(securityService))

	chatAppService := chatapps.NewService(chatapps.NewRepository(a.db), a.vault)
	chatapps.RegisterRoutes(apiGroup, authGroup, chatapps.NewHandler(chatAppService, chatapps.NewDispatcher(chatAppService, nil)))

	dns.RegisterRoutes(authGroup, dns.NewHandler(s.dnsService), a.authGuard.RequireRole("admin"))

	ingressAdapter := &adapters.IngressServiceAdapter{
		Repo:  ingress.NewIngressRepository(a.db),
		CfSvc: ingress.NewCloudflareIngressService(nil),
	}
	ingress.RegisterRoutes(authGroup, ingress.NewIngressHandler(ingressAdapter), a.authGuard.RequireRole("admin"))

	actions.RegisterRoutes(authGroup, a.authGuard.RequireRole("admin"), actions.NewActionsHandler())
	support.RegisterRoutes(authGroup, support.NewSupportHandler())

	scheduledTaskHandler := tasks.NewScheduledTaskHandler(
		&adapters.ScheduledTaskServiceAdapter{Svc: s.scheduledTaskService},
		&adapters.ScheduledTaskAccessAdapter{AppService: s.appService, ProjectService: s.projectService},
	)
	tasks.RegisterRoutes(authGroup, scheduledTaskHandler)

	auditLogHandler := audit.NewAuditLogHandler(s.auditService)
	audit.RegisterRoutes(authGroup, auditLogHandler, a.authGuard.RequireRole("admin"))

	a.registerMcpAndAiRoutes(authGroup, s)
}
