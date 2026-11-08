package dns

import "errors"

var (
	ErrRecordNotFound = errors.New("dns record not found")
	ErrDomainRequired = errors.New("domainName is required")
	ErrFieldsRequired = errors.New("domain, type, name, and value are required")
	ErrInvalidPayload = errors.New("invalid request payload")
	ErrMissingID      = errors.New("missing record id")
)
