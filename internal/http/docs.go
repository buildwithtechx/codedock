package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/http/apischema"
)

func serveOpenAPISpec(c echo.Context) error {
	rendered, err := apischema.Marshal()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSONBlob(http.StatusOK, rendered)
}
