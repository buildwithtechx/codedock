import { ArrowRight, Building2, Check, ChevronDown, Flame, Layers, Rocket, Sparkles, Zap } from 'lucide-react';
import { Button } from '@codedock/ui/button';

export interface CloudPlan {
  id: string;
  name: string;
  price: string;
  popular?: boolean;
  icon: typeof Zap;
  description: string;
  vcpu: number;
  ram: string;
  storage: string;
  projects: number;
  services: number;
  builds: number;
  support: string;
}

export const CLOUD_PLANS: CloudPlan[] = [
  {
    id: 'hobby',
    name: 'Hobby',
    price: '$5',
    icon: Zap,
    description: 'A small managed server for personal projects and everyday apps.',
    vcpu: 1,
    ram: '2 GB',
    storage: '40 GB',
    projects: 3,
    services: 4,
    builds: 1,
    support: 'Community support',
  },
  {
    id: 'starter',
    name: 'Starter',
    price: '$20',
    popular: true,
    icon: Rocket,
    description: 'For your first real app — containers, databases, the works.',
    vcpu: 2,
    ram: '8 GB',
    storage: '128 GB',
    projects: 10,
    services: 15,
    builds: 4,
    support: 'Email support',
  },
  {
    id: 'pro',
    name: 'Pro',
    price: '$39',
    icon: Sparkles,
    description: 'For production apps that need room to grow.',
    vcpu: 4,
    ram: '16 GB',
    storage: '250 GB',
    projects: 25,
    services: 40,
    builds: 8,
    support: 'Priority support',
  },
  {
    id: 'scale',
    name: 'Scale',
    price: '$99',
    icon: Layers,
    description: 'For teams running serious traffic on shared infrastructure.',
    vcpu: 8,
    ram: '32 GB',
    storage: '600 GB',
    projects: 100,
    services: 100,
    builds: 20,
    support: 'Priority support',
  },
];

interface AddServerCloudPlansProps {
  onSelectPlan?: (plan: CloudPlan) => void;
  onClose: () => void;
}

