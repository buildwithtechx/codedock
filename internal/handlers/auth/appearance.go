package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	handlerutils "codedock/internal/handlers/utils"
	"codedock/internal/models"
	"codedock/internal/repositories"
	"codedock/internal/utils"
)

type AppearanceHandler struct {
	userRepo repositories.UserRepository
}

func NewAppearanceHandler(userRepo repositories.UserRepository) *AppearanceHandler {
	return &AppearanceHandler{userRepo: userRepo}
}

func (h *AppearanceHandler) GetAppearance(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	pref, err := h.userRepo.GetAppearance(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if pref == nil {
		pref = &models.AppearancePreferences{
			ActiveTheme: "default",
			Mode:        "dark",
		}
	}
	return c.JSON(http.StatusOK, map[string]any{"appearance": pref})
}

func (h *AppearanceHandler) SetAppearance(c echo.Context) error {
	userID := handlerutils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req struct {
		Appearance *models.AppearancePreferences `json:"appearance"`
	}
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid appearance payload")
	}

	if req.Appearance != nil {
		if err := models.ValidateAppearance(req.Appearance); err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
	}

	if err := h.userRepo.SetAppearance(c.Request().Context(), userID, req.Appearance); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{"appearance": req.Appearance})
}
