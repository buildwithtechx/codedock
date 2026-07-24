package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/services"
	"codedock.run/codedock/internal/utils"
)

type ServerHandler struct {
	serverService services.ServerService
}

func NewServerHandler(serverService services.ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

type CreateServerRequest struct {
	Name      string `json:"name"`
	IPAddress string `json:"ipAddress"`
}

func (h *ServerHandler) Create(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req CreateServerRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	if req.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "server name is required")
	}

	server, err := h.serverService.CreateServer(c.Request().Context(), userClaims.UserID, req.Name, req.IPAddress)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Server created", server)
}

func (h *ServerHandler) List(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	servers, err := h.serverService.ListServersByUser(c.Request().Context(), userClaims.UserID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Operation successful", servers)
}
