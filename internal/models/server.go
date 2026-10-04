package models

import "time"

type ServerStatus string

const (
	ServerStatusOnline       ServerStatus = "online"
	ServerStatusOffline      ServerStatus = "offline"
	ServerStatusProvisioning ServerStatus = "provisioning"
)

type Server struct {
	ID             string       `json:"id" db:"id"`
	OrganizationID string       `json:"organizationId,omitempty" db:"organization_id"`
	UserID         string       `json:"userId" db:"user_id"`
	Name           string       `json:"name" db:"name"`
	IPAddress      string       `json:"ipAddress" db:"ip_address"`
	IsLocal        bool         `json:"isLocal" db:"is_local"`
	SSHHost        string       `json:"sshHost" db:"ssh_host"`
	SSHPort        int          `json:"sshPort" db:"ssh_port"`
	SSHUser        string       `json:"sshUser" db:"ssh_user"`
	SSHAuthMethod  string       `json:"sshAuthMethod" db:"ssh_auth_method"`
	SSHKey         string       `json:"-" db:"ssh_key"`
	SSHPrivateKey  string       `json:"-" db:"ssh_private_key"`
	SSHPassword    string       `json:"-" db:"ssh_password"`
	SSHTransport   string       `json:"sshTransport" db:"ssh_transport"`
	SSHJumpHost    string       `json:"sshJumpHost,omitempty" db:"ssh_jump_host"`
	IsControlPlane bool         `json:"isControlPlane" db:"-"`
	Status         ServerStatus `json:"status" db:"status"`
	WorkerToken    string       `json:"workerToken,omitempty" db:"worker_token"`
	LastSeenAt     *time.Time   `json:"lastSeenAt,omitempty" db:"last_seen_at"`

	Metrics []byte `json:"metrics,omitempty" db:"metrics"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateServerRequest struct {
	Name          string `json:"name"`
	IPAddress     string `json:"ipAddress,omitempty"`
	IsLocal       bool   `json:"isLocal,omitempty"`
	SSHHost       string `json:"sshHost,omitempty"`
	SSHPort       int    `json:"sshPort,omitempty"`
	SSHUser       string `json:"sshUser,omitempty"`
	SSHAuthMethod string `json:"sshAuthMethod,omitempty"`
	SSHKey        string `json:"sshKey,omitempty"`
	SSHPrivateKey string `json:"sshPrivateKey,omitempty"`
	SSHPassword   string `json:"sshPassword,omitempty"`
	SSHTransport  string `json:"sshTransport,omitempty"`
	SSHJumpHost   string `json:"sshJumpHost,omitempty"`
}

type TestSSHRequest struct {
	SSHHost     string `json:"sshHost"`
	SSHPort     int    `json:"sshPort"`
	SSHUser     string `json:"sshUser"`
	SSHKey      string `json:"sshKey"`
	SSHPassword string `json:"sshPassword"`
}
