import { Activity, Boxes, CheckCircle2, FolderKanban } from 'lucide-react';
import type { CanvasSummary } from '#/features/projects';

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
  const health = totalServices === 0 ? 100 : Math.round((onlineServices / totalServices) * 100);

  const items = [
    {
      label: 'Projects',
      value: projects.length,
      icon: FolderKanban,
      tone: 'text-primary bg-primary/12',
    },
    { label: 'Services', value: totalServices, icon: Boxes, tone: 'text-sky-500 bg-sky-500/12' },
    {
      label: 'Online now',
      value: onlineServices,
      icon: Activity,
      tone: 'text-emerald-500 bg-emerald-500/12',
    },
  ];

  return (
    <aside className="rounded-2xl border border-border/80 bg-card p-5 shadow-sm">
      <div className="flex items-center gap-2">
        <Activity className="h-4 w-4 text-muted-foreground" />
        <h2 className="font-semibold text-sm">Runtime pulse</h2>
      </div>
      <div className="mt-5 space-y-3.5">
        {items.map((item) => (
          <div key={item.label} className="flex items-center justify-between">
            <div className="flex items-center gap-2.5">
              <div className={`flex h-8 w-8 items-center justify-center rounded-lg ${item.tone}`}>
                <item.icon className="h-4 w-4" />
              </div>
              <span className="text-muted-foreground text-sm">{item.label}</span>
            </div>
            <span className="font-semibold text-base">
              {isLoading || isUnavailable ? '–' : item.value}
            </span>
          </div>
        ))}
      </div>
      <div className="mt-5 border-border/70 border-t pt-4">
        <div className="flex items-center justify-between text-sm">
          <span className="flex items-center gap-2 text-muted-foreground">
            <CheckCircle2 className="h-4 w-4 text-emerald-500" />
            Service health
          </span>
          <span className="font-semibold text-emerald-500">
            {isLoading || isUnavailable ? '–' : `${health}%`}
          </span>
        </div>
      </div>
    </aside>
  );
}
