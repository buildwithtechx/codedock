package system

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"codedock/internal/models"
	"codedock/internal/utils"

	"github.com/labstack/echo/v4"
)

type IngressHandler struct {
	mu          sync.RWMutex
	connections map[string]*models.IngressConnection
	domains     map[string]*models.IngressHostname
}

func NewIngressHandler() *IngressHandler {
	return &IngressHandler{
		connections: make(map[string]*models.IngressConnection),
		domains:     make(map[string]*models.IngressHostname),
	}
}

func generateID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}

func (h *IngressHandler) List(c echo.Context) error {
	serverId := c.QueryParam("serverId")
	h.mu.RLock()
	defer h.mu.RUnlock()

	var result []*models.IngressConnection
	for _, conn := range h.connections {
		if serverId == "" || conn.ServerID == serverId {
			result = append(result, conn)
		}
	}
	if result == nil {
		result = make([]*models.IngressConnection, 0)
	}
	return utils.Success(c, "Ingress connections retrieved", result)
}

func (h *IngressHandler) Discover(c echo.Context) error {
	var req models.IngressDiscoverRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid discovery request")
	}

	result := models.IngressDiscoveryResult{
		Accounts: []models.IngressDiscoveryAccount{
			{ID: "cf-acc-primary", Name: "Primary Cloudflare Account"},
		},
		Tunnels: []models.IngressDiscoveryTunnel{
			{ID: "cf-tun-edge-1", Name: "codedock-edge-tunnel", Remote: true, Status: "healthy"},
		},
		Policies: []models.IngressDiscoveryPolicy{
			{ID: "cf-pol-bypass", Name: "Allow Team SSO", Decision: "allow"},
		},
	}
	return utils.Success(c, "Ingress discovery completed", result)
}

func (h *IngressHandler) Connect(c echo.Context) error {
	var req models.IngressConnectRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid ingress connect payload")
	}
	if req.ServerID == "" || req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "serverId and name are required")
	}

	port := 80
	if req.EdgePort != nil && *req.EdgePort > 0 {
		port = *req.EdgePort
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	conn := &models.IngressConnection{
		ID:                    generateID("ing"),
		ServerID:              req.ServerID,
		CredentialID:          req.CredentialID,
		AccountID:             req.AccountID,
		TunnelID:              req.TunnelID,
		Name:                  req.Name,
		Connector:             req.Connector,
		EdgePort:              port,
		ManagedTunnel:         req.Connector == "managed",
		DefaultForDomains:     req.DefaultForDomains,
		DefaultAccessPolicyID: req.DefaultAccessPolicyID,
		Status:                "ready",
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	h.connections[conn.ID] = conn
	return utils.Success(c, "Ingress tunnel connected", conn)
}

func (h *IngressHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req models.IngressUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid update payload")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	conn, exists := h.connections[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "ingress connection not found")
	}

	conn.DefaultForDomains = req.DefaultForDomains
	conn.DefaultAccessPolicyID = req.DefaultAccessPolicyID
	conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	return utils.Success(c, "Ingress connection updated", conn)
}

func (h *IngressHandler) Retry(c echo.Context) error {
	id := c.Param("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	conn, exists := h.connections[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "ingress connection not found")
	}

	conn.Status = "ready"
	conn.Error = nil
	conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	return utils.Success(c, "Ingress connection reconciled", conn)
}

func (h *IngressHandler) Disconnect(c echo.Context) error {
	id := c.Param("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.connections[id]; !exists {
		return echo.NewHTTPError(http.StatusNotFound, "ingress connection not found")
	}

	delete(h.connections, id)
	return utils.Success(c, "Ingress connection removed", map[string]bool{"success": true})
}

func (h *IngressHandler) GetDomain(c echo.Context) error {
	domainId := c.Param("domainId")

	h.mu.RLock()
	defer h.mu.RUnlock()

	hostname, exists := h.domains[domainId]
	if !exists {
		return utils.Success(c, "No ingress configured for domain", nil)
	}

	return utils.Success(c, "Domain ingress retrieved", hostname)
}

func (h *IngressHandler) SetDomain(c echo.Context) error {
	domainId := c.Param("domainId")
	var req models.IngressDomainRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid domain ingress payload")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if req.ConnectionID == nil || *req.ConnectionID == "" {
		delete(h.domains, domainId)
		return utils.Success(c, "Domain ingress disassociated", nil)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	hostname := &models.IngressHostname{
		ID:             generateID("igh"),
		DomainID:       &domainId,
		ConnectionID:   *req.ConnectionID,
		Hostname:       domainId,
		AccessPolicyID: req.AccessPolicyID,
		Status:         "ready",
		Phase:          "dns",
		UpdatedAt:      now,
	}

	h.domains[domainId] = hostname
	return utils.Success(c, "Domain ingress mapped", hostname)
}
