import { Check, Cloud, Laptop, Server } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';

type PlatformMode = 'self-hosted' | 'cloud' | 'desktop';

const MODES = [
  {
    id: 'self-hosted' as PlatformMode,
    title: 'Self-Hosted Mode',
    description: 'Control plane runs entirely on your local Docker daemon or self-managed VPS.',
    icon: Server,
  },
  {
    id: 'cloud' as PlatformMode,
    title: 'Cloud Mode',
    description:
      'Managed control plane with edge ingress, automatic backups, and high availability.',
    icon: Cloud,
  },
  {
    id: 'desktop' as PlatformMode,
    title: 'Desktop App Shell',
    description: 'Embedded native runner for offline and local development workspaces.',
    icon: Laptop,
  },
];

export function PlatformModeCard({ initialMode = 'self-hosted' }: { initialMode?: PlatformMode }) {
  const [mode, setMode] = useState<PlatformMode>(initialMode);

  const handleSelect = (next: PlatformMode) => {
    setMode(next);
    toast.success(`Platform operating mode updated to ${next}`);
  };

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5">
      <div className="mb-5">
        <h2 className="font-semibold text-foreground text-sm">Deployment & Operating Mode</h2>
        <p className="text-muted-foreground text-xs">
          Select the runtime topology this Codedock control plane targets.
        </p>
      </div>

      <div className="grid gap-3 sm:grid-cols-3">
        {MODES.map((item) => {
          const active = mode === item.id;
          const Icon = item.icon;
          return (
            <button
              key={item.id}
              type="button"
              onClick={() => handleSelect(item.id)}
              className={`flex flex-col gap-3 rounded-xl border p-4 text-left transition-all ${
                active
                  ? 'border-primary/60 bg-primary/5 ring-1 ring-primary/20'
                  : 'border-border/50 bg-muted/15 hover:border-border hover:bg-muted/30'
              }`}
            >
              <div className="flex items-center justify-between">
                <div
                  className={`flex size-8 items-center justify-center rounded-lg ${
                    active ? 'bg-primary/15 text-primary' : 'bg-muted text-muted-foreground'
                  }`}
                >
                  <Icon className="size-4" />
                </div>
                {active && (
                  <span className="flex size-4.5 items-center justify-center rounded-full bg-primary text-primary-foreground">
                    <Check className="size-3" />
                  </span>
                )}
              </div>
              <div>
                <p className="font-semibold text-foreground text-xs">{item.title}</p>
                <p className="mt-1 text-[11px] text-muted-foreground leading-relaxed">
                  {item.description}
                </p>
              </div>
            </button>
          );
        })}
      </div>
    </section>
  );
}
