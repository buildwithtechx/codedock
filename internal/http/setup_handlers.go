package http

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/engine/deploy"
	"codedock/internal/handlers/system"
	"codedock/internal/http/mcp"
	"codedock/internal/http/middleware"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/chatapps"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	mailmodule "codedock/internal/modules/mail"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/security"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/users"
	compose "codedock/internal/platform/docker"
	"codedock/internal/platform/networking"
	"codedock/internal/platform/notifications"
	core "codedock/internal/platform/notifications"
	"codedock/internal/platform/observability"
	cron "codedock/internal/platform/scheduler"
	leadership "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
	backup "codedock/internal/platform/storage"
	"codedock/internal/platform/vault"
	"codedock/internal/repositories"
	"codedock/internal/services/operations"
	projectservices "codedock/internal/services/projects"
	systemservices "codedock/internal/services/system"
)

type serverSetupDeps struct {
	e                      *echo.Echo
	db                     *sql.DB
	v                      *vault.Vault
	dataDir                string
	deployer               *deploy.Deployer
	traefikManager         *networking.TraefikManager
	dockerClient           *client.Client
	bridge                 *mcp.Bridge
	authGuard              *middleware.AuthGuard
	tokenService           *auth.TokenService
	cronManager            *cron.CronManager
	backupManager          *backup.BackupManager
	sshManager             *ssh.SSHManager
	statsMonitor           *observability.StatsMonitor
	volumeOperations       *deploy.VolumeGate
	serviceLinker          *projectservices.ServiceLinker
	dispatcherService      *core.DispatcherService
	projectService         *projectservices.ProjectService
	projectAppService      *projectservices.ProjectAppService
	appService             *projectservices.AppService
	databaseService        *databases.Service
	notifSettingsService   *systemservices.NotificationSettingsService
	aiSettingsService      *projectservices.AISettingsService
	aiAnalysisService      *projectservices.AIAnalysisService
	userService            *users.UserService
	settingsService        *settingsmodule.SettingsService
	auditService           *audit.AuditService
	orgService             *organizations.OrganizationService
	authMod                *auth.Module
	projectSettingsService *projectservices.ProjectSettingsService
	scheduledTaskService   *systemservices.ScheduledTaskService
	canvasService          *projectservices.CanvasService
	gitService             *gitmod.ProviderService
	deploymentService      *deploymentsmod.DeploymentService
	autoscalingService     *appsmod.AutoscalingService
	operationService       *operations.Service
	backupService          *backupmodule.Service
	mailerService          *notifications.MailerService
	environmentRepo        *repositories.EnvironmentRepo
	projectRepo            *repositories.ProjectRepo
	appRepo                *repositories.AppServiceRepo
	dbRepo                 *repositories.DatabaseRepo
	settingsRepo           *repositories.SettingsRepo
	envVarRepo             *repositories.EnvRepo
	serverlessRepository   repositories.ServerlessRepository
	userRepo               *repositories.UserRepo
	userModRepo            users.UserRepository
	deployRepo             *repositories.DeploymentRepo
	prPreviewRepository    repositories.PRPreviewRepository
	domainRepo             *repositories.DomainRepo
	gitAppRepo             *repositories.GitAppRepo
	dnsRepo                dns.Repository
	orgRepo                repositories.OrganizationRepository
	webhookDeliveryRepo    repositories.WebhookDeliveryRepository
	serverRepo             repositories.ServerRepository
	runtimeRepository      *repositories.RuntimeRepo
}

