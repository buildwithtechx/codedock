package monitoring

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Metrics       *MetricsHandler
	Logs          *LogHandler
	Analytics     *AnalyticsHandler
	Attention     *AttentionHandler
	Example       *ExampleHandler
	ServiceLogsWS *ServiceLogsWSHandler
	AdminAuth     echo.MiddlewareFunc
}

func RegisterRoutes(authGroup, apiGroup *echo.Group, cfg RoutesConfig) {
	if cfg.Metrics != nil {
		authGroup.GET("/services/:serviceId/metrics/historical", cfg.Metrics.GetHistoricalMetrics)
	}
	if cfg.Logs != nil {
		authGroup.GET("/services/:serviceId/logs/historical", cfg.Logs.GetHistoricalLogs)
	}
	if cfg.Example != nil {
		authGroup.GET("/examples", cfg.Example.List)
	}
	if cfg.ServiceLogsWS != nil {
		apiGroup.GET("/services/:serviceId/logs", cfg.ServiceLogsWS.Handle)
	}
}
