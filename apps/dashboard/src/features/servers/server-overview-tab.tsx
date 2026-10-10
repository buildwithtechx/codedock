import type { LucideIcon } from 'lucide-react';
import { Activity, Cpu, HardDrive, MemoryStick, WifiOff } from 'lucide-react';
import { Button } from '#/components/ui/button';
import { MetricChartCard } from '#/features/services/metric-chart-card';
import { parseServerMetrics, type Server } from '#/interfaces/server';
import { cn } from '#/lib/utils';
import {
  formatBytes,
  type LiveMetricsSnapshot,
  snapshotFromStored,
  useServerLiveMetrics,
} from './server-live-metrics';

function UsageBar({ pct }: { pct: number }) {
  const tone = pct >= 90 ? 'bg-destructive' : pct >= 70 ? 'bg-warning' : 'bg-foreground/60';
  return (
    <div className="mt-3 h-1.5 overflow-hidden rounded-full bg-muted">
      <div
        className={cn('h-full rounded-full transition-all duration-700 ease-out', tone)}
        style={{ width: `${Math.min(pct, 100)}%` }}
      />
    </div>
  );
}

function StatCard({
  icon: Icon,
  label,
  value,
  sub,
  pct,
}: {
  icon: LucideIcon;
  label: string;
  value: string;
  sub?: string;
  pct?: number;
}) {
  return (
    <div className="min-w-0 rounded-2xl bg-card p-5">
      <div className="mb-3 flex items-center gap-2">
        <Icon className="size-4 text-muted-foreground" />
        <span className="font-medium text-muted-foreground text-xs uppercase tracking-wide">
          {label}
        </span>
      </div>
      <p className="font-semibold text-2xl text-foreground tabular-nums tracking-tight">{value}</p>
      {sub && <p className="mt-1 text-muted-foreground text-xs tabular-nums">{sub}</p>}
      {pct != null && <UsageBar pct={pct} />}
    </div>
  );
}

function formatHeartbeat(lastSeenAt?: string, isLocal?: boolean): string {
  if (isLocal) return 'Just now';
  if (!lastSeenAt) return 'Never';
  const seen = new Date(lastSeenAt);
  if (Number.isNaN(seen.getTime())) return 'Never';
  const minutes = Math.max(0, Math.round((Date.now() - seen.getTime()) / 60000));
  if (minutes < 1) return 'Just now';
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  return seen.toLocaleDateString();
}

export function ServerOverviewTab({ server }: { server: Server }) {
  const stored = snapshotFromStored(parseServerMetrics(server.metrics));
  const live = useServerLiveMetrics(server.id, stored);
  const stats: LiveMetricsSnapshot | null = live.snapshot ?? stored;

  return (
    <div className="space-y-6">
      {!live.connected && (
        <div
          role="alert"
          className="flex flex-wrap items-center justify-between gap-3 rounded-xl bg-muted/40 p-4 text-sm"
        >
          <p className="inline-flex min-w-0 items-center gap-2 break-words text-muted-foreground">
            <WifiOff className="size-4 shrink-0" />
            Live telemetry is disconnected. Showing the last reported snapshot.
          </p>
          <Button size="sm" variant="secondary" onClick={live.reconnect}>
            Retry
          </Button>
        </div>
      )}

      <div className="grid grid-cols-2 gap-4 xl:grid-cols-4">
        <StatCard
          icon={Cpu}
          label="CPU"
          value={stats ? `${stats.cpu.toFixed(1)}%` : '-'}
          sub={live.connected ? 'Live from node' : 'Last reported'}
          pct={stats?.cpu}
        />
        <StatCard
          icon={MemoryStick}
          label="Memory"
          value={stats ? `${stats.memory.toFixed(1)}%` : '-'}
          sub={
            stats
              ? `${formatBytes(stats.memoryUsedBytes)} of ${formatBytes(stats.memoryLimitBytes)}`
              : undefined
          }
          pct={stats?.memory}
        />
        <StatCard
          icon={HardDrive}
          label="Disk"
          value={stats ? `${stats.disk.toFixed(1)}%` : '-'}
          sub={
            stats
              ? `${formatBytes(stats.diskUsedBytes)} of ${formatBytes(stats.diskTotalBytes)}`
              : undefined
          }
          pct={stats?.disk}
        />
        <StatCard
          icon={Activity}
          label="Heartbeat"
          value={formatHeartbeat(server.lastSeenAt, server.isLocal)}
          sub={live.connected ? 'Stream connected' : 'Stream idle'}
        />
      </div>

      <div className="grid grid-cols-1 gap-6 xl:grid-cols-3">
        <MetricChartCard
          title="CPU Usage"
          icon={<Cpu className="h-4 w-4" />}
          badge={live.connected ? 'Live' : 'Idle'}
          data={live.cpuSeries}
          isLoading={false}
          color="#8b5cf6"
          formatY={(v) => `${Math.round(v)}%`}
          formatTooltip={(v) => `${v.toFixed(2)}%`}
        />
        <MetricChartCard
          title="Memory Usage"
          icon={<MemoryStick className="h-4 w-4" />}
          badge={live.connected ? 'Live' : 'Idle'}
          data={live.memorySeries}
          isLoading={false}
          color="#3b82f6"
          formatY={(v) => `${Math.round(v)}%`}
          formatTooltip={(v) => `${v.toFixed(2)}%`}
        />
        <MetricChartCard
          title="Disk Usage"
          icon={<HardDrive className="h-4 w-4" />}
          badge={live.connected ? 'Live' : 'Idle'}
          data={live.diskSeries}
          isLoading={false}
          color="#10b981"
          formatY={(v) => `${Math.round(v)}%`}
          formatTooltip={(v) => `${v.toFixed(2)}%`}
        />
      </div>
    </div>
  );
}
