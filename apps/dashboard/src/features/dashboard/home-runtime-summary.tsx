import { Activity, AlertTriangle, Boxes, CheckCircle2, FolderKanban } from 'lucide-react';
import type { CanvasSummary } from '#/interfaces/projects';

export function HomeRuntimeSummary({
  projects,
  isLoading,
  isUnavailable = false,
}: {
  projects: CanvasSummary[];
  isLoading: boolean;
  isUnavailable?: boolean;
}) {
  const totalServices = projects.reduce((total, project) => total + project.totalServices, 0);
  const onlineServices = projects.reduce((total, project) => total + project.onlineServices, 0);
  const health = totalServices === 0 ? null : Math.round((onlineServices / totalServices) * 100);
  const isHealthy = totalServices > 0 && onlineServices === totalServices;
  const isDegraded = totalServices > 0 && onlineServices < totalServices;

  const items = [
    {
      label: 'Projects',
      value: projects.length,
      icon: FolderKanban,
      tone: 'text-primary bg-primary/10',
    },
    {
      label: 'Total Services',
      value: totalServices,
      icon: Boxes,
      tone: 'text-sky-500 bg-sky-500/10',
    },
    {
      label: 'Online Now',
      value: onlineServices,
      icon: Activity,
      tone: 'text-emerald-500 bg-emerald-500/10',
    },
  ];

  return (
    <aside className="rounded-2xl border border-border/50 bg-card p-5 shadow-xs">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Activity className="size-4 text-primary" />
          <h2 className="font-semibold text-foreground text-sm">Runtime Pulse</h2>
        </div>
        {!isLoading && !isUnavailable && (
          <div className="flex items-center gap-1.5">
            {isHealthy && (
              <>
                <span className="size-2 animate-pulse rounded-full bg-emerald-500" />
                <span className="font-medium text-emerald-500 text-xs">Healthy</span>
              </>
            )}
            {isDegraded && (
              <>
                <span className="size-2 animate-pulse rounded-full bg-amber-500" />
                <span className="font-medium text-amber-500 text-xs">Degraded</span>
              </>
            )}
            {totalServices === 0 && (
              <span className="font-medium text-muted-foreground text-xs">Standby</span>
            )}
          </div>
        )}
      </div>

      <div className="mt-4 space-y-2.5">
        {items.map((item) => (
          <div
            key={item.label}
            className="flex items-center justify-between rounded-xl border border-border/30 bg-muted/20 px-3 py-2.5"
          >
            <div className="flex items-center gap-2.5">
              <div className={`flex size-7 items-center justify-center rounded-lg ${item.tone}`}>
                <item.icon className="size-3.5" />
              </div>
              <span className="font-medium text-muted-foreground text-xs">{item.label}</span>
            </div>
            <span className="font-bold text-foreground text-sm">
              {isLoading || isUnavailable ? '–' : item.value}
            </span>
          </div>
        ))}
      </div>

      <div className="mt-4 border-border/40 border-t pt-3">
        <div className="flex items-center justify-between text-xs">
          <span className="flex items-center gap-1.5 text-muted-foreground">
            {isDegraded ? (
              <AlertTriangle className="size-3.5 text-amber-500" />
            ) : (
              <CheckCircle2 className="size-3.5 text-emerald-500" />
            )}
            System Health
          </span>
          <span
            className={
              isDegraded ? 'font-semibold text-amber-500' : 'font-semibold text-emerald-500'
            }
          >
            {isLoading || isUnavailable
              ? '–'
              : health === null
                ? 'No services'
                : `${health}%`}
          </span>
        </div>
      </div>
    </aside>
  );
}
