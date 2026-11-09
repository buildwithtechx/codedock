package system

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	System    *SystemHandler
	Updater   *UpdaterHandler
	AdminAuth echo.MiddlewareFunc
}

func RegisterRoutes(authGroup, apiGroup *echo.Group, cfg RoutesConfig) {
	admin := cfg.AdminAuth

	if cfg.System != nil {
		authGroup.GET("/system/stats", cfg.System.GetStats)
		if admin != nil {
			apiGroup.POST("/system/restart", cfg.System.Restart, admin)
			apiGroup.POST("/system/maintenance/cleanup", cfg.System.Cleanup, admin)
		} else {
			apiGroup.POST("/system/restart", cfg.System.Restart)
			apiGroup.POST("/system/maintenance/cleanup", cfg.System.Cleanup)
		}
		apiGroup.POST("/system/local-cli-session", cfg.System.CreateLocalCLISession)
	}

	if cfg.Updater != nil {
		authGroup.GET("/settings/updates/status", cfg.Updater.GetUpdateStatus)
		if admin != nil {
			apiGroup.POST("/settings/updates/check", cfg.Updater.CheckUpdate, admin)
			apiGroup.POST("/settings/updates/deploy", cfg.Updater.DeployUpdate, admin)
		} else {
			apiGroup.POST("/settings/updates/check", cfg.Updater.CheckUpdate)
			apiGroup.POST("/settings/updates/deploy", cfg.Updater.DeployUpdate)
		}
	}
}
