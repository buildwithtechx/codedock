package http

import (
	"fmt"
	"net/http"

	"codedock/internal/config"
	"codedock/internal/http/ai"
	"codedock/internal/http/mcp"
	"codedock/internal/utils"
	"codedock/internal/version"

	"github.com/labstack/echo/v4"
)

func (s *Server) getMeta(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"name":    "codedock",
		"version": "v1.0.0",
		"capabilities": map[string]bool{
			"sseLogs": true,
			"push":    true,
			"actions": true,
			"apps":    true,
			"mcp":     true,
			"ai":      true,
		},
	})
}

func (s *Server) registerMcpAndAiRoutes(authGroup *echo.Group) {
	s.mcpBridge.SetDispatcher(s.router)

	// OAuth 2.1 Protected Resource Discovery for MCP clients
	s.router.GET("/.well-known/oauth-protected-resource", s.handleOAuthProtectedResource)
	s.router.GET("/.well-known/oauth-protected-resource/*", s.handleOAuthProtectedResource)

	// Stateless Streamable-HTTP JSON-RPC 2.0 (POST /api/mcp & /api/v1/mcp)
	authGroup.POST("/mcp", s.handleMCPStreamableHTTP)
	authGroup.GET("/mcp", func(c echo.Context) error {
		return c.JSON(http.StatusMethodNotAllowed, map[string]any{
			"error": "Streamable-HTTP MCP endpoint accepts POST JSON-RPC 2.0 requests; use /api/mcp/sse for SSE streams.",
		})
	})
	authGroup.GET("/mcp/catalog", s.handleMCPCatalog)

	// Persistent SSE MCP endpoints (/api/mcp/sse + /api/v1/mcp/sse)
	authGroup.GET("/mcp/sse", s.HandleMCPSSE)
	authGroup.POST("/mcp/messages", s.HandleMCPMessage)
	authGroup.POST("/mcp/message", s.HandleMCPMessage)

	v1Group := s.router.Group("/api/v1", s.authGuard.RequireAuth())
	v1Group.POST("/mcp", s.handleMCPStreamableHTTP)
	v1Group.GET("/mcp/catalog", s.handleMCPCatalog)
	v1Group.GET("/mcp/sse", s.HandleMCPSSE)
	v1Group.POST("/mcp/messages", s.HandleMCPMessage)
	v1Group.POST("/mcp/message", s.HandleMCPMessage)

	// Instance-Aware AI Co-Pilot endpoints
	authGroup.GET("/ai/resources", s.handleListAiResources)
	authGroup.POST("/ai/chat", s.handleAiChat)
	authGroup.POST("/ai/actions/confirm", s.handleConfirmAiAction)
	v1Group.GET("/ai/resources", s.handleListAiResources)
	v1Group.POST("/ai/chat", s.handleAiChat)
	v1Group.POST("/ai/actions/confirm", s.handleConfirmAiAction)
}

func (s *Server) handleOAuthProtectedResource(c echo.Context) error {
	base := fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host)
	return c.JSON(http.StatusOK, map[string]any{
		"resource":                 base + "/api/mcp",
		"authorization_servers":    []string{base},
		"bearer_methods_supported": []string{"header"},
		"scopes_supported":         []string{"read", "write", "admin"},
	})
}

func (s *Server) handleMCPStreamableHTTP(c echo.Context) error {
	claims := GetUserClaimsFromContext(c.Request().Context())
	return s.mcpBridge.HandleStreamableHTTP(c, claims)
}

func (s *Server) handleMCPCatalog(c echo.Context) error {
	claims := GetUserClaimsFromContext(c.Request().Context())
	cloudEnabled := config.Get().Cloud.Enabled
	tools := mcp.FilterToolsForClaims(s.mcpBridge.Tools(), claims, cloudEnabled)
	return utils.Success(c, "MCP catalog loaded", map[string]any{
		"serverInfo": map[string]any{
			"name":     "codedock-mcp",
			"version":  version.Version,
			"protocol": "2025-06-18",
		},
		"tools":     tools,
		"prompts":   mcp.ListPrompts(),
		"resources": mcp.ListResources(),
	})
}

func (s *Server) aiDeps(c echo.Context) ai.Dependencies {
	deps := ai.Dependencies{
		Bridge:         s.mcpBridge,
		ProjectService: s.projectService,
		AppService:     s.appService,
	}
	if s.mcpBridge != nil {
		deps.DBService = s.mcpBridge.DatabaseService()
	}
	if s.aiSettingsHandler != nil {
		deps.Settings, _ = s.aiSettingsHandler.GetRawSettings(c.Request().Context())
	}
	return deps
}

func dispatchMetaFromEcho(c echo.Context) mcp.DispatchMeta {
	return mcp.DispatchMeta{
		AuthHeader: c.Request().Header.Get("Authorization"),
		Cookies:    c.Request().Cookies(),
		OrgID:      c.Request().Header.Get("X-Organization-Id"),
		ClientIP:   c.RealIP(),
		UserAgent:  c.Request().UserAgent(),
	}
}

func (s *Server) handleListAiResources(c echo.Context) error {
	items := ai.ListResources(c.Request().Context(), s.aiDeps(c))
	return utils.Success(c, "Resources loaded", items)
}

func (s *Server) handleAiChat(c echo.Context) error {
	var req ai.ChatRequest
	if err := c.Bind(&req); err != nil || len(req.Messages) == 0 {
		return utils.Error(c, http.StatusBadRequest, "Valid messages are required")
	}
	resp, err := ai.HandleChat(c.Request().Context(), req, dispatchMetaFromEcho(c), s.aiDeps(c))
	if err != nil {
		return utils.Error(c, http.StatusBadGateway, err.Error())
	}
	return utils.Success(c, "OK", resp)
}

func (s *Server) handleConfirmAiAction(c echo.Context) error {
	var body struct {
		ActionID string `json:"actionId"`
	}
	if err := c.Bind(&body); err != nil || body.ActionID == "" {
		return utils.Error(c, http.StatusBadRequest, "actionId is required")
	}
	res, err := ai.ConfirmAction(c.Request().Context(), body.ActionID, dispatchMetaFromEcho(c), s.aiDeps(c))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Action executed", res)
}
