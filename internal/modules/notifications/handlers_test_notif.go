package notifications

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type NotificationHandler struct {
	service NotificationService
}

func NewNotificationHandler(ns NotificationService) *NotificationHandler {
	return &NotificationHandler{service: ns}
}

func (h *NotificationHandler) TestNotification(c echo.Context) error {
	if c.Request().Method != http.MethodPost {
		return utils.Error(c, http.StatusMethodNotAllowed, "Method not allowed")
	}
	var req TestNotificationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.Provider == "" {
		return utils.Error(c, http.StatusBadRequest, "Provider required")
	}
	if err := h.service.TestNotification(c.Request().Context(), req.Provider); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Test notification queued", nil)
}
