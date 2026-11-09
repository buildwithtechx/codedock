package ssh

import (
	"encoding/json"
	"time"
)

type ServerStatus string

const (
	ServerStatusOnline       ServerStatus = "online"
	ServerStatusOffline      ServerStatus = "offline"
	ServerStatusProvisioning ServerStatus = "provisioning"
)

type Server struct {
	ID                    string          `json:"id" db:"id"`
	OrganizationID        string          `json:"organizationId,omitempty" db:"organization_id"`
	UserID                string          `json:"userId" db:"user_id"`
	Name                  string          `json:"name" db:"name"`
	IPAddress             string          `json:"ipAddress" db:"ip_address"`
	IsLocal               bool            `json:"isLocal" db:"is_local"`
	SSHHost               string          `json:"sshHost" db:"ssh_host"`
	SSHPort               int             `json:"sshPort" db:"ssh_port"`
	SSHUser               string          `json:"sshUser" db:"ssh_user"`
	SSHAuthMethod         string          `json:"sshAuthMethod" db:"ssh_auth_method"`
	SSHKey                string          `json:"-" db:"ssh_key"`
	SSHPrivateKey         string          `json:"-" db:"ssh_private_key"`
	SSHPassword           string          `json:"-" db:"ssh_password"`
	SSHTransport          string          `json:"sshTransport" db:"ssh_transport"`
	SSHJumpHost           string          `json:"sshJumpHost,omitempty" db:"ssh_jump_host"`
	IsControlPlane        bool            `json:"isControlPlane" db:"-"`
	Status                ServerStatus    `json:"status" db:"status"`
	Provider              string          `json:"provider,omitempty" db:"provider"`
	ExternalID            string          `json:"externalId,omitempty" db:"external_id"`
	Region                string          `json:"region,omitempty" db:"region"`
	ServerType            string          `json:"serverType,omitempty" db:"server_type"`
	WorkerToken           string          `json:"workerToken,omitempty" db:"worker_token"`
	LastSeenAt            *time.Time      `json:"lastSeenAt,omitempty" db:"last_seen_at"`
	BuildCacheBudgetBytes *int64          `json:"buildCacheBudgetBytes,omitempty" db:"build_cache_budget_bytes"`
	Metrics               json.RawMessage `json:"metrics,omitempty" db:"metrics"`
	CreatedAt             time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time       `json:"updatedAt" db:"updated_at"`
}

type WorkerMetricsPayload struct {
	CPUUsagePercentage float64 `json:"cpu_usage_percentage"`
	MemoryUsageBytes   uint64  `json:"memory_usage_bytes"`
	MemoryLimitBytes   uint64  `json:"memory_limit_bytes"`
	DiskUsageBytes     uint64  `json:"disk_usage_bytes"`
	DiskTotalBytes     uint64  `json:"disk_total_bytes"`
}
