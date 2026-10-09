import { Badge } from '#/components/ui/badge';
import { cn } from '#/lib/utils';
import type { AttentionIssue } from './interfaces';
import { kindMeta, SEVERITY_TONE, toDisplaySeverity } from './issue-meta';
import { IssueRow } from './issue-row';
import { KindIcon } from './kind-icon';

export function IssueGroup({
  kind,
  issues,
  standAlone,
  busyIds,
  onAcknowledge,
  onResolve,
  onAct,
}: {
  kind: string;
  issues: AttentionIssue[];
  standAlone: boolean;
  busyIds: ReadonlySet<string>;
  onAcknowledge: (issue: AttentionIssue) => void;
  onResolve: (issue: AttentionIssue) => void;
  onAct: (issue: AttentionIssue) => void;
}) {
  if (issues.length === 0) return null;
  const first = issues[0] as AttentionIssue;
  const worst = toDisplaySeverity(first.severity);
  const muted = worst === 'advisory' && !standAlone;
  const tone = SEVERITY_TONE[worst];
  const meta = kindMeta(kind);

  return (
    <section className="overflow-hidden rounded-2xl border border-border/60 bg-card">
      <header className="flex items-center gap-3 border-border/60 border-b px-5 py-4">
        <span
          className={cn(
            'flex size-9 shrink-0 items-center justify-center rounded-xl',
            muted ? 'bg-muted text-muted-foreground' : cn(tone.soft, tone.text)
          )}
        >
          <KindIcon kind={kind} className="size-[18px]" />
        </span>
        <div className="min-w-0 flex-1">
          <h2 className="font-semibold text-[15px]">{meta.title}</h2>
          <p className="truncate text-muted-foreground text-xs">{meta.subtitle}</p>
        </div>
        <Badge variant="secondary" className="shrink-0 tabular-nums">
          {issues.length}
        </Badge>
      </header>
      <ul className="divide-y divide-border/50">
        {issues.map((issue) => (
          <IssueRow
            key={issue.id}
            issue={issue}
            busy={busyIds.has(issue.id)}
            onAcknowledge={onAcknowledge}
            onResolve={onResolve}
            onAct={onAct}
          />
        ))}
      </ul>
    </section>
  );
}
