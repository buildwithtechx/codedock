package auth

import "errors"

var (
	ErrInvalidPasscode       = errors.New("invalid passcode")
	ErrDeviceSessionNotFound = errors.New("device session not found")
	ErrDeviceSessionExpired  = errors.New("device session has expired")
	ErrDeviceAlreadyDecided  = errors.New("device session has already been decided")
)
