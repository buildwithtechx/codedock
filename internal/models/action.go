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

type ActionRunner struct {
	ID            string   `json:"id" db:"id"`
	Name          string   `json:"name" db:"name"`
	Mode          string   `json:"mode" db:"mode"` // container | native
	Image         string   `json:"image,omitempty" db:"image"`
	Status        string   `json:"status" db:"status"` // online | offline | busy
	Labels        []string `json:"labels"`
	Enabled       bool     `json:"enabled" db:"enabled"`
	LastHeartbeat string   `json:"lastHeartbeatAt,omitempty" db:"last_heartbeat"`
	CreatedAt     string   `json:"createdAt" db:"created_at"`
	UpdatedAt     string   `json:"updatedAt" db:"updated_at"`
}

type CreateWorkflowRequest struct {
	ProjectID   string `json:"projectId,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Trigger     string `json:"trigger,omitempty"`
	YAML        string `json:"yaml,omitempty"`
}

type CreateRunnerRequest struct {
	Name    string   `json:"name"`
	Mode    string   `json:"mode"`
	Image   string   `json:"image,omitempty"`
	Labels  []string `json:"labels,omitempty"`
	Enabled bool     `json:"enabled"`
}
