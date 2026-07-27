import { createFileRoute, Link } from '@tanstack/react-router';
import { Cpu, HardDrive, Loader2, MemoryStick, ServerIcon } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { env } from '#/env';
import { MetricChartCard } from '#/features/services/metric-chart-card';
import { useListServers } from '#/hooks/use-servers';
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

function ServerDetailsPage() {
  const { serverId } = Route.useParams();
  const { data: servers, isLoading } = useListServers();
  const server = servers?.find((s) => s.id === serverId);

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
          // Keep last 30 data points (assuming 5s interval -> 2.5 minutes)
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

  const statusColor = server.status === 'online' ? 'bg-green-500' : 'bg-yellow-500';

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-12 w-12 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
            <ServerIcon className="h-6 w-6" />
          </div>
          <div>
            <h1 className="font-bold text-xl">{server.name}</h1>
            <p className="text-muted-foreground text-sm">{server.ipAddress}</p>
          </div>
        </div>
        <div className="flex items-center gap-2 rounded-full border bg-background px-3 py-1 shadow-sm">
          <div
            className={`h-2.5 w-2.5 rounded-full ${statusColor} ${server.status === 'online' ? 'animate-pulse' : ''} shadow-sm`}
          />
          <span className="font-semibold text-xs uppercase tracking-wider">{server.status}</span>
        </div>
      </div>

      <div className="flex items-center gap-2">
        <span className={`h-2 w-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
        <span className="text-muted-foreground text-xs">
          {isConnected ? 'Real-time metrics connected' : 'Disconnected from metrics stream'}
        </span>
      </div>

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
  );
}
