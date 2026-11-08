package backups

type ConfigStatus string

const (
	ConfigStatusActive   ConfigStatus = "active"
	ConfigStatusInactive ConfigStatus = "inactive"
)

const (
	StatusActive   = string(ConfigStatusActive)
	StatusInactive = string(ConfigStatusInactive)
)

type RecordStatus string

const (
	RecordStatusRunning   RecordStatus = "running"
	RecordStatusCompleted RecordStatus = "completed"
	RecordStatusFailed    RecordStatus = "failed"
	RecordStatusExpired   RecordStatus = "expired"
	RecordStatusExpiring  RecordStatus = "expiring"
)

const (
	ProviderAWS        = "aws"
	ProviderWasabi     = "wasabi"
	ProviderBackblaze  = "backblaze"
	ProviderCloudflare = "cloudflare"
	ProviderMinIO      = "minio"
	ProviderCustom     = "custom"

	DefaultSchedule = "0 2 * * *"
	DefaultTimeout  = 3600
	DefaultLimit    = 50
)
