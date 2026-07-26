package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/utils"

	systemservices "codedock.run/codedock/internal/services/system"
)

type NotificationHandler struct {
	notificationService *systemservices.NotificationService
}

func NewNotificationHandler(ns *systemservices.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: ns}
}

type TestNotificationRequest struct {
	Provider string `json:"provider"`
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

	err := h.notificationService.TestGlobalNotification(c.Request().Context(), req.Provider)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Test notification queued", nil)
}
