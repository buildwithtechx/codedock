package models

import (
	"regexp"
	"time"
)

type CancelUnavailableReason string

const (
	CancelUnavailableUnsupported         CancelUnavailableReason = "unsupported"
	CancelUnavailableProviderUnavailable CancelUnavailableReason = "provider_unavailable"
)

type PendingCheckout struct {
	CheckoutID              string                   `json:"checkoutId"`
	IdempotencyKey          string                   `json:"idempotencyKey"`
	PurchaseIntent          string                   `json:"purchaseIntent,omitempty"`
	PriceID                 string                   `json:"priceId"`
	URL                     string                   `json:"url"`
	Status                  string                   `json:"status"`
	CanCancel               bool                     `json:"canCancel"`
	CancelUnavailableReason *CancelUnavailableReason `json:"cancelUnavailableReason,omitempty"`
	CreatedAt               time.Time                `json:"createdAt"`
}

var purchaseIntentRegex = regexp.MustCompile(`(?i)^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)

func IsPurchaseIntent(value string) bool {
	return purchaseIntentRegex.MatchString(value)
}
