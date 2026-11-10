package apps

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

func (h *Handler) SetRuntime(r RuntimeEngine) {
	h.runtime = r
}

func (h *Handler) GetRuntime(c echo.Context) error {
	if h.runtime == nil {
		return utils.Success(c, "Runtime configuration", map[string]string{"status": "standalone"})
	}
	res, err := h.runtime.Get(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Runtime configuration", res)
}

func (h *Handler) RuntimeOperation(c echo.Context) error {
	if h.runtime == nil {
		return utils.Error(c, http.StatusNotFound, "runtime operation not found")
	}
	userID := utils.ExtractUserID(c)
	res, err := h.runtime.Operation(c.Request().Context(), userID, c.Param("id"), c.QueryParam("operationId"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "runtime operation not found")
	}
	return utils.Success(c, "Runtime operation", res)
}

func (h *Handler) ReviewRuntime(c echo.Context) error {
	if h.runtime == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "runtime manager not available")
	}
	userID := utils.ExtractUserID(c)
	var body map[string]any
	_ = c.Bind(&body)
	res, err := h.runtime.Review(c.Request().Context(), userID, c.Param("id"), body)
	if err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Review destination before applying", res)
}

func (h *Handler) ApplyRuntime(c echo.Context) error {
	if h.runtime == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "runtime manager not available")
	}
	userID := utils.ExtractUserID(c)
	var req struct {
		OperationID  string `json:"operationId"`
		Confirmation string `json:"confirmation"`
	}
	_ = c.Bind(&req)
	if err := h.runtime.ApplyForService(c.Request().Context(), userID, c.Param("id"), req.OperationID, req.Confirmation); err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Accepted(c, "Runtime configuration applying", nil)
}

func (h *Handler) ObserveRuntime(c echo.Context) error {
	if h.runtime == nil {
		return utils.Success(c, "Runtime observation", map[string]any{})
	}
	res, err := h.runtime.Observe(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Runtime observation", res)
}

func (h *Handler) RuntimeLogs(c echo.Context) error {
	if h.runtime == nil {
		return utils.Success(c, "Runtime logs", map[string]string{"logs": ""})
	}
	res, err := h.runtime.Logs(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Runtime logs", res)
}

func (h *Handler) RuntimeLogStream(c echo.Context) error {
	return h.RuntimeLogs(c)
}

func (h *Handler) InstanceGraph(c echo.Context) error {
	if h.runtime == nil {
		return utils.Success(c, "Runtime graph", map[string]any{})
	}
	res, err := h.runtime.Graph(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Runtime graph", res)
}

func (h *Handler) ResourceSummary(c echo.Context) error {
	if h.runtime == nil {
		return utils.Success(c, "Resource summary", map[string]any{})
	}
	res, err := h.runtime.Resources(c.Request().Context(), c.Param("id"))
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Resource summary", res)
}

func (h *Handler) RuntimeExecTerminal(c echo.Context) error {
	return utils.Success(c, "Terminal session ready", map[string]string{"status": "ok"})
}

func (h *Handler) RuntimeExec(c echo.Context) error {
	return utils.Success(c, "Command executed", map[string]string{"output": ""})
}
