import { History } from 'lucide-react';
import { useMemo } from 'react';
import type { AuditCategoryFacet } from './api';

export function AuditSummaryCard({
  total,
  facets,
}: {
  total: number;
  facets: AuditCategoryFacet[];
}) {
  const top = useMemo(() => [...facets].sort((a, b) => b.count - a.count).slice(0, 5), [facets]);

  return (
    <div className="rounded-2xl border border-border/50 bg-card p-5">
      <div className="flex items-center gap-3">
        <div className="flex size-9 items-center justify-center rounded-xl bg-muted">
          <History className="size-[18px] text-muted-foreground" />
        </div>
        <div>
          <h2 className="font-semibold text-[15px] text-foreground">About this log</h2>
          <p className="text-muted-foreground text-xs">Who did what, newest first</p>
        </div>
      </div>
      <div className="mt-4 space-y-2.5 border-border/40 border-t pt-4 text-sm">
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Events</span>
          <span className="font-medium tabular-nums">{total}</span>
        </div>
        {top.map((entry) => (
          <div key={entry.id} className="flex items-center justify-between">
            <span className="text-muted-foreground">{entry.label}</span>
            <span className="font-medium tabular-nums">{entry.count}</span>
          </div>
        ))}
      </div>
    </div>
  );
}
