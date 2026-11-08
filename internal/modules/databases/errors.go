package databases

import "errors"

var (
	ErrDatabaseNotFound      = errors.New("database not found")
	ErrProjectNotFound       = errors.New("project not found")
	ErrAccessDenied          = errors.New("access denied")
	ErrDeployerNotConfigured = errors.New("database deployer not configured")
	ErrPrimaryKeyRequired    = errors.New("at least one primary key is required")
)
