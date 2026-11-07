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

type ServerNetworksHandler struct {
	mu        sync.RWMutex
	networks  map[string]*models.ServerNetwork
	watchdogs map[string]*models.NetworkRollbackWatchdog
}

func NewServerNetworksHandler() *ServerNetworksHandler {
	return &ServerNetworksHandler{
		networks:  make(map[string]*models.ServerNetwork),
		watchdogs: make(map[string]*models.NetworkRollbackWatchdog),
	}
}

func (h *ServerNetworksHandler) List(c echo.Context) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var list []*models.ServerNetwork
	for _, net := range h.networks {
		list = append(list, net)
	}
	if list == nil {
		list = make([]*models.ServerNetwork, 0)
	}
	return utils.Success(c, "Server networks retrieved", list)
}

func (h *ServerNetworksHandler) Get(c echo.Context) error {
	id := c.Param("id")
	h.mu.RLock()
	defer h.mu.RUnlock()

	net, exists := h.networks[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "server network not found")
	}
	return utils.Success(c, "Server network retrieved", net)
}

func (h *ServerNetworksHandler) Create(c echo.Context) error {
	var req models.CreateNetworkRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid create network payload")
	}
	if req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "network name is required")
	}

	mode := req.Mode
	if mode == "" {
		mode = string(models.NetworkModeWireguard)
	}
	mtu := req.MTU
	if mtu <= 0 {
		mtu = 1420
	}
	probePort := req.VerificationPort
	if probePort <= 0 {
		probePort = 51821
	}
	subnet := req.SubnetCIDR
	if subnet == "" {
		subnet = "10.42.0.0/16"
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	opID := "op-" + hex.EncodeToString(generateBytes(6))
	netID := "net-" + hex.EncodeToString(generateBytes(6))

	// Arm 120s rollback watchdog
	expiresAt := time.Now().Add(120 * time.Second).UTC().Format(time.RFC3339)
	watchdog := &models.NetworkRollbackWatchdog{
		OperationID:      opID,
		Armed:            true,
		TimeoutSeconds:   120,
		RemainingSeconds: 120,
		ExpiresAt:        expiresAt,
		CanCommit:        true,
		CanRollback:      true,
		Status:           "armed",
	}

	net := &models.ServerNetwork{
		ID:               netID,
		Name:             req.Name,
		Mode:             mode,
		SubnetCIDR:       subnet,
		MTU:              mtu,
		VerificationPort: probePort,
		Status:           "applying",
		Revision:         1,
		Members:          req.Members,
		AccessPolicy:     req.AccessPolicy,
		Watchdog:         watchdog,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	h.networks[netID] = net
	h.watchdogs[opID] = watchdog

	return utils.Success(c, "Server network created and rollback watchdog armed", net)
}

func (h *ServerNetworksHandler) ReviseAccess(c echo.Context) error {
	id := c.Param("id")
	var req models.ReviseNetworkAccessRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid access policy payload")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	net, exists := h.networks[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "server network not found")
	}

	net.AccessPolicy = req.AccessPolicy
	net.Revision++
	net.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	return utils.Success(c, "Network access policy revised", net)
}

func (h *ServerNetworksHandler) GetWatchdog(c echo.Context) error {
	opID := c.Param("operationId")
	h.mu.RLock()
	defer h.mu.RUnlock()

	watchdog, exists := h.watchdogs[opID]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "rollback watchdog not found")
	}

	// Calculate remaining seconds
	expTime, err := time.Parse(time.RFC3339, watchdog.ExpiresAt)
	if err == nil {
		rem := int(time.Until(expTime).Seconds())
		if rem < 0 {
			rem = 0
			if watchdog.Status == "armed" {
				watchdog.Status = "expired"
				watchdog.Armed = false
				watchdog.CanCommit = false
			}
		}
		watchdog.RemainingSeconds = rem
	}

	return utils.Success(c, "Rollback watchdog status retrieved", watchdog)
}

func (h *ServerNetworksHandler) CommitWatchdog(c echo.Context) error {
	opID := c.Param("operationId")
	h.mu.Lock()
	defer h.mu.Unlock()

	watchdog, exists := h.watchdogs[opID]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "rollback watchdog not found")
	}
	if !watchdog.CanCommit {
		return echo.NewHTTPError(http.StatusConflict, "watchdog cannot be committed")
	}

	watchdog.Armed = false
	watchdog.CanCommit = false
	watchdog.CanRollback = false
	watchdog.Status = "committed"

	// Mark associated network ready
	for _, net := range h.networks {
		if net.Watchdog != nil && net.Watchdog.OperationID == opID {
			net.Status = "verified"
			net.Watchdog.Status = "committed"
			net.Watchdog.Armed = false
			net.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			break
		}
	}

	return utils.Success(c, "Network changes committed successfully", watchdog)
}

func (h *ServerNetworksHandler) RollbackWatchdog(c echo.Context) error {
	opID := c.Param("operationId")
	h.mu.Lock()
	defer h.mu.Unlock()

	watchdog, exists := h.watchdogs[opID]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "rollback watchdog not found")
	}

	watchdog.Armed = false
	watchdog.CanCommit = false
	watchdog.CanRollback = false
	watchdog.Status = "rolled_back"

	// Mark associated network rolled back
	for _, net := range h.networks {
		if net.Watchdog != nil && net.Watchdog.OperationID == opID {
			net.Status = "rolled_back"
			net.Watchdog.Status = "rolled_back"
			net.Watchdog.Armed = false
			net.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
			break
		}
	}

	return utils.Success(c, "Emergency network rollback executed", watchdog)
}

func (h *ServerNetworksHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.networks[id]; !exists {
		return echo.NewHTTPError(http.StatusNotFound, "server network not found")
	}

	delete(h.networks, id)
	return utils.Success(c, "Server network removed", map[string]bool{"success": true})
}

func generateBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
