package backups

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *Handler) ListS3Destinations(c echo.Context) error {
	list, err := h.backupService.ListS3Destinations(c.Request().Context())
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	for _, destination := range list {
		RedactS3Destination(destination)
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *Handler) CreateS3Destination(c echo.Context) error {
	var req S3DestinationPayload
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, ErrInvalidPayload.Error())
	}
	isDef := false
	if req.IsDefault != nil {
		isDef = *req.IsDefault
	}
	dest := S3Destination{
		Name:            req.Name,
		Description:     req.Description,
		Provider:        req.Provider,
		Endpoint:        req.Endpoint,
		Bucket:          req.Bucket,
		Region:          req.Region,
		PathPrefix:      req.PathPrefix,
		IsDefault:       isDef,
		AccessKeyID:     req.AccessKeyID,
		SecretAccessKey: req.SecretAccessKey,
	}
	if err := h.backupService.CreateS3Destination(c.Request().Context(), &dest); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.s3_destination_create", dest.ID, c.RealIP())
	}
	RedactS3Destination(&dest)
	return utils.Created(c, "S3 destination created successfully", dest)
}

func (h *Handler) VerifyS3Draft(c echo.Context) error {
	var req S3DestinationPayload
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, ErrInvalidPayload.Error())
	}
	dest := S3Destination{
		Name:            req.Name,
		Provider:        req.Provider,
		Endpoint:        req.Endpoint,
		Bucket:          req.Bucket,
		Region:          req.Region,
		AccessKeyID:     req.AccessKeyID,
		SecretAccessKey: req.SecretAccessKey,
	}
	verified, err := h.backupService.VerifyS3Draft(c.Request().Context(), &dest)
	if !verified || err != nil {
		errMsg := "verification failed"
		if err != nil {
			errMsg = err.Error()
		}
		return utils.Error(c, http.StatusBadRequest, errMsg)
	}
	return utils.Success(c, "S3 destination verified successfully", map[string]bool{"verified": true})
}

func (h *Handler) UpdateS3Destination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	dest, err := h.backupService.GetS3Destination(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || utils.IsNotFound(err) {
			return utils.Error(c, http.StatusNotFound, ErrS3DestinationNotFound.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	var req S3DestinationPayload
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, ErrInvalidPayload.Error())
	}
	if req.Name != "" {
		dest.Name = req.Name
	}
	if req.Description != "" {
		dest.Description = req.Description
	}
	if req.Provider != "" {
		dest.Provider = req.Provider
	}
	if req.Endpoint != "" {
		dest.Endpoint = req.Endpoint
	}
	if req.Bucket != "" {
		dest.Bucket = req.Bucket
	}
	if req.Region != "" {
		dest.Region = req.Region
	}
	if req.PathPrefix != "" {
		dest.PathPrefix = req.PathPrefix
	}
	if req.IsDefault != nil {
		dest.IsDefault = *req.IsDefault
	}
	if req.AccessKeyID != "" {
		dest.AccessKeyID = req.AccessKeyID
	}
	if req.SecretAccessKey != "" && req.SecretAccessKey != "********" {
		dest.SecretAccessKey = req.SecretAccessKey
	}
	if err := h.backupService.UpdateS3Destination(c.Request().Context(), dest); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.s3_destination_update", id, c.RealIP())
	}
	RedactS3Destination(dest)
	return utils.Success(c, "S3 destination updated successfully", dest)
}

func (h *Handler) VerifyS3Destination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	verified, err := h.backupService.VerifyS3Destination(c.Request().Context(), id)
	if !verified || err != nil {
		errMsg := "verification failed"
		if err != nil {
			errMsg = err.Error()
		}
		return utils.Error(c, http.StatusBadRequest, errMsg)
	}
	return utils.Success(c, "S3 destination verified successfully", map[string]bool{"verified": true})
}

func (h *Handler) SetDefaultS3Destination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	if err := h.backupService.SetDefaultS3Destination(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.s3_destination_default", id, c.RealIP())
	}
	return utils.Success(c, "Default S3 destination updated successfully", nil)
}

func (h *Handler) DeleteS3Destination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	if err := h.backupService.DeleteS3Destination(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.s3_destination_delete", id, c.RealIP())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) ListSFTPDestinations(c echo.Context) error {
	projectID := c.QueryParam("projectId")
	list, err := h.backupService.ListSFTPDestinations(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	for i := range list {
		RedactSFTPDestination(&list[i])
	}
	return utils.Success(c, "Operation successful", list)
}

func (h *Handler) CreateSFTPDestination(c echo.Context) error {
	var req SFTPDestinationPayload
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, ErrInvalidPayload.Error())
	}
	dest := SFTPDestination{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
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
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.sftp_destination_create", dest.ID, c.RealIP())
	}
	RedactSFTPDestination(&dest)
	return utils.Created(c, "SFTP destination created successfully", dest)
}

func (h *Handler) VerifySFTPDestination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	verified, err := h.backupService.VerifySFTPDestination(c.Request().Context(), id)
	if !verified || err != nil {
		errMsg := "verification failed"
		if err != nil {
			errMsg = err.Error()
		}
		return utils.Error(c, http.StatusBadRequest, errMsg)
	}
	return utils.Success(c, "SFTP destination verified successfully", map[string]bool{"verified": true})
}

func (h *Handler) DeleteSFTPDestination(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, ErrMissingID.Error())
	}
	if err := h.backupService.DeleteSFTPDestination(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if h.auditLogger != nil {
		h.auditLogger(c.Request().Context(), auditActorID(c), "backup.sftp_destination_delete", id, c.RealIP())
	}
	return c.NoContent(http.StatusNoContent)
}
