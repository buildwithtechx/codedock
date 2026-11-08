package databases

import "time"

type Database struct {
	ID                 string         `json:"id" db:"id"`
	ProjectID          string         `json:"projectId" db:"project_id"`
	EnvironmentID      string         `json:"environmentId,omitempty" db:"environment_id"`
	Name               string         `json:"name" db:"name"`
	Engine             DatabaseEngine `json:"engine" db:"engine"`
	Version            string         `json:"version" db:"version"`
	Username           string         `json:"username" db:"username"`
	Password           string         `json:"-" db:"password"`
	EncryptedPassword  string         `json:"-" db:"encrypted_password"`
	DatabaseName       string         `json:"databaseName" db:"database_name"`
	InternalHost       string         `json:"internalHost" db:"internal_host"`
	InternalDNS        string         `json:"internalDns,omitempty" db:"internal_dns"`
	ExternalDNS        string         `json:"externalDns,omitempty" db:"external_dns"`
	VolumePath         string         `json:"volumePath,omitempty" db:"volume_path"`
	ContainerID        string         `json:"containerId,omitempty" db:"container_id"`
	CustomArgs         string         `json:"customArgs,omitempty" db:"custom_args"`
	LogicalReplication bool           `json:"logicalReplication,omitempty" db:"logical_replication"`
	CPULimit           float64        `json:"cpuLimit,omitempty" db:"cpu_limit"`
	MemoryLimit        int            `json:"memoryLimit,omitempty" db:"memory_limit"`
	Port               int            `json:"port" db:"port"`
	ExternalPort       int            `json:"externalPort,omitempty" db:"external_port"`
	Status             DatabaseStatus `json:"status" db:"status"`
	CreatedAt          time.Time      `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time      `json:"updatedAt" db:"updated_at"`
}

type TableSchema struct {
	Name    string         `json:"name"`
	Columns []ColumnSchema `json:"columns"`
}

type ColumnSchema struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsNullable bool   `json:"isNullable,omitempty"`
	IsPrimary  bool   `json:"isPrimary,omitempty"`
}

type AuditLogEntry struct {
	UserID    string
	Action    string
	Resource  string
	IPAddress string
	Details   map[string]any
}
