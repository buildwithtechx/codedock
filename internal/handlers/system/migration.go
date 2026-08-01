package system

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

type MigrationHandler struct {
	service     *systemservices.MigrationService
	userCounter interface {
		CountUsers(context.Context) (int, error)
	}
}

func NewMigrationHandler(s *systemservices.MigrationService, userCounter interface {
	CountUsers(context.Context) (int, error)
}) *MigrationHandler {
	return &MigrationHandler{service: s, userCounter: userCounter}
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
	return h.importBundle(c)
}

func (h *MigrationHandler) ImportDuringSetup(c echo.Context) error {
	count, err := h.userCounter.CountUsers(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to check user count")
	}
	if count != 0 {
		return utils.Error(c, http.StatusForbidden, "setup import is only available before the first account is created")
	}
	return h.importBundle(c)
}

func (h *MigrationHandler) importBundle(c echo.Context) error {
	passphrase := c.FormValue("passphrase")
	if passphrase == "" {
		return utils.Error(c, http.StatusBadRequest, "passphrase form value is required")
	}

	file, err := c.FormFile("bundle")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "bundle file is required (multipart field: bundle)")
	}

	src, err := file.Open()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to open uploaded bundle")
	}
	defer src.Close()

	const maxBundleUploadBytes = 500 * 1024 * 1024
	bundleData, err := io.ReadAll(io.LimitReader(src, maxBundleUploadBytes+1))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to read bundle data")
	}
	if len(bundleData) > maxBundleUploadBytes {
		return utils.Error(c, http.StatusRequestEntityTooLarge, "migration bundle file exceeds maximum allowed size (500MB)")
	}

	manifest, err := h.service.Import(c.Request().Context(), bundleData, passphrase)
	if err != nil {
		return utils.Error(c, http.StatusUnprocessableEntity, fmt.Sprintf("import failed: %v", err))
	}

	return utils.Success(c, "Import completed successfully", manifest)
}
