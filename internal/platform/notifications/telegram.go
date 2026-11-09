package notifications

import (
	"fmt"
)

func SendTelegramNotification(botToken, chatID string, event *NotificationEvent) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	return postJSON(url, map[string]string{
		"chat_id":    chatID,
		"text":       fmt.Sprintf("*%s*\n%s\n%s", event.Title, event.Message, event.URL),
		"parse_mode": "Markdown",
	})
}
