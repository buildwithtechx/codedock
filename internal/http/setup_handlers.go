package http

import (
	"context"
	"log/slog"
	"time"

	"codedock/internal/config"
	"codedock/internal/http/middleware"
	actions "codedock/internal/modules/actions"
	aimod "codedock/internal/modules/ai"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/billing"
	"codedock/internal/modules/canvas"
	"codedock/internal/modules/chatapps"
	modclusters "codedock/internal/modules/clusters"
	composemod "codedock/internal/modules/compose"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	mailmodule "codedock/internal/modules/mail"
	modmigrations "codedock/internal/modules/migrations"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/security"
	"codedock/internal/modules/servers/managed"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/support"
	"codedock/internal/modules/system"
	templatesmod "codedock/internal/modules/templates"
	"codedock/internal/modules/users"
	compose "codedock/internal/platform/docker"
	"codedock/internal/platform/networking"
	leadership "codedock/internal/platform/scheduler"
)

func buildServerInstance(d *serverSetupDeps) (*Server, error) {
	prPreviewRepo := gitmod.NewPRPreviewRepository(d.db)
	prPreviewService := gitmod.NewPRPreviewService(prPreviewRepo, &gitServiceLookupAdapter{appRepo: d.appRepo})
	dnsProviderService := dns.NewProviderService(&dnsSettingsAdapter{repo: d.settingsRepo})
	dnsService := dns.NewService(d.dnsRepo, dnsProviderService)
	gitAppsRepo := gitmod.NewGitAppRepository(d.db, d.v)
	gitAppsService := gitmod.NewGitAppsService(gitAppsRepo)
	webhookRepo := gitmod.NewWebhookDeliveryRepository(d.db)
	webhookDeliveryService := gitmod.NewWebhookDeliveryService(webhookRepo)
	webhookDispatcher := gitmod.NewWebhookDispatcher(webhookDeliveryService, d.deploymentService, &gitServiceLookupAdapter{appRepo: d.appRepo}, gitAppsService)
	webhookHandler := gitmod.NewWebhookHandler(webhookDispatcher, webhookDeliveryService, prPreviewService, &gitServiceLookupAdapter{appRepo: d.appRepo})
	gitProvidersHandler := gitmod.NewProvidersHandler(d.gitService)
	gitInspectionHandler := gitmod.NewInspectionHandler(d.gitService)
	gitAppsHandler := gitmod.NewAppsHandler(gitAppsService)

	retentionService := deploymentsmod.NewRetentionService(deploymentsmod.NewDeploymentRepository(d.db))
	_ = retentionService
	deploymentHandler := deploymentsmod.NewDeploymentHandler(d.deploymentService, &appLookupAdapter{appService: d.appService}, &projectLookupAdapter{projectService: d.projectService}, &aiExplainerAdapter{svc: d.aiAnalysisService}, &prPreviewLookupAdapter{pr: prPreviewService})
	lifecycleHandler := deploymentsmod.NewLifecycleHandler(d.deploymentService, &appLookupAdapter{appService: d.appService}, &auditLoggerAdapter{audit: d.auditService})
	registryRepo := deploymentsmod.NewRegistryRepository(d.db)
	registryService := deploymentsmod.NewRegistryService(registryRepo)
	registryHandler := deploymentsmod.NewRegistryHandler(registryService)
	archiveService := deploymentsmod.NewArchiveService(d.deploymentService)
	archiveHandler := deploymentsmod.NewArchiveHandler(archiveService, &projectLookupAdapter{projectService: d.projectService})
	terminalHandler := deploymentsmod.NewTerminalHandler(d.dockerClient, d.tokenService, &appLookupAdapter{appService: d.appService}, &projectLookupAdapter{projectService: d.projectService}, &userStatusAdapter{us: d.userService})
	terminalHandler.Targets = d.deploymentService

	appsService := appsmod.NewService(d.appRepo, appsmod.NewServiceVarRepo(d.db), appsmod.NewServiceVolumeRepo(d.db))
	appAccess := &appAccessChecker{ps: d.projectService}
	envLookup := &envLookupAdapter{repo: d.environmentRepo}
	appsHandler := appsmod.NewHandler(appsService, appAccess, envLookup)
	appsHandler.SetDeployer(&appLifecycleDeployerAdapter{ds: d.deploymentService, appService: d.appService})
	appVarHandler := appsmod.NewVariableHandler(appsService, appAccess)
	appVolumeHandler := appsmod.NewVolumeHandler(appsService, appAccess)
	appWebhookHandler := appsmod.NewWebhookHandler(appsService)
	autoscalingService := appsmod.NewAutoscalingService(appsmod.NewAutoscalingRepository(d.db), appsmod.NewServerlessRepository(d.db), d.appRepo)
	autoscalingHandler := appsmod.NewAutoscalingHandler(autoscalingService)

	projService := projectsmod.NewService(d.projectRepo, d.orgRepo)
	projSettingsService := projectsmod.NewProjectSettingsService(projectsmod.NewProjectSettingsRepo(d.db))
	projectHandler := projectsmod.NewHandler(projService, projSettingsService, d.orgService)
	projectAppHandler := projectsmod.NewProjectAppHandler(d.projectAppService, projService)
	projectSettingsHandler := projectsmod.NewProjectSettingsHandler(projSettingsService)
	checksHandler := projectsmod.NewChecksHandler()
	canvasHandler := canvas.NewCanvasHandler(&canvasServiceAdapter{cs: d.canvasService})

	stackService := composemod.NewComposeStackService(composemod.NewComposeStackRepo(d.db, d.v), &composeStackRuntimeAdapter{rt: compose.NewStackRuntime(d.dockerClient)})
	stackService.BeforeDeployment = d.backupService.BeforeProjectDeployment
	_ = stackService.Recover(context.Background())
	composeStackHandler := composemod.NewComposeStackHandler(&composeStackAdapter{cs: stackService})
	tmplMgr, _ := compose.NewTemplateManager()
	oneClickService := templatesmod.NewOneClickService(tmplMgr, stackService, &templateEnvAdapter{repo: d.environmentRepo})
	oneClickHandler := templatesmod.NewOneClickHandler(&oneClickAdapter{oc: oneClickService})

	envService := environments.NewEnvironmentService(d.environmentRepo, environments.NewEnvVarRepo())
	envHandler := environments.NewEnvironmentHandler(envService, &projectLookupAdapter{projectService: d.projectService}, environments.NewEnvSuggestionService(&envGitProviderAdapter{ps: d.gitService}))
	protectionService := environments.NewRouteProtectionService(environments.NewRouteProtectionRepository(d.db), &projectLookupAdapter{projectService: d.projectService})
	protectionHandler := environments.NewRouteProtectionHandler(protectionService)
	edgeAuthHandler := environments.NewEdgeAuthHandler(protectionService)
	routeRuleRepo := environments.NewRouteRuleRepository(d.db)
	routeRuleHandler := environments.NewRouteRuleHandler(environments.NewRouteRuleService(routeRuleRepo), &serviceCheckerAdapter{appService: d.appService})
	domainHandler := domainsmod.NewDomainHandler(domainsmod.NewDomainService(d.domainRepo, dnsProviderService), appAccess, &domainServiceLookupAdapter{appRepo: d.appRepo}, &serverSettingsAdapter{repo: d.settingsRepo}, &auditDomainAdapter{audit: d.auditService})

	metricsService := monitoring.NewMetricsService()
	logService := monitoring.NewLogService()
	updaterService := system.NewUpdaterService(d.settingsRepo)
	schedulerElector := leadership.NewElector(d.db, "scheduled-execution")
	go schedulerElector.Run(context.Background(), d.cronManager.ServeElected, d.backupManager.ServeElected, updaterService.ServeElected)

	backupHandler := backupmodule.NewHandler(d.backupService)
	backupHandler.SetOperations(d.operationService)
	migrationService := modmigrations.NewBundleService(d.db, d.dbRepo, d.dataDir, config.Get().Database.URL)
	serversHandler, serversWSHandler, serversTakeoverHandler, serversNetworksHandler := createServerHandlers(d)
	securityVerifier, err := security.NewVerifier("")
	if err != nil {
		slog.Error("failed to initialize security verifier", "error", err)
	}
	securityService := security.NewSecurityService(d.db, security.NewFeedManager(d.db, securityVerifier, ""), nil, nil, nil, nil)
	chatAppService := chatapps.NewService(chatapps.NewRepository(d.db), d.v)
	monHandlers := createMonitoringHandlers(d, updaterService, metricsService, logService)
	var billingHandler *billing.BillingHandler
	if config.Get().Cloud.Enabled {
		billingHandler = billing.NewBillingHandler(billing.NewBillingService(&billingUserAdapter{repo: d.userRepo}))
	}
	supportHandler := support.NewSupportHandler()
	_ = networking.EnsureTraefikAccessRoute(networking.DynamicDir, config.Get().Server.Port)

	srv := &Server{
		operationHandler:           monHandlers.operationHandler,
		router:                     d.e,
		mcpBridge:                  d.bridge,
		authRateLimiter:            middleware.NewRateLimiter(10, time.Minute),
		otpRateLimiter:             middleware.NewRateLimiter(5, time.Minute),
		aiRateLimiter:              middleware.NewRateLimiter(5, time.Minute),
		deployer:                   d.deployer,
		traefikManager:             d.traefikManager,
		dockerClient:               d.dockerClient,
		tokenService:               d.tokenService,
		authGuard:                  d.authGuard,
		cronManager:                d.cronManager,
		serviceLinker:              d.serviceLinker,
		dispatcherService:          d.dispatcherService,
		projectService:             d.projectService,
		appService:                 d.appService,
		appServiceHandler:          appsHandler,
		appVarHandler:              appVarHandler,
		appVolumeHandler:           appVolumeHandler,
		appWebhookHandler:          appWebhookHandler,
		autoscalingHandler:         autoscalingHandler,
		dbHandler:                  databases.NewHandler(d.databaseService),
		scheduledTaskHandler:       monHandlers.scheduledTaskHandler,
		canvasHandler:              canvasHandler,
		terminalHandler:            terminalHandler,
		deploymentHandler:          deploymentHandler,
		deploymentLifecycleHandler: lifecycleHandler,
		archiveHandler:             archiveHandler,
		registryHandler:            registryHandler,
		projectSettingsHandler:     projectSettingsHandler,
		backupHandler:              backupHandler,
		settingsHandler:            settingsmodule.NewSettingsHandler(d.settingsService),
		appearanceHandler:          users.NewAppearanceHandler(d.userModRepo),
		notifSettingsHandler:       monHandlers.notifSettingsHandler,
		aiSettingsHandler:          aimod.NewHandler(d.aiSettingsService),
		updaterHandler:             monHandlers.updaterHandler,
		userHandler:                users.NewUserHandler(d.userService, d.mailerService, d.auditService),
		authHandler:                d.authMod.AuthHandler,
		oauthHandler:               d.authMod.OAuthHandler,
		deviceFlowHandler:          d.authMod.DeviceFlowHandler,
		passkeyHandler:             d.authMod.PasskeyHandler,
		gitHandler:                 gitProvidersHandler,
		gitInspectionHandler:       gitInspectionHandler,
		gitAppsHandler:             gitAppsHandler,
		webhookHandler:             webhookHandler,
		projectHandler:             projectHandler,
		projectAppHandler:          projectAppHandler,
		orgHandler:                 organizations.NewOrganizationHandler(d.orgService, d.auditService),
		environmentHandler:         envHandler,
		domainHandler:              domainHandler,
		notificationHandler:        monHandlers.notificationHandler,
		systemHandler:              monHandlers.systemHandler,
		composeStackHandler:        composeStackHandler,
		chatAppHandler:             chatapps.NewHandler(chatAppService, chatapps.NewDispatcher(chatAppService, nil)),
		oneClickHandler:            oneClickHandler,
		migrationHandler:           modmigrations.NewBundleHandler(&migrationBundleAdapter{service: migrationService}, &migrationUserCounterAdapter{users: d.userRepo}, middleware.NewRateLimiter(5, time.Minute)),
		onboardingHandler:          users.NewOnboardingHandler(d.userService),
		dnsHandler:                 dns.NewHandler(dnsService),
		metricsHandler:             monHandlers.metricsHandler,
		logHandler:                 monHandlers.logHandler,
		auditLogHandler:            audit.NewAuditLogHandler(d.auditService),
		exampleHandler:             monHandlers.exampleHandler,
		serverHandler:              serversHandler,
		sshManager:                 d.sshManager,
		billingHandler:             billingHandler,
		serverMetricsWSHandler:     serversWSHandler,
		serviceLogsWSHandler:       monHandlers.serviceLogsWSHandler,
		takeoverHandler:            serversTakeoverHandler,
		routeRuleHandler:           routeRuleHandler,
		mailHandler:                mailmodule.NewHandler(mailmodule.NewService(mailmodule.NewRepository(d.db, d.v)), nil, nil),
		ingressHandler:             monHandlers.ingressHandler,
		serverNetworksHandler:      serversNetworksHandler,
		actionsHandler:             actions.NewActionsHandler(),
		checksHandler:              checksHandler,
		supportHandler:             supportHandler,
		routeProtectionHandler:     protectionHandler,
		edgeAuthHandler:            edgeAuthHandler,
		securityHandler:            security.NewHandler(securityService),
	}

	srv.managedHandler = managed.NewManagedHandler(nil)
	srv.migrationLifecycleHandler = modmigrations.NewLifecycleHandler(nil)
	accessLogPath := ""
	if d.traefikManager != nil {
		accessLogPath = d.traefikManager.AccessLogPath()
	}
	analyticsServices := monitoring.InitAnalyticsServices(
		d.db,
		d.v,
		accessLogPath,
		&analyticsDeployRepoAdapter{repo: d.deployRepo},
		&analyticsProjectRepoAdapter{repo: d.projectRepo},
		&analyticsAppRepoAdapter{repo: d.appRepo},
		&analyticsDomainListerAdapter{repo: d.domainRepo},
		&analyticsUserRepoAdapter{repo: d.userRepo},
		&analyticsOrgRepoAdapter{repo: d.orgRepo},
		metricsService,
		d.statsMonitor,
		&attentionDeployerAdapter{deploymentService: d.deploymentService, appService: d.appService},
		d.backupService,
		nil,
	)
	srv.analyticsHandler = monitoring.NewAnalyticsHandler(&analyticsServiceAdapter{svc: analyticsServices.AnalyticsService})
	srv.attentionHandler = monitoring.NewAttentionHandler(&attentionServiceAdapter{svc: analyticsServices.AttentionService})
	clusterServices, err := modclusters.InitClusterServices(d.db, d.v, d.projectRepo, d.serverRepo, d.operationService, d.volumeOperations)
	if err != nil {
		return nil, err
	}
	srv.clusterHandler = modclusters.NewHandler(&clusterServiceAdapter{
		service:    clusterServices.Service,
		operations: d.operationService,
	})
	srv.clusterDataHandler = modclusters.NewDataHandler(&clusterDataServiceAdapter{
		service:    clusterServices.DataService,
		operations: d.operationService,
		apps:       d.appRepo,
		vars:       appsmod.NewServiceVarRepo(d.db),
	})
	configureDeploymentBindings(srv, environments.NewRouteRuleRepository(d.db), d.domainRepo, environments.NewRouteProtectionRepository(d.db))

	srv.registerRoutes()
	return srv, nil
}
