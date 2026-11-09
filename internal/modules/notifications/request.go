package notifications

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

type RegisterDeviceRequest struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

type UpdateNotificationSettingsRequest struct {
	NotificationAlerts *bool `json:"notificationAlerts"`
}

type TestNotificationRequest struct {
	Provider string `json:"provider"`
}
