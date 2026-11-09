package managed

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ControlsHandler struct {
	mu           sync.Mutex
	networkStore map[string]*ServerNetworkSettings
	sshStore     map[string]*ManagedSshStatus
	runtimeStore map[string]*ManagedRuntimeStatus
}

func NewControlsHandler() *ControlsHandler {
	return &ControlsHandler{
		networkStore: make(map[string]*ServerNetworkSettings),
		sshStore:     make(map[string]*ManagedSshStatus),
		runtimeStore: make(map[string]*ManagedRuntimeStatus),
	}
}

func randomTokenHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *ControlsHandler) GetNetworkSettings(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	settings, exists := h.networkStore[id]
	if !exists {
		t := true
		settings = &ServerNetworkSettings{
			InternetAccess: &t,
			IngressPorts:   []int{80, 443, 22},
			Egress:         []string{"*"},
			PrivateIP:      "10.0.0.4",
			OutboundIP:     "198.51.100.10",
			OutboundMode:   "managed",
			Revision:       randomTokenHex(32),
		}
		h.networkStore[id] = settings
	}

	return utils.Success(c, "Network settings retrieved", settings)
}

func (h *ControlsHandler) UpdateNetworkSettings(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	var req UpdateServerNetworkSettingsRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	t := req.InternetAccess
	egress := req.Egress
	if len(egress) == 0 {
		egress = []string{"*"}
	}

	settings := &ServerNetworkSettings{
		InternetAccess: &t,
		IngressPorts:   []int{80, 443, 22},
		Egress:         egress,
		PrivateIP:      "10.0.0.4",
		OutboundIP:     "198.51.100.10",
		OutboundMode:   "managed",
		Revision:       randomTokenHex(32),
	}
	h.networkStore[id] = settings

	return utils.Success(c, "Network settings updated", settings)
}

func (h *ControlsHandler) GetManagedInfo(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	info := ManagedServerInfo{
		WorkspaceID:     fmt.Sprintf("ws-%s", id),
		Image:           "codedock-runner:v1.4",
		State:           "running",
		Mode:            "permanent",
		OperatingSystem: "Debian GNU/Linux 12 (bookworm)",
		RestartPolicy:   "always",
		Resources: ManagedServerResources{
			CPUCores: 4,
			MemoryMB: 8192,
			DiskMB:   80000,
		},
	}
	return utils.Success(c, "Managed server info retrieved", info)
}

func (h *ControlsHandler) GetManagedBootLogs(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	bootLogs := fmt.Sprintf("[%s] systemd 252.12-1~deb12u1 running in system mode\n[%s] Detected virtualization kvm\n[%s] Initialized codedock runtime engine\n[%s] Starting Docker Application Container Engine...\n[%s] Started Docker Application Container Engine\n[%s] Connected to Codedock management plane\n[%s] Ready for application workloads",
		time.Now().Add(-5*time.Minute).Format(time.RFC3339),
		time.Now().Add(-5*time.Minute+time.Second).Format(time.RFC3339),
		time.Now().Add(-4*time.Minute).Format(time.RFC3339),
		time.Now().Add(-4*time.Minute+2*time.Second).Format(time.RFC3339),
		time.Now().Add(-3*time.Minute).Format(time.RFC3339),
		time.Now().Add(-2*time.Minute).Format(time.RFC3339),
		time.Now().Add(-1*time.Minute).Format(time.RFC3339),
	)
	return utils.Success(c, "Boot logs retrieved", ManagedLogsResponse{Logs: bootLogs, Truncated: false})
}
