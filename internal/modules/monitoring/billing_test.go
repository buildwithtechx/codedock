package monitoring

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

type mockBillingService struct{}

func (m *mockBillingService) GetConfig(ctx context.Context) (any, error) {
	return &BillingConfig{
		Enabled:        true,
		PublishableKey: "pk_test_123",
		Mode:           "test",
	}, nil
}
func (m *mockBillingService) CreateCheckoutSession(ctx context.Context, userID string, opts CreateCheckoutSessionOpts) (string, error) {
	return "https://stripe.com/checkout", nil
}
func (m *mockBillingService) GetPendingCheckout(ctx context.Context, userID string) (any, error) {
	return map[string]string{"status": "pending"}, nil
}
func (m *mockBillingService) CancelCheckoutSession(ctx context.Context, userID, id, idempotencyKey string) (any, error) {
	return map[string]string{"id": id}, nil
}
func (m *mockBillingService) HandleWebhook(ctx context.Context, payload []byte, sig string) error {
	return nil
}

func TestBillingHandler_GetConfig(t *testing.T) {
	e := echo.New()
	handler := NewBillingHandler(&mockBillingService{})

	req := httptest.NewRequest(http.MethodGet, "/billing/config", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := handler.GetConfig(c); err != nil {
		t.Fatalf("unexpected GetConfig error: %v", err)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}
}

func TestBillingHandler_PendingRequiresAuth(t *testing.T) {
	e := echo.New()
	handler := NewBillingHandler(&mockBillingService{})

	req := httptest.NewRequest(http.MethodGet, "/billing/pending", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_ = handler.GetPendingCheckout(c)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized, got %d", rec.Code)
	}
}
