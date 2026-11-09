package notifications

import (
	"fmt"
)

func SendDiscordNotification(webhookURL string, event *NotificationEvent) error {
	return postJSON(webhookURL, map[string]string{
		"content": fmt.Sprintf("**%s**\n%s\n%s", event.Title, event.Message, event.URL),
	})
}
