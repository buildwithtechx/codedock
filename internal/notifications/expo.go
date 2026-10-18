package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"codedock/internal/models"
)

var expoPushAPIURL = "https://exp.host/--/api/v2/push/send"

type ExpoPushMessage struct {
	To        string         `json:"to"`
	Title     string         `json:"title,omitempty"`
	Body      string         `json:"body"`
	Data      map[string]any `json:"data,omitempty"`
	Sound     string         `json:"sound,omitempty"`
	Priority  string         `json:"priority,omitempty"`
	ChannelID string         `json:"channelId,omitempty"`
}

type expoPushResponseItem struct {
	Status  string         `json:"status"`
	Message string         `json:"message,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

type expoPushResponseBody struct {
	Data   []expoPushResponseItem `json:"data,omitempty"`
	Errors []struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	} `json:"errors,omitempty"`
}

func IsExpoPushToken(token string) bool {
	return strings.HasPrefix(token, "ExponentPushToken[") || strings.HasPrefix(token, "ExpoPushToken[")
}

func SendExpoNotification(token string, event *models.NotificationEvent) error {
	client := &http.Client{Timeout: 10 * time.Second}
	return SendExpoNotificationWithAuth(client, expoPushAPIURL, token, "", event)
}

func SendExpoNotificationWithClient(httpClient *http.Client, apiURL, token string, event *models.NotificationEvent) error {
	return SendExpoNotificationWithAuth(httpClient, apiURL, token, "", event)
}

func SendExpoNotificationWithAuth(httpClient *http.Client, apiURL, token, accessToken string, event *models.NotificationEvent) error {
	if !IsExpoPushToken(token) {
		return fmt.Errorf("invalid expo push token: %s", token)
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	if apiURL == "" {
		apiURL = expoPushAPIURL
	}

	msg := ExpoPushMessage{
		To:       token,
		Title:    event.Title,
		Body:     event.Message,
		Sound:    "default",
		Priority: "high",
		Data: map[string]any{
			"url":       event.URL,
			"eventType": event.EventType,
			"level":     event.Level,
			"projectId": event.ProjectID,
		},
	}

	payloadBytes, err := json.Marshal([]ExpoPushMessage{msg})
	if err != nil {
		return fmt.Errorf("marshal expo push message: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("create expo push request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send expo push request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return fmt.Errorf("expo push returned HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var parsed expoPushResponseBody
	if err := json.Unmarshal(bodyBytes, &parsed); err == nil {
		if len(parsed.Errors) > 0 {
			return fmt.Errorf("expo push error: %s", parsed.Errors[0].Message)
		}
		if len(parsed.Data) > 0 && parsed.Data[0].Status == "error" {
			return fmt.Errorf("expo ticket error: %s", parsed.Data[0].Message)
		}
	}

	return nil
}
