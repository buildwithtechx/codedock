package canvas

import "errors"

var (
	ErrCanvasNotFound      = errors.New("environment canvas not found")
	ErrOrgIDRequired       = errors.New("organization is required")
	ErrCanvasIDRequired    = errors.New("canvas id is required")
	ErrEnvironmentRequired = errors.New("environment id is required")
)
