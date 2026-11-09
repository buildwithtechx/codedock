package monitoring

import (
	"time"

	"codedock/internal/platform/notifications"
)

type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

type NotificationDelivery = notifications.NotificationDelivery

type NotificationSubscription = notifications.NotificationSubscription
type NotificationDefault = notifications.NotificationDefault

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

type IngressConnection struct {
	ID                    string  `json:"id" db:"id"`
	ServerID              string  `json:"serverId" db:"server_id"`
	CredentialID          string  `json:"credentialId" db:"credential_id"`
	AccountID             string  `json:"accountId" db:"account_id"`
	TunnelID              *string `json:"tunnelId" db:"tunnel_id"`
	Name                  string  `json:"name" db:"name"`
	Connector             string  `json:"connector" db:"connector"`
	EdgePort              int     `json:"edgePort" db:"edge_port"`
	ManagedTunnel         bool    `json:"managedTunnel" db:"managed_tunnel"`
	DefaultForDomains     bool    `json:"defaultForDomains" db:"default_for_domains"`
	DefaultAccessPolicyID *string `json:"defaultAccessPolicyId" db:"default_access_policy_id"`
	Status                string  `json:"status" db:"status"`
	Error                 *string `json:"error" db:"error"`
	CreatedAt             string  `json:"createdAt" db:"created_at"`
	UpdatedAt             string  `json:"updatedAt" db:"updated_at"`
}

type IngressHostname struct {
	ID             string  `json:"id" db:"id"`
	DomainID       *string `json:"domainId" db:"domain_id"`
	ConnectionID   string  `json:"connectionId" db:"connection_id"`
	Hostname       string  `json:"hostname" db:"hostname"`
	AccessPolicyID *string `json:"accessPolicyId" db:"access_policy_id"`
	Status         string  `json:"status" db:"status"`
	Phase          string  `json:"phase" db:"phase"`
	Error          *string `json:"error" db:"error"`
	UpdatedAt      string  `json:"updatedAt" db:"updated_at"`
}

type NotificationSettings = notifications.NotificationSettings

type NotificationEvent struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	Level     string `json:"level"`
	EventType string `json:"eventType"`
	ProjectID string `json:"projectId,omitempty"`
	URL       string `json:"url,omitempty"`
}

type CloudSupportTicketRecord struct {
	ID        string
	AccountID string
	Subject   string
	Category  CloudSupportCategory
	Status    CloudSupportStatus
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CloudSupportMessageRecord struct {
	ID        string
	TicketID  string
	Author    CloudSupportAuthor
	Body      string
	CreatedAt time.Time
}

type CloudSupportCustomerTicket struct {
	ID        string               `json:"id"`
	Subject   string               `json:"subject"`
	Category  CloudSupportCategory `json:"category"`
	Status    CloudSupportStatus   `json:"status"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
}

type Operation struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"-" db:"user_id"`
	ProjectID string `json:"projectId" db:"project_id"`
	Kind      string `json:"kind" db:"kind"`
	Target    string `json:"target" db:"target"`
	Status    string `json:"status" db:"status"`
	Phase     string `json:"phase" db:"phase"`
	Effects   string `json:"effects" db:"effects"`
	Payload   string `json:"-" db:"encrypted_payload"`
	Snapshot  string `json:"-" db:"snapshot"`
	TokenHash string `json:"-" db:"token_hash"`
	Error     string `json:"error" db:"error"`
	Logs      string `json:"logs" db:"logs"`
	ExpiresAt int64  `json:"expiresAt" db:"expires_at"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type OperationReview struct {
	Operation    *Operation `json:"operation"`
	Confirmation string     `json:"confirmation"`
}

type SystemStats struct {
	CPU       CPUStats    `json:"cpu"`
	Memory    MemoryStats `json:"memory"`
	Disk      DiskStats   `json:"disk"`
	Uptime    int64       `json:"uptimeSeconds"`
	LoadAvg   [3]float64  `json:"loadAvg"`
	Processes int         `json:"processes"`
	Docker    DockerStats `json:"docker"`
}

type DockerStats struct {
	Images        DockerLayerStat `json:"images"`
	Containers    DockerLayerStat `json:"containers"`
	Volumes       DockerLayerStat `json:"volumes"`
	BuildCache    DockerLayerStat `json:"buildCache"`
	ReclaimableGB float64         `json:"reclaimableGb"`
}

type DockerLayerStat struct {
	Active      string `json:"active"`
	TotalCount  string `json:"totalCount"`
	Size        string `json:"size"`
	Reclaimable string `json:"reclaimable"`
}

type CPUStats struct {
	Percent float64 `json:"percent"`
	Cores   int     `json:"cores"`
}

type MemoryStats struct {
	TotalMB int64   `json:"totalMb"`
	UsedMB  int64   `json:"usedMb"`
	FreeMB  int64   `json:"freeMb"`
	Percent float64 `json:"percent"`
}

type DiskStats struct {
	TotalGB int64   `json:"totalGb"`
	UsedGB  int64   `json:"usedGb"`
	FreeGB  int64   `json:"freeGb"`
	Percent float64 `json:"percent"`
}

type UpdateStatus struct {
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	UpdateAvail    bool   `json:"updateAvailable"`
	Updating       bool   `json:"updating"`
}
