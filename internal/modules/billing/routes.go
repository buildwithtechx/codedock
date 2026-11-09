package billing

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(apiGroup, authGroup *echo.Group, h *BillingHandler, cloudEnabled bool) {
	if !cloudEnabled {
		if apiGroup != nil {
			unavailable := func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusNotFound, "billing is unavailable on self-hosted instances")
			}
			apiGroup.Any("/billing", unavailable)
			apiGroup.Any("/billing/*", unavailable)
		}
		return
	}
	if h == nil {
		h = &BillingHandler{}
	}
	if apiGroup != nil {
		apiGroup.POST("/billing/webhook", h.Webhook)
	}
	if authGroup != nil {
		authGroup.GET("/billing/config", h.GetConfig)
		authGroup.POST("/billing/checkout", h.CreateCheckoutSession)
		authGroup.GET("/billing/pending", h.GetPendingCheckout)
		authGroup.POST("/billing/checkouts/:id/cancel", h.CancelCheckoutSession)
	}
}
