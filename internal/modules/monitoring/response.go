package monitoring

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

type IngressDiscoveryAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IngressDiscoveryTunnel struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Remote bool   `json:"remote"`
	Status string `json:"status"`
}

type IngressDiscoveryPolicy struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Decision string `json:"decision"`
}

type IngressDiscoveryResult struct {
	Accounts      []IngressDiscoveryAccount `json:"accounts"`
	Tunnels       []IngressDiscoveryTunnel  `json:"tunnels"`
	Policies      []IngressDiscoveryPolicy  `json:"policies"`
	PoliciesError *string                   `json:"policiesError,omitempty"`
}

type SystemInfoResponse struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}
