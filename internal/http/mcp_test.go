package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"codedock/internal/http/mcp"
	"codedock/internal/models"

	"github.com/labstack/echo/v4"
)

func TestMCPToolGenerationAndSafety(t *testing.T) {
	tools := mcp.AllTools()
	if len(tools) < 50 {
		t.Fatalf("expected at least 50 route-backed MCP tools, got %d", len(tools))
	}

	var hasGetProjects, hasRedeployApp bool
	for _, tool := range tools {
		if strings.HasPrefix(tool.Path, "/api/auth") || strings.HasPrefix(tool.Path, "/api/profile/tokens") || strings.HasPrefix(tool.Path, "/api/mcp") {
			t.Fatalf("forbidden route exposed as MCP tool: %s (%s)", tool.Name, tool.Path)
		}
		if tool.Name == "get_projects" {
			hasGetProjects = true
			if !tool.Annotations.ReadOnlyHint {
				t.Fatalf("expected get_projects to be marked readOnlyHint=true")
			}
		}
		if tool.Name == "post_apps_by_id_redeploy" {
			hasRedeployApp = true
			if tool.Annotations.ReadOnlyHint {
				t.Fatalf("expected post_apps_by_id_redeploy to have readOnlyHint=false")
			}
		}
	}
	if !hasGetProjects || !hasRedeployApp {
		t.Fatalf("expected route-backed tools get_projects and post_apps_by_id_redeploy to exist")
	}

	readOnlyClaims := &models.UserClaims{Role: "readonly"}
	filtered := mcp.FilterToolsForClaims(tools, readOnlyClaims, false)
	for _, ft := range filtered {
		if !ft.Annotations.ReadOnlyHint {
			t.Fatalf("readonly claims must not see mutating tool %s", ft.Name)
		}
	}
}

func TestMCPStreamableHTTPProtocol(t *testing.T) {
	e := echo.New()
	e.GET("/api/system/stats", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"status": "ok", "cpu": 12})
	})

	bridge := mcp.NewBridge(nil, nil)
	bridge.SetDispatcher(e)
	srv := &Server{router: e, mcpBridge: bridge}

	initBody := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewBufferString(initBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := srv.handleMCPStreamableHTTP(c); err != nil {
		t.Fatalf("handleMCPStreamableHTTP initialize failed: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"protocolVersion":"2025-06-18"`) {
		t.Fatalf("unexpected initialize response: %s", rec.Body.String())
	}

	callBody := `{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_system_stats","arguments":{}}}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/mcp", bytes.NewBufferString(callBody))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	c2 := e.NewContext(req2, rec2)

	if err := srv.handleMCPStreamableHTTP(c2); err != nil {
		t.Fatalf("handleMCPStreamableHTTP tools/call failed: %v", err)
	}
	var rpcResp map[string]any
	if err := json.Unmarshal(rec2.Body.Bytes(), &rpcResp); err != nil {
		t.Fatalf("failed to decode jsonrpc response: %v", err)
	}
	if !strings.Contains(rec2.Body.String(), `\"cpu\":12`) {
		t.Fatalf("expected dispatched tool output containing cpu:12, got: %s", rec2.Body.String())
	}
}
