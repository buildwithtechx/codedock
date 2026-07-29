package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	systemservices "codedock.run/codedock/internal/services/system"
	"codedock.run/codedock/internal/utils"
)

type ServerHandler struct {
	serverService systemservices.ServerService
}

func NewServerHandler(serverService systemservices.ServerService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
	}
}

func (h *ServerHandler) TestSSH(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req models.TestSSHRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	if err := h.serverService.TestSSH(c.Request().Context(), req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}

	return utils.Success(c, "SSH connection successful", nil)
}

func (h *ServerHandler) Create(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req models.CreateServerRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	if req.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "server name is required")
	}

	server, err := h.serverService.CreateServer(c.Request().Context(), userClaims.UserID, req)
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

	for _, s := range servers {
		s.WorkerToken = "********"
		s.SSHPassword = "********"
		if s.SSHKey != "" {
			s.SSHKey = "********"
		}
	}

	return utils.Success(c, "Operation successful", servers)
}

func (h *ServerHandler) Delete(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	if err := h.serverService.DeleteServer(c.Request().Context(), id); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	return utils.Success(c, "Server deleted", nil)
}
