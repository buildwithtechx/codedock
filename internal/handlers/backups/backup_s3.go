package backups

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

type CreateS3DestinationRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Provider        string `json:"provider"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type UpdateS3DestinationRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Provider        string `json:"provider"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

func (h *BackupHandler) ListS3Destinations(c echo.Context) error {
	list, err := h.backupService.ListS3Destinations(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	for _, destination := range list {
		redactS3Destination(destination)
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *BackupHandler) CreateS3Destination(c echo.Context) error {
	var req CreateS3DestinationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	dest := models.S3Destination{
		Name:            req.Name,
		Description:     req.Description,
		Provider:        req.Provider,
		Endpoint:        req.Endpoint,
		Bucket:          req.Bucket,
		Region:          req.Region,
		AccessKeyID:     req.AccessKeyID,
		SecretAccessKey: req.SecretAccessKey,
	}
	if err := h.backupService.CreateS3Destination(c.Request().Context(), &dest); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	redactS3Destination(&dest)
	return utils.Created(c, "Created successfully", dest)
}

func (h *BackupHandler) UpdateS3Destination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id")
	}
	var req UpdateS3DestinationRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	dest := models.S3Destination{
		ID:              id,
		Name:            req.Name,
		Description:     req.Description,
		Provider:        req.Provider,
		Endpoint:        req.Endpoint,
		Bucket:          req.Bucket,
		Region:          req.Region,
		AccessKeyID:     req.AccessKeyID,
		SecretAccessKey: req.SecretAccessKey,
	}
	if err := h.backupService.UpdateS3Destination(c.Request().Context(), &dest); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	redactS3Destination(&dest)
	return utils.Success(c, "Updated successfully", dest)
}

func redactS3Destination(destination *models.S3Destination) {
	if destination != nil {
		destination.SecretAccessKey = "********"
	}
}

func (h *BackupHandler) DeleteS3Destination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id")
	}
	if err := h.backupService.DeleteS3Destination(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
