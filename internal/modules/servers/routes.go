package servers

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/modules/servers/managed"
	"codedock/internal/modules/servers/networks"
	"codedock/internal/modules/servers/takeover"
)

type RoutesConfig struct {
	Handler         *Handler
	ManagedHandler  *managed.ManagedHandler
	TakeoverHandler *takeover.TakeoverHandler
	NetworksHandler *networks.ServerNetworksHandler
	WSHandler       *ServerMetricsWSHandler
	Middleware      RouteMiddleware
}

type RouteMiddleware struct {
	ServerRead  echo.MiddlewareFunc
	ServerWrite echo.MiddlewareFunc
	Admin       echo.MiddlewareFunc
	OrgMember   echo.MiddlewareFunc
	OrgAdmin    echo.MiddlewareFunc
}

func RegisterRoutes(apiGroup, authGroup *echo.Group, cfg RoutesConfig) {
	mw := cfg.Middleware

	if cfg.Handler != nil {
		authGroup.GET("/system/monitor/stream", cfg.Handler.StreamMonitor, mw.ServerRead)
		authGroup.GET("/servers", cfg.Handler.List, mw.ServerRead)
		authGroup.GET("/servers/:id", cfg.Handler.Get, mw.ServerRead)
		authGroup.POST("/servers", cfg.Handler.Create, mw.ServerWrite)
		authGroup.PATCH("/servers/:id", cfg.Handler.Update, mw.ServerWrite)
		authGroup.POST("/servers/test-ssh", cfg.Handler.TestSSH, mw.ServerWrite)
		authGroup.DELETE("/servers/:id", cfg.Handler.Delete, mw.ServerWrite)
		authGroup.GET("/servers/:id/components", cfg.Handler.GetComponents, mw.ServerRead)
		authGroup.POST("/servers/:id/scan-ports", cfg.Handler.ScanPorts, mw.ServerRead)
		authGroup.GET("/servers/:id/rate-limit", cfg.Handler.GetRateLimit, mw.ServerRead)
		authGroup.PUT("/servers/:id/rate-limit", cfg.Handler.UpdateRateLimit, mw.ServerWrite)
		authGroup.GET("/servers/:id/network-settings", cfg.Handler.GetNetworkSettings, mw.ServerRead)
		authGroup.PATCH("/servers/:id/network-settings", cfg.Handler.UpdateNetworkSettings, mw.ServerWrite)
		authGroup.GET("/servers/:id/managed/info", cfg.Handler.GetManagedInfo, mw.ServerRead)
		authGroup.POST("/servers/:id/managed/boot-logs", cfg.Handler.GetManagedBootLogs, mw.ServerRead)
		authGroup.GET("/servers/:id/managed/ssh", cfg.Handler.GetManagedSshStatus, mw.ServerRead)
		authGroup.PATCH("/servers/:id/managed/ssh", cfg.Handler.SetManagedSsh, mw.ServerWrite)
		authGroup.PUT("/servers/:id/managed/ssh/key", cfg.Handler.SetManagedSshKey, mw.ServerWrite)
		authGroup.PUT("/servers/:id/managed/ssh/password", cfg.Handler.SetManagedSshPassword, mw.ServerWrite)
		authGroup.GET("/servers/:id/managed/runtime-api", cfg.Handler.GetManagedRuntimeStatus, mw.ServerRead)
		authGroup.POST("/servers/:id/managed/runtime-api/enable", cfg.Handler.EnableManagedRuntime, mw.ServerWrite)
		authGroup.POST("/servers/:id/managed/runtime-api/credential", cfg.Handler.GetManagedRuntimeCredential, mw.ServerWrite)
		authGroup.POST("/servers/:id/managed/runtime-api/rotate", cfg.Handler.RotateManagedRuntimeCredential, mw.ServerWrite)
		authGroup.GET("/servers/:id/managed/workloads", cfg.Handler.GetManagedWorkloads, mw.ServerRead)
		authGroup.POST("/servers/:id/managed/workloads", cfg.Handler.CreateManagedWorkload, mw.ServerWrite)
		authGroup.POST("/servers/:id/managed/workloads/control", cfg.Handler.ControlManagedWorkload, mw.ServerWrite)
		authGroup.POST("/servers/:id/managed/workloads/logs", cfg.Handler.GetManagedWorkloadLogs, mw.ServerRead)
	}

	if cfg.NetworksHandler != nil {
		networks.RegisterRoutes(authGroup, cfg.NetworksHandler, mw.Admin)
	}

	if cfg.WSHandler != nil {
		apiGroup.GET("/ws/servers/:serverId/metrics", cfg.WSHandler.Handle)
	}

	if cfg.TakeoverHandler != nil {
		takeover.RegisterRoutes(authGroup, apiGroup, cfg.TakeoverHandler, mw.Admin)
	}

	if cfg.ManagedHandler != nil {
		managed.RegisterRoutes(authGroup, cfg.ManagedHandler, managed.RouteMiddleware{
			OrgMember:   mw.OrgMember,
			OrgAdmin:    mw.OrgAdmin,
			ServerRead:  mw.ServerRead,
			ServerWrite: mw.ServerWrite,
		})
	}
}
