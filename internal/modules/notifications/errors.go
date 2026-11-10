package notifications

import "errors"

var (
	ErrNotificationSettingsNotFound = errors.New("notification settings not found")
	ErrInvalidProvider              = errors.New("invalid notification provider")
	ErrDeviceTokenRequired          = errors.New("device token is required")
	ErrPlatformRequired             = errors.New("device platform is required")
	ErrDeliveryRecordFailed         = errors.New("failed to record notification delivery")
)
