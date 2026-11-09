package dns

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type IngressService interface {
	List(ctx context.Context, serverID string) (any, error)
	Discover(ctx context.Context, req IngressDiscoverRequest) (*IngressDiscoveryResult, error)
	Connect(ctx context.Context, req IngressConnectRequest) (*IngressConnection, error)
	Update(ctx context.Context, id string, req IngressUpdateRequest) (*IngressConnection, error)
	Retry(ctx context.Context, id string) (*IngressConnection, error)
	Disconnect(ctx context.Context, id string) error
	GetDomain(ctx context.Context, domainID string) (any, error)
	SetDomain(ctx context.Context, domainID string, req IngressDomainRequest) (any, error)
}

type IngressHandler struct {
	service IngressService
}

func NewIngressHandler(service ...IngressService) *IngressHandler {
	if len(service) > 0 && service[0] != nil {
		return &IngressHandler{service: service[0]}
	}
	return &IngressHandler{service: newInMemoryIngressService()}
}

func (h *IngressHandler) List(c echo.Context) error {
	serverID := c.QueryParam("serverId")
	conns, err := h.service.List(c.Request().Context(), serverID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, "failed to list ingress connections: "+err.Error())
	}
	return utils.Success(c, "Ingress connections retrieved", conns)
}

func (h *IngressHandler) Discover(c echo.Context) error {
	var req IngressDiscoverRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid discovery request")
	}
	res, err := h.service.Discover(c.Request().Context(), req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Ingress discovery completed", res)
}

func (h *IngressHandler) Connect(c echo.Context) error {
	var req IngressConnectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ingress connect payload")
	}
	if req.ServerID == "" || req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "serverId and name are required")
	}
	conn, err := h.service.Connect(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save ingress connection: "+err.Error())
	}
	return utils.Success(c, "Ingress tunnel connected", conn)
}

func (h *IngressHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req IngressUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid update payload")
	}
	conn, err := h.service.Update(c.Request().Context(), id, req)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Ingress connection updated", conn)
}

func (h *IngressHandler) Retry(c echo.Context) error {
	id := c.Param("id")
	conn, err := h.service.Retry(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Ingress connection retry initiated", conn)
}

func (h *IngressHandler) Disconnect(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.Disconnect(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Ingress tunnel disconnected", map[string]string{"id": id})
}

func (h *IngressHandler) GetDomain(c echo.Context) error {
	domainID := c.Param("domainId")
	mapping, err := h.service.GetDomain(c.Request().Context(), domainID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Domain ingress retrieved", mapping)
}

func (h *IngressHandler) SetDomain(c echo.Context) error {
	domainID := c.Param("domainId")
	var req IngressDomainRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ingress domain payload")
	}
	mapping, err := h.service.SetDomain(c.Request().Context(), domainID, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Domain ingress configured", mapping)
}
