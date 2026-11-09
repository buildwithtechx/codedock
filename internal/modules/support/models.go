package support

import "time"

type CloudSupportTicketRecord struct {
	ID        string
	AccountID string
	Subject   string
	Category  CloudSupportCategory
	Status    CloudSupportStatus
	Message   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CloudSupportMessageRecord struct {
	ID        string
	TicketID  string
	Author    CloudSupportAuthor
	Body      string
	CreatedAt time.Time
}

type CloudSupportCustomerTicket struct {
	ID        string               `json:"id"`
	Subject   string               `json:"subject"`
	Category  CloudSupportCategory `json:"category"`
	Status    CloudSupportStatus   `json:"status"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
}
