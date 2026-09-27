package system

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/repositories"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

const (
	maxSetupBundleUploadBytes = 32 * 1024 * 1024  // 32MB pre-auth setup upload limit
	maxBundleUploadBytes      = 500 * 1024 * 1024 // 500MB authenticated upload limit
)

type MigrationHandler struct {
	service     *systemservices.MigrationService
	userCounter interface {
		CountUsers(context.Context) (int, error)
	}
	setupLimiter *middleware.RateLimiter
}

func NewMigrationHandler(s *systemservices.MigrationService, userCounter interface {
	CountUsers(context.Context) (int, error)
}) *MigrationHandler {
	return &MigrationHandler{
		service:      s,
		userCounter:  userCounter,
		setupLimiter: middleware.NewRateLimiter(5, time.Minute),
	}
}

func (h *MigrationHandler) Export(c echo.Context) error {
	var req struct {
		Passphrase string `json:"passphrase"`
	}
	if err := c.Bind(&req); err != nil || req.Passphrase == "" {
		return utils.Error(c, http.StatusBadRequest, "passphrase is required in request body")
	}

	bundleData, err := h.service.Export(c.Request().Context(), req.Passphrase)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, fmt.Sprintf("export failed: %v", err))
	}

	filename := fmt.Sprintf("codedock-bundle-%s.codedock", time.Now().UTC().Format("20060102-150405"))
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().WriteHeader(http.StatusOK)
	_, _ = c.Response().Write(bundleData)
	return nil
}

func (h *MigrationHandler) Import(c echo.Context) error {
	return h.importBundle(c, maxBundleUploadBytes)
}

func (h *MigrationHandler) ImportDuringSetup(c echo.Context) error {
	if h.setupLimiter != nil {
		key := c.RealIP()
		if !h.setupLimiter.Allow(key) {
			return c.JSON(http.StatusTooManyRequests, utils.RateLimitError{
				Message:    "setup import rate limit exceeded",
				RetryAfter: 60,
			})
		}
	}

	c.Request().Body = http.MaxBytesReader(c.Response().Writer, c.Request().Body, maxSetupBundleUploadBytes)

	repositories.FirstUserSetupLock.Lock()
	defer repositories.FirstUserSetupLock.Unlock()

	count, err := h.userCounter.CountUsers(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to check user count")
	}
	if count != 0 {
		return utils.Error(c, http.StatusForbidden, "setup import is only available before the first account is created")
	}

	return h.importBundle(c, maxSetupBundleUploadBytes)
}

func (h *MigrationHandler) importBundle(c echo.Context, maxAllowedBytes int64) error {
	passphrase := c.FormValue("passphrase")
	if passphrase == "" {
		return utils.Error(c, http.StatusBadRequest, "passphrase form value is required")
	}

	file, err := c.FormFile("bundle")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "bundle file is required (multipart field: bundle)")
	}

	if file.Size > maxAllowedBytes {
		return utils.Error(c, http.StatusRequestEntityTooLarge, fmt.Sprintf("migration bundle file exceeds maximum allowed size (%dMB)", maxAllowedBytes/(1024*1024)))
	}

	src, err := file.Open()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to open uploaded bundle")
	}
	defer src.Close()

	bundleData, err := io.ReadAll(io.LimitReader(src, maxAllowedBytes+1))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to read bundle data")
	}
	if int64(len(bundleData)) > maxAllowedBytes {
		return utils.Error(c, http.StatusRequestEntityTooLarge, fmt.Sprintf("migration bundle file exceeds maximum allowed size (%dMB)", maxAllowedBytes/(1024*1024)))
	}

	manifest, err := h.service.Import(c.Request().Context(), bundleData, passphrase)
	if err != nil {
		return utils.Error(c, http.StatusUnprocessableEntity, fmt.Sprintf("import failed: %v", err))
	}

	return utils.Success(c, "Import completed successfully", manifest)
}
