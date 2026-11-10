package networks

import "errors"

var (
	ErrNetworkNotFound  = errors.New("server network not found")
	ErrNetworkNameEmpty = errors.New("network name is required")
	ErrWatchdogNotFound = errors.New("network rollback watchdog not found")
)
