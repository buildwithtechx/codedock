package system

import (
	"context"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type NotificationSettingsService struct {
	repo         repositories.NotificationSettingsRepository
	deliveries   repositories.NotificationDeliveryRepository
	deviceTokens repositories.DeviceTokenRepository
}

func NewNotificationSettingsService(repo repositories.NotificationSettingsRepository, deliveries repositories.NotificationDeliveryRepository, deviceTokens repositories.DeviceTokenRepository) *NotificationSettingsService {
	return &NotificationSettingsService{repo: repo, deliveries: deliveries, deviceTokens: deviceTokens}
}

func (s *NotificationSettingsService) GetNotificationSettings(ctx context.Context) (*models.NotificationSettings, error) {
	return s.repo.GetNotificationSettings(ctx)
}

func (s *NotificationSettingsService) UpdateNotificationSettings(ctx context.Context, cfg *models.NotificationSettings) error {
	return s.repo.UpdateNotificationSettings(ctx, cfg)
}

func (s *NotificationSettingsService) ListSubscriptions(ctx context.Context, userID, orgID string) ([]models.NotificationSubscription, error) {
	return s.repo.ListSubscriptions(ctx, userID, orgID)
}

func (s *NotificationSettingsService) UpsertSubscription(ctx context.Context, sub *models.NotificationSubscription) error {
	return s.repo.UpsertSubscription(ctx, sub)
}

func (s *NotificationSettingsService) ListDefaults(ctx context.Context, orgID string) ([]models.NotificationDefault, error) {
	return s.repo.ListDefaults(ctx, orgID)
}

func (s *NotificationSettingsService) UpsertDefault(ctx context.Context, def *models.NotificationDefault) error {
	return s.repo.UpsertDefault(ctx, def)
}

func (s *NotificationSettingsService) ListDeliveries(ctx context.Context, limit int) ([]models.NotificationDelivery, error) {
	if s.deliveries == nil {
		return []models.NotificationDelivery{}, nil
	}
	return s.deliveries.ListDeliveries(ctx, limit)
}

func (s *NotificationSettingsService) RegisterDeviceToken(ctx context.Context, userID, token, platform string) (*models.DeviceToken, error) {
	if s.deviceTokens == nil {
		return nil, nil
	}
	return s.deviceTokens.RegisterToken(ctx, userID, token, platform)
}

func (s *NotificationSettingsService) DeleteDeviceToken(ctx context.Context, userID, token string) error {
	if s.deviceTokens == nil {
		return nil
	}
	return s.deviceTokens.DeleteToken(ctx, userID, token)
}

func (s *NotificationSettingsService) ListDeviceTokens(ctx context.Context, userID string) ([]*models.DeviceToken, error) {
	if s.deviceTokens == nil {
		return []*models.DeviceToken{}, nil
	}
	return s.deviceTokens.ListTokensByUser(ctx, userID)
}
