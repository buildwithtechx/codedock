package templates

import "errors"

var (
	ErrTemplateNotFound     = errors.New("one-click template not found")
	ErrInvalidTemplate      = errors.New("invalid one-click template")
	ErrTemplateReviewFailed = errors.New("template review failed")
	ErrTemplateDeployFailed = errors.New("template deployment failed")
	ErrEnvironmentNotFound  = errors.New("target environment not found")
	ErrAppIDRequired        = errors.New("appId is required")
	ErrProjectIDRequired    = errors.New("projectId is required")
)
