package monitoring

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	System        *SystemHandler
	Metrics       *MetricsHandler
	Logs          *LogHandler
	Updater       *UpdaterHandler
	Notification  *NotificationHandler
	NotifSettings *NotificationSettingsHandler
	ScheduledTask *ScheduledTaskHandler
	Ingress       *IngressHandler
	Analytics     *AnalyticsHandler
	Attention     *AttentionHandler
	Operation     *OperationHandler
	Example       *ExampleHandler
	ServiceLogsWS *ServiceLogsWSHandler
	AdminAuth     echo.MiddlewareFunc
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

	if cfg.Metrics != nil {
		authGroup.GET("/services/:serviceId/metrics/historical", cfg.Metrics.GetHistoricalMetrics)
	}
	if cfg.Logs != nil {
		authGroup.GET("/services/:serviceId/logs/historical", cfg.Logs.GetHistoricalLogs)
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
	if cfg.Notification != nil {
		authGroup.POST("/notifications/test", cfg.Notification.TestNotification)
	}
	if cfg.NotifSettings != nil {
		authGroup.GET("/settings/notifications", cfg.NotifSettings.GetNotificationSettings)
		authGroup.POST("/settings/notifications", cfg.NotifSettings.UpdateNotificationSettings)
		authGroup.GET("/settings/notifications/deliveries", cfg.NotifSettings.ListDeliveries)
		authGroup.POST("/settings/notifications/devices", cfg.NotifSettings.RegisterDeviceToken)
	}
	if cfg.ScheduledTask != nil {
		authGroup.GET("/scheduled-tasks", cfg.ScheduledTask.ListProjectScheduledTasks)
		authGroup.POST("/scheduled-tasks", cfg.ScheduledTask.Create)
		authGroup.GET("/scheduled-tasks/:id", cfg.ScheduledTask.Get)
		authGroup.PUT("/scheduled-tasks/:id", cfg.ScheduledTask.Update)
		authGroup.DELETE("/scheduled-tasks/:id", cfg.ScheduledTask.Delete)
		authGroup.POST("/scheduled-tasks/:id/trigger", cfg.ScheduledTask.Run)
		authGroup.GET("/jobs", cfg.ScheduledTask.ListProjectScheduledTasks)
		authGroup.POST("/jobs", cfg.ScheduledTask.Create)
		authGroup.GET("/jobs/:id", cfg.ScheduledTask.Get)
		authGroup.PUT("/jobs/:id", cfg.ScheduledTask.Update)
		authGroup.DELETE("/jobs/:id", cfg.ScheduledTask.Delete)
		authGroup.POST("/jobs/:id/trigger", cfg.ScheduledTask.Run)
	}
	if cfg.Ingress != nil {
		authGroup.GET("/dns/ingress", cfg.Ingress.List)
		authGroup.POST("/dns/ingress/discover", cfg.Ingress.Discover)
		authGroup.POST("/dns/ingress", cfg.Ingress.Connect)
		authGroup.PUT("/dns/ingress/:id", cfg.Ingress.Update)
		authGroup.POST("/dns/ingress/:id/retry", cfg.Ingress.Retry)
		authGroup.DELETE("/dns/ingress/:id", cfg.Ingress.Disconnect)
		authGroup.GET("/dns/ingress/domains/:domainId", cfg.Ingress.GetDomain)
		if admin != nil {
			authGroup.PUT("/dns/ingress/domains/:domainId", cfg.Ingress.SetDomain, admin)
		} else {
			authGroup.PUT("/dns/ingress/domains/:domainId", cfg.Ingress.SetDomain)
		}
	}
	if cfg.Operation != nil {
		authGroup.GET("/operations", cfg.Operation.List)
		authGroup.GET("/operations/:operationId", cfg.Operation.Get)
		authGroup.POST("/operations/:operationId/cancel", cfg.Operation.Cancel)
	}
	if cfg.Example != nil {
		authGroup.GET("/examples", cfg.Example.List)
	}
	if cfg.ServiceLogsWS != nil {
		apiGroup.GET("/services/:serviceId/logs", cfg.ServiceLogsWS.Handle)
	}
}
