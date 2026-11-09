package templates

import "github.com/labstack/echo/v4"

type RoutesConfig struct {
	Handler *OneClickHandler
}

func RegisterRoutes(g *echo.Group, cfg RoutesConfig) {
	if cfg.Handler == nil {
		cfg.Handler = NewOneClickHandler(nil)
	}
	g.GET("/one-click", cfg.Handler.List)
	g.GET("/one-click/:id", cfg.Handler.Get)
	g.POST("/one-click/review", cfg.Handler.Review)
	g.POST("/one-click/deploy", cfg.Handler.Deploy)
	g.POST("/one-click/custom", cfg.Handler.AddCustom)
	g.DELETE("/one-click/custom/:id", cfg.Handler.DeleteCustom)
}
