package deployments

import "time"

type Deployment struct {
	ID                 string           `json:"id" db:"id"`
	OrganizationID     string           `json:"organizationId,omitempty" db:"organization_id"`
	ProjectID          string           `json:"projectId" db:"project_id"`
	ServiceID          string           `json:"serviceId,omitempty" db:"service_id"`
	EnvironmentID      string           `json:"environmentId,omitempty" db:"environment_id"`
	Status             DeploymentStatus `json:"status" db:"status"`
	Branch             string           `json:"branch,omitempty" db:"branch"`
	CommitHash         string           `json:"commitHash,omitempty" db:"commit_hash"`
	CommitMessage      string           `json:"commitMessage,omitempty" db:"commit_message"`
	CommitSHABefore    string           `json:"commitShaBefore,omitempty" db:"commit_sha_before"`
	Trigger            string           `json:"trigger,omitempty" db:"trigger"`
	Framework          string           `json:"framework,omitempty" db:"framework"`
	ImageRef           string           `json:"imageRef,omitempty" db:"image_ref"`
	BuildDurationMS    int              `json:"buildDurationMs,omitempty" db:"build_duration_ms"`
	Version            int              `json:"version" db:"version"`
	BuildLogs          string           `json:"buildLogs,omitempty" db:"build_logs"`
	ContainerID        string           `json:"containerId,omitempty" db:"container_id"`
	ArtifactRetainedAt *time.Time       `json:"artifactRetainedAt,omitempty" db:"artifact_retained_at"`
	Pinned             bool             `json:"pinned" db:"pinned"`
	ResourceIssue      string           `json:"resourceIssue,omitempty" db:"resource_issue"`
	RestartLoop        string           `json:"restartLoop,omitempty" db:"restart_loop"`
	CreatedAt          time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time        `json:"updatedAt" db:"updated_at"`
	FinishedAt         *time.Time       `json:"finishedAt,omitempty" db:"finished_at"`
}

type DeploymentListItem struct {
	ID                 string           `json:"id" db:"id"`
	ServiceID          string           `json:"serviceId" db:"service_id"`
	ServiceName        string           `json:"serviceName" db:"service_name"`
	EnvironmentID      string           `json:"environmentId" db:"environment_id"`
	ProjectID          string           `json:"projectId" db:"project_id"`
	ProjectName        string           `json:"projectName" db:"project_name"`
	Status             DeploymentStatus `json:"status" db:"status"`
	Branch             string           `json:"branch,omitempty" db:"branch"`
	CommitHash         string           `json:"commitHash,omitempty" db:"commit_hash"`
	CommitMessage      string           `json:"commitMessage,omitempty" db:"commit_message"`
	Trigger            string           `json:"trigger,omitempty" db:"trigger"`
	ContainerID        string           `json:"containerId,omitempty" db:"container_id"`
	ArtifactRetainedAt *time.Time       `json:"artifactRetainedAt,omitempty" db:"artifact_retained_at"`
	Pinned             bool             `json:"pinned" db:"pinned"`
	ResourceIssue      string           `json:"resourceIssue,omitempty" db:"resource_issue"`
	RestartLoop        string           `json:"restartLoop,omitempty" db:"restart_loop"`
	CreatedAt          time.Time        `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time        `json:"updatedAt" db:"updated_at"`
	FinishedAt         *time.Time       `json:"finishedAt,omitempty" db:"finished_at"`
}

type DeploymentListFilter struct {
	OrganizationID string `json:"organizationId"`
	ProjectID      string `json:"projectId,omitempty"`
	ServiceID      string `json:"serviceId,omitempty"`
	Status         string `json:"status,omitempty"`
	Search         string `json:"search,omitempty"`
	Days           int    `json:"days,omitempty"`
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
}

type ServiceMetric struct {
	Timestamp  string  `json:"timestamp"`
	CPUPercent float64 `json:"cpuPercent"`
	MemoryMB   float64 `json:"memoryMB"`
	NetworkRx  float64 `json:"networkRxKB"`
	NetworkTx  float64 `json:"networkTxKB"`
}

type Registry struct {
	ID            string    `json:"id" db:"id"`
	ProjectID     string    `json:"projectId" db:"project_id"`
	Name          string    `json:"name" db:"name"`
	RegistryURL   string    `json:"registryUrl" db:"registry_url"`
	Username      string    `json:"username" db:"username"`
	PasswordToken string    `json:"passwordToken,omitempty" db:"password_token"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type RollbackPlan struct {
	TargetDeploymentID   string `json:"targetDeploymentId"`
	PreviousDeploymentID string `json:"previousDeploymentId"`
	ImageRef             string `json:"imageRef"`
	Feasible             bool   `json:"feasible"`
	Reason               string `json:"reason,omitempty"`
}

type RetentionPolicy struct {
	MaxRetainedDeployments int  `json:"maxRetainedDeployments"`
	MaxRetainedDays        int  `json:"maxRetainedDays"`
	PruneArtifacts         bool `json:"pruneArtifacts"`
	PruneContainers        bool `json:"pruneContainers"`
}

type RetentionReclaimResult struct {
	PrunedDeployments int   `json:"prunedDeployments"`
	PrunedContainers  int   `json:"prunedContainers"`
	PrunedImages      int   `json:"prunedImages"`
	ReclaimedBytes    int64 `json:"reclaimedBytes"`
}

type BuildAdmissionStatus struct {
	ActiveBuilds        int `json:"activeBuilds"`
	MaxConcurrentBuilds int `json:"maxConcurrentBuilds"`
	QueueLength         int `json:"queueLength"`
}

type DeployStatusOpts struct {
	ID            string
	Status        DeploymentStatus
	BuildLogs     string
	ContainerID   string
	ResourceIssue string
}

type AppServiceSummary struct {
	ID            string      `json:"id"`
	ProjectID     string      `json:"projectId"`
	EnvironmentID string      `json:"environmentId"`
	Name          string      `json:"name"`
	RepositoryURL string      `json:"repositoryUrl"`
	Branch        string      `json:"branch"`
	ImageRef      string      `json:"imageRef,omitempty"`
	ContainerID   string      `json:"containerId,omitempty"`
	Status        string      `json:"status"`
	RuntimeMode   RuntimeMode `json:"runtimeMode"`
	DeployToken   string      `json:"deployToken,omitempty"`
}
