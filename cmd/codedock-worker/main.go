package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codedock.run/codedock/internal/models"
	"github.com/gorilla/websocket"
)

func main() {
	log.Println("Starting codedock-worker...")

	token := os.Getenv("CODEDOCK_WORKER_TOKEN")
	if token == "" {
		log.Fatal("CODEDOCK_WORKER_TOKEN is required")
	}

	apiHost := os.Getenv("CODEDOCK_API_HOST")
	if apiHost == "" {
		apiHost = "localhost:8080"
	}

	u := url.URL{Scheme: "ws", Host: apiHost, Path: "/ws/worker", RawQuery: "token=" + token}
	if os.Getenv("CODEDOCK_USE_WSS") == "true" {
		u.Scheme = "wss"
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	for {
		log.Printf("Connecting to %s", u.String())
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("Dial error: %v", err)
			log.Println("Retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		done := make(chan struct{})
		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Printf("Read error: %v", err)
					return
				}
				log.Printf("Received: %s", message)
			}
		}()

		// Start heartbeat/metrics loop
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					// Send heartbeat or metrics
					metricsPayload := models.WorkerMetricsPayload{
						CPUUsagePercentage: 5.0,                     // Placeholder
						MemoryUsageBytes:   1024 * 1024 * 100,       // Placeholder
						DiskUsageBytes:     1024 * 1024 * 1024 * 10, // Placeholder
						MemoryLimitBytes:   1024 * 1024 * 1024,
						DiskTotalBytes:     1024 * 1024 * 1024 * 100,
					}
					payloadBytes, _ := json.Marshal(metricsPayload)

					msg := models.WorkerMessage{
						ID:        "heartbeat",
						Type:      models.WorkerMessageTypeMetrics,
						Timestamp: time.Now(),
						Payload:   payloadBytes,
					}

					err := c.WriteJSON(msg)
					if err != nil {
						log.Printf("Heartbeat error: %v", err)
						return
					}
				}
			}
		}()

		select {
		case <-done:
			log.Println("Connection closed, reconnecting...")
			time.Sleep(5 * time.Second)
		case <-interrupt:
			log.Println("Interrupt received, closing connection")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Printf("Write close error: %v", err)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
