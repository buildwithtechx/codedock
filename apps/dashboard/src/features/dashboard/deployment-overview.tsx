import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts';
import { Button } from '#/components/ui/button';

export interface DeploymentHistorySummary {
  total: number;
  success: number;
  failed: number;
  active: number;
  canceled: number;
  avgBuildDurationMs: number | null;
  daily: Array<{
    date: string;
    total: number;
    success: number;
    failed: number;
    active: number;
    canceled: number;
  }>;
  projects: Array<{
    projectId: string;
    name?: string;
    total: number;
    success: number;
    failed: number;
    active: number;
    canceled: number;
  }>;
}

interface DeploymentOverviewProps {
  summary: DeploymentHistorySummary | null;
  loading: boolean;
  days: number;
  onDaysChange: (days: number) => void;
  onFilterChange: (filter: string) => void;
  onProjectChange: (id: string) => void;
  onRefresh: () => void;
}

const SERIES = [
  { key: 'success', label: 'Success', color: '#22c55e' },
  { key: 'failed', label: 'Failed', color: '#ef4444' },
  { key: 'active', label: 'In Progress', color: '#f59e0b' },
  { key: 'canceled', label: 'Canceled', color: '#94a3b8' },
] as const;

export function DeploymentOverview({
  summary,
  loading,
  days,
  onDaysChange,
  onFilterChange,
  onProjectChange,
  onRefresh,
}: DeploymentOverviewProps) {
  const completed = summary ? summary.success + summary.failed : 0;
  const metrics = [
    {
      label: 'Deployment attempts',
      value: summary ? summary.total.toLocaleString() : '0',
      filter: 'all',
      detail: 'Across the selected window',
    },
    {
      label: 'Success rate',
      value: completed > 0 ? `${Math.round((summary!.success / completed) * 100)}%` : '—',
      filter: 'success',
      detail: 'Successful / (successful + failed); excludes cancellations',
    },
    {
      label: 'Failed',
      value: summary ? summary.failed.toLocaleString() : '0',
      filter: 'failed',
      detail: 'Includes partial and blocked releases',
    },
    {
      label: 'In Progress',
      value: summary ? summary.active.toLocaleString() : '0',
      filter: 'active',
      detail: 'Running or queued attempts in this window',
    },
  ] as const;

  return (
    <section aria-label="Deployment overview" className="space-y-4" aria-busy={loading}>
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h2 className="text-base font-medium text-foreground">Deployment overview</h2>
          <p className="mt-1 text-xs text-muted-foreground">
            Selected project and search · UTC calendar days · retained deployment history
          </p>
        </div>
        <div className="flex items-center gap-2">
          <div className="flex rounded-xl bg-muted/60 p-1" role="group" aria-label="Time window">
            {[7, 30, 90].map((value) => (
              <button
                key={value}
                type="button"
                aria-pressed={days === value}
                onClick={() => onDaysChange(value)}
                className={`rounded-lg px-3 py-1.5 text-xs font-medium transition-colors ${
                  days === value
                    ? 'bg-background text-foreground shadow-sm'
                    : 'text-muted-foreground hover:text-foreground'
                }`}
              >
                {value} days
              </button>
            ))}
          </div>
          <Button variant="secondary" size="sm" disabled={loading} onClick={onRefresh}>
            Refresh
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-3 xl:grid-cols-4">
        {metrics.map((metric) => (
          <button
            key={metric.filter}
            type="button"
            onClick={() => onFilterChange(metric.filter)}
            className="flex flex-col items-start rounded-2xl border border-border/50 bg-card p-5 text-start transition-colors hover:bg-muted/60 focus-visible:ring-2 focus-visible:ring-ring"
          >
            <span className="text-xs text-muted-foreground">{metric.label}</span>
            <span className="mt-2 block text-3xl font-medium tabular-nums text-foreground">
              {metric.value}
            </span>
            <span className="mt-2 block text-xs text-muted-foreground">{metric.detail}</span>
          </button>
        ))}
      </div>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_340px]">
        <div className="min-w-0 rounded-2xl border border-border/50 bg-card p-5">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <h3 className="text-sm font-medium text-foreground">Daily outcomes</h3>
            <span className="text-xs text-muted-foreground">
              Average successful build:{' '}
              {summary?.avgBuildDurationMs != null
                ? `${Math.round(summary.avgBuildDurationMs / 1000)}s`
                : '—'}
            </span>
          </div>

          <div className="mt-4 h-48" aria-label="Daily outcomes">
            {summary && summary.total > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={summary.daily}>
                  <CartesianGrid vertical={false} stroke="currentColor" strokeOpacity={0.1} />
                  <XAxis
                    dataKey="date"
                    tickLine={false}
                    axisLine={false}
                    minTickGap={28}
                    tick={{ fontSize: 11 }}
                    tickFormatter={(date) => {
                      const d = new Date(`${date}T00:00:00Z`);
                      return d.toLocaleDateString(undefined, {
                        month: 'short',
                        day: 'numeric',
                        timeZone: 'UTC',
                      });
                    }}
                  />
                  <YAxis
                    allowDecimals={false}
                    tickLine={false}
                    axisLine={false}
                    width={30}
                    tick={{ fontSize: 11 }}
                  />
                  <Tooltip
                    contentStyle={{
                      backgroundColor: 'hsl(var(--card))',
                      borderColor: 'hsl(var(--border))',
                      borderRadius: '0.75rem',
                      fontSize: '0.75rem',
                    }}
                  />
                  {SERIES.map((item) => (
                    <Bar
                      key={item.key}
                      dataKey={item.key}
                      name={item.label}
                      stackId="outcomes"
                      fill={item.color}
                      maxBarSize={24}
                    />
                  ))}
                </BarChart>
              </ResponsiveContainer>
            ) : (
              <div className="flex h-full items-center justify-center text-sm text-muted-foreground">
                No deployments in this window.
              </div>
            )}
          </div>

          <div className="mt-3 flex flex-wrap gap-x-4 gap-y-2 text-xs text-muted-foreground">
            {SERIES.map((item) => (
              <span key={item.key} className="flex items-center gap-1.5">
                <span className="h-2 w-2 rounded-full" style={{ background: item.color }} />
                {item.label}
              </span>
            ))}
          </div>
        </div>

        <div className="rounded-2xl border border-border/50 bg-card p-5">
          <h3 className="text-sm font-medium text-foreground">Project breakdown</h3>
          <p className="mt-1 text-xs text-muted-foreground">
            Top 6 by failures, then active attempts
          </p>
          <div className="mt-3 space-y-1">
            {summary && summary.projects.length > 0 ? (
              summary.projects.slice(0, 6).map((project) => (
                <button
                  key={project.projectId}
                  type="button"
                  onClick={() => onProjectChange(project.projectId)}
                  className="flex w-full items-center justify-between gap-3 rounded-xl px-2 py-2.5 text-start hover:bg-muted/60 focus-visible:ring-2 focus-visible:ring-ring"
                >
                  <span className="truncate text-sm">
                    {project.name || project.projectId}
                  </span>
                  <span
                    className={`shrink-0 text-xs tabular-nums ${
                      project.failed > 0 ? 'text-destructive font-medium' : 'text-muted-foreground'
                    }`}
                  >
                    {project.failed}/{project.total} failed
                  </span>
                </button>
              ))
            ) : (
              <p className="py-4 text-sm text-muted-foreground">No deployments in this window.</p>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
