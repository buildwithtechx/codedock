package system

import (
	"context"
	"errors"
	"testing"

	"github.com/stripe/stripe-go/v78"

	"codedock/internal/models"
)

func TestPurchaseIntentValidation(t *testing.T) {
	tests := []struct {
		intent string
		valid  bool
	}{
		{"c3b12345-6789-4abc-8def-0123456789ab", true},
		{"C3B12345-6789-4ABC-8DEF-0123456789AB", true},
		{"not-a-uuid", false},
		{"", false},
		{"12345", false},
		{"c3b12345-6789-1abc-8def-0123456789ab", false}, // version 1, not 4
	}

	for _, tt := range tests {
		got := models.IsPurchaseIntent(tt.intent)
		if got != tt.valid {
			t.Errorf("IsPurchaseIntent(%q) = %v, want %v", tt.intent, got, tt.valid)
		}
	}
}

func TestBillingCheckoutHygiene(t *testing.T) {
	svc := NewBillingService(nil)

	// Mock Stripe session functions
	sessCounter := 0
	svc.SetSessionHandlers(
		func(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
			sessCounter++
			return &stripe.CheckoutSession{
				ID:  fmtSessionID(sessCounter),
				URL: "https://checkout.stripe.com/pay/" + fmtSessionID(sessCounter),
			}, nil
		},
		func(id string, params *stripe.CheckoutSessionExpireParams) (*stripe.CheckoutSession, error) {
			return &stripe.CheckoutSession{
				ID:     id,
				Status: stripe.CheckoutSessionStatusExpired,
			}, nil
		},
	)

	ctx := context.Background()
	userID := "user-alice"
	validIntent := "a1b2c3d4-e5f6-4a1b-8c2d-3e4f5a6b7c8d"

	t.Run("rejects invalid purchase intent", func(t *testing.T) {
		_, err := svc.CreateCheckoutSession(ctx, userID, "price_pro", "https://app/ok", "https://app/cancel", "key-1", "invalid-intent", false)
		if err == nil || err.Error() != "invalid purchase intent format" {
			t.Fatalf("expected invalid purchase intent format error, got: %v", err)
		}
	})

	t.Run("creates first checkout and tracks pending", func(t *testing.T) {
		url, err := svc.CreateCheckoutSession(ctx, userID, "price_pro", "https://app/ok", "https://app/cancel", "key-1", validIntent, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if url == "" {
			t.Errorf("expected non-empty url")
		}

		pending := svc.GetPendingCheckout(ctx, userID)
		if pending == nil {
			t.Fatalf("expected pending checkout to be tracked")
		}
		if pending.Status != "open" || !pending.CanCancel {
			t.Errorf("expected open and cancellable, got %+v", pending)
		}
	})

	t.Run("refuses second checkout while pending is open", func(t *testing.T) {
		_, err := svc.CreateCheckoutSession(ctx, userID, "price_enterprise", "https://app/ok", "https://app/cancel", "key-2", validIntent, false)
		if err == nil {
			t.Fatalf("expected conflict error, got nil")
		}
		var pendingErr *ErrCheckoutPending
		if !errors.As(err, &pendingErr) {
			t.Fatalf("expected ErrCheckoutPending, got: %v", err)
		}
		if pendingErr.Pending.CheckoutID != "cs_test_1" {
			t.Errorf("expected pending checkout cs_test_1, got: %s", pendingErr.Pending.CheckoutID)
		}
	})

	t.Run("idempotent replay returns same checkout URL", func(t *testing.T) {
		url, err := svc.CreateCheckoutSession(ctx, userID, "price_pro", "https://app/ok", "https://app/cancel", "key-1", validIntent, false)
		if err != nil {
			t.Fatalf("unexpected error on idempotent replay: %v", err)
		}
		if url != "https://checkout.stripe.com/pay/cs_test_1" {
			t.Errorf("expected original checkout URL, got: %s", url)
		}
	})

	t.Run("cancelPending settles previous checkout and creates new session", func(t *testing.T) {
		url, err := svc.CreateCheckoutSession(ctx, userID, "price_enterprise", "https://app/ok", "https://app/cancel", "key-3", validIntent, true)
		if err != nil {
			t.Fatalf("unexpected error with cancelPending: %v", err)
		}
		if url != "https://checkout.stripe.com/pay/cs_test_2" {
			t.Errorf("expected new checkout URL cs_test_2, got %s", url)
		}

		current := svc.GetPendingCheckout(ctx, userID)
		if current == nil || current.CheckoutID != "cs_test_2" {
			t.Fatalf("expected active checkout to be cs_test_2, got %+v", current)
		}
	})

	t.Run("explicit cancellation marks session canceled", func(t *testing.T) {
		canceled, err := svc.CancelCheckoutSession(ctx, userID, "cs_test_2", "cancel-idem-1")
		if err != nil {
			t.Fatalf("unexpected error canceling checkout: %v", err)
		}
		if canceled.Status != "canceled" || canceled.CanCancel {
			t.Errorf("expected canceled status and CanCancel=false, got %+v", canceled)
		}

		openPending := svc.GetPendingCheckout(ctx, userID)
		if openPending != nil {
			t.Errorf("expected no open pending checkout, got %+v", openPending)
		}
	})
}

func fmtSessionID(n int) string {
	switch n {
	case 1:
		return "cs_test_1"
	case 2:
		return "cs_test_2"
	case 3:
		return "cs_test_3"
	default:
		return "cs_test_n"
	}
}
