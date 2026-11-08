package apps

import "errors"

var (
	ErrAppServiceNotFound  = errors.New("app service not found")
	ErrInvalidAppService   = errors.New("valid app service with projectId and name required")
	ErrVariableNotFound    = errors.New("variable not found")
	ErrVolumeNotFound      = errors.New("volume not found")
	ErrAutoscalingNotFound = errors.New("autoscaling policy not found")
	ErrWebhookNotFound     = errors.New("webhook not found")
	ErrLogDrainNotFound    = errors.New("log drain not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("insufficient permissions")
)
