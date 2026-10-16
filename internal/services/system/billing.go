package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/webhook"

	"codedock/internal/config"
	"codedock/internal/models"
	"codedock/internal/repositories"
)

type ErrCheckoutPending struct {
	Pending *models.PendingCheckout
}

func (e *ErrCheckoutPending) Error() string {
	return "This server already has an unfinished checkout. Resume that payment or wait for it to expire before choosing another plan."
}

type BillingService struct {
	userRepo         repositories.UserRepository
	pendingMu        sync.Mutex
	pendingCheckouts map[string][]*models.PendingCheckout
	sessionNew       func(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error)
	sessionExpire    func(id string, params *stripe.CheckoutSessionExpireParams) (*stripe.CheckoutSession, error)
}

func NewBillingService(userRepo repositories.UserRepository) *BillingService {
	stripe.Key = config.Get().Stripe.SecretKey
	return &BillingService{
		userRepo:         userRepo,
		pendingCheckouts: make(map[string][]*models.PendingCheckout),
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
	if purchaseIntent != "" && !models.IsPurchaseIntent(purchaseIntent) {
		return "", errors.New("invalid purchase intent format")
	}

	s.pendingMu.Lock()
	userPending := s.pendingCheckouts[userID]
	var activeOpen *models.PendingCheckout
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

	newPending := &models.PendingCheckout{
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

func (s *BillingService) GetPendingCheckout(ctx context.Context, userID string) *models.PendingCheckout {
	s.pendingMu.Lock()
	defer s.pendingMu.Unlock()

	userPending := s.pendingCheckouts[userID]
	for i := len(userPending) - 1; i >= 0; i-- {
		if userPending[i].Status == "open" {
			return userPending[i]
		}
	}
	return nil
}

func (s *BillingService) CancelCheckoutSession(
	ctx context.Context,
	userID, checkoutID, idempotencyKey string,
) (*models.PendingCheckout, error) {
	s.pendingMu.Lock()
	defer s.pendingMu.Unlock()

	userPending := s.pendingCheckouts[userID]
	var target *models.PendingCheckout
	for _, item := range userPending {
		if item.CheckoutID == checkoutID {
			target = item
			break
		}
	}

	if target == nil {
		return nil, fmt.Errorf("checkout not found or already settled")
	}

	if target.Status != "open" {
		return target, nil
	}

	if s.sessionExpire != nil && stripe.Key != "" {
		var expireParams *stripe.CheckoutSessionExpireParams
		if idempotencyKey != "" {
			expireParams = &stripe.CheckoutSessionExpireParams{
				Params: stripe.Params{
					IdempotencyKey: stripe.String(idempotencyKey),
				},
			}
		}
		if _, err := s.sessionExpire(checkoutID, expireParams); err != nil {
			reason := models.CancelUnavailableProviderUnavailable
			target.CancelUnavailableReason = &reason
			return nil, fmt.Errorf("provider cancel failed: %w", err)
		}
	}

	target.Status = "canceled"
	target.CanCancel = false
	target.CancelUnavailableReason = nil
	return target, nil
}

func (s *BillingService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	webhookSecret := config.Get().Stripe.WebhookSecret
	if webhookSecret == "" {
		return fmt.Errorf("stripe webhook secret not configured")
	}

	event, err := webhook.ConstructEvent(payload, signature, webhookSecret)
	if err != nil {
		return fmt.Errorf("invalid webhook signature: %w", err)
	}

	switch event.Type {
	case "customer.subscription.created", "customer.subscription.updated":
		return s.handleSubscriptionUpdated(ctx, event)
	case "customer.subscription.deleted":
		return s.handleSubscriptionDeleted(ctx, event)
	}

	return nil
}

func (s *BillingService) getOrCreateStripeCustomer(ctx context.Context, u *models.User) (string, error) {
	if u.StripeCustomerID != nil && *u.StripeCustomerID != "" {
		return *u.StripeCustomerID, nil
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(u.Email),
		Name:  stripe.String(u.Name),
		Metadata: map[string]string{
			"user_id": u.ID,
		},
	}
	c, err := customer.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create stripe customer: %w", err)
	}
	customerID := c.ID
	u.StripeCustomerID = &customerID
	if err := s.userRepo.UpdateUser(ctx, u); err != nil {
		return "", fmt.Errorf("failed to save stripe customer id: %w", err)
	}
	return customerID, nil
}

func (s *BillingService) handleSubscriptionUpdated(ctx context.Context, event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	u, err := s.userRepo.GetUserByStripeCustomerID(ctx, sub.Customer.ID)
	if err != nil {
		return fmt.Errorf("user not found for customer: %w", err)
	}

	if len(sub.Items.Data) > 0 {
		priceID := sub.Items.Data[0].Price.ID
		u.StripeSubscriptionID = &sub.ID
		u.StripePriceID = &priceID

		switch sub.Status {
		case stripe.SubscriptionStatusActive, stripe.SubscriptionStatusTrialing:
			u.PlanType = "pro"
		case stripe.SubscriptionStatusCanceled, stripe.SubscriptionStatusUnpaid, stripe.SubscriptionStatusIncompleteExpired:
			u.PlanType = "free"
		}

		if err := s.userRepo.UpdateUser(ctx, u); err != nil {
			return fmt.Errorf("failed to update user plan: %w", err)
		}
	}

	return nil
}

func (s *BillingService) handleSubscriptionDeleted(ctx context.Context, event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	u, err := s.userRepo.GetUserByStripeCustomerID(ctx, sub.Customer.ID)
	if err == nil {
		u.PlanType = "free"
		if err := s.userRepo.UpdateUser(ctx, u); err != nil {
			return fmt.Errorf("failed to update user plan: %w", err)
		}
	}

	return nil
}
