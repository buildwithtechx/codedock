package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/http/apischema"
)

// ServeOpenAPISpec renders and serves the OpenAPI spec in JSON format.
func ServeOpenAPISpec(c echo.Context) error {
	rendered, err := apischema.Marshal()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSONBlob(http.StatusOK, rendered)
}

var serveOpenAPISpec = ServeOpenAPISpec
