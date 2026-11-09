package http

import (
	"context"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/server"

	"codedock/internal/engine/deploy"
	"codedock/internal/http/mcp"
	"codedock/internal/http/middleware"
	actions "codedock/internal/modules/actions"
	aimod "codedock/internal/modules/ai"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/audit"
	"codedock/internal/modules/auth"
	backupmodule "codedock/internal/modules/backups"
	"codedock/internal/modules/billing"
	"codedock/internal/modules/canvas"
	"codedock/internal/modules/chatapps"
	modclusters "codedock/internal/modules/clusters"
	composemod "codedock/internal/modules/compose"
	"codedock/internal/modules/databases"
	deploymentsmod "codedock/internal/modules/deployments"
	"codedock/internal/modules/dns"
	"codedock/internal/modules/dns/ingress"
	domainsmod "codedock/internal/modules/domains"
	"codedock/internal/modules/environments"
	gitmod "codedock/internal/modules/git"
	mailmodule "codedock/internal/modules/mail"
	modmigrations "codedock/internal/modules/migrations"
	modmonitoring "codedock/internal/modules/monitoring"
	notifsmod "codedock/internal/modules/notifications"
	"codedock/internal/modules/operations"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	securitymod "codedock/internal/modules/security"
	"codedock/internal/modules/servers"
	"codedock/internal/modules/servers/managed"
	"codedock/internal/modules/servers/networks"
	"codedock/internal/modules/servers/takeover"
	"codedock/internal/modules/settings"
	"codedock/internal/modules/support"
	systemmod "codedock/internal/modules/system"
	"codedock/internal/modules/tasks"
	templatesmod "codedock/internal/modules/templates"
	"codedock/internal/modules/users"
	"codedock/internal/platform/networking"
	core "codedock/internal/platform/notifications"
	cron "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
)

type Server struct {
	clusterDataHandler         *modclusters.DataHandler
	clusterHandler             *modclusters.Handler
	operationHandler           *operations.OperationHandler
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
	serviceLinker              *appsmod.ServiceLinker
	dispatcherService          *core.DispatcherService
	projectService             *projectsmod.Service
	appService                 *appsmod.Service
	appServiceHandler          *appsmod.Handler
	appVarHandler              *appsmod.VariableHandler
	appVolumeHandler           *appsmod.VolumeHandler
	appWebhookHandler          *appsmod.WebhookHandler
	autoscalingHandler         *appsmod.AutoscalingHandler
	dbHandler                  *databases.DatabaseHandler
	scheduledTaskHandler       *tasks.ScheduledTaskHandler
	canvasHandler              *canvas.CanvasHandler
	terminalHandler            *deploymentsmod.TerminalHandler
	deploymentHandler          *deploymentsmod.DeploymentHandler
	deploymentLifecycleHandler *deploymentsmod.LifecycleHandler
	archiveHandler             *deploymentsmod.ArchiveHandler
	registryHandler            *deploymentsmod.RegistryHandler
	projectSettingsHandler     *projectsmod.ProjectSettingsHandler
	backupHandler              *backupmodule.Handler
	settingsHandler            *settings.SettingsHandler
	appearanceHandler          *users.AppearanceHandler
	notifSettingsHandler       *notifsmod.NotificationSettingsHandler
	aiSettingsHandler          *aimod.Handler
	updaterHandler             *systemmod.UpdaterHandler
	userHandler                *users.UserHandler
	authHandler                *auth.AuthHandler
	oauthHandler               *auth.OAuthHandler
	deviceFlowHandler          *auth.DeviceFlowHandler
	passkeyHandler             *auth.PasskeyHandler
	gitHandler                 *gitmod.ProvidersHandler
	gitInspectionHandler       *gitmod.InspectionHandler
	gitAppsHandler             *gitmod.AppsHandler
	webhookHandler             *gitmod.WebhookHandler
	projectHandler             *projectsmod.Handler
	projectAppHandler          *projectsmod.ProjectAppHandler
	orgHandler                 *organizations.OrganizationHandler
	environmentHandler         *environments.EnvironmentHandler
	domainHandler              *domainsmod.DomainHandler
	notificationHandler        *notifsmod.NotificationHandler
	systemHandler              *systemmod.SystemHandler
	composeStackHandler        *composemod.ComposeStackHandler
	chatAppHandler             *chatapps.Handler
	oneClickHandler            *templatesmod.OneClickHandler
	migrationHandler           *modmigrations.BundleHandler
	onboardingHandler          *users.OnboardingHandler
	dnsHandler                 *dns.Handler
	metricsHandler             *modmonitoring.MetricsHandler
	logHandler                 *modmonitoring.LogHandler
	auditLogHandler            *audit.AuditLogHandler
	exampleHandler             *modmonitoring.ExampleHandler
	serverHandler              *servers.Handler
	managedHandler             *managed.ManagedHandler
	migrationLifecycleHandler  *modmigrations.LifecycleHandler
	analyticsHandler           *modmonitoring.AnalyticsHandler
	attentionHandler           *modmonitoring.AttentionHandler
	billingHandler             *billing.BillingHandler
	serverMetricsWSHandler     *servers.ServerMetricsWSHandler
	serviceLogsWSHandler       *modmonitoring.ServiceLogsWSHandler
	takeoverHandler            *takeover.TakeoverHandler
	routeRuleHandler           *environments.RouteRuleHandler
	mailHandler                *mailmodule.Handler
	ingressHandler             *ingress.IngressHandler
	serverNetworksHandler      *networks.ServerNetworksHandler
	actionsHandler             *actions.ActionsHandler
	checksHandler              *projectsmod.ChecksHandler
	supportHandler             *support.SupportHandler
	routeProtectionHandler     *environments.RouteProtectionHandler
	edgeAuthHandler            *environments.EdgeAuthHandler
	securityHandler            *securitymod.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Handler() http.Handler {
	return s.router
}

func GetUserClaimsFromContext(ctx context.Context) *middleware.UserClaims {
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

func (s *Server) RequireProjectRole(permission organizations.MemberPermission) echo.MiddlewareFunc {
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
			if !s.projectService.HasPermission(c.Request().Context(), project, user.UserID, users.UserRole(user.Role), permission) {
				return echo.NewHTTPError(http.StatusForbidden, "project permission required")
			}
			return next(c)
		}
	}
}
