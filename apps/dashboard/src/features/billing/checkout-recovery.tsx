import { AlertCircle, ExternalLink, Loader2, RefreshCw, XCircle } from 'lucide-react';
import { useState } from 'react';
import { Button } from '@codedock/ui/button';
import type { PendingCheckout } from '#/interfaces/billing';
import { billingService } from '#/services/billing';

export interface CheckoutRecoveryProps {
  pending?: PendingCheckout | null;
  onSettled?: () => void;
  onCancel?: (checkoutId: string, idempotencyKey: string) => Promise<void> | void;
  onRefresh?: () => Promise<void> | void;
}

export function CheckoutRecovery({
  pending,
  onSettled,
  onCancel,
  onRefresh,
}: CheckoutRecoveryProps) {
  const [canceling, setCanceling] = useState(false);
  const [refreshing, setRefreshing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!pending || pending.status === 'complete') {
    return null;
  }

  const handleCancel = async () => {
    if (!pending.canCancel || canceling) return;

    const checkoutId = pending.checkoutId;
    const idempotencyKey = pending.idempotencyKey || `cancel_${checkoutId}`;

    setCanceling(true);
    setError(null);

    try {
      if (onCancel) {
        await onCancel(checkoutId, idempotencyKey);
      } else {
        const result = await billingService.cancelCheckoutSession(checkoutId, idempotencyKey);
        if (result.checkout.checkoutId !== checkoutId) {
          throw new Error('Mismatched checkout response during cancellation');
        }
      }
      onSettled?.();
    } catch {
      // Per v0.9.2 hardening rule: never clear on uncertain outcomes, keep tracked
      setError('Unable to verify cancellation status. Retaining checkout session for retry.');
    } finally {
      setCanceling(false);
    }
  };

  const handleRefresh = async () => {
    setRefreshing(true);
    setError(null);
    try {
      if (onRefresh) {
        await onRefresh();
      } else {
        const res = await billingService.getPendingCheckout();
        if (!res.pending || res.pending.status === 'complete') {
          onSettled?.();
        }
      }
    } catch {
      setError('Failed to refresh settlement status.');
    } finally {
      setRefreshing(false);
    }
  };

  const cancelNotice =
    !pending.canCancel && pending.cancelUnavailableReason === 'provider_unavailable'
      ? 'Cancellation is temporarily unavailable from the payment provider.'
      : !pending.canCancel && pending.cancelUnavailableReason === 'unsupported'
        ? 'This checkout cannot be canceled online.'
        : null;

  return (
    <div
      role="region"
      aria-label="Pending Checkout Recovery"
      className="rounded-2xl border border-amber-500/30 bg-amber-500/5 p-5 shadow-sm"
    >
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div className="flex items-start gap-3">
          <div className="rounded-xl bg-amber-500/10 p-2 text-amber-500">
            <AlertCircle className="size-5" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h3 className="font-semibold text-base text-foreground">
                Unsettled Payment in Progress
              </h3>
              <span className="inline-flex items-center rounded-full bg-amber-500/10 px-2 py-0.5 font-medium text-amber-500 text-xs uppercase">
                {pending.status}
              </span>
            </div>
            <p className="mt-1 text-muted-foreground text-sm">
              A checkout session is awaiting settlement. You can resume payment or safely cancel the
              session.
            </p>
            {cancelNotice && <p className="mt-1 text-muted-foreground text-xs">{cancelNotice}</p>}
          </div>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          {pending.url && pending.status === 'open' && (
            <Button size="sm" asChild className="gap-1.5">
              <a href={pending.url} target="_blank" rel="noopener noreferrer">
                <ExternalLink className="size-3.5" />
                Resume Payment
              </a>
            </Button>
          )}

          <Button
            variant="outline"
            size="sm"
            disabled={!pending.canCancel || canceling}
            onClick={handleCancel}
            className="gap-1.5"
          >
            {canceling ? (
              <Loader2 className="size-3.5 animate-spin" />
            ) : (
              <XCircle className="size-3.5" />
            )}
            Cancel Checkout
          </Button>

          <Button
            variant="ghost"
            size="sm"
            disabled={refreshing}
            onClick={handleRefresh}
            className="gap-1.5"
            aria-label="Refresh Status"
          >
            <RefreshCw className={`size-3.5 ${refreshing ? 'animate-spin' : ''}`} />
          </Button>
        </div>
      </div>

      {error && (
        <div
          className="mt-3 rounded-lg bg-destructive/10 p-2.5 text-destructive text-xs"
          role="alert"
        >
          {error}
        </div>
      )}
    </div>
  );
}
