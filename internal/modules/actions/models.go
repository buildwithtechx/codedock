package actions

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
	Status      string      `json:"status" db:"status"`
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
	OS                 string   `json:"os"`
	Architecture       string   `json:"architecture"`
	Docker             bool     `json:"docker"`
	DockerArchitecture string   `json:"dockerArchitecture,omitempty"`
	DockerPlatforms    []string `json:"dockerPlatforms,omitempty"`
	Git                bool     `json:"git"`
	Node               bool     `json:"node"`
	Distribution       *string  `json:"distribution,omitempty"`
	Version            *string  `json:"version,omitempty"`
}

type ActionRunnerConfig struct {
	Mode              string   `json:"mode"`
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
	Kind          string              `json:"kind" db:"kind"`
	Mode          string              `json:"mode" db:"mode"`
	Image         string              `json:"image,omitempty" db:"image"`
	Status        string              `json:"status" db:"status"`
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
