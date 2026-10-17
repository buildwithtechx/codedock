import { Link, useNavigate, useSearch } from '@tanstack/react-router';
import { Building2, Loader2, RefreshCw, Trash2 } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Button } from '@codedock/ui/button';
import { QueryErrorState } from '#/components/query-error-state';
import { Skeleton } from '@codedock/ui/skeleton';
import { DeploymentDirectory } from '#/features/dashboard/deployment-directory';
import { NotificationsSettings } from '#/features/notifications/notifications-settings';
import {
  useAcknowledgeIssue,
  useActiveOrganizationId,
  useActOnIssue,
  useClearIssueHistory,
  useEvaluateIssues,
  useListAttentionIssues,
  useResolveIssue,
} from './hooks';
import type { AttentionIssue, MonitoringTab } from './interfaces';
import { IssuesFeedView } from './issues-feed-view';
import { MonitoringHealth } from './monitoring-health';
import { MonitoringNavigation } from './monitoring-navigation';

function useMonitoringTab(): [MonitoringTab, (tab: MonitoringTab) => void] {
  let urlTab: MonitoringTab | undefined;
  let nav: ((opts: any) => void) | undefined;
  try {
    const s = useSearch({ strict: false }) as { tab?: MonitoringTab };
    if (s?.tab && ['open', 'health', 'resolved', 'notifications', 'deployments'].includes(s.tab)) {
      urlTab = s.tab;
    }
    nav = useNavigate();
  } catch {}
  const [localTab, setLocalTab] = useState<MonitoringTab>(urlTab || 'open');
  const activeTab = urlTab || localTab;
  const setTab = (t: MonitoringTab) => {
    setLocalTab(t);
    try {
      nav?.({ to: '/monitoring', search: (prev: any) => ({ ...prev, tab: t }), replace: true });
    } catch {}
  };
  return [activeTab, setTab];
}

export function IssuesView() {
  const activeOrganizationId = useActiveOrganizationId();
  const [tab, setTab] = useMonitoringTab();

  const listQuery = useListAttentionIssues();
  const evaluate = useEvaluateIssues();
  const acknowledge = useAcknowledgeIssue();
  const resolve = useResolveIssue();
  const act = useActOnIssue();
  const clearHistory = useClearIssueHistory();
  const [confirmingClear, setConfirmingClear] = useState(false);
  const [clearError, setClearError] = useState<string | null>(null);

  const feed = listQuery.data?.data;
  const issues = useMemo(() => feed?.issues ?? [], [feed]);
  const canClearHistory = feed?.canClearHistory ?? false;
  const tabIssues = useMemo(
    () =>
      tab === 'resolved'
        ? issues.filter((issue) => issue.status === 'resolved')
        : issues.filter((issue) => issue.status !== 'resolved'),
    [issues, tab]
  );

  useEffect(() => {
    setConfirmingClear(false);
    setClearError(null);
  }, [activeOrganizationId, tab]);

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

  const handleAct = (issue: AttentionIssue) => {
    if (!issue.action) return;
    act.mutate({ issueId: issue.id, action: issue.action });
  };

  const runClearHistory = () => {
    clearHistory.mutate(undefined, {
      onSuccess: () => setConfirmingClear(false),
      onError: (error) =>
        setClearError(
          error instanceof Error ? error.message : 'Could not clear monitoring history'
        ),
    });
  };

  const showRescan = tab === 'open' || tab === 'resolved';
  const showClearHistory = tab === 'resolved' && canClearHistory && tabIssues.length > 0;
  const headerAction =
    showRescan || showClearHistory ? (
      <div className="flex items-center gap-2">
        {showRescan && (
          <Button variant="outline" onClick={() => evaluate.mutate()} disabled={evaluate.isPending}>
            {evaluate.isPending ? (
              <Loader2 className="size-4 animate-spin" />
            ) : (
              <RefreshCw className="size-4" />
            )}
            {evaluate.isPending ? 'Rescanning' : 'Rescan'}
          </Button>
        )}
        {showClearHistory && (
          <Button
            variant="destructive"
            onClick={() => {
              setClearError(null);
              setConfirmingClear(true);
            }}
            disabled={clearHistory.isPending}
          >
            <Trash2 className="size-4" />
            Clear history
          </Button>
        )}
      </div>
    ) : undefined;

  if (!activeOrganizationId) {
    return (
      <div className="space-y-6">
        <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
          <div>
            <h1 className="font-medium text-2xl text-foreground/80 tracking-tight">Monitoring</h1>
            <p className="mt-1 text-muted-foreground/70 text-sm">
              Health, updates and activity across your projects, servers and domains.
            </p>
          </div>
        </div>
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
      <div className="mb-6 flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="font-medium text-2xl text-foreground/80 tracking-tight">Monitoring</h1>
          <p className="mt-1 text-muted-foreground/70 text-sm">
            Health, updates and activity across your projects, servers and domains.
          </p>
        </div>
        {headerAction}
      </div>

      <MonitoringNavigation value={tab} onChange={setTab} />

      <section
        role="tabpanel"
        id={`monitoring-panel-${tab}`}
        aria-labelledby={`monitoring-tab-${tab}`}
      >
        {tab === 'deployments' ? (
          <DeploymentDirectory hideHeader />
        ) : tab === 'notifications' ? (
          <NotificationsSettings />
        ) : tab === 'health' ? (
          <MonitoringHealth />
        ) : (
          <>
            {!loading && tab === 'resolved' && (
              <p className="mb-4 rounded-xl border border-border/50 bg-muted/25 px-4 py-3 text-[12px] text-muted-foreground leading-relaxed">
                History covers container and server incidents. Other checks report current state, so
                a fixed certificate or deploy simply stops appearing here.
              </p>
            )}

            {tab === 'resolved' && confirmingClear && canClearHistory && (
              <section
                aria-label="Clear monitoring history"
                className="mb-5 space-y-3 rounded-2xl border border-border/50 bg-card p-5"
              >
                <p className="text-foreground text-sm">
                  Clear resolved monitoring history for this organization?
                </p>
                <p className="text-muted-foreground text-sm">
                  This permanently removes every resolved issue. Open and acknowledged issues are
                  never deleted.
                </p>
                {clearError && (
                  <p role="alert" className="text-destructive text-sm">
                    {clearError}
                  </p>
                )}
                <div className="flex flex-wrap gap-2">
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={runClearHistory}
                    disabled={clearHistory.isPending}
                  >
                    {clearHistory.isPending && <Loader2 className="size-4 animate-spin" />}
                    {clearHistory.isPending ? 'Clearing' : 'Clear history'}
                  </Button>
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => setConfirmingClear(false)}
                    disabled={clearHistory.isPending}
                  >
                    Cancel
                  </Button>
                </div>
              </section>
            )}

            {loading ? (
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
            ) : (
              <IssuesFeedView
                tabIssues={tabIssues}
                tab={tab}
                busyIds={busyIds}
                onAcknowledge={(id) => acknowledge.mutate(id)}
                onResolve={(id) => resolve.mutate(id)}
                onAct={handleAct}
              />
            )}
          </>
        )}
      </section>
    </div>
  );
}
