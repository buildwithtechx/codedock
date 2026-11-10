package mcp

import "errors"

var (
	ErrToolNotFound       = errors.New("mcp tool not found")
	ErrExecutionFailed    = errors.New("mcp tool execution failed")
	ErrInvalidToolPayload = errors.New("invalid mcp tool payload")
	ErrUnauthorized       = errors.New("unauthorized mcp access")
)
