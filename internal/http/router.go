package http

import (
	"github.com/labstack/echo/v4"
)

// NewRouter constructs an Echo router configured with all global transport middleware.
func NewRouter() *echo.Echo {
	e := echo.New()
	configureEchoMiddleware(e)
	return e
}
