package billing

import "github.com/labstack/echo/v4"

func RegisterRoutes(apiGroup, authGroup *echo.Group, h *BillingHandler) {
	if h == nil {
		return
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
