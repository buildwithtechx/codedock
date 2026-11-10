package system

import "errors"

var (
	ErrFoundingAdminNotFound = errors.New("founding admin user not found")
	ErrLocalCLIUnauthorized  = errors.New("local CLI access denied: must be run from localhost")
	ErrUpdaterNotInitialized = errors.New("updater service not initialized")
	ErrSystemRestartFailed   = errors.New("system restart failed")
	ErrSystemCleanupFailed   = errors.New("system cleanup failed")
)
