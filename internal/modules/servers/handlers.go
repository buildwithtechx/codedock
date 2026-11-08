package servers

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type AuditLogger interface {
	Log(ctx context.Context, userID, action, resource, ip string)
}

type Handler struct {
	service Service
	audit   AuditLogger
}

func NewHandler(service Service, audit AuditLogger) *Handler {
	return &Handler{service: service, audit: audit}
}

func (h *Handler) List(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	search := c.QueryParam("search")
	servers, err := h.service.ListServersByUser(c.Request().Context(), userID, search)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	for _, s := range servers {
		h.redactCredentials(s)
	}

	if servers == nil {
		servers = make([]*Server, 0)
	}
	return utils.Success(c, "Servers retrieved", servers)
}

func (h *Handler) Get(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	server, err := h.service.GetServer(c.Request().Context(), id, userID)
	if err != nil {
		if utils.IsForbidden(err) || err == ErrServerUnauthorized {
			return utils.Error(c, http.StatusForbidden, "server access denied")
		}
		if utils.IsNotFound(err) || err == ErrServerNotFound {
			return utils.Error(c, http.StatusNotFound, "server not found")
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if server == nil {
		return utils.Error(c, http.StatusNotFound, "server not found")
	}

	h.redactCredentials(server)
	return utils.Success(c, "Server retrieved", server)
}

func (h *Handler) Create(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req CreateServerRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	if req.Name == "" {
		return utils.Error(c, http.StatusBadRequest, "server name is required")
	}

	server, err := h.service.CreateServer(c.Request().Context(), userID, req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if h.audit != nil {
		h.audit.Log(c.Request().Context(), userID, "server.create", server.ID, c.RealIP())
	}

	h.redactCredentials(server)
	return utils.Success(c, "Server created", server)
}

func (h *Handler) Update(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	var req UpdateServerRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	server, err := h.service.UpdateServer(c.Request().Context(), id, userID, req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if h.audit != nil {
		h.audit.Log(c.Request().Context(), userID, "server.update", server.ID, c.RealIP())
	}

	h.redactCredentials(server)
	return utils.Success(c, "Server updated", server)
}

func (h *Handler) Delete(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	if err := h.service.DeleteServer(c.Request().Context(), id, userID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	if h.audit != nil {
		h.audit.Log(c.Request().Context(), userID, "server.delete", id, c.RealIP())
	}
	return utils.Success(c, "Server deleted", nil)
}

func (h *Handler) TestSSH(c echo.Context) error {
	var req TestSSHRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	if err := h.service.TestSSH(c.Request().Context(), req); err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "SSH connection successful", nil)
}

func (h *Handler) redactCredentials(s *Server) {
	if s == nil {
		return
	}
	s.WorkerToken = "********"
	s.SSHPassword = "********"
	s.SSHPrivateKey = "********"
	if s.SSHKey != "" {
		s.SSHKey = "********"
	}
}