export function AddServerCloudPlans({ onSelectPlan, onClose }: AddServerCloudPlansProps) {
  return (
    <div className="space-y-6">
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {CLOUD_PLANS.map((plan) => {
          const Icon = plan.icon;
          return (
            <article
              key={plan.id}
              className={`flex flex-col justify-between rounded-2xl border p-5 ${
                plan.popular
                  ? 'border-border bg-card shadow-sm ring-1 ring-border/80'
                  : 'border-border/60 bg-card'
              }`}
            >
              <div className="space-y-4">
                <div className="flex items-center gap-3">
                  <div
                    className={`flex size-10 shrink-0 items-center justify-center rounded-xl ${
                      plan.popular ? 'bg-primary/10 text-primary' : 'bg-muted/60 text-foreground'
                    }`}
                  >
                    <Icon className="size-5" />
                  </div>
                  <h3 className="font-semibold text-foreground text-lg tracking-tight">
                    {plan.name}
                  </h3>
                  {plan.popular && (
                    <span className="ms-auto rounded-full bg-foreground px-2.5 py-0.5 font-semibold text-[11px] text-background">
                      Most popular
                    </span>
                  )}
                </div>

                <div>
                  <div className="flex items-baseline gap-1">
                    <span className="font-medium text-3xl text-foreground tracking-tight tabular-nums">
                      {plan.price}
                    </span>
                    <span className="font-medium text-muted-foreground text-sm">/mo</span>
                  </div>
                  <p className="mt-0.5 text-muted-foreground text-xs">billed monthly, cancel anytime</p>
                  <p className="mt-2 min-h-[40px] text-muted-foreground text-sm leading-5">
                    {plan.description}
                  </p>
                </div>

                <Button
                  type="button"
                  variant={plan.popular ? 'default' : 'secondary'}
                  onClick={() => {
                    onSelectPlan?.(plan);
                    onClose();
                  }}
                  className={`h-10 w-full rounded-xl text-sm ${
                    plan.popular
                      ? 'bg-foreground text-background font-semibold hover:bg-foreground/90'
                      : 'border border-border/60 bg-muted/30 font-medium text-foreground hover:bg-muted/60'
                  }`}
                >
                  Choose {plan.name} →
                </Button>

                <div className="space-y-3 pt-2">
                  <p className="font-medium text-muted-foreground text-xs">Plan capacity</p>
                  <dl className="grid grid-cols-3 gap-2 rounded-xl bg-muted/40 p-3 text-center">
                    <div>
                      <dd className="font-semibold text-foreground text-base tabular-nums">{plan.vcpu}</dd>
                      <dt className="text-muted-foreground text-xs">vCPU</dt>
                    </div>
                    <div>
                      <dd className="font-semibold text-foreground text-base tabular-nums">{plan.ram}</dd>
                      <dt className="text-muted-foreground text-xs">RAM</dt>
                    </div>
                    <div>
                      <dd className="font-semibold text-foreground text-base tabular-nums">{plan.storage}</dd>
                      <dt className="text-muted-foreground text-xs">Storage</dt>
                    </div>
                  </dl>

                  <div className="space-y-2 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Projects</span>
                      <span className="font-medium text-foreground tabular-nums">{plan.projects}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Running services</span>
                      <span className="font-medium text-foreground tabular-nums">{plan.services}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Concurrent builds</span>
                      <span className="font-medium text-foreground tabular-nums">{plan.builds}</span>
                    </div>
                    <div className="flex items-center gap-2 pt-1 text-muted-foreground">
                      <Check className="size-3.5 shrink-0" />
                      <span>{plan.support}</span>
                    </div>
                  </div>
                </div>
              </div>
            </article>
          );
        })}
      </div>

      <div className="border-border/50 border-t pt-5">
        <h4 className="mb-3 font-semibold text-foreground text-sm">In every plan</h4>
        <div className="grid grid-cols-1 gap-x-6 gap-y-2.5 text-muted-foreground text-sm md:grid-cols-3">
          <div className="space-y-2.5">
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>CPU, RAM and storage covered 24/7</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>Full available capacity for any service</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>Live build and deploy logs</span>
            </div>
          </div>
          <div className="space-y-2.5">
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>No compute credits or top-ups</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>Deploy from GitHub, GitLab, or any Git remote</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>Audit log of every change</span>
            </div>
          </div>
          <div className="space-y-2.5">
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>Shared capacity for apps and builds</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>Unlimited custom domains, TLS included</span>
            </div>
            <div className="flex items-center gap-2">
              <Check className="size-4 shrink-0 text-muted-foreground" />
              <span>Unlimited team members — never pay per seat</span>
            </div>
          </div>
        </div>
        <button
          type="button"
          className="mt-3 flex items-center gap-1 text-muted-foreground text-xs hover:text-foreground"
        >
          <span>Usage, transfer and storage terms</span>
          <ChevronDown className="size-3.5" />
        </button>
      </div>

      <div className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-border/50 bg-card p-4">
        <div className="flex items-center gap-3.5">
          <div className="flex size-10 shrink-0 items-center justify-center rounded-xl bg-muted/60">
            <Building2 className="size-5 text-muted-foreground" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <h4 className="font-semibold text-foreground text-sm">Enterprise</h4>
              <span className="rounded-full bg-muted px-2 py-0.5 font-medium text-[11px] text-muted-foreground">
                Custom
              </span>
            </div>
            <p className="mt-0.5 text-muted-foreground text-xs">
              Custom limits, dedicated infrastructure, and an SLA.
            </p>
          </div>
        </div>
        <Button variant="secondary" className="gap-1.5 rounded-xl text-xs">
          Talk to sales <ArrowRight className="size-3.5" />
        </Button>
      </div>
    </div>
  );
}
