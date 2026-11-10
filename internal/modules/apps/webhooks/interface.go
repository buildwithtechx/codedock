package webhooks

import "context"

// Repository defines database operations for webhooks and log drains.
type Repository interface {
	CreateWebhook(ctx context.Context, w *Webhook) error
	ListWebhooksByService(ctx context.Context, serviceID string) ([]*Webhook, error)
	DeleteWebhook(ctx context.Context, id, serviceID string) error
	CreateLogDrain(ctx context.Context, d *LogDrain) error
	ListLogDrainsByService(ctx context.Context, serviceID string) ([]*LogDrain, error)
	DeleteLogDrain(ctx context.Context, id, serviceID string) error
}

// WebhookService provides application business logic for webhooks and log drains.
type WebhookService interface {
	CreateWebhook(ctx context.Context, w *Webhook) (*Webhook, error)
	ListWebhooks(ctx context.Context, serviceID string) ([]*Webhook, error)
	DeleteWebhook(ctx context.Context, id, serviceID string) error
	CreateLogDrain(ctx context.Context, d *LogDrain) (*LogDrain, error)
	ListLogDrains(ctx context.Context, serviceID string) ([]*LogDrain, error)
	DeleteLogDrain(ctx context.Context, id, serviceID string) error
}
