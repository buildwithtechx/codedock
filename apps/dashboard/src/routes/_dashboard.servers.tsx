import { createFileRoute, Link } from '@tanstack/react-router';
import { Cpu, HardDrive, Loader2, MemoryStick, Plus, ServerIcon } from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '#/components/ui/card';
import { Progress } from '#/components/ui/progress';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { ServerEmptyState } from '#/features/servers/server-empty-state';
import { useListServers } from '#/hooks/use-servers';
import { parseServerMetrics, type Server } from '#/interfaces/server';

export const Route = createFileRoute('/_dashboard/servers')({
  component: ServersPage,
});

function ServersPage() {
  const { data: servers, isLoading, isError, refetch } = useListServers();

  return (
    <div className="space-y-6">
      <PageHeader
        title="Servers"
        description={
          isLoading
            ? 'Loading infrastructure...'
            : `${servers?.length ?? 0} connected server${servers?.length === 1 ? '' : 's'}`
        }
        action={
          <Link to="/servers/new">
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              New server
            </Button>
          </Link>
        }
      />

      {isLoading ? (
        <div className="flex justify-center p-12">
          <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
        </div>
      ) : isError ? (
        <QueryErrorState
          title="Servers are unavailable"
          description="Codedock could not load the current runtime fleet."
          onRetry={() => void refetch()}
        />
      ) : servers?.length === 0 ? (
        <ServerEmptyState />
      ) : (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {servers?.map((server) => (
            <ServerCard key={server.id} server={server} />
          ))}
        </div>
      )}
    </div>
  );
}

function formatBytes(bytes: number) {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return `${Number.parseFloat((bytes / k ** i).toFixed(1))} ${sizes[i]}`;
}

function ServerCard({ server }: { server: Server }) {
  const statusColor = server.status === 'online' ? 'bg-green-500' : 'bg-yellow-500';
  const m = parseServerMetrics(server.metrics);

  const memPercent =
    m && m.memory_limit_bytes > 0 ? (m.memory_usage_bytes / m.memory_limit_bytes) * 100 : 0;
  const diskPercent =
    m && m.disk_total_bytes > 0 ? (m.disk_usage_bytes / m.disk_total_bytes) * 100 : 0;
  const cpuPercent = m ? m.cpu_usage_percentage : 0;

  return (
    <Link to="/servers/$serverId" params={{ serverId: server.id }} className="block">
      <Card className="flex h-full flex-col overflow-hidden transition-all hover:shadow-md">
        <CardHeader className="flex flex-row items-center justify-between border-b bg-muted/20 px-6 py-4 pb-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-full bg-primary/10">
              <ServerIcon className="h-5 w-5 text-primary" />
            </div>
            <div>
              <CardTitle className="font-semibold text-base">{server.name}</CardTitle>
              <CardDescription className="text-xs">
                {server.isControlPlane ? 'Local deployment runtime' : server.ipAddress}
              </CardDescription>
            </div>
          </div>
          <div className="flex items-center gap-2 rounded-full border bg-background px-3 py-1 shadow-sm">
            <div className={`h-2.5 w-2.5 rounded-full ${statusColor} animate-pulse shadow-sm`} />
            <span className="font-semibold text-xs uppercase tracking-wider">{server.status}</span>
          </div>
        </CardHeader>

        <CardContent className="flex-1 space-y-5 p-6">
          <div className="grid grid-cols-3 gap-4">
            <div className="space-y-2">
              <div className="flex items-center justify-between text-xs">
                <span className="flex items-center gap-1.5 font-medium text-muted-foreground">
                  <Cpu className="h-3.5 w-3.5" /> CPU
                </span>
                <span className="font-semibold">{cpuPercent.toFixed(1)}%</span>
              </div>
              <Progress value={cpuPercent} className="h-1.5" />
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between text-xs">
                <span className="flex items-center gap-1.5 font-medium text-muted-foreground">
                  <MemoryStick className="h-3.5 w-3.5" /> RAM
                </span>
                <span className="font-semibold">{memPercent.toFixed(1)}%</span>
              </div>
              <Progress value={memPercent} className="h-1.5" />
              <div className="text-right text-[10px] text-muted-foreground">
                {m
                  ? `${formatBytes(m.memory_usage_bytes)} / ${formatBytes(m.memory_limit_bytes)}`
                  : 'N/A'}
              </div>
            </div>

            <div className="space-y-2">
              <div className="flex items-center justify-between text-xs">
                <span className="flex items-center gap-1.5 font-medium text-muted-foreground">
                  <HardDrive className="h-3.5 w-3.5" /> Disk
                </span>
                <span className="font-semibold">{diskPercent.toFixed(1)}%</span>
              </div>
              <Progress value={diskPercent} className="h-1.5" />
              <div className="text-right text-[10px] text-muted-foreground">
                {m
                  ? `${formatBytes(m.disk_usage_bytes)} / ${formatBytes(m.disk_total_bytes)}`
                  : 'N/A'}
              </div>
            </div>
          </div>

          <div className="pt-2">
            <CardDescription className="text-center text-[11px]">
              Last heartbeat:{' '}
              {server.lastSeenAt ? new Date(server.lastSeenAt).toLocaleString() : 'Never'}
            </CardDescription>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}
