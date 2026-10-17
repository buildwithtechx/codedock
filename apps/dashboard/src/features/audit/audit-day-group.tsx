import { User } from 'lucide-react';
import type { AuditEventRow } from './api';
import { resolveResourceDisplay } from './audit-resource-helper';
import type { AuditTone } from '#/interfaces/audit';
import { actorDisplay, describeAuditEvent } from './audit-taxonomy';

const TONE_DOT: Record<AuditTone, string> = {
  neutral: 'bg-muted-foreground/60',
  info: 'bg-blue-500',
  success: 'bg-emerald-500',
  warning: 'bg-amber-500',
  danger: 'bg-rose-500',
};

function formatDayHeader(dateIso?: string): string {
  if (!dateIso) return '';
  const date = new Date(dateIso);
  if (Number.isNaN(date.getTime())) return dateIso;
  return date
    .toLocaleDateString(undefined, {
      weekday: 'short',
      month: 'short',
      day: 'numeric',
    })
    .toUpperCase();
}

function formatClockTime(dateIso?: string): string {
  if (!dateIso) return '';
  const date = new Date(dateIso);
  if (Number.isNaN(date.getTime())) return '';
  return date.toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
  });
}

function formatRelativeTime(dateIso?: string): string {
  if (!dateIso) return '';
  const date = new Date(dateIso);
  if (Number.isNaN(date.getTime())) return '';
  const diffMs = Date.now() - date.getTime();
  const mins = Math.floor(diffMs / 60000);
  if (mins < 1) return 'just now';
  if (mins < 60) return `${mins}m ago`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
}

function resolveActionName(event?: Partial<AuditEventRow>): string {
  if (!event) return 'system.event';
  return event.eventType || event.action || 'system.event';
}

function resolveEventSentence(event: AuditEventRow): string {
  const actionName = resolveActionName(event);
  const desc = describeAuditEvent(actionName);
  const actor = actorDisplay(event);
  const resource = resolveResourceDisplay(event, actionName);
  const action = desc.action || actionName.replace('.', ' ');

  return resource ? `${actor} ${action} ${resource}` : `${actor} ${action}`;
}

interface AuditDayGroupProps {
  dayKey: string;
  events: AuditEventRow[];
  onSelectEvent: (event: AuditEventRow) => void;
}

export function AuditDayGroup({ dayKey, events, onSelectEvent }: AuditDayGroupProps) {
  const headerDate = events[0]?.createdAt ?? dayKey;

  return (
    <section className="space-y-2">
      <h3 className="px-1 font-semibold text-[11px] text-muted-foreground uppercase tracking-wider">
        {formatDayHeader(headerDate)}
      </h3>

      <div className="divide-y divide-border/40 overflow-hidden rounded-2xl border border-border/50 bg-card">
        {events.map((event) => {
          const actionName = resolveActionName(event);
          const desc = describeAuditEvent(actionName);
          const tone = desc.tone ?? 'neutral';

          return (
            <button
              key={event.id}
              type="button"
              onClick={() => onSelectEvent(event)}
              className="flex w-full items-center gap-3.5 px-4 py-3 text-start transition-colors hover:bg-muted/40 focus-visible:outline-2 focus-visible:outline-ring"
            >
              <div className="relative flex size-8 shrink-0 items-center justify-center rounded-full bg-muted/60 text-muted-foreground">
                <User className="size-4" />
                <span
                  className={`absolute -right-0.5 -bottom-0.5 size-2.5 rounded-full ring-2 ring-card ${TONE_DOT[tone]}`}
                />
              </div>

              <div className="min-w-0 flex-1">
                <p className="truncate font-medium text-foreground text-sm">
                  {resolveEventSentence(event)}
                </p>
                <p className="mt-0.5 text-muted-foreground text-xs">
                  {formatClockTime(event.createdAt)}
                </p>
              </div>

              <span className="shrink-0 text-muted-foreground text-xs">
                {formatRelativeTime(event.createdAt)}
              </span>
            </button>
          );
        })}
      </div>
    </section>
  );
}
