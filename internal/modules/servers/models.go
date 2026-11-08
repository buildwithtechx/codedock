package servers

import (
	"encoding/json"
	"time"
)

type Server struct {
	ID                    string                 `json:"id" db:"id"`
	OrganizationID        string                 `json:"organizationId,omitempty" db:"organization_id"`
	UserID                string                 `json:"userId" db:"user_id"`
	Name                  string                 `json:"name" db:"name"`
	IPAddress             string                 `json:"ipAddress" db:"ip_address"`
	IsLocal               bool                   `json:"isLocal" db:"is_local"`
	SSHHost               string                 `json:"sshHost" db:"ssh_host"`
	SSHPort               int                    `json:"sshPort" db:"ssh_port"`
	SSHUser               string                 `json:"sshUser" db:"ssh_user"`
	SSHAuthMethod         string                 `json:"sshAuthMethod" db:"ssh_auth_method"`
	SSHKey                string                 `json:"-" db:"ssh_key"`
	SSHPrivateKey         string                 `json:"-" db:"ssh_private_key"`
	SSHPassword           string                 `json:"-" db:"ssh_password"`
	SSHTransport          string                 `json:"sshTransport" db:"ssh_transport"`
	SSHJumpHost           string                 `json:"sshJumpHost,omitempty" db:"ssh_jump_host"`
	IsControlPlane        bool                   `json:"isControlPlane" db:"-"`
	Status                ServerStatus           `json:"status" db:"status"`
	Provider              string                 `json:"provider,omitempty" db:"provider"`
	ExternalID            string                 `json:"externalId,omitempty" db:"external_id"`
	Region                string                 `json:"region,omitempty" db:"region"`
	ServerType            string                 `json:"serverType,omitempty" db:"server_type"`
	WorkerToken           string                 `json:"workerToken,omitempty" db:"worker_token"`
	LastSeenAt            *time.Time             `json:"lastSeenAt,omitempty" db:"last_seen_at"`
	BuildCacheBudgetBytes *int64                 `json:"buildCacheBudgetBytes,omitempty" db:"build_cache_budget_bytes"`
	Managed               *CloudWorkspaceSummary `json:"managed,omitempty" db:"-"`
	Capabilities          *ServerCapabilities    `json:"capabilities,omitempty" db:"-"`
	Metrics               json.RawMessage        `json:"metrics,omitempty" db:"metrics"`
	CreatedAt             time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time              `json:"updatedAt" db:"updated_at"`
}

type ServerCapabilities struct {
	Monitor           bool `json:"monitor"`
	Terminal          bool `json:"terminal"`
	Exec              bool `json:"exec"`
	HostConfiguration bool `json:"hostConfiguration"`
	SSH               bool `json:"ssh"`
	NetworkSettings   bool `json:"networkSettings"`
	BuildCacheBudget  bool `json:"buildCacheBudget"`
}

type ServerComponentStatus struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Installable bool   `json:"installable"`
	Installed   bool   `json:"installed"`
	Version     string `json:"version,omitempty"`
	Healthy     bool   `json:"healthy"`
	Message     string `json:"message"`
}

type ServerListener struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	State    string `json:"state"`
	Exposed  bool   `json:"exposed"`
	Process  string `json:"process"`
}

type ServerPortScanResult struct {
	Scanned   bool             `json:"scanned"`
	ServerID  string           `json:"serverId"`
	Listeners []ServerListener `json:"listeners"`
}

type ServerRateLimitConfig struct {
	RPS       int      `json:"rps"`
	Burst     int      `json:"burst"`
	Whitelist []string `json:"whitelist"`
}

type ServerNetworkSettings struct {
	InternetAccess *bool    `json:"internetAccess"`
	IngressPorts   []int    `json:"ingressPorts"`
	IngressAll     *bool    `json:"ingressAll,omitempty"`
	Egress         []string `json:"egress,omitempty"`
	PrivateIP      string   `json:"privateIp,omitempty"`
	OutboundIP     string   `json:"outboundIp,omitempty"`
	OutboundMode   string   `json:"outboundMode,omitempty"`
	Revision       string   `json:"revision,omitempty"`
}

type ServerNetwork struct {
	ID               string                   `json:"id"`
	Name             string                   `json:"name"`
	Mode             string                   `json:"mode"`
	SubnetCIDR       string                   `json:"subnetCidr"`
	MTU              int                      `json:"mtu"`
	VerificationPort int                      `json:"verificationPort"`
	Status           string                   `json:"status"`
	Revision         int                      `json:"revision"`
	Members          any                      `json:"members,omitempty"`
	AccessPolicy     any                      `json:"accessPolicy,omitempty"`
	Watchdog         *NetworkRollbackWatchdog `json:"watchdog,omitempty"`
	CreatedAt        string                   `json:"createdAt"`
	UpdatedAt        string                   `json:"updatedAt"`
}

type NetworkRollbackWatchdog struct {
	OperationID      string `json:"operationId"`
	Armed            bool   `json:"armed"`
	TimeoutSeconds   int    `json:"timeoutSeconds"`
	RemainingSeconds int    `json:"remainingSeconds"`
	ExpiresAt        string `json:"expiresAt"`
	CanCommit        bool   `json:"canCommit"`
	CanRollback      bool   `json:"canRollback"`
	Status           string `json:"status"`
}
