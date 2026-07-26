package http

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/server"

	"codedock.run/codedock/internal/core"
	"codedock.run/codedock/internal/engine"
	"codedock.run/codedock/internal/engine/cron"
	"codedock.run/codedock/internal/engine/networking"
	"codedock.run/codedock/internal/handlers/auth"
	"codedock.run/codedock/internal/handlers/backups"
	"codedock.run/codedock/internal/handlers/databases"
	"codedock.run/codedock/internal/handlers/deployments"
	"codedock.run/codedock/internal/handlers/projects"
	"codedock.run/codedock/internal/handlers/system"
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	authservices "codedock.run/codedock/internal/services/auth"
	projectservices "codedock.run/codedock/internal/services/projects"
)

type Server struct {
	router                 *echo.Echo
	mcpBridge              *Bridge
	authRateLimiter        *middleware.RateLimiter
	otpRateLimiter         *middleware.RateLimiter
	aiRateLimiter          *middleware.RateLimiter
	deployer               *engine.Deployer
	traefikManager         *networking.TraefikManager
	dockerClient           *client.Client
	tokenService           *authservices.TokenService
	authGuard              *middleware.AuthGuard
	cronManager            *cron.CronManager
	serviceLinker          *projectservices.ServiceLinker
	dispatcherService      *core.DispatcherService
	projectService         *projectservices.ProjectService
	appService             *projectservices.AppService
	appServiceHandler      *projects.AppHandler
	dbHandler              *databases.DatabaseHandler
	scheduledTaskHandler   *system.ScheduledTaskHandler
	canvasHandler          *projects.CanvasHandler
	terminalHandler        *deployments.TerminalHandler
	deploymentHandler      *deployments.DeploymentHandler
	serviceVarHandler      *projects.ServiceVarHandler
	projectSettingsHandler *projects.ProjectSettingsHandler
	backupHandler          *backups.BackupHandler
	settingsHandler        *auth.SettingsHandler
	notifSettingsHandler   *system.NotificationSettingsHandler
	aiSettingsHandler      *system.AISettingsHandler
	updaterHandler         *system.UpdaterHandler
	userHandler            *auth.UserHandler
	authHandler            *auth.AuthHandler
	oauthHandler           *auth.OAuthHandler
	gitHandler             *deployments.GitHandler
	webhookHandler         *deployments.WebhookHandler
	projectHandler         *projects.ProjectHandler
	orgHandler             *auth.OrganizationHandler
	environmentHandler     *projects.EnvironmentHandler
	domainHandler          *projects.DomainHandler
	projectEnvHandler      *projects.ProjectEnvHandler
	notificationHandler    *system.NotificationHandler
	gitAppsHandler         *deployments.GitAppsHandler
	serverlessHandler      *projects.ServerlessHandler
	systemHandler          *system.SystemHandler
	composeHandler         *projects.ComposeHandler
	oneClickHandler        *projects.OneClickHandler
	archiveHandler         *deployments.ArchiveHandler
	migrationHandler       *system.MigrationHandler
	onboardingHandler      *auth.OnboardingHandler
	dnsHandler             *system.DNSHandler
	metricsHandler         *system.MetricsHandler
	logHandler             *system.LogHandler
	auditLogHandler        *auth.AuditLogHandler
	exampleHandler         *system.ExampleHandler
	serverHandler          *system.ServerHandler
	workerWSHandler        *system.WorkerWSHandler
	registryHandler        *deployments.RegistryHandler
	billingHandler         *system.BillingHandler
	serverMetricsWSHandler *system.ServerMetricsWSHandler
	serviceLogsWSHandler   *system.ServiceLogsWSHandler
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
