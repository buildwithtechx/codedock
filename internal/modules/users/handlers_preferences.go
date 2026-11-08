package users

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AppearanceHandler struct {
	userRepo UserRepository
}

func NewAppearanceHandler(userRepo UserRepository) *AppearanceHandler {
	return &AppearanceHandler{userRepo: userRepo}
}

func (h *AppearanceHandler) GetAppearance(c echo.Context) error {
	userID := extractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	pref, err := h.userRepo.GetAppearance(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if pref == nil {
		pref = &AppearancePreferences{
			ActiveTheme: "default",
			Mode:        "dark",
		}
	}
	return c.JSON(http.StatusOK, AppearanceResponse{Appearance: pref})
}

func (h *AppearanceHandler) SetAppearance(c echo.Context) error {
	userID := extractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req SetAppearanceRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid appearance payload")
	}

	if req.Appearance != nil {
		if err := ValidateAppearance(req.Appearance); err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
	}

	if err := h.userRepo.SetAppearance(c.Request().Context(), userID, req.Appearance); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, AppearanceResponse{Appearance: req.Appearance})
}

type OnboardingHandler struct {
	userService *UserService
}

func NewOnboardingHandler(userService *UserService) *OnboardingHandler {
	return &OnboardingHandler{
		userService: userService,
	}
}

func (h *OnboardingHandler) SetupStatus(c echo.Context) error {
	count, err := h.userService.CountUsers(c.Request().Context())
	if err != nil {
		return utils.Error(c, 500, "failed to check user count")
	}
	cwd, _ := os.Getwd()
	return utils.Success(c, "Setup status", SetupStatusResponse{
		SetupRequired: count == 0,
		Cwd:           cwd,
	})
}
