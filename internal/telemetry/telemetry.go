package telemetry

import (
	"log"
	"os"
	"sync"

	"github.com/posthog/posthog-go"
)

var (
	client posthog.Client
	mu     sync.Mutex
)

func Init() {
	mu.Lock()
	defer mu.Unlock()

	apiKey := os.Getenv("POSTHOG_API_KEY")
	if apiKey == "" {
		return // telemetry disabled
	}

	host := os.Getenv("POSTHOG_HOST")
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
