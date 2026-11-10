package chatapps

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

// Handler handles HTTP requests for Chat App connections, identities, pairing, and webhooks
type Handler struct {
	service    *Service
	dispatcher *Dispatcher
}

// NewHandler creates a new chat apps HTTP handler
func NewHandler(service *Service, dispatcher *Dispatcher) *Handler {
	return &Handler{
		service:    service,
		dispatcher: dispatcher,
	}
}

// List returns all configured chat app connections for the authenticated user
func (h *Handler) List(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	conns, err := h.service.ListConnections(c.Request().Context(), userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Chat connections retrieved", conns)
}

// Create provisions a new chat app connection with encrypted credentials
func (h *Handler) Create(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req CreateChatAppRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	conn, err := h.service.CreateConnection(c.Request().Context(), userID, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Chat connection created", conn)
}

// Get returns details for a specific chat app connection
func (h *Handler) Get(c echo.Context) error {
	id := c.Param("id")
	conn, err := h.service.GetConnection(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Chat connection retrieved", conn)
}

// Update modifies connection settings or credentials
func (h *Handler) Update(c echo.Context) error {
	id := c.Param("id")
	var req UpdateChatAppRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request body")
	}
	conn, err := h.service.UpdateConnection(c.Request().Context(), id, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Chat connection updated", conn)
}

// Delete removes a chat app connection and revokes credentials
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteConnection(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Chat connection deleted", nil)
}

// Test verifies connection health and bot API connectivity
func (h *Handler) Test(c echo.Context) error {
	id := c.Param("id")
	res, err := h.service.TestConnection(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Connection test completed", res)
}
