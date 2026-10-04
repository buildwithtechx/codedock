import { createFileRoute, Link } from '@tanstack/react-router';
import {
  Activity,
  AlertTriangle,
  ArrowRight,
  Box,
  CheckCircle2,
  Cpu,
  HardDrive,
  Loader2,
  MemoryStick,
  Server,
} from 'lucide-react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '#/components/ui/card';
import { Progress } from '#/components/ui/progress';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { useGetSystemStats } from '#/features/settings/hooks';
import { useListCanvasSummaries } from '#/hooks/use-canvas';

export const Route = createFileRoute('/_dashboard/monitoring')({
  component: MonitoringPage,
});

export function MonitoringPage() {
  const {
    data: statsRes,
    isLoading: isLoadingStats,
    isError: isStatsError,
    refetch: refetchStats,
  } = useGetSystemStats();
  const { data: canvasRes, isLoading: isLoadingProjects } = useListCanvasSummaries();

  const stats = statsRes?.data;
  const projects = canvasRes?.data || [];

  const totalServices = projects.reduce((acc, p) => acc + p.totalServices, 0);

  const cpuPercent = stats?.cpu?.percent ?? 0;
  const memPercent = stats?.memory?.percent ?? 0;
  const diskPercent = stats?.disk?.percent ?? 0;

  const issues = [];
  if (diskPercent > 85) {
    issues.push({
      level: 'warning',
      message: `Disk utilization is critically high (${diskPercent}% used). Consider pruning unused Docker layers.`,
    });
  }
  if (cpuPercent > 90) {
    issues.push({
      level: 'warning',
      message: `CPU utilization is elevated (${cpuPercent}%).`,
    });
  }
  if (memPercent > 90) {
    issues.push({
      level: 'warning',
      message: `System memory is running low (${memPercent}% used).`,
    });
  }

  const formatUptime = (seconds?: number) => {
    if (!seconds) return 'N/A';
    const days = Math.floor(seconds / 86400);
    const hours = Math.floor((seconds % 86400) / 3600);
    const mins = Math.floor((seconds % 3600) / 60);
    if (days > 0) return `${days}d ${hours}h ${mins}m`;
    if (hours > 0) return `${hours}h ${mins}m`;
    return `${mins}m`;
  };

  return (
    <div className="space-y-6">
      <PageHeader
        title="Monitoring"
        description="Global container health, system resource utilization, and runtime capacity."
      />

      {isStatsError ? (
        <QueryErrorState
          title="System metrics are unavailable"
          description="Could not connect to the host runtime to retrieve telemetry."
          onRetry={() => void refetchStats()}
        />
      ) : isLoadingStats ? (
        <div className="flex min-h-[18rem] items-center justify-center">
          <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
      ) : (
        <>
          {issues.length > 0 && (
            <div className="space-y-2">
              {issues.map((issue, idx) => (
                <div
                  key={idx}
                  className="flex items-center gap-3 rounded-xl border border-amber-500/20 bg-amber-500/10 p-4 text-amber-500 text-sm"
                >
                  <AlertTriangle className="h-4 w-4 shrink-0" />
                  <span className="flex-1">{issue.message}</span>
                </div>
              ))}
            </div>
          )}

          <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <Card className="border-border/80 bg-card shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                  CPU Utilization
                </CardTitle>
                <Cpu className="h-4 w-4 text-primary" />
              </CardHeader>
              <CardContent>
                <div className="font-bold text-2xl">{cpuPercent}%</div>
                <p className="mt-1 text-muted-foreground text-xs">{stats?.cpu?.cores || 1} Cores</p>
                <Progress value={cpuPercent} className="mt-3 h-1.5" />
              </CardContent>
            </Card>

            <Card className="border-border/80 bg-card shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                  Memory Usage
                </CardTitle>
                <MemoryStick className="h-4 w-4 text-primary" />
              </CardHeader>
              <CardContent>
                <div className="font-bold text-2xl">{memPercent}%</div>
                <p className="mt-1 text-muted-foreground text-xs">
                  {stats?.memory ? `${stats.memory.usedMB} / ${stats.memory.totalMB} MB` : 'N/A'}
                </p>
                <Progress value={memPercent} className="mt-3 h-1.5" />
              </CardContent>
            </Card>

            <Card className="border-border/80 bg-card shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                  Host Storage
                </CardTitle>
                <HardDrive className="h-4 w-4 text-primary" />
              </CardHeader>
              <CardContent>
                <div className="font-bold text-2xl">{diskPercent}%</div>
                <p className="mt-1 text-muted-foreground text-xs">
                  {stats?.disk ? `${stats.disk.usedGB} / ${stats.disk.totalGB} GB` : 'N/A'}
                </p>
                <Progress value={diskPercent} className="mt-3 h-1.5" />
              </CardContent>
            </Card>

            <Card className="border-border/80 bg-card shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="font-medium text-muted-foreground text-xs uppercase tracking-wider">
                  Host Health
                </CardTitle>
                <Server className="h-4 w-4 text-emerald-500" />
              </CardHeader>
              <CardContent>
                <div className="flex items-center gap-1.5 font-bold text-emerald-500 text-2xl">
                  <CheckCircle2 className="h-5 w-5" />
                  <span>Healthy</span>
                </div>
                <p className="mt-1 text-muted-foreground text-xs">
                  Uptime: {formatUptime(stats?.uptimeSeconds)}
                </p>
                <p className="mt-1 text-muted-foreground text-xs">
                  Load Avg: {stats?.loadAvg?.map((l) => l.toFixed(2)).join(', ') || 'N/A'}
                </p>
              </CardContent>
            </Card>
          </div>

          <div className="grid gap-6 lg:grid-cols-2">
            <Card className="border-border/80 bg-card shadow-sm">
              <CardHeader className="border-border/70 border-b pb-4">
                <div className="flex items-center gap-2">
                  <Activity className="h-4 w-4 text-primary" />
                  <CardTitle>Docker Daemon Telemetry</CardTitle>
                </div>
              </CardHeader>
              <CardContent className="pt-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="rounded-xl border border-border/70 p-3.5">
                    <p className="text-muted-foreground text-xs">Containers</p>
                    <p className="mt-1 font-semibold text-lg">
                      {stats?.docker?.containers?.active || '0'} active
                    </p>
                    <p className="text-muted-foreground text-xs">
                      {stats?.docker?.containers?.totalCount || '0'} total
                    </p>
                  </div>
                  <div className="rounded-xl border border-border/70 p-3.5">
                    <p className="text-muted-foreground text-xs">Images</p>
                    <p className="mt-1 font-semibold text-lg">
                      {stats?.docker?.images?.totalCount || '0'}
                    </p>
                    <p className="text-muted-foreground text-xs">
                      {stats?.docker?.images?.size || '0 B'}
                    </p>
                  </div>
                  <div className="rounded-xl border border-border/70 p-3.5">
                    <p className="text-muted-foreground text-xs">Volumes</p>
                    <p className="mt-1 font-semibold text-lg">
                      {stats?.docker?.volumes?.totalCount || '0'}
                    </p>
                    <p className="text-muted-foreground text-xs">
                      {stats?.docker?.volumes?.size || '0 B'}
                    </p>
                  </div>
                  <div className="rounded-xl border border-border/70 p-3.5">
                    <p className="text-muted-foreground text-xs">Reclaimable Cache</p>
                    <p className="mt-1 font-semibold text-lg">
                      {stats?.docker?.reclaimableGB ? `${stats.docker.reclaimableGB} GB` : '0 GB'}
                    </p>
                    <p className="text-muted-foreground text-xs">Build cache & unattached</p>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card className="border-border/80 bg-card shadow-sm">
              <CardHeader className="flex flex-row items-center justify-between border-border/70 border-b pb-4">
                <div className="flex items-center gap-2">
                  <Box className="h-4 w-4 text-primary" />
                  <CardTitle>Workload Deployments ({totalServices} services)</CardTitle>
                </div>
                <Button variant="ghost" size="sm" asChild>
                  <Link to="/projects">View all</Link>
                </Button>
              </CardHeader>
              <CardContent className="pt-4">
                {isLoadingProjects ? (
                  <div className="flex justify-center p-8">
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  </div>
                ) : projects.length === 0 ? (
                  <div className="p-8 text-center text-muted-foreground text-sm">
                    No active projects deployed.
                  </div>
                ) : (
                  <div className="divide-y divide-border/60">
                    {projects.map((project) => (
                      <div
                        key={project.id}
                        className="flex items-center justify-between py-2.5 first:pt-0 last:pb-0"
                      >
                        <div className="min-w-0 flex-1">
                          <Link
                            to="/projects/$projectId"
                            params={{ projectId: project.id }}
                            className="font-medium text-sm hover:underline"
                          >
                            {project.name}
                          </Link>
                          <p className="text-muted-foreground text-xs">
                            {project.totalServices} service{project.totalServices === 1 ? '' : 's'}
                          </p>
                        </div>
                        <Link
                          to="/projects/$projectId"
                          params={{ projectId: project.id }}
                          className="flex h-7 w-7 items-center justify-center rounded-lg text-muted-foreground hover:bg-muted hover:text-foreground"
                        >
                          <ArrowRight className="h-3.5 w-3.5" />
                        </Link>
                      </div>
                    ))}
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </>
      )}
    </div>
  );
}
