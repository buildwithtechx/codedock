import {
  AlertTriangle,
  ArchiveRestore,
  CheckCircle2,
  Loader2,
  ShieldAlert,
  XCircle,
} from 'lucide-react';
import { useEffect, useState } from 'react';
import { Button } from '#/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogTitle } from '#/components/ui/dialog';
import { Input } from '#/components/ui/input';
import { cn } from '#/lib/utils';
import { formatBytes } from './destination-display';
import { useApplyRestore, useCancelRestore, usePrepareRestore, useRestoreOperation } from './hooks';
import type { BackupRecord, RestoreOperation } from './interfaces';

type WizardStep = 'review' | 'preparing' | 'prepared' | 'applying' | 'done';

const stepOrder: WizardStep[] = ['review', 'preparing', 'prepared', 'applying', 'done'];
const stepLabels: Record<WizardStep, string> = {
  review: 'Review',
  preparing: 'Prepare',
  prepared: 'Confirm',
  applying: 'Apply',
  done: 'Done',
};

function isTerminal(status: string): boolean {
  return ['COMPLETED', 'FAILED', 'ERROR', 'CANCELLED', 'CANCELED', 'EXPIRED'].includes(status);
}

function terminalTone(status: string): 'success' | 'failed' | 'cancelled' {
  if (status === 'COMPLETED') return 'success';
  if (status === 'CANCELLED' || status === 'CANCELED') return 'cancelled';
  return 'failed';
}

