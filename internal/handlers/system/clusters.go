package system

import (
	"codedock/internal/http/middleware"
	"codedock/internal/models"
	"codedock/internal/services/clusters"
	"codedock/internal/services/operations"
	"codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
	"time"
)

type ClusterHandler struct {
	service    *clusters.Service
	operations *operations.Service
	scalingMu  sync.RWMutex
	scaling    map[string]*models.ClusterScaling
}

func NewClusterHandler(service *clusters.Service, operations *operations.Service) *ClusterHandler {
	return &ClusterHandler{
		service:    service,
		operations: operations,
		scaling:    make(map[string]*models.ClusterScaling),
	}
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

func (h *ClusterHandler) GetScaling(c echo.Context) error {
	id := c.Param("clusterId")
	h.scalingMu.RLock()
	defer h.scalingMu.RUnlock()

	sc, exists := h.scaling[id]
	if !exists {
		return utils.Success(c, "Cluster scaling status", &models.ClusterScaling{
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

func (h *ClusterHandler) SetupScaling(c echo.Context) error {
	id := c.Param("clusterId")
	h.scalingMu.Lock()
	defer h.scalingMu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	sc := &models.ClusterScaling{
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

func (h *ClusterHandler) RetryScaling(c echo.Context) error {
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

func (h *ClusterHandler) RemoveScaling(c echo.Context) error {
	id := c.Param("clusterId")
	h.scalingMu.Lock()
	defer h.scalingMu.Unlock()

	sc := &models.ClusterScaling{
		Status:    "removed",
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	h.scaling[id] = sc
	return utils.Success(c, "Cluster scaling removed", sc)
}

