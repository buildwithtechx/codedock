package tasks

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, h *ScheduledTaskHandler) {
	if authGroup == nil {
		return
	}
	if h == nil {
		h = &ScheduledTaskHandler{}
	}

	authGroup.GET("/scheduled-tasks", h.ListProjectScheduledTasks)
	authGroup.POST("/scheduled-tasks", h.Create)
	authGroup.GET("/scheduled-tasks/:id", h.Get)
	authGroup.PUT("/scheduled-tasks/:id", h.Update)
	authGroup.DELETE("/scheduled-tasks/:id", h.Delete)
	authGroup.POST("/scheduled-tasks/:id/trigger", h.Run)

	// Legacy /jobs alias routes
	authGroup.GET("/jobs", h.ListProjectScheduledTasks)
	authGroup.POST("/jobs", h.Create)
	authGroup.GET("/jobs/:id", h.Get)
	authGroup.PUT("/jobs/:id", h.Update)
	authGroup.DELETE("/jobs/:id", h.Delete)
	authGroup.POST("/jobs/:id/trigger", h.Run)
}
