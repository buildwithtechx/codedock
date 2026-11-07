import { AlertCircle, CheckCircle2, Loader2, Sparkles, X } from 'lucide-react';
import { useCallback, useEffect, useId, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { useBillingConfig, useCreateCheckoutSession } from '#/hooks/use-billing';
import { useAuthStore } from '#/stores/auth-store';
import { BillingPlansCards } from './billing-plans-cards';

export interface CloudDeployRestriction {
  code: 'CLOUD_BILLING_BLOCKED' | 'CLOUD_CAPACITY_REQUIRED' | 'PLAN_UPGRADE_REQUIRED';
  reason?: string;
  message?: string;
}

export interface CloudDeployPlanModalProps {
  open?: boolean;
  onClose: () => void;
  restriction?: CloudDeployRestriction;
  workspaceId?: string;
  serverName?: string;
  initialCheckoutUrl?: string;
  onRetry?: () => Promise<unknown>;
  onSelectPlan?: (plan: any) => void;
}

export function CloudDeployPlanModal({
  open = true,
  onClose,
  restriction = { code: 'CLOUD_BILLING_BLOCKED' },
  workspaceId,
  serverName,
  initialCheckoutUrl,
  onRetry,
  onSelectPlan,
}: CloudDeployPlanModalProps) {
  const titleId = useId();
  const descId = useId();
  const user = useAuthStore((s) => s.user);
  const { data: config } = useBillingConfig();
  const { mutateAsync: createCheckoutSession, isPending } = useCreateCheckoutSession();

  const [checkoutStarted, setCheckoutStarted] = useState(Boolean(initialCheckoutUrl));
  const [retrying, setRetrying] = useState(false);
  const [retryError, setRetryError] = useState<string | null>(null);
  const [ready, setReady] = useState(false);

  // Poll or detect when user returns from payment window
  const checkState = useCallback(() => {
    if (checkoutStarted && !ready) {
      // If plan has upgraded in auth store
      if (user?.planType && user.planType !== 'free' && user.planType !== 'hobby') {
        setReady(true);
      }
    }
  }, [checkoutStarted, ready, user?.planType]);

  useEffect(() => {
    if (ready) return;
    const handleFocus = () => {
      if (document.visibilityState !== 'hidden') checkState();
    };
    window.addEventListener('focus', handleFocus);
    const interval = checkoutStarted ? window.setInterval(checkState, 4000) : undefined;
    return () => {
      window.removeEventListener('focus', handleFocus);
      if (interval !== undefined) window.clearInterval(interval);
    };
  }, [checkoutStarted, ready, checkState]);

  if (!open) return null;

  const handleSelectPlan = async (planId: string) => {
    const plan = config?.plans.find((p) => p.id === planId) ?? config?.plans[0];
    if (plan && onSelectPlan) {
      onSelectPlan(plan);
    }
    if (!plan?.priceId) return;
    try {
      const { url } = await createCheckoutSession({
        priceId: plan.priceId,
        successUrl: `${window.location.origin}${window.location.pathname}?billing=success&workspace=${workspaceId ?? ''}`,
        cancelUrl: `${window.location.origin}${window.location.pathname}?billing=canceled`,
      });
      setCheckoutStarted(true);
      window.open(url, '_blank', 'noopener,noreferrer');
    } catch (e) {
      setRetryError(e instanceof Error ? e.message : 'Failed to create checkout session');
    }
  };

  const handleRetryAction = async () => {
    if (!onRetry || retrying) return;
    setRetrying(true);
    setRetryError(null);
    try {
      await onRetry();
      onClose();
    } catch (err) {
      setRetryError(err instanceof Error ? err.message : 'Retry failed. Check plan limits.');
    } finally {
      setRetrying(false);
    }
  };

  const modalTitle = ready
    ? 'Plan ready'
    : serverName
      ? `Configure capacity for ${serverName}`
      : restriction.code === 'CLOUD_CAPACITY_REQUIRED'
        ? 'Compute capacity limit reached'
        : 'Choose a Cloud plan';

  const modalDesc = ready
    ? 'Your subscription is active and capacity has been allocated. You can now retry your deployment.'
    : (restriction.message ??
      'Your workspace requires an active dedicated plan to provision managed servers, deploy full-stack apps, or execute background actions.');

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby={titleId}
      aria-describedby={descId}
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-xs"
    >
      <div className="relative flex max-h-[90vh] w-full max-w-4xl flex-col overflow-hidden rounded-2xl border border-border bg-card shadow-2xl">
        <header className="flex shrink-0 items-center justify-between border-border/60 border-b px-6 py-4">
          <div className="flex items-center gap-2.5">
            <div className="flex size-9 items-center justify-center rounded-xl bg-primary/10 text-primary">
              <Sparkles className="size-4" />
            </div>
            <div>
              <h2 id={titleId} className="font-semibold text-foreground text-lg tracking-tight">
                {modalTitle}
              </h2>
            </div>
          </div>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            onClick={onClose}
            aria-label="Close dialog"
            className="rounded-xl text-muted-foreground hover:text-foreground"
          >
            <X className="size-4" />
          </Button>
        </header>

        <div className="min-h-0 flex-1 overflow-y-auto p-6">
          <p id={descId} className="mb-6 text-muted-foreground text-sm">
            {modalDesc}
          </p>

          {retryError && (
            <div
              role="alert"
              className="mb-6 flex items-center gap-2 rounded-xl bg-destructive/10 p-3 text-destructive text-sm"
            >
              <AlertCircle className="size-4 shrink-0" />
              <span>{retryError}</span>
            </div>
          )}

          {ready ? (
            <div className="flex flex-col items-center justify-center py-10 text-center">
              <div className="flex size-12 items-center justify-center rounded-full bg-emerald-500/10 text-emerald-500">
                <CheckCircle2 className="size-6" />
              </div>
              <h3 className="mt-4 font-semibold text-foreground text-lg">Capacity unlocked</h3>
              <p className="mt-1 max-w-sm text-muted-foreground text-sm">
                Your account is verified and ready. Resume your action or deployment.
              </p>
              <div className="mt-6 flex gap-3">
                <Button onClick={onClose} variant="outline" className="rounded-xl">
                  Close
                </Button>
                {onRetry && (
                  <Button onClick={handleRetryAction} disabled={retrying} className="rounded-xl">
                    {retrying && <Loader2 className="mr-2 size-4 animate-spin" />}
                    Retry action
                  </Button>
                )}
              </div>
            </div>
          ) : (
            <div className="space-y-6">
              <BillingPlansCards busy={isPending} onSelectPlan={handleSelectPlan} />
              {checkoutStarted && (
                <div className="flex items-center justify-between rounded-xl border border-primary/20 bg-primary/5 p-4">
                  <div className="flex items-center gap-3">
                    <Loader2 className="size-4 animate-spin text-primary" />
                    <div>
                      <p className="font-medium text-foreground text-sm">Checkout in progress...</p>
                      <p className="text-muted-foreground text-xs">
                        Complete payment in the opened tab. This window will update automatically.
                      </p>
                    </div>
                  </div>
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => setReady(true)}
                    className="text-xs"
                  >
                    I completed payment
                  </Button>
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
