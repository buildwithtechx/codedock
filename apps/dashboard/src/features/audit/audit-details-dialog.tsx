import { Check, ChevronDown, Copy } from 'lucide-react';
import { useMemo, useState } from 'react';
import { Button } from '#/components/ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '#/components/ui/dialog';
import type { AuditLog } from '#/interfaces/audit';
import type { AuditLogRow } from './api';
import {
  absoluteTime,
  actorDisplay,
  auditCategoryOf,
  categoryLabel,
  describeAuditAction,
  relativeTime,
  resourceDisplay,
  toneDot,
} from './audit-taxonomy';

function CopyButton({ value }: { value: string }) {
  const [copied, setCopied] = useState(false);
  return (
    <button
      type="button"
      onClick={() => {
        if (typeof navigator !== 'undefined' && navigator.clipboard) {
          void navigator.clipboard.writeText(value).then(() => {
            setCopied(true);
            setTimeout(() => setCopied(false), 1200);
          });
        }
      }}
      className="inline-flex items-center gap-1 rounded-md px-1.5 py-0.5 text-[10px] text-muted-foreground transition-colors hover:bg-muted/60 hover:text-foreground"
      aria-label={copied ? 'Copied' : 'Copy'}
      title={copied ? 'Copied' : 'Copy'}
    >
      {copied ? <Check className="size-3" /> : <Copy className="size-3" />}
    </button>
  );
}

function DetailRow({
  label,
  value,
  copyable,
  mono,
}: {
  label: string;
  value: string;
  copyable?: boolean;
  mono?: boolean;
}) {
  return (
    <div className="flex items-start gap-3 py-1.5">
      <span className="w-24 shrink-0 text-muted-foreground text-xs">{label}</span>
      <span
        className={`flex-1 break-all text-foreground text-sm ${mono ? 'font-mono text-xs' : ''}`}
      >
        {value}
      </span>
      {copyable && value !== '—' && <CopyButton value={value} />}
    </div>
  );
}

function JsonBlock({ value, label }: { value: unknown; label: string }) {
  const text = useMemo(() => {
    try {
      return JSON.stringify(value, null, 2) ?? 'null';
    } catch {
      return String(value);
    }
  }, [value]);
  return (
    <div>
      <div className="mb-1 flex items-center justify-between">
        <span className="font-semibold text-muted-foreground text-xs uppercase tracking-wide">
          {label}
        </span>
        <CopyButton value={text} />
      </div>
      <pre className="max-h-60 overflow-auto whitespace-pre-wrap break-words rounded-lg border border-border/50 bg-muted/30 p-3 font-mono text-foreground text-xs leading-relaxed">
        {text}
      </pre>
    </div>
  );
}

export interface AuditDetailsDialogProps {
  event: AuditLogRow | AuditLog | null;
  open?: boolean;
  onClose: () => void;
}

export function AuditDetailsDialog({ event, open, onClose }: AuditDetailsDialogProps) {
  const [showTechnical, setShowTechnical] = useState(false);
  const info = event ? describeAuditAction(event.action) : null;
  const actor = event ? actorDisplay(event) : '';
  const resource = event ? resourceDisplay(event) : '';

  const row = event as AuditLogRow | null;
  const hasBefore = row?.before !== undefined && row?.before !== null;
  const hasAfter = row?.after !== undefined && row?.after !== null;
  const isOpen = open !== undefined ? open : event !== null;

  return (
    <Dialog open={isOpen} onOpenChange={(v) => !v && onClose()}>
      <DialogContent className="max-h-[85vh] max-w-2xl overflow-y-auto">
        {event && info && (
          <div className="flex flex-col gap-5">
            <DialogHeader className="text-left">
              <div className="flex items-center gap-2">
                <span className={`size-2 shrink-0 rounded-full ${toneDot(info.tone)}`} />
                <DialogTitle className="text-lg">{info.label}</DialogTitle>
              </div>
              <p className="mt-1.5 text-foreground/80 text-sm">
                <span className="font-medium">{actor}</span> {info.action}
                {resource && <span className="font-medium"> {resource}</span>}
              </p>
              <p className="text-muted-foreground text-xs" title={absoluteTime(event.createdAt)}>
                {relativeTime(event.createdAt)} · {absoluteTime(event.createdAt)}
              </p>
            </DialogHeader>

            <section className="rounded-xl bg-card p-4">
              <DetailRow label="Who" value={actor} />
              {row?.actorEmail && <DetailRow label="Email" value={row.actorEmail} />}
              <DetailRow
                label="Came from"
                value={row?.source ? row.source.toUpperCase() : 'Dashboard'}
              />
              {row?.sourceClientName && (
                <DetailRow label="Via agent" value={row.sourceClientName} />
              )}
              <DetailRow label="Category" value={categoryLabel(auditCategoryOf(event))} />
              <DetailRow label="When" value={absoluteTime(event.createdAt)} />
              <DetailRow label="IP address" value={event.ipAddress || '—'} copyable mono />
              <DetailRow label="Resource" value={resource || '—'} />
              {event.details && <DetailRow label="Details" value={event.details} />}
            </section>

            <section className="rounded-xl bg-card">
              <button
                type="button"
                onClick={() => setShowTechnical((value) => !value)}
                className="flex w-full items-center justify-between gap-2 px-4 py-3 text-start"
                aria-expanded={showTechnical}
              >
                <span className="font-semibold text-muted-foreground text-xs uppercase tracking-wide">
                  Technical details
                </span>
                <ChevronDown
                  className={`size-4 text-muted-foreground transition-transform ${showTechnical ? 'rotate-180' : ''}`}
                />
              </button>
              {showTechnical && (
                <div className="space-y-4 border-border/40 border-t px-4 pt-2 pb-4">
                  <div>
                    <DetailRow label="Event ID" value={event.id} copyable mono />
                    <DetailRow label="Action" value={event.action} copyable mono />
                    <DetailRow label="Resource" value={event.resource || '—'} mono />
                    <DetailRow label="User ID" value={event.userId || '—'} copyable mono />
                  </div>

                  {(hasBefore || hasAfter) && (
                    <div className="border-border/40 border-t pt-3">
                      <h4 className="mb-2 font-semibold text-muted-foreground text-xs uppercase tracking-wide">
                        Payload Changes
                      </h4>
                      {hasBefore && hasAfter ? (
                        <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                          <JsonBlock value={row?.before} label="Before" />
                          <JsonBlock value={row?.after} label="After" />
                        </div>
                      ) : hasAfter ? (
                        <JsonBlock value={row?.after} label="After" />
                      ) : (
                        <JsonBlock value={row?.before} label="Removed" />
                      )}
                    </div>
                  )}
                </div>
              )}
            </section>

            <div className="flex justify-end">
              <Button variant="outline" onClick={onClose}>
                Close
              </Button>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
