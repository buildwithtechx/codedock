package deployments

import (
	"context"
	"time"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type WebhookDeliveryService struct {
	repo repositories.WebhookDeliveryRepository
}

func NewWebhookDeliveryService(repo repositories.WebhookDeliveryRepository) *WebhookDeliveryService {
	return &WebhookDeliveryService{repo: repo}
}

func (s *WebhookDeliveryService) ClaimGithub(
	ctx context.Context,
	deliveryID string,
	event string,
	clientIP *string,
	userAgent *string,
) (models.ClaimResult, error) {
	input := models.WebhookDeliveryInput{
		DeliveryID: &deliveryID,
		Event:      event,
		Source:     "github",
		ClientIP:   clientIP,
		UserAgent:  userAgent,
		Outcome:    "received",
	}
	return s.repo.ClaimGithub(ctx, input)
}

func (s *WebhookDeliveryService) Record(ctx context.Context, input models.WebhookDeliveryInput) (string, error) {
	return s.repo.Record(ctx, input)
}

func (s *WebhookDeliveryService) MarkProcessed(ctx context.Context, id string, res models.WebhookDeliveryResult) error {
	return s.repo.MarkProcessed(ctx, id, res)
}

func (s *WebhookDeliveryService) ListByProject(
	ctx context.Context,
	projectID string,
	limit int,
	cursor string,
) ([]models.WebhookDelivery, string, error) {
	return s.repo.ListByProject(ctx, projectID, limit, cursor)
}

func (s *WebhookDeliveryService) Prune(ctx context.Context, olderThan time.Duration) (int64, error) {
	return s.repo.Prune(ctx, olderThan)
}
