package networks

import "github.com/labstack/echo/v4"

// NetworksController handles server-defined network lifecycles.
type NetworksController interface {
	List(c echo.Context) error
	Create(c echo.Context) error
	Get(c echo.Context) error
	Delete(c echo.Context) error
	ReviseAccess(c echo.Context) error
	Commit(c echo.Context) error
	Rollback(c echo.Context) error
}
