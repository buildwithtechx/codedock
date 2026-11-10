package tasks

import "errors"

var (
	ErrTaskNotFound    = errors.New("scheduled task not found")
	ErrInvalidCronSpec = errors.New("invalid cron schedule specification")
	ErrCommandRequired = errors.New("command is required")
	ErrExecutionFailed = errors.New("scheduled task execution failed")
)
