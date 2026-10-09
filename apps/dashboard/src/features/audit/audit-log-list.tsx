import { useQuery } from '@tanstack/react-query';
import { Search, Shield, User } from 'lucide-react';
import { useMemo, useState } from 'react';
import { Input } from '#/components/ui/input';
import { QueryErrorState } from '#/components/ui/query-error-state';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '#/components/ui/select';
import { Skeleton } from '#/components/ui/skeleton';
import { Tabs, TabsList, TabsTrigger } from '#/components/ui/tabs';
import type { AuditLog } from '#/interfaces/audit';
import { auditApi } from './api';
import { AuditDetailsDialog } from './audit-details-dialog';
import { AuditSummaryCard } from './audit-summary-card';
import {
  AUDIT_CATEGORIES,
  type AuditCategoryId,
  absoluteTime,
  actorDisplay,
  auditCategoryOf,
  clockTime,
  dayKey,
  describeAuditAction,
  relativeTime,
  resourceDisplay,
  toneDot,
} from './audit-taxonomy';

const PER_PAGE = 50;

type CategoryKey = 'all' | AuditCategoryId;
type PeriodKey = 'all' | 'today' | '7d' | '30d';

function periodStart(period: PeriodKey): number | null {
  if (period === 'all') return null;
  const start = new Date();
  start.setHours(0, 0, 0, 0);
  const back = period === 'today' ? 0 : period === '7d' ? 6 : 29;
  start.setDate(start.getDate() - back);
  return start.getTime();
}

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
  const {
    data: response,
    isLoading,
    isError,
    refetch,
  } = useQuery({
    queryKey: ['auditLogs'],
    queryFn: () => auditApi.list(),
  });
  const [category, setCategory] = useState<CategoryKey>('all');
  const [search, setSearch] = useState('');
  const [actor, setActor] = useState('all');
  const [period, setPeriod] = useState<PeriodKey>('all');
  const [page, setPage] = useState(1);
  const [selected, setSelected] = useState<AuditLog | null>(null);

  const logs = useMemo(() => response?.data ?? [], [response]);

  const actors = useMemo(() => {
    const seen = new Map<string, number>();
    for (const log of logs) {
      const name = actorDisplay(log);
      seen.set(name, (seen.get(name) ?? 0) + 1);
    }
    return [...seen.entries()].sort((a, b) => b[1] - a[1]);
  }, [logs]);

  const counts = useMemo(() => {
    const map = new Map<CategoryKey, number>([['all', logs.length]]);
    for (const log of logs) {
      const key = auditCategoryOf(log);
      map.set(key, (map.get(key) ?? 0) + 1);
    }
    return map;
  }, [logs]);

  const filtered = useMemo(() => {
    const query = search.trim().toLowerCase();
    const from = periodStart(period);
    return logs
      .filter((log) => category === 'all' || auditCategoryOf(log) === category)
      .filter((log) => actor === 'all' || actorDisplay(log) === actor)
      .filter((log) => from === null || new Date(log.createdAt).getTime() >= from)
      .filter(
        (log) =>
          !query ||
          `${log.action} ${log.resource} ${log.details} ${log.userId} ${log.ipAddress}`
            .toLowerCase()
            .includes(query)
      )
      .sort((a, b) => +new Date(b.createdAt) - +new Date(a.createdAt));
  }, [logs, category, actor, period, search]);

  const totalPages = Math.max(1, Math.ceil(filtered.length / PER_PAGE));
  const safePage = Math.min(page, totalPages);
  const pageRows = filtered.slice((safePage - 1) * PER_PAGE, safePage * PER_PAGE);

  const days = useMemo(() => {
    const groups: Array<{ key: string; label: string; rows: AuditLog[] }> = [];
    for (const row of pageRows) {
      const key = dayKey(row.createdAt);
      const last = groups[groups.length - 1];
      if (last && last.key === key) {
        last.rows.push(row);
        continue;
      }
      groups.push({ key, label: dayLabel(row.createdAt), rows: [row] });
    }
    return groups;
  }, [pageRows]);

  const filtersActive =
    category !== 'all' || actor !== 'all' || period !== 'all' || search.trim() !== '';

  const resetPage = () => setPage(1);

  const clearFilters = () => {
    setCategory('all');
    setActor('all');
    setPeriod('all');
    setSearch('');
    setPage(1);
  };

  if (isLoading) {
    return (
      <div className="space-y-3">
        {Array.from({ length: 6 }).map((_, index) => (
          <Skeleton key={index} className="h-16 w-full rounded-2xl" />
        ))}
      </div>
    );
  }

  if (isError) {
    return (
      <QueryErrorState
        title="Audit logs are unavailable"
        description="Codedock could not load activity logs for your workspace."
        onRetry={() => void refetch()}
      />
    );
  }

  return (
    <div>
      <Tabs
        value={category}
        onValueChange={(value) => {
          setCategory(value as CategoryKey);
          resetPage();
        }}
      >
        <TabsList variant="line" className="flex w-full flex-wrap justify-start">
          <TabsTrigger value="all">
            All{counts.get('all') ? ` (${counts.get('all')})` : ''}
          </TabsTrigger>
          {AUDIT_CATEGORIES.filter((cat) => (counts.get(cat.id) ?? 0) > 0).map((cat) => (
            <TabsTrigger key={cat.id} value={cat.id}>
              {cat.label} ({counts.get(cat.id)})
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      <div className="mt-6 grid grid-cols-1 gap-6 lg:grid-cols-[1fr_320px]">
        <div className="min-w-0 space-y-6">
          <div className="flex flex-wrap items-center gap-3">
            <div className="relative min-w-52 flex-1">
              <Search className="absolute top-1/2 left-3 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                type="search"
                value={search}
                onChange={(event) => {
                  setSearch(event.target.value);
                  resetPage();
                }}
                placeholder="Search actions, resources, users…"
                className="ps-10"
              />
            </div>
            <Select
              value={actor}
              onValueChange={(value) => {
                setActor(value);
                resetPage();
              }}
            >
              <SelectTrigger className="w-48">
                <SelectValue placeholder="Anyone" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">Anyone</SelectItem>
                {actors.map(([name, count]) => (
                  <SelectItem key={name} value={name}>
                    {name} ({count})
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            <Select
              value={period}
              onValueChange={(value) => {
                setPeriod(value as PeriodKey);
                resetPage();
              }}
            >
              <SelectTrigger className="w-40">
                <SelectValue placeholder="All time" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All time</SelectItem>
                <SelectItem value="today">Today</SelectItem>
                <SelectItem value="7d">Last 7 days</SelectItem>
                <SelectItem value="30d">Last 30 days</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {filtered.length === 0 ? (
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

          {filtered.length > PER_PAGE && (
            <div className="flex items-center justify-between text-muted-foreground text-sm">
              <span>
                Showing {(safePage - 1) * PER_PAGE + 1}–
                {Math.min(safePage * PER_PAGE, filtered.length)} of {filtered.length}
              </span>
              <div className="flex gap-2">
                <button
                  type="button"
                  onClick={() => setPage(Math.max(1, safePage - 1))}
                  disabled={safePage === 1}
                  className="rounded-lg border border-border/50 px-3 py-1.5 transition-colors hover:bg-muted/40 disabled:opacity-50"
                >
                  Previous
                </button>
                <button
                  type="button"
                  onClick={() => setPage(Math.min(totalPages, safePage + 1))}
                  disabled={safePage >= totalPages}
                  className="rounded-lg border border-border/50 px-3 py-1.5 transition-colors hover:bg-muted/40 disabled:opacity-50"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>

        <aside className="space-y-4 lg:sticky lg:top-6 lg:self-start">
          <AuditSummaryCard logs={logs} />
        </aside>
      </div>

      <AuditDetailsDialog event={selected} onClose={() => setSelected(null)} />
    </div>
  );
}
