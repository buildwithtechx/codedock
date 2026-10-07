package system

import (
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *ClusterHandler) Storage(c echo.Context) error {
	result, err := h.service.StorageInventory(c.Request().Context(), c.Param("id"), c.Param("clusterId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cluster storage inventory", result)
}

func (h *ClusterHandler) StorageHealth(c echo.Context) error {
	if err := h.service.StorageHealth(c.Request().Context(), c.Param("id"), c.Param("clusterId")); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Shared storage provisioner is installed and claims are bound", nil)
}

func (h *ClusterHandler) Preflight(c echo.Context) error {
	result, err := h.service.NetworkPreflight(c.Request().Context(), c.Param("id"), c.Param("clusterId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cluster network preflight", result)
}
