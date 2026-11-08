package ai

import "errors"

var (
	ErrFailedLoadSettings   = errors.New("failed to load AI settings")
	ErrPromptRequired       = errors.New("prompt is required")
	ErrInvalidRequest       = errors.New("invalid request payload")
	ErrDeploymentNotFound   = errors.New("deployment not found")
	ErrDeploymentNotFailed  = errors.New("AI explanations are only available after a deployment fails")
	ErrServiceNotFound      = errors.New("service not found")
	ErrProviderNotSupported = errors.New("unsupported AI provider")
)
