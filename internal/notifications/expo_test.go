package notifications

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"codedock/internal/models"
)

func TestIsExpoPushToken(t *testing.T) {
	if !IsExpoPushToken("ExponentPushToken[xxxxxxxxxxxxxxxxxxxxxx]") {
		t.Fatal("expected valid ExponentPushToken")
	}
	if !IsExpoPushToken("ExpoPushToken[xxxxxxxxxxxxxxxxxxxxxx]") {
		t.Fatal("expected valid ExpoPushToken")
	}
	if IsExpoPushToken("invalid-token") {
		t.Fatal("expected invalid token to return false")
	}
}

func TestSendExpoNotificationSuccess(t *testing.T) {
	var received []ExpoPushMessage
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"status":"ok","id":"ticket-1"}]}`))
	}))
	defer server.Close()

	event := &models.NotificationEvent{
		Title:     "Deploy Succeeded",
		Message:   "Application web deployed successfully",
		Level:     "success",
		EventType: "deploy.succeeded",
		ProjectID: "proj-1",
		URL:       "https://dashboard.example.com/projects/proj-1",
	}

	err := SendExpoNotificationWithClient(server.Client(), server.URL, "ExponentPushToken[test-token]", event)
	if err != nil {
		t.Fatalf("unexpected error sending expo notification: %v", err)
	}

	if len(received) != 1 {
		t.Fatalf("expected 1 message sent, got %d", len(received))
	}
	msg := received[0]
	if msg.To != "ExponentPushToken[test-token]" || msg.Title != event.Title || msg.Body != event.Message {
		t.Fatalf("unexpected message content: %+v", msg)
	}
	if msg.Data["eventType"] != "deploy.succeeded" || msg.Data["projectId"] != "proj-1" {
		t.Fatalf("unexpected message data payload: %+v", msg.Data)
	}
}

func TestSendExpoNotificationErrorHandling(t *testing.T) {
	// Rejects invalid token format
	err := SendExpoNotificationWithClient(http.DefaultClient, "http://localhost", "bad-token", &models.NotificationEvent{})
	if err == nil {
		t.Fatal("expected error on invalid token format")
	}

	// Server error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"errors":[{"code":"PUSH_TOO_MANY_EXPERIENCE_IDS","message":"Multiple experiences"}]}`))
	}))
	defer server.Close()

	err = SendExpoNotificationWithClient(server.Client(), server.URL, "ExponentPushToken[test-token]", &models.NotificationEvent{Title: "hi"})
	if err == nil {
		t.Fatal("expected error on server error response")
	}
}
