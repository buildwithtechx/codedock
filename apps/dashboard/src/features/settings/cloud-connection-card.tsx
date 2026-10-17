import { Check, Cloud, ExternalLink, Globe, Loader2, LogOut, Rocket, Shield } from 'lucide-react';
import { useState } from 'react';
import { Button } from '#/components/ui/button';
import { cn } from '#/lib/utils';

export interface CloudUser {
  name?: string;
  email?: string;
  image?: string;
}

export interface CloudConnectionCardProps {
  connected?: boolean;
  loading?: boolean;
  connecting?: boolean;
  cloudUser?: CloudUser | null;
  baseDomain?: string;
  onConnect?: () => Promise<void> | void;
  onDisconnect?: () => Promise<void> | void;
  className?: string;
}

export function CloudConnectionCard({
  connected = false,
  loading = false,
  connecting = false,
  cloudUser,
  baseDomain = 'codedock.app',
  onConnect,
  onDisconnect,
  className,
}: CloudConnectionCardProps) {
  const [disconnecting, setDisconnecting] = useState(false);

  const perks = [
    {
      icon: Globe,
      text: `Free wildcard domains on *.${baseDomain}`,
    },
    {
      icon: Rocket,
      text: 'Managed build servers and distributed workers',
    },
    {
      icon: Shield,
      text: 'Automated SSL certificates and edge routing',
    },
  ];

  async function handleDisconnect() {
    if (!window.confirm('Are you sure you want to disconnect from CodeDock Cloud?')) {
      return;
    }
    if (!onDisconnect) return;
    try {
      setDisconnecting(true);
      await onDisconnect();
    } finally {
      setDisconnecting(false);
    }
  }

  return (
    <div className={cn('rounded-2xl border border-border/50 bg-card p-6 text-left', className)}>
      {loading ? (
        <div className="flex items-center justify-center gap-2 py-6">
          <Loader2 className="size-4 animate-spin text-muted-foreground" />
          <p className="text-muted-foreground text-xs">Checking cloud connection...</p>
        </div>
      ) : connected ? (
        <div className="space-y-5">
          <div className="flex items-center gap-3">
            {cloudUser?.image ? (
              <img
                src={cloudUser.image}
                alt=""
                className="size-10 shrink-0 rounded-full ring-1 ring-border"
              />
            ) : (
              <div className="flex size-10 shrink-0 items-center justify-center rounded-full bg-primary/10 text-primary ring-1 ring-border">
                <span className="font-semibold text-sm">
                  {cloudUser?.name?.charAt(0)?.toUpperCase() ||
                    cloudUser?.email?.charAt(0)?.toUpperCase() ||
                    'C'}
                </span>
              </div>
            )}
            <div className="min-w-0 flex-1">
              <p className="truncate font-medium text-foreground text-sm">
                {cloudUser?.name || 'Connected Account'}
              </p>
              <p className="truncate text-muted-foreground text-xs">
                {cloudUser?.email || 'CodeDock Cloud User'}
              </p>
            </div>
            <div className="inline-flex items-center gap-1.5 rounded-full bg-emerald-500/15 px-2.5 py-1 font-semibold text-emerald-600 text-xs dark:text-emerald-400">
              <Check className="size-3" aria-hidden="true" />
              <span>Connected</span>
            </div>
          </div>

          <div className="h-px bg-border/60" />

          <div className="flex items-center justify-between gap-4">
            <p className="text-muted-foreground text-xs leading-relaxed">
              Cloud connection active. Managed edge domains and build infrastructure are enabled.
            </p>
            {onDisconnect && (
              <Button
                variant="ghost"
                size="sm"
                onClick={handleDisconnect}
                disabled={disconnecting}
                className="text-muted-foreground hover:text-destructive"
              >
                {disconnecting ? (
                  <Loader2 className="size-3.5 animate-spin" />
                ) : (
                  <LogOut className="size-3.5" />
                )}
                <span>Disconnect</span>
              </Button>
            )}
          </div>
        </div>
      ) : (
        <div className="flex flex-col items-center text-center">
          <div className="mb-4 flex size-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
            <Cloud className="size-6" />
          </div>

          <h3 className="font-semibold text-foreground text-lg tracking-tight">
            Connect to CodeDock Cloud
          </h3>
          <p className="mt-1.5 max-w-md text-muted-foreground text-sm leading-relaxed">
            Pair your self-hosted instance with CodeDock Cloud to unlock managed domain routing,
            cloud builds, and distributed edge workers.
          </p>

          <ul className="mt-5 grid w-full max-w-lg gap-2 text-start sm:grid-cols-2">
            {perks.map((perk, i) => {
              const Icon = perk.icon;
              return (
                <li key={i} className="flex items-start gap-2.5">
                  <span className="mt-0.5 flex size-5 shrink-0 items-center justify-center rounded-md bg-primary/10 text-primary">
                    <Icon className="size-3" />
                  </span>
                  <span className="text-muted-foreground text-xs leading-snug">{perk.text}</span>
                </li>
              );
            })}
          </ul>

          <div className="mt-6 flex w-full max-w-xs flex-col gap-2.5 sm:flex-row">
            {onConnect && (
              <Button onClick={onConnect} disabled={connecting} className="flex-1">
                {connecting ? (
                  <>
                    <Loader2 className="size-3.5 animate-spin" />
                    <span>Connecting...</span>
                  </>
                ) : (
                  <>
                    <ExternalLink className="size-3.5" />
                    <span>Connect Cloud</span>
                  </>
                )}
              </Button>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
