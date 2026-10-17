import { CreditCard, Loader2 } from 'lucide-react';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';
import { env } from '#/env';
import { SettingsRow } from '#/features/settings/settings-row';
import { SettingsSection } from '#/features/settings/settings-section';
import { useBillingConfig, useCreateCheckoutSession } from '#/hooks/use-billing';
import { refreshAuthSession } from '#/lib/auth-refresh';
import { billingService } from '#/services/billing';
import { useAuthStore } from '#/stores/auth-store';

export function BillingSection() {
  const user = useAuthStore((s) => s.user);
  const { data: config, isLoading } = useBillingConfig();
  const { mutateAsync: createCheckoutSession, isPending } = useCreateCheckoutSession();
  const [loadingPortal, setLoadingPortal] = useState(false);

  useEffect(() => {
    if (typeof window === 'undefined' || !window.location.search.includes('billing=success')) {
      return;
    }

    let attempts = 0;
    const maxAttempts = 6;
    const interval = setInterval(async () => {
      attempts++;
      await refreshAuthSession(env.VITE_API_URL);
      const currentUser = useAuthStore.getState().user;

      if (currentUser?.planType === 'pro' || attempts >= maxAttempts) {
        clearInterval(interval);
        const url = new URL(window.location.href);
        url.searchParams.delete('billing');
        window.history.replaceState({}, '', url.toString());
      }
    }, 1500);

    return () => clearInterval(interval);
  }, []);

  const handleUpgrade = async (priceId: string) => {
    try {
      const { url } = await createCheckoutSession({
        priceId,
        successUrl: `${window.location.origin}${window.location.pathname}?billing=success`,
        cancelUrl: `${window.location.origin}${window.location.pathname}?billing=canceled`,
      });
      window.location.href = url;
    } catch (e) {
      console.error(e);
    }
  };

  const handleManagePortal = async () => {
    setLoadingPortal(true);
    try {
      const { url } = await billingService.createPortalSession();
      if (url) window.location.href = url;
    } catch {
      toast.error('Unable to open billing portal');
    } finally {
      setLoadingPortal(false);
    }
  };

  const isPro = user?.planType === 'pro';

  return (
    <SettingsSection
      icon={<CreditCard className="size-4 text-primary" />}
      title="Billing & Plan"
      description="Manage your subscription and billing details."
    >
      {isLoading ? (
        <div className="flex items-center justify-center py-6">
          <Loader2 className="size-5 animate-spin text-muted-foreground" />
        </div>
      ) : (
        <SettingsRow
          label={isPro ? 'Pro Plan' : 'Hobby Plan'}
          description={
            isPro
              ? 'You are currently on the Pro plan with all features unlocked.'
              : 'You are on the free Hobby plan. Upgrade to unlock more features.'
          }
          bordered={false}
        >
          <div className="flex items-center gap-3">
            <span
              className={`inline-flex items-center rounded-full px-2.5 py-0.5 font-medium text-xs ${
                isPro ? 'bg-emerald-500/10 text-emerald-500' : 'bg-muted text-muted-foreground'
              }`}
            >
              {isPro ? 'Active' : 'Free'}
            </span>
            {!isPro && config?.plans.find((p) => p.id === 'pro')?.priceId && (
              <Button
                size="sm"
                onClick={() => handleUpgrade(config.plans.find((p) => p.id === 'pro')!.priceId!)}
                disabled={isPending}
              >
                {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Upgrade to Pro
              </Button>
            )}
            {isPro && (
              <Button
                variant="outline"
                size="sm"
                onClick={handleManagePortal}
                disabled={loadingPortal}
              >
                {loadingPortal && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Manage Subscription
              </Button>
            )}
          </div>
        </SettingsRow>
      )}
    </SettingsSection>
  );
}
