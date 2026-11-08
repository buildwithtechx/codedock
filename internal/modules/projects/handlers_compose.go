package projects

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ComposeStackService interface {
	List(ctx context.Context, projectID string) (any, error)
	Config(ctx context.Context, projectID, stackID string) (any, error)
	Review(ctx context.Context, projectID string, req any) (any, error)
	Save(ctx context.Context, projectID string, req any) (any, error)
	Deploy(ctx context.Context, projectID, stackID string) (any, error)
	Cancel(ctx context.Context, projectID, stackID string) error
	Restart(ctx context.Context, projectID, stackID string) error
	Stop(ctx context.Context, projectID, stackID string) error
}

type ComposeStackHandler struct {
	stacks ComposeStackService
}

func NewComposeStackHandler(stacks ComposeStackService) *ComposeStackHandler {
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

// Deprecated direct compose handlers
type ComposeHandler struct{}

func NewComposeHandler() *ComposeHandler { return &ComposeHandler{} }

func (h *ComposeHandler) Analyze(c echo.Context) error {
	return utils.Error(c, http.StatusGone, "Use the project-scoped stack review API; direct analysis discards runtime fields.")
}

func (h *ComposeHandler) Deploy(c echo.Context) error {
	return utils.Error(c, http.StatusGone, "Partial Compose import has been replaced. Open the project's Compose page to review and save a complete stack before deployment.")
}
