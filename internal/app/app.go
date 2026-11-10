package app

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/server"

	"codedock/internal/deploy"
	cdkhttp "codedock/internal/http"
	"codedock/internal/http/middleware"
	appsmod "codedock/internal/modules/apps"
	"codedock/internal/modules/mcp"
	"codedock/internal/modules/organizations"
	projectsmod "codedock/internal/modules/projects"
	"codedock/internal/modules/users"
	"codedock/internal/platform/networking"
	cron "codedock/internal/platform/scheduler"
	"codedock/internal/platform/ssh"
	backup "codedock/internal/platform/storage"
	"codedock/internal/platform/vault"
	"codedock/internal/utils"
)

// Config encapsulates parameters required to bootstrap the Codedock application.
type Config struct {
	DB             *sql.DB
	Vault          *vault.Vault
	Deployer       *deploy.Deployer
	TraefikManager *networking.TraefikManager
	DockerClient   *client.Client
	DataDir        string
}

// App represents the central composition container for Codedock.
type App struct {
	router         *echo.Echo
	db             *sql.DB
	vault          *vault.Vault
	dataDir        string
	deployer       *deploy.Deployer
	traefikManager *networking.TraefikManager
	dockerClient   *client.Client
	mcpBridge      *mcp.Bridge
	authGuard      *middleware.AuthGuard
	cronManager    *cron.CronManager
	backupManager  *backup.BackupManager
	sshManager     *ssh.SSHManager
	projectService *projectsmod.Service
	appService     *appsmod.Service
}

// New initializes the entire application dependency graph, executes startup
// recovery, sets up event listeners, and mounts domain HTTP routes.
func New(cfg Config) (*App, error) {
	if cfg.Deployer != nil && cfg.DataDir != "" {
		if err := recoverRollouts(cfg.Deployer, cfg.DataDir); err != nil {
			return nil, err
		}
	}

	services, err := initServices(cfg)
	if err != nil {
		return nil, err
	}

	if err := runRecovery(context.Background(), services); err != nil {
		return nil, err
	}

	a := &App{
		router:         cdkhttp.NewRouter(),
		db:             cfg.DB,
		vault:          cfg.Vault,
		dataDir:        cfg.DataDir,
		deployer:       cfg.Deployer,
		traefikManager: cfg.TraefikManager,
		dockerClient:   cfg.DockerClient,
		mcpBridge:      services.bridge,
		authGuard:      services.authGuard,
		cronManager:    services.cronManager,
		backupManager:  services.backupManager,
		sshManager:     services.sshManager,
		projectService: services.projectService,
		appService:     services.appService,
	}

	configureBindings(a, services)
	accessLogPath := ""
	if cfg.TraefikManager != nil {
		accessLogPath = cfg.TraefikManager.AccessLogPath()
	}
	configureListeners(services, cfg.Deployer, accessLogPath)
	a.wireRoutes(services)

	return a, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.router.ServeHTTP(w, r)
}

func (a *App) Handler() http.Handler {
	return a.router
}

func (a *App) Router() *echo.Echo {
	return a.router
}

func (a *App) Bridge() *mcp.Bridge {
	return a.mcpBridge
}

func (a *App) StartMCPStdio() error {
	mcpServer := a.mcpBridge.MCPServer()
	return server.ServeStdio(mcpServer)
}

func (a *App) HandleMCPSSE(c echo.Context) error {
	a.mcpBridge.SSEServer().SSEHandler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (a *App) HandleMCPMessage(c echo.Context) error {
	a.mcpBridge.SSEServer().MessageHandler().ServeHTTP(c.Response().Writer, c.Request())
	return nil
}

func (a *App) RequireProjectRole(permission organizations.MemberPermission) echo.MiddlewareFunc {
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
			if !a.projectService.HasPermission(c.Request().Context(), project, user.UserID, users.UserRole(user.Role), permission) {
				return echo.NewHTTPError(http.StatusForbidden, "project permission required")
			}
			return next(c)
		}
	}
}

func (a *App) RequireServiceRole(minPermission organizations.MemberPermission) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			serviceID := c.Param("serviceId")
			if serviceID == "" {
				serviceID = c.Param("id")
			}
			if serviceID == "" {
				return next(c)
			}

			userClaims := middleware.GetUserClaimsFromContext(c.Request().Context())
			if userClaims == nil {
				return utils.Error(c, 401, "unauthorized")
			}
			if userClaims.Role == users.UserRoleAdmin || userClaims.Role == users.UserRoleOwner {
				return next(c)
			}

			svc, err := a.appService.GetAppService(c.Request().Context(), serviceID)
			if err != nil || svc == nil {
				return utils.Error(c, 404, "service not found")
			}

			if userClaims.Role == "api" {
				if c.Get("project_id") != svc.ProjectID {
					return utils.Error(c, 403, "api token not authorized for this project")
				}
				if minPermission != "" && minPermission != organizations.MemberPermissionMember {
					return utils.Error(c, 403, "api tokens cannot perform admin/owner actions")
				}
				return next(c)
			}

			if !a.projectService.HasPermission(c.Request().Context(), svc.ProjectID, userClaims.UserID, userClaims.Role, minPermission) {
				return utils.Error(c, 403, "insufficient project permissions")
			}
			return next(c)
		}
	}
}
