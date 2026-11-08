package monitoring

import (
	"codedock/internal/config"
	"context"
	"fmt"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
)

func (s *BillingService) GetPendingCheckout(ctx context.Context, userID string) *PendingCheckout {
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
) (*PendingCheckout, error) {
	s.pendingMu.Lock()
	defer s.pendingMu.Unlock()

	userPending := s.pendingCheckouts[userID]
	var target *PendingCheckout
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
			reason := CancelUnavailableProviderUnavailable
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

func (s *BillingService) getOrCreateStripeCustomer(ctx context.Context, u *User) (string, error) {
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
