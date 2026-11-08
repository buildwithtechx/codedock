package servers

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

var (
	managedStoreMu sync.Mutex
	networkStore   = make(map[string]*ServerNetworkSettings)
	sshStore       = make(map[string]*ManagedSshStatus)
	runtimeStore   = make(map[string]*ManagedRuntimeStatus)
)

func randomTokenHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *Handler) GetNetworkSettings(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	settings, exists := networkStore[id]
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
		networkStore[id] = settings
	}

	return utils.Success(c, "Network settings retrieved", settings)
}

func (h *Handler) UpdateNetworkSettings(c echo.Context) error {
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

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

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
	networkStore[id] = settings

	return utils.Success(c, "Network settings updated", settings)
}

func (h *Handler) GetManagedInfo(c echo.Context) error {
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

func (h *Handler) GetManagedBootLogs(c echo.Context) error {
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

func (h *Handler) GetManagedSshStatus(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	status, exists := sshStore[id]
	if !exists {
		t := true
		f := false
		status = &ManagedSshStatus{
			Enabled:                &t,
			KeyConfigured:          &t,
			PasswordConfigured:     &f,
			RequiresIdentityAccess: false,
			Connection: &ManagedSshConnection{
				User:    "root",
				Host:    "198.51.100.10",
				Bastion: "ssh.codedock.run",
				Command: "ssh -J ssh.codedock.run root@198.51.100.10",
			},
		}
		sshStore[id] = status
	}

	return utils.Success(c, "SSH status retrieved", status)
}

func (h *Handler) SetManagedSsh(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	var req SetManagedSshRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	status := &ManagedSshStatus{
		Enabled:                &req.Enabled,
		KeyConfigured:          &req.Enabled,
		PasswordConfigured:     &req.Enabled,
		RequiresIdentityAccess: false,
		Connection: &ManagedSshConnection{
			User:    "root",
			Host:    "198.51.100.10",
			Bastion: "ssh.codedock.run",
			Command: "ssh -J ssh.codedock.run root@198.51.100.10",
		},
	}
	sshStore[id] = status

	return utils.Success(c, "SSH configuration updated", ManagedSshResult{Status: *status})
}

func (h *Handler) SetManagedSshKey(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	status := sshStore[id]
	if status != nil {
		t := true
		status.KeyConfigured = &t
	}
	return utils.Success(c, "SSH key updated", map[string]bool{"ok": true})
}

func (h *Handler) SetManagedSshPassword(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	status := sshStore[id]
	if status != nil {
		t := true
		status.PasswordConfigured = &t
	}
	return utils.Success(c, "SSH password updated", map[string]bool{"ok": true})
}

func (h *Handler) GetManagedRuntimeStatus(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	status, exists := runtimeStore[id]
	if !exists {
		t := true
		status = &ManagedRuntimeStatus{
			Enabled: &t,
			Running: &t,
		}
		runtimeStore[id] = status
	}

	return utils.Success(c, "Runtime status retrieved", status)
}

func (h *Handler) EnableManagedRuntime(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	t := true
	status := &ManagedRuntimeStatus{
		Enabled: &t,
		Running: &t,
	}
	runtimeStore[id] = status

	return utils.Success(c, "Runtime enabled", status)
}

func (h *Handler) GetManagedRuntimeCredential(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	cred := ManagedRuntimeCredential{
		Endpoint:  fmt.Sprintf("https://api.codedock.run/runtime/%s", id),
		Token:     fmt.Sprintf("cd_rt_%s", randomTokenHex(24)),
		Revision:  randomTokenHex(32),
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
	return utils.Success(c, "Runtime credential generated", cred)
}

func (h *Handler) RotateManagedRuntimeCredential(c echo.Context) error {
	return h.GetManagedRuntimeCredential(c)
}
