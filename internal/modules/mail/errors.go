package mail

import "errors"

var (
	ErrMailNotFound   = errors.New("mail resource not found")
	ErrMailExists     = errors.New("mail resource already exists")
	ErrServerRequired = errors.New("organization, domain, maild url, and api token are required")
	ErrServerNotFound = errors.New("mail server not found")
	ErrDomainRequired = errors.New("domainId is required")
	ErrAliasRequired  = errors.New("domainId and aliasId are required")
	ErrInvalidPayload = errors.New("invalid request payload")
	ErrUnauthorized   = errors.New("unauthorized mail request")
	ErrForbidden      = errors.New("insufficient permissions for this organization")
	ErrMissingOrg     = errors.New("missing organization id")
)
