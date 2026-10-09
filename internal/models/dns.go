package models

import (
	"time"
)

const (
	DNSProvisionStatusPending = "pending"
	DNSProvisionStatusSuccess = "provisioned"
	DNSProvisionStatusFailed  = "failed"
)

type Domain struct {
	ID                  string    `json:"id" db:"id"`
	ServiceID           string    `json:"serviceId" db:"service_id"`
	DomainName          string    `json:"domainName" db:"domain_name"`
	RedirectTo          string    `json:"redirectTo,omitempty" db:"redirect_to"`
	PathPrefix          string    `json:"pathPrefix,omitempty" db:"path_prefix"`
	IsCustom            bool      `json:"isCustom" db:"is_custom"`
	SSLStatus           string    `json:"sslStatus" db:"ssl_status"`
	SSLCertStatus       string    `json:"sslCertStatus,omitempty" db:"ssl_cert_status"`
	DNSProvisionStatus  string    `json:"dnsProvisionStatus,omitempty" db:"dns_provision_status"`
	DNSProvider         string    `json:"dnsProvider,omitempty" db:"dns_provider"`
	DNSProvisionedIP    string    `json:"dnsProvisionedIp,omitempty" db:"dns_provisioned_ip"`
	DNSVerificationCode string    `json:"dnsVerificationCode,omitempty" db:"dns_verification_code"`
	Verified            bool      `json:"verified" db:"verified"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
}

type DomainConfig = Domain

type CreateDomainRequest struct {
	DomainName string `json:"domainName"`
	RedirectTo string `json:"redirectTo,omitempty"`
	PathPrefix string `json:"pathPrefix,omitempty"`
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
