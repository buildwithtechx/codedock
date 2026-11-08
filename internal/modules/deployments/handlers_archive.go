package deployments

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ArchiveHandler struct {
	archiveService *ArchiveService
	projectLookup  ProjectLookup
}

func NewArchiveHandler(s *ArchiveService, pl ProjectLookup) *ArchiveHandler {
	return &ArchiveHandler{
		archiveService: s,
		projectLookup:  pl,
	}
}

func (h *ArchiveHandler) DeployArchive(c echo.Context) error {
	userID, role, ok := utils.ExtractUser(c)
	if !ok || userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.FormValue("projectId")
	if projectID == "" {
		projectID = c.FormValue("project_id")
	}
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "projectId is required")
	}

	if role != "admin" && role != "owner" && h.projectLookup != nil {
		if !h.projectLookup.HasPermission(c.Request().Context(), projectID, userID, role, "") {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions for this project")
		}
	}

	file, err := c.FormFile("file")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "archive file is required")
	}

	appName := c.FormValue("name")
	if appName == "" {
		base := strings.TrimSuffix(file.Filename, ".tar.gz")
		base = strings.TrimSuffix(base, ".tar")
		base = strings.TrimSuffix(base, ".zip")
		appName = base
	}

	src, err := file.Open()
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to read uploaded file")
	}
	defer src.Close()

	const maxArchiveBytes = 500 * 1024 * 1024
	if file.Size > maxArchiveBytes {
		return utils.Error(c, http.StatusRequestEntityTooLarge, "archive exceeds 500 MB limit")
	}

	tmpDir := filepath.Join(os.TempDir(), "codedock-upload")
	_ = os.MkdirAll(tmpDir, 0o755)
	tmpPath := filepath.Join(tmpDir, uuid.New().String()+".tar.gz")

	dst, err := os.Create(tmpPath)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to save archive")
	}
	_, err = io.Copy(dst, io.LimitReader(src, maxArchiveBytes))
	dst.Close()
	if err != nil {
		_ = os.Remove(tmpPath)
		return utils.Error(c, http.StatusInternalServerError, "failed to write archive")
	}
	defer os.Remove(tmpPath)

	serviceID := c.FormValue("serviceId")
	result, err := h.archiveService.Deploy(c.Request().Context(), serviceID, projectID, tmpPath)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Archive deployed", result)
}
