package webhooks

type WebhookResponse struct {
	Webhook *Webhook `json:"webhook"`
}

type WebhookListResponse struct {
	Webhooks []*Webhook `json:"webhooks"`
	Total    int        `json:"total"`
}

type LogDrainResponse struct {
	LogDrain *LogDrain `json:"logDrain"`
}

type LogDrainListResponse struct {
	LogDrains []*LogDrain `json:"logDrains"`
	Total     int         `json:"total"`
}
