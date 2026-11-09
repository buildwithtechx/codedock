package mcp

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type RoutesConfig struct {
	StreamableHTTP echo.HandlerFunc
	Catalog        echo.HandlerFunc
	SSE            echo.HandlerFunc
	Message        echo.HandlerFunc
}

func RegisterRoutes(groups []*echo.Group, cfg RoutesConfig) {
	for _, grp := range groups {
		if grp == nil {
			continue
		}
		if cfg.StreamableHTTP != nil {
			grp.POST("/mcp", cfg.StreamableHTTP)
			grp.GET("/mcp", func(c echo.Context) error {
				return c.JSON(http.StatusMethodNotAllowed, map[string]any{
					"error": "Streamable-HTTP MCP endpoint accepts POST JSON-RPC 2.0 requests; use /api/mcp/sse for SSE streams.",
				})
			})
		}
		if cfg.Catalog != nil {
			grp.GET("/mcp/catalog", cfg.Catalog)
		}
		if cfg.SSE != nil {
			grp.GET("/mcp/sse", cfg.SSE)
		}
		if cfg.Message != nil {
			grp.POST("/mcp/messages", cfg.Message)
			grp.POST("/mcp/message", cfg.Message)
		}
	}
}

func HandleOAuthProtectedResource(c echo.Context) error {
	base := c.Scheme() + "://" + c.Request().Host
	return c.JSON(http.StatusOK, map[string]any{
		"resource":                 base + "/api/mcp",
		"authorization_servers":    []string{base},
		"bearer_methods_supported": []string{"header"},
		"scopes_supported":         []string{"read", "write", "admin"},
	})
}

func RegisterOAuthProtectedResourceRoutes(e *echo.Echo) {
	if e == nil {
		return
	}
	e.GET("/.well-known/oauth-protected-resource", HandleOAuthProtectedResource)
	e.GET("/.well-known/oauth-protected-resource/*", HandleOAuthProtectedResource)
}
