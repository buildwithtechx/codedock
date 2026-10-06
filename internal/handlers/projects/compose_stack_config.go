package projects

import (
	"codedock.run/codedock/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *ComposeStackHandler) Config(c echo.Context) error {
	if err := h.authorize(c, true); err != nil {
		return err
	}
	stack, err := h.stacks.Get(c.Request().Context(), c.Param("id"), c.Param("stackId"))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "stack not found")
	}
	return utils.Success(c, "Stack configuration retrieved", map[string]any{"stack": stack, "content": stack.Config})
}
