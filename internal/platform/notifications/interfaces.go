package notifications

import (
	"context"
)

type SettingsRepository interface {
	GetServerSettings(ctx context.Context) (*ServerSettings, error)
}

type NotificationSettingsRepository interface {
	GetNotificationSettings(ctx context.Context) (*NotificationSettings, error)
	ListSubscriptions(ctx context.Context, userID, orgID string) ([]NotificationSubscription, error)
	ListDefaults(ctx context.Context, orgID string) ([]NotificationDefault, error)
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type NotificationDeliveryRepository interface {
	RecordDelivery(ctx context.Context, d *NotificationDelivery) error
}

type DeviceTokenRepository interface {
	ListTokensByUser(ctx context.Context, userID string) ([]*DeviceToken, error)
	ListAllTokens(ctx context.Context) ([]*DeviceToken, error)
}

type AppServiceRepository interface {
	GetByID(ctx context.Context, id string) (*AppService, error)
	ListWebhooksByService(ctx context.Context, serviceID string) ([]*Webhook, error)
}
