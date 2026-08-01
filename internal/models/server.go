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
	UserID         string       `json:"userId" db:"user_id"`
	Name           string       `json:"name" db:"name"`
	IPAddress      string       `json:"ipAddress" db:"ip_address"`
	SSHHost        string       `json:"sshHost" db:"ssh_host"`
	SSHPort        int          `json:"sshPort" db:"ssh_port"`
	SSHUser        string       `json:"sshUser" db:"ssh_user"`
	SSHKey         string       `json:"-" db:"ssh_key"`
	SSHPassword    string       `json:"-" db:"ssh_password"`
	IsControlPlane bool         `json:"isControlPlane" db:"-"`
	Status         ServerStatus `json:"status" db:"status"`
	WorkerToken    string       `json:"workerToken,omitempty" db:"worker_token"`
	LastSeenAt     *time.Time   `json:"lastSeenAt,omitempty" db:"last_seen_at"`

	Metrics []byte `json:"metrics,omitempty" db:"metrics"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CreateServerRequest struct {
	Name        string `json:"name"`
	IPAddress   string `json:"ipAddress"`
	SSHHost     string `json:"sshHost"`
	SSHPort     int    `json:"sshPort"`
	SSHUser     string `json:"sshUser"`
	SSHKey      string `json:"sshKey"`
	SSHPassword string `json:"sshPassword"`
}

type TestSSHRequest struct {
	SSHHost     string `json:"sshHost"`
	SSHPort     int    `json:"sshPort"`
	SSHUser     string `json:"sshUser"`
	SSHKey      string `json:"sshKey"`
	SSHPassword string `json:"sshPassword"`
}
