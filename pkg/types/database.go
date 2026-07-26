package types

import "time"

type DatabaseEngine = string
type DatabaseStatus = string

type Database struct {
	ID                 string         `json:"id" db:"id"`
	ProjectID          string         `json:"projectId" db:"project_id"`
	EnvironmentID      string         `json:"environmentId,omitempty" db:"environment_id"`
	Name               string         `json:"name" db:"name"`
	Engine             DatabaseEngine `json:"engine" db:"engine"`
	Version            string         `json:"version" db:"version"`
	Username           string         `json:"username" db:"username"`
	Password           string         `json:"password,omitempty" db:"password"`
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

type CreateDatabaseRequest struct {
	ProjectID          string         `json:"projectId"`
	EnvironmentID      string         `json:"environmentId,omitempty"`
	Name               string         `json:"name"`
	Engine             DatabaseEngine `json:"engine"`
	Version            string         `json:"version"`
	DatabaseName       string         `json:"databaseName"`
	Username           string         `json:"username"`
	Password           string         `json:"password"`
	Port               int            `json:"port,omitempty"`
	VolumePath         string         `json:"volumePath,omitempty"`
	CustomArgs         string         `json:"customArgs,omitempty"`
	LogicalReplication bool           `json:"logicalReplication,omitempty"`
}

type ImportDatabaseRequest struct {
	SQL       string `json:"sql,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
}
