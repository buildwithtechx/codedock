package monitoring

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/utils"
)

type FoundingAdminFinder interface {
	FindFoundingAdmin(ctx context.Context) (id, email, role string, verified, totp bool, err error)
}

type TokenGenerator interface {
	GenerateCLIToken(id, email, role string, verified, totp bool) (string, error)
}

type SystemHandler struct {
	service     SystemService
	adminFinder FoundingAdminFinder
	tokenGen    TokenGenerator
}

func NewSystemHandler(s SystemService, adminFinder FoundingAdminFinder, tokenGen TokenGenerator) *SystemHandler {
	return &SystemHandler{
		service:     s,
		adminFinder: adminFinder,
		tokenGen:    tokenGen,
	}
}

func (h *SystemHandler) GetStats(c echo.Context) error {
	stats, err := h.service.GetStats(c.Request().Context())
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
	if err := h.service.Cleanup(c.Request().Context()); err != nil {
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

	if h.adminFinder == nil || h.tokenGen == nil {
		return utils.Error(c, http.StatusInternalServerError, "authentication service unavailable")
	}

	id, email, role, verified, totp, err := h.adminFinder.FindFoundingAdmin(c.Request().Context())
	if err != nil || id == "" || role != "admin" {
		return utils.Error(c, http.StatusConflict, "No local administrator is available. Complete Codedock setup first.")
	}

	token, err := h.tokenGen.GenerateCLIToken(id, email, role, verified, totp)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to mint CLI session token")
	}

	c.Response().Header().Set("Cache-Control", "no-store")
	return c.JSON(http.StatusOK, map[string]string{
		"token":     token,
		"expiresAt": time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339),
	})
}
