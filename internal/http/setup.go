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
	aimod "codedock/internal/modules/ai"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/monitoring"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/servers"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/users"
	"codedock/internal/platform/networking"
	"codedock/internal/platform/notifications"
	"codedock/internal/platform/observability"
	cron "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
	backup "codedock/internal/platform/storage"
	"codedock/internal/platform/vault"
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

	environmentRepo := environments.NewEnvironmentRepo(db)
	projectRepo := projectsmod.NewProjectRepo(db)
	projectAppRepo := projectsmod.NewProjectAppRepo(db)
	appRepo := appsmod.NewAppServiceRepo(db)
	serviceVarRepo := appsmod.NewServiceVarRepo(db)
	dbRepo := databases.NewRepository(db, v)
	settingsRepo := settingsmodule.NewSettingsRepo(db, v)
	notifRepo := monitoring.NewNotificationSettingsRepo(db)
	aiRepo := aimod.NewRepository(db)
	envVarRepo := projectsmod.NewEnvRepo(db, v)
	scheduledTaskRepo := monitoring.NewScheduledTaskRepo(db)
	backupRepo := backupmodule.NewRepository(db, v)
	backupRepo.SetSnapshotDumper(backupmodule.NewSystemSnapshotDumper(dataDir))
	serverlessRepository := appsmod.NewServerlessRepository(db)
	projectSettingsRepo := projectsmod.NewProjectSettingsRepo(db)
	userRepo := users.NewUserRepo(db)
	canvasRepo := projectsmod.NewCanvasRepo(db, &projectEnvAdapter{repo: environmentRepo})
	canvasRepo.SetVault(v)
	canvasRepo.SetClusterData(&canvasClusterDataAdapter{repo: modclusters.NewClusterDataRepo(db, v)})
	deployRepo := deploymentsmod.NewDeploymentRepository(db)
	if err := deployRepo.RecoverInterrupted(context.Background()); err != nil {
		return nil, fmt.Errorf("recover deployments: %w", err)
	}
	prPreviewRepository := gitmod.NewPRPreviewRepository(db)
	domainRepo := domainsmod.NewDomainRepo(db)
	gitAppRepo := gitmod.NewGitAppRepository(db, v)
	dnsRepo := dns.NewRepository(db)
	volumeRepo := appsmod.NewServiceVolumeRepo(db)
	orgRepo := organizations.NewOrganizationRepository(db)
	webhookDeliveryRepo := gitmod.NewWebhookDeliveryRepository(db)

	httpEngineAdapter := newEngineAdapter(appRepo, dbRepo, projectRepo, scheduledTaskRepo, backupRepo)
	httpEngineAdapter.sftpRepo = backupmodule.NewSFTPRepository(db)
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

	projectService := projectsmod.NewService(projectRepo, orgRepo)
	projectAppService := projectsmod.NewProjectAppService(projectAppRepo)
	appService := appsmod.NewService(appRepo, serviceVarRepo, volumeRepo)
	databaseService := databases.NewService(dbRepo, databases.NewDockerDeployer(dockerClient, dbRepo, volumeOperations))
	deviceTokenRepo := monitoring.NewDeviceTokenRepo(db)
	notifSettingsService := monitoring.NewNotificationSettingsService(notifRepo, monitoring.NewNotificationDeliveryRepo(db), deviceTokenRepo)
	aiSettingsService := aimod.NewService(aiRepo)
	serviceLinker := appsmod.NewServiceLinker(&appDatabaseAdapter{repo: dbRepo})
	serviceLinker.SetApplications(appRepo, serviceVarRepo)
	serviceLinker.SetClusterData(&appClusterPlansAdapter{repo: modclusters.NewClusterDataRepo(db, v)})
	mailerService, err := notifications.NewMailerService(notifSettingsService)
	if err != nil {
		return nil, fmt.Errorf("mailer service: %w", err)
	}
	userService := users.NewUserService(userRepo)
	settingsService := settingsmodule.NewSettingsService(settingsRepo)
	auditRepo := audit.NewAuditLogRepo(db)
	auditService := audit.NewAuditService(auditRepo, audit.NewSQLAuditResourceSearcher(db))
	orgService := organizations.NewOrganizationService(orgRepo)

	authMod, err := auth.NewModule(db, userRepo, userService, settingsRepo, auditService, mailerService)
	if err != nil {
		return nil, err
	}
	tokenService := authMod.TokenService
	projectSettingsService := projectsmod.NewProjectSettingsService(projectSettingsRepo)
	dispatcherService := notifications.NewDispatcherService(
		&notifSettingsAdapter{repo: settingsRepo},
		notifRepo,
		&notifUserLookupAdapter{repo: userRepo},
		monitoring.NewNotificationDeliveryRepo(db),
		mailerService,
		&notifDeviceTokenAdapter{repo: deviceTokenRepo},
	)

	deploymentListeners := notifications.NewDeploymentListeners(dispatcherService, &notifAppRepoAdapter{repo: appRepo})
	deploymentListeners.Register()

	serverRepo := servers.NewRepository(db, v)
	sshManager := ssh.NewSSHManager(&sshServerRepoAdapter{repo: serverRepo})
	sshManager.OnServerUnreachable = func(server *ssh.Server, err error) {
		baseURL := strings.TrimRight(config.Get().Server.DashboardURL, "/")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
		dispatcherService.Dispatch(&notifications.NotificationEvent{
			EventType: notifications.EventServerUnreachable,
			Title:     "Server unreachable: " + server.Name,
			Message:   fmt.Sprintf("Cannot reach Docker daemon on server %s (%s): %v", server.Name, server.IPAddress, err),
			Level:     "error",
			URL:       fmt.Sprintf("%s/servers/%s", baseURL, server.ID),
		})
	}

	scheduledTaskService := monitoring.NewScheduledTaskService(scheduledTaskRepo, cronManager)
	canvasService := projectsmod.NewCanvasService(canvasRepo)
	canvasService.SetRuntime(&canvasRuntimeAdapter{rt: observability.NewCanvasRuntime(dockerClient)})

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

	aiAnalysisService := aimod.NewAnalysisService(aiRepo)
	runtimeRepository := appsmod.NewRuntimeRepo(db, v)
	if err := runtimeRepository.SyncKinds(context.Background()); err != nil {
		return nil, fmt.Errorf("recover runtime kinds: %w", err)
	}
	deploymentService.SetRuntime(runtimeRepository)

	autoscalingRepo := appsmod.NewWorkerAutoscalingRepo(db)
	autoscalingService := appsmod.NewAutoscalingService(appsmod.NewAutoscalingRepository(db), serverlessRepository, appRepo)
	autoscaler := deploy.NewAutoscalerWorker(appRepo, statsMonitor, &autoscalerDeploymentAdapter{ds: deploymentService}, autoscalingRepo)

	operationService := monitoring.NewOperationService(monitoring.NewOperationRepo(db, v))
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
		dispatcherService.Dispatch(&notifications.NotificationEvent{
			EventType: notifications.EventBackupFailed,
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
	authGuard.PersonalResources = users.NewPersonalTokenResources(db)

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
		userModRepo:            userRepo,
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
