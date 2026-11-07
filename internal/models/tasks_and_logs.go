package models

import "time"

type ScheduledTask struct {
	ID         string    `json:"id" db:"id"`
	ServiceID  string    `json:"serviceId" db:"service_id"`
	Name       string    `json:"name" db:"name"`
	Schedule   string    `json:"schedule" db:"schedule"`
	Command    string    `json:"command" db:"command"`
	Status     string    `json:"status" db:"status"`
	LastRunAt  time.Time `json:"lastRunAt" db:"last_run_at"`
	LastOutput string    `json:"lastOutput" db:"last_output"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateScheduledTaskRequest struct {
	ServiceID string `json:"serviceId"`
	Name      string `json:"name"`
	Schedule  string `json:"schedule"`
	Command   string `json:"command"`
}

type UpdateScheduledTaskRequest struct {
	Name     *string `json:"name,omitempty"`
	Schedule *string `json:"schedule,omitempty"`
	Command  *string `json:"command,omitempty"`
	Status   *string `json:"status,omitempty"`
}

type PRPreview struct {
	ID            string          `json:"id" db:"id"`
	ServiceID     string          `json:"serviceId" db:"service_id"`
	ProjectID     string          `json:"projectId" db:"project_id"`
	PRNumber      int             `json:"prNumber" db:"pr_number"`
	Branch        string          `json:"branch" db:"branch"`
	CommitHash    string          `json:"commitHash" db:"commit_hash"`
	Status        PRPreviewStatus `json:"status" db:"status"`
	PreviewDomain string          `json:"previewDomain" db:"preview_domain"`
	ContainerID   string          `json:"containerId" db:"container_id"`
	CreatedAt     time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time       `json:"updatedAt" db:"updated_at"`
}

type LogDrainType string

const (
	LogDrainTypeAxiom    LogDrainType = "axiom"
	LogDrainTypeNewRelic LogDrainType = "new_relic"
	LogDrainTypeDatadog  LogDrainType = "datadog"
	LogDrainTypeWebhook  LogDrainType = "webhook"
)

type LogDrain struct {
	ID          string       `json:"id" db:"id"`
	ServiceID   string       `json:"serviceId" db:"service_id"`
	ProjectID   string       `json:"projectId" db:"project_id"`
	DrainType   LogDrainType `json:"drainType" db:"drain_type"`
	EndpointURL string       `json:"endpointUrl" db:"endpoint_url"`
	AuthToken   string       `json:"authToken" db:"auth_token"`
	CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time    `json:"updatedAt" db:"updated_at"`
}

type CreateLogDrainRequest struct {
	DrainType   LogDrainType `json:"drainType"`
	EndpointURL string       `json:"endpointUrl"`
	AuthToken   string       `json:"authToken"`
}
