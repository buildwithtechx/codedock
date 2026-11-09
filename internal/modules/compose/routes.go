package compose

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Handler *ComposeStackHandler
}

func RegisterRoutes(g *echo.Group, cfg RoutesConfig) {
	if cfg.Handler == nil {
		cfg.Handler = NewComposeStackHandler(nil)
	}
	g.GET("/projects/:id/stacks", cfg.Handler.List)
	g.GET("/projects/:id/stacks/:stackId/config", cfg.Handler.Config)
	g.POST("/projects/:id/stacks/review", cfg.Handler.Review)
	g.POST("/projects/:id/stacks", cfg.Handler.Save)
	g.POST("/projects/:id/stacks/:stackId/deploy", cfg.Handler.Deploy)
	g.POST("/projects/:id/stacks/:stackId/cancel", cfg.Handler.Cancel)
	g.POST("/projects/:id/stacks/:stackId/restart", cfg.Handler.Restart)
	g.POST("/projects/:id/stacks/:stackId/stop", cfg.Handler.Stop)
}
