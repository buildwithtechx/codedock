package http

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/server"

	"codedock/internal/engine/deploy"
	"codedock/internal/handlers/projects"
	"codedock/internal/handlers/system"
	"codedock/internal/http/mcp"
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/chatapps"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	mailmodule "codedock/internal/modules/mail"
	"codedock/internal/modules/organizations"
	"codedock/internal/modules/settings"
	"codedock/internal/modules/users"
	"codedock/internal/platform/networking"
	core "codedock/internal/platform/notifications"
	cron "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
	"codedock/internal/services/migrations"
	projectservices "codedock/internal/services/projects"
)

type Server struct {
	clusterDataHandler         *system.ClusterDataHandler
	clusterHandler             *system.ClusterHandler
	operationHandler           *system.OperationHandler
	router                     *echo.Echo
	mcpBridge                  *mcp.Bridge
	authRateLimiter            *middleware.RateLimiter
	otpRateLimiter             *middleware.RateLimiter
	aiRateLimiter              *middleware.RateLimiter
	deployer                   *deploy.Deployer
	traefikManager             *networking.TraefikManager
	dockerClient               *client.Client
	sshManager                 *ssh.SSHManager
	tokenService               *auth.TokenService
	authGuard                  *middleware.AuthGuard
	cronManager                *cron.CronManager
	serviceLinker              *projectservices.ServiceLinker
	dispatcherService          *core.DispatcherService
	projectService             *projectservices.ProjectService
	appService                 *projectservices.AppService
	appServiceHandler          *projects.AppHandler
	autoscalingHandler         *projects.AutoscalingHandler
	dbHandler                  *databases.DatabaseHandler
	scheduledTaskHandler       *system.ScheduledTaskHandler
	canvasHandler              *projects.CanvasHandler
	terminalHandler            *deploymentsmod.TerminalHandler
	deploymentHandler          *deploymentsmod.DeploymentHandler
	deploymentLifecycleHandler *deploymentsmod.LifecycleHandler
	serviceVarHandler          *projects.ServiceVarHandler
	projectSettingsHandler     *projects.ProjectSettingsHandler
	backupHandler              *backupmodule.Handler
	settingsHandler            *settings.SettingsHandler
	appearanceHandler          *users.AppearanceHandler
	notifSettingsHandler       *system.NotificationSettingsHandler
	aiSettingsHandler          *system.AISettingsHandler
	updaterHandler             *system.UpdaterHandler
	userHandler                *users.UserHandler
	authHandler                *auth.AuthHandler
	oauthHandler               *auth.OAuthHandler
	deviceFlowHandler          *auth.DeviceFlowHandler
	passkeyHandler             *auth.PasskeyHandler
	gitHandler                 *gitmod.ProvidersHandler
	gitInspectionHandler       *gitmod.InspectionHandler
	webhookHandler             *gitmod.WebhookHandler
	projectHandler             *projects.ProjectHandler
	projectAppHandler          *projects.ProjectAppHandler
	orgHandler                 *organizations.OrganizationHandler
	environmentHandler         *projects.EnvironmentHandler
	domainHandler              *projects.DomainHandler
	projectEnvHandler          *projects.ProjectEnvHandler
	notificationHandler        *system.NotificationHandler
	gitAppsHandler             *gitmod.AppsHandler
	serverlessHandler          *projects.ServerlessHandler
	systemHandler              *system.SystemHandler
	composeHandler             *projects.ComposeHandler
	composeStackHandler        *projects.ComposeStackHandler
	chatAppHandler             *chatapps.Handler
	oneClickHandler            *projects.OneClickHandler
	archiveHandler             *deploymentsmod.ArchiveHandler
	migrationHandler           *system.MigrationHandler
	onboardingHandler          *users.OnboardingHandler
	dnsHandler                 *dns.Handler
	metricsHandler             *system.MetricsHandler
	logHandler                 *system.LogHandler
	auditLogHandler            *audit.AuditLogHandler
	exampleHandler             *system.ExampleHandler
	serverHandler              *system.ServerHandler
	managedHandler             *system.ManagedHandler
	migrationLifecycleHandler  *system.MigrationLifecycleHandler
	migrationService           *migrations.Service
	analyticsHandler           *system.AnalyticsHandler
	attentionHandler           *system.AttentionHandler
	registryHandler            *deploymentsmod.RegistryHandler
	billingHandler             *system.BillingHandler
	serverMetricsWSHandler     *system.ServerMetricsWSHandler
	serviceLogsWSHandler       *system.ServiceLogsWSHandler
	takeoverHandler            *system.TakeoverHandler
	routeRuleHandler           *projects.RouteRuleHandler
	mailHandler                *mailmodule.Handler
	ingressHandler             *system.IngressHandler
	serverNetworksHandler      *system.ServerNetworksHandler
	actionsHandler             *system.ActionsHandler
	checksHandler              *projects.ChecksHandler
	supportHandler             *system.SupportHandler
	routeProtectionHandler     *environments.RouteProtectionHandler
	edgeAuthHandler            *environments.EdgeAuthHandler
	securityHandler            *system.SecurityHandler
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
	s.mcpBridge.SSEServer().SSEHandler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (s *Server) HandleMCPMessage(c echo.Context) error {
	s.mcpBridge.SSEServer().MessageHandler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (s *Server) RequireProjectRole(permission models.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := middleware.GetUserClaimsFromContext(c.Request().Context())
			if user == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
			}
			project := c.Param("id")
			if selected, ok := c.Get("project_id").(string); ok && selected != "" && selected != project {
				return echo.NewHTTPError(http.StatusForbidden, "token does not have access to this project")
			}
			if !s.projectService.HasPermission(c.Request().Context(), project, user.UserID, models.UserRole(user.Role), permission) {
				return echo.NewHTTPError(http.StatusForbidden, "project permission required")
			}
			return next(c)
		}
	}
}
