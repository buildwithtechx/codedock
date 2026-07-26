package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type ArchiveHandler struct {
	service        *services.ArchiveService
	projectService *services.ProjectService
}

func NewArchiveHandler(s *services.ArchiveService, ps *services.ProjectService) *ArchiveHandler {
	return &ArchiveHandler{service: s, projectService: ps}
}

func (h *ArchiveHandler) DeployArchive(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	projectID := c.FormValue("projectId")
	if projectID == "" {
		projectID = c.FormValue("project_id")
	}
	if projectID == "" {
		return utils.Error(c, http.StatusBadRequest, "projectId is required")
	}
	if user.Role != "admin" {
		if !h.projectService.HasPermission(c.Request().Context(), projectID, user.UserID, models.UserRole(user.Role), "") {
			return utils.Error(c, http.StatusForbidden, "insufficient permissions for this project")
		}
	}
	appName := c.FormValue("name")

	file, err := c.FormFile("file")
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, "archive file is required")
	}

	if appName == "" {
		base := strings.TrimSuffix(file.Filename, ".tar.gz")
		base = strings.TrimSuffix(base, ".tar")
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

	tmpPath := filepath.Join(os.TempDir(), "codedock-upload", uuid.New().String()+".tar.gz")
	if err := writeFile(tmpPath, io.LimitReader(src, maxArchiveBytes)); err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to save archive")
	}
	defer os.Remove(tmpPath)

	result, err := h.service.Deploy(c.Request().Context(), projectID, appName, tmpPath)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "Archive deployed", result)
}

func writeFile(path string, r io.Reader) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
