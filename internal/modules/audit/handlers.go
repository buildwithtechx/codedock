package audit

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AuditLogHandler struct {
	auditService *AuditService
}

func NewAuditLogHandler(as *AuditService) *AuditLogHandler {
	return &AuditLogHandler{auditService: as}
}

func (h *AuditLogHandler) List(c echo.Context) error {
	limitParam := c.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = DefaultPageLimit
	}

	offsetParam := c.QueryParam("offset")
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		offset = 0
	}

	logs, err := h.auditService.ListLogsByCategory(c.Request().Context(), c.QueryParam("category"), c.QueryParam("q"), limit, offset)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if logs == nil {
		logs = []AuditLog{}
	}

	return utils.Success(c, "Audit logs fetched", logs)
}

func (h *AuditLogHandler) Facets(c echo.Context) error {
	facets, err := h.auditService.Facets(c.Request().Context(), c.QueryParam("category"), c.QueryParam("q"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Audit facets fetched", facets)
}

func (h *AuditLogHandler) GetSettings(c echo.Context) error {
	orgID := c.QueryParam("organizationId")
	if orgID == "" {
		orgID = c.Request().Header.Get("X-Organization-ID")
	}
	if orgID == "" {
		orgID = "default"
	}
	settings, err := h.auditService.GetSettings(c.Request().Context(), orgID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Audit settings fetched", settings)
}

func (h *AuditLogHandler) UpdateSettings(c echo.Context) error {
	var req UpdateAuditSettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request")
	}
	if req.OrganizationID == "" {
		req.OrganizationID = c.QueryParam("organizationId")
	}
	if req.OrganizationID == "" {
		req.OrganizationID = c.Request().Header.Get("X-Organization-ID")
	}
	if req.OrganizationID == "" {
		req.OrganizationID = "default"
	}
	if err := h.auditService.UpdateSettings(c.Request().Context(), &req); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Audit settings updated", req)
}
