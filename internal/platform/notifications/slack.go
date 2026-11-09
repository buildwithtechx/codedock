package notifications

import (
	"fmt"
)

func SendSlackNotification(webhookURL string, event *NotificationEvent) error {
	return postJSON(webhookURL, map[string]string{
		"text": fmt.Sprintf("*%s*\n%s\n%s", event.Title, event.Message, event.URL),
	})
}
