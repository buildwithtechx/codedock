package auth

import (
	"net/http"
	"strconv"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	authservices "codedock.run/codedock/internal/services/auth"
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
)

type AuditLogHandler struct {
	auditService *authservices.AuditService
}

func NewAuditLogHandler(as *authservices.AuditService) *AuditLogHandler {
	return &AuditLogHandler{auditService: as}
}

func (h *AuditLogHandler) List(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil || user.Role != "admin" {
		return utils.Error(c, http.StatusForbidden, "insufficient permissions to view audit logs")
	}

	limitParam := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 100
	}

	offsetParam := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	logs, err := h.auditService.ListLogs(c.Request().Context(), limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if logs == nil {
		logs = []models.AuditLog{}
	}

	return utils.Success(c, "Audit logs fetched", logs)
}
