package notifications

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
)
