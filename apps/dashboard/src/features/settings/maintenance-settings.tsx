import { AlertTriangle, HardDrive, RefreshCw, Trash2, Wrench } from 'lucide-react';
import { useState } from 'react';
import { toast } from 'sonner';
import { Badge } from '#/components/ui/badge';
import { Button } from '#/components/ui/button';
import { useCleanupSystem, useGetSystemStats, useRestartSystem } from '#/features/settings';
import { DockerStorageCard } from './docker-storage-card';
import { MaintenanceDialogs } from './maintenance-dialogs';
import { SettingsSection } from './settings-section';

const ProgressBar = ({ value, colorClass }: { value: number; colorClass: string }) => (
  <div className="flex h-1.5 w-full overflow-hidden rounded-full bg-background">
    <div className={`h-full ${colorClass}`} style={{ width: `${value}%` }} />
  </div>
);

export const MaintenancePage = () => {
  const { data: statsData, refetch } = useGetSystemStats();
  const { mutateAsync: cleanup, isPending: cleaning } = useCleanupSystem();
  const { mutateAsync: restart, isPending: restarting } = useRestartSystem();
  const [confirmCleanup, setConfirmCleanup] = useState(false);
  const [confirmRestart, setConfirmRestart] = useState(false);

  const stats = statsData?.data;

  const handleCleanup = async () => {
    try {
      await cleanup();
      toast.success('Docker cleanup completed');
    } catch {
      toast.error('Cleanup failed');
    } finally {
      setConfirmCleanup(false);
    }
  };

  const handleRestart = async () => {
    try {
      await restart();
      toast.success('Restart initiated');
    } catch {
      toast.error('Restart failed');
    } finally {
      setConfirmRestart(false);
    }
  };

  const usedPercent = stats?.disk.percent ? Number(stats.disk.percent.toFixed(1)) : 0;
  const freeGb = stats?.disk.freeGb ? stats.disk.freeGb.toFixed(1) : '0';
  const usedGb = stats?.disk.usedGb ? stats.disk.usedGb.toFixed(1) : '0';
  const totalGb = stats?.disk.totalGb ? stats.disk.totalGb.toFixed(1) : '0';

  const dockerPercent =
    stats?.disk?.totalGb && stats?.docker?.reclaimableGb
      ? Number(((stats.docker.reclaimableGb / stats.disk.totalGb) * 100).toFixed(1))
      : 0;
  const reclaimableGb = stats?.docker?.reclaimableGb ? stats.docker.reclaimableGb.toFixed(2) : '0';

  const buildCacheReclaimableStr = stats?.docker?.buildCache?.reclaimable || '0';
  const buildCacheGb = parseFloat(buildCacheReclaimableStr)
    ? parseFloat(buildCacheReclaimableStr).toFixed(2)
    : '0';

  return (
    <div className="space-y-6">
      <SettingsSection
        icon={<Wrench className="size-4 text-primary" />}
        title="System Maintenance"
        description="Monitor root disk capacity, Docker cache, and run cleanup routines."
        action={
          <div className="flex shrink-0 items-center gap-2.5">
            {Number(reclaimableGb) > 3 ? (
              <Badge
                variant="outline"
                className="border-destructive/50 bg-destructive/10 px-2.5 py-0.5 font-semibold text-[10px] text-destructive uppercase tracking-wider"
              >
                Attention
              </Badge>
            ) : (
              <Badge
                variant="outline"
                className="border-primary/50 bg-primary/10 px-2.5 py-0.5 font-semibold text-[10px] text-primary uppercase tracking-wider"
              >
                Healthy
              </Badge>
            )}
            <Button variant="outline" size="sm" onClick={() => refetch()} className="gap-1.5">
              <RefreshCw className="h-3.5 w-3.5" />
              Refresh
            </Button>
          </div>
        }
      >
        <div className="space-y-6">
          {stats?.docker?.reclaimableGb && stats.docker.reclaimableGb > 3 ? (
            <div className="flex w-full items-center gap-3 rounded-lg border border-destructive/30 bg-destructive/10 p-4 font-medium text-destructive text-sm">
              <AlertTriangle className="h-4 w-4" /> Docker has more than 3 GB reclaimable.
            </div>
          ) : null}

          <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
            <div className="flex flex-col justify-between space-y-6 rounded-2xl border border-border/80 bg-card p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  ROOT DISK FREE
                </p>
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  {usedPercent}% USED
                </p>
              </div>
              <h2 className="font-bold text-3xl">{freeGb} GB</h2>
              <div className="space-y-2">
                <ProgressBar value={usedPercent} colorClass="bg-primary" />
                <p className="text-muted-foreground text-xs">
                  {usedGb} GB used of {totalGb} GB on /.
                </p>
              </div>
            </div>

            <div className="flex flex-col justify-between space-y-6 rounded-2xl border border-border/80 bg-card p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <p className="w-32 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  DOCKER CLEANUP CANDIDATES
                </p>
                <p className="w-20 text-right font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  {dockerPercent}% OF DISK
                </p>
              </div>
              <h2 className="font-bold text-3xl">{reclaimableGb} GB</h2>
              <div className="space-y-2">
                <ProgressBar value={dockerPercent} colorClass="bg-yellow-500" />
                <p className="text-muted-foreground text-xs">
                  {buildCacheGb} GB is build cache. Safe cleanup can usually clear that.
                </p>
              </div>
            </div>

            <div className="flex flex-col justify-between space-y-6 rounded-2xl border border-border/80 bg-card p-6 shadow-sm">
              <div className="flex items-center justify-between">
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  BUILD ARTIFACTS
                </p>
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  0%
                </p>
              </div>
              <h2 className="font-bold text-3xl">0 GB</h2>
              <div className="space-y-2">
                <ProgressBar value={0} colorClass="bg-primary" />
                <p className="text-muted-foreground text-xs">No build artifact directory yet.</p>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-6 md:grid-cols-3">
            <div className="flex h-32 flex-col rounded-2xl border border-border/80 bg-card p-6 shadow-sm">
              <p className="mb-4 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                DISK TREND
              </p>
              <div className="flex flex-1 items-end gap-1">
                {Array.from({ length: 15 }).map((_, i) => (
                  <div
                    key={i}
                    className="flex-1 rounded-sm bg-primary/40"
                    style={{ height: `${40 + i * 2}%` }}
                  />
                ))}
              </div>
            </div>
            <div className="flex h-32 flex-col rounded-2xl border border-border/80 bg-card p-6 shadow-sm">
              <p className="mb-4 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                DOCKER RECLAIMABLE TREND
              </p>
              <div className="flex flex-1 items-end gap-1">
                {Array.from({ length: 15 }).map((_, i) => (
                  <div
                    key={i}
                    className="flex-1 rounded-sm bg-yellow-500/40"
                    style={{ height: `${60 - i}%` }}
                  />
                ))}
              </div>
            </div>
            <div className="flex h-32 flex-col rounded-2xl border border-border/80 bg-card p-6 shadow-sm">
              <p className="mb-4 font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                BUILD ARTIFACT TREND
              </p>
              <div className="flex flex-1 items-end gap-1">
                {Array.from({ length: 15 }).map((_, i) => (
                  <div key={i} className="h-1 flex-1 rounded-sm bg-primary/10" />
                ))}
              </div>
            </div>
          </div>

          <DockerStorageCard stats={stats} />

          <div className="flex flex-col space-y-6 rounded-2xl border border-border/80 bg-card p-6 shadow-sm">
            <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg border border-primary/20 bg-primary/10 text-primary">
                  <HardDrive className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="font-bold text-xl">Cleanup Targets</h3>
                  <p className="mt-1 text-muted-foreground text-sm">
                    Reclaimable space across Docker builds, service logs, and packages.
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-2.5">
                <Button
                  variant="outline"
                  onClick={() => setConfirmRestart(true)}
                  disabled={restarting}
                  className="border-destructive/30 bg-destructive/5 font-semibold text-destructive text-xs uppercase tracking-wider hover:bg-destructive/10 hover:text-destructive"
                >
                  <Trash2 className="mr-2 h-4 w-4" /> RESTART DAEMON
                </Button>
                <Button
                  onClick={() => setConfirmCleanup(true)}
                  disabled={cleaning}
                  className="font-semibold text-xs uppercase tracking-wider"
                >
                  <RefreshCw className="mr-2 h-4 w-4" /> {cleaning ? 'CLEANING...' : 'SAFE CLEANUP'}
                </Button>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
              <div className="rounded-xl border border-border/70 bg-background/50 p-4">
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  TOP DOCKER CANDIDATE
                </p>
                <p className="mt-2 font-mono font-semibold text-foreground text-lg">
                  {stats?.docker?.buildCache?.reclaimable || '0 B'}
                </p>
                <p className="mt-1 text-muted-foreground text-xs">Build Cache</p>
              </div>
              <div className="rounded-xl border border-border/70 bg-background/50 p-4">
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  CODEDOCK DATA
                </p>
                <p className="mt-2 font-mono font-semibold text-foreground text-lg">1.81 MB</p>
                <p className="mt-1 text-muted-foreground text-xs">Runtime configs</p>
              </div>
              <div className="rounded-xl border border-border/70 bg-background/50 p-4">
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  BACKUPS
                </p>
                <p className="mt-2 font-mono font-semibold text-foreground text-lg">
                  {stats?.backups?.size || '0 B'}
                </p>
                <p className="mt-1 text-muted-foreground text-xs">Archived snapshots</p>
              </div>
              <div className="rounded-xl border border-border/70 bg-background/50 p-4">
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  PACKAGE CACHE
                </p>
                <p className="mt-2 font-mono font-semibold text-foreground text-lg">
                  {stats?.packageCache?.size || '0 B'}
                </p>
                <p className="mt-1 text-muted-foreground text-xs">Package downloads</p>
              </div>
              <div className="rounded-xl border border-border/70 bg-background/50 p-4">
                <p className="font-bold text-[10px] text-muted-foreground uppercase tracking-[0.15em]">
                  SYSTEM LOGS
                </p>
                <p className="mt-2 font-mono font-semibold text-foreground text-lg">
                  {stats?.systemLogs?.size || '0 B'}
                </p>
                <p className="mt-1 text-muted-foreground text-xs">Daemon journal logs</p>
              </div>
            </div>

            <div className="rounded-xl border border-border/60 bg-background/30 p-4 text-muted-foreground text-xs leading-relaxed">
              Safe cleanup purges dangling image layers, unused build stages, and temporary package
              caches without interrupting running containers or deleting production volumes.
            </div>
          </div>
        </div>
      </SettingsSection>

      <MaintenanceDialogs
        confirmCleanup={confirmCleanup}
        setConfirmCleanup={setConfirmCleanup}
        cleaning={cleaning}
        handleCleanup={handleCleanup}
        confirmRestart={confirmRestart}
        setConfirmRestart={setConfirmRestart}
        restarting={restarting}
        handleRestart={handleRestart}
      />
    </div>
  );
};
