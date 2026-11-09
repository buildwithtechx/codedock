package http

import (
	"database/sql"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"

	"codedock/internal/engine/deploy"
	"codedock/internal/http/mcp"
	"codedock/internal/http/middleware"
	aimod "codedock/internal/modules/ai"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/canvas"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	"codedock/internal/modules/monitoring"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/servers"
	settingsmodule "codedock/internal/modules/settings"
	"codedock/internal/modules/tasks"
	"codedock/internal/modules/users"
	"codedock/internal/platform/networking"
	"codedock/internal/platform/notifications"
	"codedock/internal/platform/observability"
	cron "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
	backup "codedock/internal/platform/storage"
	"codedock/internal/platform/vault"
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
	serviceLinker          *appsmod.ServiceLinker
	dispatcherService      *notifications.DispatcherService
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
	scheduledTaskService   *tasks.ScheduledTaskService
	canvasService          *canvas.CanvasService
	gitService             *gitmod.ProviderService
	deploymentService      *deploymentsmod.DeploymentService
	autoscalingService     *appsmod.AutoscalingService
	operationService       *monitoring.OperationService
	backupService          *backupmodule.Service
	mailerService          *notifications.MailerService
	environmentRepo        *environments.EnvironmentRepo
	projectRepo            *projectsmod.ProjectRepo
	appRepo                *appsmod.AppServiceRepo
	dbRepo                 databases.Repository
	settingsRepo           settingsmodule.Repository
	envVarRepo             *projectsmod.EnvRepo
	serverlessRepository   appsmod.ServerlessRepository
	userRepo               users.UserRepository
	userModRepo            users.UserRepository
	deployRepo             deploymentsmod.DeploymentRepository
	prPreviewRepository    gitmod.PRPreviewRepository
	domainRepo             domainsmod.DomainRepository
	gitAppRepo             gitmod.GitAppRepository
	dnsRepo                dns.Repository
	orgRepo                organizations.OrganizationRepository
	webhookDeliveryRepo    gitmod.WebhookDeliveryRepository
	serverRepo             servers.Repository
	runtimeRepository      *appsmod.RuntimeRepo
}
