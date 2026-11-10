package takeover

import "errors"

var (
	ErrRunNotFound     = errors.New("takeover run not found")
	ErrHostRequired    = errors.New("host is required")
	ErrSSHUserRequired = errors.New("sshUser is required")
	ErrScanFailed      = errors.New("server takeover scan failed")
	ErrAdoptionFailed  = errors.New("server takeover adoption failed")
)
