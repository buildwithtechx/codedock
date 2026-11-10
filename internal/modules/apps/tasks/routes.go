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
}
