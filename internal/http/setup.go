package http

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/engine/deploy"
	"codedock/internal/http/mcp"
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/organizations"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/users"
	"codedock/internal/platform/networking"
	"codedock/internal/platform/notifications"
	core "codedock/internal/platform/notifications"
	"codedock/internal/platform/observability"
	cron "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
	backup "codedock/internal/platform/storage"
	"codedock/internal/platform/vault"
	"codedock/internal/repositories"
	"codedock/internal/services/operations"
	projectservices "codedock/internal/services/projects"
	systemservices "codedock/internal/services/system"
)

func NewServer(db *sql.DB, v *vault.Vault, deployer *deploy.Deployer, traefikManager *networking.TraefikManager, dockerClient *client.Client, dataDir string) (*Server, error) {
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
	backupRepo.SetSnapshotDumper(backupmodule.NewSystemSnapshotDumper(dataDir))
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
	prPreviewRepository := repositories.NewPRPreviewRepository(db)
	domainRepo := repositories.NewDomainRepo(db)
	gitAppRepo := repositories.NewGitAppRepo(db, v)
	dnsRepo := dns.NewRepository(db)
	volumeRepo := repositories.NewServiceVolumeRepo(db)
	orgRepo := repositories.NewOrganizationRepository(db)
	webhookDeliveryRepo := repositories.NewWebhookDeliveryRepo(db)

	httpEngineAdapter := newEngineAdapter(settingsRepo, appRepo, envVarRepo, dbRepo, projectRepo, scheduledTaskRepo, backupRepo, s3DestinationRepo, serviceVarRepo, serverlessRepository)
	httpEngineAdapter.sftpRepo = repositories.NewSFTPDestinationRepo(db)
	volumeOperations := deploy.NewVolumeGate()
	if deployer != nil {
		deployer.SetVolumeOperations(volumeOperations)
	}

	cronManager := cron.NewCronManager(dockerClient, httpEngineAdapter)
	serverSettings, _ := settingsRepo.GetServerSettings(context.Background())
	if serverSettings != nil && serverSettings.DockerCleanupCron != "" {
		_ = cronManager.ScheduleDockerCleanup(serverSettings.DockerCleanupCron)
	}
	if serverSettings != nil && serverSettings.DiskUsageCron != "" {
		_ = cronManager.ScheduleDiskUsageCheck(serverSettings.DiskUsageCron, serverSettings.DiskUsageThreshold)
	}

	backupManager := backup.NewBackupManager(dockerClient, httpEngineAdapter, "")
	backupManager.SetVolumeOperations(volumeOperations)

	projectService := projectservices.NewProjectService(projectRepo, environmentRepo, appRepo, serviceVarRepo, settingsRepo, orgRepo)
	projectAppService := projectservices.NewProjectAppService(projectAppRepo)
	appService := projectservices.NewAppService(appRepo, serviceVarRepo, volumeRepo)
	dbModRepo := databases.NewRepository(db, v)
	databaseService := databases.NewService(dbModRepo, databases.NewDockerDeployer(dockerClient, dbModRepo, volumeOperations))
	deviceTokenRepo := repositories.NewDeviceTokenRepo(db)
	notifSettingsService := systemservices.NewNotificationSettingsService(notifRepo, repositories.NewNotificationDeliveryRepo(db), deviceTokenRepo)
	aiSettingsService := projectservices.NewAISettingsService(aiRepo)
	serviceLinker := projectservices.NewServiceLinker(dbRepo)
	serviceLinker.SetApplications(appRepo, serviceVarRepo)
	serviceLinker.SetClusterData(repositories.NewClusterDataRepo(db, v))
	mailerService, err := notifications.NewMailerService(notifSettingsService)
	if err != nil {
		return nil, fmt.Errorf("mailer service: %w", err)
	}
	userModRepo := users.NewUserRepo(db)
	userService := users.NewUserService(userModRepo)
	settingsModRepo := settingsmodule.NewSettingsRepo(db, v)
	settingsService := settingsmodule.NewSettingsService(settingsModRepo)
	auditRepo := audit.NewAuditLogRepo(db)
	auditService := audit.NewAuditService(auditRepo, audit.NewSQLAuditResourceSearcher(db))
	orgModRepo := organizations.NewOrganizationRepository(db)
	orgService := organizations.NewOrganizationService(orgModRepo)

	authMod, err := auth.NewModule(db, userModRepo, userService, settingsModRepo, auditService, mailerService)
	if err != nil {
		return nil, err
	}
	tokenService := authMod.TokenService
	projectSettingsService := projectservices.NewProjectSettingsService(projectSettingsRepo, userRepo, authMod.AuthService)
	dispatcherService := core.NewDispatcherService(settingsRepo, notifRepo, userRepo, repositories.NewNotificationDeliveryRepo(db), mailerService, deviceTokenRepo)

	deploymentListeners := core.NewDeploymentListeners(dispatcherService, appRepo)
	deploymentListeners.Register()

	serverRepo := repositories.NewServerRepository(db, v)
	sshManager := ssh.NewSSHManager(serverRepo)
	sshManager.OnServerUnreachable = func(server *models.Server, err error) {
		baseURL := strings.TrimRight(config.Get().Server.DashboardURL, "/")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		dispatcherService.Dispatch(&models.NotificationEvent{
			EventType: models.EventServerUnreachable,
			Title:     "Server unreachable: " + server.Name,
			Message:   fmt.Sprintf("Cannot reach Docker daemon on server %s (%s): %v", server.Name, server.IPAddress, err),
			Level:     "error",
			URL:       fmt.Sprintf("%s/servers/%s", baseURL, server.ID),
		})
	}

	scheduledTaskService := systemservices.NewScheduledTaskService(scheduledTaskRepo, cronManager)
	canvasService := projectservices.NewCanvasService(canvasRepo)
	canvasService.SetRuntime(observability.NewCanvasRuntime(dockerClient))

	gitProviderRepo := gitmod.NewGitProviderRepository(db, v)
	gitService := gitmod.NewProviderService(gitProviderRepo)
	statsMonitor := observability.NewStatsMonitor(dockerClient)
	deploymentRepo := deploymentsmod.NewDeploymentRepository(db)
	deploymentService := deploymentsmod.NewDeploymentService(
		deploymentRepo,
		&appLookupAdapter{appService: appService},
		&projectLookupAdapter{projectService: projectService},
		gitService,
		&deployEngineAdapter{deployer: deployer, appService: appService},
		statsMonitor,
	)
	deploymentService.SetLocalDocker(dockerClient)

	aiAnalysisService := projectservices.NewAIAnalysisService(deployRepo, appRepo, aiRepo)
	runtimeRepository := repositories.NewRuntimeRepo(db, v)
	if err := runtimeRepository.SyncKinds(context.Background()); err != nil {
		return nil, fmt.Errorf("recover runtime kinds: %w", err)
	}
	deploymentService.SetRuntime(runtimeRepository)

	autoscalingRepo := repositories.NewAutoscalingRepo(db)
	autoscalingService := appsmod.NewAutoscalingService(appsmod.NewAutoscalingRepository(db), appsmod.NewServerlessRepository(db), appsmod.NewAppServiceRepo(db))
	autoscaler := deploy.NewAutoscalerWorker(appRepo, statsMonitor, &autoscalerDeploymentAdapter{ds: deploymentService}, autoscalingRepo)

	operationService := operations.NewService(repositories.NewOperationRepo(db, v))
	if err := operationService.Recover(context.Background()); err != nil {
		return nil, fmt.Errorf("recover operations: %w", err)
	}
	backupModRepo := backupmodule.NewRepository(db, v)
	backupModRepo.SetSnapshotDumper(backupmodule.NewSystemSnapshotDumper(dataDir))
	backupService := backupmodule.NewService(backupModRepo, backupModRepo, backupmodule.NewManager(backupManager))
	backupService.SetBatches(backupModRepo)
	backupService.SetSFTP(backupmodule.NewSFTPRepository(db))
	backupManager.OnBackupFailed = func(recordID, errStr string) {
		baseURL := strings.TrimRight(config.Get().Server.DashboardURL, "/")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		dispatcherService.Dispatch(&models.NotificationEvent{
			EventType: models.EventBackupFailed,
			Title:     "Backup failed: " + recordID,
			Message:   fmt.Sprintf("Backup execution failed: %s", errStr),
			Level:     "error",
			URL:       fmt.Sprintf("%s/backups", baseURL),
		})
	}
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
	deploymentService.BeforeProjectDeployment = func(ctx context.Context, projectID string) error {
		return backupService.BeforeProjectDeployment(ctx, projectID, "")
	}
	deploymentService.BeforeDeployment = backupService.BeforeDeployment
	autoscaler.Start()

	authGuard := middleware.NewAuthGuard(tokenService, settingsService, projectSettingsService, orgRepo, projectRepo, userRepo)
	authGuard.PersonalTokens = userRepo
	authGuard.PersonalResources = repositories.NewPersonalTokenResources(db)

	return buildServerInstance(&serverSetupDeps{
		e:                      e,
		db:                     db,
		v:                      v,
		dataDir:                dataDir,
		deployer:               deployer,
		traefikManager:         traefikManager,
		dockerClient:           dockerClient,
		bridge:                 mcp.NewBridge(projectService, databaseService),
		authGuard:              authGuard,
		tokenService:           tokenService,
		cronManager:            cronManager,
		backupManager:          backupManager,
		sshManager:             sshManager,
		statsMonitor:           statsMonitor,
		volumeOperations:       volumeOperations,
		serviceLinker:          serviceLinker,
		dispatcherService:      dispatcherService,
		projectService:         projectService,
		projectAppService:      projectAppService,
		appService:             appService,
		databaseService:        databaseService,
		notifSettingsService:   notifSettingsService,
		aiSettingsService:      aiSettingsService,
		aiAnalysisService:      aiAnalysisService,
		userService:            userService,
		settingsService:        settingsService,
		auditService:           auditService,
		orgService:             orgService,
		authMod:                authMod,
		projectSettingsService: projectSettingsService,
		scheduledTaskService:   scheduledTaskService,
		canvasService:          canvasService,
		gitService:             gitService,
		deploymentService:      deploymentService,
		autoscalingService:     autoscalingService,
		operationService:       operationService,
		backupService:          backupService,
		mailerService:          mailerService,
		environmentRepo:        environmentRepo,
		projectRepo:            projectRepo,
		appRepo:                appRepo,
		dbRepo:                 dbRepo,
		settingsRepo:           settingsRepo,
		envVarRepo:             envVarRepo,
		serverlessRepository:   serverlessRepository,
		userRepo:               userRepo,
		userModRepo:            userModRepo,
		deployRepo:             deployRepo,
		prPreviewRepository:    prPreviewRepository,
		domainRepo:             domainRepo,
		gitAppRepo:             gitAppRepo,
		dnsRepo:                dnsRepo,
		orgRepo:                orgRepo,
		webhookDeliveryRepo:    webhookDeliveryRepo,
		serverRepo:             serverRepo,
		runtimeRepository:      runtimeRepository,
	})
}
