package chatapps

import "github.com/labstack/echo/v4"

// RegisterRoutes registers all Chat Apps endpoints onto Echo groups
func RegisterRoutes(apiGroup, authGroup *echo.Group, h *Handler) {
	authGroup.GET("/chat-apps", h.List)
	authGroup.POST("/chat-apps", h.Create)
	authGroup.GET("/chat-apps/:id", h.Get)
	authGroup.PATCH("/chat-apps/:id", h.Update)
	authGroup.DELETE("/chat-apps/:id", h.Delete)
	authGroup.POST("/chat-apps/:id/test", h.Test)
	authGroup.GET("/chat-apps/:id/identities", h.ListIdentities)
	authGroup.DELETE("/chat-apps/identities/:id", h.DeleteIdentity)
	authGroup.GET("/chat-apps/pairing/challenge/:token", h.GetChallenge)
	authGroup.POST("/chat-apps/pairing/verify", h.VerifyPairing)

	apiGroup.POST("/chat-apps/webhooks/:id", h.WebhookIngress)
}
