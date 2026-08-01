package models

import (
	"time"

	"codedock.run/codedock/pkg/types"
)

type Domain = types.Domain
type DomainConfig = types.DomainConfig
type CreateDomainRequest = types.CreateDomainRequest
type DomainVerifyResult = types.DomainVerifyResult

const (
	DNSProvisionStatusPending = types.DNSProvisionStatusPending
	DNSProvisionStatusSuccess = types.DNSProvisionStatusSuccess
	DNSProvisionStatusFailed  = types.DNSProvisionStatusFailed
)

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
