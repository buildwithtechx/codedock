package mcp

import (
	"regexp"
	"strings"

	"codedock/internal/http/apischema"
	"codedock/internal/models"
)

type ToolAnnotations struct {
	ReadOnlyHint    bool `json:"readOnlyHint"`
	DestructiveHint bool `json:"destructiveHint"`
}

type ToolDef struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Method      string          `json:"method"`
	Path        string          `json:"path"`
	Category    string          `json:"category"`
	Auth        string          `json:"auth"`
	PathParams  []string        `json:"pathParams"`
	HasBody     bool            `json:"hasBody"`
	Annotations ToolAnnotations `json:"annotations"`
	InputSchema map[string]any  `json:"inputSchema"`
}

var nonAlphaNumRe = regexp.MustCompile(`[^a-zA-Z0-9]+`)

func routeToToolName(method, path string) string {
	trimmed := strings.TrimPrefix(path, "/api/")
	trimmed = strings.TrimPrefix(trimmed, "/")
	parts := strings.Split(trimmed, "/")
	tokens := []string{strings.ToLower(method)}
	for _, p := range parts {
		if p == "" {
			continue
		}
		if strings.HasPrefix(p, ":") {
			param := strings.TrimPrefix(p, ":")
			tokens = append(tokens, "by", nonAlphaNumRe.ReplaceAllString(param, "_"))
		} else if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			param := strings.TrimSuffix(strings.TrimPrefix(p, "{"), "}")
			tokens = append(tokens, "by", nonAlphaNumRe.ReplaceAllString(param, "_"))
		} else {
			clean := nonAlphaNumRe.ReplaceAllString(strings.ToLower(p), "_")
			tokens = append(tokens, strings.Trim(clean, "_"))
		}
	}
	return strings.Join(tokens, "_")
}

func isBlockedOperation(op apischema.Operation) bool {
	if op.Content == apischema.ContentSSE || op.Content == apischema.ContentBinary {
		return true
	}
	p := op.Path
	if strings.HasPrefix(p, "/api/auth") ||
		strings.HasPrefix(p, "/api/mcp") ||
		strings.HasPrefix(p, "/api/ws/") ||
		strings.HasPrefix(p, "/api/webhooks/") ||
		strings.HasPrefix(p, "/api/profile/tokens") ||
		strings.HasPrefix(p, "/api/profile/password") ||
		strings.HasPrefix(p, "/api/profile/email") ||
		strings.HasSuffix(p, "/stream") ||
		strings.HasSuffix(p, "/exec-terminal") ||
		strings.HasSuffix(p, "/download") {
		return true
	}
	return false
}

func isDestructiveOperation(method, path string) bool {
	if method == "DELETE" {
		return true
	}
	lower := strings.ToLower(path)
	for _, suffix := range []string{"/stop", "/cancel", "/rollback", "/empty", "/cleanup", "/restart", "/reset"} {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

func schemaToJSONMap(s apischema.Schema) map[string]any {
	out := map[string]any{}
	if s.Type != "" {
		out["type"] = s.Type
	} else {
		out["type"] = "object"
	}
	if s.Description != "" {
		out["description"] = s.Description
	}
	if len(s.Enum) > 0 {
		out["enum"] = s.Enum
	}
	if s.Items != nil {
		out["items"] = schemaToJSONMap(*s.Items)
	}
	if len(s.Properties) > 0 {
		props := map[string]any{}
		for _, prop := range s.Properties {
			props[prop.Name] = schemaToJSONMap(prop.Schema)
		}
		out["properties"] = props
	}
	if len(s.Required) > 0 {
		out["required"] = s.Required
	}
	return out
}

func buildOperationInputSchema(op apischema.Operation) (map[string]any, []string, bool) {
	props := map[string]any{
		"organizationId": map[string]any{
			"type":        "string",
			"description": "Optional workspace/organization ID scope for this call.",
		},
	}
	var required []string
	var pathParams []string

	for _, pp := range op.PathParams {
		pathParams = append(pathParams, pp.Name)
		t := pp.Type
		if t == "" {
			t = "string"
		}
		props[pp.Name] = map[string]any{
			"type":        t,
			"description": pp.Description,
		}
		required = append(required, pp.Name)
	}

	if len(op.Query) > 0 {
		qProps := map[string]any{}
		for _, qp := range op.Query {
			qt := qp.Type
			if qt == "" {
				qt = "string"
			}
			qProps[qp.Name] = map[string]any{
				"type":        qt,
				"description": qp.Description,
			}
		}
		props["query"] = map[string]any{
			"type":        "object",
			"description": "Query string parameters.",
			"properties":  qProps,
		}
	}

	hasBody := op.Request != nil && op.Request.Content == "application/json"
	if hasBody {
		props["body"] = schemaToJSONMap(op.Request.Schema)
	}

	schema := map[string]any{
		"type":       "object",
		"properties": props,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema, pathParams, hasBody
}

func AllTools() []ToolDef {
	ops := apischema.All()
	seen := map[string]bool{}
	var tools []ToolDef

	for _, op := range ops {
		if isBlockedOperation(op) {
			continue
		}
		name := routeToToolName(op.Method, op.Path)
		if seen[name] {
			continue
		}
		seen[name] = true

		category := "General"
		if len(op.Tags) > 0 && op.Tags[0] != "" {
			category = op.Tags[0]
		}
		inputSchema, pathParams, hasBody := buildOperationInputSchema(op)
		desc := op.Summary
		if desc == "" {
			desc = op.Method + " " + op.Path
		}

		tools = append(tools, ToolDef{
			Name:        name,
			Description: desc,
			Method:      op.Method,
			Path:        op.Path,
			Category:    category,
			Auth:        op.Auth,
			PathParams:  pathParams,
			HasBody:     hasBody,
			Annotations: ToolAnnotations{
				ReadOnlyHint:    op.Method == "GET",
				DestructiveHint: isDestructiveOperation(op.Method, op.Path),
			},
			InputSchema: inputSchema,
		})
	}

	return tools
}

func FilterToolsForClaims(tools []ToolDef, claims *models.UserClaims, cloudEnabled bool) []ToolDef {
	var out []ToolDef
	isAdmin := claims == nil || claims.Role == models.UserRoleAdmin || claims.Role == models.UserRoleOwner
	isReadOnly := claims != nil && claims.Role == "readonly"

	for _, t := range tools {
		if isReadOnly && !t.Annotations.ReadOnlyHint {
			continue
		}
		if t.Auth == apischema.AuthAdmin && !isAdmin {
			continue
		}
		if !cloudEnabled && strings.HasPrefix(t.Path, "/api/billing") {
			continue
		}
		out = append(out, t)
	}
	return out
}
