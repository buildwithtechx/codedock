package organizations

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, orgHandler *OrganizationHandler, requireMember, requireAdmin, requireOwner echo.MiddlewareFunc) {
	if orgHandler == nil {
		orgHandler = &OrganizationHandler{}
	}
	authGroup.GET("/organizations", orgHandler.List)
	authGroup.POST("/organizations", orgHandler.Create)
	authGroup.GET("/organizations/:id", orgHandler.Get, requireMember)
	authGroup.PUT("/organizations/:id", orgHandler.Update, requireAdmin)
	authGroup.DELETE("/organizations/:id", orgHandler.Delete, requireOwner)
	authGroup.GET("/organizations/:id/members", orgHandler.ListMembers, requireMember)
	authGroup.POST("/organizations/:id/members", orgHandler.InviteMember, requireAdmin)
	authGroup.PUT("/organizations/:id/members/:userId", orgHandler.UpdateMember, requireAdmin)
	authGroup.DELETE("/organizations/:id/members/:memberId", orgHandler.RemoveMember, requireAdmin)
}
