package types

import "time"

const (
	DNSProvisionStatusPending = "pending"
	DNSProvisionStatusSuccess = "success"
	DNSProvisionStatusFailed  = "failed"
	DNSProvisionStatusManual  = "manual"
)

const (
	RecordTypeA     = "A"
	RecordTypeCNAME = "CNAME"
	RecordTypeTXT   = "TXT"
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
	DNSVerificationCode string    `json:"dnsVerificationCode,omitempty" db:"dns_verification_code"`
	Verified            bool      `json:"verified" db:"verified"`
	CreatedAt           time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt" db:"updated_at"`
}

type DomainConfig = Domain
