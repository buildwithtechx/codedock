package http

import (
	"context"
	"errors"

	"codedock/internal/modules/monitoring"
)

type billingServiceAdapter struct {
	svc *monitoring.BillingService
}

func (b *billingServiceAdapter) GetConfig(ctx context.Context) (any, error) {
	return nil, nil
}

func (b *billingServiceAdapter) CreateCheckoutSession(ctx context.Context, userID string, opts monitoring.CreateCheckoutSessionOpts) (string, error) {
	if b.svc == nil {
		return "", errors.New("billing service unavailable")
	}
	return b.svc.CreateCheckoutSession(
		ctx,
		userID,
		opts.PriceID,
		opts.SuccessURL,
		opts.CancelURL,
		opts.IdempotencyKey,
		opts.PurchaseIntent,
		opts.CancelPending,
	)
}

func (b *billingServiceAdapter) GetPendingCheckout(ctx context.Context, userID string) (any, error) {
	if b.svc == nil {
		return nil, nil
	}
	return b.svc.GetPendingCheckout(ctx, userID), nil
}

func (b *billingServiceAdapter) CancelCheckoutSession(ctx context.Context, userID, id, idempotencyKey string) (any, error) {
	if b.svc == nil {
		return nil, errors.New("billing service unavailable")
	}
	return b.svc.CancelCheckoutSession(ctx, userID, id, idempotencyKey)
}

func (b *billingServiceAdapter) HandleWebhook(ctx context.Context, payload []byte, sig string) error {
	if b.svc == nil {
		return errors.New("billing service unavailable")
	}
	return b.svc.HandleWebhook(ctx, payload, sig)
}
