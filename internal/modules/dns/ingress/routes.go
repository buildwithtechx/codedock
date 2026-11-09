package ingress

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, h *IngressHandler, adminMw ...echo.MiddlewareFunc) {
	if authGroup == nil {
		return
	}
	if h == nil {
		h = &IngressHandler{}
	}

	authGroup.GET("/dns/ingress", h.List)
	authGroup.POST("/dns/ingress/discover", h.Discover)
	authGroup.POST("/dns/ingress", h.Connect)
	authGroup.PUT("/dns/ingress/:id", h.Update)
	authGroup.POST("/dns/ingress/:id/retry", h.Retry)
	authGroup.DELETE("/dns/ingress/:id", h.Disconnect)
	authGroup.GET("/dns/ingress/domains/:domainId", h.GetDomain)
	authGroup.PUT("/dns/ingress/domains/:domainId", h.SetDomain)
}
