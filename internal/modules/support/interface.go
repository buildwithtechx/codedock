package support

import "github.com/labstack/echo/v4"

// SupportController defines the support HTTP endpoints contract.
type SupportController interface {
	ListTickets(c echo.Context) error
	CreateTicket(c echo.Context) error
	GetTicket(c echo.Context) error
	ReplyTicket(c echo.Context) error
	CloseTicket(c echo.Context) error
}
