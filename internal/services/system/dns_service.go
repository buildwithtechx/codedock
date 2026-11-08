package system

import "codedock/internal/modules/dns"

type DNSService = dns.Service
type DNSProviderService = dns.ProviderService

var (
	NewDNSService         = dns.NewService
	NewDNSProviderService = dns.NewProviderService
	VerifyDomain          = dns.VerifyDomain
)
