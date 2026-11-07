import { Search } from 'lucide-react';
import { useMemo, useRef, useState } from 'react';
import { Input } from '@codedock/ui/input';
import { cn } from '#/lib/utils';
import { EmptyIssues } from './empty-issues';
import type { AttentionIssue, MonitoringTab, SeverityFilter } from '#/interfaces/issues';
import { IssueList } from './issue-list';
import { SEVERITY_FILTERS, SEVERITY_LABELS, toDisplaySeverity } from './issue-meta';
import { IssueSummary } from './issue-summary';

export interface IssuesFeedViewProps {
  tabIssues: AttentionIssue[];
  tab: MonitoringTab;
  busyIds: Set<string>;
  onAcknowledge: (id: string) => void;
  onResolve: (id: string) => void;
  onAct: (issue: AttentionIssue) => void;
}

export function IssuesFeedView({
  tabIssues,
  tab,
  busyIds,
  onAcknowledge,
  onResolve,
  onAct,
}: IssuesFeedViewProps) {
  const [severity, setSeverity] = useState<SeverityFilter>('all');
  const [query, setQuery] = useState('');
  const [queryDraft, setQueryDraft] = useState('');
  const debounce = useRef<ReturnType<typeof setTimeout> | null>(null);

  const q = query.trim().toLowerCase();
  const filtered = useMemo(
    () =>
      tabIssues.filter((issue) => {
        if (severity !== 'all' && toDisplaySeverity(issue.severity) !== severity) return false;
        if (!q) return true;
        return `${issue.title} ${issue.detail} ${issue.subject} ${issue.kind}`
          .toLowerCase()
          .includes(q);
      }),
    [tabIssues, severity, q]
  );

  const facetCount = (filter: SeverityFilter): number =>
    filter === 'all'
      ? tabIssues.length
      : tabIssues.filter((issue) => toDisplaySeverity(issue.severity) === filter).length;

  const onQueryChange = (value: string) => {
    setQueryDraft(value);
    if (debounce.current) clearTimeout(debounce.current);
    debounce.current = setTimeout(() => setQuery(value), 300);
  };

  const clearFilters = () => {
    if (debounce.current) clearTimeout(debounce.current);
    setQueryDraft('');
    setQuery('');
    setSeverity('all');
  };

  if (tabIssues.length === 0) {
    return <EmptyIssues filtered={false} resolved={tab === 'resolved'} />;
  }

  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_340px]">
      <div className="min-w-0">
        <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-center">
          <div className="relative w-full sm:min-w-55 sm:flex-1">
            <Search className="pointer-events-none absolute inset-s-3.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="text"
              placeholder="Search issues"
              value={queryDraft}
              onChange={(event) => onQueryChange(event.target.value)}
              className="h-10 ps-10"
            />
          </div>
          <div className="inline-flex max-w-full shrink-0 flex-wrap items-center gap-1 rounded-xl bg-muted/40 p-1">
            {SEVERITY_FILTERS.map((filter) => {
              const count = facetCount(filter);
              return (
                <button
                  key={filter}
                  type="button"
                  onClick={() => setSeverity(filter)}
                  className={cn(
                    'inline-flex h-8 items-center gap-1.5 rounded-lg px-3.5 font-medium text-[12px] transition-colors',
                    severity === filter
                      ? 'border border-border/60 bg-card text-foreground'
                      : 'text-muted-foreground hover:bg-background/60 hover:text-foreground'
                  )}
                >
                  {filter === 'all' ? 'All' : SEVERITY_LABELS[filter]}
                  {count > 0 && (
                    <span className="text-muted-foreground/60 tabular-nums">{count}</span>
                  )}
                </button>
              );
            })}
          </div>
        </div>

        {filtered.length === 0 ? (
          <EmptyIssues filtered resolved={tab === 'resolved'} />
        ) : (
          <IssueList
            issues={filtered}
            busyIds={busyIds}
            onAcknowledge={(issue) => onAcknowledge(issue.id)}
            onResolve={(issue) => onResolve(issue.id)}
            onAct={onAct}
          />
        )}

        {filtered.length !== tabIssues.length && (
          <div className="mt-4 flex items-center gap-3">
            <p className="text-[12px] text-muted-foreground/70 tabular-nums">
              {filtered.length} / {tabIssues.length}
            </p>
            <button
              type="button"
              onClick={clearFilters}
              className="font-medium text-[12px] text-muted-foreground underline-offset-4 transition-colors hover:text-foreground hover:underline"
            >
              Clear filters
            </button>
          </div>
        )}
      </div>

      <div className="space-y-4 lg:sticky lg:top-6 lg:self-start">
        <IssueSummary issues={tabIssues} tab={tab} />
      </div>
    </div>
  );
}
