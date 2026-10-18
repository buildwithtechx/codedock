package models

type NotificationSettings struct {
	ID                    string `json:"id"`
	DiscordWebhookURL     string `json:"discordWebhookUrl,omitempty"`
	DiscordPingEnabled    bool   `json:"discordPingEnabled"`
	DiscordEnabled        bool   `json:"discordEnabled"`
	SlackWebhookURL       string `json:"slackWebhookUrl,omitempty"`
	SlackEnabled          bool   `json:"slackEnabled"`
	TelegramBotToken      string `json:"telegramBotToken,omitempty"`
	TelegramChatID        string `json:"telegramChatId,omitempty"`
	TelegramEnabled       bool   `json:"telegramEnabled"`
	SMTPHost              string `json:"smtpHost,omitempty"`
	SMTPPort              int    `json:"smtpPort,omitempty"`
	SMTPUser              string `json:"smtpUser,omitempty"`
	SMTPPassword          string `json:"smtpPassword,omitempty"`
	SMTPFromName          string `json:"smtpFromName,omitempty"`
	SMTPFromAddress       string `json:"smtpFromAddress,omitempty"`
	SMTPEnabled           bool   `json:"smtpEnabled"`
	ResendAPIKey          string `json:"resendApiKey,omitempty"`
	ResendEnabled         bool   `json:"resendEnabled"`
	PushoverUserKey       string `json:"pushoverUserKey,omitempty"`
	PushoverAPIToken      string `json:"pushoverApiToken,omitempty"`
	PushoverEnabled       bool   `json:"pushoverEnabled"`
	GenericWebhookURL     string `json:"genericWebhookUrl,omitempty"`
	GenericWebhookEnabled bool   `json:"genericWebhookEnabled"`
	NotificationAlerts    bool   `json:"notificationAlerts"`
	CreatedAt             string `json:"createdAt"`
	UpdatedAt             string `json:"updatedAt"`
}

type UpdateNotificationSettingsRequest struct {
	NotificationSettings
}

type NotificationEvent struct {
	Title     string `json:"title"`
	Message   string `json:"message"`
	Level     string `json:"level"`
	EventType string `json:"eventType"`
	ProjectID string `json:"projectId,omitempty"`
	URL       string `json:"url,omitempty"`
}

type NotificationSubscription struct {
	ID             string   `json:"id"`
	UserID         string   `json:"userId"`
	OrganizationID string   `json:"organizationId"`
	Category       string   `json:"category"`
	Channels       []string `json:"channels"`
	Enabled        bool     `json:"enabled"`
	CreatedAt      string   `json:"createdAt"`
	UpdatedAt      string   `json:"updatedAt"`
}

type NotificationDefault struct {
	OrganizationID string   `json:"organizationId"`
	Category       string   `json:"category"`
	Channels       []string `json:"channels"`
	Enabled        bool     `json:"enabled"`
	UpdatedAt      string   `json:"updatedAt"`
}

type UpsertSubscriptionRequest struct {
	Category string   `json:"category"`
	Channels []string `json:"channels"`
	Enabled  bool     `json:"enabled"`
}

type UpsertDefaultRequest struct {
	Category string   `json:"category"`
	Channels []string `json:"channels"`
	Enabled  bool     `json:"enabled"`
}

type NotificationDelivery struct {
	ID        string `json:"id" db:"id"`
	EventType string `json:"eventType" db:"event_type"`
	Channel   string `json:"channel" db:"channel"`
	Title     string `json:"title" db:"title"`
	Success   bool   `json:"success" db:"success"`
	Error     string `json:"error" db:"error"`
	CreatedAt string `json:"createdAt" db:"created_at"`
}

type DeviceToken struct {
	ID        string `json:"id" db:"id"`
	UserID    string `json:"userId" db:"user_id"`
	Token     string `json:"token" db:"token"`
	Platform  string `json:"platform" db:"platform"`
	CreatedAt string `json:"createdAt" db:"created_at"`
	UpdatedAt string `json:"updatedAt" db:"updated_at"`
}

type RegisterDeviceTokenRequest struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

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

func GetNotificationCategory(eventType string) string {
	switch {
	case eventType == EventDeploySucceeded || eventType == EventDeployFailed || eventType == EventDeployCancelled:
		return NotificationCategoryDeployment
	case eventType == EventServiceDown || eventType == EventServiceCrashLoop || eventType == EventServiceUnhealthy || eventType == EventServiceRecovered || eventType == EventServerUnreachable:
		return NotificationCategoryAppHealth
	case eventType == EventBackupFailed || eventType == EventBackupSucceeded || eventType == EventBackupRestoreCompleted:
		return NotificationCategoryBackups
	case eventType == EventDomainExpiring || eventType == EventDomainVerificationFailed:
		return NotificationCategoryDomains
	default:
		return "general"
	}
}
