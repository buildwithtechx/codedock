package migrations

type MigrationMode string

const (
	MigrationModeMove MigrationMode = "move"
	MigrationModeCopy MigrationMode = "copy"
)

type MigrationStatus string

const (
	MigrationStatusPending         MigrationStatus = "pending"
	MigrationStatusScanning        MigrationStatus = "scanning"
	MigrationStatusPreviewed       MigrationStatus = "previewed"
	MigrationStatusTransferring    MigrationStatus = "transferring"
	MigrationStatusVerifying       MigrationStatus = "verifying"
	MigrationStatusAwaitingCutover MigrationStatus = "awaiting_cutover"
	MigrationStatusCompleted       MigrationStatus = "completed"
	MigrationStatusFailed          MigrationStatus = "failed"
	MigrationStatusCancelled       MigrationStatus = "cancelled"
)

type MigrationPromptKind string

const (
	MigrationPromptVolumeConflict MigrationPromptKind = "volume-conflict"
	MigrationPromptPortConflict   MigrationPromptKind = "port-conflict"
)

const (
	MaxSetupBundleUploadBytes = 32 * 1024 * 1024  // 32MB pre-auth setup upload limit
	MaxBundleUploadBytes      = 500 * 1024 * 1024 // 500MB authenticated upload limit
)
