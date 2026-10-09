import { useQuery } from '@tanstack/react-query';
import { Search, Shield, User } from 'lucide-react';
import { useMemo, useState } from 'react';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import { Skeleton } from '#/components/ui/skeleton';
import { Tabs, TabsList, TabsTrigger } from '#/components/ui/tabs';
import { type AuditLogRow, auditApi } from './api';
import { AuditDetailsDialog } from './audit-details-dialog';
import { AuditSummaryCard } from './audit-summary-card';
import {
  AUDIT_CATEGORIES,
  absoluteTime,
  actorDisplay,
  clockTime,
  dayKey,
  describeAuditAction,
  relativeTime,
  resourceDisplay,
  toneDot,
} from './audit-taxonomy';

const PER_PAGE = 50;

function dayLabel(iso: string): string {
  const key = dayKey(iso);
  const today = dayKey(new Date().toISOString());
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);
  if (key === today) return 'Today';
  if (key === dayKey(yesterday.toISOString())) return 'Yesterday';
  return new Date(iso).toLocaleDateString(undefined, {
    weekday: 'short',
    month: 'short',
    day: 'numeric',
  });
}

export function AuditLogList() {
  const [category, setCategory] = useState('all');
  const [search, setSearch] = useState('');
  const [page, setPage] = useState(1);
  const [selected, setSelected] = useState<AuditLogRow | null>(null);

  const facetsQuery = useQuery({
    queryKey: ['auditLogFacets'],
    queryFn: () => auditApi.facets(),
  });
  const facets = facetsQuery.data?.data ?? null;

  const listQuery = useQuery({
    queryKey: ['auditLogs', category, page],
    queryFn: () =>
      auditApi.list({
        category: category === 'all' ? undefined : category,
        limit: PER_PAGE,
        offset: (page - 1) * PER_PAGE,
      }),
  });
  const rows = useMemo(() => listQuery.data?.data ?? [], [listQuery.data]);

  const tabDefs = useMemo(() => {
    if (facets) {
      return [
        { id: 'all', label: 'All', count: facets.total },
        ...facets.categories.map((entry) => ({
          id: entry.id,
          label: entry.label,
          count: entry.count,
        })),
      ];
    }
    return [
      { id: 'all', label: 'All', count: 0 },
      ...AUDIT_CATEGORIES.map((entry) => ({ id: entry.id, label: entry.label, count: 0 })),
    ];
  }, [facets]);

  const total =
    category === 'all'
      ? (facets?.total ?? rows.length)
      : (facets?.categories.find((entry) => entry.id === category)?.count ?? rows.length);

  const searched = useMemo(() => {
    const query = search.trim().toLowerCase();
    const matched = query
      ? rows.filter((log) =>
          `${log.action} ${log.resource} ${log.details} ${log.userId} ${log.ipAddress}`
            .toLowerCase()
            .includes(query)
        )
      : rows;
    return [...matched].sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt));
  }, [rows, search]);

  const days = useMemo(() => {
    const groups: Array<{ key: string; label: string; rows: AuditLogRow[] }> = [];
    for (const row of searched) {
      const key = dayKey(row.createdAt);
      const last = groups[groups.length - 1];
      if (last && last.key === key) {
        last.rows.push(row);
        continue;
      }
      groups.push({ key, label: dayLabel(row.createdAt), rows: [row] });
    }
    return groups;
  }, [searched]);

  const totalPages = Math.max(1, Math.ceil(total / PER_PAGE));
  const filtersActive = category !== 'all' || search.trim() !== '';

  const pickCategory = (value: string) => {
    setCategory(value);
    setPage(1);
  };

  const clearFilters = () => {
    setCategory('all');
    setSearch('');
    setPage(1);
  };

  if (listQuery.isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 6 }).map((_, index) => (
          <Skeleton key={index} className="h-16 w-full rounded-2xl" />
        ))}
      </div>
    );
  }

  if (listQuery.isError) {
    return (
      <QueryErrorState
        title="Audit logs are unavailable"
        description="Codedock could not load activity logs for your workspace."
        onRetry={() => void listQuery.refetch()}
      />
    );
  }

  return (
    <div>
      <Tabs value={category} onValueChange={pickCategory}>
        <TabsList variant="line" className="flex w-full flex-wrap justify-start">
          {tabDefs.map((tab) => (
            <TabsTrigger key={tab.id} value={tab.id}>
              {tab.label}
              {facets && <span className="text-muted-foreground tabular-nums">({tab.count})</span>}
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      <div className="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-[1fr_320px]">
        <div className="min-w-0 space-y-6">
          <div className="relative min-w-52 flex-1">
            <Search className="absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
            <Input
              type="search"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              placeholder="Search actions, resources, users…"
              className="ps-10"
            />
          </div>

          {searched.length === 0 ? (
            <div className="rounded-2xl border border-border/50 bg-card py-16 text-center">
              <Shield className="mx-auto mb-4 size-8 text-muted-foreground/30" />
              <p className="text-muted-foreground text-sm">
                {filtersActive ? 'No events match these filters.' : 'No audit events recorded yet.'}
              </p>
              {filtersActive && (
                <button
                  type="button"
                  onClick={clearFilters}
                  className="mt-3 font-medium text-primary text-xs hover:underline"
                >
                  Clear filters
                </button>
              )}
            </div>
          ) : (
            <div className="space-y-4">
              {days.map((day) => (
                <div
                  key={day.key}
                  className="overflow-hidden rounded-2xl border border-border/50 bg-card"
                >
                  <div className="border-border/40 border-b px-5 py-2.5">
                    <span className="font-semibold text-muted-foreground text-xs uppercase tracking-wide">
                      {day.label}
                    </span>
                  </div>
                  <div className="divide-y divide-border/40">
                    {day.rows.map((log) => {
                      const info = describeAuditAction(log.action);
                      const actor = actorDisplay(log);
                      const resource = resourceDisplay(log);
                      return (
                        <button
                          key={log.id}
                          type="button"
                          onClick={() => setSelected(log)}
                          className="flex w-full items-center gap-4 px-5 py-3.5 text-start transition-colors hover:bg-muted/30 focus:bg-muted/40 focus:outline-none"
                        >
                          <span className="relative flex size-8 shrink-0 items-center justify-center rounded-full bg-muted/60">
                            <User className="size-3.5 text-muted-foreground" />
                            <span
                              className={`absolute -end-0.5 -bottom-0.5 size-2.5 rounded-full ring-2 ring-card ${toneDot(info.tone)}`}
                            />
                          </span>
                          <div className="min-w-0 flex-1">
                            <p className="truncate text-foreground text-sm">
                              <span className="font-medium">{actor}</span> {info.action}
                              {resource && <span className="font-medium"> {resource}</span>}
                            </p>
                            <p className="mt-0.5 text-muted-foreground text-xs">
                              {clockTime(log.createdAt)}
                            </p>
                          </div>
                          <time
                            className="shrink-0 text-muted-foreground text-xs"
                            title={absoluteTime(log.createdAt)}
                          >
                            {relativeTime(log.createdAt)}
                          </time>
                        </button>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>
          )}

          {total > PER_PAGE && (
            <div className="flex items-center justify-between text-muted-foreground text-sm">
              <span>
                Showing {(page - 1) * PER_PAGE + 1}–{Math.min(page * PER_PAGE, total)} of {total}
              </span>
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={() => setPage(Math.max(1, page - 1))}
                  disabled={page === 1 || listQuery.isFetching}
                  className="rounded-lg border border-border/50 px-3 py-1.5 transition-colors hover:bg-muted/40 disabled:opacity-50"
                >
                  Previous
                </button>
                <button
                  type="button"
                  onClick={() => setPage(Math.min(totalPages, page + 1))}
                  disabled={page >= totalPages || listQuery.isFetching}
                  className="rounded-lg border border-border/50 px-3 py-1.5 transition-colors hover:bg-muted/40 disabled:opacity-50"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>

        <aside className="space-y-4 lg:sticky lg:top-6 lg:self-start">
          <AuditSummaryCard
            total={facets?.total ?? rows.length}
            facets={facets?.categories ?? []}
          />
        </aside>
      </div>

      <AuditDetailsDialog event={selected} onClose={() => setSelected(null)} />
    </div>
  );
}
