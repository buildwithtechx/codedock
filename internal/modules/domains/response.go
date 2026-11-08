package domains

type DomainResponse struct {
	Domain *DomainConfig `json:"domain"`
}

type DomainListResponse struct {
	Domains []DomainConfig `json:"domains"`
}
