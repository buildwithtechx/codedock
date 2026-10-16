import { useState } from 'react';
import { CreditCard, ExternalLink, Loader2, Sparkles } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { billingService } from '#/services/billing';

export interface BillingOverviewProps {
  planType?: 'hobby' | 'pro' | 'custom' | string;
  planName?: string;
  status?: 'active' | 'canceled' | 'past_due' | 'trialing';
  renewalDate?: string;
  cancelAtPeriodEnd?: boolean;
  onUpgradeClick?: () => void;
  onPortalClick?: () => Promise<void> | void;
}

export function BillingOverview({
  planType = 'hobby',
  planName,
  status = 'active',
  renewalDate,
  cancelAtPeriodEnd = false,
  onUpgradeClick,
  onPortalClick,
}: BillingOverviewProps) {
  const [loadingPortal, setLoadingPortal] = useState(false);
  const [portalError, setPortalError] = useState<string | null>(null);

  const isPro = planType === 'pro';
  const isCustom = planType === 'custom';
  const isPaid = isPro || isCustom;
  const displayName = planName ?? (isCustom ? 'Custom Cloud Dedicated' : isPro ? 'Pro Platform' : 'Hobby Developer');

  const handleOpenPortal = async () => {
    setLoadingPortal(true);
    setPortalError(null);
    try {
      if (onPortalClick) {
        await onPortalClick();
      } else {
        const { url } = await billingService.createPortalSession();
        if (url) window.location.href = url;
      }
    } catch {
      setPortalError('Unable to open billing portal. Please try again or contact support.');
    } finally {
      setLoadingPortal(false);
    }
  };

  return (
    <div className="rounded-2xl border border-border bg-card p-6 shadow-sm">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-start gap-3.5">
          <div className="rounded-xl bg-primary/10 p-2.5 text-primary">
            <CreditCard className="size-5" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h2 className="font-semibold text-foreground text-lg">{displayName}</h2>
              <span
                className={`inline-flex items-center rounded-full px-2.5 py-0.5 font-medium text-xs ${
                  isPaid && !cancelAtPeriodEnd
                    ? 'bg-emerald-500/10 text-emerald-500'
                    : cancelAtPeriodEnd
                      ? 'bg-amber-500/10 text-amber-500'
                      : 'bg-muted text-muted-foreground'
                }`}
              >
                {cancelAtPeriodEnd ? 'Canceling' : status === 'active' ? (isPaid ? 'Active' : 'Free') : status}
              </span>
            </div>
            <p className="mt-1 text-muted-foreground text-sm">
              {cancelAtPeriodEnd
                ? renewalDate
                  ? `Access remains active until ${renewalDate}`
                  : 'Subscription will not renew at the end of the billing period'
                : renewalDate
                  ? `Renews on ${renewalDate}`
                  : isPaid
                    ? 'Automatic monthly renewal active'
                    : 'Free community plan for personal workloads and testing'}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {isPaid ? (
            <Button
              variant="outline"
              size="sm"
              disabled={loadingPortal}
              onClick={handleOpenPortal}
              className="gap-1.5"
            >
              {loadingPortal ? (
                <Loader2 className="size-3.5 animate-spin" />
              ) : (
                <ExternalLink className="size-3.5" />
              )}
              Manage Billing
            </Button>
          ) : (
            <Button size="sm" onClick={onUpgradeClick} className="gap-1.5">
              <Sparkles className="size-3.5" />
              Upgrade to Pro
            </Button>
          )}
        </div>
      </div>

      {portalError && (
        <div className="mt-4 rounded-lg bg-destructive/10 p-3 text-destructive text-sm" role="alert">
          {portalError}
        </div>
      )}
    </div>
  );
}
