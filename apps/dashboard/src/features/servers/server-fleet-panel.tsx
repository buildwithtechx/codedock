import { Activity, Layers, MapPin, Server as ServerIcon } from 'lucide-react';
import { Skeleton } from '#/components/ui/skeleton';
import type { Server } from '#/interfaces/server';
import { reachabilityOf } from './server-list-row';

interface ServerFleetPanelProps {
  servers: Server[];
  loading: boolean;
}

export function ServerFleetPanel({ servers, loading }: ServerFleetPanelProps) {
  const online = servers.filter((server) => reachabilityOf(server) === 'online').length;
  const offline = servers.filter((server) => reachabilityOf(server) === 'offline').length;
  const regions = new Set(servers.map((server) => server.region).filter(Boolean)).size;
  const pct = servers.length > 0 ? Math.round((online / servers.length) * 100) : 0;

  const rows = [
    { icon: ServerIcon, label: 'Total servers', value: servers.length },
    { icon: Activity, label: 'Online', value: online },
    {
      icon: Layers,
      label: 'Providers',
      value: new Set(servers.map((s) => s.provider).filter(Boolean)).size,
    },
    ...(regions > 0 ? [{ icon: MapPin, label: 'Regions', value: regions }] : []),
  ];

  return (
    <div className="rounded-2xl bg-card" aria-busy={loading}>
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        <div className="flex h-9 w-9 items-center justify-center rounded-xl bg-muted">
          <Activity className="size-[18px] text-muted-foreground" />
        </div>
        <div>
          <h2 className="font-semibold text-[15px] text-foreground">Quick info</h2>
          <p className="text-muted-foreground text-xs">Fleet overview</p>
        </div>
      </div>
      <div className="space-y-5 p-5">
        <div>
          <div className="flex items-baseline justify-between">
            {loading ? (
              <div
                aria-hidden="true"
                className="flex h-8 items-center gap-2 motion-safe:animate-pulse"
              >
                <Skeleton className="h-6 w-9" />
                <Skeleton className="h-3.5 w-20" />
              </div>
            ) : (
              <div className="flex items-baseline gap-1.5">
                <span className="font-semibold text-2xl text-foreground tabular-nums">
                  {online}
                </span>
                <span className="text-muted-foreground text-sm">/ {servers.length} online</span>
              </div>
            )}
            {loading ? (
              <Skeleton className="h-3 w-8 self-center" />
            ) : (
              <span className="font-medium text-muted-foreground text-xs tabular-nums">{pct}%</span>
            )}
          </div>
          <div className="mt-2.5 h-1.5 w-full overflow-hidden rounded-full bg-muted">
            {!loading && (
              <div
                className="h-full rounded-full bg-success transition-[width] duration-500"
                style={{ width: `${pct}%` }}
              />
            )}
          </div>
          {!loading && offline > 0 && (
            <p className="mt-2 inline-flex items-center gap-1.5 text-destructive text-xs">
              <span className="size-1.5 rounded-full bg-destructive" />
              {offline} offline
            </p>
          )}
        </div>

        <div className="space-y-0.5 border-border/50 border-t pt-4">
          {rows.map((row) => (
            <div key={row.label} className="flex items-center justify-between py-1.5">
              <span className="inline-flex items-center gap-2.5 text-muted-foreground text-sm">
                <row.icon className="size-4 text-muted-foreground/60" />
                {row.label}
              </span>
              {loading ? (
                <Skeleton className="h-4 w-5" />
              ) : (
                <span className="font-medium text-foreground text-sm tabular-nums">
                  {row.value}
                </span>
              )}
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
