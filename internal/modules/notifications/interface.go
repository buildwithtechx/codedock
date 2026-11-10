package notifications

import "context"

// NotificationService executes immediate test notification dispatches.
type NotificationService interface {
	TestNotification(ctx context.Context, provider string) error
}

// NotificationSettingsOperations orchestrates settings, deliveries, and device tokens.
type NotificationSettingsOperations interface {
	GetSettings(ctx context.Context) (*NotificationSettings, error)
	UpdateSettings(ctx context.Context, req UpdateNotificationSettingsRequest) (*NotificationSettings, error)
	ListDeliveries(ctx context.Context, limit int) (any, error)
	RegisterDeviceToken(ctx context.Context, userID, token, platform string) error
}

// NotificationSettingsRepository provides persistence for notification configurations and subscriptions.
type NotificationSettingsRepository interface {
	GetNotificationSettings(ctx context.Context) (*NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, cfg *NotificationSettings) error
	ListSubscriptions(ctx context.Context, userID, orgID string) ([]NotificationSubscription, error)
	UpsertSubscription(ctx context.Context, sub *NotificationSubscription) error
	ListDefaults(ctx context.Context, orgID string) ([]NotificationDefault, error)
	UpsertDefault(ctx context.Context, def *NotificationDefault) error
}

// NotificationDeliveryRepository provides persistence for notification audit deliveries.
type NotificationDeliveryRepository interface {
	RecordDelivery(ctx context.Context, delivery *NotificationDelivery) error
	ListDeliveries(ctx context.Context, limit int) ([]NotificationDelivery, error)
}

// DeviceTokenRepository manages user mobile/desktop push device registration tokens.
type DeviceTokenRepository interface {
	RegisterToken(ctx context.Context, userID, token, platform string) (*DeviceToken, error)
	DeleteToken(ctx context.Context, userID, token string) error
	ListTokensByUser(ctx context.Context, userID string) ([]*DeviceToken, error)
	ListAllTokens(ctx context.Context) ([]*DeviceToken, error)
}
