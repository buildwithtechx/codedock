package chatapps

// DiscordInteractionResponse represents Discord's expected interaction callback shape
type DiscordInteractionResponse struct {
	Type int                        `json:"type"`
	Data *DiscordInteractionDataOut `json:"data,omitempty"`
}

type DiscordInteractionDataOut struct {
	Content string `json:"content"`
}

// IngressIgnoredResponse represents the response when an incoming webhook is safely ignored
type IngressIgnoredResponse struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

// SimpleTextResponse represents a standard plain-text chat reply
type SimpleTextResponse struct {
	Text string `json:"text"`
}
