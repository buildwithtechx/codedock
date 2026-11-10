package app

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/adapters"
	"codedock/internal/config"
	cdkhttp "codedock/internal/http"
	aimod "codedock/internal/modules/ai"
	"codedock/internal/modules/auth"
	"codedock/internal/modules/billing"
	composemod "codedock/internal/modules/compose"
	"codedock/internal/modules/dns"
	"codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/projects"
	"codedock/internal/modules/projects/canvas"
	settingsmodule "codedock/internal/modules/settings"
	systemmod "codedock/internal/modules/system"
	templatesmod "codedock/internal/modules/templates"
	"codedock/internal/modules/users"
	compose "codedock/internal/platform/docker"
)

func (a *App) registerAuthRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	uh := users.NewUserHandler(s.userService, s.mailerService, s.auditService)
	auth.RegisterAuthRoutes(apiGroup, authGroup, auth.AuthRouteDeps{
		AuthHandler:       s.authMod.AuthHandler,
		GetProfile:        uh.GetProfile,
		OAuthHandler:      s.authMod.OAuthHandler,
		DeviceFlowHandler: s.authMod.DeviceFlowHandler,
		PasskeyHandler:    s.authMod.PasskeyHandler,
		BootstrapCSRF:     cdkhttp.BootstrapCSRF,
		AuthRateLimit:     s.authRateLimiter.Middleware,
		OTPRateLimit:      s.otpRateLimiter.Middleware,
		AdminAuth:         a.authGuard.RequireRole("admin"),
	})
}

func (a *App) registerUserRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	uh := users.NewUserHandler(s.userService, s.mailerService, s.auditService)
	users.RegisterRoutes(authGroup, uh, s.authMod.AuthHandler.AdminInviteUser, a.authGuard.RequireRole("admin"), s.otpRateLimiter.Middleware)
	users.RegisterOnboardingRoutes(apiGroup, users.NewOnboardingHandler(s.userService))
}

func (a *App) registerProjectRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	projectHandler := projects.NewHandler(s.projectService, s.projectSettingsService, s.orgService)
	projectAppHandler := projects.NewProjectAppHandler(s.projectAppService, s.projectService)
	projectSettingsHandler := projects.NewProjectSettingsHandler(s.projectSettingsService)
	checksHandler := projects.NewChecksHandler()

	projects.RegisterRoutes(authGroup, projects.RoutesConfig{
		Handler:         projectHandler,
		AppHandler:      projectAppHandler,
		SettingsHandler: projectSettingsHandler,
		ChecksHandler:   checksHandler,
	})

	stackService := composemod.NewComposeStackService(composemod.NewComposeStackRepo(a.db, a.vault), &adapters.ComposeStackRuntimeAdapter{Rt: compose.NewStackRuntime(a.dockerClient)})
	stackService.BeforeDeployment = s.backupService.BeforeProjectDeployment
	composemod.RegisterRoutes(authGroup, composemod.RoutesConfig{
		Handler: composemod.NewComposeStackHandler(&adapters.ComposeStackAdapter{Cs: stackService}),
	})

	tmplMgr, _ := compose.NewTemplateManager()
	oneClickService := templatesmod.NewOneClickService(tmplMgr, stackService, &adapters.TemplateEnvAdapter{Repo: s.environmentRepo})
	templatesmod.RegisterRoutes(authGroup, templatesmod.RoutesConfig{
		Handler: templatesmod.NewOneClickHandler(&adapters.OneClickAdapter{Oc: oneClickService}),
	})

	canvas.RegisterRoutes(authGroup, canvas.NewCanvasHandler(&adapters.CanvasServiceAdapter{Cs: s.canvasService}))

	envService := environments.NewEnvironmentService(s.environmentRepo, environments.NewEnvVarRepo())
	envHandler := environments.NewEnvironmentHandler(envService, &adapters.ProjectLookupAdapter{ProjectService: s.projectService}, environments.NewEnvSuggestionService(&adapters.EnvGitProviderAdapter{Ps: s.gitService}))
	protectionService := environments.NewRouteProtectionService(s.routeProtectionRepo, &adapters.ProjectLookupAdapter{ProjectService: s.projectService})
	routeRuleHandler := environments.NewRouteRuleHandler(environments.NewRouteRuleService(s.routeRuleRepo), &adapters.ServiceCheckerAdapter{AppService: s.appService})
	environments.RegisterRoutes(authGroup, environments.RoutesConfig{
		EnvHandler:        envHandler,
		ProtectionHandler: environments.NewRouteProtectionHandler(protectionService),
		RouteRuleHandler:  routeRuleHandler,
	})

	appAccess := &adapters.AppAccessChecker{Ps: s.projectService}
	dnsProviderService := dns.NewProviderService(&adapters.DNSSettingsAdapter{Repo: s.settingsRepo})
	domainHandler := domains.NewDomainHandler(domains.NewDomainService(s.domainRepo, dnsProviderService), appAccess, &adapters.DomainServiceLookupAdapter{AppRepo: s.appRepo}, &adapters.ServerSettingsAdapter{Repo: s.settingsRepo}, &adapters.AuditDomainAdapter{Audit: s.auditService})
	domains.RegisterRoutes(authGroup, domains.RoutesConfig{
		Handler: domainHandler,
	})
}

