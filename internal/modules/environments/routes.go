package environments

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	EnvHandler        *EnvironmentHandler
	ProtectionHandler *RouteProtectionHandler
	EdgeAuthHandler   *EdgeAuthHandler
	RouteRuleHandler  *RouteRuleHandler
}

func RegisterRoutes(g *echo.Group, cfg RoutesConfig) {
	if cfg.EnvHandler != nil {
		g.POST("/projects/:id/environments", cfg.EnvHandler.Create)
		g.GET("/projects/:id/environments", cfg.EnvHandler.ListByProject)
		g.DELETE("/environments/:id", cfg.EnvHandler.Delete)
		g.GET("/projects/:id/env/vars", cfg.EnvHandler.GetVars)
		g.POST("/projects/:id/env/vars", cfg.EnvHandler.SetVars)
		g.GET("/projects/:id/env/suggest", cfg.EnvHandler.SuggestVars)
	}

	if cfg.ProtectionHandler != nil {
		g.GET("/projects/:id/route-protection", cfg.ProtectionHandler.List)
		g.POST("/projects/:id/route-protection", cfg.ProtectionHandler.Change)
	}

	if cfg.RouteRuleHandler != nil {
		g.GET("/services/:serviceId/route-rules", cfg.RouteRuleHandler.List)
		g.POST("/services/:serviceId/route-rules", cfg.RouteRuleHandler.Create)
		g.PATCH("/services/:serviceId/route-rules/:ruleId", cfg.RouteRuleHandler.Update)
		g.DELETE("/services/:serviceId/route-rules/:ruleId", cfg.RouteRuleHandler.Delete)
	}
}

func RegisterEdgeAuthRoutes(e *echo.Echo, handler *EdgeAuthHandler) {
	if handler != nil {
		e.GET("/api/v1/edge/auth", handler.Verify)
		e.GET("/_codedock/access", handler.LoginForm)
		e.POST("/_codedock/access", handler.Login)
	}
}
