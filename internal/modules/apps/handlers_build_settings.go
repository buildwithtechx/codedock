package apps

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type UpdateBuildSettingsRequest struct {
	StaticOutputDir string `json:"staticOutputDir"`
	InstallCmd      string `json:"installCmd"`
	BuildCmd        string `json:"buildCmd"`
	StartCmd        string `json:"startCmd"`
	RuntimeVersion  string `json:"runtimeVersion"`
}

var runtimeVersionPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,31}$`)

func normalizeRuntimeVersion(raw string) (string, error) {
	version := strings.TrimSpace(raw)
	if version != "" && !runtimeVersionPattern.MatchString(version) {
		return "", errors.New("runtime version must look like 8.4 or 22.11.0")
	}
	return version, nil
}

func isValidHealthCheckPath(path string) bool {
	if path == "" {
		return true
	}
	if path[0] != '/' {
		return false
	}
	for _, ch := range path {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') || ch == '/' || ch == '-' || ch == '_' || ch == '.' || ch == '?' || ch == '=' || ch == '&') {
			return false
		}
	}
	return true
}

func (h *Handler) UpdateBuildSettings(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.service.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}

	var req UpdateBuildSettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	version, err := normalizeRuntimeVersion(req.RuntimeVersion)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	existing.StaticOutput = strings.TrimSpace(req.StaticOutputDir)
	existing.InstallCommand = strings.TrimSpace(req.InstallCmd)
	existing.BuildCommand = strings.TrimSpace(req.BuildCmd)
	existing.StartCommand = strings.TrimSpace(req.StartCmd)
	existing.RuntimeVersion = version

	if err := h.service.UpdateAppService(c.Request().Context(), existing); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Build settings updated successfully", existing)
}
