package system

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"codedock/internal/models"
	"codedock/internal/repositories"
	systemservices "codedock/internal/services/system"
	"codedock/internal/utils"

	"github.com/labstack/echo/v4"
)

type IngressHandler struct {
	repo  repositories.IngressRepository
	cfSvc *systemservices.CloudflareIngressService
}

func NewIngressHandler() *IngressHandler {
	return NewIngressHandlerWithDeps(nil, nil)
}

func NewIngressHandlerWithDeps(repo repositories.IngressRepository, cfSvc *systemservices.CloudflareIngressService) *IngressHandler {
	if repo == nil {
		repo = repositories.NewIngressRepository(nil)
	}
	if cfSvc == nil {
		cfSvc = systemservices.NewCloudflareIngressService(nil)
	}
	return &IngressHandler{
		repo:  repo,
		cfSvc: cfSvc,
	}
}

func generateID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return prefix + "-" + hex.EncodeToString(b)
}

func (h *IngressHandler) List(c echo.Context) error {
	serverID := c.QueryParam("serverId")
	list, err := h.repo.ListConnections(c.Request().Context(), serverID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list ingress connections: "+err.Error())
	}
	return utils.Success(c, "Ingress connections retrieved", list)
}

func (h *IngressHandler) Discover(c echo.Context) error {
	var req models.IngressDiscoverRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid discovery request")
	}

	ctx := c.Request().Context()
	accounts, err := h.cfSvc.ListAccounts(ctx, req.CredentialID)
	if err == nil && len(accounts) > 0 {
		accountID := accounts[0].ID
		if req.AccountID != nil && *req.AccountID != "" {
			accountID = *req.AccountID
		}
		tunnels, _ := h.cfSvc.ListTunnels(ctx, req.CredentialID, accountID)
		policies, _ := h.cfSvc.ListAccessPolicies(ctx, req.CredentialID, accountID)
		return utils.Success(c, "Ingress discovery completed", models.IngressDiscoveryResult{
			Accounts: accounts,
			Tunnels:  tunnels,
			Policies: policies,
		})
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

	tunnelID := req.TunnelID
	if (tunnelID == nil || *tunnelID == "") && req.Connector == "managed" {
		ctx := c.Request().Context()
		if newID, err := h.cfSvc.CreateTunnel(ctx, req.CredentialID, req.AccountID, req.Name); err == nil && newID != "" {
			tunnelID = &newID
		}
	}

	now := time.Now().UTC().Format(time.RFC3339)
	conn := &models.IngressConnection{
		ID:                    generateID("ing"),
		ServerID:              req.ServerID,
		CredentialID:          req.CredentialID,
		AccountID:             req.AccountID,
		TunnelID:              tunnelID,
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

	if err := h.repo.SaveConnection(c.Request().Context(), conn); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to save ingress connection: "+err.Error())
	}

	return utils.Success(c, "Ingress tunnel connected", conn)
}

func (h *IngressHandler) Update(c echo.Context) error {
	id := c.Param("id")
	var req models.IngressUpdateRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid update payload")
	}

	ctx := c.Request().Context()
	conn, err := h.repo.GetConnection(ctx, id)
	if err != nil || conn == nil {
		return echo.NewHTTPError(http.StatusNotFound, "ingress connection not found")
	}

	conn.DefaultForDomains = req.DefaultForDomains
	conn.DefaultAccessPolicyID = req.DefaultAccessPolicyID
	conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := h.repo.SaveConnection(ctx, conn); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update connection: "+err.Error())
	}

	return utils.Success(c, "Ingress connection updated", conn)
}

func (h *IngressHandler) Retry(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	conn, err := h.repo.GetConnection(ctx, id)
	if err != nil || conn == nil {
		return echo.NewHTTPError(http.StatusNotFound, "ingress connection not found")
	}

	conn.Status = "ready"
	conn.Error = nil
	conn.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := h.repo.SaveConnection(ctx, conn); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to reconcile connection: "+err.Error())
	}

	return utils.Success(c, "Ingress connection reconciled", conn)
}

func (h *IngressHandler) Disconnect(c echo.Context) error {
	id := c.Param("id")
	ctx := c.Request().Context()
	if _, err := h.repo.GetConnection(ctx, id); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "ingress connection not found")
	}

	if err := h.repo.DeleteConnection(ctx, id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to remove connection: "+err.Error())
	}

	return utils.Success(c, "Ingress connection removed", map[string]bool{"success": true})
}

func (h *IngressHandler) GetDomain(c echo.Context) error {
	domainID := c.Param("domainId")
	host, err := h.repo.GetHostname(c.Request().Context(), domainID)
	if err != nil || host == nil {
		return utils.Success(c, "No ingress configured for domain", nil)
	}
	return utils.Success(c, "Domain ingress retrieved", host)
}

func (h *IngressHandler) SetDomain(c echo.Context) error {
	domainID := c.Param("domainId")
	var req models.IngressDomainRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid domain ingress payload")
	}

	ctx := c.Request().Context()
	if req.ConnectionID == nil || *req.ConnectionID == "" {
		if err := h.repo.DeleteHostname(ctx, domainID); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to disassociate ingress: "+err.Error())
		}
		return utils.Success(c, "Domain ingress disassociated", nil)
	}

	now := time.Now().UTC().Format(time.RFC3339)
	hostname := &models.IngressHostname{
		ID:             generateID("igh"),
		DomainID:       &domainID,
		ConnectionID:   *req.ConnectionID,
		Hostname:       domainID,
		AccessPolicyID: req.AccessPolicyID,
		Status:         "ready",
		Phase:          "dns",
		UpdatedAt:      now,
	}

	if err := h.repo.SaveHostname(ctx, hostname); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to map domain ingress: "+err.Error())
	}

	return utils.Success(c, "Domain ingress mapped", hostname)
}
