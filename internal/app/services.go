package app

import (
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
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/mcp"
	modmonitoring "codedock/internal/modules/monitoring"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/projects/canvas"
	"codedock/internal/modules/servers"
	"codedock/internal/modules/servers/managed"
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
	managedRepo         managed.ManagedRepository
	runtimeRepo         *appsmod.RuntimeRepo
	routeRuleRepo       environments.RouteRuleRepository
	routeProtectionRepo environments.RouteProtectionRepository
	deviceTokenRepo     notifsmod.DeviceTokenRepository
	notifDeliveryRepo   notifsmod.NotificationDeliveryRepository
	clusterRepo         *modclusters.ClusterRepo
	clusterDataRepo     *modclusters.ClusterDataRepo
	trafficRepo         modmonitoring.TrafficRepository
	attentionRepo       modmonitoring.AttentionRepository

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
	clusterRunner          *kubernetes.ClusterRunner
	clusterService         *modclusters.ClusterService
	clusterDataService     *modclusters.ClusterDataService
	metricsService         *modmonitoring.MetricsService
	analyticsService       *modmonitoring.Service
	attentionService       *modmonitoring.AttentionManager
	logService             *modmonitoring.LogService
	exampleService         *modmonitoring.ExampleService

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
	s.managedRepo = managed.NewManagedRepository(cfg.DB, cfg.Vault)
	s.runtimeRepo = appsmod.NewRuntimeRepo(cfg.DB, cfg.Vault)
	s.routeRuleRepo = environments.NewRouteRuleRepository(cfg.DB)
	s.routeProtectionRepo = environments.NewRouteProtectionRepository(cfg.DB)
	s.deviceTokenRepo = notifsmod.NewDeviceTokenRepo(cfg.DB)
	s.notifDeliveryRepo = notifsmod.NewNotificationDeliveryRepo(cfg.DB)
	s.clusterRepo = modclusters.NewClusterRepo(cfg.DB, cfg.Vault)
	s.clusterDataRepo = modclusters.NewClusterDataRepo(cfg.DB, cfg.Vault)
	s.trafficRepo = modmonitoring.NewTrafficRepository(cfg.DB)
	s.attentionRepo = modmonitoring.NewAttentionRepository(cfg.DB)

	// Canvas Repo
	s.canvasRepo = canvas.NewCanvasRepo(cfg.DB, &adapters.CanvasEnvAdapter{Repo: s.environmentRepo})
	s.canvasRepo.SetVault(cfg.Vault)
	s.canvasRepo.SetClusterData(&adapters.CanvasClusterDataAdapter{Repo: s.clusterDataRepo})

	return initDomainServices(s, cfg)
}
