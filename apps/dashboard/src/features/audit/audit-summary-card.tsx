import { History } from 'lucide-react';
import { useMemo } from 'react';
import type { AuditLog } from '#/interfaces/audit';
import { relativeTime } from './audit-taxonomy';

export function AuditSummaryCard({ logs }: { logs: AuditLog[] }) {
  const actors = useMemo(() => new Set(logs.map((log) => log.userId)).size, [logs]);
  const oldest = useMemo(() => {
    if (logs.length === 0) return null;
    return logs.reduce((a, b) => (+new Date(a.createdAt) < +new Date(b.createdAt) ? a : b));
  }, [logs]);

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
          <span className="font-medium tabular-nums">{logs.length}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Actors</span>
          <span className="font-medium tabular-nums">{actors}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-muted-foreground">Oldest</span>
          <span className="font-medium">{oldest ? relativeTime(oldest.createdAt) : '—'}</span>
        </div>
      </div>
    </div>
  );
}
