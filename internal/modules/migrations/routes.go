package migrations

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	LifecycleHandler *LifecycleHandler
	BundleHandler    *BundleHandler
	Middleware       RouteMiddleware
}

type RouteMiddleware struct {
	OrgMember   echo.MiddlewareFunc
	OrgAdmin    echo.MiddlewareFunc
	ServerRead  echo.MiddlewareFunc
	ServerWrite echo.MiddlewareFunc
	Admin       echo.MiddlewareFunc
}

func RegisterRoutes(apiGroup, authGroup *echo.Group, cfg RoutesConfig) {
	mw := cfg.Middleware

	if cfg.LifecycleHandler != nil {
		authGroup.GET("/organizations/:id/migration/sources", cfg.LifecycleHandler.ListSources, mw.OrgMember)
		authGroup.POST("/organizations/:id/migration/sources/test", cfg.LifecycleHandler.TestSource, mw.OrgMember, mw.ServerWrite)
		authGroup.POST("/organizations/:id/migration/sources", cfg.LifecycleHandler.CreateSource, mw.OrgAdmin)
		authGroup.DELETE("/organizations/:id/migration/sources/:sourceId", cfg.LifecycleHandler.DeleteSource, mw.OrgAdmin)
		authGroup.POST("/organizations/:id/migration/scan", cfg.LifecycleHandler.Scan, mw.OrgMember, mw.ServerWrite)
		authGroup.GET("/organizations/:id/migration/scan/stream", cfg.LifecycleHandler.ScanStream, mw.OrgMember, mw.ServerWrite)
		authGroup.POST("/organizations/:id/migration/reveal-env", cfg.LifecycleHandler.RevealEnv, mw.OrgAdmin)
		authGroup.POST("/organizations/:id/migration/adopt", cfg.LifecycleHandler.Adopt, mw.OrgMember, mw.ServerWrite)
		authGroup.POST("/organizations/:id/migration/reimport", cfg.LifecycleHandler.Reimport, mw.OrgMember, mw.ServerWrite)
		authGroup.POST("/organizations/:id/migration/repo-compose", cfg.LifecycleHandler.RepoCompose, mw.OrgMember)
		authGroup.POST("/organizations/:id/migration/preview", cfg.LifecycleHandler.Preview, mw.OrgMember, mw.ServerWrite)
		authGroup.POST("/organizations/:id/migration/migrate", cfg.LifecycleHandler.StartMigration, mw.OrgMember, mw.ServerWrite)
		authGroup.POST("/organizations/:id/migration/project", cfg.LifecycleHandler.StartProjectMove, mw.OrgMember, mw.ServerWrite)
		authGroup.GET("/organizations/:id/migration/active", cfg.LifecycleHandler.Active, mw.OrgMember, mw.ServerRead)
		authGroup.GET("/organizations/:id/migration/runs", cfg.LifecycleHandler.ListRuns, mw.OrgMember, mw.ServerRead)
		authGroup.GET("/migrations/:runId", cfg.LifecycleHandler.GetRun, mw.ServerRead)
		authGroup.GET("/migrations/:runId/stream", cfg.LifecycleHandler.StreamRun, mw.ServerRead)
		authGroup.POST("/migrations/:runId/cutover", cfg.LifecycleHandler.Cutover, mw.ServerWrite)
		authGroup.POST("/migrations/:runId/cancel", cfg.LifecycleHandler.Cancel, mw.ServerWrite)
		authGroup.POST("/migrations/:runId/respond", cfg.LifecycleHandler.Respond, mw.ServerWrite)
		authGroup.POST("/migrations/:runId/resume", cfg.LifecycleHandler.Resume, mw.ServerWrite)
		authGroup.POST("/migrations/:runId/cleanup-target", cfg.LifecycleHandler.CleanupTarget, mw.ServerWrite)
		authGroup.DELETE("/migrations/:runId", cfg.LifecycleHandler.DeleteRun, mw.ServerWrite)
	}

	if cfg.BundleHandler != nil {
		if apiGroup != nil {
			apiGroup.POST("/system/setup/import", cfg.BundleHandler.ImportDuringSetup)
			apiGroup.POST("/system/export", cfg.BundleHandler.Export, mw.Admin)
			apiGroup.POST("/system/import", cfg.BundleHandler.Import, mw.Admin)
		}
	}
}
