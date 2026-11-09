package notifications

import ()

func SendGenericWebhook(webhookURL string, event *NotificationEvent) error {
	return postJSON(webhookURL, event)
}
