package mail

// Response envelopes and view models for mail domain
type MailServerListResponse struct {
	Servers []*MailServer `json:"servers"`
}

type MailDomainListResponse struct {
	Domains []MailDomain `json:"domains"`
}

type MailboxListResponse struct {
	Mailboxes []Mailbox `json:"mailboxes"`
}

type MailAliasListResponse struct {
	Aliases []MailAlias `json:"aliases"`
}