export function RestoreWizard({
  record,
  sourceName,
  onClose,
  onDone,
}: {
  record: BackupRecord;
  sourceName: string;
  onClose: () => void;
  onDone?: () => void;
}) {
  const [step, setStep] = useState<WizardStep>('review');
  const [operationId, setOperationId] = useState<string | null>(null);
  const [confirmation, setConfirmation] = useState<string | null>(null);
  const [typed, setTyped] = useState('');
  const [busy, setBusy] = useState(false);
  const [cancelRequested, setCancelRequested] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);

  const prepare = usePrepareRestore();
  const apply = useApplyRestore();
  const cancel = useCancelRestore();
  const operationQuery = useRestoreOperation(operationId, step === 'applying' || step === 'done');
  const operation: RestoreOperation | null = operationQuery.data?.data ?? null;

  useEffect(() => {
    if (!operation) return;
    if (step === 'applying' && isTerminal(operation.status)) setStep('done');
  }, [operation, step]);

  useEffect(() => {
    if (step === 'done' && operation && terminalTone(operation.status) === 'success') onDone?.();
  }, [step, operation, onDone]);

  const startPrepare = async () => {
    if (busy) return;
    setBusy(true);
    setActionError(null);
    setStep('preparing');
    try {
      const res = await prepare.mutateAsync({ recordId: record.id });
      setOperationId(res.data.operation.id);
      setConfirmation(res.data.confirmation);
      setStep(res.data.operation.status === 'REVIEWED' ? 'prepared' : 'preparing');
      if (res.data.operation.status !== 'REVIEWED') {
        await operationQuery.refetch();
        setStep('prepared');
      }
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Failed to prepare restore');
      setStep('review');
    } finally {
      setBusy(false);
    }
  };

  const applyRestore = async () => {
    if (busy || !operationId || !confirmation) return;
    if (!sourceName || typed !== sourceName) {
      setActionError(`Type "${sourceName}" to confirm`);
      return;
    }
    setBusy(true);
    setActionError(null);
    try {
      await apply.mutateAsync({ operationId, confirmation });
      setStep('applying');
      void operationQuery.refetch();
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Failed to apply restore');
    } finally {
      setBusy(false);
    }
  };

  const cancelRestore = async () => {
    if (busy) return;
    if (!operationId) {
      onClose();
      return;
    }
    setBusy(true);
    setActionError(null);
    try {
      await cancel.mutateAsync(operationId);
      setCancelRequested(true);
      await operationQuery.refetch();
      const status = operationQuery.data?.data?.status ?? '';
      if (step !== 'applying' || isTerminal(status)) {
        if (step !== 'applying') onClose();
        else setStep('done');
      }
    } catch (error) {
      setActionError(error instanceof Error ? error.message : 'Failed to cancel restore');
    } finally {
      setBusy(false);
    }
  };

  const takenAt = record.startedAt ? new Date(record.startedAt).toLocaleString() : 'Unknown date';
  const currentIndex = stepOrder.indexOf(step);

  return (
    <Dialog open onOpenChange={(openState) => !openState && !busy && onClose()}>
      <DialogContent className="sm:max-w-xl">
        <DialogTitle className="flex items-center gap-2">
          <ArchiveRestore className="h-5 w-5 text-primary" />
          Restore snapshot
        </DialogTitle>
        <DialogDescription>
          Restoring <span className="font-medium text-foreground">{sourceName}</span> from {takenAt}
          .
        </DialogDescription>

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
              {index < stepOrder.length - 1 && (
                <span className="h-px min-w-2 flex-1 bg-border/50" />
              )}
            </div>
          ))}
        </div>

        {actionError && (
          <p role="alert" className="rounded-lg bg-destructive/10 p-3 text-destructive text-sm">
            {actionError}
          </p>
        )}

        {step === 'review' && (
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
                    Current tables and stored data are replaced with the snapshot contents. The
                    restore is verified before anything is applied.
                  </p>
                </div>
              </div>
            </div>
            <div className="flex items-center justify-end gap-2">
              <Button variant="ghost" onClick={onClose} disabled={busy}>
                Cancel
              </Button>
              <Button onClick={startPrepare} disabled={busy}>
                {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
                Verify and continue
              </Button>
            </div>
          </div>
        )}

        {step === 'preparing' && (
          <div className="space-y-3">
            <div className="flex items-center gap-3 rounded-xl bg-muted/40 p-4 text-sm">
              <Loader2 className="h-4 w-4 animate-spin text-primary" />
              <div>
                <p className="font-medium">Verifying snapshot integrity</p>
                <p className="text-muted-foreground text-xs">
                  Checksums and restore targets are reviewed before anything changes.
                </p>
              </div>
            </div>
            <div className="flex items-center justify-end">
              <Button variant="ghost" onClick={cancelRestore} disabled={busy}>
                Cancel restore
              </Button>
            </div>
          </div>
        )}

        {step === 'prepared' && (
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
                  <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">
                    {sourceName}
                  </code>{' '}
                  to confirm the overwrite.
                </p>
              </div>
            </div>
            <Input
              value={typed}
              onChange={(event) => setTyped(event.target.value)}
              placeholder={sourceName}
              className="font-mono"
              autoFocus
            />
            <div className="flex items-center justify-end gap-2">
              <Button variant="ghost" onClick={cancelRestore} disabled={busy}>
                Cancel
              </Button>
              <Button
                variant="destructive"
                onClick={applyRestore}
                disabled={busy || typed !== sourceName}
              >
                {busy ? <Loader2 className="h-4 w-4 animate-spin" /> : null}
                Apply restore
              </Button>
            </div>
          </div>
        )}

        {step === 'applying' && (
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
                <span className="font-mono text-[11px] text-muted-foreground">
                  {operation.status}
                </span>
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
                <Button variant="outline" onClick={cancelRestore} disabled={busy}>
                  Interrupt restore
                </Button>
              </div>
            )}
          </div>
        )}

        {step === 'done' && (
          <div className="space-y-4">
            {(!operation || terminalTone(operation.status) === 'success') && (
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
            {operation && terminalTone(operation.status) === 'cancelled' && (
              <div className="rounded-xl border border-amber-500/30 bg-amber-500/10 p-4">
                <div className="flex items-start gap-2">
                  <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0 text-amber-500" />
                  <div className="text-sm">
                    <p className="font-medium">Restore cancelled</p>
                    <p className="mt-1 text-muted-foreground text-xs">
                      {operation.error ||
                        'The restore was interrupted. The target may hold partial data.'}
                    </p>
                  </div>
                </div>
              </div>
            )}
            {operation && terminalTone(operation.status) === 'failed' && (
              <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-4">
                <div className="flex items-start gap-2">
                  <XCircle className="mt-0.5 h-4 w-4 shrink-0 text-destructive" />
                  <div className="text-sm">
                    <p className="font-medium">Restore failed</p>
                    <p className="mt-1 text-muted-foreground text-xs">
                      {operation.error || 'The restore did not complete.'}
                    </p>
                  </div>
                </div>
              </div>
            )}
            <div className="flex items-center justify-end">
              <Button onClick={onClose}>Close</Button>
            </div>
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
