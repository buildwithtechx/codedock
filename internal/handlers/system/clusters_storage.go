package system

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/utils"
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

func (h *ClusterHandler) ReviewStorage(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request models.StorageSetupRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid storage version")
	}
	result, err := h.service.ReviewStorageSetup(c.Request().Context(), user.UserID, c.Param("id"), c.Param("clusterId"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review storage installation", result)
}

func (h *ClusterHandler) ApplyStorage(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request models.OperationApply
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.operations.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || op.Kind != "cluster-storage" || op.UserID != user.UserID {
		return echo.NewHTTPError(http.StatusNotFound, "storage operation not found")
	}
	if err := h.service.ApplyStorageSetup(c.Request().Context(), user.UserID, op.ID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Storage installation applied", op)
}
