import { useMemo } from 'react';
import type { AttentionIssue } from './interfaces';
import { IssueGroup } from './issue-group';
import { KIND_ORDER, SEVERITY_RANK, toDisplaySeverity } from './issue-meta';

export function IssueList({
  issues,
  busyIds,
  onAcknowledge,
  onResolve,
  onAct,
}: {
  issues: AttentionIssue[];
  busyIds: ReadonlySet<string>;
  onAcknowledge: (issue: AttentionIssue) => void;
  onResolve: (issue: AttentionIssue) => void;
  onAct: (issue: AttentionIssue) => void;
}) {
  const grouped = useMemo(() => {
    const buckets = new Map<string, AttentionIssue[]>();
    for (const issue of issues) {
      const list = buckets.get(issue.kind);
      if (list) list.push(issue);
      else buckets.set(issue.kind, [issue]);
    }
    const rank = (kind: string) => {
      const index = KIND_ORDER.indexOf(kind);
      return index === -1 ? KIND_ORDER.length : index;
    };
    const groups = [...buckets.entries()]
      .map(([kind, rows]) => ({
        kind,
        rows: [...rows].sort((a, b) => {
          const bySeverity =
            SEVERITY_RANK[toDisplaySeverity(a.severity)] -
            SEVERITY_RANK[toDisplaySeverity(b.severity)];
          if (bySeverity !== 0) return bySeverity;
          return b.lastSeen.localeCompare(a.lastSeen);
        }),
      }))
      .sort((a, b) => rank(a.kind) - rank(b.kind));
    const needsAttention = (rows: AttentionIssue[]) =>
      rows.some((issue) => toDisplaySeverity(issue.severity) !== 'advisory');
    return groups.sort((a, b) => Number(needsAttention(b.rows)) - Number(needsAttention(a.rows)));
  }, [issues]);

  const advisoriesStandAlone = useMemo(
    () => !issues.some((issue) => toDisplaySeverity(issue.severity) !== 'advisory'),
    [issues]
  );

  return (
    <div className="space-y-4">
      {grouped.map(({ kind, rows }) => (
        <IssueGroup
          key={kind}
          kind={kind}
          issues={rows}
          standAlone={advisoriesStandAlone}
          busyIds={busyIds}
          onAcknowledge={onAcknowledge}
          onResolve={onResolve}
          onAct={onAct}
        />
      ))}
    </div>
  );
}
