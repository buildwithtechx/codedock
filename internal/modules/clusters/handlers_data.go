package clusters

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type DataHandler struct {
	service DataService
}

func NewDataHandler(service DataService) *DataHandler {
	return &DataHandler{service: service}
}

func (h *DataHandler) List(c echo.Context) error {
	result, err := h.service.List(c.Request().Context(), c.Param("id"), c.Param("clusterId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Cluster databases", result)
}

func (h *DataHandler) Review(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	var request ClusterDataRequest
	if err := c.Bind(&request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid database specification")
	}
	result, err := h.service.Review(c.Request().Context(), userID, c.Param("id"), c.Param("clusterId"), request)
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review database effects", result)
}

func (h *DataHandler) Apply(c echo.Context) error {
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

func (h *DataHandler) Operation(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	op, err := h.service.GetOperation(c.Request().Context(), userID, c.Param("operationId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "database operation not found")
	}
	return utils.Success(c, "Database operation", op)
}

func (h *DataHandler) Cancel(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	if err := h.service.CancelOperation(c.Request().Context(), userID, c.Param("operationId")); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Execution interrupted; owned data remains for reviewed recovery", nil)
}

func (h *DataHandler) Credentials(c echo.Context) error {
	result, err := h.service.Credentials(c.Request().Context(), c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	c.Response().Header().Set("Cache-Control", "no-store")
	return utils.Success(c, "Private cluster connection", result)
}

func (h *DataHandler) Backups(c echo.Context) error {
	result, err := h.service.Backups(c.Request().Context(), c.Param("id"), c.Param("clusterId"), c.Param("databaseId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "PostgreSQL base backups", result)
}

func (h *DataHandler) Verify(c echo.Context) error {
	if err := h.service.VerifyConnection(c.Request().Context(), c.Param("id"), c.Param("clusterId"), c.Param("databaseId")); err != nil {
		return echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Database connection verified with stored credentials", nil)
}
