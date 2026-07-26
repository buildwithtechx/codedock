package handlers

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

func (h *AppHandler) CreateLogDrain(c echo.Context) error {
	id := c.Param("id")
	var req models.CreateLogDrainRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}

	if req.EndpointURL != "" {
		if err := validateDrainURL(req.EndpointURL); err != nil {
			return utils.Error(c, http.StatusBadRequest, err.Error())
		}
	}

	drain := &models.LogDrain{
		ServiceID:   existing.ID,
		ProjectID:   existing.ProjectID,
		DrainType:   req.DrainType,
		EndpointURL: req.EndpointURL,
		AuthToken:   req.AuthToken,
	}
	created, err := h.appService.CreateLogDrain(c.Request().Context(), drain)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	go func() {
		_ = h.appService.UpdateAppService(context.Background(), existing)
	}()
	created.AuthToken = ""
	return utils.Created(c, "Log drain created successfully", created)
}

func (h *AppHandler) ListLogDrains(c echo.Context) error {
	id := c.Param("id")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	drains, err := h.appService.ListLogDrains(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	for _, drain := range drains {
		drain.AuthToken = ""
	}
	return utils.Success(c, "Operation successful", drains)
}

func (h *AppHandler) DeleteLogDrain(c echo.Context) error {
	id := c.Param("id")
	drainID := c.Param("drainId")
	existing, err := h.appService.GetAppService(c.Request().Context(), id)
	if err != nil || existing == nil {
		return utils.Error(c, http.StatusNotFound, "app service not found")
	}
	if err := h.verifyProjectOwnership(c, existing.ProjectID); err != nil {
		return err
	}
	if err := h.appService.DeleteLogDrain(c.Request().Context(), drainID, id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	go func() {
		_ = h.appService.UpdateAppService(context.Background(), existing)
	}()

	return c.NoContent(http.StatusNoContent)
}

func validateDrainURL(u string) error {
	parsed, err := url.Parse(u)
	if err != nil {
		return errors.New("invalid url format")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return errors.New("url must use http or https")
	}
	host := parsed.Hostname()
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
			return errors.New("internal or private IPs are not allowed for log drains")
		}
	}
	return nil
}
