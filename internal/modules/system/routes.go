package system

import "github.com/labstack/echo/v4"

func withMw(mws ...echo.MiddlewareFunc) []echo.MiddlewareFunc {
	var active []echo.MiddlewareFunc
	for _, m := range mws {
		if m != nil {
			active = append(active, m)
		}
	}
	return active
}

type RoutesConfig struct {
	System    *SystemHandler
	Updater   *UpdaterHandler
	AdminAuth echo.MiddlewareFunc
}

func RegisterRoutes(authGroup, apiGroup *echo.Group, cfg RoutesConfig) {
	admin := withMw(cfg.AdminAuth)

	if apiGroup != nil {
		apiGroup.GET("/github-stars", GetGithubStars)
	}

	if cfg.System != nil {
		authGroup.GET("/system/stats", cfg.System.GetStats)
		apiGroup.POST("/system/restart", cfg.System.Restart, admin...)
		apiGroup.POST("/system/maintenance/cleanup", cfg.System.Cleanup, admin...)
		apiGroup.POST("/system/local-cli-session", cfg.System.CreateLocalCLISession)
	}

	if cfg.Updater != nil {
		authGroup.GET("/settings/updates/status", cfg.Updater.GetUpdateStatus)
		apiGroup.POST("/settings/updates/check", cfg.Updater.CheckUpdate, admin...)
		apiGroup.POST("/settings/updates/deploy", cfg.Updater.DeployUpdate, admin...)
	}
}
