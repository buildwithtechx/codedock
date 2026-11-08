package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/modules/security"
	"codedock/internal/utils"
)

type SecurityHandler struct {
	securityService security.SecurityService
}

func NewSecurityHandler(service security.SecurityService) *SecurityHandler {
	return &SecurityHandler{securityService: service}
}

func (h *SecurityHandler) GetAdvisories(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	resp, err := h.securityService.GetAdvisories(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Advisories retrieved", resp)
}

func (h *SecurityHandler) GetProjectSecurity(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}

	resp, err := h.securityService.GetProjectSecurity(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Project security retrieved", resp)
}

func (h *SecurityHandler) CheckProject(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "project id is required")
	}

	resp, err := h.securityService.CheckProject(c.Request().Context(), projectID, true)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Project security checked", resp)
}

type SetPolicyRequest struct {
	AutoRefresh bool `json:"autoRefresh"`
}

func (h *SecurityHandler) SetProjectPolicy(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
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

	if err := h.securityService.SetProjectPolicy(c.Request().Context(), projectID, req.AutoRefresh); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Security policy updated", nil)
}

type ApplyRemediationRequest struct {
	AdvisoryID     string  `json:"advisoryId"`
	DeploymentID   string  `json:"deploymentId"`
	ServiceID      *string `json:"serviceId,omitempty"`
	RuntimePatchID *string `json:"runtimePatchId,omitempty"`
}

func (h *SecurityHandler) ApplyRemediation(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
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

	resp, err := h.securityService.ApplyRemediation(c.Request().Context(), projectID, req.AdvisoryID, req.DeploymentID, req.ServiceID, req.RuntimePatchID)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Security remediation queued", resp)
}
