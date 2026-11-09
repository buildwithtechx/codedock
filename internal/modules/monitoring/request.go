package monitoring

type UpsertSubscriptionRequest struct {
	Category string   `json:"category"`
	Channels []string `json:"channels"`
	Enabled  bool     `json:"enabled"`
}

type UpdateNotificationSettingsRequest struct {
	NotificationSettings
}

type TestNotificationRequest struct {
	Provider string `json:"provider"`
}

type AttentionActionRequest struct {
	Action string            `json:"action"`
	Params map[string]string `json:"params"`
}

type PathsToggleRequest struct {
	Enabled bool `json:"enabled"`
}
