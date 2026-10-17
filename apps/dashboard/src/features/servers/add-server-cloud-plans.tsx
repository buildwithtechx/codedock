import { Check, Flame, Layers, Sparkles, Zap } from 'lucide-react';
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
    icon: Flame,
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
    <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {CLOUD_PLANS.map((plan) => {
        const Icon = plan.icon;
        return (
          <div
            key={plan.id}
            className={`flex flex-col justify-between rounded-2xl border p-5 ${
              plan.popular
                ? 'border-border bg-card shadow-sm ring-1 ring-border'
                : 'border-border/60 bg-card/60'
            }`}
          >
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <span className="flex size-9 items-center justify-center rounded-xl bg-muted/80 text-foreground">
                  <Icon className="size-4" />
                </span>
                {plan.popular && (
                  <span className="rounded-full bg-foreground/10 px-2.5 py-0.5 font-medium text-[11px] text-foreground">
                    Most popular
                  </span>
                )}
              </div>
              <div>
                <h3 className="font-semibold text-foreground text-lg">{plan.name}</h3>
                <div className="mt-1 flex items-baseline gap-1">
                  <span className="font-bold text-2xl text-foreground tracking-tight">
                    {plan.price}
                  </span>
                  <span className="text-muted-foreground text-xs">/mo</span>
                </div>
                <p className="text-[11px] text-muted-foreground">billed monthly, cancel anytime</p>
                <p className="mt-3 text-muted-foreground text-xs leading-relaxed">
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
                className="w-full justify-center rounded-xl font-medium text-xs"
              >
                Choose {plan.name} →
              </Button>

              <div className="space-y-3 pt-2">
                <p className="font-medium text-[11px] text-muted-foreground uppercase tracking-wider">
                  Plan capacity
                </p>
                <div className="grid grid-cols-3 gap-1 rounded-xl bg-muted/40 p-2.5 text-center">
                  <div>
                    <div className="font-semibold text-foreground text-xs">{plan.vcpu}</div>
                    <div className="text-[10px] text-muted-foreground">vCPU</div>
                  </div>
                  <div>
                    <div className="font-semibold text-foreground text-xs">{plan.ram}</div>
                    <div className="text-[10px] text-muted-foreground">RAM</div>
                  </div>
                  <div>
                    <div className="font-semibold text-foreground text-xs">{plan.storage}</div>
                    <div className="text-[10px] text-muted-foreground">Storage</div>
                  </div>
                </div>

                <div className="space-y-2 text-xs">
                  <div className="flex justify-between text-muted-foreground">
                    <span>Projects</span>
                    <span className="font-medium text-foreground">{plan.projects}</span>
                  </div>
                  <div className="flex justify-between text-muted-foreground">
                    <span>Running services</span>
                    <span className="font-medium text-foreground">{plan.services}</span>
                  </div>
                  <div className="flex justify-between text-muted-foreground">
                    <span>Concurrent builds</span>
                    <span className="font-medium text-foreground">{plan.builds}</span>
                  </div>
                  <div className="flex items-center gap-1.5 pt-1 text-foreground text-xs">
                    <Check className="size-3 text-emerald-500" />
                    <span>{plan.support}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}
