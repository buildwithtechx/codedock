package migrations

import "errors"

var (
	ErrSourceNotFound = errors.New("migration source not found")
	ErrRunNotFound    = errors.New("migration run not found")
	ErrInvalidPayload = errors.New("invalid request payload")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrPassphraseReq  = errors.New("passphrase is required")
	ErrBundleFileReq  = errors.New("bundle file is required")
	ErrBundleTooLarge = errors.New("migration bundle file exceeds maximum allowed size")
	ErrSetupForbidden = errors.New("setup import is only available before the first account is created")
)
