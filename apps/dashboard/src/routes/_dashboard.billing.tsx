import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '#/components/layout/page-header';
import {
  BillingOverview,
  BillingTopups,
  CustomPlanConfigurator,
  type CustomPlanQuote,
} from '#/features/billing';
import { useAuthStore } from '#/stores/auth-store';
import { useBillingConfig, useCreateCheckoutSession } from '#/hooks/use-billing';

export const Route = createFileRoute('/_dashboard/billing')({
  component: BillingPage,
});

function BillingPage() {
  const user = useAuthStore((s) => s.user);
  const { data: config } = useBillingConfig();
  const { mutateAsync: createCheckoutSession, isPending } = useCreateCheckoutSession();

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

  return (
    <div className="space-y-6">
      <PageHeader
        title="Billing & Subscription"
        description="Manage your subscription, compute plans, and prepaid credit balance."
      />
      <BillingOverview
        planType={user?.planType ?? 'hobby'}
        onUpgradeClick={handleUpgradeToPro}
      />
      <CustomPlanConfigurator
        busy={isPending}
        onSelectPlan={handleSelectCustomPlan}
      />
      <BillingTopups />
    </div>
  );
}

