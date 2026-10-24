package projects

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *AppHandler) guardAppOperation(c echo.Context, id string) (func(), error) {
	if h.deployer == nil {
		return func() {}, nil
	}
	ctx, release, err := h.deployer.BeginServiceOperation(c.Request().Context(), id)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusConflict, err.Error())
	}
	c.SetRequest(c.Request().WithContext(ctx))
	return release, nil
}
