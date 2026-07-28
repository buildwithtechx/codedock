package auth

import (
	"os"

	"github.com/labstack/echo/v4"

	authservices "codedock.run/codedock/internal/services/auth"
	"codedock.run/codedock/internal/utils"
)

type OnboardingHandler struct {
	userService *authservices.UserService
}

func NewOnboardingHandler(userService *authservices.UserService) *OnboardingHandler {
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
	return utils.Success(c, "Setup status", map[string]any{
		"setupRequired": count == 0,
		"cwd":           cwd,
	})
}
