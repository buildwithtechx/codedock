package billing

type CreateCheckoutResponse struct {
	URL string `json:"url"`
}

type PendingCheckoutResponse struct {
	Pending *PendingCheckout `json:"pending"`
}
