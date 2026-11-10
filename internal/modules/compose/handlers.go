package compose

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ComposeStackHandler struct {
	stacks ComposeStackOperations
}

func NewComposeStackHandler(stacks ComposeStackOperations) *ComposeStackHandler {
	return &ComposeStackHandler{stacks: stacks}
}

func (h *ComposeStackHandler) List(c echo.Context) error {
	projectID := c.Param("id")
	if h.stacks == nil {
		return utils.Success(c, "Compose stacks", []any{})
	}
	res, err := h.stacks.List(c.Request().Context(), projectID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Compose stacks", res)
}

func (h *ComposeStackHandler) Config(c echo.Context) error {
	projectID := c.Param("id")
	stackID := c.Param("stackId")
	if h.stacks == nil {
		return utils.Error(c, http.StatusNotFound, "stack not found")
	}
	res, err := h.stacks.Config(c.Request().Context(), projectID, stackID)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Compose configuration", res)
}

func (h *ComposeStackHandler) Review(c echo.Context) error {
	projectID := c.Param("id")
	var req any
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid stack request")
	}
	if h.stacks == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "compose stack service unavailable")
	}
	res, err := h.stacks.Review(c.Request().Context(), projectID, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Stack reviewed", res)
}

func (h *ComposeStackHandler) Save(c echo.Context) error {
	projectID := c.Param("id")
	var req any
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid stack request")
	}
	if h.stacks == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "compose stack service unavailable")
	}
	res, err := h.stacks.Save(c.Request().Context(), projectID, req)
	if err != nil {
		return utils.Error(c, http.StatusBadRequest, err.Error())
	}
	return utils.Success(c, "Stack saved", res)
}

func (h *ComposeStackHandler) Deploy(c echo.Context) error {
	projectID := c.Param("id")
	stackID := c.Param("stackId")
	if h.stacks == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "compose stack service unavailable")
	}
	res, err := h.stacks.Deploy(c.Request().Context(), projectID, stackID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Accepted(c, "Stack deployment started", res)
}

func (h *ComposeStackHandler) Cancel(c echo.Context) error {
	projectID := c.Param("id")
	stackID := c.Param("stackId")
	if h.stacks == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "compose stack service unavailable")
	}
	if err := h.stacks.Cancel(c.Request().Context(), projectID, stackID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Stack deployment cancelled", nil)
}

func (h *ComposeStackHandler) Restart(c echo.Context) error {
	projectID := c.Param("id")
	stackID := c.Param("stackId")
	if h.stacks == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "compose stack service unavailable")
	}
	if err := h.stacks.Restart(c.Request().Context(), projectID, stackID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Stack restarted", nil)
}

func (h *ComposeStackHandler) Stop(c echo.Context) error {
	projectID := c.Param("id")
	stackID := c.Param("stackId")
	if h.stacks == nil {
		return utils.Error(c, http.StatusServiceUnavailable, "compose stack service unavailable")
	}
	if err := h.stacks.Stop(c.Request().Context(), projectID, stackID); err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	return utils.Success(c, "Stack stopped", nil)
}
