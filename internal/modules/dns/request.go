package dns

type CreateDNSRecordRequest struct {
	DomainName  string `json:"domainName"`
	RecordType  string `json:"recordType"`
	RecordName  string `json:"recordName"`
	RecordValue string `json:"recordValue"`
	TTL         int    `json:"ttl,omitempty"`
}

type UpdateDNSRecordRequest struct {
	RecordType  string `json:"recordType"`
	RecordName  string `json:"recordName"`
	RecordValue string `json:"recordValue"`
	TTL         int    `json:"ttl,omitempty"`
}

