package chatapps

import "errors"

var (
	ErrConnectionNotFound = errors.New("chat connection not found")
	ErrConnectionDisabled = errors.New("chat connection is disabled")
	ErrChallengeNotFound  = errors.New("pairing challenge not found")
	ErrChallengeExpired   = errors.New("pairing challenge has expired")
	ErrInvalidPairingCode = errors.New("invalid pairing verification code")
	ErrIdentityNotFound   = errors.New("paired chat identity not found")
	ErrTransportFailed    = errors.New("chat transport communication failed")
	ErrUnauthorizedChat   = errors.New("unauthorized chat command execution")
)
