package users

import "errors"

var (
	ErrUserEmailRequired  = errors.New("user email is required")
	ErrUserIDRequired     = errors.New("user id is required")
	ErrInvalidOldPassword = errors.New("invalid old password")
	ErrHashPasswordFailed = errors.New("failed to hash new password")
)
