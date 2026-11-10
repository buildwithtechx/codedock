package operations

import "errors"

var (
	ErrOperationNotFound  = errors.New("operation not found")
	ErrOperationExpired   = errors.New("operation has expired")
	ErrOperationCancelled = errors.New("operation was cancelled")
	ErrUnauthorized       = errors.New("unauthorized to manage operation")
)
