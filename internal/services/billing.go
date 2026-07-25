package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
	"github.com/stripe/stripe-go/v78/webhook"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/repositories"
)

type BillingService struct {
	userRepo *repositories.UserRepo
}

func NewBillingService(userRepo *repositories.UserRepo) *BillingService {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &BillingService{
		userRepo: userRepo,
	}
}

func (s *BillingService) CreateCheckoutSession(ctx context.Context, userID, priceID, successURL, cancelURL string) (string, error) {
	if stripe.Key == "" {
		return "", fmt.Errorf("billing is not configured")
	}

	u, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	customerID, err := s.getOrCreateStripeCustomer(ctx, u)
	if err != nil {
		return "", err
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
	}

	sess, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return sess.URL, nil
}

func (s *BillingService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
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
