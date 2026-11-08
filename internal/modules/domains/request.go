package domains

type CreateDomainRequest struct {
	DomainName string `json:"domainName"`
	RedirectTo string `json:"redirectTo,omitempty"`
	PathPrefix string `json:"pathPrefix,omitempty"`
}

type CreateDNSRecordRequest struct {
	DomainName  string `json:"domainName"`
	RecordType  string `json:"recordType"`
	RecordName  string `json:"recordName"`
	RecordValue string `json:"recordValue"`
	TTL         int    `json:"ttl,omitempty"`
}
