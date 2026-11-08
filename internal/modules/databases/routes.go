package databases

import "github.com/labstack/echo/v4"

func RegisterRoutes(authGroup *echo.Group, h *Handler, requireScope func(string) echo.MiddlewareFunc) {
	if h == nil {
		h = &Handler{}
	}
	if authGroup == nil {
		return
	}
	authGroup.GET("/databases", h.ListDatabases, requireScope("database:manage"))
	authGroup.POST("/databases", h.CreateDatabase, requireScope("database:manage"))
	authGroup.GET("/databases/:id", h.GetDatabase, requireScope("database:manage"))
	authGroup.PUT("/databases/:id", h.UpdateDatabase, requireScope("database:manage"))
	authGroup.DELETE("/databases/:id", h.DeleteDatabase, requireScope("database:manage"))
	authGroup.POST("/databases/:id/start", h.StartDatabase, requireScope("database:manage"))
	authGroup.POST("/databases/:id/stop", h.StopDatabase, requireScope("database:manage"))
	authGroup.POST("/databases/:id/restart", h.RestartDatabase, requireScope("database:manage"))
	authGroup.POST("/databases/:id/credentials/reveal", h.RevealCredentials, requireScope("database:manage"))
	authGroup.POST("/databases/:id/query", h.QueryDatabase, requireScope("database:manage"))
	authGroup.POST("/databases/:id/import", h.ImportData, requireScope("database:manage"))

	authGroup.GET("/databases/:id/schemas", h.GetSchemas, requireScope("database:manage"))
	authGroup.GET("/databases/:id/redis/keys", h.GetRedisKeys, requireScope("database:manage"))
	authGroup.DELETE("/databases/:id/redis/keys/:key", h.DeleteRedisKey, requireScope("database:manage"))
	authGroup.GET("/databases/:id/data/:table", h.GetTableData, requireScope("database:manage"))
	authGroup.POST("/databases/:id/data/:table", h.InsertTableRow, requireScope("database:manage"))
	authGroup.PUT("/databases/:id/data/:table", h.UpdateTableRow, requireScope("database:manage"))
	authGroup.DELETE("/databases/:id/data/:table", h.DeleteTableRow, requireScope("database:manage"))
}
