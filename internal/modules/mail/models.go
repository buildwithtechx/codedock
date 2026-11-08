package mail

type MailServer struct {
	ID             string `json:"id" db:"id"`
	OrganizationID string `json:"organizationId" db:"organization_id"`
	ServerID       string `json:"serverId,omitempty" db:"server_id"`
	Domain         string `json:"domain" db:"domain"`
	MaildURL       string `json:"maildUrl" db:"maild_url"`
	WebmailURL     string `json:"webmailUrl,omitempty" db:"webmail_url"`
	APIToken       string `json:"-" db:"api_token"`
	CreatedAt      string `json:"createdAt" db:"created_at"`
	UpdatedAt      string `json:"updatedAt" db:"updated_at"`
}

type MailDomain struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Mailbox struct {
	ID        string `json:"id"`
	DomainID  string `json:"domainId"`
	LocalPart string `json:"localPart"`
	Address   string `json:"address"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type MailboxCredentials struct {
	Mailbox
	Password string `json:"password,omitempty"`
}

type DKIMKey struct {
	ID        string `json:"id"`
	DomainID  string `json:"domainId"`
	Selector  string `json:"selector"`
	Algorithm string `json:"algorithm"`
	KeyID     string `json:"keyId,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
	DNSRecord string `json:"dnsRecord,omitempty"`
	CreatedAt string `json:"createdAt"`
}

type MailDNSRecord struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

type DNSResponse struct {
	Domain      string          `json:"domain"`
	Records     []MailDNSRecord `json:"records"`
	EngineError string          `json:"engineError,omitempty"`
}

type MailCheck struct {
	Name   string `json:"name"`
	Passed bool   `json:"passed"`
	Detail string `json:"detail"`
}

type MailCheckReport struct {
	Domain string      `json:"domain"`
	Checks []MailCheck `json:"checks"`
}

type MailCert struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	ExpiresAt string `json:"expiresAt"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type Fail2banJail struct {
	Jail            string   `json:"jail"`
	CurrentlyFailed int      `json:"currentlyFailed"`
	TotalFailed     int      `json:"totalFailed"`
	CurrentlyBanned int      `json:"currentlyBanned"`
	TotalBanned     int      `json:"totalBanned"`
	BannedIPs       []string `json:"bannedIps"`
}

type Fail2banStatusResponse struct {
	Jails     []string                `json:"jails"`
	Details   map[string]Fail2banJail `json:"details"`
	CheckedAt string                  `json:"checkedAt"`
	CanManage bool                    `json:"canManage"`
}

type Fail2banActionResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output,omitempty"`
}

type MailAlias struct {
	ID          string `json:"id"`
	DomainID    string `json:"domainId"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Active      bool   `json:"active"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
