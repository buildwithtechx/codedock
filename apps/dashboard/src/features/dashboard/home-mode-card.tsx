import { Link } from '@tanstack/react-router';
import { ArrowRight, Cloud, Server } from 'lucide-react';
import { Badge } from '@codedock/ui/badge';
import { Button } from '@codedock/ui/button';
import { useListServers } from '#/hooks/use-servers';
import { useGetPublicSettings } from '#/hooks/use-settings';
import { useAuthStore } from '#/stores/auth-store';

export function HomeModeCard() {
  const { data: publicSettings } = useGetPublicSettings();
  const { data: servers = [], isLoading: serversLoading } = useListServers();
  const user = useAuthStore((s) => s.user);

  const isCloud = Boolean(publicSettings?.data?.cloudMode);

  if (isCloud) {
    const planType = user?.planType || 'free';
    const planLabel =
      planType === 'pro'
        ? 'Pro Workspace'
        : planType === 'payg'
          ? 'Pay-As-You-Go'
          : 'Starter Cloud';

    return (
      <div className="relative overflow-hidden rounded-2xl border border-primary/20 bg-linear-to-br from-primary/10 via-card to-card p-5 shadow-xs">
        <div className="flex items-center justify-between gap-3">
          <div className="flex items-center gap-2.5">
            <div className="flex size-9 items-center justify-center rounded-xl bg-primary/15 text-primary">
              <Cloud className="size-4" />
            </div>
            <div>
              <p className="font-semibold text-[10px] text-primary uppercase tracking-wider">
                Codedock Cloud
              </p>
              <h3 className="font-semibold text-foreground text-sm">{planLabel}</h3>
            </div>
          </div>
          <Badge
            variant="outline"
            className="border-primary/30 text-[11px] text-primary capitalize"
          >
            {planType}
          </Badge>
        </div>

        <p className="mt-3 text-muted-foreground text-xs leading-relaxed">
          Managed infrastructure with global edge ingress, automated TLS, and continuous
          deployments.
        </p>

        <div className="mt-4 border-border/40 border-t pt-3">
          <Button
            asChild
            variant="secondary"
            size="sm"
            className="h-8 w-full justify-between font-medium text-xs"
          >
            <Link to="/billing">
              <span>Manage Subscription & Quotas</span>
              <ArrowRight className="size-3.5 opacity-70" />
            </Link>
          </Button>
        </div>
      </div>
    );
  }

  // Self-Hosted Mode
  const onlineServers = servers.filter((s) => s.status?.toLowerCase() === 'online');
  const serverCount = servers.length;

  return (
    <div className="relative overflow-hidden rounded-2xl border border-border/60 bg-linear-to-br from-muted/30 via-card to-card p-5 shadow-xs">
      <div className="flex items-center justify-between gap-3">
        <div className="flex items-center gap-2.5">
          <div className="flex size-9 items-center justify-center rounded-xl bg-muted text-foreground">
            <Server className="size-4" />
          </div>
          <div>
            <p className="font-semibold text-[10px] text-muted-foreground uppercase tracking-wider">
              Infrastructure
            </p>
            <h3 className="font-semibold text-foreground text-sm">Self-Hosted Cluster</h3>
          </div>
        </div>
        <div className="flex items-center gap-1.5">
          <span className="size-2 animate-pulse rounded-full bg-emerald-500" />
          <span className="font-medium text-emerald-500 text-xs">Active</span>
        </div>
      </div>

      <div className="mt-4 grid grid-cols-2 gap-2 text-xs">
        <div className="rounded-xl border border-border/40 bg-muted/20 p-2.5">
          <p className="text-[11px] text-muted-foreground">Connected Nodes</p>
          <p className="mt-0.5 font-bold text-foreground text-sm">
            {serversLoading ? '...' : `${onlineServers.length}/${serverCount || 1}`}
          </p>
        </div>
        <div className="rounded-xl border border-border/40 bg-muted/20 p-2.5">
          <p className="text-[11px] text-muted-foreground">Runtime Engine</p>
          <p className="mt-0.5 font-bold text-foreground text-sm">Docker Host</p>
        </div>
      </div>

      <div className="mt-4 border-border/40 border-t pt-3">
        <Button
          asChild
          variant="outline"
          size="sm"
          className="h-8 w-full justify-between font-medium text-xs"
        >
          <Link to="/servers">
            <span>Manage Servers & Nodes</span>
            <ArrowRight className="size-3.5 opacity-70" />
          </Link>
        </Button>
      </div>
    </div>
  );
}
