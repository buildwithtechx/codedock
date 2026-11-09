package webhooks

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type WebhookService interface {
	CreateWebhook(ctx context.Context, w *Webhook) (*Webhook, error)
	ListWebhooks(ctx context.Context, serviceID string) ([]*Webhook, error)
	DeleteWebhook(ctx context.Context, id, serviceID string) error
	CreateLogDrain(ctx context.Context, d *LogDrain) (*LogDrain, error)
	ListLogDrains(ctx context.Context, serviceID string) ([]*LogDrain, error)
	DeleteLogDrain(ctx context.Context, id, serviceID string) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateWebhook(ctx context.Context, w *Webhook) (*Webhook, error) {
	if w == nil {
		return nil, errors.New("webhook is nil")
	}
	if w.ServiceID == "" {
		return nil, errors.New("serviceId is required")
	}
	if w.URL == "" {
		return nil, errors.New("URL is required")
	}
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	if err := s.repo.CreateWebhook(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) ListWebhooks(ctx context.Context, serviceID string) ([]*Webhook, error) {
	if serviceID == "" {
		return nil, errors.New("serviceId is required")
	}
	return s.repo.ListWebhooksByService(ctx, serviceID)
}

func (s *Service) DeleteWebhook(ctx context.Context, id, serviceID string) error {
	if id == "" || serviceID == "" {
		return errors.New("id and serviceId are required")
	}
	return s.repo.DeleteWebhook(ctx, id, serviceID)
}

func (s *Service) CreateLogDrain(ctx context.Context, d *LogDrain) (*LogDrain, error) {
	if d == nil {
		return nil, errors.New("log drain is nil")
	}
	if d.ServiceID == "" || d.ProjectID == "" {
		return nil, errors.New("serviceId and projectId are required")
	}
	if d.DrainType == "" || d.EndpointURL == "" {
		return nil, errors.New("drainType and endpointUrl are required")
	}
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	d.CreatedAt = now
	d.UpdatedAt = now
	if err := s.repo.CreateLogDrain(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) ListLogDrains(ctx context.Context, serviceID string) ([]*LogDrain, error) {
	if serviceID == "" {
		return nil, errors.New("serviceId is required")
	}
	return s.repo.ListLogDrainsByService(ctx, serviceID)
}

func (s *Service) DeleteLogDrain(ctx context.Context, id, serviceID string) error {
	if id == "" || serviceID == "" {
		return errors.New("id and serviceId are required")
	}
	return s.repo.DeleteLogDrain(ctx, id, serviceID)
}
