package actions

import "github.com/labstack/echo/v4"

// ActionsController exposes all action workflow, run, and runner endpoints.
type ActionsController interface {
	ListWorkflows(c echo.Context) error
	CreateWorkflow(c echo.Context) error
	DispatchWorkflow(c echo.Context) error
	DeleteWorkflow(c echo.Context) error

	ListRuns(c echo.Context) error
	GetRun(c echo.Context) error
	RetryRun(c echo.Context) error
	CancelRun(c echo.Context) error
	DeleteRun(c echo.Context) error

	ListRunners(c echo.Context) error
	InspectDestination(c echo.Context) error
	CreateRunner(c echo.Context) error
	DeleteRunner(c echo.Context) error
}
