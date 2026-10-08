package http

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"time"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/core"
	"codedock.run/codedock/internal/engine/backup"
	"codedock.run/codedock/internal/engine/compose"
	"codedock.run/codedock/internal/engine/cron"
	"codedock.run/codedock/internal/engine/deploy"
	"codedock.run/codedock/internal/engine/leadership"
	"codedock.run/codedock/internal/engine/networking"
	"codedock.run/codedock/internal/engine/observability"
	"codedock.run/codedock/internal/engine/ssh"
	"codedock.run/codedock/internal/handlers/auth"
	"codedock.run/codedock/internal/handlers/backups"
	"codedock.run/codedock/internal/handlers/databases"
	"codedock.run/codedock/internal/handlers/deployments"
	"codedock.run/codedock/internal/handlers/projects"
	"codedock.run/codedock/internal/handlers/system"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/notifications"
	"codedock.run/codedock/internal/repositories"
	authservices "codedock.run/codedock/internal/services/auth"
	backupservices "codedock.run/codedock/internal/services/backups"
	databaseservices "codedock.run/codedock/internal/services/databases"
	deploymentservices "codedock.run/codedock/internal/services/deployments"
	"codedock.run/codedock/internal/services/operations"
	projectservices "codedock.run/codedock/internal/services/projects"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

