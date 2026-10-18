import { AlertTriangle, CheckCircle2, Loader2, ShieldAlert, XCircle } from 'lucide-react';
import { Button } from '@codedock/ui/button';
import { Input } from '@codedock/ui/input';
import { cn } from '#/lib/utils';
import { formatBytes } from './destination-display';
import type { BackupRecord, RestoreOperation } from '#/interfaces/backups';

export type WizardStep = 'review' | 'preparing' | 'prepared' | 'applying' | 'done';

export const stepOrder: WizardStep[] = ['review', 'preparing', 'prepared', 'applying', 'done'];
export const stepLabels: Record<WizardStep, string> = {
  review: 'Review',
  preparing: 'Prepare',
  prepared: 'Confirm',
  applying: 'Apply',
  done: 'Done',
};

export function isTerminal(status: string): boolean {
  return [
    'COMPLETED',
    'FAILED',
    'ERROR',
    'CANCELLED',
    'CANCELED',
    'EXPIRED',
    'INTERRUPTED',
  ].includes(status);
}

export function terminalTone(status: string): 'success' | 'failed' | 'cancelled' {
  if (status === 'COMPLETED') return 'success';
  if (status === 'CANCELLED' || status === 'CANCELED' || status === 'INTERRUPTED') {
    return 'cancelled';
  }
  return 'failed';
}

export function RestoreStepIndicator({ currentStep }: { currentStep: WizardStep }) {
  const currentIndex = stepOrder.indexOf(currentStep);

  return (
    <div className="mt-2 flex items-center gap-2">
      {stepOrder.map((item, index) => (
        <div key={item} className="flex min-w-0 flex-1 items-center gap-1.5 last:flex-none">
          <span
            className={cn(
              'flex h-5 w-5 shrink-0 items-center justify-center rounded-full text-[10px]',
              index < currentIndex && 'bg-emerald-500/15 text-emerald-500',
              index === currentIndex && 'bg-primary text-primary-foreground',
              index > currentIndex && 'bg-muted text-muted-foreground'
            )}
          >
            {index + 1}
          </span>
          <span
            className={cn(
              'truncate text-[11px]',
              index <= currentIndex ? 'text-foreground' : 'text-muted-foreground/60'
            )}
          >
            {stepLabels[item]}
          </span>
          {index < stepOrder.length - 1 && <span className="h-px min-w-2 flex-1 bg-border/50" />}
        </div>
      ))}
    </div>
  );
}

