package billing

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

func (s *BillingService) GetConfig(ctx context.Context) (any, error) {
	cfg := config.Get()
	return &BillingConfig{
		PublishableKey: cfg.Stripe.PublishableKey,
		IsConfigured:   cfg.Stripe.PublishableKey != "",
		Plans: []BillingPlan{
			{
				ID:           "free",
				Name:         "Free Tier",
				Description:  "Self-hosted basic version",
				PriceMonthly: 0,
				Features:     []string{"1 Project", "Basic Support"},
			},
			{
				ID:           "pro",
				Name:         "Pro Tier",
				Description:  "Production servers with full managed features",
				PriceMonthly: 29,
				PriceID:      cfg.Stripe.PriceIDPro,
				Features:     []string{"Unlimited Projects", "Priority Support", "Analytics"},
			},
		},
	}, nil
}

func (s *BillingService) CreateCheckoutSession(
	ctx context.Context,
	userID string,
	opts CreateCheckoutSessionOpts,
) (string, error) {
	if opts.PurchaseIntent != "" && !IsPurchaseIntent(opts.PurchaseIntent) {
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
		if opts.CancelPending {
			if s.sessionExpire != nil && stripe.Key != "" {
				_, _ = s.sessionExpire(activeOpen.CheckoutID, nil)
			}
			activeOpen.Status = "canceled"
			activeOpen.CanCancel = false
		} else if opts.IdempotencyKey != "" && activeOpen.IdempotencyKey == opts.IdempotencyKey && activeOpen.PriceID == opts.PriceID {
			s.pendingMu.Unlock()
			return activeOpen.URL, nil
		} else {
			s.pendingMu.Unlock()
			return "", &ErrCheckoutPending{Pending: activeOpen}
		}
	}
	s.pendingMu.Unlock()

	if stripe.Key == "" && s.sessionNew == nil {
		return "", ErrBillingNotConfigured
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

	idempotencyKey := opts.IdempotencyKey
	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	}

	metadata := map[string]string{
		"user_id": userID,
	}
	if opts.PurchaseIntent != "" {
		metadata["purchase_intent"] = opts.PurchaseIntent
	}

	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customerID),
		Mode:     stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(opts.PriceID),
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(opts.SuccessURL),
		CancelURL:  stripe.String(opts.CancelURL),
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
		PurchaseIntent: opts.PurchaseIntent,
		PriceID:        opts.PriceID,
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
