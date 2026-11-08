package security

import (
	"errors"
)

// Domain sentinel errors for the security module
var (
	ErrAdvisoryNotFound    = errors.New("security advisory not found")
	ErrProjectNotFound     = errors.New("project not found")
	ErrRemediationFailed   = errors.New("failed to apply security remediation")
	ErrInvalidSignature    = errors.New("invalid advisory feed signature")
	ErrInvalidProof        = errors.New("invalid security advisory proof")
	ErrSignatureInvalid    = errors.New("advisory signature verification failed")
	ErrPayloadTooLarge     = errors.New("advisory payload exceeds maximum size limit")
	ErrRollbackDetected    = errors.New("advisory manifest revision is older than stored minimum revision")
	ErrClockSkew           = errors.New("advisory manifest issued in the future")
	ErrManifestExpired     = errors.New("advisory manifest has expired")
	ErrInvalidFormat       = errors.New("unsupported advisory manifest format")
	ErrRuntimePatchInvalid = errors.New("invalid runtime patch specification")
)