func NewServer(db *sql.DB, v *utils.Vault, deployer *deploy.Deployer, traefikManager *networking.TraefikManager, dockerClient *client.Client, dataDir string) (*Server, error) {

	e := echo.New()
	configureEchoMiddleware(e)
	if deployer != nil {
		if err := deployer.SetRolloutDirectory(filepath.Join(dataDir, "rollouts")); err != nil {
			return nil, err
		}
		recoveryCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		err := deployer.RecoverRollouts(recoveryCtx)
		cancel()
		if err != nil {
			return nil, fmt.Errorf("recover interrupted rollout: %w", err)
		}
	}

	environmentRepo := repositories.NewEnvironmentRepo(db)
	projectRepo := repositories.NewProjectRepo(db, environmentRepo)
	projectAppRepo := repositories.NewProjectAppRepo(db)
	appRepo := repositories.NewAppServiceRepo(db)
	serviceVarRepo := repositories.NewServiceVarRepo(db)
	dbRepo := repositories.NewDatabaseRepo(db, v)
	settingsRepo := repositories.NewSettingsRepo(db, v)
	notifRepo := repositories.NewNotificationSettingsRepo(db)
	aiRepo := repositories.NewAISettingsRepo(db)
	envVarRepo := repositories.NewEnvRepo(db, v)
	scheduledTaskRepo := repositories.NewScheduledTaskRepo(db)
	backupRepo := repositories.NewBackupRepo(db, v)
	backupRepo.SetSnapshotDumper(newSystemSnapshotDumper(dataDir))
	s3DestinationRepo := repositories.NewS3DestinationRepo(db, v)
	serverlessRepository := repositories.NewServerlessRepository(db)
	projectSettingsRepo := repositories.NewProjectSettingsRepo(db)
	userRepo := repositories.NewUserRepo(db)
	canvasRepo := repositories.NewCanvasRepo(db, environmentRepo)
	canvasRepo.SetVault(v)
	canvasRepo.SetClusterData(repositories.NewClusterDataRepo(db, v))
	deployRepo := repositories.NewDeploymentRepo(db)
	if err := deployRepo.RecoverInterrupted(context.Background()); err != nil {
		return nil, fmt.Errorf("recover deployments: %w", err)
	}
	oauthRepo := repositories.NewOAuthRepo(db)
	gitRepo := repositories.NewGitRepo(db, v)
	prPreviewRepository := repositories.NewPRPreviewRepository(db)
	domainRepo := repositories.NewDomainRepo(db)
	gitAppRepo := repositories.NewGitAppRepo(db, v)
	dnsRepo := repositories.NewDNSRepo(db)
	auditRepository := repositories.NewAuditLogRepo(db)
	volumeRepo := repositories.NewServiceVolumeRepo(db)
	orgRepo := repositories.NewOrganizationRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepo(db)

	httpEngineAdapter := newEngineAdapter(settingsRepo, appRepo, envVarRepo, dbRepo, projectRepo, scheduledTaskRepo, backupRepo, s3DestinationRepo, serviceVarRepo, serverlessRepository)
	httpEngineAdapter.sftpRepo = repositories.NewSFTPDestinationRepo(db)
	volumeOperations := deploy.NewVolumeGate()
	if deployer != nil {
		deployer.SetVolumeOperations(volumeOperations)
	}
	databaseDeployer := deploy.NewDatabaseDeployer(dockerClient, httpEngineAdapter)
	databaseDeployer.SetVolumeOperations(volumeOperations)

	cronManager := cron.NewCronManager(dockerClient, httpEngineAdapter)

	settings, _ := settingsRepo.GetServerSettings(context.Background())
	if settings != nil && settings.DockerCleanupCron != "" {
		_ = cronManager.ScheduleDockerCleanup(settings.DockerCleanupCron)
	}
	if settings != nil && settings.DiskUsageCron != "" {
		_ = cronManager.ScheduleDiskUsageCheck(settings.DiskUsageCron, settings.DiskUsageThreshold)
	}

	backupManager := backup.NewBackupManager(dockerClient, httpEngineAdapter, "")
	backupManager.SetVolumeOperations(volumeOperations)

	projectService := projectservices.NewProjectService(projectRepo, environmentRepo, appRepo, serviceVarRepo, settingsRepo, orgRepo)
	projectAppService := projectservices.NewProjectAppService(projectAppRepo)
	appService := projectservices.NewAppService(appRepo, serviceVarRepo, volumeRepo)
	databaseService := databaseservices.NewDatabaseService(dbRepo, databaseDeployer)
	tokenService, err := authservices.NewTokenService()
	if err != nil {
		return nil, fmt.Errorf("token service: %w", err)
	}
	settingsService := authservices.NewSettingsService(settingsRepo)
	notifSettingsService := systemservices.NewNotificationSettingsService(notifRepo)
	aiSettingsService := projectservices.NewAISettingsService(aiRepo)
	serviceLinker := projectservices.NewServiceLinker(dbRepo)
	serviceLinker.SetApplications(appRepo, serviceVarRepo)
	serviceLinker.SetClusterData(repositories.NewClusterDataRepo(db, v))
	mailerService, err := notifications.NewMailerService(notifSettingsService)
	if err != nil {
		return nil, fmt.Errorf("mailer service: %w", err)
	}
	authService := authservices.NewAuthService(userRepo, settingsRepo, notifRepo, projectSettingsRepo, tokenService, mailerService, refreshTokenRepo)
	projectSettingsService := projectservices.NewProjectSettingsService(projectSettingsRepo, userRepo, authService)
	dispatcherService := core.NewDispatcherService(settingsRepo, notifRepo, userRepo, mailerService)

	deploymentListeners := core.NewDeploymentListeners(dispatcherService, appRepo)
	deploymentListeners.Register()

	serverRepo := repositories.NewServerRepository(db, v)
	sshManager := ssh.NewSSHManager(serverRepo)

	scheduledTaskService := systemservices.NewScheduledTaskService(scheduledTaskRepo, cronManager)
	canvasService := projectservices.NewCanvasService(canvasRepo)
	canvasService.SetRuntime(observability.NewCanvasRuntime(dockerClient))
	orgService := authservices.NewOrganizationService(orgRepo, userRepo)
	gitService := deploymentservices.NewGitService(gitRepo)
	statsMonitor := observability.NewStatsMonitor(dockerClient)
	deploymentService := deploymentservices.NewDeploymentService(deployRepo, appRepo, projectRepo, deployer, gitService, statsMonitor, volumeRepo, sshManager)
	deploymentService.LocalDocker = dockerClient
	deploymentService.Databases = dbRepo
	aiAnalysisService := projectservices.NewAIAnalysisService(deployRepo, appRepo, aiRepo)

	runtimeRepository := repositories.NewRuntimeRepo(db, v)
	if err := runtimeRepository.SyncKinds(context.Background()); err != nil {
		return nil, fmt.Errorf("recover runtime kinds: %w", err)
	}
	deploymentService.RuntimeKinds = runtimeRepository
	autoscalingRepo := repositories.NewAutoscalingRepo(db)
	autoscalingService := projectservices.NewAutoscalingService(autoscalingRepo, appRepo, projectRepo)
	autoscalingService.SetRuntimes(runtimeRepository)
	autoscalingHandler := projects.NewAutoscalingHandler(autoscalingService)
	autoscaler := deploy.NewAutoscalerWorker(appRepo, statsMonitor, deploymentService, autoscalingRepo)

	operationService := operations.NewService(repositories.NewOperationRepo(db, v))
	if err := operationService.Recover(context.Background()); err != nil {
		return nil, fmt.Errorf("recover operations: %w", err)
	}
	backupService := backupservices.NewBackupService(backupRepo, s3DestinationRepo, backupManager)
	backupService.SetBatches(repositories.NewPolicyBatchRepo(db))
	backupService.SetSFTP(repositories.NewSFTPDestinationRepo(db))
	backupManager.SetServiceRuntime(func(ctx context.Context, serviceID string) (string, error) {
		runtime, err := runtimeRepository.Get(ctx, serviceID)
		if err != nil {
			return "", err
		}
		return runtime.Target.Kind, nil
	})
	if err := backupRepo.RecoverRecords(context.Background()); err != nil {
		return nil, fmt.Errorf("recover backup records: %w", err)
	}
	configureRemoteBackups(backupManager, projectRepo, serverRepo, sshManager)
	backupService.SetOperations(operationService)
	configureScheduledBackupAuthorization(backupService, userRepo, projectService)
	deploymentService.SetRemoteTargets(serverRepo, volumeOperations)
	deploymentService.BeforeProjectDeployment = backupService.BeforeProjectDeployment
	deploymentService.BeforeDeployment = backupService.BeforeDeployment
	autoscaler.Start()
	userService := authservices.NewUserService(userRepo)
	oAuthService := authservices.NewOAuthService(oauthRepo, userRepo, tokenService)
	prPreviewService := deploymentservices.NewPRPreviewService(prPreviewRepository, appService, gitService, deployer, sshManager, projectRepo)
	dnsProviderService := systemservices.NewDNSProviderService(settingsRepo)
	environmentService := projectservices.NewEnvironmentService(environmentRepo, domainRepo, envVarRepo, dnsProviderService)
	notificationService := systemservices.NewNotificationService(dispatcherService)
	gitAppsService := deploymentservices.NewGitAppsService(gitAppRepo)
	serverlessService := projectservices.NewServerlessService(serverlessRepository)
	dnsService := systemservices.NewDNSService(dnsRepo, dnsProviderService)
	envSuggestionService := projectservices.NewEnvSuggestionService(gitService)
	metricsService := systemservices.NewMetricsService()
	logService := systemservices.NewLogService()
	auditService := authservices.NewAuditService(auditRepository)

	updaterService := systemservices.NewUpdaterService(settingsRepo)
	schedulerElector := leadership.NewElector(db, "scheduled-execution")
	go schedulerElector.Run(context.Background(), cronManager.ServeElected, backupManager.ServeElected, updaterService.ServeElected)

	bridge := NewBridge(projectService, appService, databaseService, deploymentService)

	authGuard := middleware.NewAuthGuard(tokenService, settingsService, projectSettingsService, orgRepo, projectRepo, userRepo)
	authGuard.PersonalTokens = userRepo
	authGuard.PersonalResources = repositories.NewPersonalTokenResources(db)

	appHandler := projects.NewAppHandler(appService, projectService, deployer, deploymentService, environmentService)
	databaseHandler := databases.NewDatabaseHandler(databaseService, projectService, auditService)
	scheduledTaskHandler := system.NewScheduledTaskHandler(scheduledTaskService, appService, projectService)
	canvasHandler := projects.NewCanvasHandler(canvasService, projectService)
	terminalHandler := deployments.NewTerminalHandler(dockerClient, tokenService, appService, projectService, userRepo)
	terminalHandler.Targets = deploymentService
	projectHandler := projects.NewProjectHandler(projectService, projectSettingsService)
	projectAppHandler := projects.NewProjectAppHandler(projectAppService, projectService)
	orgHandler := auth.NewOrganizationHandler(orgService)
	environmentHandler := projects.NewEnvironmentHandler(environmentService, projectService)
	deploymentHandler := deployments.NewDeploymentHandler(deploymentService, appService, auditService, aiAnalysisService, prPreviewService, projectService)
	serviceVarHandler := projects.NewServiceVarHandler(appService, auditService, envSuggestionService)
	projectSettingsHandler := projects.NewProjectSettingsHandler(projectSettingsService)
	backupHandler := backups.NewBackupHandler(backupService, appService, databaseService, projectService)
	backupHandler.SetOperations(operationService)
	settingsHandler := auth.NewSettingsHandler(settingsService, notifSettingsService)
	notifSettingsHandler := system.NewNotificationSettingsHandler(notifSettingsService)
	aiSettingsHandler := system.NewAISettingsHandler(aiSettingsService)
	updaterHandler := system.NewUpdaterHandler(updaterService)
	userHandler := auth.NewUserHandler(userService, mailerService)
	authHandler := auth.NewAuthHandler(authService)
	oAuthHandler := auth.NewOAuthHandler(oAuthService)
	gitHandler := deployments.NewGitHandler(gitService)
	webhookHandler := deployments.NewWebhookHandler(gitService, projectService, appService, deploymentService, prPreviewService, gitAppsService)

	domainHandler := projects.NewDomainHandler(environmentService, appService, projectService, settingsRepo)
	projectEnvHandler := projects.NewProjectEnvHandler(environmentService)
	notificationHandler := system.NewNotificationHandler(notificationService)
	gitAppsHandler := deployments.NewGitAppsHandler(gitAppsService)
	tmplMgr, err := compose.NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("load catalogue templates: %w", err)
	}
	stackService := projectservices.NewComposeStackService(repositories.NewComposeStackRepo(db, v), compose.NewStackRuntime(dockerClient))
	stackService.BeforeDeployment = backupService.BeforeProjectDeployment
	if err := stackService.Recover(context.Background()); err != nil {
		return nil, fmt.Errorf("recover compose stacks: %w", err)
	}
	composeStackHandler := projects.NewComposeStackHandler(stackService, projectService, environmentService)
	composeParserService := projectservices.NewComposeParserService()
	composeHandler := projects.NewComposeHandler(projectService, appService, databaseService, environmentRepo, appRepo, composeParserService)
	oneClickService := projectservices.NewOneClickService(tmplMgr, stackService, environmentRepo)
	oneClickHandler := projects.NewOneClickHandler(oneClickService, projectService)
	archiveService := deploymentservices.NewArchiveService(appService, deploymentService)
	archiveHandler := deployments.NewArchiveHandler(archiveService, projectService)
	serverlessHandler := projects.NewServerlessHandler(serverlessService)
	systemService := systemservices.NewSystemService()
	systemHandler := system.NewSystemHandler(systemService)
	migrationService := systemservices.NewMigrationService(db, dbRepo, dataDir, config.Get().Database.URL)
	migrationHandler := system.NewMigrationHandler(migrationService, userService)
	onboardingHandler := auth.NewOnboardingHandler(userService)
	dnsHandler := system.NewDNSHandler(dnsService)
	metricsHandler := system.NewMetricsHandler(metricsService)
	logHandler := system.NewLogHandler(logService)
	auditLogHandler := auth.NewAuditLogHandler(auditService)
	exampleService := systemservices.NewExampleService()
	exampleHandler := system.NewExampleHandler(exampleService)

	serverService := systemservices.NewServerService(serverRepo, userRepo, sshManager)
	serverHandler := system.NewServerHandler(serverService)
	serverMetricsWSHandler := system.NewServerMetricsWSHandler(tokenService, serverService, userRepo)
	serviceLogsWSHandler := system.NewServiceLogsWSHandler(tokenService, appService, projectService, userRepo)
	serviceLogsWSHandler.Streams = deploymentService

	registryRepo := repositories.NewRegistryRepository(db)
	registryService := deploymentservices.NewRegistryService(registryRepo)
	registryHandler := deployments.NewRegistryHandler(registryService)

	var billingHandler *system.BillingHandler
	if config.Get().Cloud.Enabled {
		billingService := systemservices.NewBillingService(userRepo)
		billingHandler = system.NewBillingHandler(billingService)
	}

	takeoverRepo := repositories.NewTakeoverRepository(db, v)
	takeoverScanner := systemservices.NewTakeoverScanner()
	takeoverAdopter := systemservices.NewTakeoverAdopter(projectRepo, appRepo)
	takeoverHandler := system.NewTakeoverHandler(takeoverScanner, takeoverAdopter, takeoverRepo)

	routeRuleRepo := repositories.NewRouteRuleRepository(db)
	routeRuleService := systemservices.NewRouteRuleService(routeRuleRepo)
	routeRuleHandler := projects.NewRouteRuleHandler(routeRuleService, appRepo)

	authLimiter := middleware.NewRateLimiter(10, time.Minute)
	otpLimiter := middleware.NewRateLimiter(5, time.Minute)
	aiLimiter := middleware.NewRateLimiter(5, time.Minute)

	srv := &Server{
		operationHandler:       system.NewOperationHandler(operationService),
		router:                 e,
		mcpBridge:              bridge,
		authRateLimiter:        authLimiter,
		otpRateLimiter:         otpLimiter,
		aiRateLimiter:          aiLimiter,
		deployer:               deployer,
		traefikManager:         traefikManager,
		dockerClient:           dockerClient,
		tokenService:           tokenService,
		authGuard:              authGuard,
		cronManager:            cronManager,
		serviceLinker:          serviceLinker,
		dispatcherService:      dispatcherService,
		projectService:         projectService,
		appService:             appService,
		autoscalingHandler:     autoscalingHandler,
		appServiceHandler:      appHandler,
		dbHandler:              databaseHandler,
		scheduledTaskHandler:   scheduledTaskHandler,
		canvasHandler:          canvasHandler,
		terminalHandler:        terminalHandler,
		deploymentHandler:      deploymentHandler,
		serviceVarHandler:      serviceVarHandler,
		projectSettingsHandler: projectSettingsHandler,
		backupHandler:          backupHandler,
		settingsHandler:        settingsHandler,
		notifSettingsHandler:   notifSettingsHandler,
		aiSettingsHandler:      aiSettingsHandler,
		updaterHandler:         updaterHandler,
		userHandler:            userHandler,
		authHandler:            authHandler,
		oauthHandler:           oAuthHandler,
		gitHandler:             gitHandler,
		webhookHandler:         webhookHandler,
		projectHandler:         projectHandler,
		projectAppHandler:      projectAppHandler,
		orgHandler:             orgHandler,
		environmentHandler:     environmentHandler,
		domainHandler:          domainHandler,
		projectEnvHandler:      projectEnvHandler,
		notificationHandler:    notificationHandler,
		gitAppsHandler:         gitAppsHandler,
		serverlessHandler:      serverlessHandler,
		systemHandler:          systemHandler,
		composeHandler:         composeHandler,
		composeStackHandler:    composeStackHandler,
		oneClickHandler:        oneClickHandler,
		archiveHandler:         archiveHandler,
		migrationHandler:       migrationHandler,
		onboardingHandler:      onboardingHandler,
		dnsHandler:             dnsHandler,
		metricsHandler:         metricsHandler,
		logHandler:             logHandler,
		auditLogHandler:        auditLogHandler,
		exampleHandler:         exampleHandler,
		serverHandler:          serverHandler,
		sshManager:             sshManager,
		registryHandler:        registryHandler,
		billingHandler:         billingHandler,
		serverMetricsWSHandler: serverMetricsWSHandler,
		serviceLogsWSHandler:   serviceLogsWSHandler,
		takeoverHandler:        takeoverHandler,
		routeRuleHandler:       routeRuleHandler,
	}

	if err := configureManaged(srv, db, v, serverRepo, userRepo, orgRepo, sshManager, operationService); err != nil {
		return nil, err
	}
	if err := configureMigration(srv, db, v, projectRepo, appService, appRepo, environmentRepo, serverRepo, userRepo, orgRepo, gitService, composeParserService, deploymentService, sshManager, runtimeRepository); err != nil {
		return nil, err
	}
	configureAnalytics(srv, db, v, traefikManager, deployRepo, projectRepo, appRepo, domainRepo, userRepo, orgRepo, metricsService, statsMonitor, deploymentService, backupService, srv.migrationService)
	if err := configureClusters(srv, db, v, projectRepo, serverRepo, operationService, volumeOperations, appRepo, deploymentService, canvasService, backupManager); err != nil {
		return nil, err
	}
	configureDeploymentBindings(srv, routeRuleRepo)

	srv.registerRoutes()
	return srv, nil
}
