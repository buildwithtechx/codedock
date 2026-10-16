package system

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/models"
	"codedock/internal/utils"
)

func (h *ServerHandler) StreamMonitor(c echo.Context) error {
	userClaims, ok := c.Get("user").(*models.UserClaims)
	if !ok || userClaims == nil {
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

	server, err := h.serverService.GetServer(c.Request().Context(), serverID, userClaims.UserID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if server == nil {
		return utils.Error(c, http.StatusNotFound, "server not found")
	}

	writer, flush, ok := startEventStream(c)
	if !ok {
		return utils.Error(c, http.StatusInternalServerError, "streaming unsupported")
	}

	emitSample := func() error {
		raw := GetHostMetricsPayload()
		if format == "stats" {
			writeStreamEvent(writer, flush, "stats", string(raw))
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
		writeStreamEvent(writer, flush, "usage", string(encoded))
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
