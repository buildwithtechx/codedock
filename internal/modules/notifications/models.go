package notifications

import (
	"time"

	platformnotifs "codedock/internal/platform/notifications"
)

type NotificationDelivery = platformnotifs.NotificationDelivery
type NotificationSubscription = platformnotifs.NotificationSubscription
type NotificationDefault = platformnotifs.NotificationDefault
type NotificationSettings = platformnotifs.NotificationSettings

type NotificationEvent struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
}

type DeviceToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"userId" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	Platform  string    `json:"platform" db:"platform"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
