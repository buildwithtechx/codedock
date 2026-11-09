package billing

import (
	"context"
	"regexp"
	"time"
)

type BillingConfig struct {
	PublishableKey string        `json:"publishableKey"`
	Plans          []BillingPlan `json:"plans"`
	IsConfigured   bool          `json:"isConfigured"`
	Enabled        bool          `json:"enabled,omitempty"`
	Mode           string        `json:"mode,omitempty"`
}

type BillingPlan struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	PriceMonthly int      `json:"priceMonthly"`
	Features     []string `json:"features"`
	PriceID      string   `json:"priceId,omitempty"`
}

type PendingCheckout struct {
	CheckoutID              string                   `json:"checkoutId"`
	URL                     string                   `json:"url"`
	PriceID                 string                   `json:"priceId"`
	PlanName                string                   `json:"planName,omitempty"`
	Status                  string                   `json:"status"`
	IdempotencyKey          string                   `json:"idempotencyKey,omitempty"`
	PurchaseIntent          string                   `json:"purchaseIntent,omitempty"`
	CreatedAt               time.Time                `json:"createdAt"`
	ExpiresAt               time.Time                `json:"expiresAt,omitempty"`
	CanCancel               bool                     `json:"canCancel"`
	CancelUnavailableReason *CancelUnavailableReason `json:"cancelUnavailableReason,omitempty"`
	PaymentMethodTypes      []string                 `json:"paymentMethodTypes,omitempty"`
}

var purchaseIntentRegex = regexp.MustCompile(`(?i)^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`)

func IsPurchaseIntent(value string) bool {
	return purchaseIntentRegex.MatchString(value)
}

type CreateCheckoutSessionOpts struct {
	PriceID        string
	SuccessURL     string
	CancelURL      string
	IdempotencyKey string
	PurchaseIntent string
	CancelPending  bool
}

type User struct {
	ID                   string  `json:"id"`
	Email                string  `json:"email"`
	Name                 string  `json:"name"`
	Role                 string  `json:"role"`
	PlanType             string  `json:"planType"`
	StripeCustomerID     *string `json:"stripeCustomerId"`
	StripeSubscriptionID *string `json:"stripeSubscriptionId"`
	StripePriceID        *string `json:"stripePriceId"`
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByStripeCustomerID(ctx context.Context, customerID string) (*User, error)
	UpdateUser(ctx context.Context, u *User) error
}

type BillingOperations interface {
	GetConfig(ctx context.Context) (any, error)
	CreateCheckoutSession(ctx context.Context, userID string, opts CreateCheckoutSessionOpts) (string, error)
	GetPendingCheckout(ctx context.Context, userID string) *PendingCheckout
	CancelCheckoutSession(ctx context.Context, userID, id, idempotencyKey string) (*PendingCheckout, error)
	HandleWebhook(ctx context.Context, payload []byte, sig string) error
}
