import { Activity } from 'lucide-react';
import { cn } from '#/lib/utils';
import type { AttentionIssue, DisplaySeverity, MonitoringTab } from './interfaces';
import {
  KIND_ORDER,
  kindMeta,
  SEVERITY_LABELS,
  SEVERITY_TONE,
  toDisplaySeverity,
} from './issue-meta';
import { KindIcon } from './kind-icon';

const SEVERITY_ROWS: DisplaySeverity[] = ['outage', 'action_required', 'advisory'];

export function IssueSummary({ issues, tab }: { issues: AttentionIssue[]; tab: MonitoringTab }) {
  const total = issues.length;
  const sevCount: Record<DisplaySeverity, number> = {
    outage: 0,
    action_required: 0,
    advisory: 0,
  };
  const kindCount = new Map<string, number>();
  for (const issue of issues) {
    sevCount[toDisplaySeverity(issue.severity)] += 1;
    kindCount.set(issue.kind, (kindCount.get(issue.kind) ?? 0) + 1);
  }
  const rank = (kind: string) => {
    const index = KIND_ORDER.indexOf(kind);
    return index === -1 ? KIND_ORDER.length : index;
  };
  const areas = [...kindCount.entries()]
    .map(([kind, count]) => ({ kind, count }))
    .sort((a, b) => b.count - a.count || rank(a.kind) - rank(b.kind));

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        <div className="flex size-9 shrink-0 items-center justify-center rounded-xl bg-muted">
          <Activity className="size-[18px] text-muted-foreground" />
        </div>
        <div className="min-w-0">
          <h2 className="font-semibold text-[15px]">Summary</h2>
          <p className="text-muted-foreground text-xs">Severity mix and blast radius</p>
        </div>
      </div>

      <div className="space-y-5 p-5">
        <div>
          <div className="flex items-baseline gap-1.5">
            <span className="font-semibold text-2xl tabular-nums">{total}</span>
            <span className="text-muted-foreground text-sm">
              {tab === 'resolved' ? 'resolved' : 'open'}
            </span>
          </div>
          <div className="mt-2.5 flex h-1.5 w-full overflow-hidden rounded-full bg-muted">
            {SEVERITY_ROWS.map((severity) =>
              sevCount[severity] > 0 ? (
                <div
                  key={severity}
                  className={cn(
                    'h-full transition-[width] duration-500',
                    SEVERITY_TONE[severity].bar
                  )}
                  style={{ width: `${(sevCount[severity] / total) * 100}%` }}
                />
              ) : null
            )}
          </div>
        </div>

        <div>
          <p className="mb-1.5 font-medium text-[11px] text-muted-foreground/60 uppercase tracking-wider">
            By severity
          </p>
          <div className="space-y-0.5">
            {SEVERITY_ROWS.map((severity) => {
              const count = sevCount[severity];
              return (
                <div key={severity} className="flex items-center justify-between py-1.5">
                  <span className="inline-flex items-center gap-2.5 text-sm">
                    <span
                      className={cn(
                        'size-2 rounded-full',
                        count > 0 ? SEVERITY_TONE[severity].dot : 'bg-muted-foreground/25'
                      )}
                    />
                    <span
                      className={count > 0 ? 'text-muted-foreground' : 'text-muted-foreground/50'}
                    >
                      {SEVERITY_LABELS[severity]}
                    </span>
                  </span>
                  <span
                    className={cn(
                      'font-medium text-sm tabular-nums',
                      count > 0 ? 'text-foreground' : 'text-muted-foreground/40'
                    )}
                  >
                    {count}
                  </span>
                </div>
              );
            })}
          </div>
        </div>

        {areas.length > 0 && (
          <div className="border-border/50 border-t pt-4">
            <p className="mb-1.5 font-medium text-[11px] text-muted-foreground/60 uppercase tracking-wider">
              By area
            </p>
            <div className="space-y-0.5">
              {areas.map(({ kind, count }) => (
                <div key={kind} className="flex items-center justify-between py-1.5">
                  <span className="inline-flex items-center gap-2.5 text-muted-foreground text-sm">
                    <KindIcon kind={kind} className="size-4 text-muted-foreground/60" />
                    {kindMeta(kind).title}
                  </span>
                  <span className="font-medium text-foreground text-sm tabular-nums">{count}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
