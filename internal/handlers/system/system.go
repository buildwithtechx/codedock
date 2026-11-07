package system

import (
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/models"
	"codedock/internal/repositories"
	authservices "codedock/internal/services/auth"
	systemservices "codedock/internal/services/system"
	"codedock/internal/utils"
)

type SystemHandler struct {
	service      *systemservices.SystemService
	userRepo     repositories.UserRepository
	tokenService *authservices.TokenService
}

func NewSystemHandler(s *systemservices.SystemService, userRepo repositories.UserRepository, tokenService *authservices.TokenService) *SystemHandler {
	return &SystemHandler{
		service:      s,
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

func (h *SystemHandler) GetStats(c echo.Context) error {
	stats, err := h.service.GetStats()
	if err != nil {
		return utils.Error(c, 500, "failed to get system stats")
	}
	return utils.Success(c, "System stats", stats)
}

func (h *SystemHandler) Restart(c echo.Context) error {
	go func() {
		if _, err := exec.LookPath("docker"); err == nil {
			exec.Command("docker", "restart", "codedock-control-plane").Start()
		} else {
			if p, err := os.FindProcess(os.Getpid()); err == nil {
				_ = p.Signal(syscall.SIGTERM)
			}
		}
	}()
	return utils.Success(c, "Restart initiated", map[string]string{"status": "restarting"})
}

func (h *SystemHandler) Cleanup(c echo.Context) error {
	if err := h.service.Cleanup(); err != nil {
		return utils.Error(c, 500, "Cleanup failed")
	}
	return utils.Success(c, "System cleanup completed successfully", nil)
}

func (h *SystemHandler) CreateLocalCLISession(c echo.Context) error {
	if config.Get().Cloud.Enabled {
		return utils.Error(c, http.StatusNotFound, "Not available in cloud mode")
	}

	internalToken := os.Getenv("CODEDOCK_INTERNAL_TOKEN")
	if internalToken == "" {
		internalToken = config.Get().Security.JWTSecret
	}
	reqToken := c.Request().Header.Get("X-Internal-Token")
	if internalToken == "" || reqToken == "" || reqToken != internalToken {
		return utils.Error(c, http.StatusUnauthorized, "The installation's operator credential is required")
	}

	if c.Request().Header.Get("Origin") != "" || c.Request().Header.Get("Sec-Fetch-Site") != "" {
		return utils.Error(c, http.StatusForbidden, "Local CLI authentication is not available to browsers")
	}

	if h.userRepo == nil || h.tokenService == nil {
		return utils.Error(c, http.StatusInternalServerError, "authentication service unavailable")
	}

	admin, err := h.userRepo.FindFoundingAdmin(c.Request().Context())
	if err != nil || admin == nil || admin.Role != models.UserRoleAdmin {
		return utils.Error(c, http.StatusConflict, "No local administrator is available. Complete Codedock setup first.")
	}

	token, err := h.tokenService.GenerateToken(admin)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to mint CLI session token")
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, map[string]string{
		"token":     token,
		"expiresAt": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
	})
}
