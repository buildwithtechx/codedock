package http

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"codedock.run/codedock/internal/config"
	"codedock.run/codedock/internal/core"
	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/engine/backup"
	"codedock.run/codedock/internal/engine/compose"
	"codedock.run/codedock/internal/engine/cron"
	"codedock.run/codedock/internal/engine/deploy"
	"codedock.run/codedock/internal/engine/networking"
	"codedock.run/codedock/internal/engine/observability"
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
	projectservices "codedock.run/codedock/internal/services/projects"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

func NewServer(db *sql.DB, v *utils.Vault, deployer *deploy.Deployer, traefikManager *networking.TraefikManager, dockerClient *client.Client, dataDir string) (*Server, error) {

	e := echo.New()
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			slog.Info("request", "method", v.Method, "uri", v.URI, "status", v.Status)
			return nil
		},
	}))
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.GzipWithConfig(echomiddleware.GzipConfig{
		Level: 5,
	}))

	allowOrigins := []string{"http://localhost:3000", "http://localhost:8080"}
	if dashboardURL := config.Get().Server.DashboardURL; dashboardURL != "" && !slices.Contains(allowOrigins, dashboardURL) {
		allowOrigins = append(allowOrigins, dashboardURL)
	}

	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	e.Use(echomiddleware.CSRFWithConfig(echomiddleware.CSRFConfig{
		TokenLength:  32,
		TokenLookup:  "header:X-CSRF-Token",
		CookieName:   "csrf_token",
		CookiePath:   "/",
		CookieMaxAge: 86400,
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			if strings.HasPrefix(path, "/api/auth/signin") ||
				strings.HasPrefix(path, "/api/auth/signup") ||
				strings.HasPrefix(path, "/api/auth/refresh") ||
				strings.HasPrefix(path, "/api/auth/oauth") ||
				strings.HasPrefix(path, "/api/v1/auth/") {
				_, err := c.Cookie("csrf_token")
				return err != nil
			}
			return false
		},
	}))

	environmentRepo := repositories.NewEnvironmentRepo(db)
	projectRepo := repositories.NewProjectRepo(db, environmentRepo)
	appRepo := repositories.NewAppServiceRepo(db)
	serviceVarRepo := repositories.NewServiceVarRepo(db)
	dbRepo := repositories.NewDatabaseRepo(db, v)
	settingsRepo := repositories.NewSettingsRepo(db)
	notifRepo := repositories.NewNotificationSettingsRepo(db)
	aiRepo := repositories.NewAISettingsRepo(db)
	envVarRepo := repositories.NewEnvRepo(db, v)
	scheduledTaskRepo := repositories.NewScheduledTaskRepo(db)
	backupRepo := repositories.NewBackupRepo(db, v)
	s3DestinationRepo := repositories.NewS3DestinationRepo(db, v)
	serverlessRepository := repositories.NewServerlessRepository(db)
	projectSettingsRepo := repositories.NewProjectSettingsRepo(db)
	userRepo := repositories.NewUserRepo(db)
	canvasRepo := repositories.NewCanvasRepo(db, environmentRepo)
	deployRepo := repositories.NewDeploymentRepo(db)
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
	databaseDeployer := deploy.NewDatabaseDeployer(dockerClient, httpEngineAdapter)

	cronManager := cron.NewCronManager(dockerClient, httpEngineAdapter)

	settings, _ := settingsRepo.GetServerSettings(context.Background())
	if settings != nil && settings.DockerCleanupCron != "" {
		_ = cronManager.ScheduleDockerCleanup(settings.DockerCleanupCron)
	}
	if settings != nil && settings.DiskUsageCron != "" {
		_ = cronManager.ScheduleDiskUsageCheck(settings.DiskUsageCron, settings.DiskUsageThreshold)
	}

	_ = cronManager.Start()

	backupManager := backup.NewBackupManager(dockerClient, httpEngineAdapter, "")
	_ = backupManager.Start()

	projectService := projectservices.NewProjectService(projectRepo, environmentRepo, appRepo, serviceVarRepo, settingsRepo, orgRepo)
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
	mailerService, err := notifications.NewMailerService(notifSettingsService)
	if err != nil {
		return nil, fmt.Errorf("mailer service: %w", err)
	}
	authService := authservices.NewAuthService(userRepo, settingsRepo, notifRepo, projectSettingsRepo, tokenService, mailerService, refreshTokenRepo)
	projectSettingsService := projectservices.NewProjectSettingsService(projectSettingsRepo, userRepo, authService)
	dispatcherService := core.NewDispatcherService(settingsRepo, notifRepo, userRepo, mailerService)

	deploymentListeners := core.NewDeploymentListeners(dispatcherService, appRepo)
	deploymentListeners.Register()

	serverRepo := repositories.NewServerRepository(db)
	workerHub := engine.NewWorkerHub(serverRepo)

	scheduledTaskService := systemservices.NewScheduledTaskService(scheduledTaskRepo, cronManager)
	canvasService := projectservices.NewCanvasService(canvasRepo)
	orgService := authservices.NewOrganizationService(orgRepo, userRepo)
	gitService := deploymentservices.NewGitService(gitRepo)
	statsMonitor := observability.NewStatsMonitor(dockerClient)
	deploymentService := deploymentservices.NewDeploymentService(deployRepo, appRepo, projectRepo, deployer, gitService, statsMonitor, volumeRepo, workerHub)
	aiAnalysisService := projectservices.NewAIAnalysisService(deployRepo, appRepo, aiRepo)

	autoscaler := deploy.NewAutoscalerWorker(appRepo, statsMonitor, deploymentService)
	autoscaler.Start()

	backupService := backupservices.NewBackupService(backupRepo, s3DestinationRepo, backupManager)
	userService := authservices.NewUserService(userRepo)
	oAuthService := authservices.NewOAuthService(oauthRepo, userRepo, tokenService)
	prPreviewService := deploymentservices.NewPRPreviewService(prPreviewRepository, appService, gitService, deployer, workerHub, projectRepo)
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
	updaterService.Start(context.Background())

	bridge := NewBridge(projectService, appService, databaseService, deploymentService)

	authGuard := middleware.NewAuthGuard(tokenService, settingsService, projectSettingsService, orgRepo, projectRepo, userRepo)

	appHandler := projects.NewAppHandler(appService, projectService, deployer, deploymentService, environmentService)
	databaseHandler := databases.NewDatabaseHandler(databaseService, projectService, auditService)
	scheduledTaskHandler := system.NewScheduledTaskHandler(scheduledTaskService, appService, projectService)
	canvasHandler := projects.NewCanvasHandler(canvasService, projectService)
	terminalHandler := deployments.NewTerminalHandler(dockerClient, tokenService, appService, projectService, userRepo)
	projectHandler := projects.NewProjectHandler(projectService, projectSettingsService)
	orgHandler := auth.NewOrganizationHandler(orgService)
	environmentHandler := projects.NewEnvironmentHandler(environmentService, projectService)
	deploymentHandler := deployments.NewDeploymentHandler(deploymentService, appService, auditService, aiAnalysisService, prPreviewService, projectService)
	serviceVarHandler := projects.NewServiceVarHandler(appService, auditService, envSuggestionService)
	projectSettingsHandler := projects.NewProjectSettingsHandler(projectSettingsService)
	backupHandler := backups.NewBackupHandler(backupService, appService, databaseService, projectService)
	settingsHandler := auth.NewSettingsHandler(settingsService, notifSettingsService)
	notifSettingsHandler := system.NewNotificationSettingsHandler(notifSettingsService)
	aiSettingsHandler := system.NewAISettingsHandler(aiSettingsService)
	updaterHandler := system.NewUpdaterHandler(updaterService)
	userHandler := auth.NewUserHandler(userService, mailerService)
	authHandler := auth.NewAuthHandler(authService)
	oAuthHandler := auth.NewOAuthHandler(oAuthService)
	gitHandler := deployments.NewGitHandler(gitService)
	webhookHandler := deployments.NewWebhookHandler(gitService, projectService, appService, deploymentService, prPreviewService, gitAppsService)

	domainHandler := projects.NewDomainHandler(environmentService, appService, projectService)
	projectEnvHandler := projects.NewProjectEnvHandler(environmentService)
	notificationHandler := system.NewNotificationHandler(notificationService)
	gitAppsHandler := deployments.NewGitAppsHandler(gitAppsService)
	tmplMgr, _ := compose.NewTemplateManager()
	composeParserService := projectservices.NewComposeParserService()
	composeHandler := projects.NewComposeHandler(projectService, appService, databaseService, environmentRepo, appRepo, composeParserService)
	oneClickService := projectservices.NewOneClickService(tmplMgr, databaseDeployer, environmentRepo, dbRepo)
	oneClickHandler := projects.NewOneClickHandler(oneClickService, projectService)
	archiveService := deploymentservices.NewArchiveService(appService, deploymentService)
	archiveHandler := deployments.NewArchiveHandler(archiveService, projectService)
	serverlessHandler := projects.NewServerlessHandler(serverlessService)
	systemService := systemservices.NewSystemService()
	systemHandler := system.NewSystemHandler(systemService)
	migrationService := systemservices.NewMigrationService(dbRepo, dataDir)
	migrationHandler := system.NewMigrationHandler(migrationService)
	onboardingHandler := auth.NewOnboardingHandler(userService)
	dnsHandler := system.NewDNSHandler(dnsService)
	metricsHandler := system.NewMetricsHandler(metricsService)
	logHandler := system.NewLogHandler(logService)
	auditLogHandler := auth.NewAuditLogHandler(auditService)
	exampleService := systemservices.NewExampleService()
	exampleHandler := system.NewExampleHandler(exampleService)

	serverService := systemservices.NewServerService(serverRepo, userRepo)
	serverHandler := system.NewServerHandler(serverService)
	workerWSHandler := system.NewWorkerWSHandler(workerHub, serverRepo, userRepo)
	serverMetricsWSHandler := system.NewServerMetricsWSHandler(tokenService, serverService, userRepo)
	serviceLogsWSHandler := system.NewServiceLogsWSHandler(tokenService, appService, projectService, userRepo)

	registryRepo := repositories.NewRegistryRepository(db)
	registryService := deploymentservices.NewRegistryService(registryRepo)
	registryHandler := deployments.NewRegistryHandler(registryService)

	billingService := systemservices.NewBillingService(userRepo)
	billingHandler := system.NewBillingHandler(billingService)

	takeoverRepo := repositories.NewTakeoverRepository(db)
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
		orgHandler:             orgHandler,
		environmentHandler:     environmentHandler,
		domainHandler:          domainHandler,
		projectEnvHandler:      projectEnvHandler,
		notificationHandler:    notificationHandler,
		gitAppsHandler:         gitAppsHandler,
		serverlessHandler:      serverlessHandler,
		systemHandler:          systemHandler,
		composeHandler:         composeHandler,
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
		workerWSHandler:        workerWSHandler,
		registryHandler:        registryHandler,
		billingHandler:         billingHandler,
		serverMetricsWSHandler: serverMetricsWSHandler,
		serviceLogsWSHandler:   serviceLogsWSHandler,
		takeoverHandler:        takeoverHandler,
		routeRuleHandler:       routeRuleHandler,
	}

	if srv.deployer != nil {
		srv.deployer.EnvProvider = func(projectID string) (map[string]string, error) {
			return srv.serviceLinker.GetLinkedEnvironmentVariables(context.Background(), projectID)
		}
		srv.deployer.EnvInterpolator = func(projectID string) (map[string]map[string]string, error) {
			return srv.serviceLinker.GetNamespacedVariables(context.Background(), projectID)
		}
		srv.deployer.RouteRuleFetcher = func(ctx context.Context, serviceID, serviceName string) (map[string]string, error) {
			rules, err := routeRuleRepo.ListByService(ctx, serviceID)
			if err != nil {
				return nil, err
			}
			return networking.BuildMiddlewareLabels(serviceName, rules), nil
		}
	}

	srv.registerRoutes()
	return srv, nil
}
