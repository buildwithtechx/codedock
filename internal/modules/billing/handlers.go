package billing

import (
	"errors"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"codedock/internal/config"
	"codedock/internal/utils"
)

type BillingHandler struct {
	service BillingOperations
}

func NewBillingHandler(service BillingOperations) *BillingHandler {
	return &BillingHandler{service: service}
}

func (h *BillingHandler) GetConfig(c echo.Context) error {
	if h.service != nil {
		cfg, err := h.service.GetConfig(c.Request().Context())
		if err == nil && cfg != nil {
			return c.JSON(http.StatusOK, cfg)
		}
	}
	cfg := config.Get()
	return c.JSON(http.StatusOK, map[string]any{
		"publishableKey": cfg.Stripe.PublishableKey,
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
				"priceId":  cfg.Stripe.PriceIDPro,
				"features": []string{"Unlimited Projects", "Priority Support", "Analytics"},
			},
		},
	})
}

func (h *BillingHandler) CreateCheckoutSession(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	var req CreateCheckoutRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if h.service == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "billing service not initialized")
	}

	opts := CreateCheckoutSessionOpts{
		PriceID:        req.PriceID,
		SuccessURL:     req.SuccessURL,
		CancelURL:      req.CancelURL,
		IdempotencyKey: req.IdempotencyKey,
		PurchaseIntent: req.PurchaseIntent,
		CancelPending:  req.CancelPending,
	}

	url, err := h.service.CreateCheckoutSession(c.Request().Context(), userID, opts)
	if err != nil {
		var pendErr CheckoutPendingError
		if errors.As(err, &pendErr) {
			return c.JSON(http.StatusConflict, map[string]any{
				"code":    "CLOUD_WORKSPACE_CHECKOUT_PENDING",
				"error":   pendErr.Error(),
				"pending": pendErr.PendingSession(),
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, CreateCheckoutResponse{URL: url})
}

func (h *BillingHandler) GetPendingCheckout(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	if h.service == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "billing service not initialized")
	}

	pending := h.service.GetPendingCheckout(c.Request().Context(), userID)
	return c.JSON(http.StatusOK, map[string]any{
		"pending": pending,
	})
}

func (h *BillingHandler) CancelCheckoutSession(c echo.Context) error {
	userID := utils.ExtractUserID(c)
	if userID == "" {
		return utils.Error(c, http.StatusUnauthorized, "unauthorized")
	}

	if h.service == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "billing service not initialized")
	}

	checkoutID := c.Param("id")
	if checkoutID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "checkout id is required")
	}

	var req CancelCheckoutRequest
	_ = c.Bind(&req)

	checkout, err := h.service.CancelCheckoutSession(c.Request().Context(), userID, checkoutID, req.IdempotencyKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"status":   "canceled",
		"checkout": checkout,
	})
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

	if h.service == nil {
		return echo.NewHTTPError(http.StatusServiceUnavailable, "billing service not initialized")
	}

	if err := h.service.HandleWebhook(c.Request().Context(), payload, sigHeader); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusOK)
}
