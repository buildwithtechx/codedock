package system

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/utils"
)

var (
	managedStoreMu sync.RWMutex
	workloadStore  = make(map[string][]models.ManagedWorkload)
	networkStore   = make(map[string]*models.ServerNetworkSettings)
	sshStore       = make(map[string]*models.ManagedSshStatus)
	runtimeStore   = make(map[string]*models.ManagedRuntimeStatus)
)

func randomTokenHex(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *ServerHandler) GetNetworkSettings(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
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
		settings = &models.ServerNetworkSettings{
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

func (h *ServerHandler) UpdateNetworkSettings(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	var req models.UpdateServerNetworkSettingsRequest
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

	settings := &models.ServerNetworkSettings{
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

func (h *ServerHandler) GetManagedInfo(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "server id required")
	}

	info := models.ManagedServerInfo{
		WorkspaceID:     fmt.Sprintf("ws-%s", id),
		Image:           "codedock-runner:v1.4",
		State:           "running",
		Mode:            "permanent",
		OperatingSystem: "Debian GNU/Linux 12 (bookworm)",
		RestartPolicy:   "always",
		Resources: models.ManagedServerResources{
			CPUCores: 4,
			MemoryMB: 8192,
			DiskMB:   80000,
		},
	}
	return utils.Success(c, "Managed server info retrieved", info)
}

func (h *ServerHandler) GetManagedBootLogs(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
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
	return utils.Success(c, "Boot logs retrieved", models.ManagedLogsResponse{Logs: bootLogs, Truncated: false})
}

func (h *ServerHandler) GetManagedSshStatus(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	status, exists := sshStore[id]
	if !exists {
		t := true
		f := false
		status = &models.ManagedSshStatus{
			Enabled:                &t,
			KeyConfigured:          &t,
			PasswordConfigured:     &f,
			RequiresIdentityAccess: false,
			Connection: &models.ManagedSshConnection{
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

func (h *ServerHandler) SetManagedSsh(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	var req models.SetManagedSshRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	status := &models.ManagedSshStatus{
		Enabled:                &req.Enabled,
		KeyConfigured:          &req.Enabled,
		PasswordConfigured:     &req.Enabled,
		RequiresIdentityAccess: false,
		Connection: &models.ManagedSshConnection{
			User:    "root",
			Host:    "198.51.100.10",
			Bastion: "ssh.codedock.run",
			Command: "ssh -J ssh.codedock.run root@198.51.100.10",
		},
	}
	sshStore[id] = status

	return utils.Success(c, "SSH configuration updated", models.ManagedSshResult{Status: *status})
}

func (h *ServerHandler) SetManagedSshKey(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
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

func (h *ServerHandler) SetManagedSshPassword(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
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

func (h *ServerHandler) GetManagedRuntimeStatus(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	status, exists := runtimeStore[id]
	if !exists {
		t := true
		status = &models.ManagedRuntimeStatus{
			Enabled: &t,
			Running: &t,
		}
		runtimeStore[id] = status
	}

	return utils.Success(c, "Runtime status retrieved", status)
}

func (h *ServerHandler) EnableManagedRuntime(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	t := true
	status := &models.ManagedRuntimeStatus{
		Enabled: &t,
		Running: &t,
	}
	runtimeStore[id] = status

	return utils.Success(c, "Runtime enabled", status)
}

func (h *ServerHandler) GetManagedRuntimeCredential(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	cred := models.ManagedRuntimeCredential{
		Endpoint:  fmt.Sprintf("https://api.codedock.run/runtime/%s", id),
		Token:     fmt.Sprintf("cd_rt_%s", randomTokenHex(24)),
		Revision:  randomTokenHex(32),
		ExpiresAt: time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}
	return utils.Success(c, "Runtime credential generated", cred)
}

func (h *ServerHandler) RotateManagedRuntimeCredential(c echo.Context) error {
	return h.GetManagedRuntimeCredential(c)
}

func (h *ServerHandler) GetManagedWorkloads(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	items, exists := workloadStore[id]
	if !exists {
		items = []models.ManagedWorkload{
			{
				ID:            "wk-cron-cleaner",
				Name:          "Nightly Maintenance Worker",
				State:         "running",
				RestartPolicy: "always",
				Source:        "system",
				Manageable:    true,
			},
		}
		workloadStore[id] = items
	}

	return utils.Success(c, "Workloads retrieved", models.ManagedWorkloadsResponse{
		Workloads: items,
		Truncated: false,
	})
}

func (h *ServerHandler) CreateManagedWorkload(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	var req models.CreateManagedWorkloadRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	item := models.ManagedWorkload{
		ID:            fmt.Sprintf("wk-%s", randomTokenHex(6)),
		Name:          req.Name,
		State:         "running",
		RestartPolicy: req.RestartPolicy,
		Source:        "manual",
		Manageable:    true,
	}

	workloadStore[id] = append(workloadStore[id], item)

	return utils.Success(c, "Workload created", item)
}

func (h *ServerHandler) ControlManagedWorkload(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	var req models.ControlManagedWorkloadRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	managedStoreMu.Lock()
	defer managedStoreMu.Unlock()

	var updated *models.ManagedWorkload
	items := workloadStore[id]
	newItems := make([]models.ManagedWorkload, 0, len(items))
	for _, it := range items {
		if it.ID == req.WorkloadID {
			if req.Action == "delete" {
				continue
			}
			if req.Action == "stop" {
				it.State = "stopped"
			} else if req.Action == "start" || req.Action == "restart" {
				it.State = "running"
			}
			updated = &it
		}
		newItems = append(newItems, it)
	}
	workloadStore[id] = newItems

	return utils.Success(c, "Workload controlled", models.ControlManagedWorkloadResult{
		Ok:       true,
		Workload: updated,
	})
}

func (h *ServerHandler) GetManagedWorkloadLogs(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	logs := fmt.Sprintf("[%s] process started with PID 14209\n[%s] listening for scheduled job triggers\n[%s] job execution succeeded (duration: 412ms)\n[%s] waiting for next interval...",
		time.Now().Add(-10*time.Minute).Format(time.RFC3339),
		time.Now().Add(-9*time.Minute).Format(time.RFC3339),
		time.Now().Add(-4*time.Minute).Format(time.RFC3339),
		time.Now().Add(-1*time.Minute).Format(time.RFC3339),
	)

	return utils.Success(c, "Workload logs retrieved", models.ManagedLogsResponse{
		Logs:      logs,
		Truncated: false,
	})
}
