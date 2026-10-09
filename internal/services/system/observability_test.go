package system

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestHistoricalLogsRequireObservability(t *testing.T) {
	service := &LogService{enabled: false}
	_, err := service.GetHistoricalLogs(context.Background(), HistoricalLogsOpts{
		ServiceID: "svc",
		Start:     time.Now().Add(-time.Hour),
		End:       time.Now(),
		Limit:     10,
	})
	if !errors.Is(err, ErrObservabilityDisabled) {
		t.Fatalf("expected observability disabled error, got %v", err)
	}
}

func TestServiceMetricsRequireObservability(t *testing.T) {
	service := &MetricsService{enabled: false}
	_, err := service.GetServiceMetrics(context.Background(), ServiceMetricsOpts{
		ServiceID: "svc",
		Start:     time.Now().Add(-time.Hour),
		End:       time.Now(),
		Step:      "5m",
	})
	if !errors.Is(err, ErrObservabilityDisabled) {
		t.Fatalf("expected observability disabled error, got %v", err)
	}
}
