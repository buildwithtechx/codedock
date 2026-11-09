package monitoring

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"

	"codedock/internal/config"
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

type BillingService struct {
	userRepo         UserRepository
	pendingMu        sync.Mutex
	pendingCheckouts map[string][]*PendingCheckout
	sessionNew       func(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error)
	sessionExpire    func(id string, params *stripe.CheckoutSessionExpireParams) (*stripe.CheckoutSession, error)
}

func NewBillingService(userRepo UserRepository) *BillingService {
	stripe.Key = config.Get().Stripe.SecretKey
	return &BillingService{
		userRepo:         userRepo,
		pendingCheckouts: make(map[string][]*PendingCheckout),
		sessionNew:       session.New,
		sessionExpire:    session.Expire,
	}
}

func (s *BillingService) SetSessionHandlers(
	newFn func(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error),
	expireFn func(id string, params *stripe.CheckoutSessionExpireParams) (*stripe.CheckoutSession, error),
) {
	s.sessionNew = newFn
	s.sessionExpire = expireFn
}

func (s *BillingService) CreateCheckoutSession(
	ctx context.Context,
	userID, priceID, successURL, cancelURL, idempotencyKey, purchaseIntent string,
	cancelPending bool,
) (string, error) {
	if purchaseIntent != "" && !IsPurchaseIntent(purchaseIntent) {
		return "", errors.New("invalid purchase intent format")
	}

	s.pendingMu.Lock()
	userPending := s.pendingCheckouts[userID]
	var activeOpen *PendingCheckout
	for _, item := range userPending {
		if item.Status == "open" {
			activeOpen = item
			break
		}
	}

	if activeOpen != nil {
		if cancelPending {
			if s.sessionExpire != nil && stripe.Key != "" {
				_, _ = s.sessionExpire(activeOpen.CheckoutID, nil)
			}
			activeOpen.Status = "canceled"
			activeOpen.CanCancel = false
		} else if idempotencyKey != "" && activeOpen.IdempotencyKey == idempotencyKey && activeOpen.PriceID == priceID {
			s.pendingMu.Unlock()
			return activeOpen.URL, nil
		} else {
			s.pendingMu.Unlock()
			return "", &ErrCheckoutPending{Pending: activeOpen}
		}
	}
	s.pendingMu.Unlock()

	if stripe.Key == "" && s.sessionNew == nil {
		return "", fmt.Errorf("billing is not configured")
	}

	var customerID string
	if s.userRepo != nil {
		u, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			return "", fmt.Errorf("failed to get user: %w", err)
		}
		customerID, err = s.getOrCreateStripeCustomer(ctx, u)
		if err != nil {
			return "", err
		}
	} else {
		customerID = "cus_mock_" + userID
	}

	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	}

	metadata := map[string]string{
		"user_id": userID,
	}
	if purchaseIntent != "" {
		metadata["purchase_intent"] = purchaseIntent
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		Metadata:   metadata,
	}
	params.IdempotencyKey = stripe.String(idempotencyKey)

	sess, err := s.sessionNew(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	newPending := &PendingCheckout{
		CheckoutID:     sess.ID,
		IdempotencyKey: idempotencyKey,
		PurchaseIntent: purchaseIntent,
		PriceID:        priceID,
		URL:            sess.URL,
		Status:         "open",
		CanCancel:      true,
		CreatedAt:      time.Now(),
	}

	s.pendingMu.Lock()
	s.pendingCheckouts[userID] = append(s.pendingCheckouts[userID], newPending)
	s.pendingMu.Unlock()

	return sess.URL, nil
}
