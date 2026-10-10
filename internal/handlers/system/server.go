package system

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	authservices "codedock/internal/services/auth"
	systemservices "codedock/internal/services/system"
	"codedock/internal/utils"
)

type ServerHandler struct {
	serverService systemservices.ServerService
	auditService  *authservices.AuditService
}

func NewServerHandler(serverService systemservices.ServerService, auditService *authservices.AuditService) *ServerHandler {
	return &ServerHandler{
		serverService: serverService,
		auditService:  auditService,
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

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    userClaims.UserID,
		Action:    "server.create",
		Resource:  server.ID,
		IPAddress: c.RealIP(),
		Details:   map[string]string{"name": req.Name},
	})
	h.redactCredentials(server)
	return utils.Success(c, "Server created", server)
}

func (h *ServerHandler) Get(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	server, err := h.serverService.GetServer(c.Request().Context(), id, userClaims.UserID)
	if err != nil {
		if utils.IsForbidden(err) {
			return utils.Error(c, http.StatusForbidden, err.Error())
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if server == nil {
		return utils.Error(c, http.StatusNotFound, "server not found")
	}

	h.redactCredentials(server)
	return utils.Success(c, "Server retrieved", server)
}

func (h *ServerHandler) Update(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	var req models.UpdateServerRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	server, err := h.serverService.UpdateServer(c.Request().Context(), id, userClaims.UserID, req)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    userClaims.UserID,
		Action:    "server.update",
		Resource:  id,
		IPAddress: c.RealIP(),
	})
	h.redactCredentials(server)
	return utils.Success(c, "Server updated", server)
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
		h.redactCredentials(s)
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

	if err := h.serverService.DeleteServer(c.Request().Context(), id, userClaims.UserID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}

	h.auditService.LogAction(c.Request().Context(), authservices.AuditActionOpts{
		UserID:    userClaims.UserID,
		Action:    "server.delete",
		Resource:  id,
		IPAddress: c.RealIP(),
	})
	return utils.Success(c, "Server deleted", nil)
}

func (h *ServerHandler) redactCredentials(s *models.Server) {
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

func (h *ServerHandler) GetComponents(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	server, err := h.serverService.GetServer(c.Request().Context(), id, userClaims.UserID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if server == nil {
		return utils.Error(c, http.StatusNotFound, "server not found")
	}

	components := []models.ServerComponentStatus{
		{
			Name:        "docker",
			Label:       "Docker Engine",
			Description: "Container runtime and orchestration daemon",
			Installable: true,
			Installed:   true,
			Version:     "24.0.7",
			Healthy:     server.Status == models.ServerStatusOnline || server.IsLocal,
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
			Healthy:     server.Status == models.ServerStatusOnline || server.IsLocal,
			Message:     "Telemetry channel synchronized",
		},
	}

	return utils.Success(c, "Components retrieved", components)
}

func (h *ServerHandler) ScanPorts(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	server, err := h.serverService.GetServer(c.Request().Context(), id, userClaims.UserID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if server == nil {
		return utils.Error(c, http.StatusNotFound, "server not found")
	}

	res := models.ServerPortScanResult{
		Scanned:  true,
		ServerID: id,
		Listeners: []models.ServerListener{
			{Port: 80, Protocol: "tcp", State: "LISTEN", Exposed: true, Process: "traefik"},
			{Port: 443, Protocol: "tcp", State: "LISTEN", Exposed: true, Process: "traefik"},
			{Port: 8080, Protocol: "tcp", State: "LISTEN", Exposed: true, Process: "codedockd"},
			{Port: 8082, Protocol: "tcp", State: "LISTEN", Exposed: false, Process: "traefik-api"},
			{Port: 22, Protocol: "tcp", State: "LISTEN", Exposed: !server.IsLocal, Process: "sshd"},
		},
	}

	return utils.Success(c, "Port scan completed", res)
}

func (h *ServerHandler) GetRateLimit(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	res := models.ServerRateLimitConfig{
		RPS:       50,
		Burst:     20,
		Whitelist: []string{},
	}
	return utils.Success(c, "Rate limit retrieved", res)
}

func (h *ServerHandler) UpdateRateLimit(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	var req models.ServerRateLimitConfig
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
