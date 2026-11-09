package managed

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *ControlsHandler) GetManagedSshStatus(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	status, exists := h.sshStore[id]
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
		h.sshStore[id] = status
	}

	return utils.Success(c, "SSH status retrieved", status)
}

func (h *ControlsHandler) SetManagedSsh(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	var req SetManagedSshRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

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
	h.sshStore[id] = status

	return utils.Success(c, "SSH status updated", status)
}

func (h *ControlsHandler) SetManagedSshKey(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req SetManagedSshKeyRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	return utils.Success(c, "SSH public key updated", nil)
}

func (h *ControlsHandler) SetManagedSshPassword(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	var req SetManagedSshPasswordRequest
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid request payload")
	}
	return utils.Success(c, "SSH password updated", nil)
}

func (h *ControlsHandler) GetManagedRuntimeStatus(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	status, exists := h.runtimeStore[id]
	if !exists {
		t := true
		status = &ManagedRuntimeStatus{
			Enabled: &t,
			Running: &t,
		}
		h.runtimeStore[id] = status
	}

	return utils.Success(c, "Runtime API status retrieved", status)
}

func (h *ControlsHandler) EnableManagedRuntime(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	id := c.Param("id")

	h.mu.Lock()
	defer h.mu.Unlock()

	t := true
	status := &ManagedRuntimeStatus{
		Enabled: &t,
		Running: &t,
	}
	h.runtimeStore[id] = status

	return utils.Success(c, "Runtime API enabled", status)
}

func (h *ControlsHandler) GetManagedRuntimeCredential(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	return utils.Success(c, "Runtime API credential retrieved", ManagedRuntimeCredential{
		Endpoint: "https://198.51.100.10:2376",
		Token:    randomTokenHex(32),
		Revision: "rev-1",
	})
}

func (h *ControlsHandler) RotateManagedRuntimeCredential(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}
	return utils.Success(c, "Runtime API credential rotated", ManagedRuntimeCredential{
		Endpoint: "https://198.51.100.10:2376",
		Token:    randomTokenHex(32),
		Revision: "rev-2",
	})
}
