package servers

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Handler         *Handler
	ManagedHandler  *ManagedHandler
	TakeoverHandler *TakeoverHandler
	NetworksHandler *ServerNetworksHandler
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
		authGroup.GET("/system/networks", cfg.NetworksHandler.List)
		authGroup.GET("/system/networks/:id", cfg.NetworksHandler.Get)
		authGroup.POST("/system/networks", cfg.NetworksHandler.Create, mw.Admin)
		authGroup.PUT("/system/networks/:id/access", cfg.NetworksHandler.ReviseAccess, mw.Admin)
		authGroup.DELETE("/system/networks/:id", cfg.NetworksHandler.Delete, mw.Admin)
		authGroup.GET("/system/networks/operations/:operationId/watchdog", cfg.NetworksHandler.GetWatchdog)
		authGroup.POST("/system/networks/operations/:operationId/watchdog/commit", cfg.NetworksHandler.CommitWatchdog, mw.Admin)
		authGroup.POST("/system/networks/operations/:operationId/watchdog/rollback", cfg.NetworksHandler.RollbackWatchdog, mw.Admin)
	}

	if cfg.WSHandler != nil {
		apiGroup.GET("/ws/servers/:serverId/metrics", cfg.WSHandler.Handle)
	}

	if cfg.TakeoverHandler != nil {
		apiGroup.POST("/system/takeover/scan", cfg.TakeoverHandler.Scan, mw.Admin)
		apiGroup.POST("/system/takeover/adopt", cfg.TakeoverHandler.Adopt, mw.Admin)
		authGroup.GET("/system/takeover/runs", cfg.TakeoverHandler.ListRuns, mw.Admin)
		authGroup.GET("/system/takeover/runs/:id", cfg.TakeoverHandler.GetRun, mw.Admin)
	}

	if cfg.ManagedHandler != nil {
		authGroup.GET("/managed/catalog", cfg.ManagedHandler.Catalog)
		authGroup.GET("/organizations/:id/managed/credentials", cfg.ManagedHandler.ListCredentials, mw.OrgMember)
		authGroup.POST("/organizations/:id/managed/credentials", cfg.ManagedHandler.CreateCredential, mw.OrgAdmin)
		authGroup.DELETE("/organizations/:id/managed/credentials/:credentialId", cfg.ManagedHandler.DeleteCredential, mw.OrgAdmin)
		authGroup.GET("/organizations/:id/managed/credentials/:credentialId/offers", cfg.ManagedHandler.Offers, mw.OrgMember, mw.ServerRead)
		authGroup.GET("/organizations/:id/managed/quota", cfg.ManagedHandler.GetQuota, mw.OrgMember)
		authGroup.PUT("/organizations/:id/managed/quota", cfg.ManagedHandler.SetQuota, mw.OrgAdmin)
		authGroup.POST("/organizations/:id/managed/servers/review", cfg.ManagedHandler.ReviewProvision, mw.OrgMember, mw.ServerWrite)
		authGroup.POST("/organizations/:id/managed/servers/:serverId/refresh", cfg.ManagedHandler.RefreshServer, mw.OrgMember, mw.ServerRead)
		authGroup.POST("/organizations/:id/managed/servers/:serverId/resize/review", cfg.ManagedHandler.ReviewResize, mw.OrgAdmin)
		authGroup.POST("/organizations/:id/managed/servers/:serverId/delete/review", cfg.ManagedHandler.ReviewDelete, mw.OrgAdmin)
		authGroup.GET("/managed-operations/:operationId", cfg.ManagedHandler.GetOperation, mw.ServerRead)
		authGroup.GET("/managed-operations/:operationId/stream", cfg.ManagedHandler.StreamOperation, mw.ServerRead)
		authGroup.POST("/managed-operations/:operationId/apply", cfg.ManagedHandler.ApplyOperation, mw.ServerWrite)
	}
}
