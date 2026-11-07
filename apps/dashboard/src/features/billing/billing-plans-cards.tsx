import { ArrowRight, Check, Gem, Layers, Rocket, Zap } from 'lucide-react';
import { Button } from '@codedock/ui/button';

export interface PlanTier {
  id: string;
  name: string;
  badge?: string;
  price: string;
  period: string;
  description: string;
  icon: typeof Zap;
  isPopular?: boolean;
  capacity: {
    vcpu: string;
    ram: string;
    storage: string;
  };
  metrics: {
    projects: number;
    services: number;
    builds: number;
  };
  support: string;
}

export const STANDARD_PLANS: PlanTier[] = [
  {
    id: 'hobby',
    name: 'Hobby',
    price: '$5',
    period: '/mo',
    description: 'A small managed server for personal projects and everyday apps.',
    icon: Zap,
    capacity: { vcpu: '1', ram: '2 GB', storage: '40 GB' },
    metrics: { projects: 3, services: 4, builds: 1 },
    support: 'Community support',
  },
  {
    id: 'starter',
    name: 'Starter',
    badge: 'Most popular',
    isPopular: true,
    price: '$20',
    period: '/mo',
    description: 'For your first real app — containers, databases, the works.',
    icon: Rocket,
    capacity: { vcpu: '2', ram: '8 GB', storage: '128 GB' },
    metrics: { projects: 10, services: 15, builds: 4 },
    support: 'Email support',
  },
  {
    id: 'pro',
    name: 'Pro',
    price: '$39',
    period: '/mo',
    description: 'For production apps that need room to grow.',
    icon: Gem,
    capacity: { vcpu: '4', ram: '16 GB', storage: '250 GB' },
    metrics: { projects: 25, services: 40, builds: 8 },
    support: 'Priority support',
  },
  {
    id: 'scale',
    name: 'Scale',
    price: '$99',
    period: '/mo',
    description: 'For teams running serious traffic on shared infrastructure.',
    icon: Layers,
    capacity: { vcpu: '8', ram: '32 GB', storage: '600 GB' },
    metrics: { projects: 100, services: 100, builds: 20 },
    support: 'Priority support',
  },
];

interface BillingPlansCardsProps {
  onSelectPlan: (planId: string) => void;
  busy?: boolean;
}

export function BillingPlansCards({ onSelectPlan, busy = false }: BillingPlansCardsProps) {
  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
      {STANDARD_PLANS.map((plan) => {
        const Icon = plan.icon;
        return (
          <div
            key={plan.id}
            className={`relative flex flex-col justify-between rounded-2xl border p-5 transition-all ${
              plan.isPopular
                ? 'border-primary/50 bg-card shadow-sm ring-1 ring-primary/20'
                : 'border-border/60 bg-card'
            }`}
          >
            <div>
              <div className="flex items-center justify-between gap-2">
                <div className="flex items-center gap-2.5">
                  <span className="flex size-8 items-center justify-center rounded-xl bg-muted/60 text-foreground">
                    <Icon className="size-4" />
                  </span>
                  <h3 className="font-medium text-base text-foreground">{plan.name}</h3>
                </div>
                {plan.badge && (
                  <span className="rounded-full border border-border/80 bg-muted/80 px-2.5 py-0.5 font-medium text-[11px] text-foreground">
                    {plan.badge}
                  </span>
                )}
              </div>

              <div className="mt-4">
                <div className="flex items-baseline gap-1">
                  <span className="font-semibold text-2xl text-foreground tracking-tight">
                    {plan.price}
                  </span>
                  <span className="text-muted-foreground text-xs">{plan.period}</span>
                </div>
                <p className="mt-0.5 text-muted-foreground text-xs">
                  billed monthly, cancel anytime
                </p>
              </div>

              <p className="mt-3 min-h-[36px] text-muted-foreground text-xs leading-relaxed">
                {plan.description}
              </p>

              <Button
                type="button"
                variant={plan.isPopular ? 'default' : 'secondary'}
                disabled={busy}
                onClick={() => onSelectPlan(plan.id)}
                className="mt-4 w-full justify-center gap-1.5 font-medium text-xs"
              >
                Choose {plan.name}
                <ArrowRight className="size-3.5" />
              </Button>

              <div className="mt-5 border-border/40 border-t pt-4">
                <p className="font-medium text-[11px] text-muted-foreground uppercase tracking-wider">
                  Plan capacity
                </p>
                <div className="mt-2.5 grid grid-cols-3 gap-2 text-center">
                  <div className="rounded-xl bg-muted/30 p-2">
                    <p className="font-semibold text-foreground text-xs">{plan.capacity.vcpu}</p>
                    <p className="text-[10px] text-muted-foreground">vCPU</p>
                  </div>
                  <div className="rounded-xl bg-muted/30 p-2">
                    <p className="font-semibold text-foreground text-xs">{plan.capacity.ram}</p>
                    <p className="text-[10px] text-muted-foreground">RAM</p>
                  </div>
                  <div className="rounded-xl bg-muted/30 p-2">
                    <p className="font-semibold text-foreground text-xs">{plan.capacity.storage}</p>
                    <p className="text-[10px] text-muted-foreground">Storage</p>
                  </div>
                </div>
              </div>

              <dl className="mt-4 space-y-2 border-border/40 border-t pt-4 text-xs">
                <div className="flex items-center justify-between">
                  <dt className="text-muted-foreground">Projects</dt>
                  <dd className="font-medium text-foreground">{plan.metrics.projects}</dd>
                </div>
                <div className="flex items-center justify-between">
                  <dt className="text-muted-foreground">Running services</dt>
                  <dd className="font-medium text-foreground">{plan.metrics.services}</dd>
                </div>
                <div className="flex items-center justify-between">
                  <dt className="text-muted-foreground">Concurrent builds</dt>
                  <dd className="font-medium text-foreground">{plan.metrics.builds}</dd>
                </div>
              </dl>
            </div>

            <div className="mt-4 flex items-center gap-2 border-border/40 border-t pt-3 text-muted-foreground text-xs">
              <Check className="size-3.5 text-emerald-500" />
              <span>{plan.support}</span>
            </div>
          </div>
        );
      })}
    </div>
  );
}
