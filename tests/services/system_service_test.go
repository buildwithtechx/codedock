package services_test

import (
	"testing"

	"codedock.run/codedock/internal/services/system"
)

func TestSystemService_GetStats(t *testing.T) {
	svc := system.NewSystemService()

	stats, err := svc.GetStats()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if stats == nil {
		t.Fatal("expected stats to be non-nil")
	}

	_, _ = svc.GetStats()
}
