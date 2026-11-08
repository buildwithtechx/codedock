package projects

import "errors"

var (
	ErrProjectIDRequired       = errors.New("project id is required")
	ErrProjectNotFound         = errors.New("project not found")
	ErrOrganizationIDRequired  = errors.New("organization id is required")
	ErrProjectAppNil           = errors.New("project app is nil")
	ErrProjectAppNameRequired  = errors.New("name is required")
	ErrProjectAppIDRequired    = errors.New("project app id is required")
	ErrTokenInvalid            = errors.New("valid token with projectId and name required")
	ErrEnvironmentNameRequired = errors.New("project id and environment name required")
	ErrEnvironmentIDRequired   = errors.New("environment id is required")
	ErrServiceIDRequired       = errors.New("service id is required")
	ErrVariableIDRequired      = errors.New("variable id is required")
	ErrOrgRepoNotInitialized   = errors.New("organization repository is not initialized")
)
