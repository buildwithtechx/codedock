package app

import (
	"fmt"
	"net/http"
	"time"

	"codedock/internal/adapters"
	"codedock/internal/deploy"
	"codedock/internal/http/middleware"
	aimod "codedock/internal/modules/ai"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/apps/tasks"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/audit/operations"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	modclusters "codedock/internal/modules/clusters"
	"codedock/internal/modules/clusters/kubernetes"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/mcp"
	modmonitoring "codedock/internal/modules/monitoring"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/projects/canvas"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/users"
	"codedock/internal/platform/notifications"
	"codedock/internal/platform/observability"
	cron "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
	backup "codedock/internal/platform/storage"
)

func initDomainServices(s *appServices, cfg Config) (*appServices, error) {
	// Managers & Volume Gate
	volumeOperations := deploy.NewVolumeGate()
	if cfg.Deployer != nil {
		cfg.Deployer.SetVolumeOperations(volumeOperations)
	}
	httpEngineAdapter := adapters.NewEngineAdapter(s.appRepo, s.dbRepo, s.projectRepo, s.scheduledTaskRepo, s.backupRepo)
	httpEngineAdapter.SetSFTPRepo(backupmodule.NewSFTPRepository(cfg.DB))
	s.cronManager = cron.NewCronManager(cfg.DockerClient, httpEngineAdapter)
	s.backupManager = backup.NewBackupManager(cfg.DockerClient, httpEngineAdapter, "")
	s.backupManager.SetVolumeOperations(volumeOperations)
	s.volumeOperations = volumeOperations

	// Core Services
	s.projectService = projectsmod.NewService(s.projectRepo, s.orgRepo)
	s.projectAppService = projectsmod.NewProjectAppService(s.projectAppRepo)
	volumeRepo := appsmod.NewServiceVolumeRepo(cfg.DB)
	s.appService = appsmod.NewService(s.appRepo, s.serviceVarRepo, volumeRepo)
	s.databaseService = databases.NewService(s.dbRepo, databases.NewDockerDeployer(cfg.DockerClient, s.dbRepo, volumeOperations))
	s.notifSettingsService = notifsmod.NewNotificationSettingsService(s.notifRepo, s.notifDeliveryRepo, s.deviceTokenRepo)
	s.aiSettingsService = aimod.NewService(s.aiRepo)
	s.aiAnalysisService = aimod.NewAnalysisService(s.aiRepo)

	s.serviceLinker = appsmod.NewServiceLinker(&adapters.AppDatabaseAdapter{Repo: s.dbRepo})
	s.serviceLinker.SetApplications(s.appRepo, s.serviceVarRepo)
	s.serviceLinker.SetClusterData(&adapters.AppClusterPlansAdapter{Repo: s.clusterDataRepo})

	var err error
	s.mailerService, err = notifications.NewMailerService(s.notifSettingsService)
	if err != nil {
		return nil, fmt.Errorf("mailer service: %w", err)
	}
	s.userService = users.NewUserService(s.userRepo)
	s.settingsService = settingsmodule.NewSettingsService(s.settingsRepo)
	s.auditService = audit.NewAuditService(audit.NewAuditLogRepo(cfg.DB), audit.NewSQLAuditResourceSearcher(cfg.DB))
	s.orgService = organizations.NewOrganizationService(s.orgRepo)

	s.authMod, err = auth.NewModule(cfg.DB, s.userRepo, s.userService, s.settingsRepo, s.auditService, s.mailerService)
	if err != nil {
		return nil, err
	}
	s.projectSettingsService = projectsmod.NewProjectSettingsService(s.projectSettingsRepo)

	s.dispatcherService = notifications.NewDispatcherService(
		&adapters.NotifSettingsAdapter{Repo: s.settingsRepo},
		s.notifRepo,
		&adapters.NotifUserLookupAdapter{Repo: s.userRepo},
		s.notifDeliveryRepo,
		s.mailerService,
		&adapters.NotifDeviceTokenAdapter{Repo: s.deviceTokenRepo},
	)
	s.deploymentListeners = notifications.NewDeploymentListeners(s.dispatcherService, &adapters.NotifAppRepoAdapter{Repo: s.appRepo})

	s.sshManager = ssh.NewSSHManager(&adapters.SSHServerRepoAdapter{Repo: s.serverRepo})
	s.scheduledTaskService = tasks.NewScheduledTaskService(s.scheduledTaskRepo, s.cronManager)
	s.canvasService = canvas.NewCanvasService(s.canvasRepo)
	s.canvasService.SetRuntime(&adapters.CanvasRuntimeAdapter{Rt: observability.NewCanvasRuntime(cfg.DockerClient)})

	gitProviderRepo := gitmod.NewGitProviderRepository(cfg.DB, cfg.Vault)
	s.gitService = gitmod.NewProviderService(gitProviderRepo)
	statsMonitor := observability.NewStatsMonitor(cfg.DockerClient)
	s.statsMonitor = statsMonitor
	s.deploymentService = deploymentsmod.NewDeploymentService(
		s.deployRepo,
		&adapters.AppLookupAdapter{AppService: s.appService},
		&adapters.ProjectLookupAdapter{ProjectService: s.projectService},
		s.gitService,
		&adapters.DeployEngineAdapter{Deployer: cfg.Deployer, AppService: s.appService},
		statsMonitor,
	)
	s.deploymentService.SetLocalDocker(cfg.DockerClient)
	s.deploymentService.SetRuntime(s.runtimeRepo)

	serverlessRepo := appsmod.NewServerlessRepository(cfg.DB)
	autoscalingRepo := appsmod.NewWorkerAutoscalingRepo(cfg.DB)
	s.autoscalingService = appsmod.NewAutoscalingService(appsmod.NewAutoscalingRepository(cfg.DB), serverlessRepo, s.appRepo)
	s.autoscaler = deploy.NewAutoscalerWorker(s.appRepo, statsMonitor, &adapters.AutoscalerDeploymentAdapter{Ds: s.deploymentService}, autoscalingRepo)

	s.operationService = operations.NewOperationService(operations.NewOperationRepo(cfg.DB, cfg.Vault))

	backupModRepo := backupmodule.NewRepository(cfg.DB, cfg.Vault)
	backupModRepo.SetSnapshotDumper(backupmodule.NewSystemSnapshotDumper(cfg.DataDir))
	s.backupService = backupmodule.NewService(backupModRepo, backupModRepo, backupmodule.NewManager(s.backupManager))
	s.backupService.SetBatches(backupModRepo)
	s.backupService.SetSFTP(backupmodule.NewSFTPRepository(cfg.DB))
	s.backupService.SetOperations(s.operationService)

	s.prPreviewService = gitmod.NewPRPreviewService(s.prPreviewRepo, &adapters.GitServiceLookupAdapter{AppRepo: s.appRepo})
	dnsProviderService := dns.NewProviderService(&adapters.DNSSettingsAdapter{Repo: s.settingsRepo})
	s.dnsService = dns.NewService(s.dnsRepo, dnsProviderService)
	s.gitAppsService = gitmod.NewGitAppsService(s.gitAppRepo)
	webhookDeliveryRepo := gitmod.NewWebhookDeliveryRepository(cfg.DB)
	s.webhookDeliveryService = gitmod.NewWebhookDeliveryService(webhookDeliveryRepo)
	s.webhookDispatcher = gitmod.NewWebhookDispatcher(s.webhookDeliveryService, s.deploymentService, &adapters.GitServiceLookupAdapter{AppRepo: s.appRepo}, s.gitAppsService)

	// Clusters Services
	s.clusterRunner = kubernetes.NewClusterRunner(s.serverRepo)
	s.clusterService = modclusters.NewClusterService(
		s.clusterRepo,
		s.projectRepo,
		s.serverRepo,
		s.clusterRunner,
		s.operationService,
		volumeOperations,
		&http.Client{Timeout: time.Minute},
	)
	s.clusterDataService = modclusters.NewClusterDataService(
		s.clusterDataRepo,
		s.clusterRepo,
		s.environmentRepo,
		s.backupRepo,
		s.clusterRunner,
		s.operationService,
		volumeOperations,
		&http.Client{Timeout: time.Minute},
	)

	// Monitoring Services
	s.metricsService = modmonitoring.NewMetricsService()
	s.logService = modmonitoring.NewLogService()
	s.exampleService = modmonitoring.NewExampleService()
	s.analyticsService = modmonitoring.NewService(
		s.trafficRepo,
		&adapters.AnalyticsDeployRepoAdapter{Repo: s.deployRepo},
		&adapters.AnalyticsProjectRepoAdapter{Repo: s.projectRepo},
		&adapters.AnalyticsAppRepoAdapter{Repo: s.appRepo},
		&adapters.AnalyticsDomainListerAdapter{Repo: s.domainRepo},
		&adapters.AnalyticsUserRepoAdapter{Repo: s.userRepo},
		&adapters.AnalyticsOrgRepoAdapter{Repo: s.orgRepo},
		s.attentionRepo,
		modmonitoring.NewUsageReader(&adapters.AnalyticsAppRepoAdapter{Repo: s.appRepo}, s.metricsService, s.statsMonitor),
	)
	s.attentionService = modmonitoring.NewAttentionManager(
		s.attentionRepo,
		&adapters.AnalyticsDeployRepoAdapter{Repo: s.deployRepo},
		&adapters.AnalyticsAppRepoAdapter{Repo: s.appRepo},
		&adapters.AnalyticsProjectRepoAdapter{Repo: s.projectRepo},
		&adapters.AnalyticsUserRepoAdapter{Repo: s.userRepo},
		&adapters.AnalyticsOrgRepoAdapter{Repo: s.orgRepo},
		s.backupService,
		nil,
		&adapters.AttentionManagedAdapter{Repo: s.managedRepo},
		&adapters.AttentionServerAdapter{Repo: s.serverRepo},
		&adapters.StatsHealthAdapter{Monitor: s.statsMonitor},
		&adapters.AttentionDeployerAdapter{DeploymentService: s.deploymentService, AppService: s.appService},
	)

	s.authGuard = middleware.NewAuthGuard(s.authMod.TokenService, s.settingsService, s.projectSettingsService, s.orgRepo, s.projectRepo, s.userRepo)
	s.authGuard.PersonalTokens = s.userRepo
	s.authGuard.PersonalResources = users.NewPersonalTokenResources(cfg.DB)

	s.authRateLimiter = middleware.NewRateLimiter(10, time.Minute)
	s.otpRateLimiter = middleware.NewRateLimiter(5, time.Minute)
	s.aiRateLimiter = middleware.NewRateLimiter(5, time.Minute)

	s.bridge = mcp.NewBridge(s.projectService, s.databaseService)

	return s, nil
}
