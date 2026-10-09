package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/http/middleware"
	"codedock/internal/models"
	projectservices "codedock/internal/services/projects"
	"codedock/internal/utils"
)

type OneClickHandler struct {
	service        *projectservices.OneClickService
	projectService *projectservices.ProjectService
}

func NewOneClickHandler(s *projectservices.OneClickService, ps *projectservices.ProjectService) *OneClickHandler {
	return &OneClickHandler{service: s, projectService: ps}
}

type oneClickInstallRequest struct {
	AppID         string            `json:"appId" form:"appId"`
	ProjectID     string            `json:"projectId" form:"projectId"`
	EnvironmentID string            `json:"environmentId" form:"environmentId"`
	Name          string            `json:"name" form:"name"`
	Secrets       map[string]string `json:"secrets" form:"secrets"`
	Environment   map[string]string `json:"environment" form:"environment"`
	HostPort      int               `json:"hostPort" form:"hostPort"`
	Domain        string            `json:"domain" form:"domain"`
	Digest        string            `json:"digest" form:"digest"`
}

func (h *OneClickHandler) List(c echo.Context) error {
	return utils.Success(c, "Available one-click apps", h.service.ListApps())
}

func (h *OneClickHandler) Get(c echo.Context) error {
	app, err := h.service.GetApp(c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "One-click app", app)
}

func (h *OneClickHandler) Review(c echo.Context) error {
	var req oneClickInstallRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.ProjectID == "" {
		return utils.Error(c, http.StatusBadRequest, "projectId is required")
	}
	if err := h.authorize(c, req.ProjectID); err != nil {
		return err
	}
	preview, err := h.service.ReviewInstall(c.Request().Context(), installInput(req))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Install preview", preview)
}

func (h *OneClickHandler) Deploy(c echo.Context) error {
	var req oneClickInstallRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.ProjectID == "" {
		return utils.Error(c, http.StatusBadRequest, "projectId is required")
	}
	if err := h.authorize(c, req.ProjectID); err != nil {
		return err
	}
	result, err := h.service.InstallApp(c.Request().Context(), installInput(req))
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "App installed", result)
}

func (h *OneClickHandler) authorize(c echo.Context, projectID string) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if user.Role != models.UserRoleAdmin && user.Role != models.UserRoleOwner {
		if !h.projectService.HasPermission(c.Request().Context(), projectID, user.UserID, models.UserRole(user.Role), "") {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions for this project")
		}
	}
	return nil
}

func installInput(req oneClickInstallRequest) models.InstallAppInput {
	return models.InstallAppInput{
		AppID:         req.AppID,
		ProjectID:     req.ProjectID,
		EnvironmentID: req.EnvironmentID,
		Name:          req.Name,
		Secrets:       req.Secrets,
		Environment:   req.Environment,
		HostPort:      req.HostPort,
		Domain:        req.Domain,
		Digest:        req.Digest,
	}
}
