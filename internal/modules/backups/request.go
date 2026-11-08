package backups

type CreateBackupRequest struct {
	BackupConfig
}

type UpdateBackupRequest struct {
	PreDeployment   *bool  `json:"preDeployment"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	DbUser          string `json:"dbUser"`
	DbPassword      string `json:"dbPassword"`
	Schedule        string `json:"schedule"`
	Timezone        string `json:"timezone"`
	Timeout         *int   `json:"timeout"`
	RetentionDays   *int   `json:"retentionDays"`
	MaxBackups      *int   `json:"maxBackups"`
	MaxStorageGB    *int   `json:"maxStorageGB"`
	S3DestinationID string `json:"s3DestinationId"`
	DatabaseID      string `json:"databaseId"`
	BackupEnabled   *bool  `json:"backupEnabled"`
	S3Enabled       *bool  `json:"s3Enabled"`
	DisableLocal    *bool  `json:"disableLocal"`
}

type S3DestinationPayload struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Provider        string `json:"provider"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	Region          string `json:"region"`
	PathPrefix      string `json:"pathPrefix"`
	IsDefault       *bool  `json:"isDefault"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type SFTPDestinationPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectID   string `json:"projectId"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PrivateKey  string `json:"privateKey"`
	PathPrefix  string `json:"pathPrefix"`
}

type CreatePolicyBatchRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schedule    string `json:"schedule"`
	Timezone    string `json:"timezone"`
	Timeout     int    `json:"timeout"`
}

type ProtectRecordRequest struct {
	Until int64 `json:"until"`
}

type ApplyRestoreRequest struct {
	Confirmation string `json:"confirmation"`
}

type VolumeRestoreRequest struct {
	TargetServerID string `json:"targetServerId"`
	TargetVolume   string `json:"targetVolume"`
}

type RestoreReviewRequest struct {
	Database string `json:"database"`
}
