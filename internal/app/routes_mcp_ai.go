package app

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/http/middleware"
	"codedock/internal/modules/ai"
	"codedock/internal/modules/mcp"
	"codedock/internal/utils"
	"codedock/internal/version"
)

func (a *App) getMeta(c echo.Context) error {
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

func (a *App) registerMcpAndAiRoutes(authGroup *echo.Group, s *appServices) {
	if a.mcpBridge != nil {
		a.mcpBridge.SetDispatcher(a.router)
	}

	mcp.RegisterOAuthProtectedResourceRoutes(a.router)

	v1Group := a.router.Group("/api/v1", a.authGuard.RequireAuth())
	groups := []*echo.Group{authGroup, v1Group}

	mcp.RegisterRoutes(groups, mcp.RoutesConfig{
		StreamableHTTP: a.handleMCPStreamableHTTP,
		Catalog:        a.handleMCPCatalog,
		SSE:            a.HandleMCPSSE,
		Message:        a.HandleMCPMessage,
	})

	aiHandler := ai.NewHandler(s.aiSettingsService)
	ai.RegisterAssistantRoutes(groups, ai.AssistantRoutesConfig{
		ListResources: func(c echo.Context) error {
			return a.handleListAiResources(c, s, aiHandler)
		},
		Chat: func(c echo.Context) error {
			return a.handleAiChat(c, s, aiHandler)
		},
		ConfirmAction: func(c echo.Context) error {
			return a.handleConfirmAiAction(c, s, aiHandler)
		},
	})
}

func (a *App) handleMCPStreamableHTTP(c echo.Context) error {
	claims := middleware.GetUserClaimsFromContext(c.Request().Context())
	return a.mcpBridge.HandleStreamableHTTP(c, claims)
}

func (a *App) handleMCPCatalog(c echo.Context) error {
	claims := middleware.GetUserClaimsFromContext(c.Request().Context())
	cloudEnabled := config.Get().Cloud.Enabled
	tools := mcp.FilterToolsForClaims(a.mcpBridge.Tools(), claims, cloudEnabled)
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

func (a *App) aiDeps(c echo.Context, s *appServices, aiHandler *ai.Handler) ai.Dependencies {
	deps := ai.Dependencies{
		Bridge:         a.mcpBridge,
		ProjectService: a.projectService,
		AppService:     a.appService,
	}
	if a.mcpBridge != nil {
		deps.DBService = a.mcpBridge.DatabaseService()
	}
	if aiHandler != nil {
		deps.Settings, _ = aiHandler.GetRawSettings(c.Request().Context())
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

func (a *App) handleListAiResources(c echo.Context, s *appServices, aiHandler *ai.Handler) error {
	items := ai.ListResources(c.Request().Context(), a.aiDeps(c, s, aiHandler))
	return utils.Success(c, "Resources loaded", items)
}

func (a *App) handleAiChat(c echo.Context, s *appServices, aiHandler *ai.Handler) error {
	var req ai.ChatRequest
	if err := c.Bind(&req); err != nil || len(req.Messages) == 0 {
		return utils.Error(c, http.StatusBadRequest, "Valid messages are required")
	}
	resp, err := ai.HandleChat(c.Request().Context(), req, dispatchMetaFromEcho(c), a.aiDeps(c, s, aiHandler))
	if err != nil {
		return utils.Error(c, http.StatusBadGateway, err.Error())
	}
	return utils.Success(c, "OK", resp)
}

func (a *App) handleConfirmAiAction(c echo.Context, s *appServices, aiHandler *ai.Handler) error {
	var body struct {
		ActionID string `json:"actionId"`
	}
	if err := c.Bind(&body); err != nil || body.ActionID == "" {
		return utils.Error(c, http.StatusBadRequest, "actionId is required")
	}
	res, err := ai.ConfirmAction(c.Request().Context(), body.ActionID, dispatchMetaFromEcho(c), a.aiDeps(c, s, aiHandler))
	if err != nil {
		return utils.Error(c, http.StatusNotFound, err.Error())
	}
	return utils.Success(c, "Action executed", res)
}