func buildServerInstance(d *serverSetupDeps) (*Server, error) {
	prPreviewRepo := gitmod.NewPRPreviewRepository(d.db)
	prPreviewService := gitmod.NewPRPreviewService(prPreviewRepo, &gitServiceLookupAdapter{appRepo: d.appRepo})
	dnsProviderService := dns.NewProviderService(d.settingsRepo)
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

	appsService := appsmod.NewService(appsmod.NewAppServiceRepo(d.db), appsmod.NewServiceVarRepo(d.db), appsmod.NewServiceVolumeRepo(d.db))
	appAccess := &appAccessChecker{ps: d.projectService}
	envLookup := &envLookupAdapter{repo: d.environmentRepo}
	appsHandler := appsmod.NewHandler(appsService, appAccess, envLookup)
	appsHandler.SetDeployer(&appLifecycleDeployerAdapter{ds: d.deploymentService, appService: d.appService})
	appVarHandler := appsmod.NewVariableHandler(appsService, appAccess)
	appVolumeHandler := appsmod.NewVolumeHandler(appsService, appAccess)
	appWebhookHandler := appsmod.NewWebhookHandler(appsService)
	autoscalingService := appsmod.NewAutoscalingService(appsmod.NewAutoscalingRepository(d.db), appsmod.NewServerlessRepository(d.db), appsmod.NewAppServiceRepo(d.db))
	autoscalingHandler := appsmod.NewAutoscalingHandler(autoscalingService)

	projService := projectsmod.NewService(projectsmod.NewProjectRepo(d.db), organizations.NewOrganizationRepository(d.db))
	projSettingsService := projectsmod.NewProjectSettingsService(projectsmod.NewProjectSettingsRepo(d.db))
	projectHandler := projectsmod.NewHandler(projService, projSettingsService, d.orgService)
	projectAppHandler := projectsmod.NewProjectAppHandler(projectsmod.NewProjectAppService(projectsmod.NewProjectAppRepo(d.db)), projService)
	projectSettingsHandler := projectsmod.NewProjectSettingsHandler(projSettingsService)
	checksHandler := projectsmod.NewChecksHandler()
	canvasHandler := projectsmod.NewCanvasHandler(&canvasServiceAdapter{cs: d.canvasService}, projService)

	stackService := projectservices.NewComposeStackService(repositories.NewComposeStackRepo(d.db, d.v), compose.NewStackRuntime(d.dockerClient))
	stackService.BeforeDeployment = d.backupService.BeforeProjectDeployment
	_ = stackService.Recover(context.Background())
	composeStackHandler := projectsmod.NewComposeStackHandler(&composeStackAdapter{cs: stackService})
	composeHandler := projectsmod.NewComposeHandler()
	tmplMgr, _ := compose.NewTemplateManager()
	oneClickService := projectservices.NewOneClickService(tmplMgr, stackService, d.environmentRepo)
	oneClickHandler := projectsmod.NewOneClickHandler(&oneClickAdapter{oc: oneClickService})

	envService := environments.NewEnvironmentService(environments.NewEnvironmentRepo(d.db), environments.NewEnvVarRepo())
	envHandler := environments.NewEnvironmentHandler(envService, &projectLookupAdapter{projectService: d.projectService}, environments.NewEnvSuggestionService(&envGitProviderAdapter{ps: d.gitService}))
	protectionService := environments.NewRouteProtectionService(environments.NewRouteProtectionRepository(d.db), &projectLookupAdapter{projectService: d.projectService})
	protectionHandler := environments.NewRouteProtectionHandler(protectionService)
	edgeAuthHandler := environments.NewEdgeAuthHandler(protectionService)
	routeRuleRepo := environments.NewRouteRuleRepository(d.db)
	routeRuleHandler := environments.NewRouteRuleHandler(environments.NewRouteRuleService(routeRuleRepo), &serviceCheckerAdapter{appService: d.appService})
	domainHandler := domainsmod.NewDomainHandler(domainsmod.NewDomainService(domainsmod.NewDomainRepo(d.db), dnsProviderService), appAccess, &domainServiceLookupAdapter{appRepo: d.appRepo}, &serverSettingsAdapter{repo: d.settingsRepo}, &auditDomainAdapter{audit: d.auditService})

	metricsService := systemservices.NewMetricsService()
	logService := systemservices.NewLogService()
	updaterService := systemservices.NewUpdaterService(d.settingsRepo)
	schedulerElector := leadership.NewElector(d.db, "scheduled-execution")
	go schedulerElector.Run(context.Background(), d.cronManager.ServeElected, d.backupManager.ServeElected, updaterService.ServeElected)

	backupHandler := backupmodule.NewHandler(d.backupService)
	backupHandler.SetOperations(d.operationService)
	migrationService := systemservices.NewMigrationService(d.db, d.dbRepo, d.dataDir, config.Get().Database.URL)
	serverService := systemservices.NewServerService(d.serverRepo, d.userRepo, d.sshManager)
	securityVerifier, err := security.NewVerifier("")
	if err != nil {
		slog.Error("failed to initialize security verifier", "error", err)
	}
	securityService := security.NewSecurityService(d.db, security.NewFeedManager(d.db, securityVerifier, ""), nil, nil, nil, nil)
	chatAppService := chatapps.NewService(chatapps.NewRepository(d.db), d.v)
	serviceLogsWSHandler := system.NewServiceLogsWSHandler(d.tokenService, d.appService, d.projectService, d.userRepo)
	serviceLogsWSHandler.Streams = d.deploymentService

	var billingHandler *system.BillingHandler
	if config.Get().Cloud.Enabled {
		billingHandler = system.NewBillingHandler(systemservices.NewBillingService(d.userRepo))
	}
	_ = networking.EnsureTraefikAccessRoute(networking.DynamicDir, config.Get().Server.Port)

	srv := &Server{
		operationHandler:           system.NewOperationHandler(d.operationService),
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
		scheduledTaskHandler:       system.NewScheduledTaskHandler(d.scheduledTaskService, d.appService, d.projectService),
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
		notifSettingsHandler:       system.NewNotificationSettingsHandler(d.notifSettingsService, d.orgRepo),
		aiSettingsHandler:          system.NewAISettingsHandler(d.aiSettingsService),
		updaterHandler:             system.NewUpdaterHandler(updaterService),
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
		notificationHandler:        system.NewNotificationHandler(systemservices.NewNotificationService(d.dispatcherService)),
		systemHandler:              system.NewSystemHandler(systemservices.NewSystemService(), d.userRepo, d.tokenService),
		composeHandler:             composeHandler,
		composeStackHandler:        composeStackHandler,
		chatAppHandler:             chatapps.NewHandler(chatAppService, chatapps.NewDispatcher(chatAppService, nil)),
		oneClickHandler:            oneClickHandler,
		migrationHandler:           system.NewMigrationHandler(migrationService, d.userService),
		onboardingHandler:          users.NewOnboardingHandler(d.userService),
		dnsHandler:                 dns.NewHandler(dnsService),
		metricsHandler:             system.NewMetricsHandler(metricsService),
		logHandler:                 system.NewLogHandler(logService),
		auditLogHandler:            audit.NewAuditLogHandler(d.auditService),
		exampleHandler:             system.NewExampleHandler(systemservices.NewExampleService()),
		serverHandler:              system.NewServerHandler(serverService, d.auditService),
		sshManager:                 d.sshManager,
		billingHandler:             billingHandler,
		serverMetricsWSHandler:     system.NewServerMetricsWSHandler(d.tokenService, serverService, d.userRepo),
		serviceLogsWSHandler:       serviceLogsWSHandler,
		takeoverHandler:            system.NewTakeoverHandler(systemservices.NewTakeoverScanner(), systemservices.NewTakeoverAdopter(d.projectRepo, d.appRepo), repositories.NewTakeoverRepository(d.db, d.v)),
		routeRuleHandler:           routeRuleHandler,
		mailHandler:                mailmodule.NewHandler(mailmodule.NewService(mailmodule.NewRepository(d.db, d.v)), nil, nil),
		ingressHandler:             system.NewIngressHandlerWithDeps(repositories.NewIngressRepository(d.db), systemservices.NewCloudflareIngressService(nil)),
		serverNetworksHandler:      system.NewServerNetworksHandler(),
		actionsHandler:             system.NewActionsHandler(),
		checksHandler:              checksHandler,
		supportHandler:             system.NewSupportHandler(),
		routeProtectionHandler:     protectionHandler,
		edgeAuthHandler:            edgeAuthHandler,
		securityHandler:            system.NewSecurityHandler(securityService),
	}

	if err := configureManaged(srv, d.db, d.v, d.serverRepo, d.userRepo, d.orgRepo, d.sshManager, d.operationService); err != nil {
		return nil, err
	}
	composeParserService := projectservices.NewComposeParserService()
	if err := configureMigration(srv, d.db, d.v, d.projectRepo, d.appService, d.appRepo, d.environmentRepo, d.serverRepo, d.userRepo, d.orgRepo, &migrationGitAdapter{gitService: d.gitService}, composeParserService, &migrationDeployerAdapter{deploymentService: d.deploymentService}, d.sshManager, d.runtimeRepository); err != nil {
		return nil, err
	}
	configureAnalytics(srv, d.db, d.v, d.traefikManager, d.deployRepo, d.projectRepo, d.appRepo, d.domainRepo, d.userRepo, d.orgRepo, metricsService, d.statsMonitor, &attentionDeployerAdapter{deploymentService: d.deploymentService, appService: d.appService}, d.backupService, srv.migrationService, d.auditService)
	if err := configureClusters(srv, d.db, d.v, d.projectRepo, d.serverRepo, d.operationService, d.volumeOperations, d.appRepo, d.deploymentService, d.canvasService, d.backupManager); err != nil {
		return nil, err
	}
	configureDeploymentBindings(srv, environments.NewRouteRuleRepository(d.db), d.domainRepo, environments.NewRouteProtectionRepository(d.db))

	srv.registerRoutes()
	return srv, nil
}
