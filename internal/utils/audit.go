package utils

import (
	"github.com/labstack/echo/v4"

	"codedock/internal/models"
)

func AuditActorID(c echo.Context) string {
	if claims, ok := c.Get("user").(*models.UserClaims); ok && claims != nil && claims.UserID != "" {
		return claims.UserID
	}
	return "system"
}
