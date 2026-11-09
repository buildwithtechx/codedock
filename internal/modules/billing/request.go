package billing

type CreateCheckoutRequest struct {
	PriceID        string `json:"priceId"`
	SuccessURL     string `json:"successUrl"`
	CancelURL      string `json:"cancelUrl"`
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
	PurchaseIntent string `json:"purchaseIntent,omitempty"`
	CancelPending  bool   `json:"cancelPending,omitempty"`
}

type CancelCheckoutRequest struct {
	IdempotencyKey string `json:"idempotencyKey,omitempty"`
}
