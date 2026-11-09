package support

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
