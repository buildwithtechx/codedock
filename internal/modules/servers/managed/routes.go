package managed

import (
	"github.com/labstack/echo/v4"
)

type RouteMiddleware struct {
	OrgMember   echo.MiddlewareFunc
	OrgAdmin    echo.MiddlewareFunc
	ServerRead  echo.MiddlewareFunc
	ServerWrite echo.MiddlewareFunc
}

func RegisterRoutes(authGroup *echo.Group, h *ManagedHandler, mw RouteMiddleware) {
	if h == nil {
		return
	}
	authGroup.GET("/managed/catalog", h.Catalog)
	authGroup.GET("/organizations/:id/managed/credentials", h.ListCredentials, mw.OrgMember)
	authGroup.POST("/organizations/:id/managed/credentials", h.CreateCredential, mw.OrgAdmin)
	authGroup.DELETE("/organizations/:id/managed/credentials/:credentialId", h.DeleteCredential, mw.OrgAdmin)
	authGroup.GET("/organizations/:id/managed/credentials/:credentialId/offers", h.Offers, mw.OrgMember, mw.ServerRead)
	authGroup.GET("/organizations/:id/managed/quota", h.GetQuota, mw.OrgMember)
	authGroup.PUT("/organizations/:id/managed/quota", h.SetQuota, mw.OrgAdmin)
	authGroup.POST("/organizations/:id/managed/servers/review", h.ReviewProvision, mw.OrgMember, mw.ServerWrite)
	authGroup.POST("/organizations/:id/managed/servers/:serverId/refresh", h.RefreshServer, mw.OrgMember, mw.ServerRead)
	authGroup.POST("/organizations/:id/managed/servers/:serverId/resize/review", h.ReviewResize, mw.OrgAdmin)
	authGroup.POST("/organizations/:id/managed/servers/:serverId/delete/review", h.ReviewDelete, mw.OrgAdmin)
	authGroup.GET("/managed-operations/:operationId", h.GetOperation, mw.ServerRead)
	authGroup.GET("/managed-operations/:operationId/stream", h.StreamOperation, mw.ServerRead)
	authGroup.POST("/managed-operations/:operationId/apply", h.ApplyOperation, mw.ServerWrite)
}
