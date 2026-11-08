package http

import (
	mailmodule "codedock/internal/modules/mail"
	"github.com/labstack/echo/v4"
)

func (s *Server) registerMailRoutes(authGroup *echo.Group) {
	mailmodule.RegisterRoutes(authGroup, s.mailHandler)
}
