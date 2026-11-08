package clusters

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Handler     *Handler
	DataHandler *DataHandler
	Middleware  RouteMiddleware
}

type RouteMiddleware struct {
	Admin          echo.MiddlewareFunc
	DatabaseManage echo.MiddlewareFunc
}

func RegisterRoutes(group *echo.Group, cfg RoutesConfig) {
	mw := cfg.Middleware

	if cfg.DataHandler != nil {
		group.GET("/projects/:id/clusters/:clusterId/databases", cfg.DataHandler.List, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/review", cfg.DataHandler.Review, mw.Admin)
		if mw.DatabaseManage != nil {
			group.GET("/projects/:id/clusters/:clusterId/databases/:databaseId/credentials", cfg.DataHandler.Credentials, mw.Admin, mw.DatabaseManage)
		} else {
			group.GET("/projects/:id/clusters/:clusterId/databases/:databaseId/credentials", cfg.DataHandler.Credentials, mw.Admin)
		}
		group.GET("/projects/:id/clusters/:clusterId/databases/:databaseId/backups", cfg.DataHandler.Backups, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/:databaseId/verify", cfg.DataHandler.Verify, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/:databaseId/snapshots/review", cfg.DataHandler.ReviewRedisSnapshot, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/:databaseId/snapshots/apply", cfg.DataHandler.ApplyRedisSnapshot, mw.Admin)
		group.GET("/projects/:id/clusters/:clusterId/databases/:databaseId/snapshots", cfg.DataHandler.ListRedisSnapshots, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/:databaseId/restores/review", cfg.DataHandler.ReviewRedisRestore, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/:databaseId/restores/apply", cfg.DataHandler.ApplyRedisRestore, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/:databaseId/failover/review", cfg.DataHandler.ReviewRedisFailover, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/:databaseId/failover/apply", cfg.DataHandler.ApplyRedisFailover, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/bindings/review", cfg.DataHandler.ReviewBinding, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/bindings/apply", cfg.DataHandler.ApplyBinding, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/lifecycle/review", cfg.DataHandler.ReviewLifecycle, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/databases/lifecycle/apply", cfg.DataHandler.ApplyLifecycle, mw.Admin)
		group.POST("/cluster-data-operations/:operationId/apply", cfg.DataHandler.Apply, mw.Admin)
		group.GET("/cluster-data-operations/:operationId", cfg.DataHandler.Operation, mw.Admin)
		group.POST("/cluster-data-operations/:operationId/cancel", cfg.DataHandler.Cancel, mw.Admin)
	}

	if cfg.Handler != nil {
		group.GET("/projects/:id/clusters", cfg.Handler.List, mw.Admin)
		group.GET("/projects/:id/clusters/:clusterId/storage", cfg.Handler.Storage, mw.Admin)
		group.GET("/projects/:id/clusters/:clusterId/storage/health", cfg.Handler.StorageHealth, mw.Admin)
		group.GET("/projects/:id/clusters/:clusterId/preflight", cfg.Handler.Preflight, mw.Admin)
		group.POST("/projects/:id/clusters/:clusterId/storage/review", cfg.Handler.ReviewStorage, mw.Admin)
		group.POST("/cluster-storage-operations/:operationId/apply", cfg.Handler.ApplyStorage, mw.Admin)
		group.POST("/projects/:id/clusters/review", cfg.Handler.Review, mw.Admin)
		group.POST("/cluster-operations/:operationId/apply", cfg.Handler.Apply, mw.Admin)
		group.GET("/cluster-operations/:operationId", cfg.Handler.GetOperation, mw.Admin)
		group.POST("/cluster-operations/:operationId/cancel", cfg.Handler.Cancel, mw.Admin)
		group.GET("/system/compute-clusters/:clusterId/scaling", cfg.Handler.GetScaling, mw.Admin)
		group.POST("/system/compute-clusters/:clusterId/scaling/setup", cfg.Handler.SetupScaling, mw.Admin)
		group.POST("/system/compute-clusters/:clusterId/scaling/retry", cfg.Handler.RetryScaling, mw.Admin)
		group.POST("/system/compute-clusters/:clusterId/scaling/remove", cfg.Handler.RemoveScaling, mw.Admin)
	}
}
