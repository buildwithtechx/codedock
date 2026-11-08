package git

import "time"

type ConnectedProviderStatus struct {
	Provider    string     `json:"provider"`
	Connected   bool       `json:"connected"`
	AccountName string     `json:"accountName,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type PushDispatchResult struct {
	Duplicate      bool     `json:"duplicate,omitempty"`
	HandledCount   int      `json:"handledCount"`
	FailedCount    int      `json:"failedCount"`
	Errors         []string `json:"errors,omitempty"`
	HandledProjIDs []string `json:"handledProjectIds,omitempty"`
}

type DeliveriesResponse struct {
	Deliveries []WebhookDelivery `json:"deliveries"`
	NextCursor string            `json:"nextCursor,omitempty"`
}

type DisconnectResponse struct {
	Status string `json:"status"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
