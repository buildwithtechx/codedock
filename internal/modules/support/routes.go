package support

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, h *SupportHandler) {
	if h == nil || authGroup == nil {
		return
	}
	authGroup.GET("/cloud/support/session", h.GetSession)
	authGroup.GET("/cloud/support/mine", h.ListTickets)
	authGroup.POST("/cloud/support/mine", h.CreateTicket)
	authGroup.GET("/cloud/support/mine/:id", h.GetTicket)
	authGroup.POST("/cloud/support/mine/:id/replies", h.ReplyTicket)
	authGroup.PATCH("/cloud/support/mine/:id", h.UpdateStatus)
}
