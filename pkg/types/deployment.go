package types

import "time"

type Deployment struct {
	ID            string     `json:"id" db:"id"`
	ServiceID     string     `json:"serviceId" db:"service_id"`
	EnvironmentID string     `json:"environmentId" db:"environment_id"`
	ProjectID     string     `json:"projectId" db:"project_id"`
	Status        string     `json:"status" db:"status"`
	Branch        string     `json:"branch,omitempty" db:"branch"`
	CommitHash    string     `json:"commitHash,omitempty" db:"commit_hash"`
	CommitMessage string     `json:"commitMessage,omitempty" db:"commit_message"`
	Trigger       string     `json:"trigger,omitempty" db:"trigger"`
	BuildLogs     string     `json:"buildLogs,omitempty" db:"build_logs"`
	ContainerID   string     `json:"containerId,omitempty" db:"container_id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
	FinishedAt    *time.Time `json:"finishedAt,omitempty" db:"finished_at"`
}

type Variable struct {
	ID            string    `json:"id" db:"id"`
	ServiceID     string    `json:"serviceId" db:"service_id"`
	ProjectID     string    `json:"projectId" db:"project_id"`
	EnvironmentID string    `json:"environmentId" db:"environment_id"`
	Key           string    `json:"key" db:"key"`
	Value         string    `json:"value" db:"value"`
	IsSecret      bool      `json:"isSecret" db:"is_secret"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type ServiceMetric struct {
	Timestamp  string  `json:"timestamp"`
	CPUPercent float64 `json:"cpuPercent"`
	MemoryMB   float64 `json:"memoryMB"`
	NetworkRx  float64 `json:"networkRxKB"`
	NetworkTx  float64 `json:"networkTxKB"`
}

type VarsRequest map[string]string

type SetEnvVarsRequest map[string]string
