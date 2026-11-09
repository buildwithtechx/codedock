package monitoring

type AttentionSeverity string

const (
	AttentionSeverityInfo     AttentionSeverity = "info"
	AttentionSeverityWarning  AttentionSeverity = "warning"
	AttentionSeverityCritical AttentionSeverity = "critical"
)

type AttentionStatus string

const (
	AttentionStatusOpen     AttentionStatus = "open"
	AttentionStatusAcked    AttentionStatus = "acked"
	AttentionStatusResolved AttentionStatus = "resolved"
)

type ScheduledTaskStatus string

const (
	ScheduledTaskStatusActive    ScheduledTaskStatus = "active"
	ScheduledTaskStatusInactive  ScheduledTaskStatus = "inactive"
	ScheduledTaskStatusRunning   ScheduledTaskStatus = "running"
	ScheduledTaskStatusCompleted ScheduledTaskStatus = "completed"
	ScheduledTaskStatusFailed    ScheduledTaskStatus = "failed"
	ScheduledTaskStatusError     ScheduledTaskStatus = "error"
)

const (
	NotificationCategoryDeployment = "deployment"
	NotificationCategoryAppHealth  = "app_health"
	NotificationCategoryBackups    = "backups"
	NotificationCategoryDomains    = "domains"
	NotificationCategoryBilling    = "billing"
	NotificationCategoryJobs       = "jobs"
	NotificationCategoryActions    = "actions"
	NotificationCategoryMembers    = "members"

	EventDeploySucceeded          = "deploy.succeeded"
	EventDeployFailed             = "deploy.failed"
	EventDeployCancelled          = "deploy.cancelled"
	EventServiceDown              = "service.down"
	EventServiceCrashLoop         = "service.crash_loop"
	EventServiceUnhealthy         = "service.unhealthy"
	EventServiceRecovered         = "service.recovered"
	EventServerUnreachable        = "server.unreachable"
	EventBackupFailed             = "backup.failed"
	EventBackupSucceeded          = "backup.succeeded"
	EventBackupRestoreCompleted   = "backup.restore_completed"
	EventDomainExpiring           = "domain.expiring"
	EventDomainVerificationFailed = "domain.verification_failed"

	DeploymentStatusActive = "active"
	DeploymentStatusReady  = "ready"
	DeploymentStatusFailed = "failed"

	UserRoleOwner = "owner"
	UserRoleAdmin = "admin"
)
