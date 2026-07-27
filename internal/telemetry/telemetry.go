package telemetry

import (
	"log"
	"sync"

	"github.com/posthog/posthog-go"

	"codedock.run/codedock/internal/config"
)

var (
	client posthog.Client
	mu     sync.Mutex
)

func Init() {
	mu.Lock()
	defer mu.Unlock()

	cfg := config.Get()
	apiKey := cfg.Telemetry.PosthogKey
	if apiKey == "" {
		return
	}

	host := cfg.Telemetry.PosthogHost
	if host == "" {
		host = "https://us.i.posthog.com"
	}

	c, err := posthog.NewWithConfig(apiKey, posthog.Config{
		Endpoint: host,
	})
	if err != nil {
		log.Printf("failed to initialize telemetry: %v", err)
		return
	}

	client = c
}

func Track(distinctID string, event string, properties map[string]any) {
	if client == nil {
		return
	}

	err := client.Enqueue(posthog.Capture{
		DistinctId: distinctID,
		Event:      event,
		Properties: properties,
	})
	if err != nil {
		log.Printf("failed to track event %s: %v", event, err)
	}
}

func Close() {
	mu.Lock()
	defer mu.Unlock()
	if client != nil {
		if err := client.Close(); err != nil {
			log.Printf("failed to close telemetry client: %v", err)
		}
	}
}
