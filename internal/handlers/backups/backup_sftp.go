package backups

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/utils"
)

type SFTPDestinationPayload struct {
	ProjectID   string `json:"projectId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PrivateKey  string `json:"privateKey"`
	PathPrefix  string `json:"pathPrefix"`
}

func redactSFTPDestination(destination *models.SFTPDestination) {
	if destination != nil {
		destination.Password = "********"
		destination.PrivateKey = "********"
	}
}

func (h *BackupHandler) ListSFTPDestinations(c echo.Context) error {
	list, err := h.backupService.ListSFTPDestinations(c.Request().Context(), c.QueryParam("projectId"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	for i := range list {
		redactSFTPDestination(&list[i])
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *BackupHandler) CreateSFTPDestination(c echo.Context) error {
	var req SFTPDestinationPayload
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	dest := models.SFTPDestination{
		ProjectID:   req.ProjectID,
		Name:        req.Name,
		Description: req.Description,
		Host:        req.Host,
		Port:        req.Port,
		Username:    req.Username,
		Password:    req.Password,
		PrivateKey:  req.PrivateKey,
		PathPrefix:  req.PathPrefix,
	}
	if err := h.backupService.CreateSFTPDestination(c.Request().Context(), &dest); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	redactSFTPDestination(&dest)
	return utils.Created(c, "Created successfully", dest)
}

func (h *BackupHandler) DeleteSFTPDestination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id")
	}
	if err := h.backupService.DeleteSFTPDestination(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *BackupHandler) VerifySFTPDestination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id")
	}
	if err := h.backupService.VerifySFTPDestination(c.Request().Context(), id); err != nil {
		return utils.Success(c, "Verification failed", map[string]any{"ok": false, "reason": err.Error()})
	}
	return utils.Success(c, "Verification succeeded", map[string]any{"ok": true})
}
