package networks

// NetworkResponse wraps a server network payload.
type NetworkResponse struct {
	Network *ServerNetwork `json:"network"`
}

// NetworkListResponse wraps a slice of server networks.
type NetworkListResponse struct {
	Networks []*ServerNetwork `json:"networks"`
	Total    int              `json:"total"`
}