func (a *App) registerOrganizationRoutes(authGroup *echo.Group, s *appServices) {
	orgHandler := organizations.NewOrganizationHandler(s.orgService, s.auditService)
	organizations.RegisterRoutes(
		authGroup,
		orgHandler,
		a.authGuard.RequireOrgRole(organizations.MemberPermissionMember),
		a.authGuard.RequireOrgRole(organizations.MemberPermissionAdmin),
		a.authGuard.RequireOrgRole(organizations.MemberPermissionOwner),
	)
}

func (a *App) registerSystemRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	updaterService := systemmod.NewUpdaterService(s.settingsRepo)
	systemHandler := systemmod.NewSystemHandler(
		&adapters.SystemServiceAdapter{Svc: systemmod.NewSystemService()},
		&adapters.SystemAdminFinderAdapter{UserRepo: s.userRepo},
		&adapters.SystemTokenGenAdapter{TokenService: s.authMod.TokenService},
	)
	updaterHandler := systemmod.NewUpdaterHandler(&adapters.UpdaterServiceAdapter{Svc: updaterService})
	systemmod.RegisterRoutes(authGroup, apiGroup, systemmod.RoutesConfig{
		System:    systemHandler,
		Updater:   updaterHandler,
		AdminAuth: a.authGuard.RequireRole("admin"),
	})
}

func (a *App) registerSettingsRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	settingsHandler := settingsmodule.NewSettingsHandler(s.settingsService)
	appearanceHandler := users.NewAppearanceHandler(s.userRepo)
	settingsmodule.RegisterRoutes(apiGroup, authGroup, settingsHandler, a.authGuard.RequireRole("admin"), appearanceHandler)
	aimod.RegisterRoutes(authGroup, apiGroup, aimod.NewHandler(s.aiSettingsService), a.authGuard.RequireRole("admin"), s.aiRateLimiter.Middleware)

	notifSettingsHandler := notifsmod.NewNotificationSettingsHandler(&adapters.NotifSettingsServiceAdapter{
		Svc:     s.notifSettingsService,
		OrgRepo: s.orgRepo,
	})
	notifsmod.RegisterRoutes(apiGroup, authGroup, notifsmod.NewNotificationHandler(&adapters.NotificationServiceAdapter{}), notifSettingsHandler, a.authGuard.RequireRole("admin"))
}

func (a *App) registerBillingRoutes(apiGroup, authGroup *echo.Group, s *appServices) {
	var billingHandler *billing.BillingHandler
	if config.Get().Cloud.Enabled {
		billingHandler = billing.NewBillingHandler(billing.NewBillingService(adapters.NewBillingUserAdapter(s.userRepo)))
	}
	billing.RegisterRoutes(apiGroup, authGroup, billingHandler, config.Get().Cloud.Enabled)
}
