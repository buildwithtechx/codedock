package system

import (
	"context"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type NotificationSettingsService struct {
	repo repositories.NotificationSettingsRepository
}

func NewNotificationSettingsService(repo repositories.NotificationSettingsRepository) *NotificationSettingsService {
	return &NotificationSettingsService{repo: repo}
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
