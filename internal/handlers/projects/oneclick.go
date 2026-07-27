package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	projectservices "codedock.run/codedock/internal/services/projects"
	"codedock.run/codedock/internal/utils"
)

type OneClickHandler struct {
	service        *projectservices.OneClickService
	projectService *projectservices.ProjectService
}

func NewOneClickHandler(s *projectservices.OneClickService, ps *projectservices.ProjectService) *OneClickHandler {
	return &OneClickHandler{service: s, projectService: ps}
}

type oneClickDeployRequest struct {
	AppID     string `json:"appId" form:"appId"`
	ProjectID string `json:"projectId" form:"projectId"`
	Name      string `json:"name" form:"name"`
}

func (h *OneClickHandler) List(c echo.Context) error {
	return utils.Success(c, "Available one-click apps", h.service.ListApps())
}

func (h *OneClickHandler) Deploy(c echo.Context) error {
	var req oneClickDeployRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.ProjectID == "" {
		return utils.Error(c, http.StatusBadRequest, "projectId is required")
	}

	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	if user.Role != "admin" {
		if !h.projectService.HasPermission(c.Request().Context(), req.ProjectID, user.UserID, models.UserRole(user.Role), "") {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions for this project")
		}
	}

	db, err := h.service.DeployApp(c.Request().Context(), req.AppID, req.ProjectID, req.Name)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "App deployed", db)
}
