package core

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"codedock/internal/models"
	"codedock/internal/notifications"
	"codedock/internal/repositories"
)

type DispatcherService struct {
	settingsRepo repositories.SettingsRepository
	notifRepo    repositories.NotificationSettingsRepository
	userRepo     repositories.UserRepository
	deliveries   repositories.NotificationDeliveryRepository
	deviceTokens repositories.DeviceTokenRepository
	mailer       interface {
		SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error
	}
}

type Mailer interface {
	SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error
}

func NewDispatcherService(
	settingsRepo repositories.SettingsRepository,
	notifRepo repositories.NotificationSettingsRepository,
	userRepo repositories.UserRepository,
	deliveries repositories.NotificationDeliveryRepository,
	mailer Mailer,
	deviceTokens repositories.DeviceTokenRepository,
) *DispatcherService {
	return &DispatcherService{
		settingsRepo: settingsRepo,
		notifRepo:    notifRepo,
		userRepo:     userRepo,
		deliveries:   deliveries,
		deviceTokens: deviceTokens,
		mailer:       mailer,
	}
}

func (d *DispatcherService) Dispatch(event *models.NotificationEvent) {
	go func() {
		if err := d.Send(event); err != nil {
			slog.Error("failed to dispatch event", "title", event.Title, "type", event.EventType, "err", err)
		}
	}()
}

func (d *DispatcherService) Send(event *models.NotificationEvent) error {
	if strings.HasPrefix(event.EventType, "test_global_") {
		return d.sendGlobalTest(event)
	}

	ctx := context.Background()
	settings, err := d.notifRepo.GetNotificationSettings(ctx)
	if err != nil || settings == nil {
		return fmt.Errorf("could not fetch notification settings: %v", err)
	}

	if !settings.NotificationAlerts {
		return nil
	}

	var dispatchErrs []error

	// 1. Fan out to configured global notification channels
	providers := []string{"discord", "slack", "telegram", "generic", "pushover", "expo"}
	for _, provider := range providers {
		if d.isProviderConfigured(settings, provider) {
			sendErr := d.sendToProvider(settings, provider, event)
			d.recordDelivery(ctx, event, provider, sendErr)
			if sendErr != nil {
				dispatchErrs = append(dispatchErrs, fmt.Errorf("%s: %w", provider, sendErr))
			}
		}
	}

	// 2. Fan out to registered Expo push device tokens, consulting user category preferences
	if d.deviceTokens != nil {
		tokens, err := d.deviceTokens.ListAllTokens(ctx)
		if err == nil && len(tokens) > 0 {
			category := models.GetNotificationCategory(event.EventType)
			userPrefsCache := make(map[string]bool)

			for _, dt := range tokens {
				if !notifications.IsExpoPushToken(dt.Token) {
					continue
				}

				// Consult per-user category subscriptions before sending
				if dt.UserID != "" && d.notifRepo != nil {
					allowed, cached := userPrefsCache[dt.UserID]
					if !cached {
						allowed = d.isUserSubscribed(ctx, dt.UserID, event.EventType, category)
						userPrefsCache[dt.UserID] = allowed
					}
					if !allowed {
						continue
					}
				}

				pushErr := notifications.SendExpoNotification(dt.Token, event)
				d.recordDelivery(ctx, event, "expo", pushErr)
				if pushErr != nil {
					dispatchErrs = append(dispatchErrs, fmt.Errorf("expo: %w", pushErr))
				}
			}
		}
	}

	if len(dispatchErrs) > 0 {
		return fmt.Errorf("dispatch completed with %d error(s): %v", len(dispatchErrs), dispatchErrs)
	}
	return nil
}

func (d *DispatcherService) isProviderConfigured(settings *models.NotificationSettings, provider string) bool {
	switch provider {
	case "discord":
		return settings.DiscordEnabled && settings.DiscordWebhookURL != ""
	case "slack":
		return settings.SlackEnabled && settings.SlackWebhookURL != ""
	case "telegram":
		return settings.TelegramEnabled && settings.TelegramBotToken != "" && settings.TelegramChatID != ""
	case "generic":
		return settings.GenericWebhookEnabled && settings.GenericWebhookURL != ""
	case "pushover":
		return settings.PushoverEnabled && settings.PushoverAPIToken != "" && settings.PushoverUserKey != ""
	case "expo":
		return settings.ExpoEnabled && settings.ExpoPushToken != ""
	}
	return false
}

func (d *DispatcherService) sendGlobalTest(event *models.NotificationEvent) error {
	ctx := context.Background()
	settings, err := d.notifRepo.GetNotificationSettings(ctx)
	if err != nil || settings == nil {
		return fmt.Errorf("could not fetch notification settings: %v", err)
	}

	provider := event.EventType[len("test_global_"):]
	err = d.sendToProvider(settings, provider, event)
	d.recordDelivery(ctx, event, provider, err)
	return err
}

func (d *DispatcherService) recordDelivery(ctx context.Context, event *models.NotificationEvent, channel string, err error) {
	if d.deliveries == nil {
		return
	}
	entry := &models.NotificationDelivery{EventType: event.EventType, Channel: channel, Title: event.Title, Success: err == nil}
	if err != nil {
		entry.Error = err.Error()
	}
	if recErr := d.deliveries.RecordDelivery(ctx, entry); recErr != nil {
		slog.Error("failed to record notification delivery", "err", recErr)
	}
}

func (d *DispatcherService) sendToProvider(settings *models.NotificationSettings, provider string, event *models.NotificationEvent) error {
	switch provider {
	case "discord":
		if settings.DiscordEnabled && settings.DiscordWebhookURL != "" {
			return notifications.SendDiscordNotification(settings.DiscordWebhookURL, event)
		}
	case "slack":
		if settings.SlackEnabled && settings.SlackWebhookURL != "" {
			return notifications.SendSlackNotification(settings.SlackWebhookURL, event)
		}
	case "telegram":
		if settings.TelegramEnabled && settings.TelegramBotToken != "" && settings.TelegramChatID != "" {
			return notifications.SendTelegramNotification(settings.TelegramBotToken, settings.TelegramChatID, event)
		}
	case "generic":
		if settings.GenericWebhookEnabled && settings.GenericWebhookURL != "" {
			return notifications.SendGenericWebhook(settings.GenericWebhookURL, event)
		}
	case "pushover":
		if settings.PushoverEnabled && settings.PushoverAPIToken != "" && settings.PushoverUserKey != "" {
			return notifications.SendPushoverNotification(settings.PushoverAPIToken, settings.PushoverUserKey, event)
		}
	case "expo":
		if settings.ExpoEnabled && settings.ExpoPushToken != "" {
			return notifications.SendExpoNotification(settings.ExpoPushToken, event)
		}
	}

	return fmt.Errorf("provider not enabled or configured")
}

func (d *DispatcherService) isUserSubscribed(ctx context.Context, userID, eventType, category string) bool {
	if d.notifRepo == nil || userID == "" {
		return true
	}

	subs, err := d.notifRepo.ListSubscriptions(ctx, userID, "")
	if err != nil {
		slog.Warn("failed to fetch user notification subscriptions", "userID", userID, "err", err)
		return true
	}

	if len(subs) == 0 {
		return true
	}

	// 1. Check specific event match (highest precedence, e.g. "deploy.succeeded")
	for _, sub := range subs {
		if sub.Category == eventType {
			return sub.Enabled
		}
	}

	// 2. Check general category group match (e.g. "deployment", "app_health", "backups")
	for _, sub := range subs {
		if sub.Category == category {
			return sub.Enabled
		}
	}

	return true
}
