package domains

import "errors"

var (
	ErrDomainNotFound = errors.New("domain not found")
	ErrInvalidDomain  = errors.New("valid domain with serviceId and domainName required")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("insufficient permissions")
)
