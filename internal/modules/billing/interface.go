package billing

import "context"

// UserRepository provides access to user subscription and Stripe billing attributes.
type UserRepository interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByStripeCustomerID(ctx context.Context, customerID string) (*User, error)
	UpdateUser(ctx context.Context, u *User) error
}

// BillingOperations defines the contract for checkout session orchestration and webhook ingest.
type BillingOperations interface {
	GetConfig(ctx context.Context) (any, error)
	CreateCheckoutSession(ctx context.Context, userID string, opts CreateCheckoutSessionOpts) (string, error)
	GetPendingCheckout(ctx context.Context, userID string) *PendingCheckout
	CancelCheckoutSession(ctx context.Context, userID, id, idempotencyKey string) (*PendingCheckout, error)
	HandleWebhook(ctx context.Context, payload []byte, sig string) error
}
