import { CreditCard, Loader2 } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { useBillingConfig, useCreateCheckoutSession } from '#/hooks/useBilling';
import { useAuthStore } from '#/stores/authStore';

export function BillingSection() {
  const user = useAuthStore((s) => s.user);
  const { data: config, isLoading } = useBillingConfig();
  const { mutateAsync: createCheckoutSession, isPending } = useCreateCheckoutSession();

  const handleUpgrade = async (priceId: string) => {
    try {
      const { url } = await createCheckoutSession({
        priceId,
        successUrl: `${window.location.href}?billing=success`,
        cancelUrl: `${window.location.href}?billing=canceled`,
      });
      window.location.href = url;
    } catch (e) {
      console.error(e);
    }
  };

  const isPro = user?.planType === 'pro';

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <CreditCard className="h-5 w-5" />
          Billing & Plan
        </CardTitle>
        <CardDescription>Manage your subscription and billing details.</CardDescription>
      </CardHeader>
      <CardContent>
        {isLoading ? (
          <Loader2 className="h-5 w-5 animate-spin" />
        ) : (
          <div className="space-y-4">
            <div className="flex items-center justify-between rounded-lg border bg-card p-4">
              <div>
                <h3 className="font-semibold text-lg">{isPro ? 'Pro Plan' : 'Hobby Plan'}</h3>
                <p className="text-muted-foreground text-sm">
                  {isPro
                    ? 'You are currently on the Pro plan with all features unlocked.'
                    : 'You are on the free Hobby plan. Upgrade to unlock more features.'}
                </p>
              </div>
              {!isPro && config?.plans.find((p) => p.id === 'pro')?.priceId && (
                <Button
                  onClick={() => handleUpgrade(config.plans.find((p) => p.id === 'pro')!.priceId!)}
                  disabled={isPending}
                >
                  {isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                  Upgrade to Pro
                </Button>
              )}
              {isPro && (
                <Button variant="outline" disabled>
                  Manage Subscription
                </Button>
              )}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
