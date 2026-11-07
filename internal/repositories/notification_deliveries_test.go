package repositories

import (
	"context"
	"testing"

	"codedock/internal/models"
)

func TestNotificationDeliveryRoundtrip(t *testing.T) {
	db := openPGTestDB(t)
	ctx := context.Background()
	repo := NewNotificationDeliveryRepo(db)
	first := &models.NotificationDelivery{EventType: "test_global_discord", Channel: "discord", Title: "hello", Success: true}
	if err := repo.RecordDelivery(ctx, first); err != nil {
		t.Fatalf("record delivery: %v", err)
	}
	if first.ID == "" {
		t.Fatal("expected generated delivery id")
	}
	failed := &models.NotificationDelivery{EventType: "test_global_slack", Channel: "slack", Title: "hello", Success: false, Error: "connection refused"}
	if err := repo.RecordDelivery(ctx, failed); err != nil {
		t.Fatalf("record failed delivery: %v", err)
	}
	list, err := repo.ListDeliveries(ctx, 10)
	if err != nil {
		t.Fatalf("list deliveries: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("expected 2 deliveries, got %d", len(list))
	}
	if list[0].ID != failed.ID || !list[1].Success {
		t.Fatalf("expected newest-first order with success flags, got %+v", list)
	}
	if list[0].Error != "connection refused" {
		t.Fatalf("expected error to roundtrip, got %+v", list[0])
	}
	one, err := repo.ListDeliveries(ctx, 1)
	if err != nil {
		t.Fatalf("list with limit: %v", err)
	}
	if len(one) != 1 || one[0].ID != failed.ID {
		t.Fatalf("expected newest delivery only, got %+v", one)
	}
	empty, err := repo.ListDeliveries(ctx, 0)
	if err != nil {
		t.Fatalf("list with default limit: %v", err)
	}
	if len(empty) != 2 {
		t.Fatalf("expected default limit to return all, got %d", len(empty))
	}
}
