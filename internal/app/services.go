package app

import (
	"fmt"
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
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/mcp"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/projects/canvas"
	"codedock/internal/modules/servers"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/users"
	"codedock/internal/platform/notifications"
	"codedock/internal/platform/observability"
	cron "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
	backup "codedock/internal/platform/storage"
)

type appServices struct {
	// Repositories
	environmentRepo     *environments.EnvironmentRepo
	projectRepo         *projectsmod.ProjectRepo
	projectAppRepo      *projectsmod.ProjectAppRepo
	appRepo             *appsmod.AppServiceRepo
	serviceVarRepo      *appsmod.ServiceVarRepo
	dbRepo              databases.Repository
	settingsRepo        settingsmodule.Repository
	notifRepo           notifsmod.NotificationSettingsRepository
	aiRepo              aimod.Repository
	envVarRepo          *projectsmod.EnvRepo
	scheduledTaskRepo   tasks.ScheduledTaskRepository
	backupRepo          backupmodule.Repository
	projectSettingsRepo *projectsmod.ProjectSettingsRepo
	userRepo            users.UserRepository
	canvasRepo          *canvas.CanvasRepo
	deployRepo          deploymentsmod.DeploymentRepository
	prPreviewRepo       gitmod.PRPreviewRepository
	domainRepo          domainsmod.DomainRepository
	gitAppRepo          gitmod.GitAppRepository
	dnsRepo             dns.Repository
	orgRepo             organizations.OrganizationRepository
	serverRepo          servers.Repository
	runtimeRepo         *appsmod.RuntimeRepo
	routeRuleRepo       environments.RouteRuleRepository
	routeProtectionRepo environments.RouteProtectionRepository
	deviceTokenRepo     notifsmod.DeviceTokenRepository
	notifDeliveryRepo   notifsmod.NotificationDeliveryRepository

	// Services & Managers
	cronManager            *cron.CronManager
	backupManager          *backup.BackupManager
	sshManager             *ssh.SSHManager
	projectService         *projectsmod.Service
	projectAppService      *projectsmod.ProjectAppService
	appService             *appsmod.Service
	databaseService        *databases.Service
	notifSettingsService   *notifsmod.NotificationSettingsService
	aiSettingsService      *aimod.Service
	aiAnalysisService      *aimod.AnalysisService
	userService            *users.UserService
	settingsService        *settingsmodule.SettingsService
	auditService           *audit.AuditService
	orgService             *organizations.OrganizationService
	authMod                *auth.Module
	projectSettingsService *projectsmod.ProjectSettingsService
	dispatcherService      *notifications.DispatcherService
	deploymentListeners    *notifications.DeploymentListeners
	scheduledTaskService   *tasks.ScheduledTaskService
	canvasService          *canvas.CanvasService
	gitService             *gitmod.ProviderService
	deploymentService      *deploymentsmod.DeploymentService
	autoscalingService     *appsmod.AutoscalingService
	autoscaler             *deploy.AutoscalerWorker
	operationService       *operations.OperationService
	backupService          *backupmodule.Service
	mailerService          *notifications.MailerService
	prPreviewService       *gitmod.PRPreviewService
	dnsService             *dns.Service
	gitAppsService         *gitmod.GitAppsService
	webhookDeliveryService *gitmod.WebhookDeliveryService
	webhookDispatcher      *gitmod.WebhookDispatcher
	serviceLinker          *appsmod.ServiceLinker
	volumeOperations       *deploy.VolumeGate
	statsMonitor           *observability.StatsMonitor

	// Auth & Rate Limiters
	authGuard       *middleware.AuthGuard
	authRateLimiter *middleware.RateLimiter
	otpRateLimiter  *middleware.RateLimiter
	aiRateLimiter   *middleware.RateLimiter
	bridge          *mcp.Bridge
}

func initServices(cfg Config) (*appServices, error) {
	s := &appServices{}

	// Core Repositories
	s.environmentRepo = environments.NewEnvironmentRepo(cfg.DB)
	s.projectRepo = projectsmod.NewProjectRepo(cfg.DB)
	s.projectAppRepo = projectsmod.NewProjectAppRepo(cfg.DB)
	s.appRepo = appsmod.NewAppServiceRepo(cfg.DB)
	s.serviceVarRepo = appsmod.NewServiceVarRepo(cfg.DB)
	s.dbRepo = databases.NewRepository(cfg.DB, cfg.Vault)
	s.settingsRepo = settingsmodule.NewSettingsRepo(cfg.DB, cfg.Vault)
	s.notifRepo = notifsmod.NewNotificationSettingsRepo(cfg.DB)
	s.aiRepo = aimod.NewRepository(cfg.DB)
	s.envVarRepo = projectsmod.NewEnvRepo(cfg.DB, cfg.Vault)
	s.scheduledTaskRepo = tasks.NewScheduledTaskRepo(cfg.DB)
	s.backupRepo = backupmodule.NewRepository(cfg.DB, cfg.Vault)
	s.backupRepo.SetSnapshotDumper(backupmodule.NewSystemSnapshotDumper(cfg.DataDir))
	s.projectSettingsRepo = projectsmod.NewProjectSettingsRepo(cfg.DB)
	s.userRepo = users.NewUserRepo(cfg.DB)
	s.deployRepo = deploymentsmod.NewDeploymentRepository(cfg.DB)
	s.prPreviewRepo = gitmod.NewPRPreviewRepository(cfg.DB)
	s.domainRepo = domainsmod.NewDomainRepo(cfg.DB)
	s.gitAppRepo = gitmod.NewGitAppRepository(cfg.DB, cfg.Vault)
	s.dnsRepo = dns.NewRepository(cfg.DB)
	s.orgRepo = organizations.NewOrganizationRepository(cfg.DB)
	s.serverRepo = servers.NewRepository(cfg.DB, cfg.Vault)
	s.runtimeRepo = appsmod.NewRuntimeRepo(cfg.DB, cfg.Vault)
	s.routeRuleRepo = environments.NewRouteRuleRepository(cfg.DB)
	s.routeProtectionRepo = environments.NewRouteProtectionRepository(cfg.DB)
	s.deviceTokenRepo = notifsmod.NewDeviceTokenRepo(cfg.DB)
	s.notifDeliveryRepo = notifsmod.NewNotificationDeliveryRepo(cfg.DB)

	// Canvas Repo
	s.canvasRepo = canvas.NewCanvasRepo(cfg.DB, &adapters.CanvasEnvAdapter{Repo: s.environmentRepo})
	s.canvasRepo.SetVault(cfg.Vault)
	s.canvasRepo.SetClusterData(&adapters.CanvasClusterDataAdapter{Repo: modclusters.NewClusterDataRepo(cfg.DB, cfg.Vault)})

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

	// Domain Services
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
	s.serviceLinker.SetClusterData(&adapters.AppClusterPlansAdapter{Repo: modclusters.NewClusterDataRepo(cfg.DB, cfg.Vault)})

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

	s.authGuard = middleware.NewAuthGuard(s.authMod.TokenService, s.settingsService, s.projectSettingsService, s.orgRepo, s.projectRepo, s.userRepo)
	s.authGuard.PersonalTokens = s.userRepo
	s.authGuard.PersonalResources = users.NewPersonalTokenResources(cfg.DB)

	s.authRateLimiter = middleware.NewRateLimiter(10, time.Minute)
	s.otpRateLimiter = middleware.NewRateLimiter(5, time.Minute)
	s.aiRateLimiter = middleware.NewRateLimiter(5, time.Minute)

	s.bridge = mcp.NewBridge(s.projectService, s.databaseService)

	return s, nil
}
