package http

import (
	"context"
	"errors"

	"codedock/internal/modules/monitoring"
	systemservices "codedock/internal/services/system"
)

type monitoringPendingError struct {
	err     error
	pending any
}

func (m *monitoringPendingError) Error() string {
	return m.err.Error()
}

func (m *monitoringPendingError) PendingSession() any {
	return m.pending
}

type billingServiceAdapter struct {
	svc *systemservices.BillingService
}

func (b *billingServiceAdapter) GetConfig(ctx context.Context) (any, error) {
	return nil, nil
}

func (b *billingServiceAdapter) CreateCheckoutSession(ctx context.Context, userID string, opts monitoring.CreateCheckoutSessionOpts) (string, error) {
	if b.svc == nil {
		return "", errors.New("billing service unavailable")
	}
	url, err := b.svc.CreateCheckoutSession(
		ctx,
		userID,
		opts.PriceID,
		opts.SuccessURL,
		opts.CancelURL,
		opts.IdempotencyKey,
		opts.PurchaseIntent,
		opts.CancelPending,
	)
	if err != nil {
		var pendingErr *systemservices.ErrCheckoutPending
		if errors.As(err, &pendingErr) {
			return "", &monitoringPendingError{err: pendingErr, pending: pendingErr.Pending}
		}
		return "", err
	}
	return url, nil
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
