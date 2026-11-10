package volumes

import "errors"

var (
	ErrVolumeNotFound    = errors.New("service volume not found")
	ErrPathRequired      = errors.New("hostPath and containerPath are required")
	ErrServiceIDRequired = errors.New("serviceId is required")
	ErrUnauthorized      = errors.New("unauthorized volume access")
)
