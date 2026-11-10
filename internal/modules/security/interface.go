package security

import "context"

// SecurityService coordinates vulnerability checks, policy, and automated patch remediation.
type SecurityService interface {
	GetAdvisories(ctx context.Context) (*AdvisoriesResponse, error)
	CheckProject(ctx context.Context, projectID string, force bool) (*ProjectSecurityResult, error)
	GetProjectSecurity(ctx context.Context, projectID string) (*ProjectSecurityResult, error)
	SetProjectPolicy(ctx context.Context, projectID string, autoRefresh bool) error
	ApplyRemediation(ctx context.Context, projectID, advisoryID, deploymentID string, serviceID, runtimePatchID *string) (*RefreshAttempt, error)
}

// DeploymentTrigger initiates rolling deployments for patched services.
type DeploymentTrigger interface {
	TriggerRemediation(ctx context.Context, serviceID string) error
}

// ProjectProvider resolves project metadata for security auditing.
type ProjectProvider interface {
	GetProject(ctx context.Context, id string) (*SecurityProject, error)
}

// AppServiceProvider retrieves and mutates service image tags and runtime definitions.
type AppServiceProvider interface {
	ListAppServices(ctx context.Context, projectID string) ([]*SecurityAppService, error)
	UpdateAppService(ctx context.Context, svc *SecurityAppService) error
}

// DeploymentProvider queries deployment manifests for container image digest inspection.
type DeploymentProvider interface {
	GetDeployment(ctx context.Context, id string) (*SecurityDeployment, error)
}
