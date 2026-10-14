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
	mailer       interface {
		SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error
	}
}

type Mailer interface {
	SendSystemEmail(ctx context.Context, templateName string, toAddress string, subject string, data any) error
}

func NewDispatcherService(settingsRepo repositories.SettingsRepository, notifRepo repositories.NotificationSettingsRepository, userRepo repositories.UserRepository, deliveries repositories.NotificationDeliveryRepository, mailer Mailer) *DispatcherService {
	return &DispatcherService{settingsRepo: settingsRepo, notifRepo: notifRepo, userRepo: userRepo, deliveries: deliveries, mailer: mailer}
}

func (d *DispatcherService) Dispatch(event *models.NotificationEvent) {
	go func() {
		if err := d.Send(event); err != nil {
			slog.Error("failed to dispatch event", "title", event.Title, "err", err)
		}
	}()
}

func (d *DispatcherService) Send(event *models.NotificationEvent) error {
	if strings.HasPrefix(event.EventType, "test_global_") {
		return d.sendGlobalTest(event)
	}
	return nil
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
	}

	return fmt.Errorf("provider not enabled or configured")
}
