package apischema

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func Document() map[string]any {
	paths := map[string]any{}
	for _, op := range All() {
		path := echoPathToOpenAPI(op.Path)
		item, ok := paths[path].(map[string]any)
		if !ok {
			item = map[string]any{}
			paths[path] = item
		}
		item[strings.ToLower(op.Method)] = operationDocument(op)
	}
	return map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":       "Codedock API",
			"version":     "1.0.0",
			"description": "Maintained request/response schemas for the Codedock control plane. Every route is covered; the repository test suite fails when a route lacks an entry here.",
		},
		"servers": []any{map[string]any{"url": "/api"}},
		"paths":   paths,
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{"type": "http", "scheme": "bearer", "bearerFormat": "JWT"},
			},
		},
	}
}

func Marshal() ([]byte, error) {
	rendered, err := json.MarshalIndent(Document(), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal openapi document: %w", err)
	}
	return append(rendered, '\n'), nil
}

func echoPathToOpenAPI(path string) string {
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "{" + strings.TrimPrefix(part, ":") + "}"
		}
	}
	return strings.Join(parts, "/")
}

func operationDocument(op Operation) map[string]any {
	document := map[string]any{
		"summary":     op.Summary,
		"tags":        op.Tags,
		"operationId": operationID(op),
		"responses":   responsesDocument(op),
	}
	parameters := []any{}
	for _, name := range echoPathParams(op.Path) {
		parameters = append(parameters, map[string]any{
			"name":        name,
			"in":          "path",
			"required":    true,
			"description": pathParamDescription(op, name),
			"schema":      map[string]any{"type": "string"},
		})
	}
	for _, query := range op.Query {
		parameters = append(parameters, map[string]any{
			"name":        query.Name,
			"in":          "query",
			"required":    query.Required,
			"description": query.Description,
			"schema":      map[string]any{"type": queryParamType(query)},
		})
	}
	if len(parameters) > 0 {
		document["parameters"] = parameters
	}
	if op.Authenticated() {
		document["security"] = []any{map[string]any{"bearerAuth": []any{}}}
		document["x-auth"] = op.Auth
	}
	if op.Request != nil {
		document["requestBody"] = map[string]any{
			"required": true,
			"content": map[string]any{
				op.Request.Content: map[string]any{"schema": schemaDocument(op.Request.Schema)},
			},
		}
	}
	return document
}

func operationID(op Operation) string {
	cleaned := strings.NewReplacer("/api", "", ":", "", "/", " ", "-", " ", "_", " ").Replace(op.Path)
	words := strings.Fields(strings.ToLower(op.Method) + " " + cleaned)
	for i := 1; i < len(words); i++ {
		words[i] = strings.ToUpper(words[i][:1]) + words[i][1:]
	}
	seen := map[string]bool{}
	out := []string{}
	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			out = append(out, word)
		}
	}
	return strings.Join(out, "")
}

func responsesDocument(op Operation) map[string]any {
	responses := map[string]any{}
	if op.Raw && op.SuccessCode() == 204 {
		responses["204"] = map[string]any{"description": op.Summary}
	} else if op.Raw {
		responses[fmt.Sprint(op.SuccessCode())] = map[string]any{
			"description": op.Summary,
			"content": map[string]any{
				op.ResponseContent(): map[string]any{"schema": schemaDocument(op.Response)},
			},
		}
	} else {
		responses[fmt.Sprint(op.SuccessCode())] = map[string]any{
			"description": op.Summary,
			"content": map[string]any{
				"application/json": map[string]any{"schema": schemaDocument(envelopeSchema(op.Response))},
			},
		}
	}
	for _, code := range errorCodes(op) {
		responses[fmt.Sprint(code)] = map[string]any{
			"description": errorDescription(code),
			"content": map[string]any{
				"application/json": map[string]any{"schema": schemaDocument(errorSchema())},
			},
		}
	}
	return responses
}

func errorCodes(op Operation) []int {
	codes := map[int]bool{400: true, 500: true}
	if op.Authenticated() {
		codes[401] = true
		codes[403] = true
	}
	if strings.Contains(op.Path, ":") {
		codes[404] = true
	}
	for _, code := range op.Errors {
		codes[code] = true
	}
	out := []int{}
	for code := range codes {
		if code != op.SuccessCode() {
			out = append(out, code)
		}
	}
	sort.Ints(out)
	return out
}

func errorDescription(code int) string {
	switch code {
	case 400:
		return "Invalid request"
	case 401:
		return "Missing or invalid credentials"
	case 403:
		return "Insufficient permissions"
	case 404:
		return "Resource not found"
	case 409:
		return "Conflicting state"
	case 422:
		return "Unprocessable request"
	default:
		return "Server error"
	}
}

func envelopeSchema(data Schema) Schema {
	return Obj("Codedock API envelope", RF("status", StrEnum("Result status", "success")), RF("message", Str("Human-readable message")), F("data", data), F("path", Str("Request path")), F("executionTime", Num("Handler seconds")))
}

func errorSchema() Schema {
	return Obj("Codedock API error", RF("status", StrEnum("Result status", "error")), RF("message", Str("Human-readable message")), F("path", Str("Request path")))
}

func PaginatedSchema(records Schema) Schema {
	return Obj("Paginated records", RF("records", records), RF("total", Int("Total records")), RF("page", Int("Current page")), RF("totalPages", Int("Total pages")))
}

func schemaDocument(schema Schema) map[string]any {
	document := map[string]any{}
	if schema.Description != "" {
		document["description"] = schema.Description
	}
	if schema.Type != "" {
		document["type"] = schema.Type
	}
	if schema.Format != "" {
		document["format"] = schema.Format
	}
	if len(schema.Enum) > 0 {
		document["enum"] = schema.Enum
	}
	if len(schema.Properties) > 0 {
		properties := map[string]any{}
		for _, prop := range schema.Properties {
			properties[prop.Name] = schemaDocument(prop.Schema)
		}
		document["properties"] = properties
	}
	if len(schema.Required) > 0 {
		document["required"] = schema.Required
	}
	if schema.Items != nil {
		document["items"] = schemaDocument(*schema.Items)
	}
	if schema.AdditionalProperties != nil {
		document["additionalProperties"] = schemaDocument(*schema.AdditionalProperties)
	}
	return document
}

func echoPathParams(path string) []string {
	params := []string{}
	for _, part := range strings.Split(path, "/") {
		if strings.HasPrefix(part, ":") {
			params = append(params, strings.TrimPrefix(part, ":"))
		}
	}
	return params
}

func pathParamDescription(op Operation, name string) string {
	for _, param := range op.PathParams {
		if param.Name == name {
			return param.Description
		}
	}
	return name
}

func queryParamType(query Param) string {
	if query.Type != "" {
		return query.Type
	}
	return "string"
}
