package actions

import "errors"

var (
	ErrWorkflowNotFound    = errors.New("workflow not found")
	ErrRunNotFound         = errors.New("action run not found")
	ErrRunnerNotFound      = errors.New("action runner not found")
	ErrInvalidTrigger      = errors.New("invalid action trigger")
	ErrInvalidRunnerConfig = errors.New("invalid runner configuration")
	ErrServerIDRequired    = errors.New("serverId is required")
)
