package webhooks

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type WebhookHandler struct {
	service WebhookService
}

func NewWebhookHandler(s WebhookService) *WebhookHandler {
	return &WebhookHandler{service: s}
}

func (h *WebhookHandler) ListWebhooks(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	webhooks, err := h.service.ListWebhooks(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", webhooks)
}

func (h *WebhookHandler) CreateWebhook(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	var req CreateWebhookRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.URL == "" {
		return utils.Error(c, http.StatusBadRequest, "url is required")
	}

	w := &Webhook{
		ServiceID: serviceID,
		URL:       req.URL,
		Secret:    req.Secret,
		Events:    req.Events,
	}
	created, err := h.service.CreateWebhook(c.Request().Context(), w)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *WebhookHandler) DeleteWebhook(c echo.Context) error {
	serviceID := c.Param("id")
	id := c.Param("webhookId")
	if serviceID == "" || id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing parameters")
	}

	if err := h.service.DeleteWebhook(c.Request().Context(), id, serviceID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *WebhookHandler) ListLogDrains(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	drains, err := h.service.ListLogDrains(c.Request().Context(), serviceID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Operation successful", drains)
}

func (h *WebhookHandler) CreateLogDrain(c echo.Context) error {
	serviceID := c.Param("id")
	if serviceID == "" {
		return utils.Error(c, http.StatusBadRequest, "missing service id")
	}

	var req CreateLogDrainRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}
	if req.EndpointURL == "" || req.DrainType == "" {
		return utils.Error(c, http.StatusBadRequest, "endpointUrl and drainType are required")
	}

	d := &LogDrain{
		ProjectID:   req.ProjectID,
		ServiceID:   serviceID,
		DrainType:   req.DrainType,
		EndpointURL: req.EndpointURL,
		APIKey:      req.APIKey,
	}
	created, err := h.service.CreateLogDrain(c.Request().Context(), d)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Created(c, "Created successfully", created)
}

func (h *WebhookHandler) DeleteLogDrain(c echo.Context) error {
	serviceID := c.Param("id")
	id := c.Param("drainId")
	if serviceID == "" || id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing parameters")
	}

	if err := h.service.DeleteLogDrain(c.Request().Context(), id, serviceID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
