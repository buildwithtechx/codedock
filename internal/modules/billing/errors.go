package billing

import "errors"

var (
	ErrBillingNotConfigured = errors.New("billing is not configured")
	ErrCheckoutNotFound     = errors.New("checkout session not found or already completed")
)

type ErrCheckoutPending struct {
	Pending *PendingCheckout
}

func (e *ErrCheckoutPending) Error() string {
	return "This server already has an unfinished checkout. Resume that payment or wait for it to expire before choosing another plan."
}

func (e *ErrCheckoutPending) PendingSession() any {
	return e.Pending
}

type CheckoutPendingError interface {
	error
	PendingSession() any
}
