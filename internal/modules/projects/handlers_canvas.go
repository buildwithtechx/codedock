package projects

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type CanvasService interface {
	ListSummaries(ctx context.Context, organizationID string) (any, error)
	GetSummary(ctx context.Context, id string) (any, error)
	GetEnvironmentCanvas(ctx context.Context, id string) (any, error)
	ApplyTopology(ctx context.Context, envID string, req any) error
	ReadObservation(ctx context.Context, canvas any, node string) (any, error)
}

type CanvasHandler struct {
	canvas  CanvasService
	project *Service
}

func NewCanvasHandler(cs CanvasService, ps *Service) *CanvasHandler {
	return &CanvasHandler{canvas: cs, project: ps}
}

func (h *CanvasHandler) ListCanvasSummaries(c echo.Context) error {
	organizationID := c.QueryParam("organizationId")
	if organizationID == "" {
		return utils.Error(c, http.StatusBadRequest, "organizationId is required")
	}
	if h.canvas == nil {
		return utils.Success(c, "Operation successful", []any{})
	}
	summaries, err := h.canvas.ListSummaries(c.Request().Context(), organizationID)
	if err != nil {
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if summaries == nil {
		summaries = []any{}
	}
	return utils.Success(c, "Operation successful", summaries)
}

func (h *CanvasHandler) GetCanvasSummary(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if h.canvas == nil {
		return utils.Error(c, http.StatusNotFound, "canvas summary not found")
	}
	summary, err := h.canvas.GetSummary(c.Request().Context(), id)
	if err != nil {
		if utils.IsNotFound(err) {
			return utils.Error(c, http.StatusNotFound, "canvas summary not found")
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if summary == nil {
		return utils.Error(c, http.StatusNotFound, "canvas summary not found")
	}
	return utils.Success(c, "Operation successful", summary)
}

func (h *CanvasHandler) GetEnvironmentCanvas(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if h.canvas == nil {
		return utils.Error(c, http.StatusNotFound, "environment canvas not found")
	}
	canvas, err := h.canvas.GetEnvironmentCanvas(c.Request().Context(), id)
	if err != nil {
		if utils.IsNotFound(err) {
			return utils.Error(c, http.StatusNotFound, "environment canvas not found")
		}
		return utils.Error(c, http.StatusInternalServerError, err.Error())
	}
	if canvas == nil {
		return utils.Error(c, http.StatusNotFound, "environment canvas not found")
	}
	return utils.Success(c, "Operation successful", canvas)
}

func (h *CanvasHandler) ApplyTopology(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if h.canvas == nil {
		return utils.Error(c, http.StatusNotFound, "environment not found")
	}
	var req any
	if err := c.Bind(&req); err != nil {
		return utils.Error(c, http.StatusBadRequest, "invalid topology request")
	}
	if err := h.canvas.ApplyTopology(c.Request().Context(), id, req); err != nil {
		return utils.Error(c, http.StatusConflict, err.Error())
	}
	return utils.Success(c, "Dependencies saved", nil)
}

func (h *CanvasHandler) ObserveResource(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.Error(c, http.StatusBadRequest, "missing id parameter")
	}
	if h.canvas == nil {
		return utils.Error(c, http.StatusNotFound, "environment unavailable")
	}
	canvas, err := h.canvas.GetEnvironmentCanvas(c.Request().Context(), id)
	if err != nil {
		return utils.Error(c, http.StatusNotFound, "environment unavailable")
	}
	obs, err := h.canvas.ReadObservation(c.Request().Context(), canvas, c.QueryParam("node"))
	if err != nil {
		return utils.Error(c, http.StatusServiceUnavailable, err.Error())
	}
	return utils.Success(c, "Runtime observation retrieved", obs)
}
