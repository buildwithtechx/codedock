package security

// SetPolicyRequest sets automated advisory refresh policies for a project
type SetPolicyRequest struct {
	AutoRefresh bool `json:"autoRefresh"`
}

// ApplyRemediationRequest requests automated image refresh or runtime patching for an advisory
type ApplyRemediationRequest struct {
	AdvisoryID     string  `json:"advisoryId" validate:"required"`
	DeploymentID   string  `json:"deploymentId" validate:"required"`
	ServiceID      *string `json:"serviceId,omitempty"`
	RuntimePatchID *string `json:"runtimePatchId,omitempty"`
}

// CheckProjectRequest triggers an on-demand security scan for a project
type CheckProjectRequest struct {
	ProjectID string `json:"projectId"`
	Force     bool   `json:"force"`
}
