package models

type ActionWorkflow struct {
	ID          string `json:"id" db:"id"`
	ProjectID   string `json:"projectId,omitempty" db:"project_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
	Trigger     string `json:"trigger" db:"trigger"`
	YAML        string `json:"yaml,omitempty" db:"yaml"`
	Enabled     bool   `json:"enabled" db:"enabled"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
	UpdatedAt   string `json:"updatedAt" db:"updated_at"`
}

type ActionStep struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	DurationS int    `json:"durationSeconds,omitempty"`
}

type ActionJob struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Status    string       `json:"status"`
	Steps     []ActionStep `json:"steps"`
	RunnerID  string       `json:"runnerId,omitempty"`
	DurationS int          `json:"durationSeconds,omitempty"`
}

type ActionRun struct {
	ID          string      `json:"id" db:"id"`
	WorkflowID  string      `json:"workflowId" db:"workflow_id"`
	Name        string      `json:"name" db:"name"`
	Number      int         `json:"number" db:"number"`
	Status      string      `json:"status" db:"status"` // queued | running | ready | failed | cancelled
	TriggeredBy string      `json:"triggeredBy" db:"triggered_by"`
	RunnerID    string      `json:"runnerId,omitempty" db:"runner_id"`
	Logs        string      `json:"logs,omitempty" db:"logs"`
	Jobs        []ActionJob `json:"jobs,omitempty"`
	StartedAt   string      `json:"startedAt,omitempty" db:"started_at"`
	FinishedAt  string      `json:"finishedAt,omitempty" db:"finished_at"`
	DurationS   int         `json:"durationSeconds" db:"duration_seconds"`
	CreatedAt   string      `json:"createdAt" db:"created_at"`
	UpdatedAt   string      `json:"updatedAt" db:"updated_at"`
}

type ActionCapabilities struct {
	OS                 string   `json:"os"`           // linux | macos
	Architecture       string   `json:"architecture"` // x64 | arm64
	Docker             bool     `json:"docker"`
	DockerArchitecture string   `json:"dockerArchitecture,omitempty"`
	DockerPlatforms    []string `json:"dockerPlatforms,omitempty"`
	Git                bool     `json:"git"`
	Node               bool     `json:"node"`
	Distribution       *string  `json:"distribution,omitempty"`
	Version            *string  `json:"version,omitempty"`
}

type ActionRunnerConfig struct {
	Mode              string   `json:"mode"` // container | native
	Labels            []string `json:"labels"`
	Image             *string  `json:"image,omitempty"`
	MaxParallel       int      `json:"maxParallel"`
	CPU               float64  `json:"cpu"`
	MemoryMB          int      `json:"memoryMb"`
	AllowDockerSocket bool     `json:"allowDockerSocket"`
	CloudDiskGB       int      `json:"cloudDiskGb,omitempty"`
}

type ActionRunner struct {
	ID            string              `json:"id" db:"id"`
	Name          string              `json:"name" db:"name"`
	ServerID      *string             `json:"serverId,omitempty" db:"server_id"`
	CloudPoolID   *string             `json:"cloudPoolId,omitempty" db:"cloud_pool_id"`
	Kind          string              `json:"kind" db:"kind"` // server | cloud
	Mode          string              `json:"mode" db:"mode"` // container | native
	Image         string              `json:"image,omitempty" db:"image"`
	Status        string              `json:"status" db:"status"` // online | offline | busy
	Labels        []string            `json:"labels"`
	Config        ActionRunnerConfig  `json:"config"`
	Capabilities  *ActionCapabilities `json:"capabilities,omitempty"`
	Enabled       bool                `json:"enabled" db:"enabled"`
	Error         *string             `json:"error,omitempty" db:"error"`
	CheckedAt     string              `json:"checkedAt,omitempty"`
	LastHeartbeat string              `json:"lastHeartbeatAt,omitempty" db:"last_heartbeat"`
	CreatedAt     string              `json:"createdAt" db:"created_at"`
	UpdatedAt     string              `json:"updatedAt" db:"updated_at"`
}

type CreateWorkflowRequest struct {
	ProjectID   string `json:"projectId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Trigger     string `json:"trigger,omitempty"`
	YAML        string `json:"yaml,omitempty"`
}

type InspectDestinationRequest struct {
	ServerID string `json:"serverId"`
}

type CreateRunnerRequest struct {
	ServerID    string             `json:"serverId,omitempty"`
	CloudPoolID string             `json:"cloudPoolId,omitempty"`
	Kind        string             `json:"kind,omitempty"`
	Name        string             `json:"name"`
	Mode        string             `json:"mode,omitempty"`
	Image       string             `json:"image,omitempty"`
	Labels      []string           `json:"labels,omitempty"`
	Enabled     bool               `json:"enabled"`
	Config      ActionRunnerConfig `json:"config,omitempty"`
}
