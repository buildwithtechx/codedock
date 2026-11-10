package ingress

// IngressConnectionResponse wraps an ingress connection payload.
type IngressConnectionResponse struct {
	Connection *IngressConnection `json:"connection"`
}

// IngressConnectionListResponse wraps a list of ingress connections.
type IngressConnectionListResponse struct {
	Connections []*IngressConnection `json:"connections"`
	Total       int                  `json:"total"`
}

// IngressHostnameResponse wraps an ingress hostname mapping.
type IngressHostnameResponse struct {
	Hostname *IngressHostname `json:"hostname"`
}