export function RestoreReviewStep({
  record,
  takenAt,
  busy,
  onClose,
  onStartPrepare,
}: {
  record: BackupRecord;
  takenAt: string;
  busy: boolean;
  onClose: () => void;
  onStartPrepare: () => void;
}) {
  return (
    <div className="space-y-4">
      <div className="rounded-xl bg-muted/40 p-4 text-sm">
        <p className="text-foreground/80">
          Snapshot taken <strong>{takenAt}</strong> ·{' '}
          {record.fileSizeBytes ? formatBytes(record.fileSizeBytes) : '—'}
        </p>
        <p className="mt-1 truncate font-mono text-muted-foreground text-xs">
          {record.s3Url || record.filePath || record.id}
        </p>
      </div>
      <div className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-4">
        <div className="flex items-start gap-2">
          <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0 text-amber-500" />
          <div className="text-sm">
            <p className="font-medium">This overwrites live data</p>
            <p className="mt-1 text-muted-foreground text-xs">
              Current tables and stored data are replaced with the snapshot contents. The restore is
              verified before anything is applied.
            </p>
          </div>
        </div>
      </div>
      <div className="flex items-center justify-end gap-2">
        <Button variant="ghost" onClick={onClose} disabled={busy}>
          Cancel
        </Button>
        <Button onClick={onStartPrepare} disabled={busy}>
          {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
          Verify and continue
        </Button>
      </div>
    </div>
  );
}

export function RestorePreparedStep({
  sourceName,
  typed,
  busy,
  onTypedChange,
  onCancel,
  onApply,
}: {
  sourceName: string;
  typed: string;
  busy: boolean;
  onTypedChange: (val: string) => void;
  onCancel: () => void;
  onApply: () => void;
}) {
  return (
    <div className="space-y-4">
      <div className="rounded-xl border border-emerald-500/30 bg-emerald-500/10 p-4">
        <div className="flex items-start gap-2">
          <CheckCircle2 className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" />
          <p className="text-sm">Snapshot verified and ready to apply.</p>
        </div>
      </div>
      <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-4">
        <div className="flex items-start gap-2">
          <ShieldAlert className="mt-0.5 h-4 w-4 shrink-0 text-destructive" />
          <p className="text-sm">
            Type{' '}
            <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{sourceName}</code>{' '}
            to confirm the overwrite.
          </p>
        </div>
      </div>
      <Input
        value={typed}
        onChange={(event) => onTypedChange(event.target.value)}
        placeholder={sourceName}
        className="font-mono"
        autoFocus
      />
      <div className="flex items-center justify-end gap-2">
        <Button variant="ghost" onClick={onCancel} disabled={busy}>
          Cancel
        </Button>
        <Button variant="destructive" onClick={onApply} disabled={busy || typed !== sourceName}>
          {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
          Apply restore
        </Button>
      </div>
    </div>
  );
}

export function RestoreApplyingStep({
  operation,
  cancelRequested,
  busy,
  onCancel,
}: {
  operation: RestoreOperation | null;
  cancelRequested: boolean;
  busy: boolean;
  onCancel: () => void;
}) {
  return (
    <div className="space-y-3">
      <div className="flex items-center gap-3 rounded-xl bg-muted/40 p-4 text-sm">
        <Loader2 className="h-4 w-4 animate-spin text-destructive" />
        <div className="flex-1">
          <p className="font-medium">Restoring data</p>
          <p className="text-muted-foreground text-xs">
            {operation?.phase || 'Applying snapshot contents to the target.'}
          </p>
        </div>
        {operation && (
          <span className="font-mono text-[11px] text-muted-foreground">{operation.status}</span>
        )}
      </div>
      {operation?.logs && (
        <pre className="max-h-40 overflow-auto whitespace-pre-wrap rounded-xl bg-muted/40 p-3 font-mono text-xs">
          {operation.logs}
        </pre>
      )}
      {cancelRequested || operation?.cancelRequested ? (
        <p className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-3 text-amber-600 text-xs">
          Cancellation requested — the restore stops at the next safe checkpoint.
        </p>
      ) : (
        <div className="flex items-center justify-end">
          <Button variant="outline" onClick={onCancel} disabled={busy}>
            Interrupt restore
          </Button>
        </div>
      )}
    </div>
  );
}

export function RestoreDoneStep({
  operation,
  sourceName,
  takenAt,
  onClose,
}: {
  operation: RestoreOperation | null;
  sourceName: string;
  takenAt: string;
  onClose: () => void;
}) {
  const tone = operation ? terminalTone(operation.status) : 'success';

  return (
    <div className="space-y-4">
      {tone === 'success' && (
        <div className="rounded-xl border border-emerald-500/30 bg-emerald-500/10 p-4">
          <div className="flex items-start gap-2">
            <CheckCircle2 className="mt-0.5 h-4 w-4 shrink-0 text-emerald-500" />
            <div className="text-sm">
              <p className="font-medium">Restore completed</p>
              <p className="mt-1 text-muted-foreground text-xs">
                {sourceName} now holds the snapshot data from {takenAt}.
              </p>
            </div>
          </div>
        </div>
      )}
      {tone === 'cancelled' && (
        <div className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-4">
          <div className="flex items-start gap-2">
            <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0 text-amber-500" />
            <div className="text-sm">
              <p className="font-medium">Restore cancelled</p>
              <p className="mt-1 text-muted-foreground text-xs">
                {operation?.error ||
                  'The restore was interrupted. The target may hold partial data.'}
              </p>
            </div>
          </div>
        </div>
      )}
      {tone === 'failed' && (
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-4">
          <div className="flex items-start gap-2">
            <XCircle className="mt-0.5 h-4 w-4 shrink-0 text-destructive" />
            <div className="text-sm">
              <p className="font-medium">Restore failed</p>
              <p className="mt-1 text-muted-foreground text-xs">
                {operation?.error || 'The restore did not complete.'}
              </p>
            </div>
          </div>
        </div>
      )}
      <div className="flex items-center justify-end">
        <Button onClick={onClose}>Close</Button>
      </div>
    </div>
  );
}
