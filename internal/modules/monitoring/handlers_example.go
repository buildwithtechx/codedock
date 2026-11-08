package monitoring

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/utils"
)

type ExampleHandler struct {
	service ExampleService
}

func NewExampleHandler(s ExampleService) *ExampleHandler {
	return &ExampleHandler{service: s}
}

func (h *ExampleHandler) List(c echo.Context) error {
	if h.service == nil {
		return utils.Success(c, "Available examples", []any{})
	}
	examples, err := h.service.ListExamples()
	if err != nil {
		return utils.Error(c, 500, "failed to list examples")
	}
	return utils.Success(c, "Available examples", examples)
}
