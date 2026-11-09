package support

type CloudSupportCustomerInput struct {
	RequestID string               `json:"requestId"`
	Subject   string               `json:"subject"`
	Message   string               `json:"message"`
	Category  CloudSupportCategory `json:"category"`
}

type CloudSupportCustomerReply struct {
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
}

type CloudSupportCustomerStatus struct {
	Status CloudSupportStatus `json:"status"`
}
