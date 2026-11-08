package dns

import "time"

type DNSRecord struct {
	ID          string    `json:"id" db:"id"`
	DomainName  string    `json:"domainName" db:"domain_name"`
	RecordType  string    `json:"recordType" db:"record_type"`
	RecordName  string    `json:"recordName" db:"record_name"`
	RecordValue string    `json:"recordValue" db:"record_value"`
	TTL         int       `json:"ttl" db:"ttl"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type DomainVerifyResult struct {
	DomainID   string `json:"domainId"`
	DomainName string `json:"domainName"`
	Verified   bool   `json:"verified"`
	Status     string `json:"status"`
	ResolvedIP string `json:"resolvedIp"`
	ServerIP   string `json:"serverIp"`
	Message    string `json:"message"`
}
