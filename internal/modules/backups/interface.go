package backups

import "context"

// Vault handles encryption and decryption of sensitive backup configs.
type Vault interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}

// ConfigRepository manages configuration and history records for backups.
type ConfigRepository interface {
	CreateConfig(ctx context.Context, cfg *BackupConfig) error
	UpdateConfig(ctx context.Context, cfg *BackupConfig) error
	GetConfigByID(ctx context.Context, id string) (*BackupConfig, error)
	GetConfigByDatabaseID(ctx context.Context, dbID string) (*BackupConfig, error)
	ListConfigs(ctx context.Context) ([]*BackupConfig, error)
	ListAllActiveConfigs(ctx context.Context) ([]*BackupConfig, error)
	DeleteConfig(ctx context.Context, id string) error
	CreateRecord(ctx context.Context, rec *BackupRecord) error
	GetRecordByID(ctx context.Context, id string) (*BackupRecord, error)
	ListRecordsByConfig(ctx context.Context, backupConfigID string) ([]*BackupRecord, error)
	ListRecordsByDatabase(ctx context.Context, databaseID string) ([]*BackupRecord, error)
	ListAllRecords(ctx context.Context, limit int) ([]*BackupRecord, error)
	ListRecordsByConfigs(ctx context.Context, configIDs []string, limit int) ([]*BackupRecord, error)
	UpdateRecord(ctx context.Context, rec *BackupRecord) error
	DeleteRecord(ctx context.Context, id string) error
}

// S3Repository manages S3-compatible remote storage destinations.
type S3Repository interface {
	CreateS3Destination(ctx context.Context, dest *S3Destination) error
	UpdateS3Destination(ctx context.Context, dest *S3Destination) error
	ListS3Destinations(ctx context.Context) ([]*S3Destination, error)
	GetS3Destination(ctx context.Context, id string) (*S3Destination, error)
	DeleteS3Destination(ctx context.Context, id string) error
	SetDefaultDestination(ctx context.Context, id string) error
	RecordVerificationResult(ctx context.Context, id string, verified bool, errMsg string) error
}

// BatchRepository manages grouped backup execution batches.
type BatchRepository interface {
	Create(ctx context.Context, batch *BackupPolicyBatch) error
	Get(ctx context.Context, id string) (*BackupPolicyBatch, error)
	ListByProject(ctx context.Context, projectID string) ([]BackupPolicyBatch, error)
	ListMembers(ctx context.Context, batchID string) ([]*BackupConfig, error)
}

// SFTPRepository manages SFTP destination targets.
type SFTPRepository interface {
	Create(ctx context.Context, dest *SFTPDestination) error
	Get(ctx context.Context, id string) (*SFTPDestination, error)
	ListByProject(ctx context.Context, projectID string) ([]SFTPDestination, error)
	Delete(ctx context.Context, id string) error
}

// Repository aggregates all backup storage interfaces and disaster recovery methods.
type Repository interface {
	ConfigRepository
	S3Repository
	BatchRepository
	ProtectRecord(ctx context.Context, id string, until int64) error
	ClearRestoreProtection(ctx context.Context, id string, before int64) error
	ClaimRecordExpiry(ctx context.Context, id string) (bool, error)
	SetSnapshotDumper(dumper func(ctx context.Context) ([]byte, error))
	ControlPlaneSnapshot(ctx context.Context) ([]byte, error)
	RecoverRecords(ctx context.Context) error
}

// OperationsService provides approval staging and audit review for backups.
type OperationsService interface {
	Review(ctx context.Context, user, project, kind, target, payload, snapshot, effects string) (any, error)
	Apply(ctx context.Context, id, user, confirmation, snapshot string, runner func(ctx context.Context) error) error
	GetPayload(ctx context.Context, id string) (string, error)
}
