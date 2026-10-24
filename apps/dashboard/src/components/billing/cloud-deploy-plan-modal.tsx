import { Check, X, Zap } from 'lucide-react';
import { useState } from 'react';
import { Button } from '@codedock/ui/button';
import { cn } from '#/lib/utils';

export interface CloudDeployRestriction {
  code: string;
  reason?: 'resource-tier' | 'build-concurrency' | string;
  message?: string;
}

export interface CloudDeployPlanModalProps {
  workspaceId?: string;
  restriction?: CloudDeployRestriction;
  onClose: () => void;
  onSelectPlan?: (planId: string) => Promise<void> | void;
  onRetry?: () => Promise<void> | void;
  className?: string;
}

interface PlanOption {
  id: string;
  name: string;
  price: string;
  concurrency: string;
  resources: string;
  features: string[];
  recommended?: boolean;
}

const PLANS: PlanOption[] = [
  {
    id: 'starter',
    name: 'Starter',
    price: '$0 / mo',
    concurrency: '1 concurrent build',
    resources: '2 vCPU / 4 GB RAM',
    features: ['Community support', 'Standard queue priority', 'Unlimited self-hosted deploys'],
  },
  {
    id: 'pro',
    name: 'Pro',
    price: '$29 / mo',
    concurrency: '3 concurrent builds',
    resources: '4 vCPU / 8 GB RAM',
    features: ['Priority queue access', 'Fast build caches', 'Custom domains & SSL'],
    recommended: true,
  },
  {
    id: 'scale',
    name: 'Scale',
    price: '$99 / mo',
    concurrency: '10 concurrent builds',
    resources: '16 vCPU / 32 GB RAM',
    features: ['Dedicated build workers', 'Advanced telemetry', 'SLA & 24/7 priority support'],
  },
];

export function CloudDeployPlanModal({
  restriction,
  onClose,
  onSelectPlan,
  onRetry,
  className,
}: CloudDeployPlanModalProps) {
  const [selectedPlan, setSelectedPlan] = useState<string>('pro');
  const [submitting, setSubmitting] = useState(false);

  async function handleUpgrade() {
    if (!onSelectPlan) {
      onClose();
      return;
    }
    try {
      setSubmitting(true);
      await onSelectPlan(selectedPlan);
      if (onRetry) {
        await onRetry();
      }
      onClose();
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="plan-modal-title"
      className="fixed inset-0 z-50 flex items-center justify-center bg-background/80 p-4 backdrop-blur-sm"
    >
      <div
        className={cn(
          'relative w-full max-w-3xl rounded-2xl border border-border/60 bg-card p-6 shadow-2xl transition-all',
          className
        )}
      >
        <button
          onClick={onClose}
          aria-label="Close dialog"
          className="absolute top-4 right-4 rounded-lg p-1.5 text-muted-foreground hover:bg-muted hover:text-foreground"
        >
          <X className="size-4" />
        </button>

        <div className="text-left">
          <div className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-2.5 py-0.5 font-medium text-primary text-xs">
            <Zap className="size-3" />
            <span>Plan Upgrade</span>
          </div>
          <h2 id="plan-modal-title" className="mt-2 font-bold text-foreground text-xl">
            Choose your deployment capacity
          </h2>
          {restriction?.message ? (
            <p className="mt-1 text-muted-foreground text-sm">{restriction.message}</p>
          ) : (
            <p className="mt-1 text-muted-foreground text-sm">
              Upgrade your cloud workspace to unlock higher build concurrency and dedicated resource
              tiers.
            </p>
          )}
        </div>

        <div className="mt-6 grid gap-4 sm:grid-cols-3">
          {PLANS.map((plan) => {
            const isSelected = selectedPlan === plan.id;
            return (
              <div
                key={plan.id}
                onClick={() => setSelectedPlan(plan.id)}
                className={cn(
                  'cursor-pointer rounded-xl border p-4 text-left transition-all',
                  isSelected
                    ? 'border-primary bg-primary/5 ring-2 ring-primary/20'
                    : 'border-border/60 hover:border-border hover:bg-muted/30'
                )}
              >
                <div className="flex items-center justify-between">
                  <span className="font-semibold text-foreground text-sm">{plan.name}</span>
                  {plan.recommended && (
                    <span className="rounded-md bg-primary/15 px-1.5 py-0.5 font-semibold text-[10px] text-primary">
                      Popular
                    </span>
                  )}
                </div>
                <p className="mt-2 font-bold text-foreground text-lg">{plan.price}</p>
                <div className="mt-3 space-y-1.5 text-muted-foreground text-xs">
                  <p className="font-medium text-foreground">{plan.concurrency}</p>
                  <p>{plan.resources}</p>
                </div>
                <div className="mt-4 space-y-1.5 border-border/40 border-t pt-3 text-xs">
                  {plan.features.map((feature, i) => (
                    <div key={i} className="flex items-center gap-1.5 text-muted-foreground">
                      <Check className="size-3 text-primary" />
                      <span>{feature}</span>
                    </div>
                  ))}
                </div>
              </div>
            );
          })}
        </div>

        <div className="mt-6 flex items-center justify-end gap-3 border-border/60 border-t pt-4">
          <Button variant="outline" size="sm" onClick={onClose}>
            Cancel
          </Button>
          <Button size="sm" onClick={handleUpgrade} disabled={submitting}>
            {submitting ? 'Updating plan...' : 'Select plan'}
          </Button>
        </div>
      </div>
    </div>
  );
}
