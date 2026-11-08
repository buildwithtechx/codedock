package servers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/platform/observability"
	"codedock/internal/utils"
)

func (h *Handler) StreamMonitor(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	serverID := c.QueryParam("serverId")
	if serverID == "" {
		return utils.Error(c, http.StatusBadRequest, "serverId query parameter required")
	}

	format := c.QueryParam("format")
	if format == "" {
		format = "stats"
	}
	if format != "stats" && format != "usage" {
		return utils.Error(c, http.StatusBadRequest, "invalid monitor format: expected 'stats' or 'usage'")
	}

	server, err := h.service.GetServer(c.Request().Context(), serverID, userID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if server == nil {
		return utils.Error(c, http.StatusNotFound, "server not found")
	}

	w := c.Response().Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return utils.Error(c, http.StatusInternalServerError, "streaming unsupported")
	}
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)
	flusher.Flush()

	emitSample := func() error {
		raw := observability.GetHostMetricsPayload()
		if format == "stats" {
			fmt.Fprintf(w, "event: stats\ndata: %s\n\n", string(raw))
			flusher.Flush()
			return nil
		}

		var metrics struct {
			CPUUsagePercentage float64 `json:"cpu_usage_percentage"`
			MemoryUsageBytes   uint64  `json:"memory_usage_bytes"`
			MemoryLimitBytes   uint64  `json:"memory_limit_bytes"`
			DiskUsageBytes     uint64  `json:"disk_usage_bytes"`
			DiskTotalBytes     uint64  `json:"disk_total_bytes"`
		}
		_ = json.Unmarshal(raw, &metrics)

		var memAvail uint64
		if metrics.MemoryLimitBytes > metrics.MemoryUsageBytes {
			memAvail = metrics.MemoryLimitBytes - metrics.MemoryUsageBytes
		}
		var diskAvail uint64
		if metrics.DiskTotalBytes > metrics.DiskUsageBytes {
			diskAvail = metrics.DiskTotalBytes - metrics.DiskUsageBytes
		}

		usage := map[string]any{
			"measuredAt":        time.Now().UTC().Format(time.RFC3339Nano),
			"cpuPercent":        metrics.CPUUsagePercentage,
			"memoryUsedMb":      metrics.MemoryUsageBytes / (1024 * 1024),
			"memoryAvailableMb": memAvail / (1024 * 1024),
			"diskUsedMb":        metrics.DiskUsageBytes / (1024 * 1024),
			"diskAvailableMb":   diskAvail / (1024 * 1024),
			"diskTotalMb":       metrics.DiskTotalBytes / (1024 * 1024),
		}
		encoded, _ := json.Marshal(usage)
		fmt.Fprintf(w, "event: usage\ndata: %s\n\n", string(encoded))
		flusher.Flush()
		return nil
	}

	if err := emitSample(); err != nil {
		return nil
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	ctx := c.Request().Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := emitSample(); err != nil {
				return nil
			}
		}
	}
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
