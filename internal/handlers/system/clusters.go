package system

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/services/clusters"
	"codedock/internal/services/operations"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ClusterHandler struct {
	service    *clusters.Service
	operations *operations.Service
}

func NewClusterHandler(service *clusters.Service, operations *operations.Service) *ClusterHandler {
	return &ClusterHandler{service, operations}
}
func (h *ClusterHandler) List(c echo.Context) error {
	result, err := h.service.List(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return utils.Success(c, "Clusters observed", result)
}
func (h *ClusterHandler) Review(c echo.Context) error {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request models.ClusterReviewRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cluster specification")
	}
	result, err := h.service.Review(c.Request().Context(), user.UserID, c.Param("id"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review cluster effects before applying", result)
}
func (h *ClusterHandler) ownedOperation(c echo.Context) (*models.Operation, error) {
	user := middleware.GetUserClaimsFromContext(c.Request().Context())
	if user == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized)
	}
	operation, err := h.operations.Get(c.Request().Context(), c.Param("operationId"))
	if err != nil || operation.Kind != "cluster" || operation.UserID != user.UserID {
		return nil, echo.NewHTTPError(http.StatusNotFound, "cluster operation not found")
	}
	return operation, nil
}
func (h *ClusterHandler) Apply(c echo.Context) error {
	operation, err := h.ownedOperation(c)
	if err != nil {
		return err
	}
	var request models.OperationApply
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	if err := h.service.Apply(c.Request().Context(), operation.UserID, operation.ID, request.Confirmation); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]any{"status": "success", "data": operation})
}
func (h *ClusterHandler) GetOperation(c echo.Context) error {
	operation, err := h.ownedOperation(c)
	if err != nil {
		return err
	}
	return utils.Success(c, "Cluster operation observed", operation)
}
func (h *ClusterHandler) Cancel(c echo.Context) error {
	operation, err := h.ownedOperation(c)
	if err != nil {
		return err
	}
	if err := h.operations.Cancel(c.Request().Context(), operation.ID, operation.UserID); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cluster interruption requested; inspect nodes before retrying", nil)
}
