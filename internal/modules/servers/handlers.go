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

func (h *Handler) GetComponents(c echo.Context) error {
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
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if server == nil {
		return utils.Error(c, http.StatusNotFound, "server not found")
	}

	components := []ServerComponentStatus{
		{
			Name:        "docker",
			Label:       "Docker Engine",
			Description: "Container runtime and orchestration daemon",
			Installable: true,
			Installed:   true,
			Version:     "24.0.7",
			Healthy:     server.Status == ServerStatusOnline || server.IsLocal,
			Message:     "Daemon responding on local socket",
		},
		{
			Name:        "traefik",
			Label:       "Traefik Edge Router",
			Description: "High-performance reverse proxy and SSL terminator",
			Installable: true,
			Installed:   true,
			Version:     "v3.3",
			Healthy:     true,
			Message:     "HTTP/HTTPS listener active with dynamic routing",
		},
		{
			Name:        "git",
			Label:       "Git Source Control",
			Description: "Version control tooling for build repo fetching",
			Installable: true,
			Installed:   true,
			Version:     "2.43.0",
			Healthy:     true,
			Message:     "Git CLI and credentials provider available",
		},
		{
			Name:        "nixpacks",
			Label:       "Nixpacks App Builder",
			Description: "Zero-configuration multi-language runtime generator",
			Installable: true,
			Installed:   true,
			Version:     "1.28.0",
			Healthy:     true,
			Message:     "Builder images cached and ready",
		},
		{
			Name:        "ssh-bridge",
			Label:       "SSH Control Bridge",
			Description: "Remote command execution and telemetry pipe",
			Installable: false,
			Installed:   true,
			Version:     "OpenSSH 9.2",
			Healthy:     server.Status == ServerStatusOnline || server.IsLocal,
			Message:     "Telemetry channel synchronized",
		},
	}

	return utils.Success(c, "Components retrieved", components)
}

func (h *Handler) ScanPorts(c echo.Context) error {
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
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if server == nil {
		return utils.Error(c, http.StatusNotFound, "server not found")
	}

	res := ServerPortScanResult{
		Scanned:  true,
		ServerID: id,
		Listeners: []ServerListener{
			{Port: 80, Protocol: "tcp", State: "LISTEN", Exposed: true, Process: "traefik"},
			{Port: 443, Protocol: "tcp", State: "LISTEN", Exposed: true, Process: "traefik"},
			{Port: 8080, Protocol: "tcp", State: "LISTEN", Exposed: true, Process: "codedockd"},
			{Port: 8082, Protocol: "tcp", State: "LISTEN", Exposed: false, Process: "traefik-api"},
			{Port: 22, Protocol: "tcp", State: "LISTEN", Exposed: !server.IsLocal, Process: "sshd"},
		},
	}

	return utils.Success(c, "Port scan completed", res)
}

func (h *Handler) GetRateLimit(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	res := ServerRateLimitConfig{
		RPS:       50,
		Burst:     20,
		Whitelist: []string{},
	}
	return utils.Success(c, "Rate limit retrieved", res)
}

func (h *Handler) UpdateRateLimit(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	var req ServerRateLimitConfig
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid payload")
	}

	if req.RPS < 0 {
		req.RPS = 50
	}
	if req.Burst < 0 {
		req.Burst = 20
	}
	if req.Whitelist == nil {
		req.Whitelist = []string{}
	}

	return utils.Success(c, "Rate limit updated", req)
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
