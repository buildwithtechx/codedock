package clusters

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *Handler) GetScaling(c echo.Context) error {
	id := c.Param("clusterId")
	h.scalingMu.RLock()
	defer h.scalingMu.RUnlock()

	sc, exists := h.scaling[id]
	if !exists {
		return utils.Success(c, "Cluster scaling status", &ClusterScaling{
			Status:        "not_enabled",
			K3sVersion:    "v1.31.1+k3s1",
			ControlNodes:  1,
			WorkerNodes:   0,
			PodCIDR:       "10.42.0.0/16",
			ServiceCIDR:   "10.43.0.0/16",
			ClusterDomain: "cluster.local",
			UpdatedAt:     time.Now().UTC().Format(time.RFC3339),
		})
	}
	return utils.Success(c, "Cluster scaling status", sc)
}

func (h *Handler) SetupScaling(c echo.Context) error {
	id := c.Param("clusterId")
	h.scalingMu.Lock()
	defer h.scalingMu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	sc := &ClusterScaling{
		Status:        "ready",
		K3sVersion:    "v1.31.1+k3s1",
		ControlNodes:  1,
		WorkerNodes:   2,
		PodCIDR:       "10.42.0.0/16",
		ServiceCIDR:   "10.43.0.0/16",
		ClusterDomain: "cluster.local",
		UpdatedAt:     now,
	}
	h.scaling[id] = sc
	return utils.Success(c, "Cluster scaling enabled", sc)
}

func (h *Handler) RetryScaling(c echo.Context) error {
	id := c.Param("clusterId")
	h.scalingMu.Lock()
	defer h.scalingMu.Unlock()

	sc, exists := h.scaling[id]
	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "scaling not initialized")
	}
	sc.Status = "ready"
	sc.Error = ""
	sc.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return utils.Success(c, "Cluster scaling retried", sc)
}

func (h *Handler) RemoveScaling(c echo.Context) error {
	id := c.Param("clusterId")
	h.scalingMu.Lock()
	defer h.scalingMu.Unlock()

	sc := &ClusterScaling{
		Status:    "removed",
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	h.scaling[id] = sc
	return utils.Success(c, "Cluster scaling removed", sc)
}
