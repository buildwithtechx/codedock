package vars

import "errors"

var (
	ErrVariableNotFound  = errors.New("environment variable not found")
	ErrKeyRequired       = errors.New("key is required")
	ErrServiceIDRequired = errors.New("serviceId is required")
	ErrInvalidRequest    = errors.New("invalid variable request")
)
