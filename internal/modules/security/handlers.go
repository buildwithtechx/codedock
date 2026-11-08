package security

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

// Handler handles HTTP requests for platform and project security advisories
type Handler struct {
	service SecurityService
}

// NewHandler creates a new security HTTP handler
func NewHandler(service SecurityService) *Handler {
	return &Handler{service: service}
}

// GetAdvisories returns all platform security advisories and affected projects
func (h *Handler) GetAdvisories(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	resp, err := h.service.GetAdvisories(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Advisories retrieved", resp)
}

// GetProjectSecurity returns findings and remediation options for a project
func (h *Handler) GetProjectSecurity(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}

	resp, err := h.service.GetProjectSecurity(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Project security retrieved", resp)
}

// CheckProject triggers an on-demand security scan for a project
func (h *Handler) CheckProject(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}

	resp, err := h.service.CheckProject(c.Request().Context(), projectID, true)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Project security checked", resp)
}

// SetProjectPolicy updates automatic advisory remediation policy for a project
func (h *Handler) SetProjectPolicy(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}

	var req SetPolicyRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	if err := h.service.SetProjectPolicy(c.Request().Context(), projectID, req.AutoRefresh); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Security policy updated", nil)
}

// ApplyRemediation queues an automated image rebuild or runtime patch for an advisory
func (h *Handler) ApplyRemediation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}

	var req ApplyRemediationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.AdvisoryID == "" {
		return utils.Error(c, http.StatusBadRequest, "advisoryId is required")
	}

	resp, err := h.service.ApplyRemediation(c.Request().Context(), projectID, req.AdvisoryID, req.DeploymentID, req.ServiceID, req.RuntimePatchID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Security remediation queued", resp)
}
