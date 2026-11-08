package git

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AppsHandler struct {
	service *GitAppsService
}

func NewAppsHandler(s *GitAppsService) *AppsHandler {
	return &AppsHandler{service: s}
}

func (h *AppsHandler) ListGithubApps(c echo.Context) error {
	apps, err := h.service.ListGithubApps(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if apps == nil {
		apps = []GithubApp{}
	}
	return utils.Success(c, "Operation successful", apps)
}

func (h *AppsHandler) GetGithubApp(c echo.Context) error {
	id := c.Param("id")
	app, err := h.service.GetGithubApp(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if app == nil {
		return utils.Error(c, http.StatusNotFound, "App not found")
	}
	return utils.Success(c, "Operation successful", app)
}

func (h *AppsHandler) SaveGithubApp(c echo.Context) error {
	var app GithubApp
	if err := c.Bind(&app); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if err := h.service.SaveGithubApp(c.Request().Context(), &app); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", app)
}

func (h *AppsHandler) DeleteGithubApp(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteGithubApp(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", DisconnectResponse{Status: "deleted"})
}

func (h *AppsHandler) ExchangeGithubManifestCode(c echo.Context) error {
	var req GitAppsManifestRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	app, err := h.service.ExchangeGithubManifestCode(c.Request().Context(), req.Code)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Created(c, "GitHub App configured", app)
}
