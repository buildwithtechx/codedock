import { ScrollText } from 'lucide-react';

export type JobStatusFilter = 'all' | 'running' | 'failed' | 'scheduled' | 'disabled';

export function JobsOverview({
  counts,
  active,
  onSelect,
}: {
  counts: Record<JobStatusFilter, number>;
  active: JobStatusFilter;
  onSelect: (filter: JobStatusFilter) => void;
}) {
  const rows: Array<{ key: JobStatusFilter; label: string; dot?: string }> = [
    { key: 'all', label: 'All jobs' },
    { key: 'running', label: 'Running', dot: 'bg-warning' },
    { key: 'failed', label: 'Failed', dot: 'bg-destructive' },
    { key: 'scheduled', label: 'Scheduled', dot: 'bg-primary' },
    { key: 'disabled', label: 'Disabled', dot: 'bg-muted-foreground/40' },
  ];

  return (
    <div className="rounded-2xl border border-border/50 bg-card">
      <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
        <div className="flex size-9 items-center justify-center rounded-xl bg-muted">
          <ScrollText className="size-[18px] text-muted-foreground" />
        </div>
        <div>
          <h2 className="font-semibold text-[15px] text-foreground">Overview</h2>
          <p className="text-muted-foreground text-xs">Click a row to filter the list</p>
        </div>
      </div>
      <div className="p-2">
        {rows.map((row) => {
          const isActive = active === row.key;
          return (
            <button
              key={row.key}
              type="button"
              onClick={() => onSelect(isActive && row.key !== 'all' ? 'all' : row.key)}
              className={`flex w-full items-center justify-between gap-3 rounded-xl px-3 py-2.5 text-start transition-colors ${
                isActive ? 'bg-primary/10' : 'hover:bg-muted/40'
              }`}
            >
              <span className="flex min-w-0 items-center gap-2.5">
                <span className={`size-1.5 shrink-0 rounded-full ${row.dot ?? 'bg-transparent'}`} />
                <span
                  className={`truncate text-sm ${
                    isActive ? 'font-medium text-foreground' : 'text-muted-foreground'
                  }`}
                >
                  {row.label}
                </span>
              </span>
              <span
                className={`font-medium text-sm tabular-nums ${
                  isActive ? 'text-foreground' : 'text-muted-foreground/80'
                }`}
              >
                {counts[row.key]}
              </span>
            </button>
          );
        })}
      </div>
    </div>
  );
}
