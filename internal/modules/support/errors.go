package support

import "errors"

var (
	ErrTicketNotFound  = errors.New("support ticket not found")
	ErrTicketClosed    = errors.New("support ticket is closed")
	ErrSubjectRequired = errors.New("subject is required")
	ErrMessageRequired = errors.New("message is required")
	ErrInvalidPriority = errors.New("invalid support priority")
)
