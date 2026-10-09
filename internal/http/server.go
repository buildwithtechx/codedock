package http

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/server"

	"codedock/internal/core"
	"codedock/internal/engine/cron"
	"codedock/internal/engine/deploy"
	"codedock/internal/engine/networking"
	"codedock/internal/engine/ssh"
	"codedock/internal/handlers/auth"
	"codedock/internal/handlers/backups"
	"codedock/internal/handlers/databases"
	"codedock/internal/handlers/deployments"
	"codedock/internal/handlers/projects"
	"codedock/internal/handlers/system"
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	"codedock/internal/services/migrations"
	projectservices "codedock/internal/services/projects"
)

type Server struct {
	clusterDataHandler        *system.ClusterDataHandler
	clusterHandler            *system.ClusterHandler
	operationHandler          *system.OperationHandler
	router                    *echo.Echo
	mcpBridge                 *Bridge
	authRateLimiter           *middleware.RateLimiter
	otpRateLimiter            *middleware.RateLimiter
	aiRateLimiter             *middleware.RateLimiter
	deployer                  *deploy.Deployer
	traefikManager            *networking.TraefikManager
	dockerClient              *client.Client
	sshManager                *ssh.SSHManager
	tokenService              *authservices.TokenService
	authGuard                 *middleware.AuthGuard
	cronManager               *cron.CronManager
	serviceLinker             *projectservices.ServiceLinker
	dispatcherService         *core.DispatcherService
	projectService            *projectservices.ProjectService
	appService                *projectservices.AppService
	appServiceHandler         *projects.AppHandler
	autoscalingHandler        *projects.AutoscalingHandler
	dbHandler                 *databases.DatabaseHandler
	scheduledTaskHandler      *system.ScheduledTaskHandler
	canvasHandler             *projects.CanvasHandler
	terminalHandler           *deployments.TerminalHandler
	deploymentHandler         *deployments.DeploymentHandler
	serviceVarHandler         *projects.ServiceVarHandler
	projectSettingsHandler    *projects.ProjectSettingsHandler
	backupHandler             *backups.BackupHandler
	settingsHandler           *auth.SettingsHandler
	notifSettingsHandler      *system.NotificationSettingsHandler
	aiSettingsHandler         *system.AISettingsHandler
	updaterHandler            *system.UpdaterHandler
	userHandler               *auth.UserHandler
	authHandler               *auth.AuthHandler
	oauthHandler              *auth.OAuthHandler
	gitHandler                *deployments.GitHandler
	webhookHandler            *deployments.WebhookHandler
	projectHandler            *projects.ProjectHandler
	projectAppHandler         *projects.ProjectAppHandler
	orgHandler                *auth.OrganizationHandler
	environmentHandler        *projects.EnvironmentHandler
	domainHandler             *projects.DomainHandler
	projectEnvHandler         *projects.ProjectEnvHandler
	notificationHandler       *system.NotificationHandler
	gitAppsHandler            *deployments.GitAppsHandler
	serverlessHandler         *projects.ServerlessHandler
	systemHandler             *system.SystemHandler
	composeHandler            *projects.ComposeHandler
	composeStackHandler       *projects.ComposeStackHandler
	oneClickHandler           *projects.OneClickHandler
	archiveHandler            *deployments.ArchiveHandler
	migrationHandler          *system.MigrationHandler
	onboardingHandler         *auth.OnboardingHandler
	dnsHandler                *system.DNSHandler
	metricsHandler            *system.MetricsHandler
	logHandler                *system.LogHandler
	auditLogHandler           *auth.AuditLogHandler
	exampleHandler            *system.ExampleHandler
	serverHandler             *system.ServerHandler
	managedHandler            *system.ManagedHandler
	migrationLifecycleHandler *system.MigrationLifecycleHandler
	migrationService          *migrations.Service
	analyticsHandler          *system.AnalyticsHandler
	attentionHandler          *system.AttentionHandler
	registryHandler           *deployments.RegistryHandler
	billingHandler            *system.BillingHandler
	serverMetricsWSHandler    *system.ServerMetricsWSHandler
	serviceLogsWSHandler      *system.ServiceLogsWSHandler
	takeoverHandler           *system.TakeoverHandler
	routeRuleHandler          *projects.RouteRuleHandler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func GetUserClaimsFromContext(ctx context.Context) *models.UserClaims {
	return middleware.GetUserClaimsFromContext(ctx)
}

func (s *Server) StartMCPStdio() error {
	mcpServer := s.mcpBridge.MCPServer()
	return server.ServeStdio(mcpServer)
}

func (s *Server) HandleMCPSSE(c echo.Context) error {
	mcpServer := s.mcpBridge.MCPServer()
	sseServer := server.NewSSEServer(mcpServer)
	sseServer.SSEHandler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (s *Server) HandleMCPMessage(c echo.Context) error {
	mcpServer := s.mcpBridge.MCPServer()
	sseServer := server.NewSSEServer(mcpServer)
	sseServer.MessageHandler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}
