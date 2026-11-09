package dns

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, h *Handler, middleware ...echo.MiddlewareFunc) {
	if authGroup == nil {
		return
	}
	if h == nil {
		h = &Handler{}
	}

	authGroup.POST("/dns", h.Create, middleware...)
	authGroup.GET("/dns", h.List, middleware...)
	authGroup.PUT("/dns/:id", h.Update, middleware...)
	authGroup.DELETE("/dns/:id", h.Delete, middleware...)
}

func RegisterIngressRoutes(authGroup *echo.Group, h *IngressHandler, adminMw ...echo.MiddlewareFunc) {
	if authGroup == nil {
		return
	}
	if h == nil {
		h = &IngressHandler{}
	}

	authGroup.GET("/dns/ingress", h.List)
	authGroup.POST("/dns/ingress/discover", h.Discover)
	authGroup.POST("/dns/ingress", h.Connect, adminMw...)
	authGroup.PUT("/dns/ingress/:id", h.Update, adminMw...)
	authGroup.POST("/dns/ingress/:id/retry", h.Retry, adminMw...)
	authGroup.DELETE("/dns/ingress/:id", h.Disconnect, adminMw...)
	authGroup.GET("/dns/ingress/domains/:domainId", h.GetDomain)
	authGroup.PUT("/dns/ingress/domains/:domainId", h.SetDomain, adminMw...)
}
