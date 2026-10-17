import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { Coins, CreditCard, ExternalLink, Layers, Receipt, Sparkles, Zap } from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@codedock/ui/tabs';
import {
  BillingOverview,
  BillingPlansCards,
  BillingTopups,
  CheckoutRecovery,
  CustomPlanConfigurator,
  type CustomPlanQuote,
  PayAsYouGoPlan,
  ResourceMetersOverview,
} from '#/features/billing';
import { useGetPublicSettings } from '#/features/settings';
import { useBillingConfig, useCreateCheckoutSession } from '#/hooks/use-billing';
import { billingService } from '#/services/billing';
import { useAuthStore } from '#/stores/auth-store';

export const Route = createFileRoute('/_dashboard/billing')({
  component: BillingPage,
});

function BillingPage() {
  const navigate = useNavigate();
  const { data: publicRes, isLoading: isLoadingPublic } = useGetPublicSettings();
  const isCloud = publicRes?.data?.cloudMode ?? false;

  const user = useAuthStore((s) => s.user);
  const { data: config } = useBillingConfig();
  const { mutateAsync: createCheckoutSession, isPending } = useCreateCheckoutSession();
  const [activeTab, setActiveTab] = useState<
    'overview' | 'plans' | 'usage' | 'topups' | 'invoices'
  >('plans');
  const [planMode, setPlanMode] = useState<'plans' | 'custom'>('plans');
  const [billingCycle, setBillingCycle] = useState<'monthly' | 'payg'>('monthly');
  const [loadingPortal, setLoadingPortal] = useState(false);
  const [showPending, setShowPending] = useState(false);

  useEffect(() => {
    if (!isLoadingPublic && !isCloud) {
      void navigate({ to: '/' });
    }
  }, [isLoadingPublic, isCloud, navigate]);

  if (!isCloud) {
    return null;
  }

  const hasSubscription = Boolean(
    user?.planType && user.planType !== 'free' && user.planType !== 'hobby'
  );

  const handleSelectTier = async (planId: string) => {
    const plan = config?.plans.find((p) => p.id === planId) ?? config?.plans[0];
    if (!plan?.priceId) return;
    try {
      const { url } = await createCheckoutSession({
        priceId: plan.priceId,
        successUrl: `${window.location.origin}${window.location.pathname}?billing=success`,
        cancelUrl: `${window.location.origin}${window.location.pathname}?billing=canceled`,
      });
      window.location.href = url;
    } catch (e) {
      console.error(e);
    }
  };

  const handleSelectCustomPlan = async (_quote: CustomPlanQuote) => {
    const proPlan = config?.plans.find((p) => p.id === 'pro');
    if (!proPlan?.priceId) return;
    try {
      const { url } = await createCheckoutSession({
        priceId: proPlan.priceId,
        successUrl: `${window.location.origin}${window.location.pathname}?billing=success`,
        cancelUrl: `${window.location.origin}${window.location.pathname}?billing=canceled`,
      });
      window.location.href = url;
    } catch (e) {
      console.error(e);
    }
  };

  const handleAddCredits = async (cents: number) => {
    try {
      const packs = await billingService.getTopupPacks();
      const pack = packs.find((p) => p.priceCents === cents) ?? packs[0];
      const { url } = await billingService.purchaseTopup({
        packId: pack.id,
        idempotencyKey: crypto.randomUUID(),
      });
      if (url) window.location.href = url;
    } catch (e) {
      console.error(e);
    }
  };

  const handleOpenPortal = async () => {
    setLoadingPortal(true);
    try {
      const { url } = await billingService.createPortalSession();
      if (url) window.location.href = url;
    } catch (e) {
      console.error(e);
    } finally {
      setLoadingPortal(false);
    }
  };

  const plansView = (
    <div className="space-y-6">
      <div className="flex items-center gap-4 border-border/50 border-b pb-2 text-sm">
        <button
          type="button"
          onClick={() => setBillingCycle('monthly')}
          className={`relative pb-2 font-medium transition-colors ${
            billingCycle === 'monthly'
              ? 'text-foreground'
              : 'text-muted-foreground hover:text-foreground'
          }`}
        >
          Monthly plans
          {billingCycle === 'monthly' && (
            <span className="absolute inset-x-0 -bottom-2 h-0.5 rounded-full bg-primary" />
          )}
        </button>
        <button
          type="button"
          onClick={() => setBillingCycle('payg')}
          className={`relative pb-2 font-medium transition-colors ${
            billingCycle === 'payg'
              ? 'text-foreground'
              : 'text-muted-foreground hover:text-foreground'
          }`}
        >
          Prepaid PAYG
          {billingCycle === 'payg' && (
            <span className="absolute inset-x-0 -bottom-2 h-0.5 rounded-full bg-primary" />
          )}
        </button>
      </div>

      {billingCycle === 'payg' ? (
        <PayAsYouGoPlan busy={isPending} onAddCredits={handleAddCredits} />
      ) : planMode === 'plans' ? (
        <BillingPlansCards busy={isPending} onSelectPlan={handleSelectTier} />
      ) : (
        <CustomPlanConfigurator busy={isPending} onSelectPlan={handleSelectCustomPlan} />
      )}
    </div>
  );

  return (
    <div className="space-y-6">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="font-medium text-2xl text-foreground tracking-tight">
            Find the right plan
          </h1>
          <p className="mt-1 text-muted-foreground text-sm">
            Choose capacity and deploy. Codedock handles setup, routing and maintenance.
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <div className="inline-flex rounded-xl border border-border/60 bg-muted/40 p-1">
            <button
              type="button"
              onClick={() => setPlanMode('plans')}
              className={`rounded-lg px-3 py-1 font-medium text-xs transition-colors ${
                planMode === 'plans'
                  ? 'bg-card text-foreground shadow-xs'
                  : 'text-muted-foreground hover:text-foreground'
              }`}
            >
              Plans
            </button>
            <button
              type="button"
              onClick={() => setPlanMode('custom')}
              className={`rounded-lg px-3 py-1 font-medium text-xs transition-colors ${
                planMode === 'custom'
                  ? 'bg-card text-foreground shadow-xs'
                  : 'text-muted-foreground hover:text-foreground'
              }`}
            >
              Custom
            </button>
          </div>

          <Button
            variant="outline"
            size="sm"
            disabled={loadingPortal}
            onClick={() => setShowPending((v) => !v)}
            className="gap-1.5 rounded-xl text-xs"
          >
            <CreditCard className="size-3.5" />
            Pending payments
          </Button>
        </div>
      </header>

      {showPending && <CheckoutRecovery />}

      {!hasSubscription ? (
        plansView
      ) : (
        <Tabs
          value={activeTab}
          onValueChange={(val) => setActiveTab(val as typeof activeTab)}
          className="space-y-6"
        >
          <TabsList className="rounded-xl">
            <TabsTrigger value="plans" className="gap-2 rounded-lg">
              <Sparkles className="size-4" />
              Plans & Pricing
            </TabsTrigger>
            <TabsTrigger value="overview" className="gap-2 rounded-lg">
              <Layers className="size-4" />
              Overview
            </TabsTrigger>
            <TabsTrigger value="usage" className="gap-2 rounded-lg">
              <Zap className="size-4" />
              Usage & Quota
            </TabsTrigger>
            <TabsTrigger value="topups" className="gap-2 rounded-lg">
              <Coins className="size-4" />
              Prepaid Credits
            </TabsTrigger>
            <TabsTrigger value="invoices" className="gap-2 rounded-lg">
              <Receipt className="size-4" />
              Invoices & Portal
            </TabsTrigger>
          </TabsList>

          <TabsContent value="plans" className="space-y-6">
            {plansView}
          </TabsContent>

          <TabsContent value="overview" className="space-y-6">
            <BillingOverview
              planType={user?.planType ?? 'hobby'}
              onUpgradeClick={() => setActiveTab('plans')}
              onPortalClick={handleOpenPortal}
            />
            <ResourceMetersOverview />
          </TabsContent>

          <TabsContent value="usage" className="space-y-6">
            <ResourceMetersOverview />
          </TabsContent>

          <TabsContent value="topups" className="space-y-6">
            <BillingTopups />
          </TabsContent>

          <TabsContent value="invoices" className="space-y-6">
            <div className="rounded-2xl border border-border bg-card p-6 shadow-sm">
              <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
                <div>
                  <h3 className="font-semibold text-base text-foreground">
                    Payment Method & Invoices
                  </h3>
                  <p className="mt-1 max-w-lg text-muted-foreground text-xs leading-relaxed">
                    Access your Stripe customer billing portal to download official tax invoices,
                    update credit cards, and manage enterprise billing details.
                  </p>
                </div>
                <Button
                  variant="outline"
                  disabled={loadingPortal}
                  onClick={handleOpenPortal}
                  className="shrink-0 gap-1.5"
                >
                  <ExternalLink className="size-4" />
                  Open Customer Portal
                </Button>
              </div>
            </div>
          </TabsContent>
        </Tabs>
      )}
    </div>
  );
}
