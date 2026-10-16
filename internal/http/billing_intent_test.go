package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v78"

	"codedock/internal/config"
	"codedock/internal/handlers/system"
	"codedock/internal/models"
	systemservices "codedock/internal/services/system"
)

func TestBillingHandlerPendingFlow(t *testing.T) {
	cfg := config.Get()
	cfg.Cloud.Enabled = true
	t.Cleanup(func() { cfg.Cloud.Enabled = false })

	svc := systemservices.NewBillingService(nil)
	svc.SetSessionHandlers(
		func(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
			return &stripe.CheckoutSession{
				ID:  "cs_http_1",
				URL: "https://checkout.stripe.com/pay/cs_http_1",
			}, nil
		},
		func(id string, params *stripe.CheckoutSessionExpireParams) (*stripe.CheckoutSession, error) {
			return &stripe.CheckoutSession{
				ID:     id,
				Status: stripe.CheckoutSessionStatusExpired,
			}, nil
		},
	)

	handler := system.NewBillingHandler(svc)
	e := echo.New()

	authMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("user", &models.UserClaims{
				UserID: "user_test_http",
				Email:  "billing@test.local",
			})
			return next(c)
		}
	}

	e.POST("/checkout", handler.CreateCheckoutSession, authMiddleware)
	e.GET("/pending", handler.GetPendingCheckout, authMiddleware)
	e.POST("/checkouts/:id/cancel", handler.CancelCheckoutSession, authMiddleware)

	t.Run("first checkout succeeds", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{
			"priceId":        "price_pro",
			"successUrl":     "https://app/ok",
			"cancelUrl":      "https://app/cancel",
			"idempotencyKey": "key-http-1",
			"purchaseIntent": "9b1deb4d-3b7d-4bad-9bdd-2b0d7b3dcb6d",
		})
		req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
	})

	t.Run("second checkout returns 409 conflict with pending details", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{
			"priceId":        "price_other",
			"successUrl":     "https://app/ok",
			"cancelUrl":      "https://app/cancel",
			"idempotencyKey": "key-http-2",
		})
		req := httptest.NewRequest(http.MethodPost, "/checkout", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("expected 409 Conflict, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]any
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["code"] != "CLOUD_WORKSPACE_CHECKOUT_PENDING" {
			t.Errorf("expected code CLOUD_WORKSPACE_CHECKOUT_PENDING, got %v", resp["code"])
		}
	})

	t.Run("pending endpoint returns active checkout", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pending", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}
		var resp map[string]any
		json.Unmarshal(rec.Body.Bytes(), &resp)
		pending := resp["pending"].(map[string]any)
		if pending["checkoutId"] != "cs_http_1" {
			t.Errorf("expected checkoutId cs_http_1, got %v", pending["checkoutId"])
		}
	})

	t.Run("cancel endpoint cancels pending checkout", func(t *testing.T) {
		body, _ := json.Marshal(map[string]any{"idempotencyKey": "cancel-key-1"})
		req := httptest.NewRequest(http.MethodPost, "/checkouts/cs_http_1/cancel", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200 OK, got %d: %s", rec.Code, rec.Body.String())
		}

		var resp map[string]any
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if resp["status"] != "canceled" {
			t.Errorf("expected status canceled, got %v", resp["status"])
		}
	})
}
