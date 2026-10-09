import { Link } from '@tanstack/react-router';
import { Building2, Loader2, RefreshCw, Search } from 'lucide-react';
import { useMemo, useRef, useState } from 'react';
import { PageHeader } from '#/components/layout/page-header';
import { Button } from '#/components/ui/button';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Skeleton } from '#/components/ui/skeleton';
import { cn } from '#/lib/utils';
import { EmptyIssues } from './empty-issues';
import {
  useAcknowledgeIssue,
  useActiveOrganizationId,
  useActOnIssue,
  useEvaluateIssues,
  useListAttentionIssues,
  useResolveIssue,
} from './hooks';
import type { AttentionIssue, MonitoringTab, SeverityFilter } from './interfaces';
import { IssueList } from './issue-list';
import { SEVERITY_FILTERS, SEVERITY_LABELS, toDisplaySeverity } from './issue-meta';
import { IssueSummary } from './issue-summary';
import { MonitoringHealth } from './monitoring-health';
import { MonitoringNavigation } from './monitoring-navigation';

export function IssuesView() {
  const activeOrganizationId = useActiveOrganizationId();
  const [tab, setTab] = useState<MonitoringTab>('open');
  const [severity, setSeverity] = useState<SeverityFilter>('all');
  const [query, setQuery] = useState('');
  const [queryDraft, setQueryDraft] = useState('');
  const debounce = useRef<ReturnType<typeof setTimeout> | null>(null);

  const listQuery = useListAttentionIssues();
  const evaluate = useEvaluateIssues();
  const acknowledge = useAcknowledgeIssue();
  const resolve = useResolveIssue();
  const act = useActOnIssue();

  const issues = useMemo(() => listQuery.data?.data ?? [], [listQuery.data]);
  const tabIssues = useMemo(
    () =>
      tab === 'resolved'
        ? issues.filter((issue) => issue.status === 'resolved')
        : issues.filter((issue) => issue.status !== 'resolved'),
    [issues, tab]
  );

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

  const busyIds = useMemo(() => {
    const ids = new Set<string>();
    if (acknowledge.isPending && acknowledge.variables) ids.add(acknowledge.variables);
    if (resolve.isPending && resolve.variables) ids.add(resolve.variables);
    if (act.isPending && act.variables) ids.add(act.variables.issueId);
    return ids;
  }, [
    acknowledge.isPending,
    acknowledge.variables,
    resolve.isPending,
    resolve.variables,
    act.isPending,
    act.variables,
  ]);

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

  const handleAct = (issue: AttentionIssue) => {
    if (!issue.action) return;
    act.mutate({ issueId: issue.id, action: issue.action });
  };

  if (!activeOrganizationId) {
    return (
      <div className="space-y-6">
        <PageHeader
          title="Monitoring"
          description="Health, updates and activity across your projects, servers and domains."
        />
        <div className="flex flex-col items-center rounded-2xl border border-border/60 bg-card px-6 py-14 text-center">
          <span className="flex size-11 items-center justify-center rounded-2xl bg-muted text-muted-foreground">
            <Building2 className="size-5" />
          </span>
          <h3 className="mt-3 font-medium text-[15px]">Select an organization</h3>
          <p className="mx-auto mt-1.5 max-w-md text-[13px] text-muted-foreground/80 leading-relaxed">
            Issues are scoped to an organization. Pick one to load its monitoring feed.
          </p>
          <Button className="mt-5" asChild>
            <Link to="/organizations">Go to organizations</Link>
          </Button>
        </div>
      </div>
    );
  }

  const loading = listQuery.isLoading;
  const loadError = listQuery.isError;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Monitoring"
        description="Health, updates and activity across your projects, servers and domains."
        action={
          tab !== 'health' ? (
            <Button
              variant="outline"
              onClick={() => evaluate.mutate()}
              disabled={evaluate.isPending}
            >
              {evaluate.isPending ? (
                <Loader2 className="size-4 animate-spin" />
              ) : (
                <RefreshCw className="size-4" />
              )}
              {evaluate.isPending ? 'Rescanning' : 'Rescan'}
            </Button>
          ) : undefined
        }
      />

      <MonitoringNavigation value={tab} onChange={setTab} />

      <section
        role="tabpanel"
        id={`monitoring-panel-${tab}`}
        aria-labelledby={`monitoring-tab-${tab}`}
      >
        {!loading && tab === 'resolved' && (
          <p className="mb-4 rounded-xl border border-border/50 bg-muted/25 px-4 py-3 text-[12px] text-muted-foreground leading-relaxed">
            History covers container and server incidents. Other checks report current state, so a
            fixed certificate or deploy simply stops appearing here.
          </p>
        )}

        {tab === 'health' ? (
          <MonitoringHealth />
        ) : loading ? (
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_340px]">
            <div className="min-w-0 space-y-4">
              {[0, 1].map((index) => (
                <div
                  key={index}
                  className="overflow-hidden rounded-2xl border border-border/60 bg-card"
                >
                  <div className="flex items-start gap-3 border-border/60 border-b px-5 py-4">
                    <Skeleton className="size-9 shrink-0 rounded-xl" />
                    <div className="flex-1 space-y-2">
                      <Skeleton className="h-3.5 w-32" />
                      <Skeleton className="h-3 w-56" />
                    </div>
                  </div>
                  <div className="space-y-3 px-5 py-4">
                    <Skeleton className="h-3 w-3/5" />
                    <Skeleton className="h-3 w-2/5" />
                  </div>
                </div>
              ))}
            </div>
            <div className="hidden rounded-2xl border border-border/50 bg-card lg:block">
              <div className="flex items-center gap-3 border-border/50 border-b px-5 py-4">
                <Skeleton className="size-9 shrink-0 rounded-xl" />
                <div className="space-y-2">
                  <Skeleton className="h-3.5 w-24" />
                  <Skeleton className="h-3 w-32" />
                </div>
              </div>
              <div className="space-y-4 p-5">
                <Skeleton className="h-6 w-16" />
                <Skeleton className="h-1.5 w-full" />
                <div className="space-y-2 pt-1">
                  {[0, 1, 2].map((index) => (
                    <Skeleton key={index} className="h-4 w-full" />
                  ))}
                </div>
              </div>
            </div>
          </div>
        ) : loadError ? (
          <QueryErrorState
            title="Issues are unavailable"
            description="Could not load the monitoring feed for this organization."
            onRetry={() => void listQuery.refetch()}
          />
        ) : tabIssues.length === 0 ? (
          <EmptyIssues filtered={false} resolved={tab === 'resolved'} />
        ) : (
          <div className="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_340px]">
            <div className="min-w-0">
              <div className="mb-5 flex flex-col gap-3 sm:flex-row sm:flex-wrap sm:items-center">
                <div className="relative w-full sm:min-w-[220px] sm:flex-1">
                  <Search className="pointer-events-none absolute start-3.5 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
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
                  onAcknowledge={(issue) => acknowledge.mutate(issue.id)}
                  onResolve={(issue) => resolve.mutate(issue.id)}
                  onAct={handleAct}
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
        )}
      </section>
    </div>
  );
}
