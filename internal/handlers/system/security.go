package system

import (
	"codedock/internal/modules/security"
)

// SecurityHandler re-exports the modular security handler for backwards compatibility
type SecurityHandler = security.Handler

// NewSecurityHandler creates a new security handler
func NewSecurityHandler(service security.SecurityService) *SecurityHandler {
	return security.NewHandler(service)
}

// SetPolicyRequest re-exports the request DTO
type SetPolicyRequest = security.SetPolicyRequest

// ApplyRemediationRequest re-exports the request DTO
type ApplyRemediationRequest = security.ApplyRemediationRequest
