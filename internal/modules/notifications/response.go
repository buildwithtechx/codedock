package notifications

// SettingsResponse wraps notification settings.
type SettingsResponse struct {
	Settings *NotificationSettings `json:"settings"`
}

// SubscriptionsResponse wraps user subscription channel assignments.
type SubscriptionsResponse struct {
	Subscriptions []NotificationSubscription `json:"subscriptions"`
}

// DefaultsResponse wraps organization default notifications.
type DefaultsResponse struct {
	Defaults []NotificationDefault `json:"defaults"`
}

// DeliveriesResponse wraps delivery logs.
type DeliveriesResponse struct {
	Deliveries []NotificationDelivery `json:"deliveries"`
	Total      int                    `json:"total"`
}

// DeviceTokensResponse wraps registered device push tokens.
type DeviceTokensResponse struct {
	Tokens []*DeviceToken `json:"tokens"`
}
