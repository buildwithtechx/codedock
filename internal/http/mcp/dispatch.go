package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"codedock/internal/config"
	"codedock/internal/models"
	"codedock/internal/version"
	"github.com/labstack/echo/v4"
)

type DispatchMeta struct {
	AuthHeader string
	Cookies    []*http.Cookie
	OrgID      string
	ClientIP   string
	UserAgent  string
}

type ToolDispatchResult struct {
	Status int
	OK     bool
	Text   string
	Data   any
}

func resolvePathParam(args map[string]any, param string) string {
	if v, ok := args[param]; ok && v != nil && fmt.Sprint(v) != "" {
		return fmt.Sprint(v)
	}
	return ""
}

func (b *Bridge) DispatchTool(ctx context.Context, tool ToolDef, args map[string]any, meta DispatchMeta) ToolDispatchResult {
	if args == nil {
		args = map[string]any{}
	}
	if b.dispatcher == nil {
		return ToolDispatchResult{Status: 503, OK: false, Text: `{"error":"HTTP dispatcher unavailable"}`}
	}

	path := tool.Path
	for _, param := range tool.PathParams {
		val := resolvePathParam(args, param)
		if val == "" || val == "." || val == ".." {
			msg := fmt.Sprintf(`{"error":"Missing or invalid path parameter: %s","code":"INVALID_TOOL_ARGUMENTS"}`, param)
			return ToolDispatchResult{Status: 400, OK: false, Text: msg}
		}
		path = strings.ReplaceAll(path, ":"+param, url.PathEscape(val))
		path = strings.ReplaceAll(path, "{"+param+"}", url.PathEscape(val))
	}

	u := &url.URL{Scheme: "http", Host: "mcp.internal", Path: path}
	q := u.Query()
	if queryMap, ok := args["query"].(map[string]any); ok {
		for k, v := range queryMap {
			if v != nil {
				q.Set(k, fmt.Sprint(v))
			}
		}
	}
	u.RawQuery = q.Encode()

	var bodyReader *bytes.Reader
	if tool.HasBody || tool.Method == "POST" || tool.Method == "PUT" || tool.Method == "PATCH" {
		bodyPayload := args["body"]
		if bodyPayload == nil {
			bodyPayload = map[string]any{}
		}
		raw, _ := json.Marshal(bodyPayload)
		bodyReader = bytes.NewReader(raw)
	} else {
		bodyReader = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, tool.Method, u.String(), bodyReader)
	if err != nil {
		return ToolDispatchResult{Status: 500, OK: false, Text: fmt.Sprintf(`{"error":%q}`, err.Error())}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Codedock-Source", "mcp")
	if meta.AuthHeader != "" {
		req.Header.Set("Authorization", meta.AuthHeader)
	}
	for _, ck := range meta.Cookies {
		req.AddCookie(ck)
	}
	orgID := meta.OrgID
	if argOrg, ok := args["organizationId"].(string); ok && strings.TrimSpace(argOrg) != "" {
		orgID = strings.TrimSpace(argOrg)
	}
	if orgID != "" {
		req.Header.Set("X-Organization-Id", orgID)
	}
	if meta.ClientIP != "" {
		req.Header.Set("X-Real-IP", meta.ClientIP)
	}
	if meta.UserAgent != "" {
		req.Header.Set("User-Agent", meta.UserAgent)
	}

	rec := httptest.NewRecorder()
	b.dispatcher.ServeHTTP(rec, req)

	if strings.HasPrefix(rec.Header().Get("Content-Type"), "text/event-stream") {
		msg := `{"error":"Endpoint returns a live stream; read the JSON status/logs endpoint instead.","code":"MCP_STREAM_NOT_SUPPORTED"}`
		return ToolDispatchResult{Status: 502, OK: false, Text: msg}
	}

	text := strings.TrimSpace(rec.Body.String())
	if text == "" {
		text = fmt.Sprintf(`{"status":%d}`, rec.Code)
	}
	var parsed any
	_ = json.Unmarshal([]byte(text), &parsed)
	return ToolDispatchResult{
		Status: rec.Code,
		OK:     rec.Code >= 200 && rec.Code < 300,
		Text:   text,
		Data:   parsed,
	}
}

type jsonRpcRequest struct {
	JSONRPC string         `json:"jsonrpc"`
	ID      any            `json:"id"`
	Method  string         `json:"method"`
	Params  map[string]any `json:"params"`
}

func jsonRpcErr(id any, code int, message string) map[string]any {
	return map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"error":   map[string]any{"code": code, "message": message},
	}
}

