package monitoring

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

type RouteMiddleware struct {
	OrgMember   echo.MiddlewareFunc
	OrgAdmin    echo.MiddlewareFunc
	ReadScope   echo.MiddlewareFunc
	WriteScope  echo.MiddlewareFunc
	ServiceAuth echo.MiddlewareFunc
}

type RoutesConfig struct {
	Metrics       *MetricsHandler
	Logs          *LogHandler
	Analytics     *AnalyticsHandler
	Attention     *AttentionHandler
	Example       *ExampleHandler
	ServiceLogsWS *ServiceLogsWSHandler
	AdminAuth     echo.MiddlewareFunc
	Middleware    RouteMiddleware
}

func RegisterRoutes(authGroup, apiGroup *echo.Group, cfg RoutesConfig) {
	mw := cfg.Middleware

	if cfg.Metrics != nil {
		authGroup.GET("/services/:serviceId/metrics/historical", cfg.Metrics.GetHistoricalMetrics, withMw(mw.ServiceAuth)...)
	}
	if cfg.Logs != nil {
		authGroup.GET("/services/:serviceId/logs/historical", cfg.Logs.GetHistoricalLogs, withMw(mw.ServiceAuth)...)
	}
	if cfg.Example != nil {
		authGroup.GET("/examples", cfg.Example.List)
	}
	if cfg.ServiceLogsWS != nil && apiGroup != nil {
		apiGroup.GET("/services/:serviceId/logs", cfg.ServiceLogsWS.Handle)
	}

	if cfg.Analytics != nil {
		authGroup.GET("/organizations/:id/analytics/periods", cfg.Analytics.Periods, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/analytics/dashboard", cfg.Analytics.Dashboard, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/summary", cfg.Analytics.Summary, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/overview", cfg.Analytics.Overview, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/geo", cfg.Analytics.Geo, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/live", cfg.Analytics.Live, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/paths", cfg.Analytics.GetPaths, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.PUT("/organizations/:id/projects/:projectId/analytics/paths", cfg.Analytics.SetPaths, withMw(mw.OrgMember, mw.WriteScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/deployments", cfg.Analytics.DeploymentStats, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage", cfg.Analytics.Usage, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/resources", cfg.Analytics.Resources, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage/history", cfg.Analytics.UsageHistory, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage/stream", cfg.Analytics.UsageStream, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.GET("/organizations/:id/projects/:projectId/analytics/container", cfg.Analytics.Container, withMw(mw.OrgMember, mw.ReadScope)...)
	}

	if cfg.Attention != nil {
		authGroup.GET("/organizations/:id/attention", cfg.Attention.List, withMw(mw.OrgMember, mw.ReadScope)...)
		authGroup.POST("/organizations/:id/attention/evaluate", cfg.Attention.Evaluate, withMw(mw.OrgMember, mw.WriteScope)...)
		authGroup.DELETE("/organizations/:id/attention/history", cfg.Attention.ClearHistory, withMw(mw.OrgAdmin, mw.WriteScope)...)
		authGroup.POST("/organizations/:id/attention/:issueId/ack", cfg.Attention.Acknowledge, withMw(mw.OrgMember, mw.WriteScope)...)
		authGroup.POST("/organizations/:id/attention/:issueId/resolve", cfg.Attention.Resolve, withMw(mw.OrgMember, mw.WriteScope)...)
		authGroup.POST("/organizations/:id/attention/:issueId/act", cfg.Attention.Act, withMw(mw.OrgMember, mw.WriteScope)...)
	}
}
