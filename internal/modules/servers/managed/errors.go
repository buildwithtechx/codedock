package managed

import "errors"

var (
	ErrCredentialNotFound = errors.New("managed provider credential not found")
	ErrServerNotFound     = errors.New("managed server not found")
	ErrInvalidOperation   = errors.New("invalid managed operation")
	ErrQuotaExceeded      = errors.New("managed server quota exceeded")
	ErrSSHConfigFailed    = errors.New("failed to configure ssh access")
)
