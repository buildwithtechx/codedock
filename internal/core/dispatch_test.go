package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"
	"codedock/internal/repositories"
)

type stubNotifSettingsRepo struct {
	repositories.NotificationSettingsRepository
	settings *models.NotificationSettings
}

func (s *stubNotifSettingsRepo) GetNotificationSettings(context.Context) (*models.NotificationSettings, error) {
	return s.settings, nil
}

type stubDeliveriesRepo struct {
	repositories.NotificationDeliveryRepository
	recorded []*models.NotificationDelivery
}

func (s *stubDeliveriesRepo) RecordDelivery(_ context.Context, delivery *models.NotificationDelivery) error {
	s.recorded = append(s.recorded, delivery)
	return nil
}

func TestDispatcherRecordsSuccessfulDelivery(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()
	deliveries := &stubDeliveriesRepo{}
	dispatcher := NewDispatcherService(nil, &stubNotifSettingsRepo{settings: &models.NotificationSettings{
		DiscordEnabled:    true,
		DiscordWebhookURL: server.URL,
	}}, nil, deliveries, nil)
	err := dispatcher.Send(&models.NotificationEvent{EventType: "test_global_discord", Title: "hello"})
	if err != nil {
		t.Fatalf("send test notification: %v", err)
	}
	if len(deliveries.recorded) != 1 {
		t.Fatalf("expected 1 recorded delivery, got %d", len(deliveries.recorded))
	}
	recorded := deliveries.recorded[0]
	if recorded.Channel != "discord" || !recorded.Success || recorded.Error != "" {
		t.Fatalf("unexpected delivery record: %+v", recorded)
	}
}

func TestDispatcherRecordsFailedDelivery(t *testing.T) {
	deliveries := &stubDeliveriesRepo{}
	dispatcher := NewDispatcherService(nil, &stubNotifSettingsRepo{settings: &models.NotificationSettings{}}, nil, deliveries, nil)
	err := dispatcher.Send(&models.NotificationEvent{EventType: "test_global_slack", Title: "hello"})
	if err == nil {
		t.Fatal("expected send to fail for unconfigured provider")
	}
	if len(deliveries.recorded) != 1 {
		t.Fatalf("expected 1 recorded delivery, got %d", len(deliveries.recorded))
	}
	recorded := deliveries.recorded[0]
	if recorded.Channel != "slack" || recorded.Success {
		t.Fatalf("unexpected delivery record: %+v", recorded)
	}
	if recorded.Error == "" {
		t.Fatal("expected error text on failed delivery")
	}
}
