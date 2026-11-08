package users

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, userHandler *UserHandler, adminInvite echo.HandlerFunc, requireAdmin, otpRateLimit echo.MiddlewareFunc) {
	if userHandler == nil {
		userHandler = &UserHandler{}
	}
	authGroup.GET("/users", userHandler.ListUsers, requireAdmin)
	authGroup.POST("/users/invite", adminInvite, requireAdmin)
	authGroup.DELETE("/users/:id", userHandler.DeleteUser, requireAdmin)
	authGroup.GET("/profile", userHandler.GetProfile)
	authGroup.PUT("/profile", userHandler.UpdateProfile)
	authGroup.POST("/profile/email/request", userHandler.RequestEmailChange)
	authGroup.POST("/profile/email/verify", userHandler.VerifyEmailChange, otpRateLimit)
	authGroup.PUT("/profile/password", userHandler.ChangePassword)
	authGroup.GET("/profile/tokens", userHandler.ListPATs)
	authGroup.POST("/profile/tokens", userHandler.CreatePAT)
	authGroup.DELETE("/profile/tokens/:id", userHandler.DeletePAT)
}
