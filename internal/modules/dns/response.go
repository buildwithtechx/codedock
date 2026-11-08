package dns

type DNSRecordsListResponse struct {
	Records []*DNSRecord `json:"records"`
}

