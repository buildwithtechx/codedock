package servers

import "errors"

var (
	ErrServerNotFound       = errors.New("server not found")
	ErrServerUnauthorized   = errors.New("unauthorized server access")
	ErrInvalidServerInput   = errors.New("invalid server input")
	ErrSSHConnectionFailed  = errors.New("failed to connect via SSH")
	ErrServerNameRequired   = errors.New("server name is required")
	ErrControlPlaneReadOnly = errors.New("control plane server cannot be modified")
	ErrManagedQuotaExceeded = errors.New("managed server quota exceeded")
	ErrCredentialInUse      = errors.New("managed credential is still in use")
	ErrTakeoverFailed       = errors.New("takeover scan failed")
)
