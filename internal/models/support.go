package models

import "time"

type CloudSupportCategory string

const (
	CloudSupportCategoryDeployment CloudSupportCategory = "deployment"
	CloudSupportCategoryBilling    CloudSupportCategory = "billing"
	CloudSupportCategoryAccount    CloudSupportCategory = "account"
	CloudSupportCategoryGeneral    CloudSupportCategory = "general"
)

type CloudSupportStatus string

const (
	CloudSupportStatusOpen     CloudSupportStatus = "open"
	CloudSupportStatusResolved CloudSupportStatus = "resolved"
)

type CloudSupportAuthor string

const (
	CloudSupportAuthorCustomer CloudSupportAuthor = "customer"
	CloudSupportAuthorSupport  CloudSupportAuthor = "support"
)

type CloudSupportAccountInfo struct {
	Key   string  `json:"key"`
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Email string  `json:"email"`
}

type CloudSupportSession struct {
	Account *CloudSupportAccountInfo `json:"account"`
}

type CloudSupportReceipt struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type CloudSupportCustomerInput struct {
	RequestID string               `json:"requestId"`
	Subject   string               `json:"subject"`
	Message   string               `json:"message"`
	Category  CloudSupportCategory `json:"category"`
}

type CloudSupportCustomerReply struct {
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}

type CloudSupportCustomerStatus struct {
	Status CloudSupportStatus `json:"status"`
}

type CloudSupportCustomerTicket struct {
	ID        string               `json:"id"`
	Subject   string               `json:"subject"`
	Category  CloudSupportCategory `json:"category"`
	Status    CloudSupportStatus   `json:"status"`
	CreatedAt string               `json:"createdAt"`
	UpdatedAt string               `json:"updatedAt"`
}

type CloudSupportCustomerTicketDetail struct {
	CloudSupportCustomerTicket
	Message string `json:"message"`
}

type CloudSupportCustomerMessageItem struct {
	ID        string             `json:"id"`
	Author    CloudSupportAuthor `json:"author"`
	Body      string             `json:"body"`
	CreatedAt string             `json:"createdAt"`
}

type CloudSupportCustomerDetail struct {
	Ticket   CloudSupportCustomerTicketDetail  `json:"ticket"`
	Messages []CloudSupportCustomerMessageItem `json:"messages"`
}

type CloudSupportCustomerList struct {
	Tickets    []CloudSupportCustomerTicket `json:"tickets"`
	NextCursor *string                      `json:"nextCursor"`
}

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
