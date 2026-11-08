package monitoring

import "time"

type Vault interface {
	Encrypt(plaintext []byte) (string, error)
	Decrypt(ciphertext string) ([]byte, error)
}

type NotificationDelivery struct {
	ID        string `json:"id" db:"id"`
	EventType string `json:"eventType" db:"event_type"`
	Channel   string `json:"channel" db:"channel"`
	Title     string `json:"title" db:"title"`
	Success   bool   `json:"success" db:"success"`
	Error     string `json:"error" db:"error"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

type NotificationSubscription struct {
	ID             string   `json:"id"`
	UserID         string   `json:"userId"`
	OrganizationID string   `json:"organizationId"`
	Category       string   `json:"category"`
	Channels       []string `json:"channels"`
	Enabled        bool     `json:"enabled"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type NotificationDefault struct {
	OrganizationID string   `json:"organizationId"`
	Category       string   `json:"category"`
	Channels       []string `json:"channels"`
	Enabled        bool     `json:"enabled"`
	UpdatedAt      string   `json:"updatedAt"`
}

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

type NotificationSettings struct {
	ID                    string `json:"id"`
	DiscordWebhookURL     string `json:"discordWebhookUrl,omitempty"`
	DiscordPingEnabled    bool   `json:"discordPingEnabled"`
	DiscordEnabled        bool   `json:"discordEnabled"`
	SlackWebhookURL       string `json:"slackWebhookUrl,omitempty"`
	SlackEnabled          bool   `json:"slackEnabled"`
	TelegramBotToken      string `json:"telegramBotToken,omitempty"`
	TelegramChatID        string `json:"telegramChatId,omitempty"`
	TelegramEnabled       bool   `json:"telegramEnabled"`
	SMTPHost              string `json:"smtpHost,omitempty"`
	SMTPPort              int    `json:"smtpPort,omitempty"`
	SMTPUser              string `json:"smtpUser,omitempty"`
	SMTPPassword          string `json:"smtpPassword,omitempty"`
	SMTPFromName          string `json:"smtpFromName,omitempty"`
	SMTPFromAddress       string `json:"smtpFromAddress,omitempty"`
	SMTPEnabled           bool   `json:"smtpEnabled"`
	ResendAPIKey          string `json:"resendApiKey,omitempty"`
	ResendEnabled         bool   `json:"resendEnabled"`
	PushoverUserKey       string `json:"pushoverUserKey,omitempty"`
	PushoverAPIToken      string `json:"pushoverApiToken,omitempty"`
	PushoverEnabled       bool   `json:"pushoverEnabled"`
	GenericWebhookURL     string `json:"genericWebhookUrl,omitempty"`
	GenericWebhookEnabled bool   `json:"genericWebhookEnabled"`
	ExpoPushToken         string `json:"expoPushToken,omitempty"`
	ExpoEnabled           bool   `json:"expoEnabled"`
	NotificationAlerts    bool   `json:"notificationAlerts"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

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

type SystemStats struct {
	CPUUsage    float64 `json:"cpuUsage"`
	MemoryUsage float64 `json:"memoryUsage"`
	DiskUsage   float64 `json:"diskUsage"`
	Uptime      int64   `json:"uptime"`
}

type UpdateStatus struct {
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
	UpdateAvail    bool   `json:"updateAvailable"`
	Updating       bool   `json:"updating"`
}
