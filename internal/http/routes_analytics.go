package http

import "github.com/labstack/echo/v4"

func (s *Server) registerAnalyticsRoutes(authGroup *echo.Group) {
	member := s.authGuard.RequireOrgRole("member")
	admin := s.authGuard.RequireOrgRole("admin")
	read := s.authGuard.RequireScope("server:read")
	write := s.authGuard.RequireScope("server:write")

	authGroup.GET("/organizations/:id/analytics/periods", s.analyticsHandler.Periods, member, read)
	authGroup.GET("/organizations/:id/analytics/dashboard", s.analyticsHandler.Dashboard, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/summary", s.analyticsHandler.Summary, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/overview", s.analyticsHandler.Overview, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/geo", s.analyticsHandler.Geo, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/live", s.analyticsHandler.Live, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/paths", s.analyticsHandler.GetPaths, member, read)
	authGroup.PUT("/organizations/:id/projects/:projectId/analytics/paths", s.analyticsHandler.SetPaths, member, write)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/deployments", s.analyticsHandler.DeploymentStats, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage", s.analyticsHandler.Usage, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/resources", s.analyticsHandler.Resources, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage/history", s.analyticsHandler.UsageHistory, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/usage/stream", s.analyticsHandler.UsageStream, member, read)
	authGroup.GET("/organizations/:id/projects/:projectId/analytics/container", s.analyticsHandler.Container, member, read)

	authGroup.GET("/organizations/:id/attention", s.attentionHandler.List, member, read)
	authGroup.POST("/organizations/:id/attention/evaluate", s.attentionHandler.Evaluate, member, write)
	authGroup.DELETE("/organizations/:id/attention/history", s.attentionHandler.ClearHistory, admin, write)
	authGroup.POST("/organizations/:id/attention/:issueId/ack", s.attentionHandler.Acknowledge, member, write)
	authGroup.POST("/organizations/:id/attention/:issueId/resolve", s.attentionHandler.Resolve, member, write)
	authGroup.POST("/organizations/:id/attention/:issueId/act", s.attentionHandler.Act, member, write)
}
