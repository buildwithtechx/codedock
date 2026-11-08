package deployments

import "errors"

var (
	ErrDeploymentNotFound       = errors.New("deployment not found")
	ErrDeploymentCancelled      = errors.New("deployment cancelled")
	ErrConcurrentDeployment     = errors.New("another deployment is already running for this service")
	ErrAdmissionLimitExceeded   = errors.New("build admission limit exceeded")
	ErrRegistryNotFound         = errors.New("container registry not found")
	ErrInvalidArchive           = errors.New("invalid deployment archive")
	ErrArchiveTooLarge          = errors.New("archive exceeds maximum upload limit of 100MB")
	ErrTerminalUnauthorized     = errors.New("insufficient permissions to access this terminal")
	ErrUserDeactivated          = errors.New("user account not found or deactivated")
	ErrRollbackNotFound         = errors.New("no suitable previous deployment found for rollback")
	ErrDeploymentAlreadyRunning = errors.New("deployment is currently running")
	ErrDeploymentNotRunning     = errors.New("deployment is no longer running")
	ErrAppServiceRequired       = errors.New("app service is required")
	ErrProjectRequired          = errors.New("project ID is required")
)
