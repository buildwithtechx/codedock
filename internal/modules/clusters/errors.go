package clusters

import "errors"

var (
	ErrClusterNotFound     = errors.New("cluster not found")
	ErrClusterUnauthorized = errors.New("unauthorized cluster access")
	ErrClusterDataNotFound = errors.New("cluster database not found")
	ErrStorageUnavailable  = errors.New("storage engine unavailable")
	ErrInvalidPlan         = errors.New("invalid cluster plan")
)
