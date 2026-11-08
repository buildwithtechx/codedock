package monitoring

import "errors"

var (
	ErrUnauthorized          = errors.New("unauthorized")
	ErrNotFound              = errors.New("resource not found")
	ErrForbidden             = errors.New("forbidden")
	ErrObservabilityDisabled = errors.New("observability is disabled: set CODEDOCK_OBSERVABILITY=true for historical logs and metrics")
	ErrUpdaterNotInitialized = errors.New("updater service not initialized")
	ErrInvalidRequest        = errors.New("invalid request payload")
	ErrStreamingUnsupported  = errors.New("streaming unsupported")
)
