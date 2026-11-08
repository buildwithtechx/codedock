package audit

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, auditHandler *AuditLogHandler, requireAdmin echo.MiddlewareFunc) {
	if auditHandler == nil {
		auditHandler = &AuditLogHandler{}
	}
	authGroup.GET("/audit-logs", auditHandler.List, requireAdmin)
	authGroup.GET("/audit-logs/facets", auditHandler.Facets, requireAdmin)
	authGroup.GET("/audit-logs/settings", auditHandler.GetSettings, requireAdmin)
	authGroup.PUT("/audit-logs/settings", auditHandler.UpdateSettings, requireAdmin)
}
