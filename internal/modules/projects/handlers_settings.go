package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ProjectSettingsHandler struct {
	settingsService *ProjectSettingsService
}

func NewProjectSettingsHandler(s *ProjectSettingsService) *ProjectSettingsHandler {
	return &ProjectSettingsHandler{settingsService: s}
}

func (h *ProjectSettingsHandler) ListTokens(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId")
	}
	list, err := h.settingsService.ListTokens(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *ProjectSettingsHandler) CreateToken(c echo.Context) error {
	projectID := c.Param("projectId")
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId")
	}
	var req CreateTokenRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	t := &ProjectToken{
		ProjectID:     projectID,
		Name:          req.Name,
		EnvironmentID: req.EnvironmentID,
		Scopes:        req.Scopes,
		IPAllowlist:   req.IPAllowlist,
		ExpiresAt:     req.ExpiresAt,
	}
	token, raw, err := h.settingsService.CreateToken(c.Request().Context(), t)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Token created successfully", TokenCreatedPayload{
		ID:          token.ID,
		Name:        token.Name,
		Token:       raw,
		Scopes:      token.Scopes,
		IPAllowlist: token.IPAllowlist,
		ExpiresAt:   token.ExpiresAt,
		CreatedAt:   token.CreatedAt,
	})
}

func (h *ProjectSettingsHandler) DeleteToken(c echo.Context) error {
	projectID := c.Param("projectId")
	id := c.Param("id")
	if projectID == "" || id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing projectId or id")
	}
	if err := h.settingsService.DeleteToken(c.Request().Context(), id, projectID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
