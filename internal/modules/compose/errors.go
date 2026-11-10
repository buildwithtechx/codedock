package compose

import "errors"

var (
	ErrStackNotFound         = errors.New("compose stack not found")
	ErrInvalidComposeConfig  = errors.New("invalid compose configuration")
	ErrStackBusy             = errors.New("compose stack operation already in progress")
	ErrStackDeploymentFailed = errors.New("compose stack deployment failed")
	ErrProjectIDRequired     = errors.New("projectId is required")
	ErrStackIDRequired       = errors.New("stackId is required")
)
