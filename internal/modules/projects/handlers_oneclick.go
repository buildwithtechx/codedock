package projects

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type OneClickService interface {
	ListApps() any
	GetApp(id string) (any, error)
	AddCustomApp(tmpl any) (any, error)
	DeleteCustomApp(id string) error
	Review(ctx context.Context, req any) (any, error)
	Deploy(ctx context.Context, req any) (any, error)
}

type OneClickHandler struct {
	service OneClickService
}

func NewOneClickHandler(s OneClickService) *OneClickHandler {
	return &OneClickHandler{service: s}
}

func (h *OneClickHandler) List(c echo.Context) error {
	if h.service == nil {
		return utils.Success(c, "Available one-click apps", []any{})
	}
	return utils.Success(c, "Available one-click apps", h.service.ListApps())
}

func (h *OneClickHandler) Get(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusNotFound, "one-click app not found")
	}
	app, err := h.service.GetApp(c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "One-click app", app)
}

func (h *OneClickHandler) AddCustom(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "service unavailable")
	}
	var tmpl any
	if err := c.Bind(&tmpl); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid template json: "+err.Error())
	}
	app, err := h.service.AddCustomApp(tmpl)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Custom app added", app)
}

func (h *OneClickHandler) DeleteCustom(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "service unavailable")
	}
	if err := h.service.DeleteCustomApp(c.Param("id")); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Custom app deleted", nil)
}

func (h *OneClickHandler) Review(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "service unavailable")
	}
	var req any
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request: "+err.Error())
	}
	result, err := h.service.Review(c.Request().Context(), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Stack review generated", result)
}

func (h *OneClickHandler) Deploy(c echo.Context) error {
	if h.service == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "service unavailable")
	}
	var req any
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload: "+err.Error())
	}
	created, err := h.service.Deploy(c.Request().Context(), req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Accepted(c, "One-click deployment started", created)
}
