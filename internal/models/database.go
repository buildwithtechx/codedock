package models

import (
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
