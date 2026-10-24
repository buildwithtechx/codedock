package system

import (
	"context"
	"testing"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type stubDeliveryStore struct {
	repositories.NotificationDeliveryRepository
	rows []models.NotificationDelivery
}

func (s *stubDeliveryStore) ListDeliveries(context.Context, int) ([]models.NotificationDelivery, error) {
	return s.rows, nil
}

func TestListDeliveriesReturnsStoredRows(t *testing.T) {
	service := NewNotificationSettingsService(nil, &stubDeliveryStore{rows: []models.NotificationDelivery{{ID: "d-1", Channel: "discord", Success: true}}}, nil)
	rows, err := service.ListDeliveries(context.Background(), 10)
	if err != nil {
		t.Fatalf("list deliveries: %v", err)
	}
	if len(rows) != 1 || rows[0].ID != "d-1" {
		t.Fatalf("unexpected deliveries: %+v", rows)
	}
}

func TestListDeliveriesWithoutStoreReturnsEmpty(t *testing.T) {
	service := NewNotificationSettingsService(nil, nil, nil)
	rows, err := service.ListDeliveries(context.Background(), 10)
	if err != nil {
		t.Fatalf("list deliveries: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected empty deliveries, got %+v", rows)
	}
}
