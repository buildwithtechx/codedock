package http

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/engine/deploy"
	"codedock/internal/handlers/deployments"
	"codedock/internal/handlers/projects"
	"codedock/internal/handlers/system"
	"codedock/internal/http/mcp"
	"codedock/internal/http/middleware"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/chatapps"
	"codedock/internal/modules/databases"
	"codedock/internal/modules/dns"
	mailmodule "codedock/internal/modules/mail"
	"codedock/internal/modules/organizations"
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
	deploymentservices "codedock/internal/services/deployments"
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
	gitService             *deploymentservices.GitService
	deploymentService      *deploymentservices.DeploymentService
	autoscalingService     *projectservices.AutoscalingService
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
	prPreviewService := deploymentservices.NewPRPreviewService(d.prPreviewRepository, d.appService, d.gitService, d.deployer, d.sshManager, d.projectRepo)
	dnsProviderService := dns.NewProviderService(d.settingsRepo)
	environmentService := projectservices.NewEnvironmentService(d.environmentRepo, d.domainRepo, d.envVarRepo, dnsProviderService)
	notificationService := systemservices.NewNotificationService(d.dispatcherService)
	gitAppsService := deploymentservices.NewGitAppsService(d.gitAppRepo)
	serverlessService := projectservices.NewServerlessService(d.serverlessRepository)
	dnsService := dns.NewService(d.dnsRepo, dnsProviderService)
	envSuggestionService := projectservices.NewEnvSuggestionService(d.gitService)
	metricsService := systemservices.NewMetricsService()
	logService := systemservices.NewLogService()

	updaterService := systemservices.NewUpdaterService(d.settingsRepo)
	schedulerElector := leadership.NewElector(d.db, "scheduled-execution")
	go schedulerElector.Run(context.Background(), d.cronManager.ServeElected, d.backupManager.ServeElected, updaterService.ServeElected)

	retentionService := deploymentservices.NewRetentionService(d.projectRepo, d.deployRepo, d.deploymentService.DockerForProject)
	deploymentHandler := deployments.NewDeploymentHandler(d.deploymentService, d.appService, d.auditService, d.aiAnalysisService, prPreviewService, d.projectService)
	deploymentHandler.SetRetentionService(retentionService)
	rollbackSettingsHandler := projects.NewRollbackSettingsHandler(retentionService)

	routeProtectionRepo := repositories.NewRouteProtectionRepository(d.db)
	routeProtectionService := projectservices.NewRouteProtectionService(routeProtectionRepo, d.projectRepo, d.settingsRepo)
	routeProtectionHandler := projects.NewRouteProtectionHandler(routeProtectionService, d.projectService)
	edgeAuthHandler := system.NewEdgeAuthHandler(routeProtectionService)
	_ = networking.EnsureTraefikAccessRoute(networking.DynamicDir, config.Get().Server.Port)

	backupHandler := backupmodule.NewHandler(d.backupService)
	backupHandler.SetOperations(d.operationService)

	webhookDeliveryService := deploymentservices.NewWebhookDeliveryService(d.webhookDeliveryRepo)
	webhookDispatcher := deploymentservices.NewWebhookDispatcher(webhookDeliveryService, d.deploymentService, d.appService, gitAppsService, prPreviewService)
	webhookHandler := deployments.NewWebhookHandler(d.gitService, d.projectService, d.appService, d.deploymentService, prPreviewService, gitAppsService, webhookDispatcher, webhookDeliveryService)

	tmplMgr, err := compose.NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("load catalogue templates: %w", err)
	}
	stackService := projectservices.NewComposeStackService(repositories.NewComposeStackRepo(d.db, d.v), compose.NewStackRuntime(d.dockerClient))
	stackService.BeforeDeployment = d.backupService.BeforeProjectDeployment
	if err := stackService.Recover(context.Background()); err != nil {
		return nil, fmt.Errorf("recover compose stacks: %w", err)
	}
	composeParserService := projectservices.NewComposeParserService()
	oneClickService := projectservices.NewOneClickService(tmplMgr, stackService, d.environmentRepo)
	archiveService := deploymentservices.NewArchiveService(d.appService, d.deploymentService)
	migrationService := systemservices.NewMigrationService(d.db, d.dbRepo, d.dataDir, config.Get().Database.URL)
	serverService := systemservices.NewServerService(d.serverRepo, d.userRepo, d.sshManager)

	securityVerifier, err := security.NewVerifier("")
	if err != nil {
		slog.Error("failed to initialize security verifier", "error", err)
	}
	securityFeedManager := security.NewFeedManager(d.db, securityVerifier, "")
	securityService := security.NewSecurityService(d.db, securityFeedManager, nil, nil, nil, nil)

	chatAppService := chatapps.NewService(chatapps.NewRepository(d.db), d.v)
	chatAppDispatcher := chatapps.NewDispatcher(chatAppService, nil)

	terminalHandler := deployments.NewTerminalHandler(d.dockerClient, d.tokenService, d.appService, d.projectService, d.userRepo)
	terminalHandler.Targets = d.deploymentService
	serviceLogsWSHandler := system.NewServiceLogsWSHandler(d.tokenService, d.appService, d.projectService, d.userRepo)
	serviceLogsWSHandler.Streams = d.deploymentService

	var billingHandler *system.BillingHandler
	if config.Get().Cloud.Enabled {
		billingHandler = system.NewBillingHandler(systemservices.NewBillingService(d.userRepo))
	}
	routeRuleRepo := repositories.NewRouteRuleRepository(d.db)

	srv := &Server{
		operationHandler:        system.NewOperationHandler(d.operationService),
		router:                  d.e,
		mcpBridge:               d.bridge,
		authRateLimiter:         middleware.NewRateLimiter(10, time.Minute),
		otpRateLimiter:          middleware.NewRateLimiter(5, time.Minute),
		aiRateLimiter:           middleware.NewRateLimiter(5, time.Minute),
		deployer:                d.deployer,
		traefikManager:          d.traefikManager,
		dockerClient:            d.dockerClient,
		tokenService:            d.tokenService,
		authGuard:               d.authGuard,
		cronManager:             d.cronManager,
		serviceLinker:           d.serviceLinker,
		dispatcherService:       d.dispatcherService,
		projectService:          d.projectService,
		appService:              d.appService,
		autoscalingHandler:      projects.NewAutoscalingHandler(d.autoscalingService),
		appServiceHandler:       projects.NewAppHandler(d.appService, d.projectService, d.deployer, d.deploymentService, environmentService, d.auditService),
		dbHandler:               databases.NewHandler(d.databaseService),
		scheduledTaskHandler:    system.NewScheduledTaskHandler(d.scheduledTaskService, d.appService, d.projectService),
		canvasHandler:           projects.NewCanvasHandler(d.canvasService, d.projectService),
		terminalHandler:         terminalHandler,
		deploymentHandler:       deploymentHandler,
		serviceVarHandler:       projects.NewServiceVarHandler(d.appService, d.auditService, envSuggestionService),
		projectSettingsHandler:  projects.NewProjectSettingsHandler(d.projectSettingsService),
		backupHandler:           backupHandler,
		settingsHandler:         settingsmodule.NewSettingsHandler(d.settingsService),
		appearanceHandler:       users.NewAppearanceHandler(d.userModRepo),
		notifSettingsHandler:    system.NewNotificationSettingsHandler(d.notifSettingsService, d.orgRepo),
		aiSettingsHandler:       system.NewAISettingsHandler(d.aiSettingsService),
		updaterHandler:          system.NewUpdaterHandler(updaterService),
		userHandler:             users.NewUserHandler(d.userService, d.mailerService, d.auditService),
		authHandler:             d.authMod.AuthHandler,
		oauthHandler:            d.authMod.OAuthHandler,
		deviceFlowHandler:       d.authMod.DeviceFlowHandler,
		passkeyHandler:          d.authMod.PasskeyHandler,
		gitHandler:              deployments.NewGitHandler(d.gitService),
		webhookHandler:          webhookHandler,
		projectHandler:          projects.NewProjectHandler(d.projectService, d.projectSettingsService),
		projectAppHandler:       projects.NewProjectAppHandler(d.projectAppService, d.projectService),
		orgHandler:              organizations.NewOrganizationHandler(d.orgService, d.auditService),
		environmentHandler:      projects.NewEnvironmentHandler(environmentService, d.projectService),
		domainHandler:           projects.NewDomainHandler(environmentService, d.appService, d.projectService, d.settingsRepo, d.auditService),
		projectEnvHandler:       projects.NewProjectEnvHandler(environmentService, d.appService),
		notificationHandler:     system.NewNotificationHandler(notificationService),
		gitAppsHandler:          deployments.NewGitAppsHandler(gitAppsService),
		serverlessHandler:       projects.NewServerlessHandler(serverlessService),
		systemHandler:           system.NewSystemHandler(systemservices.NewSystemService(), d.userRepo, d.tokenService),
		composeHandler:          projects.NewComposeHandler(d.projectService, d.appService, d.databaseService, d.environmentRepo, d.appRepo, composeParserService),
		composeStackHandler:     projects.NewComposeStackHandler(stackService, d.projectService, environmentService),
		chatAppHandler:          chatapps.NewHandler(chatAppService, chatAppDispatcher),
		oneClickHandler:         projects.NewOneClickHandler(oneClickService, d.projectService),
		archiveHandler:          deployments.NewArchiveHandler(archiveService, d.projectService),
		migrationHandler:        system.NewMigrationHandler(migrationService, d.userService),
		onboardingHandler:       users.NewOnboardingHandler(d.userService),
		dnsHandler:              dns.NewHandler(dnsService),
		metricsHandler:          system.NewMetricsHandler(metricsService),
		logHandler:              system.NewLogHandler(logService),
		auditLogHandler:         audit.NewAuditLogHandler(d.auditService),
		exampleHandler:          system.NewExampleHandler(systemservices.NewExampleService()),
		serverHandler:           system.NewServerHandler(serverService, d.auditService),
		sshManager:              d.sshManager,
		registryHandler:         deployments.NewRegistryHandler(deploymentservices.NewRegistryService(repositories.NewRegistryRepository(d.db))),
		billingHandler:          billingHandler,
		serverMetricsWSHandler:  system.NewServerMetricsWSHandler(d.tokenService, serverService, d.userRepo),
		serviceLogsWSHandler:    serviceLogsWSHandler,
		takeoverHandler:         system.NewTakeoverHandler(systemservices.NewTakeoverScanner(), systemservices.NewTakeoverAdopter(d.projectRepo, d.appRepo), repositories.NewTakeoverRepository(d.db, d.v)),
		routeRuleHandler:        projects.NewRouteRuleHandler(systemservices.NewRouteRuleService(routeRuleRepo), d.appRepo),
		mailHandler:             mailmodule.NewHandler(mailmodule.NewService(mailmodule.NewRepository(d.db, d.v)), nil, nil),
		ingressHandler:          system.NewIngressHandlerWithDeps(repositories.NewIngressRepository(d.db), systemservices.NewCloudflareIngressService(nil)),
		serverNetworksHandler:   system.NewServerNetworksHandler(),
		actionsHandler:          system.NewActionsHandler(),
		checksHandler:           projects.NewChecksHandler(),
		supportHandler:          system.NewSupportHandler(),
		rollbackSettingsHandler: rollbackSettingsHandler,
		routeProtectionHandler:  routeProtectionHandler,
		edgeAuthHandler:         edgeAuthHandler,
		securityHandler:         system.NewSecurityHandler(securityService),
	}

	if err := configureManaged(srv, d.db, d.v, d.serverRepo, d.userRepo, d.orgRepo, d.sshManager, d.operationService); err != nil {
		return nil, err
	}
	if err := configureMigration(srv, d.db, d.v, d.projectRepo, d.appService, d.appRepo, d.environmentRepo, d.serverRepo, d.userRepo, d.orgRepo, d.gitService, composeParserService, d.deploymentService, d.sshManager, d.runtimeRepository); err != nil {
		return nil, err
	}
	configureAnalytics(srv, d.db, d.v, d.traefikManager, d.deployRepo, d.projectRepo, d.appRepo, d.domainRepo, d.userRepo, d.orgRepo, metricsService, d.statsMonitor, d.deploymentService, d.backupService, srv.migrationService, d.auditService)
	if err := configureClusters(srv, d.db, d.v, d.projectRepo, d.serverRepo, d.operationService, d.volumeOperations, d.appRepo, d.deploymentService, d.canvasService, d.backupManager); err != nil {
		return nil, err
	}
	configureDeploymentBindings(srv, routeRuleRepo, d.domainRepo, routeProtectionRepo)

	srv.registerRoutes()
	return srv, nil
}
