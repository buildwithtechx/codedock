package models

import "time"

type ServerStatus string

const (
	ServerStatusOnline       ServerStatus = "online"
	ServerStatusOffline      ServerStatus = "offline"
	ServerStatusProvisioning ServerStatus = "provisioning"
)

type Server struct {
	ID          string       `json:"id" db:"id"`
	UserID      string       `json:"userId" db:"user_id"`
	Name        string       `json:"name" db:"name"`
	IPAddress   string       `json:"ipAddress" db:"ip_address"`
	Status      ServerStatus `json:"status" db:"status"`
	WorkerToken string       `json:"workerToken" db:"worker_token"`
	LastSeenAt  *time.Time   `json:"lastSeenAt" db:"last_seen_at"`

	Metrics []byte `json:"metrics,omitempty" db:"metrics"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
