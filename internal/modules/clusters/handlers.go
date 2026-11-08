package clusters

import (
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type Handler struct {
	service   Service
	scalingMu sync.RWMutex
	scaling   map[string]*ClusterScaling
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
		scaling: make(map[string]*ClusterScaling),
	}
}

func (h *Handler) List(c echo.Context) error {
	result, err := h.service.List(c.Request().Context(), c.Param("id"))
	if err != nil {
		return err
	}
	return utils.Success(c, "Clusters observed", result)
}

func (h *Handler) Review(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request ClusterReviewRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid cluster specification")
	}
	result, err := h.service.Review(c.Request().Context(), userID, c.Param("id"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review cluster effects before applying", result)
}

func (h *Handler) Apply(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request OperationApply
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.service.Apply(c.Request().Context(), userID, c.Param("operationId"), request.Confirmation)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return c.JSON(http.StatusAccepted, map[string]any{"status": "success", "data": op})
}

func (h *Handler) GetOperation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	operation, err := h.service.GetOperation(c.Request().Context(), userID, c.Param("operationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "cluster operation not found")
	}
	return utils.Success(c, "Cluster operation observed", operation)
}

func (h *Handler) Cancel(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	if err := h.service.CancelOperation(c.Request().Context(), userID, c.Param("operationId")); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cluster interruption requested; inspect nodes before retrying", nil)
}

func (h *Handler) Storage(c echo.Context) error {
	result, err := h.service.StorageInventory(c.Request().Context(), c.Param("id"), c.Param("clusterId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cluster storage inventory", result)
}

func (h *Handler) StorageHealth(c echo.Context) error {
	if err := h.service.StorageHealth(c.Request().Context(), c.Param("id"), c.Param("clusterId")); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Shared storage provisioner is installed and claims are bound", nil)
}

func (h *Handler) Preflight(c echo.Context) error {
	result, err := h.service.NetworkPreflight(c.Request().Context(), c.Param("id"), c.Param("clusterId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cluster network preflight", result)
}

func (h *Handler) ReviewStorage(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request StorageSetupRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid storage version")
	}
	result, err := h.service.ReviewStorageSetup(c.Request().Context(), userID, c.Param("id"), c.Param("clusterId"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review storage installation", result)
}

func (h *Handler) ApplyStorage(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request OperationApply
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid confirmation")
	}
	op, err := h.service.ApplyStorageSetup(c.Request().Context(), userID, c.Param("operationId"), request.Confirmation)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Storage installation applied", op)
}
