package system

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"

	"codedock.run/codedock/internal/models"
	systemservices "codedock.run/codedock/internal/services/system"
)

type BillingHandler struct {
	billingSvc *systemservices.BillingService
}

func NewBillingHandler(billingSvc *systemservices.BillingService) *BillingHandler {
	return &BillingHandler{
		billingSvc: billingSvc,
	}
}

type CreateCheckoutRequest struct {
	PriceID    string `json:"priceId"`
	SuccessURL string `json:"successUrl"`
	CancelURL  string `json:"cancelUrl"`
}

type CreateCheckoutResponse struct {
	URL string `json:"url"`
}

func (h *BillingHandler) CreateCheckoutSession(c echo.Context) error {
	claims, ok := c.Get("user").(*models.UserClaims)
	if !ok || claims == nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
	}

	var req CreateCheckoutRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	url, err := h.billingSvc.CreateCheckoutSession(c.Request().Context(), claims.UserID, req.PriceID, req.SuccessURL, req.CancelURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, CreateCheckoutResponse{URL: url})
}

func (h *BillingHandler) Webhook(c echo.Context) error {
	payload, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "failed to read body")
	}

	sigHeader := c.Request().Header.Get("Stripe-Signature")
	if sigHeader == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing signature")
	}

	if err := h.billingSvc.HandleWebhook(c.Request().Context(), payload, sigHeader); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (h *BillingHandler) GetConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"publishableKey": os.Getenv("STRIPE_PUBLISHABLE_KEY"),
		"plans": []map[string]any{
			{
				"id":       "free",
				"name":     "Hobby",
				"price":    0,
				"features": []string{"1 Project", "Basic Support"},
			},
			{
				"id":       "pro",
				"name":     "Pro",
				"price":    9.99,
				"priceId":  os.Getenv("STRIPE_PRICE_ID_PRO"),
				"features": []string{"Unlimited Projects", "Priority Support", "Analytics"},
			},
		},
	})
}
