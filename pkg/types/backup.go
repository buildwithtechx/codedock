package types

import "time"

type BackupSchedule struct {
	ID            string    `json:"id" db:"id"`
	ProjectID     string    `json:"projectId" db:"project_id"`
	ServiceID     string    `json:"serviceId,omitempty" db:"service_id"`
	DatabaseID    string    `json:"databaseId,omitempty" db:"database_id"`
	DestinationID string    `json:"destinationId,omitempty" db:"destination_id"`
	Name          string    `json:"name,omitempty" db:"name"`
	Schedule      string    `json:"schedule,omitempty" db:"schedule"`
	Status        string    `json:"status,omitempty" db:"status"`
	Cron          string    `json:"cron" db:"cron"`
	RetentionDays int       `json:"retentionDays" db:"retention_days"`
	Enabled       bool      `json:"enabled" db:"enabled"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time `json:"updatedAt" db:"updated_at"`
}

type BackupConfig = BackupSchedule

type BackupExecution struct {
	ID            string     `json:"id" db:"id"`
	ScheduleID    string     `json:"scheduleId" db:"schedule_id"`
	Status        string     `json:"status" db:"status"`
	SizeBytes     int64      `json:"sizeBytes" db:"size_bytes"`
	FileSizeBytes int64      `json:"fileSizeBytes,omitempty" db:"file_size_bytes"`
	StorageKey    string     `json:"storageKey" db:"storage_key"`
	StartedAt     time.Time  `json:"startedAt" db:"started_at"`
	FinishedAt    *time.Time `json:"finishedAt,omitempty" db:"finished_at"`
	CompletedAt   string     `json:"completedAt,omitempty" db:"completed_at"`
	Error         string     `json:"error,omitempty" db:"error"`
}

type BackupRecord = BackupExecution

type S3Destination struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Bucket    string    `json:"bucket" db:"bucket"`
	Region    string    `json:"region" db:"region"`
	Endpoint  string    `json:"endpoint,omitempty" db:"endpoint"`
	AccessKey string    `json:"accessKey" db:"access_key"`
	SecretKey string    `json:"-" db:"secret_key"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
