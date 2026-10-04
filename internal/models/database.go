package models

import (
	"time"

	"codedock.run/codedock/pkg/types"
)

type Database = types.Database
type CreateDatabaseRequest = types.CreateDatabaseRequest
type DatabaseStatus = types.DatabaseStatus
type DatabaseEngine = types.DatabaseEngine

const (
	DatabaseEnginePostgres   DatabaseEngine = "postgres"
	DatabaseEngineMySQL      DatabaseEngine = "mysql"
	DatabaseEngineRedis      DatabaseEngine = "redis"
	DatabaseEngineMongoDB    DatabaseEngine = "mongodb"
	DatabaseEngineMariaDB    DatabaseEngine = "mariadb"
	DatabaseEngineClickhouse DatabaseEngine = "clickhouse"
)

const (
	DatabaseStatusCreated = "created"
	DatabaseStatusRunning = "running"
	DatabaseStatusStopped = "stopped"
	DatabaseStatusError   = "error"
)

type UpdateDatabaseRequest struct {
	Name               string  `json:"name,omitempty"`
	Version            string  `json:"version,omitempty"`
	ExternalDNS        string  `json:"externalDns,omitempty"`
	CPULimit           float64 `json:"cpuLimit,omitempty"`
	MemoryLimit        int     `json:"memoryLimit,omitempty"`
	LogicalReplication bool    `json:"logicalReplication,omitempty"`
	CustomArgs         string  `json:"customArgs,omitempty"`
}

type QueryDatabaseRequest struct {
	Query string `json:"query"`
}

type DatabaseQueryRequest = QueryDatabaseRequest

type QueryDatabaseResponse struct {
	Columns         []string         `json:"columns"`
	Rows            []map[string]any `json:"rows"`
	RowCount        int              `json:"rowCount"`
	ExecutionTimeMs int64            `json:"executionTimeMs"`
	Result          any              `json:"result,omitempty"`
}

type DatabaseQueryResponse = QueryDatabaseResponse

type DatabaseTableSchema struct {
	Name    string               `json:"name"`
	Columns []DatabaseColumnInfo `json:"columns"`
}

type TableSchema = DatabaseTableSchema

type DatabaseColumnInfo struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsNullable bool   `json:"isNullable,omitempty"`
	IsPrimary  bool   `json:"isPrimary,omitempty"`
}

type ColumnSchema = DatabaseColumnInfo

type TableRowPayload map[string]any

type ImportDatabaseRequest struct {
	SQL       string `json:"sql,omitempty"`
	SourceURL string `json:"sourceUrl,omitempty"`
}

type ClusterDatabase struct {
	ID              string    `json:"id" db:"id"`
	OrganizationID  string    `json:"organizationId" db:"organization_id"`
	ProjectID       string    `json:"projectId" db:"project_id"`
	ServerID        string    `json:"serverId,omitempty" db:"server_id"`
	Name            string    `json:"name" db:"name"`
	Engine          string    `json:"engine" db:"engine"`
	Version         string    `json:"version" db:"version"`
	Status          string    `json:"status" db:"status"`
	Intent          string    `json:"intent" db:"intent"`
	Port            int       `json:"port" db:"port"`
	Username        string    `json:"username" db:"username"`
	DatabaseName    string    `json:"databaseName" db:"database_name"`
	SecretEncrypted string    `json:"-" db:"secret_encrypted"`
	EnvKey          string    `json:"envKey,omitempty" db:"env_key"`
	ContainerID     string    `json:"containerId,omitempty" db:"container_id"`
	InternalDNS     string    `json:"internalDns,omitempty" db:"internal_dns"`
	ExternalDNS     string    `json:"externalDns,omitempty" db:"external_dns"`
	ProgressJSON    string    `json:"progress,omitempty" db:"progress_json"`
	BackupRequestID string    `json:"backupRequestId,omitempty" db:"backup_request_id"`
	CPULimit        float64   `json:"cpuLimit,omitempty" db:"cpu_limit"`
	MemoryLimit     int       `json:"memoryLimit,omitempty" db:"memory_limit"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}
