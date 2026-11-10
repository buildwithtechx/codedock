package webhooks

import "errors"

var (
	ErrWebhookNotFound   = errors.New("webhook not found")
	ErrLogDrainNotFound  = errors.New("log drain not found")
	ErrServiceIDRequired = errors.New("serviceId is required")
	ErrURLRequired       = errors.New("URL is required")
)
