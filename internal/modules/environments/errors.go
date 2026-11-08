package environments

import "errors"

var (
	ErrEnvironmentNotFound = errors.New("environment not found")
	ErrDomainNotFound      = errors.New("domain not found")
	ErrRouteRuleNotFound   = errors.New("route rule not found")
	ErrProtectionNotFound  = errors.New("route protection not found")
	ErrInvalidEnvironment  = errors.New("valid environment with projectId and name required")
	ErrInvalidDomain       = errors.New("valid domain with serviceId and domainName required")
	ErrInvalidRouteRule    = errors.New("invalid route rule")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
)