func (b *Bridge) HandleStreamableHTTP(c echo.Context, claims *models.UserClaims) error {
	var raw json.RawMessage
	if err := json.NewDecoder(c.Request().Body).Decode(&raw); err != nil {
		return c.JSON(http.StatusOK, jsonRpcErr(nil, -32700, "Parse error"))
	}
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		return c.JSON(http.StatusOK, jsonRpcErr(nil, -32600, "JSON-RPC batching is not supported"))
	}
	var msg jsonRpcRequest
	if err := json.Unmarshal(trimmed, &msg); err != nil || msg.JSONRPC != "2.0" || msg.Method == "" {
		return c.JSON(http.StatusOK, jsonRpcErr(msg.ID, -32600, "Invalid Request"))
	}

	cloudEnabled := config.Get().Cloud.Enabled

	switch msg.Method {
	case "initialize":
		proto := "2025-06-18"
		if reqProto, ok := msg.Params["protocolVersion"].(string); ok {
			if reqProto == "2025-06-18" || reqProto == "2025-03-26" || reqProto == "2024-11-05" {
				proto = reqProto
			}
		}
		return c.JSON(http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      msg.ID,
			"result": map[string]any{
				"protocolVersion": proto,
				"capabilities": map[string]any{
					"tools":     map[string]any{"listChanged": false},
					"prompts":   map[string]any{"listChanged": false},
					"resources": map[string]any{"listChanged": false},
				},
				"serverInfo":   map[string]any{"name": "codedock-mcp", "version": version.Version},
				"instructions": "Codedock MCP server. Use tools/list to inspect available operations and prompts/list for guided multi-step workflows.",
			},
		})
	case "notifications/initialized", "initialized":
		return c.NoContent(http.StatusAccepted)
	case "ping":
		return c.JSON(http.StatusOK, map[string]any{"jsonrpc": "2.0", "id": msg.ID, "result": map[string]any{}})
	case "tools/list":
		visible := FilterToolsForClaims(b.Tools(), claims, cloudEnabled)
		return c.JSON(http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      msg.ID,
			"result":  map[string]any{"tools": visible},
		})
	case "tools/call":
		name, _ := msg.Params["name"].(string)
		tool, ok := b.FindTool(name)
		if !ok {
			return c.JSON(http.StatusOK, jsonRpcErr(msg.ID, -32602, "Unknown tool: "+name))
		}
		args, _ := msg.Params["arguments"].(map[string]any)
		meta := DispatchMeta{
			AuthHeader: c.Request().Header.Get("Authorization"),
			Cookies:    c.Request().Cookies(),
			OrgID:      c.Request().Header.Get("X-Organization-Id"),
			ClientIP:   c.RealIP(),
			UserAgent:  c.Request().UserAgent(),
		}
		res := b.DispatchTool(c.Request().Context(), tool, args, meta)
		return c.JSON(http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      msg.ID,
			"result": map[string]any{
				"content": []map[string]any{{"type": "text", "text": res.Text}},
				"isError": !res.OK,
			},
		})
	case "prompts/list":
		return c.JSON(http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      msg.ID,
			"result":  map[string]any{"prompts": ListPrompts()},
		})
	case "prompts/get":
		name, _ := msg.Params["name"].(string)
		rawArgs, _ := msg.Params["arguments"].(map[string]any)
		strArgs := map[string]string{}
		for k, v := range rawArgs {
			strArgs[k] = fmt.Sprint(v)
		}
		prompt, ok := GetPrompt(name, strArgs)
		if !ok {
			return c.JSON(http.StatusOK, jsonRpcErr(msg.ID, -32602, "Unknown prompt: "+name))
		}
		return c.JSON(http.StatusOK, map[string]any{"jsonrpc": "2.0", "id": msg.ID, "result": prompt})
	case "resources/list":
		return c.JSON(http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      msg.ID,
			"result":  map[string]any{"resources": ListResources()},
		})
	case "resources/read":
		uri, _ := msg.Params["uri"].(string)
		content, ok := b.ReadResource(c.Request().Context(), uri)
		if !ok {
			return c.JSON(http.StatusOK, jsonRpcErr(msg.ID, -32602, "Unknown resource: "+uri))
		}
		return c.JSON(http.StatusOK, map[string]any{
			"jsonrpc": "2.0",
			"id":      msg.ID,
			"result": map[string]any{
				"contents": []map[string]any{{"uri": uri, "mimeType": "application/json", "text": content}},
			},
		})
	default:
		if msg.ID == nil {
			return c.NoContent(http.StatusAccepted)
		}
		return c.JSON(http.StatusOK, jsonRpcErr(msg.ID, -32601, "Method not found: "+msg.Method))
	}
}
