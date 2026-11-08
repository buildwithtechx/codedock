package domains

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Handler *DomainHandler
}

func RegisterRoutes(g *echo.Group, cfg RoutesConfig) {
	if cfg.Handler != nil {
		g.GET("/domains", cfg.Handler.ListAll)
		g.GET("/projects/:id/domains", cfg.Handler.ListByProject)
		g.GET("/services/:id/domains", cfg.Handler.ListByService)
		g.POST("/services/:id/domains", cfg.Handler.Create)
		g.DELETE("/domains/:id", cfg.Handler.Delete)
		g.POST("/domains/:id/verify", cfg.Handler.Verify)
		g.GET("/domains/:id/verify", cfg.Handler.Verify)
	}
}
