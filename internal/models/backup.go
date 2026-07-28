package models

type BackupConfigStatus string

const (
	BackupConfigStatusActive   BackupConfigStatus = "active"
	BackupConfigStatusInactive BackupConfigStatus = "inactive"
)

type BackupRecordStatus string

const (
	BackupRecordStatusRunning   BackupRecordStatus = "running"
	BackupRecordStatusCompleted BackupRecordStatus = "completed"
	BackupRecordStatusFailed    BackupRecordStatus = "failed"
	BackupRecordStatusExpired   BackupRecordStatus = "expired"
)

type BackupConfig struct {
	ID              string             `json:"id" db:"id"`
	DatabaseID      string             `json:"databaseId,omitempty" db:"database_id"`
	ServiceID       string             `json:"serviceId,omitempty" db:"service_id"`
	VolumeName      string             `json:"volumeName,omitempty" db:"volume_name"`
	S3DestinationID string             `json:"s3DestinationId,omitempty" db:"s3_destination_id"`
	Name            string             `json:"name" db:"name"`
	Description     string             `json:"description" db:"description"`
	DbUser          string             `json:"dbUser" db:"db_user"`
	DbPassword      string             `json:"-" db:"db_password"`
	BackupEnabled   bool               `json:"backupEnabled" db:"backup_enabled"`
	S3Enabled       bool               `json:"s3Enabled" db:"s3_enabled"`
	DisableLocal    bool               `json:"disableLocal" db:"disable_local"`
	Schedule        string             `json:"schedule" db:"schedule"`
	Timezone        string             `json:"timezone" db:"timezone"`
	Timeout         int                `json:"timeout" db:"timeout"`
	RetentionDays   int                `json:"retentionDays" db:"retention_days"`
	MaxBackups      int                `json:"maxBackups" db:"max_backups"`
	MaxStorageGB    int                `json:"maxStorageGb" db:"max_storage_gb"`
	Status          BackupConfigStatus `json:"status" db:"status"`
	CreatedAt       string             `json:"createdAt" db:"created_at"`
	UpdatedAt       string             `json:"updatedAt" db:"updated_at"`
}

type BackupRecord struct {
	ID              string             `json:"id" db:"id"`
	BackupConfigID  string             `json:"backupConfigId" db:"backup_config_id"`
	DatabaseID      string             `json:"databaseId,omitempty" db:"database_id"`
	S3DestinationID string             `json:"s3DestinationId,omitempty" db:"s3_destination_id"`
	Status          BackupRecordStatus `json:"status" db:"status"`
	FilePath        string             `json:"-" db:"file_path"`
	FileSizeBytes   int64              `json:"fileSizeBytes" db:"file_size_bytes"`
	S3URL           string             `json:"s3Url,omitempty" db:"s3_url"`
	Logs            string             `json:"logs" db:"logs"`
	StartedAt       string             `json:"startedAt" db:"started_at"`
	CompletedAt     string             `json:"completedAt" db:"completed_at"`
}

type UpdateBackupRecordOpts struct {
	ID              string             `json:"id"`
	Status          BackupRecordStatus `json:"status"`
	FilePath        string             `json:"-"`
	S3URL           string             `json:"s3_url"`
	S3DestinationID string             `json:"s3_destination_id"`
	Logs            string             `json:"logs"`
	FileSizeBytes   int64              `json:"file_size_bytes"`
	CompletedAt     string             `json:"completed_at"`
}

type S3Destination struct {
	ID              string `json:"id" db:"id"`
	Name            string `json:"name" db:"name"`
	Description     string `json:"description" db:"description"`
	Provider        string `json:"provider" db:"provider"`
	Endpoint        string `json:"endpoint" db:"endpoint"`
	Bucket          string `json:"bucket" db:"bucket"`
	Region          string `json:"region" db:"region"`
	AccessKeyID     string `json:"accessKeyId" db:"access_key_id"`
	SecretAccessKey string `json:"-" db:"secret_access_key"`
	CreatedAt       string `json:"createdAt" db:"created_at"`
}
