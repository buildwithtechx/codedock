package deployments

type TriggerDeploymentRequest struct {
	Branch *string `json:"branch,omitempty"`
}

type RollbackRequest struct {
	TargetDeploymentID string `json:"targetDeploymentId"`
	Reason             string `json:"reason,omitempty"`
}

type CreateRegistryRequest struct {
	Name          string `json:"name"`
	RegistryURL   string `json:"registryUrl"`
	Username      string `json:"username"`
	PasswordToken string `json:"passwordToken"`
}

type UpdateRetentionRequest struct {
	MaxRetainedDeployments *int  `json:"maxRetainedDeployments,omitempty"`
	MaxRetainedDays        *int  `json:"maxRetainedDays,omitempty"`
	PruneArtifacts         *bool `json:"pruneArtifacts,omitempty"`
	PruneContainers        *bool `json:"pruneContainers,omitempty"`
}
