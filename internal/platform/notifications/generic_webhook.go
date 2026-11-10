package notifications

func SendGenericWebhook(webhookURL string, event *NotificationEvent) error {
	return postJSON(webhookURL, event)
}
