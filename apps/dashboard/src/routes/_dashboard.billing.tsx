import { createFileRoute } from '@tanstack/react-router';
import { PageHeader } from '#/components/layout/page-header';
import { BillingSection } from '#/features/profile/billing-section';

export const Route = createFileRoute('/_dashboard/billing')({
  component: BillingPage,
});

function BillingPage() {
  return (
    <div className="space-y-6">
      <PageHeader
        title="Billing & Subscription"
        description="Manage your subscription, compute plans, and payment details."
      />
      <BillingSection />
    </div>
  );
}
