package webhooks

type CreateWebhookRequest struct {
	URL    string `json:"url"`
	Secret string `json:"secret,omitempty"`
	Events string `json:"events,omitempty"`
}

type CreateLogDrainRequest struct {
	ProjectID   string    `json:"projectId"`
	DrainType   DrainType `json:"drainType"`
	EndpointURL string    `json:"endpointUrl"`
	APIKey      string    `json:"apiKey,omitempty"`
}
