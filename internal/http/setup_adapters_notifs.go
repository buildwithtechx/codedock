package http

import (
	"context"

	"codedock/internal/models"
	"codedock/internal/modules/monitoring"
	"codedock/internal/repositories"
	projectservices "codedock/internal/services/projects"
	systemservices "codedock/internal/services/system"
)

type notifSettingsServiceAdapter struct {
	svc     *systemservices.NotificationSettingsService
	orgRepo repositories.OrganizationRepository
}

func (n *notifSettingsServiceAdapter) GetNotificationSettings(ctx context.Context) (*monitoring.NotificationSettings, error) {
	if n.svc == nil {
		return &monitoring.NotificationSettings{}, nil
	}
	s, err := n.svc.GetNotificationSettings(ctx)
	if err != nil || s == nil {
		return nil, err
	}
	return &monitoring.NotificationSettings{
		ID:                 s.ID,
		DiscordWebhookURL:  s.DiscordWebhookURL,
		DiscordPingEnabled: s.DiscordPingEnabled,
		DiscordEnabled:     s.DiscordEnabled,
		SlackWebhookURL:    s.SlackWebhookURL,
		SlackEnabled:       s.SlackEnabled,
		TelegramBotToken:   s.TelegramBotToken,
		TelegramChatID:     s.TelegramChatID,
		TelegramEnabled:    s.TelegramEnabled,
		SMTPHost:           s.SMTPHost,
		SMTPPort:           s.SMTPPort,
		SMTPUser:           s.SMTPUser,
		SMTPPassword:       s.SMTPPassword,
		SMTPFromName:       s.SMTPFromName,
		SMTPFromAddress:    s.SMTPFromAddress,
		SMTPEnabled:        s.SMTPEnabled,
		ResendAPIKey:       s.ResendAPIKey,
		ResendEnabled:      s.ResendEnabled,
		PushoverUserKey:    s.PushoverUserKey,
		PushoverAPIToken:   s.PushoverAPIToken,
	}, nil
}

func (n *notifSettingsServiceAdapter) UpdateNotificationSettings(ctx context.Context, settings *monitoring.NotificationSettings) error {
	if n.svc == nil {
		return nil
	}
	return n.svc.UpdateNotificationSettings(ctx, &models.NotificationSettings{
		ID:                 settings.ID,
		DiscordWebhookURL:  settings.DiscordWebhookURL,
		DiscordPingEnabled: settings.DiscordPingEnabled,
		DiscordEnabled:     settings.DiscordEnabled,
		SlackWebhookURL:    settings.SlackWebhookURL,
		SlackEnabled:       settings.SlackEnabled,
		TelegramBotToken:   settings.TelegramBotToken,
		TelegramChatID:     settings.TelegramChatID,
		TelegramEnabled:    settings.TelegramEnabled,
		SMTPHost:           settings.SMTPHost,
		SMTPPort:           settings.SMTPPort,
		SMTPUser:           settings.SMTPUser,
		SMTPPassword:       settings.SMTPPassword,
		SMTPFromName:       settings.SMTPFromName,
		SMTPFromAddress:    settings.SMTPFromAddress,
		SMTPEnabled:        settings.SMTPEnabled,
		ResendAPIKey:       settings.ResendAPIKey,
		ResendEnabled:      settings.ResendEnabled,
		PushoverUserKey:    settings.PushoverUserKey,
		PushoverAPIToken:   settings.PushoverAPIToken,
	})
}

func (n *notifSettingsServiceAdapter) ListSubscriptions(ctx context.Context, userID, orgID string) (any, error) {
	if n.svc == nil {
		return []models.NotificationSubscription{}, nil
	}
	return n.svc.ListSubscriptions(ctx, userID, orgID)
}

func (n *notifSettingsServiceAdapter) UpsertSubscription(ctx context.Context, userID, orgID, category string, channels []string, enabled bool) (any, error) {
	if n.svc == nil {
		return nil, nil
	}
	sub := &models.NotificationSubscription{
		UserID:         userID,
		OrganizationID: orgID,
		Category:       category,
		Channels:       channels,
		Enabled:        enabled,
	}
	if err := n.svc.UpsertSubscription(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (n *notifSettingsServiceAdapter) ListDefaults(ctx context.Context, orgID string) (any, error) {
	if n.svc == nil {
		return []models.NotificationDefault{}, nil
	}
	return n.svc.ListDefaults(ctx, orgID)
}

func (n *notifSettingsServiceAdapter) UpsertDefault(ctx context.Context, orgID, category string, channels []string, enabled bool) (any, error) {
	if n.svc == nil {
		return nil, nil
	}
	def := &models.NotificationDefault{
		OrganizationID: orgID,
		Category:       category,
		Channels:       channels,
		Enabled:        enabled,
	}
	if err := n.svc.UpsertDefault(ctx, def); err != nil {
		return nil, err
	}
	return def, nil
}

func (n *notifSettingsServiceAdapter) ListDeliveries(ctx context.Context, limit int) (any, error) {
	if n.svc == nil {
		return []models.NotificationDelivery{}, nil
	}
	return n.svc.ListDeliveries(ctx, limit)
}

func (n *notifSettingsServiceAdapter) RegisterDeviceToken(ctx context.Context, userID, token, platform string) (any, error) {
	if n.svc == nil {
		return nil, nil
	}
	return n.svc.RegisterDeviceToken(ctx, userID, token, platform)
}

func (n *notifSettingsServiceAdapter) DeleteDeviceToken(ctx context.Context, userID, token string) error {
	if n.svc == nil {
		return nil
	}
	return n.svc.DeleteDeviceToken(ctx, userID, token)
}

func (n *notifSettingsServiceAdapter) ListDeviceTokens(ctx context.Context, userID string) (any, error) {
	if n.svc == nil {
		return []*models.DeviceToken{}, nil
	}
	return n.svc.ListDeviceTokens(ctx, userID)
}

func (n *notifSettingsServiceAdapter) ResolveOrgID(ctx context.Context, userID string) string {
	if n.orgRepo == nil {
		return ""
	}
	orgs, err := n.orgRepo.ListByUser(ctx, userID)
	if err != nil || len(orgs) == 0 {
		return ""
	}
	return orgs[0].ID
}

type wsAccessCheckerAdapter struct {
	appService     *projectservices.AppService
	projectService *projectservices.ProjectService
}

func (w *wsAccessCheckerAdapter) HasAccess(ctx context.Context, serviceID, userID, role string) (bool, error) {
	if role == "admin" {
		return true, nil
	}
	if w.appService == nil || w.projectService == nil {
		return true, nil
	}
	svc, err := w.appService.GetAppService(ctx, serviceID)
	if err != nil || svc == nil {
		return false, err
	}
	return w.projectService.HasPermission(ctx, svc.ProjectID, userID, models.UserRole(role), ""), nil
}

type wsUserValidatorAdapter struct {
	userRepo repositories.UserRepository
}

func (w *wsUserValidatorAdapter) IsUserActive(ctx context.Context, userID string) (bool, error) {
	if w.userRepo == nil {
		return true, nil
	}
	u, err := w.userRepo.GetUserByID(ctx, userID)
	if err != nil || u == nil || !u.IsActive {
		return false, err
	}
	return true, nil
}
