package system

import (
	"codedock.run/codedock/internal/http/middleware"
	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
	"codedock.run/codedock/internal/services/clusterdata"
	"codedock.run/codedock/internal/services/operations"
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ClusterDataHandler struct {
	service    *clusterdata.Service
	operations *operations.Service
	apps       repositories.AppServiceRepository
	vars       repositories.ServiceVarRepository
}

func NewClusterDataHandler(service *clusterdata.Service, operations *operations.Service) *ClusterDataHandler {
	return &ClusterDataHandler{service: service, operations: operations}
}

func (h *ClusterDataHandler) SetBindings(apps repositories.AppServiceRepository, vars repositories.ServiceVarRepository) {
	h.apps = apps
	h.vars = vars
}
func (h *ClusterDataHandler) List(c echo.Context) error {
	result, err := h.service.List(c.Request().Context(), c.Param("id"), c.Param("clusterId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cluster databases", result)
}
func (h *ClusterDataHandler) Review(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request models.ClusterDataRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid database specification")
	}
	result, err := h.service.Review(c.Request().Context(), user.UserID, c.Param("id"), c.Param("clusterId"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review database effects", result)
}
func (h *ClusterDataHandler) ownedOperation(c echo.Context) (*models.Operation, error) {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}
	op, err := h.operations.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || op.Kind != "cluster-data" || op.UserID != user.UserID {
		return nil, echo.NewHTTPError(http.StatusNotFound, "database operation not found")
	}
	return op, nil
}
func (h *ClusterDataHandler) Apply(c echo.Context) error {
	op, err := h.ownedOperation(c)
	if err != nil {
		return err
	}
	var request models.OperationApply
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	if err := h.service.Apply(c.Request().Context(), op.UserID, op.ID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]any{"status": "success", "data": op})
}
func (h *ClusterDataHandler) Operation(c echo.Context) error {
	op, err := h.ownedOperation(c)
	if err != nil {
		return err
	}
	return utils.Success(c, "Database operation", op)
}
func (h *ClusterDataHandler) Cancel(c echo.Context) error {
	op, err := h.ownedOperation(c)
	if err != nil {
		return err
	}
	if err := h.operations.Cancel(c.Request().Context(), op.ID, op.UserID); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Execution interrupted; owned data remains for reviewed recovery", nil)
}
func (h *ClusterDataHandler) Credentials(c echo.Context) error {
	result, err := h.service.Credentials(c.Request().Context(), c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	c.Response().Header().Set("Cache-Control", "no-store")
	return utils.Success(c, "Private cluster connection", result)
}
func (h *ClusterDataHandler) Backups(c echo.Context) error {
	result, err := h.service.Backups(c.Request().Context(), c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "PostgreSQL base backups", result)
}

func (h *ClusterDataHandler) Verify(c echo.Context) error {
	if err := h.service.VerifyConnection(c.Request().Context(), c.Param("id"), c.Param("clusterId"), c.Param("databaseId")); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Database connection verified with stored credentials", nil)
}
