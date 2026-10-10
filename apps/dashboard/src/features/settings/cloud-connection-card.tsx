import { Check, Cloud, Database, GitBranch, Globe, Rocket } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Button } from '#/components/ui/button';

interface CloudConnectionCardProps {
  isCloudMode: boolean;
  baseDomain?: string;
}

export function CloudConnectionCard({
  isCloudMode,
  baseDomain = 'codedock.run',
}: CloudConnectionCardProps) {
  const [connected, setConnected] = useState(isCloudMode);
  const [busy, setBusy] = useState(false);

  const perks = [
    { icon: Globe, title: `Free managed domains (*.${baseDomain})` },
    { icon: Rocket, title: 'Edge builds and zero-downtime routing' },
    { icon: Database, title: 'Integrated managed databases and storage' },
    { icon: GitBranch, title: 'Automatic Git webhook synchronization' },
  ];

  const handleToggle = () => {
    setBusy(true);
    setTimeout(() => {
      setConnected((prev) => !prev);
      setBusy(false);
      toast.success(connected ? 'Disconnected from Cloud' : 'Connected to Codedock Cloud');
    }, 400);
  };

  return (
    <section className="rounded-2xl border border-border/60 bg-card p-5">
      <div className="mb-5 flex items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-violet-500/10 text-violet-500">
            <Cloud className="h-4 w-4" />
          </div>
          <div>
            <h2 className="font-semibold text-foreground text-sm">Codedock Cloud Integration</h2>
            <p className="text-muted-foreground text-xs">
              Connect this instance to managed Cloud edge capabilities and domain routing.
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {connected ? (
            <span className="inline-flex items-center gap-1.5 rounded-full bg-emerald-500/10 px-2.5 py-1 font-semibold text-emerald-500 text-xs ring-1 ring-emerald-500/20">
              <Check className="size-3" />
              Connected
            </span>
          ) : (
            <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2.5 py-1 font-medium text-muted-foreground text-xs">
              Disconnected
            </span>
          )}
        </div>
      </div>

      <div className="grid gap-3 sm:grid-cols-2">
        {perks.map((perk) => (
          <div
            key={perk.title}
            className="flex items-center gap-3 rounded-xl border border-border/50 bg-muted/20 p-3.5"
          >
            <div className="flex size-7 shrink-0 items-center justify-center rounded-lg bg-primary/10 text-primary">
              <perk.icon className="size-3.5" />
            </div>
            <span className="font-medium text-foreground text-xs">{perk.title}</span>
          </div>
        ))}
      </div>

      <div className="mt-6 flex flex-wrap items-center justify-between gap-4 border-border/60 border-t pt-5">
        <div className="text-muted-foreground text-xs">
          {connected
            ? 'Active cloud credentials provision edge routing and TLS certificates automatically.'
            : 'Connecting enables managed wildcards, automated certificate provisioning, and cloud backups.'}
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant={connected ? 'outline' : 'default'}
            size="sm"
            onClick={handleToggle}
            disabled={busy}
          >
            {connected ? 'Disconnect' : 'Connect to Codedock Cloud'}
          </Button>
        </div>
      </div>
    </section>
  );
}
