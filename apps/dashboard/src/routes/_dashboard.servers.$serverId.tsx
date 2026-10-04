import { createFileRoute, Link } from '@tanstack/react-router';
import {
  ArrowLeft,
  Cpu,
  HardDrive,
  Loader2,
  MemoryStick,
  ServerIcon,
  Settings,
  Terminal,
} from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Card } from '#/components/ui/card';
import { env } from '#/env';
import { ServerConnectionCard } from '#/features/servers/server-connection-card';
import { ServerSettingsTab } from '#/features/servers/server-settings-tab';
import { MetricChartCard } from '#/features/services/metric-chart-card';
import { useServer } from '#/hooks/use-servers';
import { useAuthStore } from '#/stores/auth-store';

export const Route = createFileRoute('/_dashboard/servers/$serverId')({
  component: ServerDetailsPage,
});

interface WSMetric {
  time: number;
  cpu: number;
  memory: number;
  disk: number;
}

type ServerTab = 'overview' | 'terminal' | 'settings';

function ServerDetailsPage() {
  const { serverId } = Route.useParams();
  const { data: server, isLoading, refetch } = useServer(serverId);
  const [activeTab, setActiveTab] = useState<ServerTab>('overview');

  const [metrics, setMetrics] = useState<WSMetric[]>([]);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    if (!serverId) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsHost = env.VITE_API_URL.replace(/^http(s?):\/\//, '');
    const wsUrl = `${protocol}//${wsHost}/api/ws/servers/${serverId}/metrics`;

    const token = useAuthStore.getState().token;
    const protocols = token ? ['auth', token] : undefined;
    const socket = new WebSocket(wsUrl, protocols);

    socket.onopen = () => setIsConnected(true);
    socket.onclose = () => setIsConnected(false);

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        const time = Date.now();
        const cpu = data.cpu_usage_percentage || 0;
        const memory =
          data.memory_limit_bytes > 0
            ? (data.memory_usage_bytes / data.memory_limit_bytes) * 100
            : 0;
        const disk =
          data.disk_total_bytes > 0 ? (data.disk_usage_bytes / data.disk_total_bytes) * 100 : 0;

        setMetrics((prev) => {
          const newMetrics = [...prev, { time, cpu, memory, disk }];
          if (newMetrics.length > 30) {
            newMetrics.shift();
          }
          return newMetrics;
        });
      } catch (err) {
        console.error('Failed to parse WS metrics', err);
      }
    };

    return () => socket.close();
  }, [serverId]);

  const cpuData = useMemo(() => metrics.map((m) => ({ time: m.time, value: m.cpu })), [metrics]);
  const memData = useMemo(() => metrics.map((m) => ({ time: m.time, value: m.memory })), [metrics]);
  const diskData = useMemo(() => metrics.map((m) => ({ time: m.time, value: m.disk })), [metrics]);

  if (isLoading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!server) {
    return (
      <div className="flex h-64 flex-col items-center justify-center rounded-xl border border-dashed bg-card/40">
        <ServerIcon className="mb-4 h-8 w-8 text-muted-foreground" />
        <h3 className="font-bold text-lg">Server not found</h3>
        <Link to="/servers" className="mt-4 text-primary text-sm hover:underline">
          Back to Servers
        </Link>
      </div>
    );
  }

  const statusColor = server.status === 'online' ? 'bg-emerald-500' : 'bg-amber-500';

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <div>
          <Link
            to="/servers"
            className="inline-flex items-center gap-1.5 text-muted-foreground text-xs transition-colors hover:text-foreground"
          >
            <ArrowLeft className="h-3.5 w-3.5" />
            Servers
          </Link>
          <div className="mt-2 flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-primary/20 bg-primary/10 text-primary">
              <ServerIcon className="h-5 w-5" />
            </div>
            <div>
              <h1 className="font-bold text-xl">{server.name}</h1>
              <p className="font-mono text-muted-foreground text-xs">{server.ipAddress}</p>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-3">
          <div className="flex items-center gap-2 rounded-full border bg-background px-3 py-1 shadow-sm">
            <div
              className={`h-2.5 w-2.5 rounded-full ${statusColor} ${server.status === 'online' ? 'animate-pulse' : ''}`}
            />
            <span className="font-semibold text-xs uppercase tracking-wider">{server.status}</span>
          </div>
          <div className="flex items-center gap-2 rounded-full border bg-background px-3 py-1 text-xs shadow-sm">
            <span
              className={`h-2 w-2 rounded-full ${isConnected ? 'bg-emerald-500' : 'bg-red-500'}`}
            />
            <span className="text-muted-foreground">
              {isConnected ? 'Telemetry live' : 'Telemetry disconnected'}
            </span>
          </div>
        </div>
      </div>

      <div className="flex items-center gap-2 border-border/60 border-b pb-1">
        <button
          type="button"
          onClick={() => setActiveTab('overview')}
          className={`flex items-center gap-1.5 border-b-2 px-3 py-2 font-medium text-xs transition-colors ${
            activeTab === 'overview'
              ? 'border-primary text-foreground'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Cpu className="h-3.5 w-3.5" />
          Overview
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('terminal')}
          className={`flex items-center gap-1.5 border-b-2 px-3 py-2 font-medium text-xs transition-colors ${
            activeTab === 'terminal'
              ? 'border-primary text-foreground'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Terminal className="h-3.5 w-3.5" />
          Terminal
        </button>
        <button
          type="button"
          onClick={() => setActiveTab('settings')}
          className={`flex items-center gap-1.5 border-b-2 px-3 py-2 font-medium text-xs transition-colors ${
            activeTab === 'settings'
              ? 'border-primary text-foreground'
              : 'border-transparent text-muted-foreground hover:text-foreground'
          }`}
        >
          <Settings className="h-3.5 w-3.5" />
          Settings
        </button>
      </div>

      {activeTab === 'overview' && (
        <div className="space-y-6">
          <ServerConnectionCard server={server} onRefresh={() => refetch()} />

          <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
            <MetricChartCard
              title="CPU Usage"
              icon={<Cpu className="h-4 w-4" />}
              badge="Live"
              data={cpuData}
              isLoading={false}
              color="hsl(var(--chart-1))"
              formatY={(v) => `${Math.round(v)}%`}
              formatTooltip={(v) => `${v.toFixed(2)}%`}
            />
            <MetricChartCard
              title="Memory Usage"
              icon={<MemoryStick className="h-4 w-4" />}
              badge="Live"
              data={memData}
              isLoading={false}
              color="hsl(var(--chart-2))"
              formatY={(v) => `${Math.round(v)}%`}
              formatTooltip={(v) => `${v.toFixed(2)}%`}
            />
            <MetricChartCard
              title="Disk Usage"
              icon={<HardDrive className="h-4 w-4" />}
              badge="Live"
              data={diskData}
              isLoading={false}
              color="hsl(var(--chart-3))"
              formatY={(v) => `${Math.round(v)}%`}
              formatTooltip={(v) => `${v.toFixed(2)}%`}
            />
          </div>
        </div>
      )}

      {activeTab === 'terminal' && (
        <Card className="flex h-96 flex-col overflow-hidden rounded-xl border border-border bg-black font-mono text-green-400 text-xs">
          <div className="flex items-center justify-between border-border/40 border-b bg-zinc-900/80 px-4 py-2 text-zinc-400">
            <span>
              ssh {server.sshUser || 'root'}@{server.sshHost || server.ipAddress}
            </span>
            <span className="text-[11px] text-zinc-500">Interactive Shell Session</span>
          </div>
          <div className="flex-1 space-y-2 overflow-y-auto p-4 text-zinc-300">
            <p className="text-zinc-500">
              Connected to {server.name} ({server.ipAddress})
            </p>
            <p className="text-emerald-400">
              codedock-agent: worker runtime active, health probe ok
            </p>
            <div className="flex items-center gap-2 pt-2">
              <span className="text-emerald-500">
                {server.sshUser || 'root'}@{server.name}:~$
              </span>
              <span className="animate-pulse">_</span>
            </div>
          </div>
        </Card>
      )}

      {activeTab === 'settings' && <ServerSettingsTab server={server} />}
    </div>
  );
}
