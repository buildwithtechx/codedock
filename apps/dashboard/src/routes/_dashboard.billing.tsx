import { createFileRoute } from '@tanstack/react-router';
import {
  Coins,
  CreditCard,
  ExternalLink,
  Layers,
  Receipt,
  Server,
  Sparkles,
  Zap,
} from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '#/components/ui/tabs';
import {
  BillingOverview,
  BillingTopups,
  CheckoutRecovery,
  CustomPlanConfigurator,
  ResourceMetersOverview,
  type CustomPlanQuote,
} from '#/features/billing';
import { useGetPublicSettings } from '#/features/settings';
import { useAuthStore } from '#/stores/auth-store';
import { useBillingConfig, useCreateCheckoutSession } from '#/hooks/use-billing';
import { billingService } from '#/services/billing';

export const Route = createFileRoute('/_dashboard/billing')({
  component: BillingPage,
});

function BillingPage() {
  const user = useAuthStore((s) => s.user);
  const { data: publicRes } = useGetPublicSettings();
  const isCloud = publicRes?.data?.cloudMode ?? true;
  const { data: config } = useBillingConfig();
  const { mutateAsync: createCheckoutSession, isPending } = useCreateCheckoutSession();
  const [activeTab, setActiveTab] = useState<'overview' | 'plans' | 'usage' | 'topups' | 'invoices'>('overview');
  const [loadingPortal, setLoadingPortal] = useState(false);

  const handleUpgradeToPro = async () => {
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

  const handleSelectCustomPlan = async (quote: CustomPlanQuote) => {
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

  return (
    <div className="space-y-6">
      {/* Header matching OpenShip */}
      <header className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex items-center gap-3.5">
          <span className="flex size-11 shrink-0 items-center justify-center rounded-2xl border border-border/60 bg-card text-primary shadow-xs">
            <CreditCard className="size-6" />
          </span>
          <div>
            <h1 className="text-2xl font-semibold tracking-tight text-foreground">Billing & Plans</h1>
            <p className="mt-0.5 text-sm text-muted-foreground">
              Manage your subscription, compute plans, and prepaid credit balance.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            disabled={loadingPortal}
            onClick={handleOpenPortal}
            className="gap-1.5 rounded-xl"
          >
            <ExternalLink className="size-3.5" />
            Customer Portal
          </Button>
          <Button
            size="sm"
            onClick={handleUpgradeToPro}
            className="gap-1.5 rounded-xl"
          >
            <Sparkles className="size-3.5" />
            Upgrade Plan
          </Button>
        </div>
      </header>

      <CheckoutRecovery />

      {/* Tabs matching OpenShip structure */}
      <Tabs
        value={activeTab}
        onValueChange={(val) => setActiveTab(val as typeof activeTab)}
        className="space-y-6"
      >
        <TabsList className="rounded-xl">
          <TabsTrigger value="overview" className="gap-2 rounded-lg">
            <Layers className="size-4" />
            Overview
          </TabsTrigger>
          <TabsTrigger value="plans" className="gap-2 rounded-lg">
            <Sparkles className="size-4" />
            Plans & Pricing
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

        <TabsContent value="overview" className="space-y-6">
          <BillingOverview
            planType={user?.planType ?? 'hobby'}
            onUpgradeClick={handleUpgradeToPro}
            onPortalClick={handleOpenPortal}
          />
          <ResourceMetersOverview />
        </TabsContent>

        <TabsContent value="plans" className="space-y-6">
          <CustomPlanConfigurator
            busy={isPending}
            onSelectPlan={handleSelectCustomPlan}
          />
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
                <h3 className="font-semibold text-foreground text-base">Payment Method & Invoices</h3>
                <p className="mt-1 text-muted-foreground text-xs leading-relaxed max-w-lg">
                  Access your Stripe customer billing portal to download official tax invoices, update credit cards, and manage enterprise billing details.
                </p>
              </div>
              <Button
                variant="outline"
                disabled={loadingPortal}
                onClick={handleOpenPortal}
                className="gap-1.5 shrink-0"
              >
                <ExternalLink className="size-4" />
                Open Stripe Customer Portal
              </Button>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}

