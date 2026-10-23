import { Layers } from 'lucide-react';
import { useCallback, useMemo, useState } from 'react';
import { toast } from 'sonner';
import { Button } from '@codedock/ui/button';
import { AUDIT_CATEGORIES } from '#/constants/audit';
import type { AuditEventRow, AuditFacets } from '#/interfaces/audit';
import { useAuditFacets, useAuditLogs, useUpdateAuditSettings } from '#/hooks/use-audit-logs';
import { AuditDayGroup } from './audit-day-group';
import { AuditDetailsDialog } from './audit-details-dialog';
import { AuditFilterBar, type PeriodKey, type SourceKey } from './audit-filter-bar';
import { CATEGORY_ICONS, type CategoryKey, groupByDay, periodStart } from './audit-helpers';
import { AuditRecordingCard } from './audit-recording-card';

export function AuditLogList() {
  const [category, setCategory] = useState<CategoryKey>('all');
  const [search, setSearch] = useState('');
  const [actor, setActor] = useState('');
  const [period, setPeriod] = useState<PeriodKey>('all');
  const [source, setSource] = useState<SourceKey>('all');
  const [page, setPage] = useState(1);

  const [selectedEvent, setSelectedEvent] = useState<AuditEventRow | null>(null);
  const [savingSettings, setSavingSettings] = useState(false);
  const updateSettingsMutation = useUpdateAuditSettings();

  const query = useMemo(
    () => ({
      category: category === 'all' ? undefined : category,
      source: source === 'all' ? undefined : source,
      actorUserId: actor || undefined,
      from: periodStart(period),
      q: search.trim() || undefined,
    }),
    [category, source, actor, period, search]
  );

  const { data: logsRes, isLoading: loading } = useAuditLogs({ ...query, page, perPage: 50 });
  const { data: facetsRes } = useAuditFacets(query);

  const rawList = logsRes?.data;
  const events = Array.isArray(rawList)
    ? rawList
    : Array.isArray(logsRes)
      ? (logsRes as AuditEventRow[])
      : [];
  const total = (logsRes as unknown as { total?: number })?.total ?? events.length;
  const facets = facetsRes?.data ?? (facetsRes as unknown as AuditFacets) ?? null;

  const handleToggleRecording = useCallback(
    async (enabled: boolean) => {
      setSavingSettings(true);
      try {
        await updateSettingsMutation.mutateAsync({ enabled });
        toast.success(enabled ? 'Activity recording enabled' : 'Activity recording disabled');
      } catch {
        toast.error('Failed to update recording setting');
      } finally {
        setSavingSettings(false);
      }
    },
    [updateSettingsMutation]
  );

  const handleRetentionChange = useCallback(
    async (retentionDays: number) => {
      setSavingSettings(true);
      try {
        await updateSettingsMutation.mutateAsync({ retentionDays });
        toast.success(`Retention updated to ${retentionDays} days`);
      } catch {
        toast.error('Failed to update retention days');
      } finally {
        setSavingSettings(false);
      }
    },
    [updateSettingsMutation]
  );

  const dayGroups = useMemo(() => groupByDay(events), [events]);
  const counts = useMemo(
    () => new Map(facets?.categories?.map((c) => [c.id, c.count]) ?? []),
    [facets]
  );

  const settings = facets?.settings ?? { enabled: true, retentionDays: 90 };
  const canManage = facets?.canManage ?? true;

  return (
    <div className="space-y-6">
      <header>
        <h1 className="font-medium text-2xl text-foreground tracking-tight">Audit log</h1>
        <p className="mt-1 text-muted-foreground text-sm">
          Everything that happened in this organization, who did it, and when.
        </p>
      </header>

      <div className="flex items-center gap-1 overflow-x-auto border-border/50 border-b pb-px">
        <button
          type="button"
          onClick={() => {
            setCategory('all');
            setPage(1);
          }}
          className={`relative inline-flex items-center gap-2 whitespace-nowrap px-3 py-2.5 font-medium text-sm transition-colors ${
            category === 'all'
              ? 'text-foreground'
              : 'text-muted-foreground hover:text-foreground/80'
          }`}
        >
          All {facets?.total !== undefined ? facets.total : ''}
          {category === 'all' && (
            <span className="absolute inset-x-0 -bottom-px h-0.5 rounded-full bg-primary" />
          )}
        </button>

        {AUDIT_CATEGORIES.map((cat) => {
          const Icon = CATEGORY_ICONS[cat.id] ?? Layers;
          const count = counts.get(cat.id) ?? 0;
          const isSelected = category === cat.id;

          return (
            <button
              key={cat.id}
              type="button"
              onClick={() => {
                setCategory(cat.id);
                setPage(1);
              }}
              className={`relative inline-flex items-center gap-2 whitespace-nowrap px-3 py-2.5 font-medium text-sm transition-colors ${
                isSelected ? 'text-foreground' : 'text-muted-foreground hover:text-foreground/80'
              }`}
            >
              <Icon className="size-4" />
              {cat.label} {count}
              {isSelected && (
                <span className="absolute inset-x-0 -bottom-px h-0.5 rounded-full bg-primary" />
              )}
            </button>
          );
        })}
      </div>

      <div className="space-y-6">
        <AuditRecordingCard
          enabled={settings.enabled}
          retentionDays={settings.retentionDays}
          saving={savingSettings}
          canManage={canManage}
          onToggle={handleToggleRecording}
          onRetentionChange={handleRetentionChange}
        />

        <div className="min-w-0 space-y-6">
          <AuditFilterBar
            search={search}
            onSearchChange={(val) => {
              setSearch(val);
              setPage(1);
            }}
            actor={actor}
            onActorChange={(val) => {
              setActor(val);
              setPage(1);
            }}
            period={period}
            onPeriodChange={(val) => {
              setPeriod(val);
              setPage(1);
            }}
            source={source}
            onSourceChange={(val) => {
              setSource(val);
              setPage(1);
            }}
            facets={facets}
          />

          {loading && events.length === 0 ? (
            <div className="h-64 animate-pulse rounded-2xl bg-card" />
          ) : Object.keys(dayGroups).length > 0 ? (
            Object.entries(dayGroups).map(([dKey, evts]) => (
              <AuditDayGroup
                key={dKey}
                dayKey={dKey}
                events={evts}
                onSelectEvent={(evt) => setSelectedEvent(evt)}
              />
            ))
          ) : (
            <div className="rounded-2xl border border-border/50 bg-card p-10 text-center">
              <p className="font-medium text-foreground text-sm">No events found</p>
              <p className="mt-1 text-muted-foreground text-xs">
                No activity matched your selected filter parameters.
              </p>
            </div>
          )}

          {total > 50 && (
            <div className="flex items-center justify-between border-border/40 border-t pt-4">
              <Button
                variant="outline"
                size="sm"
                disabled={page <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
              >
                Previous
              </Button>
              <span className="text-muted-foreground text-xs">
                Page {page} of {Math.ceil(total / 50)}
              </span>
              <Button
                variant="outline"
                size="sm"
                disabled={page * 50 >= total}
                onClick={() => setPage((p) => p + 1)}
              >
                Next
              </Button>
            </div>
          )}
        </div>
      </div>

      {selectedEvent && (
        <AuditDetailsDialog
          event={selectedEvent}
          open={!!selectedEvent}
          onClose={() => setSelectedEvent(null)}
        />
      )}
    </div>
  );
}
